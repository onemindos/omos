import { Box, Link, Typography } from '@material-ui/core';
import { StatusChip } from '@internal/plugin-onemind-theme';
import { age, clock, NO_PILLAR, pillarTitle, UNASSIGNED } from './format';
import { Empty, Meter, OwnerLink, Panel, Row } from './Panel';
import type { Owner, OwnerIndex } from './owners';
import type {
  BackupItem,
  LiveAlert,
  LiveResponse,
  NodeItem,
} from './types';

type Chip = 'ok' | 'warn' | 'critical' | 'unknown' | 'info';
type Tone = 'critical' | 'warn' | 'ok' | 'none';
const ALERT_CHIP: Record<LiveAlert['severity'], Chip> = { critical: 'critical', warning: 'warn', info: 'info' };
const BACKUP_TONE: Record<BackupItem['attention'], Tone> = { review: 'critical', setup: 'warn', none: 'none' };

function nodeTone(node: NodeItem): Tone {
  if (!node.ready) return 'critical';
  return node.pressure.length ? 'warn' : 'ok';
}

function rowTone(status: Chip): Tone {
  if (status === 'ok' || status === 'critical') return status;
  return 'warn';
}

// Needs review first, then needs setup, then healthy, then retired.
function backupOrder(item: BackupItem) {
  if (item.attention === 'review') return 0;
  if (item.attention === 'setup') return 1;
  return item.state === 'retired' ? 3 : 2;
}
const external = { target: '_blank', rel: 'noopener noreferrer' } as const;

/** Groups by pillar, then by owning app; items with no owner land under Unassigned. */
export function groupByPillarAndApp<T>(
  items: T[],
  ownerOf: (item: T) => Owner | undefined,
): { pillar: string; apps: { owner?: Owner; items: T[] }[] }[] {
  const pillars = new Map<string, Map<string, { owner?: Owner; items: T[] }>>();
  for (const item of items) {
    const owner = ownerOf(item);
    const pillar = owner ? (owner.pillar ? pillarTitle(owner.pillar) : NO_PILLAR) : UNASSIGNED;
    const apps = pillars.get(pillar) ?? new Map();
    const key = owner?.name ?? UNASSIGNED;
    const group = apps.get(key) ?? { owner, items: [] };
    group.items.push(item);
    apps.set(key, group);
    pillars.set(pillar, apps);
  }
  const rank = (name: string) => [NO_PILLAR, UNASSIGNED].indexOf(name) + 1 || 0;
  return [...pillars]
    .sort(([a], [b]) => rank(a) - rank(b) || a.localeCompare(b))
    .map(([pillar, apps]) => ({ pillar, apps: [...apps.values()] }));
}

export function AlertsPanel({ data, index }: { data: LiveResponse; index: OwnerIndex }) {
  const alerts = data.alerts;
  const groups = alerts
    ? groupByPillarAndApp<LiveAlert>(alerts, alert => index.resolve({ namespace: alert.namespace }))
    : [];
  return (
    <Panel
      title="Firing alerts"
      state={data.panels.alerts}
      hint="From Alertmanager. Grouped by pillar and owning app (matched by Kubernetes namespace)."
      wide
    >
      {alerts && !alerts.length && <Empty>No alerts firing.</Empty>}
      {groups.map(group => (
        <Box key={group.pillar} mt={1}>
          <Typography variant="subtitle2" component="h3">
            {group.pillar}
          </Typography>
          {group.apps.map(app => (
            <Box key={app.owner?.name ?? UNASSIGNED} ml={1} mt={0.5}>
              <Typography variant="caption" color="textSecondary">
                App: <OwnerLink owner={app.owner} />
              </Typography>
              {app.items.map(alert => (
                <Row key={alert.fingerprint} tone={alert.severity === 'critical' ? 'critical' : 'warn'}>
                  <StatusChip status={ALERT_CHIP[alert.severity]} label={alert.severity} />
                  <Typography variant="body2" style={{ fontWeight: 600 }}>
                    {alert.name}
                  </Typography>
                  <Typography variant="body2" color="textSecondary" style={{ flex: '1 1 240px', minWidth: 0 }}>
                    {alert.summary}
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    {alert.namespace ? `${alert.namespace} · ` : ''}firing {age(alert.startsAt)}
                  </Typography>
                  {alert.runbookUrl && (
                    <Link href={alert.runbookUrl} variant="caption" {...external}>
                      Runbook
                    </Link>
                  )}
                  <Link href={alert.grafanaUrl} variant="caption" {...external}>
                    Grafana
                  </Link>
                </Row>
              ))}
            </Box>
          ))}
        </Box>
      ))}
    </Panel>
  );
}

export function ArgoPanel({ data, index }: { data: LiveResponse; index: OwnerIndex }) {
  const argo = data.argo;
  return (
    <Panel
      title="Argo CD"
      state={data.panels.argo}
      hint="Out of sync, degraded, progressing more than 10 minutes, or last operation failed."
    >
      {argo && !argo.items.length && <Empty>All {argo.total} applications are healthy and in sync.</Empty>}
      {argo?.items.map(item => (
        <Row key={item.name} tone={item.severity === 'critical' ? 'critical' : 'warn'}>
          <Link href={item.url} variant="body2" style={{ fontWeight: 600 }} {...external}>
            {item.name}
          </Link>
          {item.reasons.map(reason => (
            <StatusChip key={reason} status={item.severity === 'critical' ? 'critical' : 'warn'} label={reason} />
          ))}
          <Typography variant="caption" color="textSecondary">
            {item.health} / {item.sync}
            {item.since ? ` · since ${age(item.since)} ago` : ''}
          </Typography>
          <OwnerLink owner={index.byArgoApp(item.name)} />
        </Row>
      ))}
    </Panel>
  );
}

export function WorkloadsPanel({ data, index }: { data: LiveResponse; index: OwnerIndex }) {
  const workloads = data.workloads;
  const empty = workloads && !workloads.issues.length && !workloads.notReady.length;
  return (
    <Panel
      title="Workloads"
      state={data.panels.workloads}
      hint="Crash-looping, pending or not-ready pods, recent restarts and out-of-memory kills (last hour)."
    >
      {empty && <Empty>All {workloads.total} controllers are ready and no pod is failing.</Empty>}
      {workloads?.notReady.map(name => (
        <Row key={`c:${name}`} tone="warn">
          <Typography variant="body2" style={{ fontWeight: 600 }}>
            {name}
          </Typography>
          <StatusChip status="warn" label="Fewer ready replicas than desired" />
          <OwnerLink owner={index.byNamespace(name.split('/')[0])} />
        </Row>
      ))}
      {workloads?.issues.map(issue => (
        <Row key={`${issue.namespace}/${issue.pod}`} tone={issue.severity === 'critical' ? 'critical' : 'warn'}>
          <Typography variant="body2" style={{ fontWeight: 600 }}>
            {issue.namespace}/{issue.pod}
          </Typography>
          {issue.reasons.map(reason => (
            <StatusChip key={reason} status={reason === 'CrashLoopBackOff' ? 'critical' : 'warn'} label={reason} />
          ))}
          <Typography variant="caption" color="textSecondary">
            {issue.restarts} restarts{issue.lastRestartAt ? `, last ${age(issue.lastRestartAt)} ago` : ''}
          </Typography>
          <OwnerLink owner={index.resolve({ argoApp: issue.app, namespace: issue.namespace })} />
        </Row>
      ))}
    </Panel>
  );
}

export function NodesPanel({ data }: { data: LiveResponse }) {
  return (
    <Panel
      title="Nodes"
      state={data.panels.nodes}
      hint="Readiness and pressure from Kubernetes; use from Prometheus; requested is the sum of pod requests against allocatable."
    >
      {data.nodes?.map(node => (
        <Row key={node.name} tone={nodeTone(node)}>
          <Typography variant="body2" style={{ fontWeight: 600, flex: '1 1 160px', minWidth: 0 }}>
            {node.name}
          </Typography>
          <StatusChip status={node.ready ? 'ok' : 'critical'} label={node.ready ? 'Ready' : 'Not ready'} />
          {node.pressure.map(p => (
            <StatusChip key={p} status="warn" label={p} />
          ))}
          <Meter label="CPU" value={node.cpuPercent} />
          <Meter label="Memory" value={node.memoryPercent} />
          <Meter label="Disk" value={node.diskPercent} />
          <Meter label="CPU req" value={node.cpuRequestedPercent} warnAt={85} criticalAt={100} />
          <Meter label="Mem req" value={node.memoryRequestedPercent} warnAt={85} criticalAt={100} />
          <Typography variant="caption" color="textSecondary">
            {node.pods} pods
          </Typography>
        </Row>
      ))}
    </Panel>
  );
}

export function StoragePanel({ data, index }: { data: LiveResponse; index: OwnerIndex }) {
  const storage = data.storage;
  return (
    <Panel
      title="Storage"
      state={data.panels.storage}
      hint="Volumes over 80% full or projected to fill within four days. Snapshot growth arrives with the DigitalOcean view."
    >
      {storage && !storage.pvcs.length && <Empty>No volume is over 80% full.</Empty>}
      {storage?.pvcs.map(pvc => (
        <Row key={`${pvc.namespace}/${pvc.name}`} tone={pvc.severity === 'critical' ? 'critical' : 'warn'}>
          <Typography variant="body2" style={{ fontWeight: 600, flex: '1 1 180px', minWidth: 0 }}>
            {pvc.namespace}/{pvc.name}
          </Typography>
          <Meter label="Used" value={pvc.usedPercent} warnAt={80} criticalAt={95} />
          {pvc.fillsWithin4Days && <StatusChip status="warn" label="Fills within 4 days" />}
          <OwnerLink owner={index.byNamespace(pvc.namespace)} />
        </Row>
      ))}
    </Panel>
  );
}

const BACKUP_CHIP: Record<string, [Chip, string]> = {
  ok: ['ok', 'OK'],
  failing: ['critical', 'Last run failed'],
  overdue: ['warn', 'Overdue'],
  'never-ran': ['warn', 'Never succeeded'],
  suspended: ['warn', 'Suspended'],
};

function backupChip(item: BackupItem): [Chip, string] {
  if (item.state === 'retired') return ['info', 'Retired'];
  if (item.state === 'needs-setup') return ['warn', 'Needs setup'];
  return BACKUP_CHIP[item.health] ?? ['unknown', item.health];
}

export function BackupsPanel({ data }: { data: LiveResponse }) {
  const items = data.backups?.items ?? [];
  const order = backupOrder;
  const sorted = [...items].sort((a, b) => order(a) - order(b) || a.name.localeCompare(b.name));
  const volsync = data.backups?.volsync;
  return (
    <Panel
      title="Backups"
      state={data.panels.backups}
      hint="Retired jobs were replaced on purpose. Needs setup is waiting for something only a person can do. Neither counts as an attention item."
    >
      {data.backups && (
        <Typography variant="body2" color="textSecondary">
          {data.counts.backups} to review · {data.counts.backupsNeedSetup} needs setup · {data.counts.backupsRetired} retired
        </Typography>
      )}
      {sorted.map(item => {
        const [status, label] = backupChip(item);
        return (
          <Row key={item.name} tone={BACKUP_TONE[item.attention]}>
            <Typography variant="body2" style={{ fontWeight: 600 }}>
              {item.name}
            </Typography>
            <StatusChip status={status} label={label} />
            <Typography variant="caption" color="textSecondary" style={{ flex: '1 1 200px', minWidth: 0 }}>
              {item.state === 'retired' && `Replaced by ${item.supersededBy ?? 'another mechanism'}. `}
              {item.state === 'needs-setup' && `${item.pendingReason ?? 'Waiting for setup.'} `}
              {item.lastSuccess ? `Last success ${age(item.lastSuccess)} ago.` : 'No successful run recorded.'}
            </Typography>
          </Row>
        );
      })}
      {volsync && (
        <Row tone={volsync.outOfSync.length ? 'critical' : 'ok'}>
          <Typography variant="body2" style={{ fontWeight: 600 }}>
            VolSync volume backups
          </Typography>
          <StatusChip status={volsync.outOfSync.length ? 'critical' : 'ok'} label={volsync.outOfSync.length ? `${volsync.outOfSync.length} late` : 'All on schedule'} />
          <Typography variant="caption" color="textSecondary">
            {volsync.sources} sources
            {volsync.outOfSync.length ? `: ${volsync.outOfSync.map(s => s.namespace).join(', ')}` : ''}
          </Typography>
        </Row>
      )}
      <Empty>Velero and CloudNativePG backups are not read here yet: Velero exports no metrics, and the CloudNativePG metrics carry no cluster name.</Empty>
    </Panel>
  );
}

const UPTIME_CHIP: Record<string, [Chip, string]> = {
  up: ['ok', 'Up'],
  down: ['critical', 'Down'],
  unknown: ['unknown', 'Unknown'],
  'authentication-required': ['info', 'Needs sign-in'],
};

export function UptimePanel({ data }: { data: LiveResponse }) {
  return (
    <Panel
      title="Uptime"
      state={data.panels.uptime}
      hint="Cluster-side endpoint checks. The DigitalOcean uptime check and Gatus checks are added by later work orders."
    >
      {data.uptime?.map(check => {
        const [status, label] = UPTIME_CHIP[check.state] ?? ['unknown', check.state];
        return (
          <Row key={check.name} tone={rowTone(status)}>
            <Link href={check.link} variant="body2" style={{ fontWeight: 600 }} {...external}>
              {check.name}
            </Link>
            <StatusChip status={status} label={label} />
            <Typography variant="caption" color="textSecondary">
              {check.latencyMs !== undefined ? `${check.latencyMs} ms · ` : ''}
              {check.scope}
            </Typography>
          </Row>
        );
      })}
    </Panel>
  );
}

export function EventsPanel({ data, index }: { data: LiveResponse; index: OwnerIndex }) {
  const events = data.events;
  return (
    <Panel
      title="Warning events"
      state={data.panels.events}
      hint="Kubernetes Warning events from the last 30 minutes, repeats merged."
      wide
    >
      {events && !events.items.length && <Empty>No warning events in the last {events.windowMinutes} minutes.</Empty>}
      {events?.items.map(event => (
        <Row key={`${event.namespace}|${event.object}|${event.reason}|${event.message}`} tone="warn">
          <Typography variant="caption" color="textSecondary" style={{ minWidth: 64 }}>
            {clock(event.lastSeen)}
          </Typography>
          <Typography variant="body2" style={{ fontWeight: 600 }}>
            {event.namespace ? `${event.namespace}/` : ''}
            {event.object}
          </Typography>
          <StatusChip status="warn" label={event.reason} />
          <Typography variant="body2" color="textSecondary" style={{ flex: '1 1 260px', minWidth: 0 }}>
            {event.message}
          </Typography>
          {event.count > 1 && (
            <Typography variant="caption" color="textSecondary">
              ×{event.count}
            </Typography>
          )}
          <OwnerLink owner={index.byNamespace(event.namespace)} />
        </Row>
      ))}
    </Panel>
  );
}

import { DigitalOceanClient, DoApiError } from './client';
import type {
  AccountInfo,
  AlertPolicy,
  Billing,
  Cluster,
  Domain,
  Droplet,
  Firewall,
  Invoice,
  LoadBalancer,
  Section,
  Snapshot,
  SnapshotStats,
  Snapshots,
  UptimeCheck,
  Volume,
  Vpc,
} from './types';

/* eslint-disable @typescript-eslint/no-explicit-any */
type Raw = any;

const DAY = 86_400_000;
export const OLD_SNAPSHOT_DAYS = 90;

function failure(error: unknown): { state: 'unavailable'; reason: string } {
  if (error instanceof DoApiError) {
    if (error.status === 401 || error.status === 403)
      return {
        state: 'unavailable',
        reason: 'The token does not have the read scope for this section',
      };
    if (error.status === 404)
      return {
        state: 'unavailable',
        reason: 'DigitalOcean does not offer this section for the account',
      };
    if (error.status === 429)
      return { state: 'unavailable', reason: 'DigitalOcean rate limit' };
  }
  return { state: 'unavailable', reason: 'DigitalOcean could not be reached' };
}

/** Each section fails on its own: a missing scope never takes the others down. */
async function section<T>(read: () => Promise<T[]>): Promise<Section<T>> {
  try {
    return { state: 'ok', items: await read() };
  } catch (error) {
    return { ...failure(error), items: [] };
  }
}

const str = (v: unknown, fallback = '') => {
  if (typeof v === 'string') return v;
  return v === null || v === undefined ? fallback : String(v);
};
const num = (v: unknown, fallback = 0) =>
  typeof v === 'number' && Number.isFinite(v) ? v : fallback;
const strings = (v: unknown): string[] =>
  Array.isArray(v) ? v.map(x => String(x)) : [];
const ageDays = (iso: unknown, now: number) => {
  const t = typeof iso === 'string' ? Date.parse(iso) : NaN;
  return Number.isFinite(t) ? Math.max(0, Math.floor((now - t) / DAY)) : null;
};
const money = (v: unknown): string | undefined =>
  typeof v === 'string' && /^-?\d{1,12}(\.\d{1,6})?$/.test(v) ? v : undefined;

/** usage so far, scaled to the whole month (UTC), as dollars with two decimals. */
export function projectMonth(usage: string, generatedAt: number) {
  const at = new Date(generatedAt);
  const start = Date.UTC(at.getUTCFullYear(), at.getUTCMonth(), 1);
  const daysInMonth = new Date(
    Date.UTC(at.getUTCFullYear(), at.getUTCMonth() + 1, 0),
  ).getUTCDate();
  // At least one day, so the first hours of a month do not multiply noise by 30.
  const daysElapsed = Math.min(
    daysInMonth,
    Math.max(1, (generatedAt - start) / DAY),
  );
  const projected = (Number(usage) * daysInMonth) / daysElapsed;
  return {
    projectedMonth: projected.toFixed(2),
    daysElapsed: Math.round(daysElapsed * 100) / 100,
    daysInMonth,
  };
}

async function readAccount(c: DigitalOceanClient): Promise<AccountInfo> {
  try {
    const { account } = await c.get<Raw>('/v2/account');
    return {
      state: 'ok',
      status: str(account?.status),
      teamName: account?.team?.name ? str(account.team.name) : undefined,
      dropletLimit:
        typeof account?.droplet_limit === 'number'
          ? account.droplet_limit
          : undefined,
      volumeLimit:
        typeof account?.volume_limit === 'number'
          ? account.volume_limit
          : undefined,
    };
  } catch (error) {
    return failure(error);
  }
}

async function readBilling(
  c: DigitalOceanClient,
  now: number,
): Promise<Billing> {
  const invoices: Invoice[] = [];
  let invoicePreview: Invoice | undefined;
  try {
    const body = await c.get<Raw>('/v2/customers/my/invoices?per_page=6');
    for (const i of body.invoices ?? [])
      invoices.push({ period: str(i.invoice_period), amount: str(i.amount) });
    const p = body.invoice_preview;
    if (p?.invoice_period)
      invoicePreview = {
        period: str(p.invoice_period),
        amount: str(p.amount),
      };
  } catch {
    // Invoices are extra; the balance below decides the section state.
  }
  try {
    const b = await c.get<Raw>('/v2/customers/my/balance');
    const usage = money(b.month_to_date_usage);
    const stamp = Date.parse(str(b.generated_at));
    const at = Number.isFinite(stamp) ? stamp : now;
    return {
      state: 'ok',
      usageBeforeCredits: usage,
      amountDue: money(b.month_to_date_balance),
      accountBalance: money(b.account_balance),
      generatedAt: Number.isFinite(stamp)
        ? new Date(stamp).toISOString()
        : undefined,
      ...(usage ? projectMonth(usage, at) : {}),
      invoices,
      invoicePreview,
    };
  } catch (error) {
    return { ...failure(error), invoices };
  }
}

async function readClusters(c: DigitalOceanClient): Promise<Cluster[]> {
  const raw = await c.list<Raw>(
    '/v2/kubernetes/clusters',
    'kubernetes_clusters',
  );
  const out: Cluster[] = [];
  for (const k of raw) {
    let upgradeAvailable: string | null = null;
    let upgradeCheck: Cluster['upgradeCheck'] = 'ok';
    try {
      const up = await c.get<Raw>(
        `/v2/kubernetes/clusters/${encodeURIComponent(k.id)}/upgrades`,
      );
      const versions = ((up.available_upgrade_versions ?? []) as Raw[]).map(v =>
        str(v.slug),
      );
      upgradeAvailable = versions.length ? versions[versions.length - 1] : null;
    } catch {
      upgradeCheck = 'unavailable';
    }
    out.push({
      id: str(k.id),
      name: str(k.name),
      region: str(k.region),
      version: str(k.version),
      status: k.status?.state ? str(k.status.state) : undefined,
      highAvailability: Boolean(k.ha),
      upgradeAvailable,
      upgradeCheck,
      createdAt: k.created_at ? str(k.created_at) : undefined,
      nodePools: ((k.node_pools ?? []) as Raw[]).map(p => ({
        id: str(p.id),
        name: str(p.name),
        size: str(p.size),
        count: num(p.count),
        autoScale: Boolean(p.auto_scale),
        minNodes: p.auto_scale ? num(p.min_nodes) : undefined,
        maxNodes: p.auto_scale ? num(p.max_nodes) : undefined,
      })),
    });
  }
  return out;
}

const droplets = (raw: Raw[]): Droplet[] =>
  raw.map(d => ({
    id: str(d.id),
    name: str(d.name),
    status: str(d.status),
    size: str(d.size_slug ?? d.size?.slug),
    region: str(d.region?.slug),
    vcpus: typeof d.vcpus === 'number' ? d.vcpus : undefined,
    memoryMb: typeof d.memory === 'number' ? d.memory : undefined,
    diskGb: typeof d.disk === 'number' ? d.disk : undefined,
    tags: strings(d.tags),
    createdAt: d.created_at ? str(d.created_at) : undefined,
  }));

const volumes = (raw: Raw[]): Volume[] =>
  raw.map(v => {
    const attached = strings(v.droplet_ids);
    return {
      id: str(v.id),
      name: str(v.name),
      sizeGb: num(v.size_gigabytes),
      region: str(v.region?.slug),
      attachedDropletIds: attached,
      unattached: attached.length === 0,
      createdAt: v.created_at ? str(v.created_at) : undefined,
    };
  });

const snapshots = (raw: Raw[], now: number): Snapshot[] =>
  raw.map(s => ({
    id: str(s.id),
    name: str(s.name),
    resourceType: str(s.resource_type),
    resourceId: s.resource_id ? str(s.resource_id) : undefined,
    sizeGb: num(s.size_gigabytes),
    minDiskGb:
      typeof s.min_disk_size === 'number' ? s.min_disk_size : undefined,
    regions: strings(s.regions),
    createdAt: s.created_at ? str(s.created_at) : undefined,
    ageDays: ageDays(s.created_at, now),
  }));

export function snapshotStats(items: Snapshot[]): SnapshotStats {
  const dated = items.filter(s => s.createdAt);
  const oldest =
    dated.sort(
      (a, b) => Date.parse(a.createdAt!) - Date.parse(b.createdAt!),
    )[0] ?? null;
  const largest = [...items].sort((a, b) => b.sizeGb - a.sizeGb)[0] ?? null;
  return {
    count: items.length,
    totalSizeGb:
      Math.round(items.reduce((n, s) => n + s.sizeGb, 0) * 100) / 100,
    oldest,
    largest,
    olderThanDays: OLD_SNAPSHOT_DAYS,
    olderCount: items.filter(
      s => s.ageDays !== null && s.ageDays > OLD_SNAPSHOT_DAYS,
    ).length,
  };
}

const loadBalancers = (raw: Raw[]): LoadBalancer[] =>
  raw.map(l => {
    const dropletCount = strings(l.droplet_ids).length;
    return {
      id: str(l.id),
      name: str(l.name),
      region: str(l.region?.slug),
      size: l.size_slug ?? l.size ? str(l.size_slug ?? l.size) : undefined,
      type: l.type ? str(l.type) : undefined,
      status: str(l.status),
      forwardingRules: Array.isArray(l.forwarding_rules)
        ? l.forwarding_rules.length
        : 0,
      dropletCount,
      idle: dropletCount === 0 && !l.tag,
      createdAt: l.created_at ? str(l.created_at) : undefined,
    };
  });

const firewalls = (raw: Raw[]): Firewall[] =>
  raw.map(f => ({
    id: str(f.id),
    name: str(f.name),
    status: str(f.status),
    inboundRules: Array.isArray(f.inbound_rules) ? f.inbound_rules.length : 0,
    outboundRules: Array.isArray(f.outbound_rules)
      ? f.outbound_rules.length
      : 0,
    dropletCount: strings(f.droplet_ids).length,
    tags: strings(f.tags),
  }));

const vpcs = (raw: Raw[]): Vpc[] =>
  raw.map(v => ({
    id: str(v.id),
    name: str(v.name),
    region: str(v.region),
    ipRange: str(v.ip_range),
    isDefault: Boolean(v.default),
  }));

async function readDomains(c: DigitalOceanClient): Promise<Domain[]> {
  const raw = await c.list<Raw>('/v2/domains', 'domains');
  const out: Domain[] = [];
  for (const d of raw) {
    let recordCount: number | null = null;
    try {
      recordCount = (
        await c.list<Raw>(
          `/v2/domains/${encodeURIComponent(d.name)}/records`,
          'domain_records',
        )
      ).length;
    } catch {
      recordCount = null;
    }
    out.push({ name: str(d.name), recordCount });
  }
  return out;
}

export function overallState(statuses: string[]): UptimeCheck['state'] {
  const up = statuses.map(s => s.toUpperCase());
  if (up.some(s => s === 'DOWN' || s === 'DEGRADED')) return 'down';
  if (up.length && up.every(s => s === 'UP')) return 'up';
  return 'unknown';
}

async function readUptime(c: DigitalOceanClient): Promise<UptimeCheck[]> {
  const raw = await c.list<Raw>('/v2/uptime/checks', 'checks');
  const out: UptimeCheck[] = [];
  for (const u of raw) {
    let regionStates: Record<string, string> = {};
    try {
      const body = await c.get<Raw>(
        `/v2/uptime/checks/${encodeURIComponent(u.id)}/state`,
      );
      for (const [region, r] of Object.entries<Raw>(body.state?.regions ?? {}))
        regionStates[region] = str(r?.status);
    } catch {
      regionStates = {};
    }
    out.push({
      id: str(u.id),
      name: str(u.name),
      type: str(u.type),
      target: str(u.target),
      enabled: u.enabled !== false,
      regions: strings(u.regions),
      state: overallState(Object.values(regionStates)),
      regionStates,
    });
  }
  return out;
}

const alerts = (raw: Raw[]): AlertPolicy[] =>
  raw.map(a => ({
    id: str(a.uuid ?? a.id),
    description: str(a.description),
    type: str(a.type),
    enabled: a.enabled !== false,
    compare: a.compare ? str(a.compare) : undefined,
    value: typeof a.value === 'number' ? a.value : undefined,
    window: a.window ? str(a.window) : undefined,
  }));

/** Reads every section. Never throws: a failed section carries its own state. */
export async function collect(
  client: DigitalOceanClient,
  now: number,
): Promise<Snapshots> {
  const [
    account,
    billing,
    kubernetes,
    dr,
    vol,
    snap,
    lb,
    fw,
    vpc,
    dom,
    up,
    al,
  ] = await Promise.all([
    readAccount(client),
    readBilling(client, now),
    section(() => readClusters(client)),
    section(async () =>
      droplets(await client.list<Raw>('/v2/droplets', 'droplets')),
    ),
    section(async () =>
      volumes(await client.list<Raw>('/v2/volumes', 'volumes')),
    ),
    section(async () =>
      snapshots(await client.list<Raw>('/v2/snapshots', 'snapshots'), now),
    ),
    section(async () =>
      loadBalancers(
        await client.list<Raw>('/v2/load_balancers', 'load_balancers'),
      ),
    ),
    section(async () =>
      firewalls(await client.list<Raw>('/v2/firewalls', 'firewalls')),
    ),
    section(async () => vpcs(await client.list<Raw>('/v2/vpcs', 'vpcs'))),
    section(() => readDomains(client)),
    section(() => readUptime(client)),
    section(async () =>
      alerts(await client.list<Raw>('/v2/monitoring/alerts', 'policies')),
    ),
  ]);
  return {
    generatedAt: new Date(now).toISOString(),
    account,
    billing,
    kubernetes,
    droplets: dr,
    volumes: {
      ...vol,
      unattachedCount: vol.items.filter(v => v.unattached).length,
      unattachedSizeGb: vol.items
        .filter(v => v.unattached)
        .reduce((n, v) => n + v.sizeGb, 0),
    },
    snapshots: { ...snap, stats: snapshotStats(snap.items) },
    loadBalancers: lb,
    firewalls: fw,
    vpcs: vpc,
    domains: dom,
    uptime: up,
    alerts: al,
    spaces: {
      state: 'not-available',
      reason:
        'The DigitalOcean API cannot list Spaces buckets with an API token; open the Spaces page in the DigitalOcean console.',
      items: [],
    },
  };
}

const unavailable = <T>(reason: string): Section<T> => ({
  state: 'unavailable',
  reason,
  items: [],
});

/** The answer while no token is configured: every section empty and flagged, nothing crashes. */
export function notConnected(now: number): Snapshots {
  const reason = 'Not connected: needs the read-only DigitalOcean token';
  return {
    generatedAt: new Date(now).toISOString(),
    account: { state: 'unavailable', reason },
    billing: { state: 'unavailable', reason, invoices: [] },
    kubernetes: unavailable(reason),
    droplets: unavailable(reason),
    volumes: {
      ...unavailable<Volume>(reason),
      unattachedCount: 0,
      unattachedSizeGb: 0,
    },
    snapshots: {
      ...unavailable<Snapshot>(reason),
      stats: snapshotStats([]),
    },
    loadBalancers: unavailable(reason),
    firewalls: unavailable(reason),
    vpcs: unavailable(reason),
    domains: unavailable(reason),
    uptime: unavailable(reason),
    alerts: unavailable(reason),
    spaces: {
      state: 'not-available',
      reason:
        'The DigitalOcean API cannot list Spaces buckets with an API token.',
      items: [],
    },
  };
}

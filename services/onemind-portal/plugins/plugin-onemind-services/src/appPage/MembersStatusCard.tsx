import { useEffect, useMemo, useState } from 'react';
import { InfoCard, Link } from '@backstage/core-components';
import { useApi } from '@backstage/core-plugin-api';
import { catalogApiRef, humanizeEntityRef, useEntity, useEntityRefLink } from '@backstage/plugin-catalog-react';
import { Box, List, ListItem, Typography } from '@material-ui/core';
import {
  AppStatus, EntityLike, InfoTooltip, SourceLine, StatusDot, statusRank, summarize, useAppStatuses,
} from '@internal/plugin-onemind-common';

const FIELDS = ['kind', 'metadata.name', 'metadata.namespace', 'metadata.title', 'metadata.labels', 'metadata.annotations', 'spec.system', 'relations'];
const ORDER: AppStatus[] = ['down', 'warn', 'unknown', 'ok'];
const LABEL: Record<AppStatus, string> = { down: 'Down', warn: 'Degraded', unknown: 'Unknown', ok: 'Healthy' };
const MAX_LISTED = 8;

/** Member apps of a System or Domain, found through the catalog; only Components with the `app` role. */
function useMemberApps(entity: ReturnType<typeof useEntity>['entity']) {
  const catalog = useApi(catalogApiRef);
  const [apps, setApps] = useState<EntityLike[]>();
  const [failed, setFailed] = useState(false);
  const key = `${entity.kind}:${entity.metadata.namespace}/${entity.metadata.name}`;
  useEffect(() => {
    let live = true;
    setApps(undefined);
    setFailed(false);
    (async () => {
      // A System summarises its own apps; a Domain summarises the Systems in it and in its sub-Domains.
      let systems = entity.kind.toLowerCase() === 'system' ? [entity.metadata.name] : [];
      if (entity.kind.toLowerCase() === 'domain') {
        const self = `domain:${entity.metadata.namespace ?? 'default'}/${entity.metadata.name}`.toLowerCase();
        const domainRefs = [self, ...(entity.relations ?? []).filter(r => r.type === 'hasPart' && r.targetRef.startsWith('domain:')).map(r => r.targetRef)];
        const found = await catalog.getEntities({ filter: { kind: 'system', 'spec.domain': domainRefs }, fields: ['metadata.name'] });
        systems = found.items.map(s => s.metadata.name);
      }
      if (!systems.length) return [];
      const result = await catalog.getEntities({
        filter: { kind: 'component', 'spec.system': systems, 'metadata.labels.onemind.io/catalog-role': 'app' },
        fields: FIELDS,
      });
      return result.items as unknown as EntityLike[];
    })()
      .then(items => live && setApps(items))
      .catch(() => live && setFailed(true));
    return () => { live = false; };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [catalog, key]);
  return { apps, failed };
}

function Row({ entity, status, detail }: { entity: EntityLike; status: AppStatus; detail: string }) {
  const link = useEntityRefLink();
  const ref = { kind: entity.kind, namespace: entity.metadata.namespace ?? 'default', name: entity.metadata.name };
  const title = (entity.metadata as { title?: string }).title ?? humanizeEntityRef(ref, { defaultKind: 'component' });
  return (
    <ListItem disableGutters style={{ gap: 8 }}>
      <StatusDot status={status} title={detail} />
      <Link to={link(ref)}>{title}</Link>
    </ListItem>
  );
}

/** System and Domain summary: the same status model as the app page, for every member app. */
export function MembersStatusCard() {
  const { entity } = useEntity();
  const { apps, failed } = useMemberApps(entity);
  const list = useMemo(() => apps ?? [], [apps]);
  const { results, loading } = useAppStatuses(list);
  const counts = summarize(results.map(r => r.result));
  const sorted = [...results].sort((a, b) => statusRank(b.result.status) - statusRank(a.result.status));
  const problems = sorted.filter(r => r.result.status !== 'ok');
  const shown = problems.slice(0, MAX_LISTED);
  const observed = results.map(r => r.result.observedAt).filter((x): x is string => !!x).sort()[0];
  return (
    <InfoCard
      title="App status"
      subheader={<SourceLine source="Same signals as each app page" observedAt={observed} />}
      action={
        <InfoTooltip label="app status">
          One status per member app, from Argo CD, replica readiness, uptime checks and scoped alerts. Grey means the
          app has no connected signal or a signal could not be read; it is not a healthy count.
        </InfoTooltip>
      }
    >
      {failed && <Typography role="alert" color="error" variant="body2">Member apps could not be loaded.</Typography>}
      {!failed && !apps && <Typography role="status" variant="body2">Loading member apps…</Typography>}
      {apps && !apps.length && <Typography variant="body2">No member apps are recorded here.</Typography>}
      {apps && apps.length > 0 && (
        <>
          <Box display="flex" flexWrap="wrap" style={{ gap: 16 }} mb={1}>
            {ORDER.map(s => (
              <Box key={s} display="flex" alignItems="center" style={{ gap: 6 }} data-testid={`count-${s}`}>
                <StatusDot status={s} />
                <Typography variant="body2"><b>{loading ? '…' : counts[s]}</b> {LABEL[s]}</Typography>
              </Box>
            ))}
          </Box>
          <List dense disablePadding>
            {shown.map(r => (
              <Row key={r.entity.metadata.name} entity={r.entity} status={r.result.status} detail={r.result.signals.map(s => s.detail).join(' ')} />
            ))}
          </List>
          {problems.length > shown.length && (
            <Typography variant="caption">{problems.length - shown.length} more apps need a look.</Typography>
          )}
          {!problems.length && !loading && <Typography variant="body2">Every member app reports healthy.</Typography>}
        </>
      )}
    </InfoCard>
  );
}

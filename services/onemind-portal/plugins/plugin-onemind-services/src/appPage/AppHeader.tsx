import { ReactNode, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Header, ButtonIcon, ButtonLink, Menu, MenuItem, MenuTrigger, HeaderMetadataUsers } from '@backstage/ui';
import type { HeaderProps } from '@backstage/ui';
import { RELATION_OWNED_BY } from '@backstage/catalog-model';
import { useRouteRefParams } from '@backstage/core-plugin-api';
import {
  entityRouteRef,
  getEntityRelations,
  useAsyncEntity,
  useEntityPresentation,
  useEntityRefLink,
  useStarredEntity,
} from '@backstage/plugin-catalog-react';
import StarIcon from '@material-ui/icons/Star';
import StarBorderIcon from '@material-ui/icons/StarBorder';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import { hostingBadge, StatusDot, useAppStatus, EntityLike } from '@internal/plugin-onemind-common';
import { buildAppLinks } from './links';

const PILLARS: Record<string, string> = { so: 'SO', le: 'LE', hp: 'HP', ge: 'GE', shared: 'Shared' };

/** Status dot with its text, read from the one shared hook. */
function HeaderStatus({ entity }: { entity: EntityLike }) {
  const { status, signals, loading } = useAppStatus(entity);
  const detail = loading ? 'Reading live signals' : signals.map(s => `${s.name}: ${s.detail}`).join(' ');
  return <StatusDot status={status} showLabel title={detail || 'No live signal is connected for this app.'} />;
}

function StarButton({ entity }: { entity: Parameters<typeof useStarredEntity>[0] }) {
  const { isStarredEntity, toggleStarredEntity } = useStarredEntity(entity);
  return (
    <ButtonIcon
      variant="secondary"
      aria-label={isStarredEntity ? 'Remove from favorites' : 'Add to favorites'}
      icon={isStarredEntity ? <StarIcon fontSize="small" /> : <StarBorderIcon fontSize="small" />}
      onPress={() => toggleStarredEntity()}
    />
  );
}

function MoreMenu() {
  const [, setSearchParams] = useSearchParams();
  const inspect = useCallback(
    () => setSearchParams(params => { const next = new URLSearchParams(params); next.set('inspect', 'overview'); return next; }),
    [setSearchParams],
  );
  return (
    <MenuTrigger>
      <ButtonIcon variant="secondary" aria-label="More actions" icon={<MoreVertIcon fontSize="small" />} />
      <Menu placement="bottom end">
        <MenuItem onAction={inspect}>Inspect entity</MenuItem>
        <MenuItem onAction={() => void navigator.clipboard?.writeText(window.location.toString())}>Copy link</MenuItem>
      </Menu>
    </MenuTrigger>
  );
}

function Owners({ entity }: { entity: Parameters<typeof getEntityRelations>[0] }) {
  const link = useEntityRefLink();
  const owners = getEntityRelations(entity, RELATION_OWNED_BY);
  return owners.length ? (
    <HeaderMetadataUsers users={owners.map(ref => ({ name: ref.name, href: link(ref) }))} />
  ) : null;
}

/**
 * The one header for every app, System and Domain: name, pillar and hosting badges, live status dot,
 * owner and a primary "Open app" button. The tabs are the same navigation the catalog renders.
 */
export function AppHeader({ tabs, activeTabId }: Pick<HeaderProps, 'tabs' | 'activeTabId'>) {
  const params = useRouteRefParams(entityRouteRef);
  const { entity } = useAsyncEntity();
  const presentation = useEntityPresentation(entity ?? params);
  if (!entity) return <Header title={presentation.primaryTitle} />;

  const labels = entity.metadata.labels ?? {};
  const pillar = PILLARS[labels['onemind.io/pillar'] ?? ''];
  const hosting = entity.kind.toLowerCase() === 'component' ? hostingBadge(entity) : undefined;
  const open = buildAppLinks(entity).openApp;
  const metadata: { label: string; value: ReactNode }[] = [];
  if (entity.kind.toLowerCase() !== 'domain') metadata.push({ label: 'Status', value: <HeaderStatus entity={entity as EntityLike} /> });
  const owner = <Owners entity={entity} />;
  if (getEntityRelations(entity, RELATION_OWNED_BY).length) metadata.push({ label: 'Owner', value: owner });
  const lifecycle = entity.spec?.lifecycle?.toString();
  if (lifecycle) metadata.push({ label: 'Lifecycle', value: lifecycle });

  return (
    <Header
      title={presentation.primaryTitle}
      tags={[
        { label: entity.kind },
        ...(pillar ? [{ label: pillar }] : []),
        ...(hosting ? [{ label: hosting.badge }] : []),
      ]}
      metadata={metadata}
      tabs={tabs}
      activeTabId={activeTabId}
      customActions={
        <>
          {open && <ButtonLink variant="primary" href={open}>Open app</ButtonLink>}
          <StarButton entity={entity} />
          <MoreMenu />
        </>
      }
    />
  );
}

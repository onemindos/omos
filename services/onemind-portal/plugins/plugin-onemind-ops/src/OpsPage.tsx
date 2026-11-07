import { makeStyles } from '@material-ui/core/styles';
import { Progress, WarningPanel } from '@backstage/core-components';
import { ArgoStrip, ClusterGrid, WarningsFeed } from './ClusterPanels';
import { ConnectionBar } from './ConnectionBar';
import { EventStream } from './EventStream';
import { GlobePanel } from './GlobePanel';
import { useNamespaceLinks } from './namespaceLinks';
import { cssVar } from './ui';
import { useLiveFeed } from './useLiveFeed';

const useStyles = makeStyles({
  page: {
    display: 'flex',
    flexDirection: 'column',
    gap: 8,
    padding: 12,
    minHeight: 'calc(100vh - 64px)',
    background: cssVar('surface-0'),
    color: cssVar('text-1'),
  },
  head: { display: 'flex', alignItems: 'baseline', gap: 12 },
  title: { margin: 0, fontSize: 20, fontWeight: 700 },
  sub: { fontSize: 12, color: cssVar('text-2') },
  main: {
    display: 'grid',
    gridTemplateColumns: 'minmax(0, 3fr) minmax(300px, 2fr)',
    gap: 8,
    height: 'clamp(460px, 62vh, 900px)',
    '@media (max-width: 1000px)': {
      gridTemplateColumns: '1fr',
      height: 'auto',
      '& > *': { minHeight: 420 },
    },
  },
  bottom: {
    display: 'grid',
    gridTemplateColumns: 'minmax(0, 3fr) minmax(300px, 2fr)',
    gap: 8,
    height: 280,
    '@media (max-width: 1000px)': {
      gridTemplateColumns: '1fr',
      height: 'auto',
      '& > *': { minHeight: 240 },
    },
  },
  right: {
    display: 'grid',
    gridTemplateRows: 'minmax(0, 1fr) auto',
    gap: 8,
    minHeight: 0,
  },
});

export function OpsPage() {
  const c = useStyles();
  const feed = useLiveFeed();
  const links = useNamespaceLinks();
  const { state, load, now } = feed;

  const natsUnavailable =
    state.status?.nats.state === 'down' ||
    state.status?.nats.state === 'disabled'
      ? `Bus events unavailable: NATS is ${state.status.nats.state}${
          state.status.nats.detail ? ` (${state.status.nats.detail})` : ''
        }.`
      : undefined;

  return (
    <div className={c.page}>
      <div className={c.head}>
        <h1 className={c.title}>Ops Center</h1>
        <span className={c.sub}>
          Live tracks, bus traffic and cluster health. Updates arrive over
          Signals.
        </span>
      </div>

      <ConnectionBar
        state={state}
        freshness={feed.freshness}
        now={now}
        signalsConnected={feed.signalsConnected}
      />

      {load.kind === 'not-deployed' && (
        <WarningPanel
          severity="error"
          title="The live bridge is not available"
          message="GET /api/livebridge/snapshot answered 404, so the livebridge backend plugin is not deployed in this portal. Nothing below is live; there is no data to show."
        />
      )}
      {load.kind === 'error' && (
        <WarningPanel
          severity="error"
          title="Could not load the live snapshot"
          message={`${load.message}. ${
            state.loaded
              ? 'Showing the last snapshot and deltas received.'
              : 'No data has been received yet.'
          }`}
        />
      )}
      {load.kind === 'loading' && !state.loaded && <Progress />}

      {load.kind !== 'not-deployed' && (
        <>
          <div className={c.main}>
            <GlobePanel
              tracks={state.tracks}
              now={now}
              nats={state.status?.nats}
            />
            <div className={c.right}>
              <EventStream
                events={state.events}
                dropped={state.droppedEvents}
                buckets={state.buckets}
                now={now}
                unavailable={natsUnavailable}
              />
            </div>
          </div>
          <div className={c.bottom}>
            <ClusterGrid
              cells={state.cluster}
              changedMs={state.clusterChangedMs}
              now={now}
              links={links}
              source={state.status?.kubernetes}
            />
            <div
              className={c.right}
              style={{ gridTemplateRows: 'minmax(0, 1fr) minmax(0, 1fr)' }}
            >
              <ArgoStrip
                apps={state.argo}
                source={state.status?.argo}
                now={now}
              />
              <WarningsFeed
                warnings={state.warnings}
                now={now}
                source={state.status?.kubernetes}
              />
            </div>
          </div>
        </>
      )}
    </div>
  );
}

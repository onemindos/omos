import { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  Link,
  TextField,
  Typography,
} from '@material-ui/core';
import {
  configApiRef,
  storageApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import {
  Dashboard,
  DashboardViewer,
  grafanaApiRef,
} from '@backstage-community/plugin-grafana';
const origin = 'https://grafana.tail717ba0.ts.net';
function safeDashboard(value: string) {
  try {
    const url = new URL(value);
    return url.origin === origin && url.pathname.startsWith('/d/')
      ? url
      : undefined;
  } catch {
    return undefined;
  }
}
export function GrafanaHome() {
  const embeddingEnabled =
    useApi(configApiRef).getOptionalBoolean(
      'onemind.grafanaEmbedding.enabled',
    ) === true;
  const api = useApi(grafanaApiRef),
    storage = useApi(storageApiRef),
    bucket = useMemo(() => storage.forBucket('onemind-home'), [storage]);
  const [snapshot, setSnapshot] = useState(() =>
    bucket.snapshot<string[]>('dashboards'),
  );
  const [open, setOpen] = useState(false),
    [items, setItems] = useState<Dashboard[]>([]),
    [error, setError] = useState(''),
    [query, setQuery] = useState(''),
    [loading, setLoading] = useState(false),
    [busy, setBusy] = useState(false),
    [active, setActive] = useState<string>();
  const selected = Array.isArray(snapshot.value)
    ? snapshot.value
        .filter(value => typeof value === 'string' && safeDashboard(value))
        .slice(0, 8)
    : [];
  useEffect(() => {
    const subscription = bucket
      .observe$<string[]>('dashboards')
      .subscribe(setSnapshot);
    return () => subscription.unsubscribe();
  }, [bucket]);
  const load = async () => {
    setLoading(true);
    setError('');
    try {
      const dashboards = await api.listDashboards('title != ""');
      setItems(dashboards.filter(item => safeDashboard(item.url)));
    } catch {
      setError(
        'Grafana dashboard list is unavailable. Existing selections remain saved.',
      );
    } finally {
      setLoading(false);
    }
  };
  const toggle = async (url: string) => {
    setBusy(true);
    try {
      await bucket.set(
        'dashboards',
        selected.includes(url)
          ? selected.filter(value => value !== url)
          : [...selected, url].slice(0, 8),
      );
    } catch {
      setError('Dashboard selection could not be saved.');
    } finally {
      setBusy(false);
    }
  };
  const embed = active && safeDashboard(active);
  if (embed) {
    embed.searchParams.set('kiosk', '');
    embed.searchParams.set('refresh', '1m');
    embed.searchParams.set('from', 'now-1h');
    embed.searchParams.set('to', 'now');
  }
  return (
    <Box mb={3} p={2} border={1} borderColor="divider" borderRadius={20}>
      <Box
        display="flex"
        justifyContent="space-between"
        flexWrap="wrap"
        style={{ gap: 8 }}
      >
        <Typography variant="h5" component="h2">
          Your live dashboards
        </Typography>
        <Button
          color="primary"
          variant="outlined"
          disabled={snapshot.presence === 'unknown'}
          onClick={() => {
            setOpen(true);
            load();
          }}
        >
          Choose Grafana dashboards
        </Button>
      </Box>
      <Typography variant="body2" color="textSecondary">
        Save up to eight dashboards to your account.{' '}
        {embeddingEnabled
          ? 'Charts load only when opened.'
          : 'Dashboards open in Grafana; embedded viewing awaits secure Grafana configuration.'}
      </Typography>
      {!selected.length && (
        <Typography variant="caption">No dashboards selected yet.</Typography>
      )}
      <Box mt={1} display="flex" flexWrap="wrap" style={{ gap: 8 }}>
        {selected.map(url => (
          <Box key={url}>
            <Button
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              onClick={event => {
                if (embeddingEnabled) {
                  event.preventDefault();
                  setActive(active === url ? undefined : url);
                }
              }}
              aria-pressed={embeddingEnabled && active === url}
            >
              {items.find(item => item.url === url)?.title ??
                new URL(url).pathname.split('/').pop() ??
                'Dashboard'}
            </Button>
            <Link
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Open dashboard in Grafana"
            >
              ↗
            </Link>
          </Box>
        ))}
      </Box>
      {embed && (
        <>
          <Typography variant="caption" display="block">
            Uses your separate Grafana session. If sign-in or framing is
            blocked, open the dashboard in Grafana. A frame loading is not proof
            of fresh metrics.
          </Typography>
          <Button onClick={() => setActive(undefined)}>Close chart</Button>
          <Box height={480}>
            <DashboardViewer embedUrl={embed.href} />
          </Box>
        </>
      )}
      <Dialog
        open={open}
        onClose={() => setOpen(false)}
        fullWidth
        maxWidth="md"
        aria-labelledby="grafana-picker-title"
      >
        <DialogTitle id="grafana-picker-title">
          Choose live Grafana dashboards
        </DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            label="Find a dashboard"
            value={query}
            onChange={e => setQuery(e.target.value)}
          />
          <Typography variant="caption" display="block">
            Bounded list of up to 100 dashboards · {selected.length}/8 selected.
            Dashboard access follows the configured Viewer integration.
          </Typography>
          {loading && <Typography role="status">Loading Grafana…</Typography>}
          {error && (
            <Typography color="error" role="alert">
              {error}
            </Typography>
          )}
          {items
            .filter(item =>
              item.title.toLowerCase().includes(query.toLowerCase()),
            )
            .map(item => (
              <Box key={item.url}>
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={selected.includes(item.url)}
                      disabled={
                        busy ||
                        (!selected.includes(item.url) && selected.length >= 8)
                      }
                      onChange={() => toggle(item.url)}
                    />
                  }
                  label={item.title}
                />
              </Box>
            ))}
          {selected
            .filter(url => !items.some(item => item.url === url))
            .map(url => (
              <Box key={url}>
                <Typography variant="caption">
                  Saved dashboard not in current list
                </Typography>
                <Button disabled={busy} onClick={() => toggle(url)}>
                  Remove saved selection
                </Button>
              </Box>
            ))}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Done</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

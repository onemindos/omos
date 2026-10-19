import { useEffect, useMemo, useRef, useState } from 'react';
import { storageApiRef, useApi } from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';
import {
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Link,
  MenuItem,
  TextField,
  Typography,
} from '@material-ui/core';
import { AppEntry } from './useDirectory';
import {
  AccessPoint,
  catalogPath,
  methods,
  parseAccessPoints,
  platforms,
  safeAccessUrl,
} from './access';

type Saved = { points: AccessPoint[]; preferred?: string };
type Choices = Record<string, Saved>;
type Device = {
  ref: string;
  title: string;
  type: string;
  dns?: string;
  os?: string;
  observedAt?: string;
};
export function useAppAccess(app: AppEntry) {
  const storage = useApi(storageApiRef);
  const bucket = useMemo(
    () => storage.forBucket('onemind-app-access'),
    [storage],
  );
  const [snapshot, setSnapshot] = useState(() =>
    bucket.snapshot<Choices>('choices'),
  );
  useEffect(() => {
    setSnapshot(bucket.snapshot('choices'));
    const sub = bucket.observe$<Choices>('choices').subscribe(setSnapshot);
    return () => sub.unsubscribe();
  }, [bucket]);
  const saved = snapshot.value?.[app.ref];
  const base: AccessPoint = {
    id: 'web',
    label: 'Web access',
    platform: 'Web',
    method: 'Browser',
    location:
      app.hosting === 'digitalocean-nyc3'
        ? 'OneMind Cloud — NYC3'
        : 'Website / external provider',
    url: app.url,
    evidence: 'Link only',
  };
  const points = [
    base,
    ...(app.accessPoints ?? []),
    ...parseAccessPoints(JSON.stringify(saved?.points ?? [])),
  ];
  return {
    points,
    preferred: points.find(x => x.id === saved?.preferred) ?? base,
    ready: snapshot.presence !== 'unknown',
    save: async (value: Saved) => {
      const current = bucket.snapshot<Choices>('choices').value ?? {};
      await bucket.set('choices', { ...current, [app.ref]: value });
    },
  };
}
export function AccessDeployments({
  app,
  open,
  onClose,
}: {
  app: AppEntry;
  open: boolean;
  onClose: () => void;
}) {
  const access = useAppAccess(app);
  const catalog = useApi(catalogApiRef);
  const [devices, setDevices] = useState<Device[]>([]);
  const [related, setRelated] = useState<Device[]>([]);
  const [inventoryError, setInventoryError] = useState(false);
  const [loading, setLoading] = useState(false);
  const [draft, setDraft] = useState<AccessPoint[]>([]);
  const [preferred, setPreferred] = useState('web');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const initialized = useRef(false);
  const dependenciesKey = JSON.stringify(app.dependencies);
  useEffect(() => {
    if (!open) {
      initialized.current = false;
      return;
    }
    if (!access.ready || initialized.current) return;
    initialized.current = true;
    setDraft(
      access.points
        .filter(x => x.evidence === 'User recorded')
        .map(x => ({ ...x })),
    );
    setPreferred(access.preferred.id);
    setError('');
    // Only initialize an editing session when opened; catalog refresh must not erase edits.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, access.ready]);
  useEffect(() => {
    let active = false;
    const cleanup = () => {
      active = false;
    };
    if (!open) return cleanup;
    active = true;
    setLoading(true);
    setInventoryError(false);
    (async () => {
      try {
        const result: Device[] = [];
        let cursor: string | undefined;
        for (let page = 0; page < 20; page++) {
          const response = await catalog.queryEntities(
            cursor
              ? { cursor, limit: 250 }
              : { filter: { kind: 'Resource' }, limit: 250 },
          );
          for (const e of response.items) {
            const type = String(e.spec?.type ?? '');
            if (!['mesh-endpoint', 'kubernetes-cluster'].includes(type))
              continue;
            const a = e.metadata.annotations ?? {};
            result.push({
              ref: `resource:${e.metadata.namespace ?? 'default'}/${
                e.metadata.name
              }`,
              title: e.metadata.title ?? e.metadata.name,
              type,
              dns: a['onemind.io/mesh-dns'],
              os:
                a['onemind.io/mesh-os'] ??
                e.metadata.description?.match(
                  /^(macOS|linux|windows|android|iOS)\b/i,
                )?.[1],
              observedAt: a['onemind.io/inventory-observed-at'],
            });
          }
          cursor = response.pageInfo.nextCursor;
          if (!cursor) break;
          if (page === 19) throw Error('Inventory limit');
        }
        const dependencies = await Promise.all(
          (app.dependencies ?? []).map(ref => catalog.getEntityByRef(ref)),
        );
        const relations: Device[] = dependencies.filter(Boolean).map(e => ({
          ref: `${e!.kind.toLowerCase()}:${
            e!.metadata.namespace ?? 'default'
          }/${e!.metadata.name}`,
          title: e!.metadata.title ?? e!.metadata.name,
          type: String(e!.spec?.type ?? ''),
          observedAt:
            e!.metadata.annotations?.['onemind.io/inventory-observed-at'] ??
            e!.metadata.annotations?.['onemind.io/verified-at'],
        }));
        if (active) {
          setDevices(result);
          setRelated(relations);
        }
      } catch {
        if (active) setInventoryError(true);
      } finally {
        if (active) setLoading(false);
      }
    })();
    return cleanup;
    // app.dependencies is captured by dependenciesKey to avoid needless refetches.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, catalog, app.ref, dependenciesKey]);
  const update = (index: number, value: Partial<AccessPoint>) =>
    setDraft(current =>
      current.map((x, i) => (i === index ? { ...x, ...value } : x)),
    );
  const recorded = access.points.filter(x => x.evidence !== 'User recorded');
  const save = async () => {
    if (
      draft.some(
        x =>
          !x.label.trim() ||
          !x.location.trim() ||
          (x.url && !safeAccessUrl(x.url)),
      )
    ) {
      setError(
        'Add an access name and location. Web links must use HTTPS without embedded credentials.',
      );
      return;
    }
    if (![...recorded, ...draft].some(x => x.id === preferred)) {
      setError('Choose an existing preferred access point.');
      return;
    }
    setBusy(true);
    try {
      await access.save({
        points: draft.map(x => ({
          ...x,
          evidence: 'User recorded',
          observedAt: new Date().toISOString(),
          url: safeAccessUrl(x.url),
        })),
        preferred,
      });
      onClose();
    } catch {
      setError('Access points could not be saved. Your edits are still here.');
    } finally {
      setBusy(false);
    }
  };
  const target = (point: AccessPoint) => {
    const device = devices.find(d => d.ref === point.targetRef);
    const href = catalogPath(point.targetRef);
    return (
      <>
        {href && (
          <Typography variant="body2">
            <Link href={href}>
              Device / cluster: {device?.title ?? point.targetRef}
            </Link>
          </Typography>
        )}
        {device?.os && (
          <Typography variant="caption" component="p">
            Recorded device OS: {device.os}
          </Typography>
        )}
        {device?.dns && (
          <Typography variant="caption" component="p">
            Tailscale: {device.dns}
          </Typography>
        )}
        {device?.observedAt && (
          <Typography variant="caption" component="p">
            Device discovery: {device.observedAt}
          </Typography>
        )}
      </>
    );
  };
  return (
    <Dialog
      open={open}
      onClose={() => !busy && onClose()}
      fullWidth
      maxWidth="md"
      aria-labelledby={`access-${app.name}`}
    >
      <DialogTitle id={`access-${app.name}`}>
        {app.title} · Access & Deployments
      </DialogTitle>
      <DialogContent>
        <Typography variant="body2">
          Platforms, installation methods, and actual locations are recorded
          separately. A website link or Tailscale device does not establish an
          installed app or live availability.
        </Typography>
        <Box mt={2}>
          <Typography variant="subtitle2">Available for</Typography>
          <Typography variant="body2">
            {app.supportedPlatforms?.length
              ? app.supportedPlatforms.join(', ')
              : 'Supported platforms not recorded'}
          </Typography>
        </Box>
        <Typography variant="h6" style={{ marginTop: 16 }}>
          Recorded access points
        </Typography>
        {recorded.map(point => (
          <Box
            key={point.id}
            p={2}
            my={1}
            border={1}
            borderColor="divider"
            borderRadius={8}
          >
            <Typography variant="subtitle1">{point.label}</Typography>
            <Chip size="small" label={point.platform} />{' '}
            <Chip size="small" label={point.method} />
            <Typography variant="body2">
              {point.location}
              {point.version ? ` · version ${point.version}` : ''}
            </Typography>
            <Typography variant="caption" component="p">
              Evidence: {point.evidence}
              {point.observedAt ? ` · ${point.observedAt}` : ''}
            </Typography>
            {target(point)}
            {point.instructions && (
              <Typography variant="body2">{point.instructions}</Typography>
            )}
            {point.url && (
              <Button
                size="small"
                href={point.url}
                target="_blank"
                rel="noopener noreferrer"
              >
                Open {point.label}
              </Button>
            )}
          </Box>
        ))}
        <Typography variant="h6" style={{ marginTop: 16 }}>
          Catalog deployment relationships
        </Typography>
        {loading && (
          <Typography role="status">
            Loading device and workload inventory…
          </Typography>
        )}
        {inventoryError && (
          <Typography role="alert" color="error">
            Inventory lookup failed. Previous device results may be stale.
          </Typography>
        )}
        {!loading && !related.length && (
          <Typography variant="body2">
            No deployment relationship recorded.
          </Typography>
        )}
        {related.map(item => (
          <Box key={item.ref} my={1}>
            <Link href={catalogPath(item.ref)}>{item.title}</Link>
            <Typography variant="caption" component="p">
              {item.type} · catalog relationship snapshot
              {item.observedAt ? ` · ${item.observedAt}` : ''}
            </Typography>
          </Box>
        ))}
        <Typography variant="h6" style={{ marginTop: 16 }}>
          Your additional access points
        </Typography>
        {draft.map((point, index) => (
          <Box
            key={point.id}
            p={2}
            my={2}
            border={1}
            borderColor="divider"
            borderRadius={8}
          >
            <TextField
              fullWidth
              id={`access-name-${index}`}
              label="Access name"
              value={point.label}
              onChange={e => update(index, { label: e.target.value })}
              margin="dense"
              inputProps={{ maxLength: 100 }}
              disabled={busy}
            />
            <TextField
              fullWidth
              select
              id={`access-platform-${index}`}
              label="Platform"
              value={point.platform}
              onChange={e => update(index, { platform: e.target.value })}
              margin="dense"
              disabled={busy}
            >
              {platforms.map(x => (
                <MenuItem key={x} value={x}>
                  {x}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              fullWidth
              select
              id={`access-method-${index}`}
              label="Installation method"
              value={point.method}
              onChange={e => update(index, { method: e.target.value })}
              margin="dense"
              disabled={busy}
            >
              {methods.map(x => (
                <MenuItem key={x} value={x}>
                  {x}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              fullWidth
              select
              id={`access-device-${index}`}
              label="Device or cluster (optional)"
              value={point.targetRef ?? ''}
              onChange={e => {
                const d = devices.find(x => x.ref === e.target.value);
                update(index, {
                  targetRef: e.target.value || undefined,
                  ...(d ? { location: d.title } : {}),
                });
              }}
              margin="dense"
              disabled={busy || loading}
            >
              <MenuItem value="">Enter location manually</MenuItem>
              {point.targetRef &&
                !devices.some(x => x.ref === point.targetRef) && (
                  <MenuItem value={point.targetRef}>{point.targetRef}</MenuItem>
                )}
              {devices.map(x => (
                <MenuItem key={x.ref} value={x.ref}>
                  {x.title} · {x.os ?? x.type}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              fullWidth
              id={`access-location-${index}`}
              label="Actual location"
              value={point.location}
              onChange={e => update(index, { location: e.target.value })}
              margin="dense"
              inputProps={{ maxLength: 300 }}
              disabled={busy}
            />
            <TextField
              fullWidth
              id={`access-url-${index}`}
              label="HTTPS access URL (optional)"
              value={point.url ?? ''}
              onChange={e => update(index, { url: e.target.value })}
              margin="dense"
              inputProps={{ maxLength: 2000 }}
              disabled={busy}
            />
            <TextField
              fullWidth
              id={`access-instructions-${index}`}
              label="How to open / installation notes"
              value={point.instructions ?? ''}
              onChange={e => update(index, { instructions: e.target.value })}
              margin="dense"
              inputProps={{ maxLength: 2000 }}
              disabled={busy}
            />
            <Typography variant="caption">
              User recorded · installation is not independently verified
            </Typography>
            <Button
              disabled={busy}
              onClick={() =>
                setDraft(current => current.filter(x => x.id !== point.id))
              }
            >
              Remove access point
            </Button>
          </Box>
        ))}
        <Button
          disabled={busy || !access.ready || draft.length >= 20}
          onClick={() =>
            setDraft(current => [
              ...current,
              {
                id: crypto.randomUUID(),
                label: '',
                platform: 'Web',
                method: 'Browser',
                location: '',
                evidence: 'User recorded',
              },
            ])
          }
        >
          Add access point
        </Button>
        <TextField
          fullWidth
          select
          id="access-preferred"
          label="Preferred access on this account"
          value={preferred}
          onChange={e => setPreferred(e.target.value)}
          margin="normal"
          disabled={busy}
        >
          {[...recorded, ...draft].map(point => (
            <MenuItem key={point.id} value={point.id}>
              {point.label || 'New access point'} · {point.platform}
            </MenuItem>
          ))}
        </TextField>
        <Typography variant="caption">
          Native access displays instructions. Remote device app launching is
          not enabled.
        </Typography>
        {error && (
          <Typography role="alert" color="error">
            {error}
          </Typography>
        )}
      </DialogContent>
      <DialogActions>
        <Button disabled={busy} onClick={onClose}>
          Close
        </Button>
        <Button
          disabled={busy || !access.ready}
          variant="contained"
          color="primary"
          onClick={save}
        >
          Save access preferences
        </Button>
      </DialogActions>
    </Dialog>
  );
}
export function AccessButton({ app }: { app: AppEntry }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <Button size="small" onClick={() => setOpen(true)}>
        Access & deployments
      </Button>
      {open && (
        <AccessDeployments
          app={app}
          open={open}
          onClose={() => setOpen(false)}
        />
      )}
    </>
  );
}

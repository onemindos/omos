import { useEffect, useMemo, useState } from 'react';
import { storageApiRef, useApi } from '@backstage/core-plugin-api';
import {
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  TextField,
  Typography,
} from '@material-ui/core';
import { useDirectory } from './useDirectory';
import { HostingBadge } from './HostingBadge';
import { AccessButton, useAppAccess, AccessDeployments } from './AccessDeployments';

const purposes = [
  'Notes',
  'Email',
  'Passwords',
  'Files',
  'Maps',
  'Photos',
  'Voice memos',
  'Wallet',
  'Contacts',
  'Calendar',
  'SMS',
];
type Assignment = {
  ref: string;
  title?: string;
  role: 'Primary' | 'Supporting' | 'Alternative';
  use: string;
};
type Decision = {
  at: string;
  reason: string;
  apps: Assignment[];
  data: string;
  migration: string;
};
type Choice = {
  apps: Assignment[];
  data: string;
  migration: string;
  history: Decision[];
};
type Stack = Record<string, Choice>;
const empty = (): Choice => ({
  apps: [],
  data: '',
  migration: '',
  history: [],
});

export function ActiveStack() {
  const storage = useApi(storageApiRef);
  const bucket = useMemo(
    () => storage.forBucket('onemind-active-stack'),
    [storage],
  );
  const [snapshot, setSnapshot] = useState(() =>
    bucket.snapshot<Stack>('choices'),
  );
  useEffect(() => {
    setSnapshot(bucket.snapshot<Stack>('choices'));
    const subscription = bucket
      .observe$<Stack>('choices')
      .subscribe(setSnapshot);
    return () => subscription.unsubscribe();
  }, [bucket]);
  const { all, loaded, error: catalogError } = useDirectory();
  const apps = [...all].sort((a, b) => a.title.localeCompare(b.title));
  const stack = snapshot.value ?? {};
  const ready = snapshot.presence !== 'unknown' && loaded;
  const [purpose, setPurpose] = useState('');
  const [draft, setDraft] = useState<Choice>(empty);
  const [reason, setReason] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const update = (index: number, value: Partial<Assignment>) =>
    setDraft(current => ({
      ...current,
      apps: current.apps.map((item, i) =>
        i === index ? { ...item, ...value } : item,
      ),
    }));
  const save = async () => {
    if (
      draft.apps.some(
        item =>
          !item.ref ||
          !item.use.trim() ||
          (item.ref.startsWith('manual:') && !item.title?.trim()),
      )
    ) {
      setError('Choose an app and describe what you use it for.');
      return;
    }
    if (draft.apps.filter(item => item.role === 'Primary').length > 1) {
      setError(
        'Choose one primary app. Other active apps can have supporting roles.',
      );
      return;
    }
    if (new Set(draft.apps.map(item => item.ref)).size !== draft.apps.length) {
      setError('Assign each app once per purpose.');
      return;
    }
    if (!reason.trim()) {
      setError('Add a brief reason for this decision.');
      return;
    }
    setBusy(true);
    try {
      // Re-read before saving to preserve other purposes changed by another view.
      const current = bucket.snapshot<Stack>('choices').value ?? {};
      const decision: Decision = {
        at: new Date().toISOString(),
        reason: reason.trim(),
        apps: draft.apps,
        data: draft.data,
        migration: draft.migration,
      };
      await bucket.set('choices', {
        ...current,
        [purpose]: {
          ...draft,
          history: [...(current[purpose]?.history ?? []), decision].slice(-50),
        },
      });
      setPurpose('');
    } catch {
      setError('Your decision could not be saved. Try again.');
    } finally {
      setBusy(false);
    }
  };
  return (
    <Box my={3}>
      <Typography variant="h5" component="h2" style={{ fontWeight: 700 }}>
        My Active Stack
      </Typography>
      <Typography color="textSecondary" variant="body2">
        Your chosen tools by purpose. Saved to your account across devices;
        availability is shown separately in Today.
      </Typography>
      {!ready && <Typography role="status">Loading your stack…</Typography>}
      {catalogError && (
        <Typography role="alert" color="error">
          App catalog unavailable. Saved choices remain; app links may be stale.
        </Typography>
      )}
      <Box
        mt={2}
        style={{
          display: 'grid',
          gridTemplateColumns:
            'repeat(auto-fit, minmax(min(100%, 260px), 1fr))',
          gap: 16,
        }}
      >
        {purposes.map(name => {
          const choice = stack[name] ?? empty();
          const active = choice.apps.filter(
            item => item.role !== 'Alternative',
          );
          return (
            <Box
              key={name}
              p={2}
              border={1}
              borderColor="divider"
              borderRadius={12}
            >
              <Typography variant="h6" component="h3">
                {name}
              </Typography>
              {!active.length && (
                <Typography color="textSecondary">Not assigned</Typography>
              )}
              {active.map(item => {
                const app = apps.find(candidate => candidate.ref === item.ref);
                return (
                  <Box key={item.ref} mt={1}>
                    <Chip size="small" label={item.role} />
                    <HostingBadge hosting={app?.hosting} status={app?.status} />{' '}
                    <Typography component="span">
                      {app?.title ?? item.title ?? 'App no longer in directory'}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {item.use}
                    </Typography>
                    {app && <PreferredAccess app={app} />}

                  </Box>
                );
              })}
              <Button
                disabled={!ready}
                size="small"
                onClick={() => {
                  setPurpose(name);
                  setDraft({
                    ...choice,
                    apps: choice.apps.map(item => ({ ...item })),
                  });
                  setReason('');
                  setError('');
                }}
              >
                Manage {name}
              </Button>
              {!!choice.apps.filter(item => item.role === 'Alternative')
                .length && (
                <Typography variant="caption" component="p">
                  {
                    choice.apps.filter(item => item.role === 'Alternative')
                      .length
                  }{' '}
                  alternatives saved
                </Typography>
              )}
            </Box>
          );
        })}
      </Box>
      <Dialog
        open={!!purpose}
        onClose={() => !busy && setPurpose('')}
        fullWidth
        maxWidth="md"
        aria-labelledby="active-stack-title"
      >
        <DialogTitle id="active-stack-title">
          Your {purpose.toLowerCase()} stack
        </DialogTitle>
        <DialogContent>
          <Typography variant="body2">
            Choose a primary app, supporting tools, and alternatives. Choose a directory app or enter a native app by name.
          </Typography>
          {draft.apps.map((item, index) => (
            <Box
              key={index}
              mt={2}
              p={2}
              border={1}
              borderColor="divider"
              borderRadius={8}
            >
              <TextField
                select
                fullWidth
                id={`stack-app-${index}`}
                label="App"
                value={item.ref.startsWith('manual:') ? 'manual' : item.ref}
                onChange={event =>
                  update(index, {
                    ref:
                      event.target.value === 'manual'
                        ? `manual:${crypto.randomUUID()}`
                        : event.target.value,
                    title: apps.find(app => app.ref === event.target.value)
                      ?.title,
                  })
                }
                margin="dense"
                disabled={busy}
              >
                <MenuItem value="">Choose an app</MenuItem>
                <MenuItem value="manual">
                  Native or other app — enter a name
                </MenuItem>
                {item.ref &&
                  !item.ref.startsWith('manual:') &&
                  !apps.some(app => app.ref === item.ref) && (
                    <MenuItem value={item.ref}>
                      Unavailable directory entry ({item.ref})
                    </MenuItem>
                  )}
                {apps.map(app => (
                  <MenuItem key={app.ref} value={app.ref}>
                    {app.title} · {app.groupTitle}
                  </MenuItem>
                ))}
              </TextField>
              {item.ref.startsWith('manual:') && (
                <TextField
                  fullWidth
                  id={`stack-name-${index}`}
                  label="App name"
                  value={item.title ?? ''}
                  onChange={event =>
                    update(index, { title: event.target.value })
                  }
                  inputProps={{ maxLength: 100 }}
                  disabled={busy}
                />
              )}
              <TextField
                select
                fullWidth
                id={`stack-role-${index}`}
                label="Role"
                value={item.role}
                onChange={event =>
                  update(index, {
                    role: event.target.value as Assignment['role'],
                  })
                }
                margin="dense"
                disabled={busy}
              >
                {(['Primary', 'Supporting', 'Alternative'] as const).map(
                  role => (
                    <MenuItem key={role} value={role}>
                      {role}
                    </MenuItem>
                  ),
                )}
              </TextField>
              <TextField
                fullWidth
                id={`stack-use-${index}`}
                label="What I use this app for"
                value={item.use}
                onChange={event => update(index, { use: event.target.value })}
                margin="dense"
                inputProps={{ maxLength: 500 }}
                disabled={busy}
              />
              <Button
                disabled={busy}
                onClick={() =>
                  setDraft(current => ({
                    ...current,
                    apps: current.apps.filter((_, i) => i !== index),
                  }))
                }
              >
                Remove assignment
              </Button>
            </Box>
          ))}
          <Button
            disabled={busy || draft.apps.length >= 20}
            onClick={() =>
              setDraft(current => ({
                ...current,
                apps: [
                  ...current.apps,
                  {
                    ref: '',
                    role: current.apps.length ? 'Supporting' : 'Primary',
                    use: '',
                  },
                ],
              }))
            }
          >
            Add app
          </Button>
          <TextField
            fullWidth
            multiline
            id="stack-data"
            label="Where my data lives / export or backup"
            value={draft.data}
            onChange={event =>
              setDraft(current => ({ ...current, data: event.target.value }))
            }
            margin="normal"
            inputProps={{ maxLength: 2000 }}
            disabled={busy}
          />
          <TextField
            fullWidth
            multiline
            id="stack-migration"
            label="Migration or follow-up still needed"
            value={draft.migration}
            onChange={event =>
              setDraft(current => ({
                ...current,
                migration: event.target.value,
              }))
            }
            margin="normal"
            inputProps={{ maxLength: 2000 }}
            disabled={busy}
          />
          <TextField
            fullWidth
            id="stack-reason"
            label="Reason for this decision"
            value={reason}
            onChange={event => setReason(event.target.value)}
            margin="normal"
            inputProps={{ maxLength: 1000 }}
            disabled={busy}
          />
          {error && (
            <Typography role="alert" color="error">
              {error}
            </Typography>
          )}
          {!!draft.history.length && (
            <Box mt={2}>
              <Typography variant="h6">Decision history</Typography>
              <Typography variant="caption">
                Most recent 50 decisions
              </Typography>
              {[...draft.history].reverse().map(entry => (
                <Box
                  key={entry.at}
                  mt={1}
                  p={1}
                  borderBottom={1}
                  borderColor="divider"
                >
                  <Typography variant="body2">
                    {new Date(entry.at).toLocaleString()} · {entry.reason}
                  </Typography>
                  {entry.apps.map(item => (
                    <Typography key={item.ref} variant="caption" component="p">
                      {item.role}:{' '}
                      {apps.find(app => app.ref === item.ref)?.title ??
                        item.ref}{' '}
                      — {item.use}
                    </Typography>
                  ))}
                  {entry.data && (
                    <Typography variant="caption" component="p">
                      Data: {entry.data}
                    </Typography>
                  )}
                  {entry.migration && (
                    <Typography variant="caption" component="p">
                      Follow-up: {entry.migration}
                    </Typography>
                  )}
                </Box>
              ))}
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button disabled={busy} onClick={() => setPurpose('')}>
            Cancel
          </Button>
          <Button
            disabled={busy || !ready}
            variant="contained"
            color="primary"
            onClick={save}
          >
            Save decision
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

function PreferredAccess({app}: {app: import('./useDirectory').AppEntry}) {
  const access = useAppAccess(app);
  const [open,setOpen] = useState(false);
  return <Box>{access.preferred.url ? <Button size="small" color="primary" href={access.preferred.url} target="_blank" rel="noopener noreferrer">Open {app.title} · {access.preferred.platform}</Button> : <Button size="small" onClick={()=>setOpen(true)}>Open instructions · {access.preferred.platform}</Button>}<AccessButton app={app} />{open && <AccessDeployments app={app} open={open} onClose={()=>setOpen(false)} />}</Box>;
}

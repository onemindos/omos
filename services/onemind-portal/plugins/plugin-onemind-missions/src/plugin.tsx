import { useCallback, useEffect, useState } from 'react';
import {
  Box,
  Button,
  Chip,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Typography,
} from '@material-ui/core';
import AssignmentIcon from '@material-ui/icons/Assignment';
import {
  createFrontendPlugin,
  PageBlueprint,
} from '@backstage/frontend-plugin-api';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';

type Mission = {
  id: string;
  state: string;
  actor: string;
  specification: { capability: string; prompt: string; target: string };
  created_at: string;
  updated_at: string;
  result?: { answer?: string; [key: string]: unknown };
  error?: string;
};
type ObjectRow = {
  id: string;
  kind: string;
  name: string;
  source: string;
  observed_at: string;
  properties: Record<string, unknown>;
};
type Agent = { ref: string; name: string; owner: string; agent: string };

export function Missions() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [missions, setMissions] = useState<Mission[]>([]);
  const [objects, setObjects] = useState<ObjectRow[]>([]);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [capability, setCapability] = useState('inspect-system');
  const [prompt, setPrompt] = useState('Inspect current OneMind conditions');
  const [target, setTarget] = useState('');
  const [error, setError] = useState('');
  const [canOperate, setCanOperate] = useState(false);
  const [busy, setBusy] = useState(false);
  const [selected, setSelected] = useState<ObjectRow>();
  const [objectId, setObjectId] = useState('');
  const [objectName, setObjectName] = useState('');
  const [objectKind, setObjectKind] = useState('asset');
  const request = useCallback(
    async (path: string, body?: unknown) => {
      const base = await discovery.getBaseUrl('onemind-missions');
      const response = await fetchApi.fetch(
        base + path,
        body === undefined
          ? {}
          : {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(body),
            },
      );
      const data = await response.json();
      if (!response.ok)
        throw new Error(
          data.error ?? data.detail ?? `Request failed (${response.status})`,
        );
      return data;
    },
    [discovery, fetchApi],
  );
  const refresh = useCallback(async () => {
    const results = await Promise.allSettled([
      request('/viewer'),
      request('/missions'),
      request('/objects'),
      request('/agents'),
    ]);
    if (results[0].status === 'fulfilled')
      setCanOperate(results[0].value.canOperate);
    if (results[1].status === 'fulfilled') {
      setMissions(results[1].value);
      setError('');
    } else setError('Mission observations are unavailable.');
    if (results[2].status === 'fulfilled') setObjects(results[2].value.objects);
    if (results[3].status === 'fulfilled') setAgents(results[3].value.agents);
  }, [request]);
  useEffect(() => {
    refresh();
    const timer = setInterval(refresh, 5000);
    return () => clearInterval(timer);
  }, [refresh]);
  async function run() {
    setBusy(true);
    setError('');
    try {
      await request('/missions', {
        capability,
        prompt,
        target,
        idempotency_key: crypto.randomUUID(),
      });
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Mission could not be started');
    } finally {
      setBusy(false);
    }
  }
  async function cancel(id: string) {
    try {
      await request(`/missions/${id}/cancel`, {});
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Cancellation unavailable');
    }
  }
  async function createObject() {
    try {
      await request('/objects', {
        id: objectId,
        kind: objectKind,
        name: objectName,
      });
      setObjectId('');
      setObjectName('');
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Object could not be created');
    }
  }
  return (
    <Box p={3}>
      <Typography variant="overline">ONE MIND · OPERATIONS</Typography>
      <Typography variant="h4" component="h1">
        Missions and world view
      </Typography>
      <Typography color="textSecondary">
        Run work, follow its results, and inspect the latest source
        observations.
      </Typography>
      {error && (
        <Box mt={2}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}
      <Box mt={3} p={3} border={1} borderColor="divider" borderRadius={12}>
        <Typography variant="h6">Start a mission</Typography>
        <Box mt={2} display="flex" flexWrap="wrap" style={{ gap: 16 }}>
          <FormControl variant="outlined" style={{ minWidth: 220 }}>
            <InputLabel>Capability</InputLabel>
            <Select
              label="Capability"
              value={capability}
              onChange={event => {
                setCapability(String(event.target.value));
                setTarget('');
              }}
            >
              <MenuItem value="inspect-system">Inspect infrastructure</MenuItem>
              <MenuItem value="summarize-site">
                Summarize current evidence
              </MenuItem>
              <MenuItem value="delegate-agent">Delegate to an agent</MenuItem>
              <MenuItem value="restart-worker">
                Restart a processing worker
              </MenuItem>
            </Select>
          </FormControl>
          {capability === 'delegate-agent' && (
            <FormControl variant="outlined" style={{ minWidth: 260 }}>
              <InputLabel>Discovered agent</InputLabel>
              <Select
                label="Discovered agent"
                value={target}
                onChange={event => setTarget(String(event.target.value))}
              >
                {agents.map(agent => (
                  <MenuItem key={agent.ref} value={agent.ref}>
                    {agent.name} · {agent.owner}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}
          {capability === 'restart-worker' && (
            <FormControl variant="outlined" style={{ minWidth: 260 }}>
              <InputLabel>Worker</InputLabel>
              <Select
                label="Worker"
                value={target}
                onChange={event => setTarget(String(event.target.value))}
              >
                <MenuItem value="onemind-core/onemind-ingestion">
                  Source projection worker
                </MenuItem>
                <MenuItem value="onemind-ai/docling">
                  Document processing worker
                </MenuItem>
              </Select>
            </FormControl>
          )}
        </Box>
        <Box my={2}>
          <TextField
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            label="Request"
            value={prompt}
            onChange={event => setPrompt(event.target.value)}
          />
        </Box>
        <Button
          color="primary"
          variant="contained"
          disabled={
            !canOperate ||
            busy ||
            (['delegate-agent', 'restart-worker'].includes(capability) &&
              !target)
          }
          onClick={run}
        >
          {busy ? 'Submitting…' : 'Run mission'}
        </Button>
      </Box>
      <Box mt={3}>
        <Typography variant="h6">Mission history</Typography>
        {missions.length === 0 && (
          <Typography color="textSecondary">No missions observed.</Typography>
        )}
        {missions.slice(0, 30).map(mission => (
          <Box
            key={mission.id}
            my={2}
            p={2}
            border={1}
            borderColor="divider"
            borderRadius={12}
          >
            <Box
              display="flex"
              alignItems="center"
              flexWrap="wrap"
              style={{ gap: 12 }}
            >
              <Chip size="small" label={mission.state} />
              <Typography variant="subtitle1">
                {mission.specification.capability}
              </Typography>
              {canOperate &&
                !['completed', 'failed', 'cancelled'].includes(
                  mission.state,
                ) && (
                  <Button size="small" onClick={() => cancel(mission.id)}>
                    Cancel
                  </Button>
                )}
            </Box>
            <Typography>{mission.specification.prompt}</Typography>
            <Typography variant="caption" color="textSecondary">
              {mission.actor} · updated{' '}
              {new Date(mission.updated_at).toLocaleString()} · {mission.id}
            </Typography>
            {mission.error && (
              <Typography color="error">
                Execution failed: {mission.error}
              </Typography>
            )}
            {mission.result?.answer && (
              <Typography style={{ whiteSpace: 'pre-wrap', marginTop: 12 }}>
                {mission.result.answer}
              </Typography>
            )}
            {mission.result && (
              <details>
                <summary>Execution evidence</summary>
                <pre
                  style={{
                    whiteSpace: 'pre-wrap',
                    overflowWrap: 'anywhere',
                    fontSize: 12,
                  }}
                >
                  {JSON.stringify(mission.result, null, 2)}
                </pre>
              </details>
            )}
          </Box>
        ))}
      </Box>
      <Box mt={4}>
        <Typography variant="h6">
          World objects · {objects.length} observed
        </Typography>
        <Typography color="textSecondary">
          Each source record shows when it was last observed.
        </Typography>
        {canOperate && (
          <Box display="flex" flexWrap="wrap" my={2} style={{ gap: 12 }}>
            <TextField
              label="Object ID"
              placeholder="asset:greenhouse"
              value={objectId}
              onChange={e => setObjectId(e.target.value)}
            />
            <TextField
              label="Name"
              value={objectName}
              onChange={e => setObjectName(e.target.value)}
            />
            <TextField
              label="Kind"
              value={objectKind}
              onChange={e => setObjectKind(e.target.value)}
            />
            <Button disabled={!objectId || !objectName} onClick={createObject}>
              Create object
            </Button>
          </Box>
        )}
        <Box display="flex" flexWrap="wrap" style={{ gap: 8 }}>
          {objects.map(object => (
            <Button
              key={object.id}
              variant="outlined"
              size="small"
              onClick={() => setSelected(object)}
            >
              {object.name} · {object.kind}
            </Button>
          ))}
        </Box>
        {selected && (
          <Box p={2} mt={2} border={1} borderColor="divider" borderRadius={12}>
            <Typography variant="h6">{selected.name}</Typography>
            <Typography variant="caption">
              {selected.id} · {selected.source} · observed{' '}
              {new Date(selected.observed_at).toLocaleString()}
            </Typography>
            <pre style={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere' }}>
              {JSON.stringify(selected.properties, null, 2)}
            </pre>
          </Box>
        )}
      </Box>
    </Box>
  );
}
export const missionsPlugin = createFrontendPlugin({
  pluginId: 'onemind-missions',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/missions',
        title: 'Missions',
        icon: <AssignmentIcon />,
        loader: async () => <Missions />,
      },
    }),
  ],
});

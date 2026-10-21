import { useState } from 'react';
import {
  createFrontendPlugin,
  PageBlueprint,
} from '@backstage/frontend-plugin-api';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Box, Button, Link, TextField, Typography } from '@material-ui/core';
import ChatIcon from '@material-ui/icons/Chat';

type Answer = {
  mode: string;
  answer: string;
  observedAt: string;
  unavailable: string[];
  citations?: string[];
  sources: {
    id: string;
    title: string;
    url: string;
    freshness: string;
    content: string;
  }[];
};
function AssistantPage() {
  const discovery = useApi(discoveryApiRef);
  const api = useApi(fetchApiRef);
  const [question, setQuestion] = useState('');
  const [answer, setAnswer] = useState<Answer>();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const ask = async () => {
    setBusy(true);
    setError('');
    setAnswer(undefined);
    try {
      const base = await discovery.getBaseUrl('onemind-assistant');
      const response = await api.fetch(`${base}/ask`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ question }),
        signal: AbortSignal.timeout(45000),
      });
      if (!response.ok) {
        if (response.status === 403) {
          throw new Error('Administrator access required.');
        }
        if (response.status === 429) {
          throw new Error('Assistant busy. Try again shortly.');
        }
        throw new Error(
          'Assistant unavailable. Use the catalog and monitoring pages.',
        );
      }
      setAnswer(await response.json());
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Assistant unavailable');
    } finally {
      setBusy(false);
    }
  };
  return (
    <Box p={3}>
      <Typography variant="h4">OneMind assistant</Typography>
      <Typography>
        Ask about services and current observations. This assistant cannot
        change your systems.
      </Typography>
      <Typography variant="caption">
        Questions are processed by the configured model endpoint. Responses can
        be wrong; inspect the linked evidence before acting.
      </Typography>
      <Box my={2}>
        <TextField
          label="What would you like to understand?"
          fullWidth
          multiline
          variant="outlined"
          value={question}
          inputProps={{ maxLength: 2000 }}
          onChange={e => setQuestion(e.target.value)}
        />
        <Box mt={1}>
          <Button
            variant="contained"
            color="primary"
            disabled={busy || !question.trim()}
            onClick={ask}
          >
            {busy ? 'Reading observations…' : 'Ask'}
          </Button>
        </Box>
      </Box>
      {error && <Typography color="error">{error}</Typography>}
      {answer && (
        <>
          <Typography variant="h6">
            {answer.mode === 'model'
              ? 'Suggested answer'
              : 'Source explorer — model not connected'}
          </Typography>
          <Typography style={{ whiteSpace: 'pre-wrap' }}>
            {answer.answer}
          </Typography>
          <Typography variant="caption">
            Retrieved {new Date(answer.observedAt).toLocaleString()}
          </Typography>
          {answer.unavailable.map(x => (
            <Typography key={x} color="error">
              {x}
            </Typography>
          ))}
          <Typography variant="h6">Evidence</Typography>
          {answer.sources.map(source => (
            <Box key={source.id} my={2}>
              <Link href={source.url}>{source.title}</Link>
              {answer.citations?.includes(source.id) && (
                <Typography variant="caption"> · cited by model</Typography>
              )}
              <Typography variant="caption" display="block">
                {source.freshness}
              </Typography>
              <details>
                <summary>Inspect source</summary>
                <Box
                  component="pre"
                  style={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere' }}
                >
                  {source.content}
                </Box>
              </details>
            </Box>
          ))}
        </>
      )}
    </Box>
  );
}
export const assistantPlugin = createFrontendPlugin({
  pluginId: 'onemind-assistant',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/assistant',
        title: 'Assistant',
        icon: <ChatIcon />,
        loader: async () => <AssistantPage />,
      },
    }),
  ],
});

import { ReactNode, useSyncExternalStore } from 'react';
import {
  createFrontendModule,
  PluginWrapperBlueprint,
} from '@backstage/frontend-plugin-api';
import { Box, Button, Typography } from '@material-ui/core';

/**
 * The GitHub plugins ask for an OAuth token the moment a card renders, which opens a
 * "Login Required" dialog on every app page nobody asked to sign in on. These wrappers (the
 * native plugin-wrapper extension) render a plugin's cards and tabs only after one explicit
 * "Connect GitHub" click. The choice is remembered in this browser, and shared, so one click
 * connects every GitHub card on the page.
 */
const KEY = 'onemind.github.connected';
const listeners = new Set<() => void>();
let connectedInMemory = false;

function readConnected(): boolean {
  if (connectedInMemory) return true;
  try {
    return window.localStorage.getItem(KEY) === '1';
  } catch {
    return false;
  }
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

export function connectGitHub() {
  connectedInMemory = true;
  try {
    window.localStorage.setItem(KEY, '1');
  } catch {
    // Storage can be blocked; the choice then lasts until the page reloads.
  }
  listeners.forEach(listener => listener());
}

/** Test helper: forget the choice. */
export function resetGitHubConnection() {
  connectedInMemory = false;
  try {
    window.localStorage.removeItem(KEY);
  } catch {
    // ignore
  }
  listeners.forEach(listener => listener());
}

export function useGitHubConnected(): boolean {
  return useSyncExternalStore(subscribe, readConnected, () => false);
}

export function GitHubGate({ children }: { children: ReactNode }) {
  const connected = useGitHubConnected();
  if (connected) return <>{children}</>;
  return (
    <Box p={2} data-testid="github-gate">
      <Typography variant="body2" gutterBottom>
        GitHub is not connected. Pull requests and workflow runs need your GitHub sign-in.
      </Typography>
      <Button size="small" variant="outlined" color="primary" onClick={connectGitHub}>
        Connect GitHub
      </Button>
    </Box>
  );
}

const githubGateWrapper = PluginWrapperBlueprint.make({
  params: define => define({ loader: async () => ({ component: GitHubGate }) }),
});

/** One module per GitHub plugin: every component of the plugin renders behind the gate. */
export const githubGateModules = ['github-pull-requests', 'github-actions'].map(pluginId =>
  createFrontendModule({ pluginId, extensions: [githubGateWrapper] }),
);

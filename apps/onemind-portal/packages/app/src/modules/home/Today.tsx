import { FabricPanel } from './FabricPanel';
import { EventSources } from './EventSources';
import { Box, Button, Typography, useTheme } from '@material-ui/core';
import TodayIcon from '@material-ui/icons/Today';
import { createFrontendPlugin, PageBlueprint } from '@backstage/frontend-plugin-api';
import { HomePageRequestedReviewsCard, HomePageYourOpenPullRequestsCard } from '@roadiehq/backstage-plugin-github-pull-requests';
import { ActivityInbox } from './ActivityInbox';
import { CloudPanel, StatusPanel } from './homeModule';

export function Today() {
  const theme = useTheme();
  const dark = theme.palette.type === 'dark';
  const panel = { border: `1px solid ${theme.palette.divider}`, borderRadius: 16, background: theme.palette.background.paper, minWidth: 0, overflow: 'auto' as const };
  return <Box p={2}>
    <Box p={3} mb={3} style={{ ...panel, backgroundImage: `linear-gradient(${dark ? '#111215eb' : '#ffffffed'},${dark ? '#111215eb' : '#ffffffed'}),url(/brand/legacy-background-${dark ? 'dark' : 'light'}.webp)`, backgroundSize: 'cover' }}>
      <Typography variant="overline">ONE MIND · OPERATIONS</Typography>
      <Typography variant="h4" component="h1" style={{ fontWeight: 700 }}>Today</Typography>
      <Typography color="textSecondary">See what needs attention. Review your work. Choose your next move.</Typography>
      <Box mt={2} display="flex" flexWrap="wrap" style={{ gap: 8 }}>
        <Button variant="outlined" href="/">App directory</Button>
        <Button variant="outlined" href="/notifications">Full inbox</Button>
        <Button variant="outlined" href="/create">Prepare a change</Button>
        <Button variant="outlined" href="/docs">Runbooks</Button>
      </Box>
      <Typography variant="caption" display="block" style={{ marginTop: 12 }}>Each source reports its own freshness and coverage. An empty inbox does not mean healthy infrastructure.</Typography>
    </Box>
    <Box style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(min(100%, 380px), 1fr))', gap: 20, alignItems: 'start' }}>
      <Box style={panel}><StatusPanel compact /></Box>
      <Box style={panel}><ActivityInbox /><EventSources /></Box>
      <Box style={panel}><CloudPanel /></Box>
      <Box style={panel}><FabricPanel /></Box>
      <Box style={panel}><Typography variant="caption" component="p" style={{ padding: 16 }}>GitHub authorization required · repository access is restricted to the installed GitHub App.</Typography><HomePageRequestedReviewsCard query="repo:onemindos/omos is:pr is:open review-requested:@me" /></Box>
      <Box style={panel}><HomePageYourOpenPullRequestsCard query="repo:onemindos/omos is:pr is:open author:@me" /></Box>
    </Box>

  </Box>;
}
export const todayPlugin = createFrontendPlugin({ pluginId: 'onemind-today', extensions: [
  PageBlueprint.make({ params: { path: '/today', title: 'Today', icon: <TodayIcon />, loader: async () => <Today /> } }),
] });

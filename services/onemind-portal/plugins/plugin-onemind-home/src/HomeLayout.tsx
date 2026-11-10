import { Fragment, useMemo, useState } from 'react';
import { Box, Button, Typography, useTheme } from '@material-ui/core';
import { Content } from '@backstage/core-components';
import { CustomHomepageGrid } from '@backstage/plugin-home';
import { HomePageLayoutProps } from '@backstage/plugin-home-react/alpha';
import { Directory } from './homeModule';
import { HomeHeader } from './HomeHeader';
import { CommandBar } from './CommandBar';
import { NextUpHero } from './NextUpHero';
import { PlatformScore } from './PlatformScore';
import { WaitingOnYou } from './WaitingOnYou';
import { TasksCard } from './TasksCard';
import { LiveTile } from './LiveTile';
import { ActionTiles } from './ActionTiles';
import { useBentoStyles } from './bento';
import { JumpBackIn } from './JumpBackIn';
import { PinnedChips } from './PinnedChips';
import { BrowseAllApps } from './BrowseAllApps';
import { cardGridStyle } from './HomeCard';
import { PillarTiles } from './PillarTiles';
import { RecentlyChanged } from './RecentlyChanged';
import { ActiveStack } from './ActiveStack';
import { useHomeFeature } from './features';

// Command-center home, top to bottom, with no fixed heights: header, command
// bar, Row A (Next up hero + Platform score), Row B (Waiting on you + Tasks),
// Row C (action tiles), Jump back in / Pinned, recent changes (+ pillars), a
// collapsed app launcher and the active stack. Backstage's native editor and
// saved layout stay available in the optional widget view.
export function HomeLayout({ widgets, defaultConfig }: HomePageLayoutProps) {
  const theme = useTheme();
  const bento = useBentoStyles();
  const pillars = useHomeFeature('pillars');
  const [customize, setCustomize] = useState(false);
  const gridConfig = useMemo(
    () =>
      defaultConfig?.map(item => ({
        component: item.component,
        x: item.column,
        y: item.row,
        width: item.width,
        height: item.height,
        movable: item.movable,
        deletable: item.deletable,
        resizable: item.resizable,
      })),
    [defaultConfig],
  );
  const gap = theme.spacing(3);
  return (
    <Content className="onemind-home-content">
      <Box
        className="onemind-home-layout onemind-workspace"
        display="flex"
        flexDirection="column"
        style={{ gap }}
      >
        <HomeHeader
          action={
            <Button
              variant="outlined"
              color="primary"
              size="small"
              onClick={() => setCustomize(value => !value)}
            >
              {customize ? 'Back to home' : 'Customize widgets'}
            </Button>
          }
        />
        {customize ? (
          <Box>
            <Typography variant="body2" color="textSecondary">
              Your saved widget layout. Use Edit to add, resize or rearrange
              widgets.
            </Typography>
            <CustomHomepageGrid config={gridConfig} preventDuplicateWidgets>
              {widgets.map((item, index) => (
                <Fragment key={item.name ?? index}>{item.component}</Fragment>
              ))}
            </CustomHomepageGrid>
          </Box>
        ) : (
          <>
            <CommandBar />
            <Box className={bento.rowA}>
              <NextUpHero />
              <PlatformScore />
            </Box>
            <LiveTile />
            <Box className={bento.rowB}>
              <WaitingOnYou />
              <TasksCard />
            </Box>
            <ActionTiles />
            <Box style={cardGridStyle()}>
              <JumpBackIn />
              <PinnedChips />
            </Box>
            <Box style={cardGridStyle(360)}>
              {pillars && <PillarTiles />}
              <RecentlyChanged />
            </Box>
            <BrowseAllApps>
              <Directory fluid launcherOnly />
            </BrowseAllApps>
            <ActiveStack />
          </>
        )}
      </Box>
    </Content>
  );
}

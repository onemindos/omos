import { Fragment, useMemo, useState } from 'react';
import { Box, Button, Typography, useTheme } from '@material-ui/core';
import { Content } from '@backstage/core-components';
import { CustomHomepageGrid } from '@backstage/plugin-home';
import { HomePageLayoutProps } from '@backstage/plugin-home-react/alpha';
import { Directory } from './homeModule';
import { QuickAccess } from './Launcher';
import { StatusStrip } from './StatusStrip';
import { PillarTiles } from './PillarTiles';
import { RecentlyChanged } from './RecentlyChanged';
import { ActiveStack } from './ActiveStack';
import { useHomeFeature } from './features';

// Natural page flow, top to bottom, with no fixed heights so there is no blank
// region: status strip, favorites, pillars, recent changes, launcher. Backstage's
// native editor and saved layout stay available in the optional widget view.
export function HomeLayout({ widgets, defaultConfig }: HomePageLayoutProps) {
  const theme = useTheme();
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
        <Box
          className="onemind-home-toolbar"
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          flexWrap="wrap"
          style={{ gap: theme.spacing(1) }}
        >
          <Typography variant="h5" component="h1" style={{ fontWeight: 700 }}>
            Home
          </Typography>
          <Button
            variant="outlined"
            color="primary"
            size="small"
            onClick={() => setCustomize(value => !value)}
          >
            {customize ? 'Back to home' : 'Customize widgets'}
          </Button>
        </Box>
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
            <StatusStrip />
            <QuickAccess dense />
            <Box
              style={{
                display: 'grid',
                gridTemplateColumns:
                  'repeat(auto-fit,minmax(min(100%,360px),1fr))',
                gap,
                alignItems: 'start',
              }}
            >
              {pillars && <PillarTiles />}
              <RecentlyChanged />
            </Box>
            <Directory fluid launcherOnly />
            <ActiveStack />
          </>
        )}
      </Box>
    </Content>
  );
}

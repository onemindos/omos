import { useEffect, useMemo, ReactNode } from 'react';
import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { ThemeBlueprint } from '@backstage/plugin-app-react';
import { UnifiedThemeProvider } from '@backstage/theme';
import LightIcon from '@material-ui/icons/WbSunny';
import DarkIcon from '@material-ui/icons/Brightness4';
import { designPage } from './designPlugin';
import { createOneMindTheme } from './muiTheme';
import { useBackgroundImage } from './preferences';
import { heroImage, ThemeMode } from './tokens';

function OneMindThemeProvider({
  mode,
  children,
}: {
  mode: ThemeMode;
  children: ReactNode;
}) {
  const [background] = useBackgroundImage();
  const theme = useMemo(
    () => createOneMindTheme(mode, { backgroundImage: background }),
    [mode, background],
  );
  useEffect(() => {
    document.body.dataset.onemindBackground = background ? 'on' : 'off';
    document.body.style.setProperty(
      '--onemind-hero-image',
      background ? heroImage(mode) : 'none',
    );
  }, [background, mode]);
  return <UnifiedThemeProvider theme={theme}>{children}</UnifiedThemeProvider>;
}

const themes = (['light', 'dark'] as const).map(mode =>
  ThemeBlueprint.make({
    name: `onemind-${mode}`,
    params: {
      theme: {
        id: mode,
        title: `OneMind ${mode}`,
        variant: mode,
        icon: mode === 'dark' ? <DarkIcon /> : <LightIcon />,
        Provider: ({ children }) => (
          <OneMindThemeProvider mode={mode}>{children}</OneMindThemeProvider>
        ),
      },
    },
  }),
);

export const designModule = createFrontendModule({
  pluginId: 'app',
  extensions: [...themes, designPage],
});

import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { ThemeBlueprint } from '@backstage/plugin-app-react';
import { createBaseThemeOptions, createUnifiedTheme, UnifiedThemeProvider, palettes } from '@backstage/theme';
import LightIcon from '@material-ui/icons/WbSunny';
import DarkIcon from '@material-ui/icons/Brightness4';
const extensions = (['light', 'dark'] as const).map(mode => {
  const dark = mode === 'dark';
  const theme = createUnifiedTheme({
    ...createBaseThemeOptions({ palette: {
      ...palettes[mode],
      primary: { main: dark ? '#73d9e4' : '#087e8b' },
      secondary: { main: dark ? '#bea8ff' : '#6e49ae' },
      background: { default: dark ? '#101822' : '#f1f5fa', paper: dark ? '#192432' : '#ffffff' },
      text: { primary: dark ? '#f4f4f5' : '#202226', secondary: dark ? '#b7bdc7' : '#59616c' },
      navigation: { background: '#121b29', indicator: '#73d9e4', color: '#c6c8ce', selectedColor: '#ffffff' },
    } }),
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    defaultPageTheme: 'home',
    components: {
      MuiPaper: { styleOverrides: { rounded: { borderRadius: 16 } } },
      MuiButton: { styleOverrides: { root: { borderRadius: 10, textTransform: 'none', fontWeight: 600 } } },
      MuiCssBaseline: { styleOverrides: {
        'a:focus-visible, button:focus-visible': { outline: '3px solid #3899d1', outlineOffset: 4 },
      } },
      BackstageHeader: { styleOverrides: { header: {
        backgroundImage: `linear-gradient(${dark ? '#111215e8' : '#ffffffe8'},${dark ? '#111215e8' : '#ffffffe8'}),url(/brand/legacy-background-${mode}.webp)`,
        backgroundSize: 'cover', boxShadow: 'none',
      }, title: { color: dark ? '#f4f4f5' : '#202226' }, subtitle: { color: dark ? '#b7bdc7' : '#59616c' } } },
    },
  });
  return ThemeBlueprint.make({ name: `onemind-${mode}`, params: { theme: {
    id: mode, title: `OneMind ${mode}`, variant: mode,
    icon: dark ? <DarkIcon /> : <LightIcon />,
    Provider: ({ children }) => <UnifiedThemeProvider theme={theme}>{children}</UnifiedThemeProvider>,
  } } });
});
export const designModule = createFrontendModule({ pluginId: 'app', extensions });

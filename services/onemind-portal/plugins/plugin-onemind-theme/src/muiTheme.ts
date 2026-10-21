import {
  createBaseThemeOptions,
  createUnifiedTheme,
  palettes,
  PageTheme,
} from '@backstage/theme';
import { heroImage, ThemeMode, tokens } from './tokens';

/** Touch target floor on phones (px). */
export const MOBILE_TOUCH_TARGET = 40;

const phone = '@media (max-width:599.95px)';

/**
 * The MUI (v4 + v5) theme for one mode. Everything comes from `tokens`; the
 * only raw colours in this file are derived from them.
 */
export function createOneMindTheme(
  mode: ThemeMode,
  options: { backgroundImage?: boolean } = {},
) {
  const t = tokens[mode];
  const s = t.status;
  const dark = mode === 'dark';
  const flatPage: PageTheme = {
    colors: [t.surface[1], t.surface[1]],
    shape: 'none',
    backgroundImage: 'none',
    fontColor: t.text[1],
  };
  const pageTheme = Object.fromEntries(
    [
      'home',
      'documentation',
      'tool',
      'service',
      'website',
      'library',
      'other',
      'app',
      'apis',
      'card',
    ].map(id => [id, flatPage]),
  );
  const headerBackground = options.backgroundImage
    ? `linear-gradient(${t.surface[1]}e8,${t.surface[1]}e8),${heroImage(mode)}`
    : 'none';

  return createUnifiedTheme({
    ...createBaseThemeOptions({
      palette: {
        ...palettes[mode],
        primary: {
          main: t.brand.base,
          dark: t.brand.press,
          light: t.brand.hover,
          contrastText: t.brand.on,
        },
        secondary: { main: t.accent.base, contrastText: t.accent.on },
        error: { main: s.critical.fg },
        warning: { main: s.warn.fg },
        success: { main: s.ok.fg },
        info: { main: s.info.fg },
        background: { default: t.surface[0], paper: t.surface[1] },
        text: { primary: t.text[1], secondary: t.text[2] },
        divider: t.border.subtle,
        action: {
          hover: `${t.text[1]}${dark ? '14' : '0d'}`,
          selected: t.brand.subtle,
          focus: t.focus,
        },
        status: {
          ok: s.ok.fg,
          warning: s.warn.fg,
          error: s.critical.fg,
          pending: s.unknown.fg,
          running: s.info.fg,
          aborted: s.unknown.fg,
        },
        border: t.border.subtle,
        textContrast: t.text[1],
        textVerySubtle: t.border.strong,
        textSubtle: t.text[2],
        highlight: t.brand.subtle,
        errorBackground: s.critical.bg,
        warningBackground: s.warn.bg,
        infoBackground: s.info.bg,
        errorText: s.critical.fg,
        warningText: s.warn.fg,
        infoText: s.info.fg,
        link: t.brand.text,
        linkHover: t.text[1],
        gold: s.warn.fg,
        navigation: {
          background: t.nav.bg,
          indicator: t.brand.base,
          color: t.nav.text,
          selectedColor: t.nav.selected,
          navItem: { hoverBackground: t.nav.hover },
          submenu: { background: t.nav.hover },
        },
        tabbar: { indicator: t.brand.base },
        banner: {
          info: s.info.bg,
          error: s.critical.bg,
          warning: s.warn.bg,
          text: t.text[1],
          link: t.brand.text,
        },
        pinSidebarButton: { icon: t.nav.text, background: t.nav.hover },
      },
      defaultPageTheme: 'home',
      pageTheme,
    }),
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    components: {
      MuiPaper: {
        styleOverrides: {
          root: { backgroundImage: 'none' },
          rounded: { borderRadius: 10 },
          outlined: { borderColor: t.border.subtle },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            border: `1px solid ${t.border.subtle}`,
            boxShadow: 'none',
          },
        },
      },
      MuiCardHeader: { styleOverrides: { root: { padding: '10px 16px' } } },
      MuiCardContent: {
        styleOverrides: {
          root: { padding: '12px 16px', '&:last-child': { paddingBottom: 12 } },
        },
      },
      MuiTableCell: {
        styleOverrides: {
          root: {
            padding: '6px 12px',
            borderBottomColor: t.border.subtle,
          },
          head: {
            fontWeight: 600,
            color: t.text[2],
            backgroundColor: t.surface[1],
            lineHeight: 1.3,
          },
          sizeSmall: { padding: '4px 10px' },
        },
      },
      MuiTableRow: {
        styleOverrides: {
          root: {
            '&:hover': { backgroundColor: t.surface[2] },
            '&.Mui-selected, &.Mui-selected:hover': {
              backgroundColor: t.brand.subtle,
            },
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            textTransform: 'none',
            fontWeight: 600,
            [phone]: { minHeight: MOBILE_TOUCH_TARGET },
          },
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: {
            [phone]: {
              minWidth: MOBILE_TOUCH_TARGET,
              minHeight: MOBILE_TOUCH_TARGET,
            },
          },
        },
      },
      MuiListItem: {
        styleOverrides: {
          root: { [phone]: { minHeight: MOBILE_TOUCH_TARGET } },
        },
      },
      MuiTab: {
        styleOverrides: {
          root: {
            textTransform: 'none',
            fontWeight: 600,
            [phone]: { minHeight: MOBILE_TOUCH_TARGET },
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            borderRadius: 6,
            backgroundColor: t.surface[3],
            color: t.text[1],
            height: 24,
          },
          outlined: {
            borderColor: t.border.strong,
            backgroundColor: 'transparent',
          },
        },
      },
      MuiOutlinedInput: {
        styleOverrides: {
          notchedOutline: { borderColor: t.border.strong },
        },
      },
      // Floating surfaces: raised one step above the page, with a hairline.
      MuiPopover: {
        styleOverrides: {
          paper: {
            backgroundColor: t.surface[2],
            backgroundImage: 'none',
            border: `1px solid ${t.border.subtle}`,
          },
        },
      },
      MuiMenuItem: {
        styleOverrides: {
          root: {
            '&:hover': { backgroundColor: t.surface[3] },
            '&.Mui-selected, &.Mui-selected:hover': {
              backgroundColor: t.brand.subtle,
            },
          },
        },
      },
      MuiDialog: {
        styleOverrides: {
          paper: {
            backgroundColor: t.surface[1],
            backgroundImage: 'none',
            border: `1px solid ${t.border.subtle}`,
          },
        },
      },
      MuiSnackbarContent: {
        styleOverrides: {
          root: {
            backgroundColor: t.surface[2],
            color: t.text[1],
            border: `1px solid ${t.border.subtle}`,
          },
        },
      },
      MuiAlert: {
        styleOverrides: {
          standardError: {
            backgroundColor: s.critical.bg,
            color: s.critical.fg,
          },
          standardWarning: { backgroundColor: s.warn.bg, color: s.warn.fg },
          standardInfo: { backgroundColor: s.info.bg, color: s.info.fg },
          standardSuccess: { backgroundColor: s.ok.bg, color: s.ok.fg },
        },
      },
      MuiTooltip: {
        styleOverrides: {
          tooltip: {
            backgroundColor: t.surface[3],
            color: t.text[1],
            border: `1px solid ${t.border.subtle}`,
          },
        },
      },
      MuiCssBaseline: {
        styleOverrides: {
          // Paint the base surface everywhere, so short pages never show the
          // browser's white canvas below the content.
          'html, body, #root': {
            backgroundColor: t.surface[0],
            minHeight: '100%',
          },
          'a:focus-visible, button:focus-visible, [role="button"]:focus-visible, [tabindex]:focus-visible':
            {
              outline: `3px solid ${t.focus}`,
              outlineOffset: 3,
            },
        },
      },
      BackstageHeader: {
        styleOverrides: {
          header: {
            backgroundColor: t.surface[1],
            backgroundImage: headerBackground,
            backgroundSize: 'cover',
            borderBottom: `1px solid ${t.border.subtle}`,
            boxShadow: 'none',
            padding: '14px 24px',
          },
          title: { color: t.text[1] },
          subtitle: { color: t.text[2] },
        },
      },
    },
  });
}

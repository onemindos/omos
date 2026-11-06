import { useTheme } from '@material-ui/core';
import { tokens, ModeTokens } from '@internal/plugin-onemind-theme/tokens';

/** Design tokens for the active light or dark mode (same source as `brand.css`). */
export function useHomeTokens(): ModeTokens & { dark: boolean } {
  const dark = useTheme().palette.type === 'dark';
  return { ...tokens[dark ? 'dark' : 'light'], dark };
}

/** Appends an alpha channel to a token hex colour (`#rrggbb` + 2 hex digits). */
export const withAlpha = (hex: string, alpha: string) => `${hex}${alpha}`;

import { categoryPalette } from '@internal/plugin-onemind-theme/tokens';

type Colors = { base: string; bright: string };
const c = categoryPalette;

export const groupColors: Record<string, Colors> = {
  Infrastructure: c.blue,
  'AI & Agents': c.violet,
  Communication: c.teal,
  'TAK Ecosystem': c.amber,
  Intelligence: c.orange,
  'Data & Storage': c.indigo,
  Observability: c.cyan,
  Geospatial: c.moss,
  'Zoho Apps': c.plum,
  'Personal Tools': c.indigo,
  Workspaces: c.violet,
  'Personal Finance': c.moss,
  External: c.slate,
};
const appOverrides: Record<string, Colors> = {
  grafana: c.orange,
  github: c.violet,
  'zoho-mail': c.cyan,
  '1password': c.blue,
  mercury: c.indigo,
  cloudflare: c.amber,
  digitalocean: c.blue,
};
export function appColors(name: string, group: string): Colors {
  return appOverrides[name] ?? groupColors[group] ?? groupColors.External;
}

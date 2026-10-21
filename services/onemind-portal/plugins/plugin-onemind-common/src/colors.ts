export const groupColors: Record<string, { base: string; bright: string }> = {
  'Infrastructure': {base:'#1970d4',bright:'#80b8ff'},
  'AI & Agents': {base:'#8250c4',bright:'#c5a2ff'},
  'Communication': {base:'#007f76',bright:'#67ddd0'},
  'TAK Ecosystem': {base:'#a96816',bright:'#ffce83'},
  'Intelligence': {base:'#b65b18',bright:'#ffb880'},
  'Data & Storage': {base:'#5069c5',bright:'#a7b7ff'},
  'Observability': {base:'#237f91',bright:'#77cfe4'},
  'Geospatial': {base:'#358050',bright:'#8bdb9f'},
  'Zoho Apps': {base:'#ae3d56',bright:'#ff99ad'},
  'Personal Tools': {base:'#536bc5',bright:'#a9b9ff'},
  'Workspaces': {base:'#8250c4',bright:'#cba2ff'},
  'Personal Finance': {base:'#27835f',bright:'#82dcb7'},
  'External': {base:'#68758b',bright:'#b8c5dd'},
};
export function appColors(name: string, group: string) {
  return ({
    'grafana': {base:'#b65b18',bright:'#ffb880'},
    'github': {base:'#6941ba',bright:'#baa2ff'},
    'zoho-mail': {base:'#167a9a',bright:'#77d8ef'},
    '1password': {base:'#1d70c1',bright:'#80b8ff'},
    'mercury': {base:'#5a4ea8',bright:'#b9acf7'},
    'cloudflare': {base:'#a96816',bright:'#ffce83'},
    'digitalocean': {base:'#1970d4',bright:'#80b8ff'},
  } as Record<string, {base:string;bright:string}>)[name] ?? groupColors[group] ?? groupColors.External;
}

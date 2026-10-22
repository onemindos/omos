/** "5 m", "3 h", "2 d". Dates in the future or unparsable return an empty string. */
export function age(iso: string | null | undefined, now = Date.now()): string {
  const ms = now - Date.parse(iso ?? '');
  if (!Number.isFinite(ms) || ms < 0) return '';
  const minutes = Math.floor(ms / 60000);
  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes} m`;
  const hours = Math.floor(minutes / 60);
  if (hours < 48) return `${hours} h`;
  return `${Math.floor(hours / 24)} d`;
}

export function clock(iso: string | null | undefined): string {
  if (!iso || Number.isNaN(Date.parse(iso))) return '';
  return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

export const percent = (value: number | null | undefined) =>
  value === null || value === undefined ? 'n/a' : `${value}%`;

export const PILLARS: Record<string, string> = {
  so: 'SO · Sovereign Operations',
  hp: 'HP · Holistic Performance',
  le: 'LE · Legacy Evolution',
  ge: 'GE · Generational Empire',
};
export const UNASSIGNED = 'Unassigned';
export const pillarTitle = (pillar?: string) => (pillar ? PILLARS[pillar] ?? pillar : UNASSIGNED);

// Server-only provider access. No tokens or provider response bodies reach logs/UI.
type Summary = { state: string; source: string; observedAt: string | null; expiresAt: string | null; open?: number; overdue?: number; sampled?: number; partial?: boolean; link: string };
const link = 'https://projects.zoho.com/portal/zeusdelacruz#zp/dashboard/2727104000000135028';
let cache: Summary | undefined;
let nextAttempt = 0;
let pending: Promise<Summary> | undefined;
async function load(): Promise<Summary> {
  const portal = process.env.ZOHO_PROJECTS_PORTAL_ID;
  const client = process.env.ZOHO_CLIENT_ID;
  const secret = process.env.ZOHO_CLIENT_SECRET;
  const refresh = process.env.ZOHO_REFRESH_TOKEN;
  const empty = { state: 'configuration-required', source: 'Zoho Projects V3 (US account)', observedAt: null, expiresAt: null, link };
  if (!portal || !/^\d+$/.test(portal) || !client || !secret || !refresh) return empty;
  try {
    const auth = await fetch('https://accounts.zoho.com/oauth/v2/token', { method: 'POST', redirect: 'error', signal: AbortSignal.timeout(8000), body: new URLSearchParams({ grant_type: 'refresh_token', client_id: client, client_secret: secret, refresh_token: refresh }) });
    if (!auth.ok) throw new Error('Provider unavailable');
    const token = await auth.json();
    if (typeof token.access_token !== 'string') throw new Error('Provider authentication unavailable');
    const response = await fetch(`https://projectsapi.zoho.com/api/v3/portal/${portal}/tasks?page=1&per_page=100`, { redirect: 'error', signal: AbortSignal.timeout(8000), headers: { Authorization: `Zoho-oauthtoken ${token.access_token}` } });
    if (!response.ok) throw new Error('Provider unavailable');
    const data = await response.json();
    if (!Array.isArray(data.tasks)) throw new Error('Unexpected provider schema');
    let open = 0, overdue = 0;
    for (const task of data.tasks) {
      if (typeof task.status?.is_closed_type !== 'boolean') throw new Error('Unknown task status');
      if (!task.status.is_closed_type) {
        open++;
        if (task.due_date && Number.isFinite(Date.parse(task.due_date)) && task.due_date.slice(0,10) < new Date().toISOString().slice(0,10)) overdue++;
      }
    }
    const now = Date.now();
    cache = { state: 'current', source: empty.source, observedAt: new Date(now).toISOString(), expiresAt: new Date(now + 300000).toISOString(), open, overdue, sampled: data.tasks.length, partial: data.page_info?.has_next_page !== false, link };
    return cache;
  } catch {
    return cache ? { ...cache, state: 'stale' } : { ...empty, state: 'unavailable' };
  }
}
let lastResult: Summary | undefined;
export async function cloudSummary(): Promise<Summary> {
  if (pending) return pending;
  if (lastResult && Date.now() < nextAttempt) return lastResult;
  nextAttempt = Date.now() + 300000;
  pending = load().then(result => { lastResult = result; return result; }).finally(() => { pending = undefined; });
  return pending;
}

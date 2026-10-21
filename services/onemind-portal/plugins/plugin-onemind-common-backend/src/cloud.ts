/* eslint-disable */
// Server-only provider access. No tokens or provider response bodies reach logs/UI.
import { CloudSummary as Summary, CloudTask, CloudProvider } from './cloudTypes';
const link = 'https://projects.zoho.com/portal/zeusdelacruz#zp/dashboard/2727104000000135028';
let cache: Summary | undefined;
let nextAttempt = 0;
let pending: Promise<Summary> | undefined;
async function boundedJson(response: Response, limit: number): Promise<any> {
  const reader = response.body?.getReader();
  if (!reader) throw new Error('Empty provider response');
  const chunks: Uint8Array[] = [];
  let size = 0;
  try {
    while (true) {
      const part = await reader.read();
      if (part.done) break;
      size += part.value.byteLength;
      if (size > limit) throw new Error('Provider response exceeds limit');
      chunks.push(part.value);
    }
    return JSON.parse(Buffer.concat(chunks).toString('utf8'));
  } finally { await reader.cancel().catch(() => {}); }
}
function calendarDate(value: unknown): string | null {
  if (typeof value !== 'string') return null;
  // V3 dates can include an ISO timestamp; compare the explicitly recorded calendar date.
  const date = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z$/.test(value) ? value.slice(0,10) : value;
  if (!/^\d{4}-\d{2}-\d{2}$/.test(date)) return null;
  const parsed = Date.parse(`${date}T00:00:00.000Z`);
  return Number.isFinite(parsed) && new Date(parsed).toISOString().slice(0, 10) === date ? date : null;
}
async function load(): Promise<Summary> {
  const portal = process.env.ZOHO_PROJECTS_PORTAL_ID;
  const client = process.env.ZOHO_CLIENT_ID;
  const secret = process.env.ZOHO_CLIENT_SECRET;
  const refresh = process.env.ZOHO_REFRESH_TOKEN;
  const empty: Summary = { providerId: 'zoho-projects', state: 'configuration-required', source: 'Zoho Projects V3 (US account)', observedAt: null, expiresAt: null, attemptedAt: new Date().toISOString(), coverage: 'No provider observation', tasks: [], link };
  if (!portal || !/^\d+$/.test(portal) || !client || !secret || !refresh) return empty;
  try {
    const auth = await fetch('https://accounts.zoho.com/oauth/v2/token', { method: 'POST', redirect: 'error', signal: AbortSignal.timeout(8000), body: new URLSearchParams({ grant_type: 'refresh_token', client_id: client, client_secret: secret, refresh_token: refresh }) });
    if (!auth.ok) throw new Error('Provider unavailable');
    const token = await boundedJson(auth, 32768);
    if (typeof token.access_token !== 'string') throw new Error('Provider authentication unavailable');
    const signal = AbortSignal.timeout(20000);
    const tasks: CloudTask[] = [];
    const seen = new Set<string>();
    let open = 0, overdue = 0, sampled = 0, pages = 0, partial = true;
    const today = new Date().toISOString().slice(0, 10);
    for (let page = 1; page <= 5; page++) {
      const response = await fetch(`https://projectsapi.zoho.com/api/v3/portal/${portal}/tasks?page=${page}&per_page=100`, { redirect: 'error', signal, headers: { Authorization: `Bearer ${token.access_token}` } });
      if (!response.ok) throw new Error('Provider unavailable');
      const data = await boundedJson(response, 2 * 1024 * 1024);
      if (!Array.isArray(data.tasks) || data.tasks.length > 100) throw new Error('Unexpected provider schema');
      pages++;
      for (const task of data.tasks) {
        if (typeof task.status?.is_closed_type !== 'boolean' || !['string', 'number'].includes(typeof task.id)) throw new Error('Unknown task schema');
        const id = String(task.id);
        if (seen.has(id)) continue;
        seen.add(id); sampled++;
        if (!task.status.is_closed_type) {
          open++;
          const dueDate = calendarDate(task.end_date);
          const late = !!dueDate && dueDate < today;
          if (late) overdue++;
          tasks.push({ id, title: typeof task.name === 'string' ? task.name.replace(/[\u0000-\u001f]/g, ' ').slice(0, 200) : 'Untitled task', dueDate, overdue: late });
        }
      }
      if (data.page_info?.has_next_page === false) { partial = false; break; }
      // Unknown pagination stops conservatively instead of guessing completeness.
      if (data.page_info?.has_next_page !== true || data.tasks.length === 0) break;
    }
    tasks.sort((a, b) => Number(b.overdue) - Number(a.overdue) || (a.dueDate || '9999').localeCompare(b.dueDate || '9999') || a.id.localeCompare(b.id));
    const now = Date.now();
    cache = { ...empty, state: 'current', observedAt: new Date(now).toISOString(), expiresAt: new Date(now + 300000).toISOString(), open, overdue, sampled, pages, partial, coverage: partial ? 'Bounded sample, not portal totals' : 'Complete pagination observed; provider visibility applies', tasks: tasks.slice(0, 20), link };
    return cache;
  } catch {
    return cache ? { ...cache, state: 'stale', attemptedAt: empty.attemptedAt } : { ...empty, state: 'unavailable' };
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

export const zohoProjectsProvider: CloudProvider = { id: 'zoho-projects', summary: cloudSummary };

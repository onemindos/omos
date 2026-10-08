import { createHash } from 'node:crypto';
type Summary = {
  state: 'configuration-required' | 'current' | 'stale' | 'unavailable';
  observedAt: string | null;
  expiresAt: string | null;
  attemptedAt: string;
  coverage: string;
  events: {
    id: string;
    title: string;
    start: string | null;
    allDay: boolean;
  }[];
  partial: boolean;
  link: string;
};
let cached: Summary | undefined,
  last: Summary | undefined,
  nextAttempt = 0,
  pending: Promise<Summary> | undefined;
async function json(response: Response, limit: number) {
  if (!response.ok || !response.body) throw new Error('Calendar unavailable');
  const reader = response.body.getReader();
  let size = 0;
  const chunks: Uint8Array[] = [];
  try {
    while (true) {
      const part = await reader.read();
      if (part.done) break;
      size += part.value.byteLength;
      if (size > limit) throw new Error('Calendar response limit');
      chunks.push(part.value);
    }
    return JSON.parse(Buffer.concat(chunks).toString('utf8'));
  } finally {
    await reader.cancel().catch(() => {});
  }
}
function eventStart(event: any): string | null {
  const raw = event.dateandtime?.start ?? event.start;
  if (typeof raw !== 'string') return null;
  if (event.isallday === true && /^\d{8}$/.test(raw))
    return raw.slice(0, 4) + '-' + raw.slice(4, 6) + '-' + raw.slice(6, 8);
  const match =
    /^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})(\d{2})(Z|[+-]\d{4})?$/.exec(raw);
  if (!match) return null;
  const zone =
    match[7] ?? (event.dateandtime?.timezone === 'UTC' ? 'Z' : undefined);
  if (!zone) return null;
  const offset = zone === 'Z' ? zone : zone.slice(0, 3) + ':' + zone.slice(3);
  const value = Date.parse(
    `${match[1]}-${match[2]}-${match[3]}T${match[4]}:${match[5]}:${match[6]}${offset}`,
  );
  return Number.isFinite(value) ? new Date(value).toISOString() : null;
}
export async function calendarSummary(): Promise<Summary> {
  if (pending) return pending;
  if (last && Date.now() < nextAttempt) return last;
  nextAttempt = Date.now() + 300000;
  const empty: Summary = {
    state: 'configuration-required',
    observedAt: null,
    expiresAt: null,
    attemptedAt: new Date().toISOString(),
    coverage:
      'One selected US Zoho calendar; seven-day bounded window, at most 20 events; not all calendars',
    events: [],
    partial: true,
    link: 'https://calendar.zoho.com/',
  };
  pending = (async () => {
    const client = process.env.ZOHO_CLIENT_ID,
      secret = process.env.ZOHO_CLIENT_SECRET,
      refresh = process.env.ZOHO_REFRESH_TOKEN,
      uid = process.env.ZOHO_CALENDAR_UID;
    if (
      !client ||
      !secret ||
      !refresh ||
      !uid ||
      !/^[A-Za-z0-9_+=-]{1,200}$/.test(uid)
    )
      return empty;
    try {
      const token = await json(
        await fetch('https://accounts.zoho.com/oauth/v2/token', {
          method: 'POST',
          redirect: 'error',
          signal: AbortSignal.timeout(8000),
          body: new URLSearchParams({
            grant_type: 'refresh_token',
            client_id: client,
            client_secret: secret,
            refresh_token: refresh,
          }),
        }),
        32768,
      );
      if (typeof token.access_token !== 'string')
        throw new Error('Calendar authorization unavailable');
      const now = Date.now(),
        format = (date: Date) =>
          date
            .toISOString()
            .replace(/[-:]/g, '')
            .replace(/\.\d{3}Z$/, 'Z');
      const range = JSON.stringify({
        start: format(new Date(now)),
        end: format(new Date(now + 7 * 86400000)),
      });
      const query = new URLSearchParams({
        range,
        byinstance: 'true',
        timezone: 'UTC',
      });
      const value = await json(
        await fetch(
          `https://calendar.zoho.com/api/v1/calendars/${encodeURIComponent(
            uid,
          )}/events?${query}`,
          {
            headers: { Authorization: `Zoho-oauthtoken ${token.access_token}` },
            redirect: 'error',
            signal: AbortSignal.timeout(8000),
          },
        ),
        1048576,
      );
      if (!Array.isArray(value.events) || value.events.length > 500)
        throw new Error('Unknown calendar schema');
      const events = value.events
        .map((event: any, index: number) => ({
          id: createHash('sha256')
            .update(String(event.uid ?? index))
            .digest('hex')
            .slice(0, 20),
          title:
            event.isprivate !== false
              ? 'Private event'
              : typeof event.title === 'string'
              ? event.title.replace(/[\u0000-\u001f]/g, ' ').slice(0, 200)
              : 'Untitled event',
          start: eventStart(event),
          allDay: event.isallday === true,
        }))
        .sort((a: any, b: any) =>
          (a.start ?? '9999').localeCompare(b.start ?? '9999'),
        )
        .slice(0, 20);
      cached = {
        ...empty,
        state: 'current',
        observedAt: new Date(now).toISOString(),
        expiresAt: new Date(now + 300000).toISOString(),
        events,
      };
      return cached;
    } catch {
      return cached
        ? { ...cached, state: 'stale' as const, attemptedAt: empty.attemptedAt }
        : { ...empty, state: 'unavailable' as const };
    }
  })()
    .then(value => (last = value))
    .finally(() => {
      pending = undefined;
    });
  return pending;
}

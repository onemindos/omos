import type { Track } from './contract';

/** CoT events larger than this are skipped (images and attachments ride in <detail>). */
export const MAX_COT_BYTES = 64 * 1024;

const ENTITIES: Record<string, string> = {
  '&amp;': '&',
  '&lt;': '<',
  '&gt;': '>',
  '&quot;': '"',
  '&apos;': "'",
};

function decode(value: string): string {
  return value.replace(/&(amp|lt|gt|quot|apos);/g, m => ENTITIES[m] ?? m);
}

function attrs(tag: string): Record<string, string> {
  const out: Record<string, string> = {};
  const re = /([A-Za-z_][\w:.-]*)\s*=\s*(?:"([^"]*)"|'([^']*)')/g;
  for (const m of tag.matchAll(re)) out[m[1]] = decode(m[2] ?? m[3] ?? '');
  return out;
}

function firstTag(
  xml: string,
  name: string,
): Record<string, string> | undefined {
  const m = new RegExp(`<${name}\\b([^>]*)>`).exec(xml);
  return m ? attrs(m[1]) : undefined;
}

function num(value: string | undefined): number | undefined {
  if (value === undefined || value.trim() === '') return undefined;
  const n = Number(value);
  return Number.isFinite(n) ? n : undefined;
}

function iso(value: string | undefined, fallback: Date): string {
  if (value) {
    const t = Date.parse(value);
    if (Number.isFinite(t)) return new Date(t).toISOString();
  }
  return fallback.toISOString();
}

/**
 * Parse one CoT event into a Track. Returns undefined for anything malformed (not an <event>,
 * no uid or type, no usable point). Careful regex rather than a full XML parse: this runs on
 * every tak.cot.> message, the busiest subject on the bus.
 */
export function parseCot(
  xml: string,
  subject: string,
  now: Date = new Date(),
): Track | undefined {
  if (xml.length > MAX_COT_BYTES) return undefined;
  const event = firstTag(xml, 'event');
  if (!event || !event.uid || !event.type) return undefined;
  const point = firstTag(xml, 'point');
  if (!point) return undefined;
  const lat = num(point.lat);
  const lon = num(point.lon);
  if (lat === undefined || lon === undefined) return undefined;
  if (lat < -90 || lat > 90 || lon < -180 || lon > 180) return undefined;
  // CoT uses 9999999 for "unknown" height.
  const haeRaw = num(point.hae);
  const hae =
    haeRaw === undefined || Math.abs(haeRaw) >= 9999999 ? undefined : haeRaw;
  const contact = firstTag(xml, 'contact');
  const track = firstTag(xml, 'track');
  const time = iso(event.time ?? event.start, now);
  const stale = iso(event.stale, new Date(Date.parse(time) + 2 * 60_000));
  const out: Track = {
    uid: event.uid,
    type: event.type,
    lat,
    lon,
    time,
    stale,
    subject,
  };
  if (contact?.callsign) out.callsign = contact.callsign;
  if (hae !== undefined) out.hae = hae;
  const course = num(track?.course);
  const speed = num(track?.speed);
  if (course !== undefined) out.course = course;
  if (speed !== undefined) out.speed = speed;
  return out;
}

/** Affiliation atom of a CoT type (a-f-G-U-C gives "f"); "other" when not an atom type. */
export function affiliation(type: string): string {
  const parts = type.split('-');
  return parts[0] === 'a' && parts[1] ? parts[1] : 'other';
}

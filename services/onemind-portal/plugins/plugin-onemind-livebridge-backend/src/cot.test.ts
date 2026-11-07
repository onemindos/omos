import { affiliation, parseCot } from './cot';

const EUD = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><event version="2.0" uid="ANDROID-abc123" type="a-f-G-U-C" time="2026-11-07T10:00:00.000Z" start="2026-11-07T10:00:00.000Z" stale="2026-11-07T10:05:00.000Z" how="h-g-i-g-o"><point lat="40.7128" lon="-74.0060" hae="12.5" ce="9999999" le="9999999"/><detail><contact endpoint="*:-1:stcp" callsign="ZEUS &amp; CO"/><__group name="Cyan" role="Team Member"/><track course="90.0" speed="1.5"/></detail></event>`;

const AIRCRAFT = `<event version="2.0" uid="ICAO-A1B2C3" type="a-n-A-C-F" time="2026-11-07T10:00:01Z" start="2026-11-07T10:00:01Z" stale="2026-11-07T10:01:01Z" how="m-g"><point lat="40.64" lon="-73.78" hae="10668" ce="10" le="10"/><detail><contact callsign="DAL123"/><track course="270.5" speed="231.4"/></detail></event>`;

describe('parseCot', () => {
  it('parses an EUD with a contact callsign', () => {
    const t = parseCot(EUD, 'tak.cot.out')!;
    expect(t).toMatchObject({
      uid: 'ANDROID-abc123',
      type: 'a-f-G-U-C',
      callsign: 'ZEUS & CO',
      lat: 40.7128,
      lon: -74.006,
      hae: 12.5,
      course: 90,
      speed: 1.5,
      time: '2026-11-07T10:00:00.000Z',
      stale: '2026-11-07T10:05:00.000Z',
      subject: 'tak.cot.out',
    });
  });

  it('parses an aircraft with hae and track', () => {
    const t = parseCot(AIRCRAFT, 'tak.cot.inject')!;
    expect(t.type).toBe('a-n-A-C-F');
    expect(t.callsign).toBe('DAL123');
    expect(t.hae).toBe(10668);
    expect(t.course).toBe(270.5);
    expect(t.speed).toBe(231.4);
  });

  it('treats hae 9999999 as unknown and omits absent fields', () => {
    const t = parseCot(
      '<event uid="x" type="b-m-p-s-p-i" time="2026-11-07T10:00:00Z" stale="2026-11-07T10:10:00Z"><point lat="1" lon="2" hae="9999999"/></event>',
      's',
    )!;
    expect(t.hae).toBeUndefined();
    expect(t.callsign).toBeUndefined();
    expect(t.course).toBeUndefined();
    expect('hae' in t).toBe(false);
  });

  it.each([
    ['not xml', 'hello'],
    ['empty', ''],
    ['no point', '<event uid="x" type="a-f-G"><detail/></event>'],
    ['no uid', '<event type="a-f-G"><point lat="1" lon="2"/></event>'],
    [
      'bad lat',
      '<event uid="x" type="a-f-G"><point lat="abc" lon="2"/></event>',
    ],
    [
      'lat out of range',
      '<event uid="x" type="a-f-G"><point lat="91" lon="2"/></event>',
    ],
    ['truncated', '<event uid="x" type="a-f-G"><poi'],
    [
      'too large',
      `<event uid="x" type="a-f-G"><point lat="1" lon="2"/><detail>${'a'.repeat(
        70000,
      )}</detail></event>`,
    ],
  ])('skips malformed input: %s', (_name, xml) => {
    expect(parseCot(xml, 's')).toBeUndefined();
  });

  it('falls back to now when the times are missing or invalid', () => {
    const now = new Date('2026-11-07T12:00:00Z');
    const t = parseCot(
      '<event uid="x" type="a-u-G" time="nope"><point lat="1" lon="2"/></event>',
      's',
      now,
    )!;
    expect(t.time).toBe('2026-11-07T12:00:00.000Z');
    expect(t.stale).toBe('2026-11-07T12:02:00.000Z');
  });
});

describe('affiliation', () => {
  it('reads the second atom', () => {
    expect(affiliation('a-f-G-U-C')).toBe('f');
    expect(affiliation('a-h-A')).toBe('h');
    expect(affiliation('b-m-p-s-p-i')).toBe('other');
  });
});

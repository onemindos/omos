import { catalogPath, parseAccessPoints, safeAccessUrl } from './access';
test('accepts HTTPS access without allowing credentials or executable schemes', () => {
  expect(safeAccessUrl('https://example.com/path')).toBe(
    'https://example.com/path',
  );
  for (const value of [
    'javascript:alert(1)',
    'file:///Applications/App.app',
    'https://user:secret@example.com',
  ])
    expect(safeAccessUrl(value)).toBeUndefined();
});
test('retains inspected native access and safely parses catalog metadata', () => {
  const point = {
    id: 'mac',
    label: 'Mac app',
    platform: 'macOS',
    method: 'Native app',
    location: 'MacBook',
    evidence: 'Installation inspected',
    targetRef: 'resource:default/mac',
    instructions: 'Open from Applications',
  };
  expect(parseAccessPoints(JSON.stringify([point]))[0]).toMatchObject(point);
  expect(parseAccessPoints('broken')).toEqual([]);
  expect(
    parseAccessPoints(
      JSON.stringify([
        { ...point, url: 'javascript:alert(1)', instructions: { bad: true } },
      ]),
    )[0].url,
  ).toBeUndefined();
  expect(
    parseAccessPoints(
      JSON.stringify([{ ...point, instructions: { bad: true } }]),
    )[0].instructions,
  ).toBeUndefined();
  expect(catalogPath('resource:default/mac')).toBe(
    '/catalog/default/resource/mac',
  );
  expect(catalogPath('https://example.com')).toBeUndefined();
});

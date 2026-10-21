import { isReadOnlyOpenChoreoRequest } from './hostGuard';
import { createUpstreamRouter } from './upstreamRouter';

test.each(['GET', 'HEAD'])('read method %s allowed', method => {
  expect(isReadOnlyOpenChoreoRequest(method, '/deploy')).toBe(true);
});
test('known read-only batch query is allowed', () => {
  expect(isReadOnlyOpenChoreoRequest('POST', '/dashboard/bindings-count')).toBe(
    true,
  );
});
test.each(['POST', 'PUT', 'PATCH', 'DELETE'])(
  'mutating %s is not read-only',
  method => {
    expect(isReadOnlyOpenChoreoRequest(method, '/deploy-release')).toBe(false);
  },
);
test('read-only exception cannot be extended with a suffix', () => {
  expect(
    isReadOnlyOpenChoreoRequest('POST', '/dashboard/bindings-count/other'),
  ).toBe(false);
});
test('reviewed upstream router loads from the exact installed package', async () => {
  const router = await createUpstreamRouter(
    {} as Parameters<typeof createUpstreamRouter>[0],
  );
  expect(typeof router).toBe('function');
});

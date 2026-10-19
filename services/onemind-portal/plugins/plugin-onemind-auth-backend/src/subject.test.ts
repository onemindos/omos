import { signInKnownSubject } from './subject';

test.each([undefined, null, '', ' ', 123])(
  'invalid subject cannot initiate catalog login: %p',
  async value => {
    const ctx = { signInWithCatalogUser: jest.fn() };
    await expect(signInKnownSubject(value, ctx)).rejects.toThrow(
      'subject is missing',
    );
    expect(ctx.signInWithCatalogUser).not.toHaveBeenCalled();
  },
);
test('signed subject maps only to the exact catalog annotation, no fallback', async () => {
  const identity = { token: 'fixture-only' };
  const ctx = { signInWithCatalogUser: jest.fn().mockResolvedValue(identity) };
  const result = await signInKnownSubject('known-subject', ctx);
  expect(ctx.signInWithCatalogUser).toHaveBeenCalledWith({
    annotations: { 'onemind.io/oidc-sub': 'known-subject' },
  });
  expect(result.identity).toBe(identity);
});
test('unknown catalog user fails closed instead of auto-registering', async () => {
  const ctx = {
    signInWithCatalogUser: jest.fn().mockRejectedValue(Error('User not found')),
  };
  await expect(signInKnownSubject('unregistered', ctx)).rejects.toThrow(
    'User not found',
  );
  expect(ctx.signInWithCatalogUser).toHaveBeenCalledTimes(1);
});

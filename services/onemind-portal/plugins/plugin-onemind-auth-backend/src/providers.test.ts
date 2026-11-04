import {
  ONEMIND_OIDC_PROVIDER_ID,
  oneMindSignInResolver,
  registerOneMindProviders,
} from './providers';

describe('OneMind sign-in providers', () => {
  test('registers only the standard oidc provider, never an OpenChoreo provider', () => {
    const registerProvider = jest.fn();
    registerOneMindProviders({ registerProvider });
    expect(registerProvider).toHaveBeenCalledTimes(1);
    expect(registerProvider.mock.calls[0][0].providerId).toBe('oidc');
    expect(ONEMIND_OIDC_PROVIDER_ID).toBe('oidc');
  });

  test('resolver maps the signed subject to the exact catalog annotation and calls nothing else', async () => {
    const fetchSpy = jest.spyOn(global, 'fetch');
    const identity = { token: 'fixture-only' };
    const ctx = {
      signInWithCatalogUser: jest.fn().mockResolvedValue(identity),
    };
    const result = await oneMindSignInResolver(
      { result: { fullProfile: { userinfo: { sub: 'signed-subject' } } } },
      ctx,
    );
    expect(result).toBe(identity);
    expect(ctx.signInWithCatalogUser).toHaveBeenCalledWith({
      annotations: { 'onemind.io/oidc-sub': 'signed-subject' },
    });
    expect(fetchSpy).not.toHaveBeenCalled();
    fetchSpy.mockRestore();
  });

  test('a missing subject fails closed before any catalog lookup', async () => {
    const ctx = { signInWithCatalogUser: jest.fn() };
    await expect(
      oneMindSignInResolver({ result: { fullProfile: { userinfo: {} } } }, ctx),
    ).rejects.toThrow('subject is missing');
    expect(ctx.signInWithCatalogUser).not.toHaveBeenCalled();
  });
});

import { readdirSync, readFileSync } from 'fs';
import { join } from 'path';
import { oneMindAuthApiRef, ONEMIND_AUTH_PROVIDER_ID } from './authApi';

test('sign-in uses the standard oidc provider the backend registers', () => {
  expect(ONEMIND_AUTH_PROVIDER_ID).toBe('oidc');
  expect(oneMindAuthApiRef.id).toBe('internal.auth.onemind-oidc');
});

test('the auth plugin has no OpenChoreo sign-in, token header or permission API', () => {
  const sources = readdirSync(__dirname)
    .filter(f => /\.(ts|tsx)$/.test(f) && !/\.test\./.test(f))
    .map(f => readFileSync(join(__dirname, f), 'utf8'))
    .join('\n');
  for (const forbidden of [
    'openChoreoAuthApiRef',
    'openchoreo-auth',
    'OpenChoreoPermissionApi',
    'x-openchoreo',
  ]) {
    expect(sources).not.toContain(forbidden);
  }
});

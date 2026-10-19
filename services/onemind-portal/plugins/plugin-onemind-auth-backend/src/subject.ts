import type { AuthResolverContext } from '@backstage/plugin-auth-node';

export function knownSubjectQuery(subject: unknown) {
  if (typeof subject !== 'string' || !subject.trim())
    throw new Error('OIDC subject is missing');
  return { annotations: { 'onemind.io/oidc-sub': subject } };
}
export async function signInKnownSubject(
  subject: unknown,
  ctx: Pick<AuthResolverContext, 'signInWithCatalogUser'>,
) {
  const query = knownSubjectQuery(subject);
  return { query, identity: await ctx.signInWithCatalogUser(query) };
}

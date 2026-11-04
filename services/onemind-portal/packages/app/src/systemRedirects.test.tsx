import { screen, waitFor } from '@testing-library/react';
import { PageBlueprint } from '@backstage/frontend-plugin-api';
import { renderTestApp } from '@backstage/frontend-test-utils';
import redirects from './systemRedirects.json';
import { systemRedirectExtensions } from './systemRedirects';

// Stand-in for the catalog entity page, mounted on the real catalog pattern.
const entityPage = PageBlueprint.make({
  name: 'entity-stand-in',
  params: {
    path: '/catalog/:namespace/:kind/:name',
    loader: async () => {
      const { useParams } = await import('react-router-dom');
      const Page = () => {
        const p = useParams();
        return <div>entity {p.kind}/{p.name}</div>;
      };
      return <Page />;
    },
  },
});

const [from, to] = Object.entries(redirects as Record<string, string>)[0];

test('one redirect per removed System, each pointing at a different name', () => {
  expect(systemRedirectExtensions).toHaveLength(Object.keys(redirects).length);
  for (const [a, b] of Object.entries(redirects)) expect(a).not.toBe(b);
});

test('an old System URL forwards to its successor', async () => {
  renderTestApp({
    extensions: [entityPage, ...systemRedirectExtensions],
    initialRouteEntries: [`/catalog/default/system/${from}`],
  });
  await waitFor(() =>
    expect(screen.getByText(`entity system/${to}`)).toBeInTheDocument(),
  );
});

test('a live System or a Component is never intercepted', async () => {
  renderTestApp({
    extensions: [entityPage, ...systemRedirectExtensions],
    initialRouteEntries: ['/catalog/default/system/infrastructure'],
  });
  await waitFor(() =>
    expect(screen.getByText('entity system/infrastructure')).toBeInTheDocument(),
  );
});

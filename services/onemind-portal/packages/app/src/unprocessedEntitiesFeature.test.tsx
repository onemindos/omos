import { screen } from '@testing-library/react';
import { renderTestApp } from '@backstage/frontend-test-utils';
import stock from '@backstage/plugin-catalog-unprocessed-entities/alpha';
import { unprocessedEntitiesFeature } from './unprocessedEntitiesFeature';

const heading = () => screen.queryAllByText('Unprocessed Entities');

it('serves the unprocessed-entities page on its own route', async () => {
  renderTestApp({
    features: [unprocessedEntitiesFeature],
    initialRouteEntries: ['/catalog-unprocessed-entities'],
  });
  expect(
    (await screen.findAllByText(/Unprocessed Entities/i)).length,
  ).toBeGreaterThan(0);
});

it('control: the stock plugin leaves that route disabled', async () => {
  renderTestApp({
    features: [stock],
    initialRouteEntries: ['/catalog-unprocessed-entities'],
  });
  // Nothing mounts at the route, so the app falls through to its not-found page.
  expect(await screen.findByText(/not found|404/i)).toBeTruthy();
  expect(heading()).toHaveLength(0);
});

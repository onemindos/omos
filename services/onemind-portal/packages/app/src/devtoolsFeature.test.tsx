import { screen } from '@testing-library/react';
import { renderTestApp } from '@backstage/frontend-test-utils';
import { devtoolsFeature } from './devtoolsFeature';

it('offers the Info view and no Config or Scheduled Tasks view', async () => {
  renderTestApp({
    features: [devtoolsFeature],
    initialRouteEntries: ['/devtools'],
  });
  expect(await screen.findByRole('tab', { name: 'Info' })).toBeTruthy();
  expect(screen.queryByRole('tab', { name: 'Config' })).toBeNull();
  expect(screen.queryByRole('tab', { name: 'Scheduled Tasks' })).toBeNull();
});

it('control: the stock plugin shows them, so the override is what removes them', async () => {
  const { default: stock } = await import('@backstage/plugin-devtools/alpha');
  renderTestApp({ features: [stock], initialRouteEntries: ['/devtools'] });
  expect(await screen.findByRole('tab', { name: 'Config' })).toBeTruthy();
  expect(screen.getByRole('tab', { name: 'Scheduled Tasks' })).toBeTruthy();
});

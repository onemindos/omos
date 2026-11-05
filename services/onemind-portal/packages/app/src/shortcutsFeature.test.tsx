import { screen, waitFor } from '@testing-library/react';
import { PageBlueprint, useApi } from '@backstage/frontend-plugin-api';
import { renderTestApp } from '@backstage/frontend-test-utils';
import { shortcutsApiRef } from '@backstage-community/plugin-shortcuts';
import { useEffect, useState } from 'react';
import { shortcutsFeature } from './shortcutsFeature';

const probe = PageBlueprint.make({
  name: 'shortcuts-probe',
  params: {
    path: '/',
    loader: async () => {
      const Probe = () => {
        const api = useApi(shortcutsApiRef);
        const [titles, setTitles] = useState<string[]>([]);
        useEffect(() => {
          api
            .add({ title: 'Cost', url: '/cost' })
            .then(() => setTitles(api.get().map(item => item.title)));
        }, [api]);
        return <div>shortcuts: {titles.join(',') || 'none'}</div>;
      };
      return <Probe />;
    },
  },
});

it('registers the shortcuts API so the sidebar component can store a shortcut', async () => {
  renderTestApp({ features: [shortcutsFeature], extensions: [probe] });
  await waitFor(() => expect(screen.getByText('shortcuts: Cost')).toBeTruthy());
});

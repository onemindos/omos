import { PageBlueprint } from '@backstage/frontend-plugin-api';

/**
 * `/design`. No title or icon on purpose: it is an admin/review page reached
 * by URL, so it does not clutter the sidebar.
 */
export const designPage = PageBlueprint.make({
  name: 'design',
  params: {
    path: '/design',
    loader: () => import('./DesignPage').then(m => <m.DesignPage />),
  },
});

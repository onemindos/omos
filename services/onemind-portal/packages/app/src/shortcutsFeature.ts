import { shortcutsPlugin } from '@backstage-community/plugin-shortcuts';
import { convertLegacyPlugin } from '@backstage/core-compat-api';

/**
 * @backstage-community/plugin-shortcuts has no new-frontend-system export, so the
 * legacy plugin is converted. That registers its API factory (shortcutsApiRef over
 * storageApi.forBucket('shortcuts')); the component itself is rendered by the OneMind
 * sidebar. storageApi is the account-backed one from plugin-app-module-user-settings,
 * so shortcuts follow the user across browsers and are never shared with other users.
 */
export const shortcutsFeature = convertLegacyPlugin(shortcutsPlugin, {
  extensions: [],
});

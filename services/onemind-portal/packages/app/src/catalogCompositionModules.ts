import { createFrontendModule } from '@backstage/frontend-plugin-api';
import {
  catalogPillarFilterExtension,
  catalogViewFilterExtension,
} from './catalogViewFilters';

// Catalog list filters for the OneMind pillars and saved views. The catalog's own About
// and Relations cards are used unchanged.
export const catalogCompositionModule = createFrontendModule({
  pluginId: 'catalog',
  extensions: [catalogViewFilterExtension, catalogPillarFilterExtension],
});

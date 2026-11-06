export { digitalOceanPlugin as default } from './plugin';
export { catalogModuleOnemindDigitalOcean } from './module';
export { registerDigitalOcean } from './module';
export { DigitalOceanService, sharedService } from './service';
export { DigitalOceanClient, DoApiError } from './client';
export { DigitalOceanProvider } from './provider';
export { buildEntities, entityNameFor } from './resources';
export { linkVolumes, ownerOf, listPvcs } from './links';
export type { Summary } from './types';

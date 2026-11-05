import unprocessedEntitiesPlugin from '@backstage/plugin-catalog-unprocessed-entities/alpha';

/**
 * Failed and pending catalog entities, for administrators.
 *
 * The plugin ships its standalone page disabled and attaches a copy of the view to
 * the DevTools page. This portal turns the standalone page on, so the sidebar's Admin
 * group can link to it, and turns the DevTools copy off, so the view is not offered
 * twice and does not depend on DevTools being installed. The backend routes are the
 * boundary: the OneMind policy allows the read permission to administrators only.
 */
export const unprocessedEntitiesFeature =
  unprocessedEntitiesPlugin.withOverrides({
    extensions: [
      unprocessedEntitiesPlugin
        .getExtension('page:catalog-unprocessed-entities')
        .override({ disabled: false }),
      unprocessedEntitiesPlugin
        .getExtension('sub-page:catalog-unprocessed-entities')
        .override({ disabled: true }),
    ],
  });

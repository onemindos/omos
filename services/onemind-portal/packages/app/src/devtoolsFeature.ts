import devtoolsPlugin from '@backstage/plugin-devtools/alpha';

/**
 * DevTools for administrators: the Info view (host, Node, Backstage version and the
 * installed @backstage and @internal package versions).
 *
 * The Config and Scheduled Tasks views are switched off in code rather than in
 * app-config, so they cannot appear if a config layer is lost or an older image runs.
 * The backend closes the same permissions for everyone (see
 * plugin-onemind-permissions-backend): the config view only hides values a plugin's
 * schema marks secret, and the scheduler view can trigger background jobs.
 */
export const devtoolsFeature = devtoolsPlugin.withOverrides({
  extensions: [
    devtoolsPlugin
      .getExtension('sub-page:devtools/config')
      .override({ disabled: true }),
    devtoolsPlugin
      .getExtension('sub-page:devtools/scheduled-tasks')
      .override({ disabled: true }),
  ],
});

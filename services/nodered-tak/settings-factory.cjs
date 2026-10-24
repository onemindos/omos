'use strict';

module.exports = function settings(env) {
  const required = name => {
    if (!env[name]?.trim()) throw new Error(`Required secret is missing: ${name}`);
    return env[name];
  };
  const credentialSecret = required('NODE_RED_CREDENTIAL_SECRET');
  if (credentialSecret.length < 32) throw new Error('Credential encryption key must contain at least 32 characters');
  const username = required('NODE_RED_ADMIN_USERNAME');
  const password = required('NODE_RED_ADMIN_PASSWORD_HASH');
  if (!/^\$2[aby]\$(1[0-9]|2[0-9]|3[01])\$[./A-Za-z0-9]{53}$/.test(password)) {
    throw new Error('Editor password must be a bcrypt hash with cost >= 10');
  }
  return {
    uiHost: '0.0.0.0',
    uiPort: 1880,
    userDir: '/data',
    flowFile: 'flows.json',
    flowFilePretty: true,
    credentialSecret,
    httpAdminRoot: '/',
    // Editor authentication does not protect HTTP In flows. Keep those disabled.
    httpNodeRoot: false,
    adminAuth: {
      type: 'credentials',
      users: [{ username, password, permissions: '*' }],
      sessionExpiryTime: 3600,
    },
    externalModules: {
      autoInstall: false,
      palette: { allowInstall: false, allowUpload: false },
      modules: { allowInstall: false },
    },
    functionExternalModules: false,
    contextStorage: { default: { module: 'localfilesystem' } },
    editorTheme: {
      page: { title: 'OneMind Node-RED · TAK' },
      projects: { enabled: false },
      palette: { editable: false },
    },
    diagnostics: { enabled: false, ui: false },
    runtimeState: { enabled: false, ui: false },
    logging: { console: { level: 'info', metrics: false, audit: false } },
  };
};

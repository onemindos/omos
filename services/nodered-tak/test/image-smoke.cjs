'use strict';
// Run only on a local/CI Docker daemon. No external service or user credential is used.
const assert = require('node:assert/strict');
const crypto = require('node:crypto');
const { execFileSync } = require('node:child_process');
const bcrypt = require('bcryptjs');
const password = crypto.randomBytes(24).toString('base64url');
const env = {
  ...process.env,
  NODE_RED_CREDENTIAL_SECRET: crypto.randomBytes(32).toString('base64url'),
  NODE_RED_ADMIN_USERNAME: 'image-smoke',
  NODE_RED_ADMIN_PASSWORD_HASH: bcrypt.hashSync(password, 10),
};
const docker = args => execFileSync('docker', args, { env, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();
const name = `nodered-smoke-${crypto.randomBytes(6).toString('hex')}`;
(async () => {
  try {
    docker(['run', '-d', '--name', name, '--read-only', '--cap-drop=ALL', '--security-opt=no-new-privileges',
      '--tmpfs', '/data:rw,uid=1000,gid=1000,mode=0700', '--tmpfs', '/tmp:rw,uid=1000,gid=1000,mode=0700',
      '-e', 'NODE_RED_CREDENTIAL_SECRET', '-e', 'NODE_RED_ADMIN_USERNAME', '-e', 'NODE_RED_ADMIN_PASSWORD_HASH',
      '-p', '127.0.0.1::1880', 'onemind-nodered-tak:smoke']);
    const port = docker(['port', name, '1880/tcp']).split(':').pop();
    const base = `http://127.0.0.1:${port}`;
    for (let attempt = 0; attempt < 120; attempt++) {
      try { if ((await fetch(base)).status === 200) break; } catch {}
      await new Promise(resolve => setTimeout(resolve, 250));
    }
    assert.equal((await fetch(`${base}/flows`)).status, 401);
    assert.equal((await fetch(`${base}/settings`)).status, 401);
    const response = await fetch(`${base}/auth/token`, {
      method: 'POST',
      body: new URLSearchParams({ client_id: 'node-red-editor', grant_type: 'password', scope: '*', username: env.NODE_RED_ADMIN_USERNAME, password }),
    });
    assert.equal(response.status, 200);
    const token = (await response.json()).access_token;
    const nodes = await fetch(`${base}/nodes`, { headers: { Authorization: `Bearer ${token}`, Accept: 'application/json' } });
    assert.ok((await nodes.json()).some(module => module.module === 'node-red-contrib-tak'));
    assert.match(docker(['logs', name]), /safe mode/i);
    console.log('Container: read-only filesystem, authenticated editor, TAK palette and safe mode passed');
  } finally {
    try { docker(['rm', '-f', name]); } catch {}
  }
})().catch(() => { console.error('Container smoke test failed; do not publish or activate'); process.exitCode = 1; });

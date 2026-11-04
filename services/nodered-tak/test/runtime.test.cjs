'use strict';
const test = require('node:test');
const assert = require('node:assert/strict');
const crypto = require('node:crypto');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const net = require('node:net');
const { spawn } = require('node:child_process');
const { once } = require('node:events');
const bcrypt = require('bcryptjs');
const makeSettings = require('../settings-factory.cjs');
const bootstrap = require('../bootstrap.cjs');

const password = crypto.randomBytes(24).toString('base64url');
const env = {
  NODE_RED_CREDENTIAL_SECRET: crypto.randomBytes(32).toString('base64url'),
  NODE_RED_ADMIN_USERNAME: 'smoke-test',
  NODE_RED_ADMIN_PASSWORD_HASH: bcrypt.hashSync(password, 10),
};
const root = path.resolve(__dirname, '..');

test('fail closed without a persistent key or authenticated editor; HTTP nodes and installs disabled', () => {
  for (const key of Object.keys(env)) {
    assert.throws(() => makeSettings({ ...env, [key]: '' }), /Required secret/);
  }
  assert.throws(() => makeSettings({ ...env, NODE_RED_ADMIN_PASSWORD_HASH: password }), /bcrypt/);
  assert.throws(() => makeSettings({ ...env, NODE_RED_CREDENTIAL_SECRET: 'short' }), /32 characters/);
  const settings = makeSettings(env);
  assert.equal(settings.adminAuth.default, undefined);
  assert.equal(settings.adminAuth.users.length, 1);
  assert.equal(settings.httpNodeRoot, false);
  assert.equal(settings.externalModules.palette.allowInstall, false);
  assert.equal(settings.externalModules.modules.allowInstall, false);
});

test('bootstrap retains flow bytes and refuses to replace unreviewed modules', () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'nodered-bootstrap-'));
  try {
    bootstrap(dir, path.join(root, 'node_modules'));
    assert.deepEqual(JSON.parse(fs.readFileSync(path.join(dir, 'flows.json'))), []);
    const original = '[{"id":"retained","type":"tab","disabled":true}]\n';
    fs.writeFileSync(path.join(dir, 'flows.json'), original);
    bootstrap(dir, path.join(root, 'node_modules'));
    assert.equal(fs.readFileSync(path.join(dir, 'flows.json'), 'utf8'), original);
    assert.throws(() => bootstrap(dir, '/different-review'), /Unreviewed TAK module/);
  } finally { fs.rmSync(dir, { recursive: true, force: true }); }
});

test('pinned TAK node converts synthetic CoT to JSON and both binary formats without a transport', () => {
  const handlers = {};
  const types = {};
  let sent;
  const errors = [];
  const RED = {
    nodes: {
      registerType: (name, ctor) => { types[name] = ctor; },
      createNode: node => {
        node.on = (event, fn) => { handlers[event] = fn; };
        node.status = () => {};
        node.error = error => errors.push(error);
        node.send = messages => { sent = messages; };
      },
    },
    util: { getMessageProperty: (msg, property) => msg[property] },
  };
  require('node-red-contrib-tak/tak/tak')(RED);
  new types.tak({ property: 'payload' });
  const example = require('../examples/offline-tak.json');
  assert.equal(example[0].disabled, true);
  assert.equal(example[1].once, false);
  assert.equal(example.some(node => /tcp|mqtt|http|nats/.test(node.type)), false);
  handlers.input({ payload: example[1].payload });
  assert.deepEqual(errors, []);
  assert.equal(sent[0].payload.event._attributes.uid, 'onemind-offline-example');
  assert.equal(Buffer.isBuffer(sent[1].payload), true);
  assert.equal(Buffer.isBuffer(sent[2].payload), true);
});

test('real Node-RED starts safe, loads TAK, denies unauthenticated APIs and permits authenticated editor', { timeout: 30000 }, async () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'nodered-runtime-'));
  const reservation = net.createServer();
  reservation.listen(0, '127.0.0.1');
  await once(reservation, 'listening');
  const port = reservation.address().port;
  await new Promise(resolve => reservation.close(resolve));
  bootstrap(dir, path.join(root, 'node_modules'));
  // This test-only wrapper contains paths, never credentials. Production settings remain immutable.
  fs.writeFileSync(path.join(dir, 'settings.cjs'), `module.exports={...require(${JSON.stringify(path.join(root, 'settings-factory.cjs'))})(process.env),userDir:process.env.TEST_USER_DIR,uiHost:'127.0.0.1',uiPort:Number(process.env.TEST_PORT)};\n`);
  const child = spawn(process.execPath, [require.resolve('node-red/red.js'), '--settings', path.join(dir, 'settings.cjs'), '--safe'], {
    env: { ...process.env, ...env, TEST_USER_DIR: dir, TEST_PORT: String(port) },
    stdio: ['ignore', 'pipe', 'pipe'],
  });
  let log = '';
  child.stdout.on('data', data => { log += data; });
  child.stderr.on('data', data => { log += data; });
  const base = `http://127.0.0.1:${port}`;
  try {
    for (let attempt = 0; attempt < 100; attempt++) {
      if (child.exitCode !== null) throw new Error('Node-RED exited before readiness');
      try { if ((await fetch(base)).status === 200) break; } catch {}
      await new Promise(resolve => setTimeout(resolve, 100));
    }
    assert.equal((await fetch(`${base}/settings`)).status, 401);
    assert.equal((await fetch(`${base}/flows`)).status, 401);
    const login = await fetch(`${base}/auth/token`, {
      method: 'POST',
      body: new URLSearchParams({ client_id: 'node-red-editor', grant_type: 'password', scope: '*', username: env.NODE_RED_ADMIN_USERNAME, password }),
    });
    assert.equal(login.status, 200);
    const token = (await login.json()).access_token;
    assert.ok(token);
    const headers = { Authorization: `Bearer ${token}`, Accept: 'application/json' };
    const nodes = await fetch(`${base}/nodes`, { headers });
    assert.equal(nodes.status, 200);
    assert.ok((await nodes.json()).some(module => module.module === 'node-red-contrib-tak'));
    const flows = await fetch(`${base}/flows`, { headers });
    assert.deepEqual(await flows.json(), []);
    assert.match(log, /safe mode/i);
    assert.doesNotMatch(log, /Started flows/);
  } finally {
    child.kill('SIGTERM');
    await once(child, 'exit');
    fs.rmSync(dir, { recursive: true, force: true });
  }
});

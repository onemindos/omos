const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');

const root = path.join(__dirname, '..');
const read = rel => fs.readFileSync(path.join(root, rel), 'utf8');
const files = [
  'plugins/plugin-onemind-common/directory.json',
  'plugins/plugin-onemind-common-backend/directory.json',
];
const directory = JSON.parse(read(files[0]));

// Hosts observed on live Ingress status.loadBalancer and tailscale.com/hostname
// Services (kubectl, 2026-11-07). A tailnet launcher link must use one of these.
// Apps that run without a web UI link to their catalog page on the portal host.
const LIVE_TAILNET_HOSTS = new Set(
  [
    'agent-console',
    'agent-tap',
    'argocd',
    'onemind-auth',
    'onemind-preview',
    'portal',
    'cloudtak',
    'cloudtak-tiles',
    'forgejo',
    'status',
    'onemind-hermes',
    'onemind-hermes-web',
    'karakeep',
    'hubble',
    'headlamp',
    'element',
    'matrix',
    'grafana',
    'n8n',
    'heynats-1',
    'nats-dashboard-1',
    'nats-monitor',
    'nats-ui-1',
    'nats-ws',
    'onemind-langfuse',
    'enroll',
    'opencost',
    'openremote',
    'stac',
    'stac-browser',
    'thingsboard',
    'vaultwarden',
    'fedhub',
    'tak',
    'omos-nats',
  ].map(h => `${h}.tail717ba0.ts.net`),
);

test('the two directory.json copies are identical', () => {
  assert.equal(read(files[1]), read(files[0]));
});

test('entries are well formed and unique', () => {
  const names = new Set();
  const urls = new Set();
  for (const e of directory) {
    for (const key of [
      'name',
      'title',
      'url',
      'description',
      'group',
      'groupTitle',
    ]) {
      assert.equal(typeof e[key], 'string', `${e.name}.${key}`);
    }
    assert.ok(!names.has(e.name), `duplicate name ${e.name}`);
    names.add(e.name);
    const url = new URL(e.url);
    assert.equal(url.protocol, 'https:', e.name);
    assert.equal(url.username + url.password, '', e.name);
    assert.ok(!urls.has(url.href), `duplicate url ${url.href} (${e.name})`);
    urls.add(url.href);
  }
});

test('one groupTitle per group slug', () => {
  const titles = new Map();
  for (const e of directory) {
    assert.equal(titles.get(e.group) ?? e.groupTitle, e.groupTitle, e.group);
    titles.set(e.group, e.groupTitle);
  }
});

test('tailnet links point at hosts that exist in the live cluster', () => {
  for (const e of directory) {
    const { hostname } = new URL(e.url);
    if (hostname.endsWith('.ts.net')) {
      assert.ok(
        LIVE_TAILNET_HOSTS.has(hostname),
        `${e.name}: ${hostname} has no live Ingress or Service`,
      );
    }
  }
});

test('links to hosts or apps that do not exist stay removed', () => {
  const banned = [
    'openclaw',
    'longhorn',
    'minio',
    'external-argocd',
    'external-tailscale',
    'tailscale-admin',
  ];
  for (const name of banned)
    assert.ok(!directory.some(e => e.name === name), name);
  for (const e of directory) {
    const { hostname } = new URL(e.url);
    assert.ok(
      !/^(atoc|webodm|onetak)\.onemindos\.dev$/.test(hostname),
      `${e.name}: ${hostname}`,
    );
    assert.ok(
      !/^(clickhouse|prometheus|agent-web-ui-1|agenttap-1|nats-surveyor-1)\./.test(
        hostname,
      ),
      `${e.name}: ${hostname}`,
    );
  }
});

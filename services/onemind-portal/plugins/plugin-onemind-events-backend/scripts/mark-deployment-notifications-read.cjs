#!/usr/bin/env node
/* eslint-disable */
// Backlog cleanup for the old "Deployment completed: <app>" inbox items.
//
// Marks the signed-in administrator's own unread deployment-success notifications as READ.
// It never deletes anything, never touches other users, and never marks failures read.
// Dry run by default. Not run by CI or by the portal; an administrator runs it by hand.
//
//   PORTAL_URL=https://<portal host> PORTAL_TOKEN=<your Backstage user token> \
//     node mark-deployment-notifications-read.cjs                 # dry run, prints counts
//   ... node mark-deployment-notifications-read.cjs --apply --expect=<matched count from the dry run>
//
// PORTAL_TOKEN is the Bearer token your browser sends to /api (copy it from the network tab).
// It is read from the environment only and is never printed.
'use strict';

const PAGE = 100;

// Keep in sync with isDeploymentSuccessNotification in
// plugin-onemind-common-backend/src/notificationPolicy.ts (a unit test checks parity).
function isDeploymentSuccess(payload) {
  return (
    payload?.topic === 'deployments' &&
    payload?.severity === 'low' &&
    typeof payload?.title === 'string' &&
    payload.title.startsWith('Deployment completed:')
  );
}

function parseArgs(argv) {
  const args = { apply: false, expect: undefined, max: 20000 };
  for (const arg of argv) {
    if (arg === '--apply') args.apply = true;
    else if (arg.startsWith('--expect=')) args.expect = Number(arg.slice(9));
    else if (arg.startsWith('--max=')) args.max = Number(arg.slice(6));
    else throw new Error(`Unknown argument: ${arg}`);
  }
  return args;
}

function portalBase(env) {
  const raw = env.PORTAL_URL;
  if (!raw) throw new Error('Set PORTAL_URL to the portal base URL');
  const url = new URL(raw);
  const local = ['localhost', '127.0.0.1', '[::1]'].includes(url.hostname);
  if (url.protocol !== 'https:' && !local) throw new Error('PORTAL_URL must be https (or localhost)');
  return url.origin;
}

async function run({ fetch: doFetch = fetch, env = process.env, argv = process.argv.slice(2), out = console.log } = {}) {
  const args = parseArgs(argv);
  const base = portalBase(env);
  const token = env.PORTAL_TOKEN;
  if (!token) throw new Error('Set PORTAL_TOKEN to your Backstage user token');
  const headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };

  // Pass 1: read-only. Collect ids of matching unread items.
  const matched = [];
  const kept = [];
  let unreadInTopic = 0;
  for (let offset = 0; ; offset += PAGE) {
    const url = `${base}/api/notifications/notifications?read=false&topic=deployments&limit=${PAGE}&offset=${offset}`;
    const response = await doFetch(url, { headers });
    if (!response.ok) throw new Error(`Listing failed with HTTP ${response.status}`);
    const page = await response.json();
    const items = page.notifications ?? [];
    unreadInTopic += items.length;
    for (const item of items) (isDeploymentSuccess(item.payload) ? matched : kept).push(item);
    if (items.length < PAGE || unreadInTopic >= args.max) break;
  }
  out(`Unread notifications in topic "deployments": ${unreadInTopic}`);
  out(`Matched as old deployment-success items (would be marked read): ${matched.length}`);
  out(`Kept unread (failures, anything else): ${kept.length}`);
  for (const item of kept.slice(0, 5)) out(`  keep: ${String(item.payload?.title).slice(0, 100)} [${item.payload?.severity}]`);

  if (!args.apply) {
    out(`Dry run only. To apply: add --apply --expect=${matched.length}`);
    return { matched: matched.length, kept: kept.length, applied: 0 };
  }
  if (args.expect !== matched.length) {
    throw new Error(`Refusing to apply: --expect=${args.expect} does not equal the matched count ${matched.length}. Re-run the dry run.`);
  }
  // Pass 2: mark read in chunks. read:true only; there is no delete code path in this script.
  let applied = 0;
  for (let i = 0; i < matched.length; i += PAGE) {
    const ids = matched.slice(i, i + PAGE).map(item => item.id);
    const response = await doFetch(`${base}/api/notifications/update`, { method: 'POST', headers, body: JSON.stringify({ ids, read: true }) });
    if (!response.ok) throw new Error(`Update failed with HTTP ${response.status} after ${applied} items`);
    applied += ids.length;
  }
  out(`Marked read: ${applied}. Nothing was deleted.`);
  return { matched: matched.length, kept: kept.length, applied };
}

module.exports = { run, isDeploymentSuccess, parseArgs };

if (require.main === module) {
  run().catch(error => {
    console.error(error.message);
    process.exit(1);
  });
}

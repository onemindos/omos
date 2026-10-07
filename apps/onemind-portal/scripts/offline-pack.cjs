'use strict';
const fs = require('node:fs');
const path = require('node:path');
const escape = value => String(value).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
const names = ['index', 'backstage', 'identity', 'nats', 'tak', 'monitoring', 'resilience', 'ai', 'cloud', 'workflows', 'inventory-refresh', 'assistant', 'offline'];
function makePack(directory, documents, generatedAt) {
  const cards = directory.filter(x => /^https?:\/\//.test(x.url)).map(x =>
    `<li><a href="${escape(x.url)}" rel="noreferrer">${escape(x.title)}</a> <span>${escape(x.groupTitle)}</span><p>${escape(x.description)}</p></li>`).join('');
  const index = documents.map(x => `<li><a href="#${escape(x.name)}">${escape(x.name)}</a></li>`).join('');
  const docs = documents.map(x => `<article id="${escape(x.name)}"><h2>${escape(x.name)}</h2><pre>${escape(x.content)}</pre><a href="#top">Back to index</a></article>`).join('');
  return `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'none'"><title>OneMind private offline pack</title><style>body{font:16px/1.6 system-ui,sans-serif;max-width:1000px;margin:auto;padding:24px;color:#24262a;background:#f5f5f3}h1,h2{line-height:1.2}a{color:#275b74}span{color:#555;font-size:small}li,article{margin-bottom:20px}article{padding:24px;background:white;border:1px solid #ddd;border-radius:12px}pre{white-space:pre-wrap;overflow-wrap:anywhere;font:14px/1.6 ui-monospace,monospace}.notice{padding:16px;background:#e9eef0;border-radius:8px}@media(prefers-color-scheme:dark){body{background:#181b20;color:#eceef1}article{background:#22262c;border-color:#414650}a{color:#a5cedb}span{color:#bbc}.notice{background:#29323b}}@media print{body{background:white;color:black}article{break-before:page}a{color:black}}</style></head><body id="top"><h1>OneMind offline knowledge</h1><div class="notice"><strong>Private exported copy — generated ${escape(generatedAt)}</strong><p>App directory and documentation from the deployed image or local repository. No live health, task, login, secret or conversation data is included. These documents may be outdated. External app links require connectivity and their own login. Logging out of OneMind does not delete this downloaded file. Store it securely.</p></div><p>Use your browser's Find command to search this file. No network requests or JavaScript run while reading it.</p><h2>Runbooks</h2><ul>${index}</ul><h2>Application directory</h2><ul>${cards}</ul>${docs}</body></html>`;
}
module.exports = { makePack, names };
if (require.main === module) {
  const outputIndex = process.argv.indexOf('--output');
  if (outputIndex < 0 || !process.argv[outputIndex + 1]) { console.error('Usage: node scripts/offline-pack.cjs --output /private/path/onemind-offline.html'); process.exit(1); }
  const root = path.resolve(__dirname, '..');
  const directory = JSON.parse(fs.readFileSync(path.join(root, 'packages/app/src/modules/home/directory.json'), 'utf8'));
  const documents = names.map(name => ({ name, content: fs.readFileSync(path.join(root, 'runbooks/docs', `${name}.md`), 'utf8') }));
  // Explicit output only; never overwrite another pack or export into a public app asset folder.
  const output = path.resolve(process.argv[outputIndex + 1]);
  const realOutput = path.join(fs.realpathSync(path.dirname(output)), path.basename(output));
  const publicRoot = fs.realpathSync(path.join(root, 'packages/app/public'));
  if (realOutput === publicRoot || realOutput.startsWith(publicRoot + path.sep)) throw new Error('Choose a private output folder, not public app assets');
  fs.writeFileSync(output, makePack(directory, documents, new Date().toISOString()), { flag: 'wx', mode: 0o600 });
  console.log(`Private offline pack saved: ${output}`);
}

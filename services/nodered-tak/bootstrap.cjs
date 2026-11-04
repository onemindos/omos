'use strict';
const fs = require('node:fs');
const path = require('node:path');

// Install only the immutable image's reviewed module. Never replace user content.
module.exports = function bootstrap(dataDir, nodesDir) {
  fs.mkdirSync(path.join(dataDir, 'node_modules'), { recursive: true });
  const modulePath = path.join(dataDir, 'node_modules/node-red-contrib-tak');
  const target = path.join(nodesDir, 'node-red-contrib-tak');
  try {
    if (!fs.lstatSync(modulePath).isSymbolicLink() || fs.readlinkSync(modulePath) !== target) {
      throw new Error('Unreviewed TAK module exists in /data; preserve it and reconcile offline before starting');
    }
  } catch (error) {
    if (error.code !== 'ENOENT') throw error;
    fs.symlinkSync(target, modulePath, 'dir');
  }
  try {
    fs.writeFileSync(path.join(dataDir, 'flows.json'), '[]\n', { flag: 'wx', mode: 0o600 });
  } catch (error) {
    if (error.code !== 'EEXIST') throw error;
  }
};
if (require.main === module) module.exports('/data', '/opt/onemind-nodes/node_modules');

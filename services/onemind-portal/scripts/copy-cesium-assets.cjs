#!/usr/bin/env node
// Copies the CesiumJS static assets the Ops Center needs into packages/app/public/cesium, so
// `backstage-cli package build` ships them in the app bundle (and with it, the production image).
// The folder is generated and gitignored. Run it before every app build or start; it is idempotent.
//
//   node scripts/copy-cesium-assets.cjs            copy (skips when already current)
//   node scripts/copy-cesium-assets.cjs --check    exit 1 unless the copy is present and current
const fs = require('node:fs');
const path = require('node:path');

const ASSET_DIRS = ['Workers', 'Assets', 'ThirdParty', 'Widgets'];
const root = path.join(__dirname, '..');
const target = path.join(root, 'packages/app/public/cesium');

function locateCesium() {
  // The package is a dependency of plugins/plugin-onemind-ops; yarn hoists it to the root node_modules.
  const pkg = require.resolve('cesium/package.json', {
    paths: [path.join(root, 'plugins/plugin-onemind-ops'), root],
  });
  return path.dirname(pkg);
}

function main(argv) {
  const check = argv.includes('--check');
  const cesiumDir = locateCesium();
  const version = JSON.parse(
    fs.readFileSync(path.join(cesiumDir, 'package.json'), 'utf8'),
  ).version;
  const build = path.join(cesiumDir, 'Build/Cesium');
  const stamp = path.join(target, '.cesium-version');
  const current =
    fs.existsSync(stamp) &&
    fs.readFileSync(stamp, 'utf8').trim() === version &&
    ASSET_DIRS.every(d => fs.existsSync(path.join(target, d)));
  if (check) {
    if (!current) {
      console.error(
        `Cesium ${version} assets are missing from ${target}. Run: node scripts/copy-cesium-assets.cjs`,
      );
      return 1;
    }
    return 0;
  }
  if (current) {
    console.log(
      `cesium ${version} assets already in packages/app/public/cesium`,
    );
    return 0;
  }
  for (const d of ASSET_DIRS) {
    if (!fs.existsSync(path.join(build, d))) {
      console.error(
        `Missing ${path.join(
          build,
          d,
        )}. Is the cesium package installed (yarn install)?`,
      );
      return 1;
    }
  }
  fs.rmSync(target, { recursive: true, force: true });
  fs.mkdirSync(target, { recursive: true });
  for (const d of ASSET_DIRS)
    fs.cpSync(path.join(build, d), path.join(target, d), { recursive: true });
  fs.writeFileSync(stamp, `${version}\n`);
  console.log(
    `copied cesium ${version} ${ASSET_DIRS.join(
      ', ',
    )} to packages/app/public/cesium`,
  );
  return 0;
}

if (require.main === module) process.exit(main(process.argv.slice(2)));
module.exports = { ASSET_DIRS, main };

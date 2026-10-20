#!/usr/bin/env bash
set -euo pipefail
node .yarn/releases/yarn-4.13.0.cjs exec techdocs-cli generate --no-docker --source-dir runbooks --output-dir techdocs/default/component/onemind-runbooks
test -s techdocs/default/component/onemind-runbooks/techdocs_metadata.json
test -s techdocs/default/component/onemind-runbooks/search/search_index.json

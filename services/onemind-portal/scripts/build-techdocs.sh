#!/usr/bin/env bash
set -euo pipefail
npx --yes @techdocs/cli@1.12.0 generate --no-docker --source-dir runbooks --output-dir techdocs/default/component/onemind-runbooks
 test -s techdocs/default/component/onemind-runbooks/techdocs_metadata.json
 test -s techdocs/default/component/onemind-runbooks/search/search_index.json

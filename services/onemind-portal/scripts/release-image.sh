#!/usr/bin/env bash
# Build a clean committed portal with the same Linux gates as GitHub Actions.
# Default is local --load. Publishing is explicit and never changes Kubernetes/Git.
set -euo pipefail
publish=false
output=''
usage() { echo 'Usage: release-image.sh [--push] [--output-dir DIRECTORY]'; }
while (($#)); do
  case "$1" in
    --push) publish=true; shift ;;
    --output-dir) [[ $# -ge 2 && -n $2 ]] || { usage >&2; exit 2; }; output=$2; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) usage >&2; exit 2 ;;
  esac
done
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
root=$(git -C "$script_dir" rev-parse --show-toplevel)
cd "$root"
project=services/onemind-portal
if [[ -n $(git status --porcelain --untracked-files=normal -- "$project" .github/workflows/build-onemind-portal.yml) ]]; then
  echo 'Commit portal/build-workflow changes before producing a commit-tagged release.' >&2
  exit 1
fi
sha=$(git rev-parse HEAD)
image="${IMAGE_REPOSITORY:-ghcr.io/onemindos/onemind-portal}:sha-$sha"
docker info >/dev/null
context=$(mktemp -d "${TMPDIR:-/tmp}/onemind-portal-context.XXXXXX")
trap 'rm -rf "$context"' EXIT
# Archive only tracked files from this exact commit: no local secrets or host binaries.
git archive "$sha" "$project" | tar -x -C "$context" --strip-components=2
if [[ -z $output ]]; then
  output=$(mktemp -d "${TMPDIR:-/tmp}/onemind-portal-release.XXXXXX")
else
  mkdir -p "$output"
  output=$(cd "$output" && pwd)
fi
# Never let an old digest masquerade as the result of a failed new build.
if [[ -e $output/image-digest.txt || -e $output/release.json ]]; then
  echo 'Output directory already contains release artifacts; choose a fresh directory.' >&2
  exit 1
fi
args=(buildx build --platform linux/amd64 --file "$context/packages/backend/Dockerfile"
  --tag "$image" --metadata-file "$output/build-metadata.json"
  --label "org.opencontainers.image.revision=$sha"
  --label 'org.opencontainers.image.source=https://github.com/onemindos/omos')
if $publish; then args+=(--push); else args+=(--load); fi
docker "${args[@]}" "$context"
python3 - "$output" "$image" "$sha" "$publish" <<'PY'
import json, pathlib, re, subprocess, sys
folder, image, commit, publish = sys.argv[1:]
out = pathlib.Path(folder)
metadata = json.loads((out / 'build-metadata.json').read_text())
digest = metadata.get('containerimage.digest')
published = publish == 'true'
if published:
    if not isinstance(digest, str) or not re.fullmatch(r'sha256:[a-f0-9]{64}', digest):
        raise SystemExit('Build did not report a valid published digest')
    reference = image.rsplit(':', 1)[0] + '@' + digest
    subprocess.run(['docker', 'buildx', 'imagetools', 'inspect', reference],
                   check=True, stdout=subprocess.DEVNULL)
    (out / 'image-digest.txt').write_text(digest + '\n')
record = {'commit': commit, 'image': image, 'published': published,
          'platform': 'linux/amd64', 'digest': digest,
          'gates': ['immutable-yarn', 'typescript', 'techdocs', 'backend-build']}
(out / 'release.json').write_text(json.dumps(record, indent=2) + '\n')
print('Published and verified' if published else 'Built locally only; no deployment digest emitted')
print('Release artifacts: ' + str(out))
PY

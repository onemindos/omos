"""Fail CI when an Argo Application lacks an explicit MCP coverage disposition."""
from pathlib import Path
import sys
import yaml

ROOT = Path(__file__).resolve().parents[1]


def applications():
    found = {}
    for directory in ('platform/argocd', 'platform/ecosystem/applications', 'platform/ecosystem/foundation-applications'):
        for path in (ROOT / directory).glob('*.yaml'):
            for doc in yaml.safe_load_all(path.read_text()):
                if isinstance(doc, dict) and doc.get('kind') == 'Application':
                    found[doc['metadata']['name']] = str(path.relative_to(ROOT))
    return found


def main():
    actual = applications()
    coverage = yaml.safe_load((ROOT / 'platform/hermes/mcp/coverage.yaml').read_text())['applications']
    apps = {app['id'] for app in yaml.safe_load((ROOT / 'platform/hermes/mcp/apps.yaml').read_text())['apps']}
    errors = []
    for name in sorted(actual.keys() - coverage.keys()):
        errors.append(f'{name}: add an explicit app mapping or supporting-infrastructure disposition')
    for name in sorted(coverage.keys() - actual.keys()):
        errors.append(f'{name}: remove stale Application coverage')
    for name, item in coverage.items():
        if not item.get('reason') or not item.get('apps') or not set(item['apps']) <= apps:
            errors.append(f'{name}: reason and valid app IDs are required')
    if errors:
        print('\n'.join(errors), file=sys.stderr)
        return 1
    print(f'All {len(actual)} Git-declared Argo Applications mapped to MCP access dispositions')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())

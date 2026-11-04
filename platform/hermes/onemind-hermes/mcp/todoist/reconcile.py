#!/usr/bin/env python3
"""Set only the Git-owned Todoist server using supported Hermes configuration CLI."""
import argparse
import json
import os
from pathlib import Path
import subprocess


def reconcile(manifest, home, cli, runner=subprocess.run):
    settings = json.loads(Path(manifest).read_text())
    # Match the working public OAuth definition, plus disable remote sampling.
    # Do not accept private URLs, headers, OAuth state or new tool restrictions.
    expected = {
        'url': 'https://ai.todoist.net/mcp', 'auth': 'oauth', 'enabled': True,
        'tools': {'exclude': ['search', 'fetch', 'export-project-template', 'import-project-template']},
        'sampling': {'enabled': False},
    }
    if settings != expected:
        raise ValueError('Expected the public Todoist OAuth configuration only')
    result = runner(
        [str(cli), 'config', 'set', 'mcp_servers.todoist',
         json.dumps(settings, separators=(',', ':'))],
        env=dict(os.environ, HERMES_HOME=str(home)),
        capture_output=True, text=True, check=False,
    )
    if result.returncode:
        # Never echo output: the CLI may include resolved credentials.
        raise RuntimeError('Hermes config update failed for todoist')
    print('todoist: configuration reconciled')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--manifest', type=Path, default=Path(__file__).with_name('server.json'))
    parser.add_argument('--home', type=Path, default=Path(os.environ.get('HERMES_HOME', '/opt/data')))
    parser.add_argument('--cli', default='/opt/hermes/bin/hermes')
    args = parser.parse_args()
    reconcile(args.manifest, args.home, args.cli)


if __name__ == '__main__':
    main()

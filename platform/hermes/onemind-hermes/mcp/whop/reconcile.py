#!/usr/bin/env python3
"""Set only the Git-owned Whop server using supported Hermes configuration CLI."""
import argparse
import json
import os
from pathlib import Path
import subprocess


def reconcile(manifest, home, cli, runner=subprocess.run):
    settings = json.loads(Path(manifest).read_text())
    expected = {
        'url': 'https://mcp.whop.com/mcp', 'auth': 'oauth',
        'connect_timeout': 300, 'oauth': {'timeout': 300},
        'sampling': {'enabled': False},
    }
    if settings != expected:
        raise ValueError('Expected the public Whop OAuth configuration only')
    result = runner(
        [str(cli), 'config', 'set', 'mcp_servers.whop',
         json.dumps(settings, separators=(',', ':'))],
        env=dict(os.environ, HERMES_HOME=str(home)),
        capture_output=True, text=True, check=False,
    )
    if result.returncode:
        # Never echo output: the CLI may include resolved credentials.
        raise RuntimeError('Hermes config update failed for whop')
    print('whop: configuration reconciled')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--manifest', type=Path, default=Path(__file__).with_name('server.json'))
    parser.add_argument('--home', type=Path, default=Path(os.environ.get('HERMES_HOME', '/opt/data')))
    parser.add_argument('--cli', default='/opt/hermes/bin/hermes')
    args = parser.parse_args()
    reconcile(args.manifest, args.home, args.cli)


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
"""Set only the Git-owned Whop server using supported Hermes configuration CLI."""
import argparse
import json
import os
import re
from pathlib import Path
import subprocess


# The Whop MCP exposes 464 tools. Git owns an allowlist (tools.include) of reads and low-risk writes. A name
# that would delete, move money or change the account must never be listed, so the reconciler refuses it.
FORBIDDEN = re.compile(r'delete|remove|payout|transfer|fees?(_|$)|refund|void|capture|charge|swap|trade|deposit|'
                       r'token|api-keys|external-accounts|cards_|payment-methods|accounts_(create|update|transfer)|'
                       r'users_update|users_preferences_update|financing|identity|verifications|cancel|reverse')


def valid_allowlist(tools):
    if not isinstance(tools, dict) or set(tools) != {'include'}:
        return False
    names = tools['include']
    return (isinstance(names, list) and 0 < len(names) <= 300 and len(set(names)) == len(names)
            and all(isinstance(n, str) and re.fullmatch(r'[a-z0-9_-]+', n) and not FORBIDDEN.search(n) for n in names))


def reconcile(manifest, home, cli, runner=subprocess.run):
    settings = json.loads(Path(manifest).read_text())
    expected = {
        'url': 'https://mcp.whop.com/mcp', 'auth': 'oauth',
        'connect_timeout': 300, 'oauth': {'timeout': 300},
        'sampling': {'enabled': False},
    }
    tools = settings.pop('tools', None) if isinstance(settings, dict) else None
    if settings != expected or not valid_allowlist(tools):
        raise ValueError('Expected the public Whop OAuth configuration only')
    settings['tools'] = tools
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

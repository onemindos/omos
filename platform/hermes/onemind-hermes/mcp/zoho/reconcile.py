#!/usr/bin/env python3
"""Reconcile only Git-owned Zoho MCP entries through the supported Hermes CLI.

The manifest contains environment references, never endpoint credentials. This
runs before the gateway on every pod start, for both existing and new PVCs.
"""
import argparse
import json
import os
from pathlib import Path
import subprocess


def reconcile(manifest, home, cli, runner=subprocess.run):
    entries = json.loads(Path(manifest).read_text())
    if not isinstance(entries, dict) or len(entries) != 9:
        raise ValueError('Expected exactly nine Zoho server definitions')
    for name, settings in entries.items():
        variable = name.upper().replace('-', '_') + '_MCP_URL'
        if not name.startswith('zoho-') or settings.get('url') != '${' + variable + '}':
            raise ValueError('Expected a named Zoho environment reference')
        if settings.get('sampling', {}).get('enabled') is not False:
            raise ValueError('Remote sampling must be disabled')
        included = settings.get('tools', {}).get('include')
        if not isinstance(included, list) or not included or any('*' in t for t in included):
            raise ValueError('Expected explicit read tool names')
    env = dict(os.environ, HERMES_HOME=str(home))
    for name, settings in entries.items():
        result = runner(
            [str(cli), 'config', 'set', 'mcp_servers.' + name,
             json.dumps(settings, separators=(',', ':'))],
            env=env, capture_output=True, text=True, check=False,
        )
        if result.returncode:
            # Never echo stderr: a provider or CLI can echo resolved credentials.
            raise RuntimeError('Hermes config update failed for ' + name)
        print(name + ': configuration reconciled')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--manifest', type=Path, default=Path(__file__).with_name('servers.json'))
    parser.add_argument('--home', type=Path, default=Path(os.environ.get('HERMES_HOME', '/opt/data')))
    parser.add_argument('--cli', default='/opt/hermes/bin/hermes')
    args = parser.parse_args()
    reconcile(args.manifest, args.home, args.cli)


if __name__ == '__main__':
    main()

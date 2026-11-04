#!/usr/bin/env python3
"""Reconcile public native-source settings only; never resolve credentials."""
import argparse
import json
import os
from pathlib import Path
import subprocess



def reconcile(manifest, home, cli, runner=subprocess.run):
    settings = json.loads(Path(manifest).read_text())
    expected = {
        'enabled': True,
        'env': {'CLICKUP_API_TOKEN':
                'op://aaeyris5z44pdqspk277u2oppi/q3jm77jbor2ymoc5xlyltrv7ue/credential'},
        'account': '', 'service_account_token_env': 'OP_SERVICE_ACCOUNT_TOKEN',
        'binary_path': '/opt/onepassword/op', 'cache_ttl_seconds': 0,
        'override_existing': True,
    }
    if settings != expected:
        raise ValueError('Expected approved native 1Password configuration only')
    result = runner(
        [str(cli), 'config', 'get', 'terminal.env_passthrough', '--json'],
        env=dict(os.environ, HERMES_HOME=str(home)),
        capture_output=True, text=True, check=False,
    )
    if result.returncode:
        raise RuntimeError('Hermes terminal configuration read failed')
    passthrough = json.loads(result.stdout)
    if not isinstance(passthrough, list) or not all(isinstance(v, str) for v in passthrough):
        raise ValueError('Expected terminal.env_passthrough to be a list of names')
    passthrough = [v for v in passthrough if v != 'OP_SERVICE_ACCOUNT_TOKEN']
    if 'CLICKUP_API_TOKEN' not in passthrough:
        passthrough.append('CLICKUP_API_TOKEN')
    # Configure without any bootstrap credential in this init container. Preserve
    # other sources, MCP, model, terminal and profile settings using the CLI.
    for key, value in [('terminal.env_passthrough', passthrough),
                       ('secrets.onepassword', settings)]:
        result = runner(
            [str(cli), 'config', 'set', key, json.dumps(value, separators=(',', ':'))],
            env=dict(os.environ, HERMES_HOME=str(home)),
            capture_output=True, text=True, check=False,
        )
        if result.returncode:
            raise RuntimeError('Hermes native 1Password configuration failed')
    print('native 1Password: public configuration reconciled')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--manifest', type=Path, default=Path(__file__).with_name('settings.json'))
    parser.add_argument('--home', type=Path, default=Path('/opt/data'))
    parser.add_argument('--cli', default='/opt/hermes/bin/hermes')
    args = parser.parse_args()
    reconcile(args.manifest, args.home, args.cli)


if __name__ == '__main__':
    main()

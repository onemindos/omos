"""Build-time, hash-locked correction for the pinned Hermes v0.22.5 source.

Never run against a live installation. Apply to a scratch source copy or during
an image build; unknown source fails closed before any write.
"""
import argparse
import hashlib
from pathlib import Path

RELATIVE_PATH = Path('tools/environments/local_env_policy.py')
SOURCE_SHA256 = '4429af03406ead5e79e5cfe711787030f40d42bf50c6253b4abc6ee0154a77d4'
OLD = '''    upper = key.upper()
    if upper.startswith("AUXILIARY_") and upper.endswith(("_API_KEY", "_BASE_URL")):
'''
NEW = '''    upper = key.upper()
    # Native 1Password resolves in-process and constructs its own minimal op env.
    # Generic children must never receive the vault bootstrap, even via a skill,
    # operator passthrough, _HERMES_FORCE_ override or inherit_credentials=True.
    if upper == "OP_SERVICE_ACCOUNT_TOKEN":
        return True
    if upper.startswith("AUXILIARY_") and upper.endswith(("_API_KEY", "_BASE_URL")):
'''


PROBE_OLD = '''        from agent.delegation_context import delegated_child_subprocess_env
        return subprocess.run(
            [python_path, "-c", code], timeout=5, capture_output=True, text=text,
            creationflags=subprocess.CREATE_NO_WINDOW if _IS_WINDOWS else 0,
            stdin=subprocess.DEVNULL, env=delegated_child_subprocess_env(),
'''
PROBE_NEW = '''        from tools.environments.local import hermes_subprocess_env
        return subprocess.run(
            [python_path, "-c", code], timeout=5, capture_output=True, text=text,
            creationflags=subprocess.CREATE_NO_WINDOW if _IS_WINDOWS else 0,
            stdin=subprocess.DEVNULL, env=hermes_subprocess_env(),
'''
PATCHES = (
    (RELATIVE_PATH, SOURCE_SHA256, OLD, NEW),
    (Path('tools/code_execution_env.py'),
     '1c6e0c42a7f4c59f0c65a02aec397c92c87860226ab181de64751e00fecf2785',
     PROBE_OLD, PROBE_NEW),
)


def apply(root: Path) -> None:
    pending = []
    # Preflight every source before writing any of them.
    for relative, digest, old, new in PATCHES:
        target = root / relative
        source = target.read_bytes()
        if hashlib.sha256(source).hexdigest() != digest:
            raise ValueError('Unrecognized Hermes source; refusing to patch')
        text = source.decode('utf-8')
        if text.count(old) != 1:
            raise ValueError('Expected exactly one source patch location')
        patched = text.replace(old, new)
        compile(patched, str(target), 'exec')
        pending.append((target, patched))
    for target, patched in pending:
        target.write_text(patched, encoding='utf-8')
    print('Applied hash-locked 1Password bootstrap child-environment denial')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--source-root', required=True, type=Path)
    apply(parser.parse_args().source_root)

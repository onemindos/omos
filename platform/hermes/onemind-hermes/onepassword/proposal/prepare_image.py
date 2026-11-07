"""Render ONLY the inactive Deployment snapshot with a real build digest.

No network access, activation, Secret generation, or writes to the repository.
Digest provenance must be verified against the successful build by the reviewer.
"""
import argparse
from pathlib import Path
import re
import sys

HERE = Path(__file__).parent
UPSTREAM = ('nousresearch/hermes-agent:v0.23.0@sha256:'
            'adb7890916de63245523492be8f68f9d2095fecbec4e2a6dab7960381f0fc002')


def render(reference):
    if not re.fullmatch(r'ghcr\.io/onemindos/hermes-agent@sha256:[0-9a-f]{64}', reference):
        raise ValueError('Expected corrected GHCR repository with a lowercase SHA256 digest')
    if reference.endswith(UPSTREAM.split('@', 1)[1]):
        raise ValueError('Uncorrected upstream digest is forbidden')
    content = (HERE / 'deployment.yaml').read_text()
    old = '        - name: hermes\n          image: ' + UPSTREAM + '\n'
    if content.count(old) != 1:
        raise ValueError('Proposal main-container image drift; review snapshot before substitution')
    return content.replace(old, '        - name: hermes\n          image: ' + reference + '\n', 1)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--image', required=True, help='Verified image reference from build artifact')
    args = parser.parse_args()
    try:
        content = render(args.image)
    except ValueError as error:
        parser.error(str(error))
    sys.stdout.write(content)

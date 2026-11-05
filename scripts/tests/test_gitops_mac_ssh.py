"""Run Mac SSH bootstrap contracts in the portable GitOps quality gate."""
import importlib.util
from pathlib import Path
import sys
import unittest

ROOT = Path(__file__).resolve().parents[2]
DIRECTORY = ROOT / 'platform/hermes/onemind-hermes/mac-access'


def load_tests(loader, tests, pattern):
    # Contract modules also import the actual helper modules during test methods.
    location = str(DIRECTORY)
    if location not in sys.path:
        sys.path.insert(0, location)
    suite = unittest.TestSuite()
    for name in ('test_mac_ssh', 'test_prepare_mac_ssh'):
        spec = importlib.util.spec_from_file_location('mac_ssh_contracts_' + name, DIRECTORY / (name + '.py'))
        assert spec is not None and spec.loader is not None
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        suite.addTests(loader.loadTestsFromModule(module))
    return suite

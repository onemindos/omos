"""Run Zoho startup/credential contracts in the existing portable GitOps CI gate."""
import importlib.util
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[2]
TARGET = ROOT / 'platform/hermes/onemind-hermes/mcp/zoho/test_reconcile.py'


def load_tests(loader, tests, pattern):
    spec = importlib.util.spec_from_file_location('zoho_gitops_contracts', TARGET)
    assert spec is not None and spec.loader is not None
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return loader.loadTestsFromModule(module)

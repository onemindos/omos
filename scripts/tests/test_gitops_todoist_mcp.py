"""Run Todoist recovery contracts in the existing portable GitOps quality gate."""
import importlib.util
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
TARGET = ROOT / 'platform/hermes/onemind-hermes/mcp/todoist/test_todoist.py'


def load_tests(loader, tests, pattern):
    spec = importlib.util.spec_from_file_location('todoist_gitops_contracts', TARGET)
    assert spec is not None and spec.loader is not None
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return loader.loadTestsFromModule(module)

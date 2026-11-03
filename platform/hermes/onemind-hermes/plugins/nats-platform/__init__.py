# OneMind change 3 (ONEMIND-CHANGES.md): put the init container's package directory on sys.path.
# The `hermes` launcher unsets PYTHONPATH and the venv has no user site, so a variable or a .pth
# cannot do it. Appended, not inserted, so the venv's own pydantic and pynacl win.
import os as _os
import sys as _sys

_deps = _os.environ.get("ONEMIND_GATEWAY_DEPS", "/opt/data/.local/gateway-deps")
if _os.path.isdir(_deps) and _deps not in _sys.path:
    _sys.path.append(_deps)

# When force-included into the wheel as ``hermes_nats_gateway/__init__.py``,
# this file is loaded as part of a proper package and ``from .adapter import
# register`` resolves. When pytest discovers this file at the rootdir of the
# source tree (no parent package), the relative import would fail — skip it
# so collection succeeds.
if __package__:
    from .adapter import register  # noqa: F401

__all__ = ["register"]

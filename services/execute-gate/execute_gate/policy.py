"""Tool classification. Default deny: anything not explicitly physical is destructive."""

from __future__ import annotations

import fnmatch
from dataclasses import dataclass, field
from pathlib import Path

DESTRUCTIVE = "destructive"
PHYSICAL = "physical"

# Word the owner must type for each class.
WORD_FOR_CLASS = {DESTRUCTIVE: "execute", PHYSICAL: "override"}


@dataclass(frozen=True)
class Policy:
    physical: tuple = ()
    # Informational: destructive is the default class; listed globs only document intent.
    destructive: tuple = ()
    allowed_subjects: frozenset = field(default_factory=lambda: frozenset({"hermes"}))

    def classify(self, resource: str) -> str:
        """resource is "target/tool". Physical wins; everything else is destructive."""
        for pat in self.physical:
            if fnmatch.fnmatchcase(resource, pat):
                return PHYSICAL
        return DESTRUCTIVE

    def subject_allowed(self, subject: str) -> bool:
        return subject in self.allowed_subjects


def load_policy(path, allowed_subjects=None) -> Policy:
    subjects = frozenset(allowed_subjects) if allowed_subjects else frozenset({"hermes"})
    if not path:
        return Policy(allowed_subjects=subjects)
    import yaml  # PyYAML, pinned in requirements.lock

    data = yaml.safe_load(Path(path).read_text(encoding="utf-8")) or {}
    if not isinstance(data, dict):
        raise ValueError("policy file must be a mapping")
    phys = tuple(str(x) for x in data.get("physical") or [])
    dest = tuple(str(x) for x in data.get("destructive") or [])
    return Policy(physical=phys, destructive=dest, allowed_subjects=subjects)

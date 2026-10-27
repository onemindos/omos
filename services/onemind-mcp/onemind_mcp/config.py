"""Explicit administrator-owned app catalog; credentials never enter the catalog."""
import os
import re
from pathlib import Path
from typing import Literal
from urllib.parse import urlsplit

import yaml
from pydantic import BaseModel, ConfigDict, Field, field_validator


class App(BaseModel):
    model_config = ConfigDict(extra="forbid")
    id: str = Field(pattern=r"^[a-z][a-z0-9_]{0,31}$")
    adapter: str | None = None
    enabled: bool = True
    blocked_reason: str | None = Field(default=None, max_length=500)
    allow_writes: bool = False
    description: str = Field(min_length=1, max_length=300)
    url_env: str | None = None
    auth_header: Literal["Authorization", "X-N8N-API-KEY"] = "Authorization"
    token_env: str | None = None

    @field_validator("url_env", "token_env")
    @classmethod
    def environment_name(cls, value: str | None) -> str | None:
        if value is not None and not re.fullmatch(r"[A-Z][A-Z0-9_]*", value):
            raise ValueError("Expected an environment variable name")
        return value

    def connection(self) -> tuple[str, str | None]:
        if not self.url_env:
            raise ValueError(f"{self.id}: enabled app requires url_env")
        url = os.environ.get(self.url_env, "")
        parsed = urlsplit(url)
        if (parsed.scheme not in {"http", "https"} or not parsed.hostname
                or parsed.username or parsed.password or parsed.query or parsed.fragment):
            raise ValueError(f"{self.id}: {self.url_env} must contain an HTTP(S) base URL without credentials, query or fragment")
        token = os.environ.get(self.token_env, "") if self.token_env else None
        if self.token_env and not token:
            raise ValueError(f"{self.id}: missing {self.token_env}")
        return url.rstrip("/"), token


class Catalog(BaseModel):
    model_config = ConfigDict(extra="forbid")
    apps: list[App] = Field(min_length=1, max_length=100)

    @field_validator("apps")
    @classmethod
    def unique_ids(cls, apps: list[App]) -> list[App]:
        if len({app.id for app in apps}) != len(apps):
            raise ValueError("Duplicate app IDs")
        return apps


def load_catalog(path: str | Path) -> Catalog:
    return Catalog.model_validate(yaml.safe_load(Path(path).read_text()))

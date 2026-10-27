"""Explicit administrator-owned app catalog; credentials never enter the catalog."""
import os
import re
from pathlib import Path
from urllib.parse import urlsplit

import yaml
from pydantic import BaseModel, ConfigDict, Field, field_validator

MAX_NODES = 8


class App(BaseModel):
    model_config = ConfigDict(extra="forbid")
    id: str = Field(pattern=r"^[a-z][a-z0-9_]{0,31}$")
    adapter: str
    description: str = Field(min_length=1, max_length=300)
    url_env: str
    token_env: str | None = None

    @field_validator("url_env", "token_env")
    @classmethod
    def environment_name(cls, value: str | None) -> str | None:
        if value is not None and not re.fullmatch(r"[A-Z][A-Z0-9_]*", value):
            raise ValueError("Expected an environment variable name")
        return value

    def _base_url(self, url: str, env: str) -> str:
        parsed = urlsplit(url)
        if (parsed.scheme not in {"http", "https"} or not parsed.hostname
                or parsed.username or parsed.password or parsed.query or parsed.fragment):
            raise ValueError(f"{self.id}: {env} must contain an HTTP(S) base URL without credentials, query or fragment")
        return url.rstrip("/")

    def _token(self) -> str | None:
        token = os.environ.get(self.token_env, "") if self.token_env else None
        if self.token_env and not token:
            raise ValueError(f"{self.id}: missing {self.token_env}")
        return token

    def connection(self) -> tuple[str, str | None]:
        return self._base_url(os.environ.get(self.url_env, ""), self.url_env), self._token()

    def connections(self) -> list[tuple[str, str | None]]:
        """Every node of a multi-node app: comma-separated `<url_env>S` when set, else `url_env`."""
        env = f"{self.url_env}S"
        urls = [url.strip() for url in os.environ.get(env, "").split(",") if url.strip()]
        if not urls:
            return [self.connection()]
        if len(urls) > MAX_NODES:
            raise ValueError(f"{self.id}: {env} lists more than {MAX_NODES} URLs")
        token = self._token()
        return [(url, token) for url in dict.fromkeys(self._base_url(url, env) for url in urls)]


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

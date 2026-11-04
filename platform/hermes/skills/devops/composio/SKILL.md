---
name: composio
description: "Use when working with Gmail, Google Calendar/Drive/Docs, GitHub, ClickUp and other apps through the Composio MCP server."
version: 4.0.0
author: Legacy
license: MIT
platforms: [linux]
metadata:
  hermes:
    tags: [composio, integrations, mcp, gmail, google, github, clickup]
    related_skills: [hermes-mcp-config]
---

# Composio MCP

Composio exposes many third-party apps through one MCP server with a handful of
meta-tools. The agent calls the meta-tools; Composio holds the OAuth connections.

## Setup facts

- The server is configured in Hermes as an HTTP MCP server (`hermes mcp list`,
  see `hermes-mcp-config`). Its consumer key is supplied by a SealedSecret and
  injected as an environment variable. Never print, log or commit it.
- Connected apps and their account IDs differ between instances and change when
  connections are re-authorized. Discover them with `COMPOSIO_SEARCH_TOOLS`
  (`toolkit_connection_statuses`); do not rely on remembered account IDs.
- The consumer key works only against the MCP endpoint, not the REST API. Do not
  call `initialize` by hand; the server initializes itself.
- Raw HTTP calls need `Accept: application/json, text/event-stream`, or the
  server answers "Not Acceptable".

## Meta-tool workflow

1. `COMPOSIO_SEARCH_TOOLS` first, always. Pass `queries` (each with a normalized
   English `use_case`, and ids in `known_fields` as a `key: value, key: value`
   string, not an array) and `session: {generate_id: true}`. It returns the exact
   tool slugs, whether a connection is active, pitfalls and the tool schemas.
   Keep the returned `session.id` and pass it on every later call.
2. `COMPOSIO_MULTI_EXECUTE_TOOL` runs up to 50 logically independent tools in
   parallel. Use only slugs from step 1 and only schema fields; never invent
   either. Pass `sync_response_to_workbench` (a boolean) and the `session_id`.
   Its `memory` field takes `{app: [strings]}`, not nested objects.
3. Responses arrive as server-sent events: read the `data:` lines. The inner
   `text` is JSON with `data`, `error` and `successful`.
4. If a connection is missing, the search result says so. Ask the owner to
   authorize the app; do not try to work around it.

## Rules

- Pull data live; never guess at an inbox, calendar or repository.
- Reading is free. Sending mail, changing calendar entries, creating or closing
  issues and PRs, or sharing files on the owner's behalf needs explicit approval
  for that action.
- Keep personal data (addresses, family, account numbers) out of tool memory
  fields and out of skills.
- Do not run `composio login` on a headless pod; it needs a browser click.

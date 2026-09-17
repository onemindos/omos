---
name: skill-forge
description: "Autonomous skill builder — generate complete OpenClaw skills from API specs, deploy them into the workspace, and validate existing skills. Use when Legacy needs to integrate a new API, build a custom tool, or audit installed skills. The MCP Builder pattern for OpenClaw."
metadata: { "forge": true, "clawdbot": { "emoji": "🔨" } }
---

# Skill Forge — Autonomous Skill Builder

Build new OpenClaw skills on demand. Research an API, define the tools you need, and deploy a working skill — no manual coding required.

## When to Use

- User asks to "integrate X API" or "add a tool for X"
- User provides an API URL or OpenAPI spec and wants a skill built
- User wants to audit or validate existing skills
- User says "build me a skill that can..."
- User wants to know what skills are installed and their health

## Workflow: Building a New Skill

### From a Manual Definition (most common)

1. Research the API — use `web_search` or `web_fetch` to find the API docs
2. Identify the base URL, auth method, and key endpoints
3. Call `skill_forge_create` with the API spec:

```json
{
  "name": "linear-issues",
  "description": "Manage Linear issues and projects",
  "base_url": "https://api.linear.app",
  "auth_type": "bearer",
  "auth_env_var": "LINEAR_API_KEY",
  "tools": [
    {
      "name": "list_issues",
      "description": "List issues in a project",
      "method": "GET",
      "path": "/issues",
      "parameters": [
        {
          "name": "project_id",
          "type": "string",
          "description": "Project ID",
          "required": true,
          "in": "query"
        }
      ]
    }
  ]
}
```

4. If npm dependencies are needed, run install in the skill directory
5. Restart gateway: `openclaw reload`

### From an OpenAPI/Swagger Spec (fastest)

If the API has a published OpenAPI spec:

```json
{
  "name": "petstore",
  "spec_url": "https://petstore3.swagger.io/api/v3/openapi.json",
  "auth_env_var": "PETSTORE_API_KEY",
  "include_paths": ["/pet"],
  "exclude_methods": ["DELETE"],
  "max_tools": 10
}
```

This auto-extracts endpoints, parameters, auth, and generates everything.

## Tools

- `skill_forge_create(name, description, base_url, tools, ...)` — Generate and deploy a complete skill
- `skill_forge_from_openapi(name, spec_url|spec_json, ...)` — Generate from OpenAPI/Swagger spec
- `skill_forge_validate(skill_name)` — Check a skill's structure and report issues
- `skill_forge_list(filter?)` — List all installed skills with health status

## Auth Types

| Type             | How it works                           |
| ---------------- | -------------------------------------- |
| `bearer`         | `Authorization: Bearer <token>` header |
| `api-key-header` | Custom header (default `X-API-Key`)    |
| `api-key-query`  | API key as query parameter             |
| `basic`          | `Authorization: Basic <base64>` header |
| `none`           | No authentication                      |

## Parameter Locations

| `in` value | Behavior                                        |
| ---------- | ----------------------------------------------- |
| `path`     | Replaces `{param}` in the URL path              |
| `query`    | Added as URL query parameter                    |
| `body`     | Sent in JSON request body (POST/PUT/PATCH only) |

## Generated Skill Structure

```
skills/<name>-<version>/
├── index.js         ← Generated Node.js tool implementations
├── package.json     ← Tool definitions + metadata
├── SKILL.md         ← Natural language instructions
└── .forge-created   ← Marker file (identifies forge-built skills)
```

## Important Notes

- Generated skills use Node.js built-in `https`/`http` — no external dependencies by default
- All generated code is synchronous-safe and follows OpenClaw skill conventions
- The `.forge-created` marker lets you distinguish auto-generated skills from manual ones
- Always add the API key/token to 1Password and `.env.tpl` before using the skill on VPS
- Skills require a gateway restart to load: `systemctl restart openclaw` on VPS

## Security

- Review generated `index.js` before deploying to production
- Never hardcode API keys — always use environment variables
- Generated skills inherit the gateway's trust boundary
- Use `skill_forge_validate` to check for structural issues before loading

---

_Skill Forge — the MCP Builder pattern, native to OpenClaw._


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

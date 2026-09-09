# NyxID

**Repo:** github.com/ChronoAIProject/NyxID  
**Language:** Rust (backend + CLI) + React 19/TypeScript (frontend) + React Native (mobile)  
**Role in OneMind:** ACTIVE — API credential gateway between agents and external services

Auth/SSO platform with MCP proxy capabilities. Agents request credentials through NyxID rather than holding API keys directly — centralized credential governance.

## Stack

- Backend: Rust, Axum 0.8, MongoDB 8.0
- Frontend: React 19, TypeScript, Vite 7, TanStack Router, Tailwind CSS 4
- Mobile: React Native 0.79 + Expo 53 (iOS + Android approval app)
- SDK: TypeScript OAuth 2.0 client (`@nyxids/oauth-core`, `@nyxids/oauth-react`)
- CLI: `nyxid` — manages services, keys, catalog, nodes, approvals, SSH, MCP

## Key capabilities

- OAuth/OIDC, MFA, credential brokering, admin management
- MCP proxy: sits between AI agents and external MCP servers
- `nyxid node` — on-premise credential nodes (no cloud dependency)
- `unified_key_service` — auto-provisions UserEndpoint + UserApiKey + UserService in one operation

## Why it matters for OneMind

Agents (Hermes, Spartan, Claude, etc.) call external APIs (weather, maps, social, finance) through NyxID instead of embedding raw keys. Approval gates on the mobile app before sensitive operations execute.

## Integration point

Agents publish to NyxID's MCP proxy endpoint. NyxID holds and rotates API keys, logs access, and can require mobile approval before releasing credentials.

# AI operational context

Connect an OAuth-capable MCP client to `https://onemind-preview.tail717ba0.ts.net/api/mcp-actions/v1` through the private mesh. Authenticate with the existing OIDC identity and approve the client request. Backstage's CIMD authentication is enabled; the default supported client metadata patterns apply. No anonymous or organization-wide machine token is installed.

The actions registry is restricted to read-only, non-destructive actions from catalog and OneMind. Scaffolder mutation actions are excluded. Available context includes catalog queries, catalog model descriptions, documented runbooks and configured live health/metrics. Health includes observation timestamps; documentation reports that it is packaged Git content, not a live assertion. Stale or unknown readings must not be presented as current facts.

Tool payload tracing is disabled. The existing read-only cluster identity cannot read Secrets or change workloads. Deployment changes remain reviewed GitOps proposals. Register additional tools deliberately and check the global action filter; a named MCP server alone does not restrict the default endpoint.

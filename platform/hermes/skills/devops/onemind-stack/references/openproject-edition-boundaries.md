# OpenProject self-hosting and edition checks

- Separate hosting from entitlement: Community and paid Enterprise On-Premises can both be self-hosted. Do not call every missing Enterprise add-on a limitation of self-hosting.
- Use the current pricing feature matrix for gates. Its plain-text extraction can omit check/cross icons and falsely imply availability. Inspect the actual table cells and column headers; current HTML uses icon-ok-circled versus icon-ion-close-circled. Accessibility translation strings can be broken, so do not use those strings alone. Cross-check feature documentation and deployed version.
- Recheck minimum paid tiers rather than using remembered plan names. Verify official MCP, SSO, team planning, portfolios/resources, external work-package sharing, form customization, and custom actions individually. Do not imply ordinary project/member permissions need a paid external-sharing add-on.
- Distinguish core REST API access from the official MCP entitlement. A separate permission-preserving REST connector is an integration design, not an automatic unlock for paid modules.
- Distinguish native Nextcloud file integration from Corporate shared-OIDC SSO; neither implies synchronization with Nextcloud Tasks or Deck.
- Do not assume Todoist-equivalent recurrence or mobile completeness. Check the current mobile beta scope, server prerequisites, offline capability and background push status; test recurring household/field workflows.
- Explain self-host operational obligations separately: server/storage cost, database and attachment backups, restores, TLS, mail delivery, upgrades and monitoring. Enterprise On-Premises still requires the operator to run the infrastructure.
- OpenProject states its software, including paid On-Premises Enterprise, is GPLv3. Do not conflate paid add-on access/support with Mattermost's proprietary Enterprise licensing. Avoid making legal conclusions about modifying entitlement code.

Primary references:
- https://www.openproject.org/pricing/
- https://www.openproject.org/community-edition/
- https://www.openproject.org/docs/enterprise-guide/enterprise-on-premises-guide/frequently-asked-questions/
- https://www.openproject.org/docs/system-admin-guide/integrations/mcp-server/
- https://www.openproject.org/docs/system-admin-guide/integrations/nextcloud/
- https://www.openproject.org/docs/mobile-app-guide/
- https://www.openproject.org/docs/installation-and-operations/

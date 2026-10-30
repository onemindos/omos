# Developer portal versus an owned operations application

Use for long-term OneMind frontend/platform decisions and comparisons involving Backstage, OpsLevel, Cycloid, Refine, React-admin, and low-code platforms.

## Evaluation procedure

1. Fix the decision boundary and gates first: unified end-user application, infrastructure operations console, context/relationship service, or commercial hosting arrangement. Keep open-source ownership, phone workflows and maintenance expectations explicit; these are different procurement and engineering decisions.
2. Refresh live inventory and inspect the deployed frontend, pinned source and existing integrations using `existing-platform-discovery.md`. Build a retained/rebuilt/missing capability comparison rather than treating either a fresh framework or an existing deployment as the default winner.
3. Open the exact supplied feature pages and follow their source-repository links. A company may publish both a commercial product and a separately licensed open-source project; its brand or comparison page does not establish the licence of either.
4. Inspect the actual main service, UI, nested licences, default deployment artifacts and permission resolver. Distinguish accessible source, permissive code licensing, supported self-hosting, production licensing and deny-by-default authorization.
5. Compare the candidate against the same full workflow: select project/location, dispatch bounded work, inspect progress, approve or deny, background/disconnect/reopen, receive evidence, and repeat with a restricted identity. Add recurring human responsibilities so an engineering-only demo cannot silently redefine OneMind.
6. Give one recommended architecture and a bounded acceptance gate, with explicit migration cost, authority boundaries, unknowns and fallback. If new evidence changes the recommendation, say exactly why; do not present a documentation review as a completed integration or physical-phone test.

## Product boundaries and decision rules

- Separate the durable domain/API commitment from the UI-shell decision. Reuse the existing operational services without forcing every future interface into the current developer portal.
- Never promise a framework is universally or '1000%' best from documentation. Separate technical extensibility, current integration value, product UX fit, and measured physical-phone acceptance.
- Treat Backstage as an Apache-2.0 framework/application host, not a finished OneMind product. Its replaceable app extensions permit a custom interface, but coordinated package/template maintenance is an ownership cost. Prefer maintaining our application and extension boundaries over forking the framework core.
- Treat Refine Core as an MIT-licensed, headless React meta-framework for CRUD-heavy applications, not a ready portal, a robotics orchestrator, or a proven mobile mission product. Its UI freedom makes it a credible challenger; that is an architectural judgment, not benchmark evidence.
- Distinguish commercial packaged SSO/RBAC adapters from free core hooks. Refine offers authProvider/accessControlProvider/live-provider interfaces; custom identity, backend authorization, and realtime integrations remain engineering work, not proven shipped Authentik/OpenFGA/NATS adapters.
- Check the full platform licence, not an SDK, runner, installer or integration licence. OpsLevel's published standard terms restrict modification/transfer. Cycloid distinguishes its commercial full platform from standalone OSS tools; do not call the complete platform open source or source-available without the actual grant.
- When considering Flanksource Mission Control, inspect the main service and the separate full UI LICENSE files, not an Apache-licensed SDK, chart or companion project. The UI's canary-checker dashboard exception does not exempt the complete Mission Control UI/backend. The examined main-service licence grants non-production use only and requires a commercial licence otherwise; its definition of production includes any benefit. Do not present the full platform as an unrestricted open-source foundation or assume homestead use is exempt. Its health/configuration/GitOps/playbook approvals can make it a close infrastructure-operations fit without proving it supplies a family/project/map product. Recheck the exact pinned version before deployment: https://raw.githubusercontent.com/flanksource/mission-control/master/LICENSE and https://flanksource.com/docs/reference/playbooks/.
- Distinguish operating Kubernetes applications from composing their domain functions into OneMind. A request to compare infrastructure portals does not, by itself, rescind the user's open-source ownership or whole-life mobile requirements. Offer the best functional commercial candidate separately from the best qualifying owned foundation; do not silently change the ranking criterion.
- Do not falsely claim vendor portals cannot self-host: OpsLevel, Cycloid, Cortex and Roadie Local document self-host routes. Self-hosting, white labelling, or a public installer does not establish independent full-product fork/redistribution rights.
- Distinguish Backstage upstream from commercial providers and distributions. RHDH upstream is Apache-2.0 and based on Backstage. Port Ocean being open source does not establish that the Port platform is.
- Inspect mixed-licence package trees and deployment defaults for low-code products. Paid login/branding/group modules and pro-package dependencies can matter more than a permissive root file. Do not assume every custom authentication path is paid.
- Accept the final UI only after a shared physical-phone workflow: recurring human work; location-linked mission; agent progress; attributable approval; background/disconnect/reopen; result delivery; and less-privileged privacy. Responsive pages, widgets, iframes, and PWA installation are not proof of that workflow.
- Keep live app records separate from desired features. Existing core missions and portal routes do not establish finished recurring tasks, robot authorization or offline operations.

## Roadie, Backstage and agent context graphs

- Distinguish Roadie's commercial Backstage-based product, its public plugin directory, the plugins it authors, and OpenRoadie. Count directory entries as ecosystem integrations, not as exclusive products or plugins Roadie necessarily wrote. Check each plugin's own licence, dependencies and self-hosted installation instructions.
- Start the Backstage relationship comparison with upstream catalog relations and `@backstage/plugin-catalog-graph`: entity cards, a graph page and reusable React graph components. Inspect enabled extensions and populated catalog relations; a dependency declaration or bundle marker does not prove the graph is enabled or useful.
- Separate a relationship visualization from a complete context service. Check ingestion and cross-source joins, provenance/freshness, context bundles, caller-filtered queries, MCP exposure and human-reviewed updates independently; a graph renderer does not supply these semantics.
- Follow the Context Graph page's GitHub link to RoadieHQ/openroadie instead of assuming the whole feature is closed-source or a Backstage plugin. The inspected project is separately Apache-2.0 licensed and describes a standalone organizational context/action service; recheck the actual revision, nested licences, artifacts and feature maturity before relying on that boundary.
- Inspect authorization before agent write access. OpenRoadie's security policy documents a pre-1.0 default with no REST/MCP authorization and a known signing key; the standalone example can interpret a bearer value as scope names rather than authenticate it. Replace the identity/scope resolver and signing key, and exercise denied callers before treating scoped API/MCP interfaces as enforced authorization. Check source-level/field-level disclosure separately: commercial datastore query scopes are documented as spanning all Data Sources. Named deterministic actions are not proof of exact-command human approval or local physical interlocks.
- Compare OpenRoadie's context and action authority with existing OneMind objects/links, OpenFGA, MCP and mission services before adding it. Prefer one authoritative record per scope and a deliberate read-model/projection boundary over independently editable competing graphs or approval engines. Start any context-service evaluation read-only with least-privilege source access.
- Verify Roadie Local's current availability, feature entitlements, licence and offline requirements from current material; a historical launch or beta announcement proves a documented route, not today's support or unrestricted redistribution rights.

## Primary references

- https://roadie.io/product/context-graph/
- https://roadie.io/product/context-graph/relationships/
- https://github.com/RoadieHQ/openroadie
- https://raw.githubusercontent.com/RoadieHQ/openroadie/main/LICENSE
- https://raw.githubusercontent.com/RoadieHQ/openroadie/main/README.md
- https://backstage.io/api/stable/modules/_backstage_plugin-catalog-graph.html
- https://backstage.io/docs/features/software-catalog/well-known-relations/
- https://backstage.io/docs/frontend-system/building-apps/built-in-extensions/
- https://backstage.io/docs/getting-started/keeping-backstage-updated/
- https://raw.githubusercontent.com/backstage/backstage/master/LICENSE
- https://www.opslevel.com/legal/t4-terms-of-service
- https://www.cycloid.io/comparison/cycloid-vs-humanitec-alternative/
- https://refine.dev/core/docs/
- https://refine.dev/core/docs/authentication/auth-provider/
- https://refine.dev/core/docs/authorization/access-control-provider/
- https://refine.dev/core/docs/realtime/live-provider/
- https://raw.githubusercontent.com/refinedev/refine/main/LICENSE

Recheck current terms and required modules before a commercial commitment; this is not a substitute for reviewing the actual agreement.

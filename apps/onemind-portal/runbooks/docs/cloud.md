# Cloud summaries

The first adapter reads Zoho Projects V3 tasks using server-side OAuth. It requests at most the first 100 tasks and explicitly labels partial results; sample counts are not full portal totals. Results are cached for five minutes, concurrent requests are shared, failures are rate limited, and previous results remain visibly stale. The original app remains the full workflow.

## Enable Zoho Projects

Create a US-region OAuth client and grant `ZohoProjects.tasks.READ` for the intended portal. Confirm your subscription/account permits API access in the provider console. Store `ZOHO_CLIENT_ID`, `ZOHO_CLIENT_SECRET`, `ZOHO_REFRESH_TOKEN` and the numeric `ZOHO_PROJECTS_PORTAL_ID` in a SealedSecret producing `backstage-zoho` in namespace backstage. Restart through a reviewed GitOps update after secret creation. Never commit plaintext credentials or put them in chat. This adapter is fixed to US endpoints; other regions require a reviewed endpoint change.

The link's portal slug is not the numeric API portal ID. Do not infer that ID from the project ID. Only an authenticated OneMind administrator can access the summary. Restrict account access to the same authorized portal users; per-user Zoho delegation is not implemented.

Until credentials and provider entitlement are confirmed, the card displays configuration required. That state is not a live account summary. No account data has been verified without authorized provider access.

Source: https://projects.zoho.com/api-docs

## Bounded work queue and extension contract

Zoho Projects now reads up to five pages of 100 tasks, deduplicates IDs and exposes at most 20 open task titles/due dates, overdue first. Counts describe the observed sample. A complete flag requires explicit `has_next_page=false`; unknown pagination stays partial. Provider permissions may limit visibility even with complete pagination. Dates use UTC calendar comparisons. Titles are rendered as text, never HTML; descriptions, assignee identities, emails and raw provider objects are omitted.

The five-minute per-process cache shares concurrent requests and retains the last successful result if any page fails. Such results show stale timestamps; failed partial reads never replace the successful snapshot. Two replicas can each make one bounded refresh per interval. Credentials remain server-side; `/api/onemind/cloud` retains authenticated administrator-only access and denies machine identities.

`cloudTypes.ts` defines `CloudProvider` and the sanitized `CloudSummary` contract for future adapters. Each adapter must independently enforce read-only scopes, bounded pagination, source/observation times, credential-required and unavailable states, caching and safe launch links. Adding a provider requires an explicit server-side route/registry integration and administrator authorization; the interface alone does not activate any account.

Activation still requires US-region Zoho Projects OAuth credentials with `ZohoProjects.tasks.READ`, the numeric portal ID and provider API entitlement. No account observation has been made as part of this PR. Official schema: https://projects.zoho.com/api-docs . Do not put tokens in chat or Git.

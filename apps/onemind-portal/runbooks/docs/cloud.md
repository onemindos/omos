# Cloud summaries

The first adapter reads Zoho Projects V3 tasks using server-side OAuth. It requests at most the first 100 tasks and explicitly labels partial results; sample counts are not full portal totals. Results are cached for five minutes, concurrent requests are shared, failures are rate limited, and previous results remain visibly stale. The original app remains the full workflow.

## Enable Zoho Projects

Create a US-region OAuth client and grant `ZohoProjects.tasks.READ` for the intended portal. Confirm your subscription/account permits API access in the provider console. Store `ZOHO_CLIENT_ID`, `ZOHO_CLIENT_SECRET`, `ZOHO_REFRESH_TOKEN` and the numeric `ZOHO_PROJECTS_PORTAL_ID` in a SealedSecret producing `backstage-zoho` in namespace backstage. Restart through a reviewed GitOps update after secret creation. Never commit plaintext credentials or put them in chat. This adapter is fixed to US endpoints; other regions require a reviewed endpoint change.

The link's portal slug is not the numeric API portal ID. Do not infer that ID from the project ID. Only an authenticated Backstage user can access the summary. Restrict account access to the same authorized portal users; per-user Zoho delegation is not implemented.

Until credentials and provider entitlement are confirmed, the card displays configuration required. That state is not a live account summary. No account data has been verified without authorized provider access.

Source: https://projects.zoho.com/api-docs

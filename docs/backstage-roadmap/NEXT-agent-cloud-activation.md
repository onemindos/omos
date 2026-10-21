# Agent and cloud source activation

## Scope

Review actual NATS monitoring/agent protocol and existing Zoho read-only grants. Connect only verified sources; identify external authorization dependencies explicitly. Expose source freshness and unavailable/configuration-required states. Calendar and agent completion must not be fabricated.

## Delivery

Use official Backstage documentation and installed package types. Keep private Tailscale access, signed Authentik identity and server-side credential boundaries. Preserve existing favorites and saved layouts. Build and review implementation in this PR; record concrete acceptance and remaining activation requirements before marking ready.

## Acceptance

- TypeScript and production build.
- Appropriate authenticated runtime behavior and source freshness.
- Light/dark and responsive layout.
- No fabricated health, contracts, events or costs.

## Verified source and implementation

The fixed NATS headless monitoring endpoint responds to /varz with an actual server ID, version and counters. The installed AGENTS stream contains heartbeat/status/discovery subjects; no completion contract was verified. A bounded, cached, read-only /api/onemind/fabric response now exposes allowlisted server counters and JetStream aggregates only. It does not expose messages, account data, client identities or subjects, and labels its single-server scope. Today and the optional native Home widget reuse this source.

Agent completion ingestion remains disabled until an actual producer contract, least-privilege consumer identity and delivery acceptance exist. Do not invent completion subjects from heartbeat activity. The existing notification event bridge can receive approved producers after activation.

The backstage-zoho Secret is absent. Existing Zoho Projects task support requires a US-region tasks.READ OAuth grant and numeric portal ID. Calendar has no configured grant or adapter. User clarification requested for OAuth activation; credentials must stay out of chat/Git. NATS monitoring is independently buildable and does not depend on that authorization.

## Local Zoho connection preparation

User authorized preparation of the read-only connection. Run python3 services/onemind-portal/scripts/prepare-zoho.py in a local terminal. It prompts privately for a US Self Client grant and stores the refresh configuration at ~/.config/onemind/zoho-readonly.json with mode 600, preserving existing files. The task grant is ZohoProjects.tasks.READ. Future Calendar grants are ZohoCalendar.calendar.READ,ZohoCalendar.event.READ; Calendar still needs selected calendar UID/timezone and actual response acceptance. Preparing credentials does not activate Calendar.

Seal the local JSON as backstage-zoho in namespace backstage through the existing cluster-context/SealedSecrets procedure. Existing deployment references this optional Secret. No OAuth authorization or token exchange has been run by this PR. The user must generate the grant in Zoho; do not request plaintext credentials in chat.

Sources: https://projects.zoho.com/api-docs and https://www.zoho.com/calendar/help/api/get-events-list.html .

## Calendar implementation prepared

A selected-calendar read-only adapter follows the official events-list endpoint with an explicit seven-day UTC range and byinstance=true. It uses event.READ, exposes at most 20 sanitized event titles/start times, redacts private or unspecified-privacy titles, omits attendees/descriptions/identities and never invents a timezone. It labels bounded coverage and unknown times. Access requires an administrator user. Five-minute caching and bounded response/deadline limits apply. The optional selected ZOHO_CALENDAR_UID and shared read-only refresh grant are stored locally by the preparation helper and sealed through the existing workflow.

The adapter and native Home/Today cards are implemented, but no real Calendar account response has been accepted yet. Activation remains pending the user's OAuth grant and selected UID. Projects V3's documented ISO end_date is now interpreted as its explicitly recorded calendar date.

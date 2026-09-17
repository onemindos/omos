# Azure Credit Extension — Active Case

**Ticket:** #2605150040000627
**Status:** Open — waiting for account ID confirmation
**Engineer:** Mohammed Akram — Microsoft for Startups Support
**Expiry:** May 19, 2026
**Credits remaining:** ~$3,800
**Subscription:** Legacy AI (`acba54ab-f643-4db6-9f77-57c97aac244f`)

---

## Step 1 — Reply to Akram (DO THIS NOW)

Reply to the support email from Mohammed Akram referencing ticket #2605150040000627:

---

> Hi Mohammed,
>
> Thank you for looking into this. The account the benefits were redeemed under is:
>
> - **Email:** zeus.ops@legacyhq.ai
> - **Subscription Name:** Legacy AI
> - **Subscription ID:** acba54ab-f643-4db6-9f77-57c97aac244f
> - **Tenant ID:** 5368e6c7-0948-4428-836a-921d1c2b8003
> - **Tenant Domain:** zeusopslegacyhq.onmicrosoft.com
>
> Note: zeus@zeusdelacruz.com may be a separate personal account. The active Azure subscription with the sponsorship credits is under zeus.ops@legacyhq.ai.
>
> The credits were activated approximately November 2025 — less than 6 months ago — and $3,800 remains. I am requesting a 3–6 month extension to allow time to fully utilize these credits for the production AI agent infrastructure I am building on Azure Foundry and Azure OpenAI.
>
> Please let me know if you need anything else.
>
> Thank you,
> Zeus Delacruz

---

## Step 2 — Microsoft for Startups Portal (parallel path)

Go to: https://aka.ms/startuphelp-mfs-portal

1. Sign in (try both zeus.ops@legacyhq.ai AND zeus@zeusdelacruz.com to find which has the benefit)
2. Click **"?" icon** in the top-right corner
3. Select **"Submit a Support Ticket"**
4. Use this message:

---

> **Subject:** Request to extend Azure sponsorship credit expiration — reference Azure ticket #2605150040000627
>
> Hi,
>
> I have an open Azure support ticket (#2605150040000627) requesting an extension of my Microsoft for Startups Azure sponsorship credits expiring May 19, 2026.
>
> I have ~$3,800 remaining in credits that were activated approximately November 2025 (less than 6 months ago). I am building a production AI agent platform (One Mind) on Azure OpenAI and Azure AI Foundry and need 3–6 additional months to fully utilize these credits.
>
> - Subscription: Legacy AI
> - Subscription ID: acba54ab-f643-4db6-9f77-57c97aac244f
> - Email: zeus.ops@legacyhq.ai
> - Azure ticket: #2605150040000627
>
> Requesting extension to August–November 2026.
>
> Thank you,
> Zeus Delacruz

---

## Account Details (for reference)

| Field | Value |
|-------|-------|
| Subscription name | Legacy AI |
| Subscription ID | acba54ab-f643-4db6-9f77-57c97aac244f |
| Tenant ID | 5368e6c7-0948-4428-836a-921d1c2b8003 |
| Tenant domain | zeusopslegacyhq.onmicrosoft.com |
| Azure login | zeus.ops@legacyhq.ai |
| MFS portal email (check) | zeus@zeusdelacruz.com OR zeus.ops@legacyhq.ai |
| Credit expiry | May 19, 2026 |
| Credits remaining | ~$3,800 |
| Activation date | ~November 2025 |
| Requested extension | 3–6 months (Aug–Nov 2026) |

---

## What to do if extension is denied

If no extension is granted before May 19, maximize spend immediately:

1. **Deploy Azure OpenAI in Foundry** — GPT-5.5 or GPT-4o (East US 2 or Sweden Central)
2. **Buy Azure OpenAI PTU reservations** — 1-year reserved throughput counts as immediate credit consumption
3. **Spin up Azure VMs** for omos infrastructure — burns credits on compute
4. **Use `azure-openai-responses` provider in OpenClaw** — route all agent traffic through Azure

---

## CLI Commands (for monitoring)

```bash
# Check ticket status
az support in-subscription tickets show \
  --subscription "acba54ab-f643-4db6-9f77-57c97aac244f" \
  --ticket-name "MFS-Credit-Extension-May2026"

# List all tickets
az support in-subscription tickets list \
  --subscription "acba54ab-f643-4db6-9f77-57c97aac244f" \
  --output table
```

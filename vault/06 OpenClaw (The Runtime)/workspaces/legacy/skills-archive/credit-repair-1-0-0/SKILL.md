# Credit Repair Skill

Generate FCRA-compliant dispute letters for credit bureau errors. Produces HTML letters ready to send via the **lob** skill (physical mail via USPS).

## When to Use

- User wants to dispute errors on their credit report
- User needs a dispute letter written for a credit bureau
- User asks about their FCRA rights or dispute timelines
- User wants to follow up on a "verified" dispute
- User needs to send a debt validation request to a collector

## Workflow

1. **Generate letter** — `credit_repair({ action: "generate_letter", ... })` → HTML output
2. **Send via Lob** — `lob({ action: "send_letter", content: html, to_name: bureau.name, ... })` → tracked USPS mail
3. **Track delivery** — `lob({ action: "get_letter", letter_id: "..." })`
4. **Follow up** — Generate method_of_verification letter if bureau claims "verified" without proof

## Letter Types

| Type                     | When to Use                                       |
| ------------------------ | ------------------------------------------------- |
| `general_dispute`        | First dispute — report an inaccuracy              |
| `method_of_verification` | Follow-up after bureau claims item was "verified" |
| `debt_validation`        | Challenge a collection agency to prove the debt   |

## Bureau Addresses (Built-in)

| Bureau                | Address                             |
| --------------------- | ----------------------------------- |
| `experian`            | P.O. Box 4500, Allen, TX 75013      |
| `equifax`             | P.O. Box 740256, Atlanta, GA 30374  |
| `transunion`          | P.O. Box 2000, Chester, PA 19016    |
| `transunion_disputes` | P.O. Box 34012, Fullerton, CA 92834 |

## Quick Reference

### Generate a dispute letter

```
credit_repair({
  action: "generate_letter",
  letter_type: "general_dispute",
  bureau: "experian",
  sender_name: "...", sender_address: "...", sender_city: "...", sender_state: "...", sender_zip: "...",
  ssn_last4: "1234",
  creditor_name: "Chase Bank",
  account_last4: "5678",
  dispute_reason: "Incorrect Payment Status",
  explanation: "Payment was made on time on March 1, 2025 but shows as 30 days late"
})
```

### Then send it

```
lob({
  action: "send_letter",
  to_name: "Experian Disputes",
  to_address_line1: "P.O. Box 4500",
  to_city: "Allen", to_state: "TX", to_zip: "75013",
  from_name: "...", from_address_line1: "...", from_city: "...", from_state: "...", from_zip: "...",
  content: "<html>...(generated HTML)...</html>"
})
```

### Get valid dispute reasons

```
credit_repair({ action: "list_dispute_reasons" })
```

### FCRA rights reference

```
credit_repair({ action: "fcra_rights" })
```

## Important Rules

- **Never include full SSN** — only last 4 digits
- **One dispute per letter** for clarity and legal strength
- **Customize the explanation** — don't use template language verbatim
- **Track everything** — save letter IDs, set 30-day reminders
- **Bureau has 30 days** to investigate (45 with additional docs)
- **No env vars needed** — this skill generates letters locally; the lob skill handles mailing


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

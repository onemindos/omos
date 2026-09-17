# HEARTBEAT.md — Guardian Proactive Checklist

## On Every Heartbeat

Run these checks in order. Stop as soon as something needs attention and surface it.

### 1. Security Alerts

- Any unusual access patterns or login attempts?
- Credential expiration warnings?
- SSL/certificate status?
- Domain registration expiry approaching?

### 2. Legal Deadlines

- Any court dates, filing deadlines, or response due dates within 30 days?
- Escalate anything within 7 days to Legacy
- Track deadline progression (30-day, 7-day, 3-day alerts)

### 3. Insurance

- Any policy renewals within 30 days?
- Outstanding claims needing follow-up?
- Coverage gap analysis (quarterly)

### 4. Credentials & Access

- API keys or tokens expiring within 30 days?
- Password rotation schedule adherence?
- Service account audit (quarterly)

### 5. Compliance

- Any regulatory filings due?
- Business license renewals approaching?
- Tax filing deadlines?

### 6. Nothing To Report

If all clear: reply `HEARTBEAT_OK`

## Schedule

- Run every 2 hours during active hours (8am-10pm Eastern)
- Critical security alerts: monitor 24/7 via system integrations
- Never spam the same alert twice — track what was sent

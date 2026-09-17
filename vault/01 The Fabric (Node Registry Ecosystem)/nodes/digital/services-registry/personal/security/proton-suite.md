---
title: "Proton Suite"
type: reference
status: active
tags: [platform, personal, security, privacy, vpn, encrypted]
---
# Proton Suite

> Full privacy suite — Swiss-based, end-to-end encrypted, zero-access.

## Products
| Product | Status | Notes |
|---------|--------|-------|
| Proton Mail | Active | E2E encrypted email (see email/proton-mail.md) |
| Proton Drive | Active | E2E encrypted storage (see storage/proton-drive.md) |
| Proton VPN | Active | No-log VPN, Swiss jurisdiction |
| Proton Calendar | Active | Encrypted calendar |
| Proton Pass | Active | Password manager (alternative to 1Password) |

## Proton VPN
```bash
# CLI (Linux)
protonvpn-cli connect --fastest
protonvpn-cli connect US
protonvpn-cli disconnect
protonvpn-cli status
```

| Feature | Notes |
|---------|-------|
| No-log | Verified by audit |
| Kill switch | Blocks traffic if VPN drops |
| Split tunneling | Choose which apps use VPN |
| Secure Core | Route through multiple countries |
| P2P | Torrent-friendly servers |
| NetShield | Built-in ad/malware blocker |

## Proton Calendar
- E2E encrypted events
- Import from Google Calendar
- Shareable calendars (encrypted)

## Proton Pass
- E2E encrypted passwords + notes
- Browser extension + mobile app
- 2FA/TOTP built-in
- Aliases (hide-my-email)

## Pricing
Proton Unlimited: ~$10/mo (all products bundled)

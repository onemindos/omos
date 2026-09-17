# Proton Mail

> **End-to-End Encrypted Email for Privacy**

## Overview

| Attribute    | Value                  |
| ------------ | ---------------------- |
| Type         | Secure Email Provider  |
| Website      | proton.me              |
| Encryption   | End-to-End (PGP-based) |
| Jurisdiction | Switzerland            |
| Status       | ✅ Active              |

## What is Proton Mail?

Proton Mail is an encrypted email service founded by CERN scientists, based in Switzerland. It offers end-to-end encryption, meaning only the sender and recipient can read messages—not even Proton can access your emails.

## Why Use Proton Mail

| Benefit                | Description                    |
| ---------------------- | ------------------------------ |
| Zero-Access Encryption | Proton cannot read your emails |
| Swiss Privacy Laws     | Strong legal protection        |
| No IP Logging          | Enhanced anonymity             |
| Open Source            | Auditable code                 |
| No Ads                 | Privacy-focused business model |

## Key Features

### Email Security

| Feature                   | Description                     |
| ------------------------- | ------------------------------- |
| End-to-End Encryption     | Automatic between Proton users  |
| Password-Protected Emails | Encrypt to non-Proton users     |
| Expiring Messages         | Auto-delete after set time      |
| PGP Support               | Compatible with other PGP users |
| Two-Factor Auth           | TOTP and security keys          |

### Organization Features

| Feature             | Description               |
| ------------------- | ------------------------- |
| Custom Domains      | Use your own domain       |
| Catch-All Addresses | Receive all domain emails |
| Multiple Addresses  | Aliases for organization  |
| Folders & Labels    | Email organization        |
| Filters & Rules     | Automatic sorting         |

### Proton Ecosystem

| Product         | Description            |
| --------------- | ---------------------- |
| Proton Mail     | Encrypted email        |
| Proton Calendar | Encrypted calendar     |
| Proton Drive    | Encrypted file storage |
| Proton VPN      | Secure VPN service     |
| Proton Pass     | Password manager       |

## Plans

| Plan             | Features                          |
| ---------------- | --------------------------------- |
| Free             | 1GB storage, 1 address            |
| Mail Plus        | 15GB, custom domain, 10 addresses |
| Proton Unlimited | 500GB, all Proton apps            |
| Business         | Team features, admin controls     |

## Use in OneMind Stack

### Personal Email

- Primary personal secure email
- Sensitive communications
- Privacy-focused correspondence
- Separate from business email (Zoho)

### Integration Considerations

| Integration         | Status                         |
| ------------------- | ------------------------------ |
| IMAP/SMTP           | ✅ Via Proton Bridge (desktop) |
| Mobile Apps         | ✅ iOS & Android native apps   |
| Web Access          | ✅ mail.proton.me              |
| Third-Party Clients | Requires Proton Bridge         |
| API Access          | Limited (security reasons)     |

### Proton Bridge

For desktop email clients:

```bash
# Proton Bridge enables IMAP/SMTP
# Required for: Thunderbird, Apple Mail, Outlook
# Runs locally, decrypts emails for client access
```

## Email Strategy

| Domain/Purpose      | Provider                      |
| ------------------- | ----------------------------- |
| Personal (Private)  | Proton Mail                   |
| Business (Zoho One) | Zoho Mail                     |
| Transactional       | Various (Zoho, AWS SES, etc.) |

## Security Best Practices

1. **Enable 2FA** - Use authenticator app or security key
2. **Recovery Methods** - Set up recovery email/phrase
3. **Check Encryption Status** - Green lock = encrypted
4. **Use Aliases** - Protect primary address
5. **Regular Key Rotation** - Update PGP keys periodically

## Resources

- [Proton Mail](https://proton.me/mail)
- [Proton Bridge](https://proton.me/mail/bridge)
- [Security Model](https://proton.me/mail/security)
- [Business Plans](https://proton.me/business)

---

Last Updated: January 25, 2026

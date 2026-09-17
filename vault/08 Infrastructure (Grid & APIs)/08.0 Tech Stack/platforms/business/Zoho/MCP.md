# Zoho MCP (Mobile Device Management)

> **Mobile Control Panel & Device Management**

## Overview

Zoho MCP (Mobile Control Panel) is a mobile device management (MDM) solution for managing, securing, and monitoring mobile devices across your organization.

## Status

- **Active**: ✅ Yes
- **Integration**: Agent-based
- **Last Updated**: January 25, 2026

## Use Cases

1. **Device Management** - Enroll and manage mobile devices
2. **Security Enforcement** - Apply security policies
3. **App Management** - Deploy and control apps
4. **Remote Wipe** - Erase lost/stolen devices
5. **Location Tracking** - Find device locations

## Key Features Used

- [ ] Device enrollment (iOS, Android)
- [ ] Security policies
- [ ] App deployment (whitelist/blacklist)
- [ ] Remote lock/wipe
- [ ] Device location tracking
- [ ] Compliance monitoring
- [ ] Kiosk mode
- [ ] Email configuration

## Device Enrollment

### Supported Platforms

| Platform    | Enrollment Method           | Status       |
| ----------- | --------------------------- | ------------ |
| **iOS**     | Apple DEP, Manual, QR Code  | ✅ Supported |
| **Android** | Android Enterprise, QR Code | ✅ Supported |
| **macOS**   | Apple DEP, Manual           | ✅ Supported |
| **Windows** | Azure AD, Manual            | ✅ Supported |

### Enrollment Workflow

1. Admin sends enrollment invite (email/SMS)
2. User downloads MDM app
3. User follows setup prompts
4. Device registers with MCP server
5. Policies automatically applied
6. Device ready for management

### Apple DEP (Device Enrollment Program)

For zero-touch iOS enrollment:

1. Purchase devices through Apple Business Manager
2. Link Apple Business Manager to Zoho MCP
3. Assign devices to Zoho MCP
4. Devices auto-enroll on first boot
5. Skip setup screens (optional)

## Security Policies

### Password Policies

| Setting                 | Value                 | Platform |
| ----------------------- | --------------------- | -------- |
| **Minimum Length**      | 8 characters          | All      |
| **Complexity**          | Alphanumeric required | All      |
| **Expiration**          | 90 days               | All      |
| **History**             | Last 5 passwords      | All      |
| **Auto-Lock**           | 5 minutes             | All      |
| **Max Failed Attempts** | 10 (then wipe)        | All      |

### Device Restrictions

| Restriction       | iOS | Android | Purpose                |
| ----------------- | --- | ------- | ---------------------- |
| **Camera**        | ✅  | ✅      | Disable camera         |
| **Screenshots**   | ✅  | ✅      | Prevent screen capture |
| **USB Transfer**  | ✅  | ✅      | Block file transfer    |
| **Bluetooth**     | ✅  | ✅      | Disable Bluetooth      |
| **App Install**   | ✅  | ✅      | Block app installs     |
| **Browser**       | ✅  | ✅      | Restrict web access    |
| **iCloud Backup** | ✅  | N/A     | Disable iCloud         |
| **Factory Reset** | ✅  | ✅      | Block reset            |

### Encryption

- iOS: FileVault (automatic, hardware-based)
- Android: Full-disk encryption (required)
- BitLocker: Windows devices

## App Management

### App Distribution

| Distribution Type   | Description                           | Use Case             |
| ------------------- | ------------------------------------- | -------------------- |
| **Store Apps**      | Public apps from App Store/Play Store | Required apps        |
| **Enterprise Apps** | Private/internal apps                 | Custom business apps |
| **Web Clips**       | Web shortcuts on home screen          | Web apps             |
| **VPP Apps**        | Volume-purchased apps (Apple VPP)     | Paid apps for org    |

### App Policies

**Whitelist Mode**

Only allowed apps can be installed:

```
Allowed Apps:
- Slack
- Zoom
- Microsoft Office
- Zoho Mail
- OneMind (internal)
```

**Blacklist Mode**

Specific apps are blocked:

```
Blocked Apps:
- TikTok
- Facebook
- Instagram
- Games (category)
```

### Managed App Configuration

Push settings to managed apps:

```json
{
  "app_id": "com.zoho.mail",
  "config": {
    "email": "{{user.email}}",
    "server": "mail.zoho.com",
    "auto_configure": true
  }
}
```

## Remote Actions

### Available Remote Actions

| Action             | Description                 | Use Case           |
| ------------------ | --------------------------- | ------------------ |
| **Lock**           | Immediately lock device     | Lost device        |
| **Wipe**           | Factory reset device        | Stolen device      |
| **Corporate Wipe** | Remove only work data       | Employee departure |
| **Locate**         | Show current location       | Find device        |
| **Ring**           | Play sound (even if silent) | Misplaced device   |
| **Reset Password** | Clear device password       | Forgotten password |
| **Restart**        | Remotely reboot device      | Troubleshooting    |
| **Sync**           | Force policy refresh        | Apply changes      |

### Wipe Options

**Full Wipe**

- Factory reset entire device
- All data deleted (personal + work)
- Device deregisters from MDM
- Use for: Stolen devices, severe policy violations

**Corporate Wipe (Selective)**

- Remove only work-related data
- Personal data preserved
- MDM profile removed
- Use for: BYOD devices, employee offboarding

## Location Tracking

### Tracking Features

- Real-time location (GPS)
- Location history (last 30 days)
- Geofencing (enter/exit alerts)
- Map view (all devices)
- Lost mode (iOS - locks and displays message)

### Geofencing

Create virtual boundaries:

**Office Geofence**

- Location: 123 Main St, NYC
- Radius: 500 meters
- Trigger: Enter/Exit

**Actions on Exit**

- Alert admin
- Lock device
- Enable camera restriction
- Disable work apps

### Privacy Considerations

- Location tracking requires user consent
- Can be disabled for personal devices
- Track only during work hours (optional)
- Data retention: 30 days (configurable)

## Compliance Monitoring

### Compliance Rules

| Rule                         | Description               | Action if Non-Compliant   |
| ---------------------------- | ------------------------- | ------------------------- |
| **Jailbreak/Root Detection** | Device is compromised     | Block access, alert admin |
| **Encryption Required**      | Disk encryption enabled   | Prompt to enable          |
| **Password Strength**        | Meets policy requirements | Force password change     |
| **OS Version**               | Minimum OS version        | Notify to update          |
| **App Blacklist**            | Blocked apps installed    | Notify/remove app         |
| **MDM Profile**              | MDM installed & active    | Revoke access             |

### Compliance Dashboard

- % of compliant devices
- Non-compliant device list
- Compliance trend over time
- Risk score per device

## Email Configuration

### Push Email Settings

Auto-configure email on devices:

```json
{
  "account_name": "Work Email",
  "email": "{{user.email}}",
  "server": "mail.zoho.com",
  "port": 993,
  "ssl": true,
  "username": "{{user.email}}",
  "password": "{{prompt}}"
}
```

### Supported Protocols

- IMAP/SMTP
- Exchange ActiveSync
- Gmail (OAuth)
- Zoho Mail (OAuth)

## Kiosk Mode

Lock device to single or limited apps:

### Single-App Kiosk

Device runs only one app:

- Point of Sale (PoS) terminal
- Customer feedback kiosk
- Digital signage
- Check-in station

### Multi-App Kiosk

Device runs allowed apps only:

- Field worker tablet (maps, inventory, forms)
- Student device (education apps only)
- Shared device (limited apps)

### Kiosk Settings

```
Enabled Apps: Zoho CRM, Zoho Forms, Camera
Disabled Features:
- Status bar
- Navigation buttons
- Settings access
- Home button
- Power button (requires PIN to exit)
```

## Device Groups

### Group Types

- **By Department**: Sales, Engineering, Support
- **By Location**: NYC, LA, Remote
- **By Device Type**: iOS, Android, Laptops
- **By Role**: Executives, Field Workers, Contractors
- **By Enrollment**: Corporate, BYOD

### Group Policies

Each group can have different policies:

**Executive Devices**

- Minimal restrictions
- Remote wipe enabled
- Location tracking disabled

**Field Worker Devices**

- App whitelist enforced
- Location tracking enabled
- Kiosk mode available

**BYOD Devices**

- Work profile only managed
- No location tracking
- Selective wipe only

## Reporting & Analytics

### Available Reports

- **Device Inventory**: All enrolled devices
- **Compliance Status**: Per-device compliance
- **App Inventory**: Installed apps across fleet
- **Location History**: Device movement
- **Security Events**: Policy violations, lost devices
- **Enrollment Trend**: New/removed devices over time

### Alerting

Alerts for:

- New device enrolled
- Device non-compliant
- Jailbreak/root detected
- Device lost (marked by user)
- Policy violation
- App blacklist detected

## API Integration

### Get All Devices

```python
GET https://mdm.zoho.com/api/v1/devices

Response:
{
  "devices": [
    {
      "id": "123",
      "name": "Zeus's iPhone",
      "platform": "iOS",
      "os_version": "17.2",
      "compliant": true,
      "last_sync": "2026-01-25T10:00:00Z"
    }
  ]
}
```

### Remote Lock Device

```python
POST https://mdm.zoho.com/api/v1/devices/{id}/lock
{
  "message": "This device has been locked by IT."
}
```

### Deploy App

```python
POST https://mdm.zoho.com/api/v1/apps/deploy
{
  "app_id": "com.zoho.crm",
  "device_ids": ["123", "456", "789"]
}
```

## Integrations

- **Zoho CRM**: Manage sales team devices
- **Zoho People**: Employee onboarding/offboarding
- **Zoho Cliq**: Alert notifications
- **Active Directory**: User sync
- **Apple Business Manager**: DEP enrollment
- **Google Admin**: Android Enterprise

## Legacy AI Integration

### Automated Actions

1. **New Employee** (Zoho People) → Auto-enroll device, push apps
2. **Employee Termination** → Corporate wipe device
3. **Lost Device Report** → Lock device, enable tracking
4. **Non-Compliance Alert** → Notify team in Cliq
5. **App Request** → Approve/deploy via chatbot

### Use Cases

- "Lock my work phone" → Trigger remote lock
- "Report my device stolen" → Lock + enable lost mode
- "Deploy CRM to sales team" → Push app to group
- "What devices are non-compliant?" → Generate report

## Best Practices

- ✅ Use DEP/Android Enterprise for corporate devices (seamless enrollment)
- ✅ Separate BYOD and corporate policies
- ✅ Enable encryption on all devices
- ✅ Test policies on pilot group first
- ✅ Document MDM enrollment in employee onboarding
- ✅ Regular compliance audits (weekly)
- ✅ Train users on security practices
- ✅ Plan for lost/stolen device procedures

## Notes

- Free tier: Up to 25 devices
- Paid plans: $1.50/device/month (Standard), $3/device/month (Enterprise)
- Apple DEP requires Apple Business Manager account
- Android Enterprise requires Google Admin account
- Location tracking subject to regional privacy laws (GDPR compliance available)
- Data stored in Zoho data centers (region-specific)

---

**Related**: [Cliq](./Cliq.md) | [Mail](./Mail.md) | [CRM](./CRM.md)

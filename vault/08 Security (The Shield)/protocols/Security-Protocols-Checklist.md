---
id: security-protocols-master
title: "Security Protocols & Checklist"
type: reference
status: active
created: 2026-01-19
modified: 2026-01-19
owner: zeus
tags: [security, protocols, checklist, reference]
---

# Security Protocols & Checklist

> **Complete security framework** - 253 items across 11 categories
> Migrated from Todoist - 2026-01-19

---

## Account Security (28 items)

- [ ] 2FA
- [ ] AUTOMATED - (NO COOKIES) - (NO PASSWORD SAVED) - (NO SAVED INFO)
- [ ] CONNECT W FINANCIAL &amp; INVESTMENT ACCOUNTS  - ICLOUD, CREDIT ACCOUNTS
- [ ] CREDIT MONITORING - IDENITITY THEFT - EXPERIAN
- [ ] FIREFOX SETTINGS &amp; MONITOR
- [ ] FREEZE AND MONITOR YOUR CREDIT (IDENTITY FORCE OR EXPERIAN)
- [ ] IMPLEMENT TWO FACTOR AUTHENTICATION
- [ ] LIST ALL IMPORTANT ACCOUNTS (CRITICAL)
- [ ] LIST TOP SOCIAL MEDIA ACCOUNTS
- [ ] LOCK YOUR SIM CARD
- [ ] LONG PASSWORD
- [ ] PASSWORD MANAGERS
- [ ] PASSWORD MANAGERS - 1 PASSWORD
- [ ] PASSWORDS
- [ ] PERSONAL IDENTITY
- [ ] PRIVACY.COM &amp; FAST MAIL 1PASSWORD
- [ ] QUESTIONS TO ASK ABOUT YOUR PASSWORDS (USE PASSWORD MANANGER)
- [ ] RATE THE STRENGTH OF CURRENT PASSWORD
- [ ] Reset All Default Passwords Don’t forget out of band management interfaces like IPMI
- [ ] SIM PIN
- [ ] STRONG SELF GENERATED PASSWORD FOR PROTECTION
- [ ] Setup Log Monitoring (Splunk, Grafana, Kibana)
- [ ] USE A PASSWORD MANAGER
- [ ] USE STRONG ACCOUNT PASSWORDS
- [ ] USING PASSWORD MANAGER APP - (1PASSWORD)
- [ ] Use 2FA where possible
- [ ] Use Named User Accounts where practical instead of having the entire family share a login
- [ ] Use a Password Manager such as KeePass or LastPass. Use strong unique passwords for every system

---

## Data Security (18 items)

- [ ] 5 COMMANDMENTS of PRIVACY (how I protect my data)
- [ ] BACKUP (ENCRYPTION ADDED)
- [ ] Backups You must have offsite versioned backups of everything you don’t want to lose
- [ ] Battery Backup APC or CyberPower offer protection against brief power failures
- [ ] CLEAR DATA - CALLS, TEXTS,BROWERS,APPLICATION
- [ ] CLOUD STORAGE
- [ ] CLOUD STORAGE (NAS W BACKUP &amp; S3)
- [ ] Continuity Give your family instructions on who to call and what to do should you perish so they don’t lose any data. Make sure an acquaintance is familiar enough with your homelab
- [ ] DEATH FILE - DRIVE
- [ ] ENCRYPT PERSONAL SENSITIVE INFO
- [ ] ENCRYPT YOUR EMAILS
- [ ] Encrypt all sensitive data at rest. Some systems like FreeNAS can encrypt the disks
- [ ] Encrypting Data at Rest and Transit - How to Protect Your Data
- [ ] SECURE DRIVE IN SAFE
- [ ] SHRED YOUR TRASH
- [ ] TRAVEL THUMB DRIVE W CLOUD
- [ ] USE ENCRYPTED MESSENGERS
- [ ] Use Let’s Encrypt to setup automatic SSL certs (requires a domain name)

---

## Device Security (46 items)

- [ ] 2FA APP ON PHONE - (1PASSWORD)
- [ ] 6 DIGIT IPHONE PASSWORD
- [ ] APPLE DEVICE SECURITY- 3 DEVICES (MAC) IPHONE) (IPAD)
- [ ] APPLE DEVICES ZONE 1 PROTECTION
- [ ] APPLE DEVICES ZONE 2 PROTECTION
- [ ] APPLE DEVICES ZONE 3 PROTECTION
- [ ] APPS
- [ ] AUTOMATED - TIME MACHINE - (ENCRYPTED)
- [ ] AUTOMATIC UPDATES
- [ ] Automatic Security Updates Make sure every device on your network is configured to install security updates automatically. Otherwise check at least once a month
- [ ] BACKUP YOUR DEVICES
- [ ] CELL PHONE PROVIDER --&gt; LINKS
- [ ] CELLULAR
- [ ] CHANGE PHONE AUTH TO YUBI KEY AUTH
- [ ] CLEAN OUT BROWSERING DATA - ALL DEVICES
- [ ] CLOUD BACKUP FOR DEVICES (NOT ICLOUD)
- [ ] DEVICES CLEAN UP
- [ ] DISABLE BIOMETRICS IN RISKY AREAS
- [ ] DISABLE GPS PERMANENTLY
- [ ] DISABLE RADIOS OFTEN (LTE,WIFI,BT,ETC)
- [ ] DISABLE SAFARI ON ALL APPLE DEVICES
- [ ] DO NOT JAILBREAK
- [ ] ENABLE AUTOMATIC UPDATES
- [ ] GPS
- [ ] HIDE YOUR HOME ADDRESS AND PHONE NUMBER (TRAVELING MAILBOX) &amp; (HUSHED)
- [ ] IOS PRIVACY &amp; SECURITY GUIDE
- [ ] IPAD
- [ ] IPHONE
- [ ] IPHONE (WATCH SAVED )
- [ ] MAC
- [ ] MALWARE PROTECTION FOR DEVICES - (BIT-DEFENDER)
- [ ] MANAGE LOCATION FOR ALL DEVICES
- [ ] MANAGING RADIOS
- [ ] NFC - APPLE PAY
- [ ] OBTAIN PHONE SERVICE WITH NO PERSONAL INFORMATION
- [ ] PHYSICAL FINDER FOR IPHONE &amp; IPAD (ENCRYPTED)
- [ ] REMOVE PASSWORDS FROM BROWERS &amp; DEVICES
- [ ] SET 6 DIGIT PASSCODE FOR MOBILE DEVICES
- [ ] SETUP YUBIKEY FOR ALL DEVICES
- [ ] SWITCH TO FOSS APPS
- [ ] SYNC 3 DEVICES
- [ ] Setup Network Monitoring (Nagios, Zabbix, Icinga)
- [ ] The ULTIMATE Mac Privacy &amp; Security Guide
- [ ] USE A STRONG DEVICE PASSWORD
- [ ] WHEN CHOOSING NEW APPS
- [ ] YUBI KEYS SETUP - (3 DEVICES)

---

## Network Security (24 items)

- [ ] ASTERISK FOR BUILDING COMMUNICATIONS APPLICATIONS
- [ ] BLUETOOTH
- [ ] DMZ Every public server should be on a network separate from your LAN
- [ ] DOUBLE CHECK YOUR HOME NETWORK
- [ ] HIDE HOME NETWORK - CREATE GUEST NETWORK
- [ ] IDSIPS Setup an Intrusion Detection  Prevention System. Some Firewalls (such as UniFi’s Security Gateway) have this built in
- [ ] LARGE LAND - NIGHT VISION CAMERAS - ON PRIVATE SERVER
- [ ] NETWORK SWITCH WITH V-LANS
- [ ] Register a real domain for your Homelab, NameCheap is a good registrar. That way you can access all your servers using a proper FQDN. Setup DNSSEC and DMARC
- [ ] SECURING HOME NETWORK
- [ ] SELF HOSTED VPN
- [ ] SSH Key Authentication Disable password authentication
- [ ] Secure DNS. Use CleanBrowsing or OpenDNS to prevent bad sites from loading
- [ ] SecurityVulnerability Scanner Regularly scan your LAN and public IP with Cloudflare’s Flan Scan to find open ports and known vulnerabilities
- [ ] Separate Guest Wi-Fi Don’t let guests onto your LAN when not needed
- [ ] Setup VLANS on Managed Switches to separate traffic DMZ, LAN, and Guest is good for a homelab
- [ ] USE A PRIVATE DNS (MAY BE VPN)
- [ ] USE A VPN
- [ ] VPN
- [ ] VPN &amp; FIREWALL - BRAVE
- [ ] VPN If you need access into your LAN remotely setup a VPN. OpenVPN is very popular, ZeroTier requires no ports to be open on your firewall
- [ ] VPN SERVICES
- [ ] WIFI
- [ ] Web Application Firewall If you have a public webserver at your house, put it behind Cloudflare for a little protection. Plus, it acts as a CDN giving you a performance boost

---

## Physical Security (16 items)

- [ ] AT HOME SECURITY SYSTEM
- [ ] CAMERA COVERS
- [ ] COVER YOUR CAMERAS
- [ ] Don’t forgot Physical Security. Get a dog
- [ ] FARADAY BAGS - SIGNAL BLOCKER
- [ ] HIDING HOME ADDRESS
- [ ] HIDING HOME ADDRESS (DO IT WHEN MOVING)
- [ ] MAILBOX MANAGEMENT
- [ ] PHYSICAL GEAR
- [ ] PRIVACY SCREEN PROTECTOR
- [ ] SLIDING CAMERA CASES (CAMERA PROTECTION)
- [ ] SNIP CAMERAS AND MICS
- [ ] TREND MICRO SERVICE - FAMILY
- [ ] USE A PRIVACY SCREEN PROTECTOR
- [ ] USE FARADAY POUCHES AND BAGS
- [ ] VIRTUAL MAILBOX - (TRAVELING MAILBOX)

---

## Privacy (30 items)

- [ ] 7 Critical Facebook Privacy Settings You Need to Change TODAY (2022 Update)
- [ ] AD BLOCKERS
- [ ] AUTOMATED COOKIE CLEANER
- [ ] BRAVE SETTING
- [ ] CHANGE APPLE ID EMAIL TO NONE KNOWN (PROTON)
- [ ] COOKIES
- [ ] DISABLE UNUSED PERMISSIONS SETTINGS
- [ ] DUCKDUCKGO SETTINGS
- [ ] FOSS APPLICATIONS (FREE AND OPEN SOURCE SOFTWARE)
- [ ] GOOGLE VOICE --&gt; FOR PROTECTION FORWARDING
- [ ] HISTORY
- [ ] HUSHED --&gt;FOR PROTECTION FORWARDING
- [ ] INSTALL A BROWSER FOR TOR USAGE
- [ ] MAKE SOCIAL MEDIA PRIVATE
- [ ] OPT OUT - FROM EVERYTHING
- [ ] PRIVACY SETTINGS
- [ ] PRIVACY TRUST - $5K
- [ ] PRIVACY.COM
- [ ] PRIVATE SEARCH ENGINE
- [ ] PRIVATE SEARCH ENGINES
- [ ] PROTON CALENDAR FOR PERSONAL CALENDAR
- [ ] Removing Your Old Personal Information from the Internet
- [ ] SIGNAL - PRIVATE MESSENGER
- [ ] SOCIAL MEDIA SECURITY
- [ ] Track Mecury Personal Information Update
- [ ] USE - IMESSAGE OR FACETIME
- [ ] USE A PRIVACY SEARCH ENGINE
- [ ] WEB BROWSERS
- [ ] WEB BROWSERS - (FOLLOW ASSIGNED CHARTED) WATCH VIDEOS
- [ ] httpswww.allthingssecured.comtutorialsfacebook-privacy-settings

---

## Uncategorized (91 items)

- [ ] 1Password Maintenance - Watch Tower (Weekly or Monthly)
- [ ] 1Password Watch Tower
- [ ] 2FA
- [ ] 2FA FOR CLICKUP - BUSINESS PLAN NEEDED
- [ ] 2FA Tracking & Management
- [ ] 2FA/OTP Management - Yubi keys & 1Password
- [ ] AWS Site To Site VPN - New video with improved steps (Part 1)
- [ ] All Apple Devices
- [ ] All Apple Devices With Finder
- [ ] All Web Browsers
- [ ] Apple Devices
- [ ] Apps - All Apple Devices
- [ ] BANK STATEMEMTS
- [ ] BEST SECURE MESSAGING SERVICES
- [ ] Backups & Disaster Management (Daily & Weekly) - (Physical & Digital)
- [ ] Bookmarks - Raindrop
- [ ] Business Protocols
- [ ] CP FOR UTILITIES
- [ ] CRITICAL SECURITY CHANGES
- [ ] Cloud drives to Nas? or Nas to Cloud drives?
- [ ] DELETE DIGITAL FOOTPRINT SECURITY & MAINTAIN / DO INVESTIGATION ON SELF TO CLEAN UP
- [ ] DISABLE ICLOUD &amp; CLOUD FEATURES
- [ ] DNS Management
- [ ] Device Management (Apple)
- [ ] Device Management (Apple)
- [ ] Domain Management - (Go daddy)
- [ ] EMPLOYEE VM'S - PASSWORD MANAGER W/YUBI KEY ACCESS
- [ ] Email - Proton & Zoho Mail
- [ ] Experian Idworks Portal
- [ ] Family Plan - kids identity watch
- [ ] Files - Local Folders/Physical Drives & Cloud
- [ ] GLOBAL COVERAGE
- [ ] HOME LAB SECURITY CHECKLIST
- [ ] Home Security
- [ ] Home Security
- [ ] ICLOUD DISABLE
- [ ] IMPLEMENT MINIMALISM (LESS IS MORE)
- [ ] Identity Protection
- [ ] Idrive
- [ ] Install Fail2Ban on all Linux Servers
- [ ] Iphone 17 Pro Max (AT&T)
- [ ] LOCAL - US,MEXICO ,CANADA
- [ ] MAKE NEW PERSONAL GAMIL EMAIL
- [ ] MANAGE ADD ONS
- [ ] MacBook Pro M2 Max (16 Inch)
- [ ] NEVER CONNECT PERSONAL
- [ ] NO - BIO METRICS
- [ ] Norton
- [ ] Organize and Keep track of (Weekly)
- [ ] PHOTOS
- [ ] PRIVACY.COM FOR BUSINESS
- [ ] PROTON SECURE BUSINESS EMAIL
- [ ] Papers - Scan
- [ ] Passkey Management
- [ ] Personal Protocols
- [ ] Personal Security
- [ ] Personal Security
- [ ] Physical Drives
- [ ] Physical Drives Redundancy
- [ ] REVIEW TO HAVE
- [ ] Raindrop - Dropbox,Google Drive
- [ ] SECURE EMAIL SERVICES
- [ ] SECURITY INFINITY FOLDER AND EDUCATION
- [ ] SEPERATE TRAFFIC ACROSS BROWERS
- [ ] SETTINGS FOR EACH APPLICATION
- [ ] SIGN OUT OF ICLOUDAPPLE ID
- [ ] SMART BROWERS
- [ ] SSH TO VM FOR SENSITIVE ACCOUNT ACCESS (HOME SERVER X2)
- [ ] Setapp Apps
- [ ] Setup an AWS Site-to-Site Virtual Private Network (VPN)
- [ ] Setup an SMTP server (Mailgun or Sendgrid can provide a free one) for servers to send out alerts
- [ ] Software Update Tracking
- [ ] Tasks - Clickup / Zoho One - Calendars
- [ ] Time Machine save - Travel & Local
- [ ] USING YOUR YUBIKEY WITH OPENPGP
- [ ] Updates Management (Weekly)
- [ ] VPN FOR BUSINESS - CLIENT, ETC - (AWS)
- [ ] YUBI KEY FOR WORD PRESS/ WORD FENCE OR ITHEMES
- [ ] YUBIKEYS (ALL DEVICES)
- [ ] Zoho Apps
- [ ] [AWS Site To Site VPN (Part 1)](https://www.youtube.com/watch?v=5YvcyBecQts)
- [ ] [Business suite plans and pricing | Proton for Business](https://proton.me/business/plans)
- [ ] [Business suite plans and pricing | Proton for Business](https://proton.me/business/plans)
- [ ] [Payment Security | Cybersecurity and Personal Data Protection](https://privacy.com/security)
- [ ] [Setup an AWS Site-to-Site Virtual Private Network (VPN)](https://www.youtube.com/watch?v=7tTrN8WXMlg)
- [ ] [Setup an AWS Site-to-Site Virtual Private Network (VPN)](https://www.youtube.com/watch?v=7tTrN8WXMlg)
- [ ] [Yubico Support Guest Portal](https://support.yubico.com/hc/en-us/articles/360013790259-Using-Your-YubiKey-with-OpenPGP)
- [ ] [Yubico Support Guest Portal](https://support.yubico.com/hc/en-us/articles/360015668799-Smart-Card-Deployment-Manually-Importing-User-Certificates)
- [ ] [Yubico Support Guest Portal](https://support.yubico.com/hc/en-us/articles/360015668799-Smart-Card-Deployment-Manually-Importing-User-Certificates)
- [ ] https://cybersecurity.pizza/guide/free-yubikey-wordpress-authentication/#2-toc-title
- [ ] https://support.yubico.com/hc/en-us/articles/360015668919-YubiKey-Smart-Card-Deployment-Considerations

---

## Implementation Guide

**Quarterly Review Cycle:**
1. Q1: Audit Account Security + 2FA
2. Q2: Review Device & Network Security
3. Q3: Check Data & Privacy Settings
4. Q4: Physical Security & Protocols

---

*Migrated from Todoist 10 Security project - 2026-01-19*

---

## Related

- [[00-24 SO (Sovereign Operations)/04 Security (Cyber & Digital Shield) (The Shield)/04 Security|04 Security]]

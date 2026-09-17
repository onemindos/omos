---
id: so-11-air-gapped-systems
title: "Air-Gapped Systems & Networks — Build Guide"
type: reference
status: active
created: 2026-04-05
modified: 2026-04-05
owner: zeus
tags: [type/reference, status/active, domain/so, security/hardening]
ai_summary: "Comprehensive guide to building air-gapped systems and networks — hardware isolation, network segmentation, secure data transfer, key management, offline AI, and OneMind-specific implementation for the sovereign infrastructure layer."
---

# Air-Gapped Systems & Networks — Build Guide

> Build guide for physically isolated computing environments. This is how you run sovereign infrastructure that cannot be reached from the internet — period.

**Use cases for OneMind:**
- Offline AI inference (Ollama + local models) with zero cloud dependency
- Cold storage for master encryption keys, seed phrases, and recovery secrets
- Secure document vault for identity documents, legal papers, estate planning
- Air-gapped signing station (GPG, SSH key generation, code signing)
- Nuclear backup: the USB-C SSD vault that works when everything online is gone

---

## What "Air-Gapped" Actually Means

An air-gapped system has **zero network connectivity** — no WiFi, no Ethernet, no Bluetooth, no NFC, no cellular. The only way data moves in or out is through a physical medium (USB drive, QR code, printed paper) under human control.

```
┌─────────────────────────────────────────────┐
│              THE AIR GAP                     │
│                                              │
│   INTERNET ──── your network ────┤ AIR GAP │──── isolated machine
│                                  │         │
│   No wire, no wireless,         │         │
│   no signal crosses this gap    │         │
└─────────────────────────────────────────────┘
```

**The gap is physical, not logical.** Firewalls and VLANs are NOT air gaps. Software-defined isolation is NOT an air gap. The machine must be physically incapable of connecting to any network.

---

## Hardware Requirements

### Minimum Air-Gapped Workstation

| Component | Recommendation | Why |
|-----------|---------------|-----|
| **Computer** | Mac Mini M2/M4 or ThinkPad X1 Carbon | Reliable, no discrete WiFi card to worry about |
| **RAM** | 16 GB minimum, 32 GB recommended | Local LLM inference needs headroom |
| **Storage** | 1 TB NVMe SSD (internal) | OS + local models + vault data |
| **WiFi/BT** | Physically disabled or removed | Cannot be re-enabled by software |
| **Ethernet** | Cable physically disconnected | No RJ45 plugged in — ever |
| **USB ports** | 2+ USB-C or USB-A | Data transfer via sneakernet |
| **Display** | Direct-connected monitor | No wireless display protocols |
| **Input** | Wired keyboard + mouse | No Bluetooth peripherals |

### Disabling Wireless Hardware

**Mac (M-series):**
- macOS doesn't let you physically remove the WiFi chip — it's on the SoC
- Software disable: `networksetup -setairportpower en0 off` + System Settings → WiFi → Off
- **Better approach:** Use a Linux install on a Mac Mini (Ubuntu Server, no NetworkManager)
- **Best approach:** Use a machine with removable wireless card (ThinkPad) and physically remove it

**ThinkPad (recommended for true air gap):**
```bash
# Remove the WiFi M.2 card physically (inside back panel)
# Verify no wireless interfaces exist
ip link show  # should only show 'lo' (loopback)
nmcli device status  # should show no wifi/bt devices
rfkill list all  # should be empty or show hard-blocked
```

**Verify the gap:**
```bash
# These should ALL fail or return nothing
ping 8.8.8.8          # No route to host
curl https://google.com  # Connection refused
ip route show          # Empty or loopback only
ss -tuln               # No listening sockets on external interfaces
```

---

## Network Architecture

### The Three Zones

```
ZONE 1: CONNECTED (your normal network)
├── Mac (daily driver)
├── VPS (onemind-main)
├── iPhone
├── Smart home devices
└── Everything with internet access

         ══════ PHYSICAL AIR GAP ══════

ZONE 2: AIR-GAPPED WORKSTATION
├── Air-gapped machine (no network)
├── Local AI (Ollama)
├── Key signing station
├── Secure document vault
└── USB transfer port only

ZONE 3: TRANSFER MEDIA (the bridge)
├── Dedicated USB-C SSD (transfer only)
├── USB-A flash drives (write-once when possible)
├── QR codes (for small data like keys/signatures)
└── Printed paper (ultimate cold storage)
```

### Data Transfer Protocol (Sneakernet)

"Sneakernet" = physically carrying data between machines. This is the ONLY way data moves across an air gap.

**Inbound (connected → air-gapped):**
```
1. On connected machine: write data to dedicated USB drive
2. Scan USB with antivirus before transfer (ClamAV)
3. Physically carry USB to air-gapped machine
4. Mount USB read-only: mount -o ro /dev/sda1 /mnt/transfer
5. Copy needed files to air-gapped storage
6. Unmount and remove USB
7. Wipe USB before returning to connected network
```

**Outbound (air-gapped → connected):**
```
1. On air-gapped machine: write ONLY the specific output to USB
2. Physically carry USB to connected machine
3. Scan USB with antivirus before reading
4. Copy the specific file(s) you need
5. Wipe USB completely after transfer
6. Never leave stale data on transfer media
```

**Rules:**
- ONE dedicated USB drive for transfers — never use it for anything else
- Always mount read-only first to inspect contents
- Never auto-mount USB drives on either machine
- Wipe transfer media between EVERY transfer cycle
- Keep a log of every transfer (date, direction, files, purpose)

### USB Drive Preparation

```bash
# Format transfer USB (destructive — wipes everything)
# Use exFAT for Mac/Linux cross-compatibility
sudo mkfs.exfat /dev/sdX1

# Wipe after transfer (secure erase)
sudo dd if=/dev/urandom of=/dev/sdX bs=1M count=100
sudo mkfs.exfat /dev/sdX1

# Or on Mac:
diskutil eraseDisk ExFAT TRANSFER /dev/diskN
```

---

## Software Stack for Air-Gapped Machine

### Operating System

| Option | Pros | Cons |
|--------|------|------|
| **Ubuntu Server 24.04** | Minimal, no GUI bloat, great hardware support | Need terminal comfort |
| **Tails OS** (USB-bootable) | Amnesic — leaves no trace, Tor-ready if needed | Not persistent by default |
| **QubesOS** | VM-based isolation within the air gap | Complex, heavy resource use |
| **macOS** (with WiFi disabled) | Familiar, Obsidian native | Can't physically remove WiFi on M-series |

**Recommended: Ubuntu Server 24.04** — lightweight, fully controllable, no hidden network services.

### Installation Checklist

```bash
# Install Ubuntu Server (no network config during setup)
# After install, verify no network:
systemctl disable --now NetworkManager
systemctl disable --now systemd-networkd
systemctl disable --now wpa_supplicant
systemctl disable --now bluetooth

# Remove network packages entirely
apt remove --purge network-manager wpasupplicant bluetooth bluez

# Install what you need (from USB package cache)
# Pre-download .deb packages on connected machine, transfer via USB
dpkg -i /mnt/transfer/packages/*.deb
```

### Local AI (Ollama)

Run LLMs entirely offline — no API calls, no telemetry, no cloud.

```bash
# Pre-download on connected machine
curl -fsSL https://ollama.com/install.sh -o ollama-install.sh
ollama pull llama3.1:8b    # 4.7 GB
ollama pull codellama:13b  # 7.4 GB
ollama pull mistral:7b     # 4.1 GB

# Transfer the model files via USB
# Models live in ~/.ollama/models/
tar -czf ollama-models.tar.gz ~/.ollama/models/
# Copy to USB, carry to air-gapped machine

# On air-gapped machine
tar -xzf /mnt/transfer/ollama-models.tar.gz -C ~/
# Start Ollama (it runs locally, no network needed)
ollama serve &
ollama run llama3.1:8b "Summarize this document..."
```

**What you can do offline:**
- Summarize documents
- Generate code
- Analyze text
- Draft communications
- Process vault notes (RAG with ChromaDB)

### Offline RAG (Retrieval-Augmented Generation)

```bash
# Pre-install on connected machine, transfer packages via USB
pip download chromadb sentence-transformers -d /mnt/transfer/pip-packages/

# On air-gapped machine
pip install --no-index --find-links /mnt/transfer/pip-packages/ chromadb sentence-transformers

# Index your vault documents
python3 index_vault.py --vault-path ~/onemind-codex/ --db-path ~/chromadb/
# Query with local embeddings + Ollama
python3 query_vault.py "What are my security protocols?"
```

Reference: `_system/infrastructure/ops/OFFLINE-SERVICES-OVERVIEW.md` has the full Ollama + ChromaDB offline architecture.

### Obsidian (Air-Gapped Vault Access)

```bash
# Pre-download Obsidian AppImage on connected machine
# Transfer via USB to air-gapped machine

# On air-gapped machine (with display)
chmod +x Obsidian-*.AppImage
./Obsidian-*.AppImage --no-sandbox

# Open the vault from the USB SSD
# Obsidian Sync won't work (no network) — that's the point
# All changes stay local on the air-gapped machine
```

---

## Key Management (The Primary Use Case)

### What Lives on the Air-Gapped Machine

| Secret Type | Format | Storage Method |
|-------------|--------|---------------|
| GPG master key | `~/.gnupg/` | Encrypted disk partition |
| SSH CA signing key | `~/.ssh/ca_key` | Encrypted, never leaves machine |
| Bitcoin/crypto seed phrases | BIP39 mnemonic | Paper + encrypted file |
| 1Password Emergency Kit | PDF | Paper + encrypted file |
| Recovery codes (all services) | Text | Encrypted file |
| Estate planning documents | PDF | Encrypted partition |
| Code signing keys | `.p12` / `.key` | Encrypted, sign on this machine |

### GPG Key Management

```bash
# Generate master key ON the air-gapped machine
gpg --full-generate-key
# Choose: RSA and RSA, 4096 bits, never expires

# Create subkeys for daily use
gpg --edit-key YOUR_KEY_ID
> addkey  # Signing subkey
> addkey  # Encryption subkey
> addkey  # Authentication subkey
> save

# Export subkeys ONLY (master stays on air-gapped machine)
gpg --export-secret-subkeys YOUR_KEY_ID > subkeys.gpg
# Transfer subkeys via USB to your daily Mac
# Import on Mac: gpg --import subkeys.gpg

# The master key NEVER leaves the air-gapped machine
# Subkeys can sign, encrypt, authenticate — but can't create new subkeys
# If a subkey is compromised, revoke it from the air-gapped machine
```

### SSH Certificate Authority

Instead of distributing SSH keys to every server, run your own SSH CA:

```bash
# On air-gapped machine: generate CA key
ssh-keygen -t ed25519 -f ~/.ssh/ca_key -C "OneMind SSH CA"

# Sign a user's public key (brought in via USB)
ssh-keygen -s ~/.ssh/ca_key -I "zeus@onemind" -n zeus -V +52w /mnt/transfer/zeus_id_ed25519.pub
# Output: zeus_id_ed25519-cert.pub (valid for 1 year)

# Transfer the signed cert back via USB
# On servers: add to /etc/ssh/sshd_config:
#   TrustedUserCAKeys /etc/ssh/ca_key.pub
```

### Code Signing

```bash
# On air-gapped machine: sign a release
gpg --armor --detach-sign release-v1.0.tar.gz
# Output: release-v1.0.tar.gz.asc

# Transfer signature via USB to connected machine
# Publish signature alongside release
# Anyone can verify: gpg --verify release-v1.0.tar.gz.asc release-v1.0.tar.gz
```

---

## OneMind-Specific Implementation

### The USB-C SSD Vault (Phase 6)

The air-gapped machine is the natural home for the "nuclear backup" USB-C SSD vault from the sync architecture plan:

```
USB-C SSD (256GB+)
├── onemind-codex/          ← Full vault snapshot (markdown + code, ~1 GB)
├── onemind-codex/.git/     ← Full git history (273 MB)
├── _docs/ subset           ← Critical PDFs only (estate, legal, identity)
├── _images/ subset         ← Critical images only
├── recovery/
│   ├── 1password-emergency-kit.pdf
│   ├── gpg-master-key.gpg.enc
│   ├── seed-phrases.gpg.enc
│   └── recovery-codes.gpg.enc
└── tools/
    ├── ollama (binary + models)
    ├── obsidian (AppImage)
    └── chromadb (offline RAG)
```

**Workflow:**
1. Periodically snapshot vault to USB-C SSD on connected Mac
2. Carry SSD to air-gapped machine to update the offline copy
3. Air-gapped machine can browse, search, and analyze vault with Ollama
4. Any changes made on air-gapped machine are carried back via separate transfer USB
5. Air-gapped machine is the LAST RESORT — when Mac, VPS, phone, and cloud are all gone

### Domain Agent Access (Offline)

From `access-contracts.md`, the air-gapped machine is a new consumer class:

| Action | How |
|--------|-----|
| Read vault | Obsidian opens vault from USB SSD |
| Query vault | Ollama + ChromaDB local RAG |
| Sign code/keys | GPG/SSH on air-gapped machine |
| Write to vault | Changes saved to transfer USB → carried to connected Mac |

**The air-gapped machine can NEVER:**
- Connect to any network
- Run OpenClaw agents (they need NATS/network)
- Sync to Obsidian Cloud
- Push to GitHub

---

## Physical Security of the Air-Gapped Machine

| Measure | Implementation |
|---------|---------------|
| **Location** | Locked room, not visible from windows |
| **Power** | UPS battery backup (prevents sudden shutdown data loss) |
| **Screen privacy** | Privacy screen filter on monitor |
| **USB policy** | Only dedicated transfer USB touches this machine |
| **Boot security** | BIOS/firmware password set, boot from internal SSD only |
| **Disk encryption** | Full-disk encryption (LUKS on Linux, FileVault on Mac) |
| **Tamper evidence** | Tamper-evident seals on case screws (glitter nail polish works) |
| **Access log** | Physical notebook next to machine: date, who, what, why |
| **Destruction plan** | If compromised: `dd if=/dev/urandom of=/dev/sda` then physically destroy drive |

---

## Transfer Log Template

Keep this as a physical notebook next to the air-gapped machine:

```
DATE       | DIRECTION | FILES              | PURPOSE                | OPERATOR
-----------+-----------+--------------------+------------------------+---------
2026-04-05 | IN        | ollama-models.tar  | Initial AI setup       | Zeus
2026-04-05 | IN        | vault-snapshot.tar  | Vault backup           | Zeus
2026-04-10 | OUT       | signed-release.asc  | Release v1.0 signature | Zeus
2026-04-15 | IN        | vault-snapshot.tar  | Monthly vault refresh  | Zeus
```

---

## Build Phases (For OneMind)

### Phase 1 — Signing Station (Start Here)

**Cost:** ~$0 (use existing hardware)

1. Take an old laptop or Mac Mini you already own
2. Disable WiFi/Bluetooth (software disable is OK for phase 1)
3. Generate GPG master key
4. Generate SSH CA key
5. Export subkeys to daily Mac via USB
6. Store recovery secrets encrypted on this machine
7. Used for: code signing, key ceremonies, emergency recovery

### Phase 2 — Offline Vault + AI

**Cost:** ~$75-150 (USB-C SSD + dedicated transfer USB)

1. Install Ollama + models from USB transfer
2. Install ChromaDB for local RAG
3. Snapshot vault to USB-C SSD monthly
4. Air-gapped machine can now browse and query your entire vault offline
5. Used for: offline research, document analysis, vault backup verification

### Phase 3 — Hardened Workstation

**Cost:** ~$300-500 (dedicated ThinkPad with WiFi card removed)

1. Dedicated hardware — not shared with any other purpose
2. Physical WiFi/BT card removal (true air gap)
3. Full-disk encryption (LUKS)
4. Tamper-evident case seals
5. UPS power protection
6. Physical access log
7. Used for: all of the above + high-security document processing

### Phase 4 — Network Segmentation (Advanced)

If you want a middle ground between "full air gap" and "connected":

```
INTERNET ─── Router ─── VLAN 1 (trusted) ─── Mac, phone
                    ├── VLAN 2 (IoT)     ─── smart home devices
                    ├── VLAN 3 (guest)   ─── visitors
                    └── VLAN 4 (secure)  ─── sensitive workstation
                                              (network access but isolated from other VLANs)

                    ═══ AIR GAP ═══

                         VLAN 5 (air-gapped) ─── does not exist on the network at all
```

**VLAN 4 is NOT an air gap** — it's network-isolated but still connected to the internet. Use it for semi-sensitive tasks. The true air-gapped machine (Phase 3) remains fully disconnected.

---

## Threat Model

| Threat | Air Gap Protects? | Notes |
|--------|:-----------------:|-------|
| Remote hacking (internet) | ✅ Yes | No network = no remote access |
| Malware via email/web | ✅ Yes | No browser, no email client |
| Supply chain attacks | ⚠️ Partial | USB transfers could carry malware — scan everything |
| Physical theft | ❌ No | Need physical security (locks, encryption) |
| Rubber-hose cryptanalysis | ❌ No | Need deniable encryption or dead-man switches |
| Side-channel (EM emissions) | ❌ No | Requires Faraday cage (extreme threat model) |
| Evil maid (hardware tampering) | ⚠️ Partial | Tamper-evident seals help detect but don't prevent |
| USB-based attacks (BadUSB) | ⚠️ Partial | Use known-good USB drives, never accept unknown devices |

**For OneMind's threat model:** The air gap protects against the most likely attacks (remote compromise, credential theft, cloud provider breach). Physical security measures handle the rest for a personal/small-business context.

---

## Resources

- **NIST SP 800-82** — Guide to ICS Security (air-gapped industrial systems)
- **Qubes OS** — qubes-os.org (VM-based isolation OS)
- **Tails** — tails.net (amnesic live OS)
- **Ollama** — ollama.com (local LLM inference)
- **Sigstore/Cosign** — sigstore.dev (keyless code signing, for connected systems)
- **OneMind offline architecture** — `_system/infrastructure/ops/OFFLINE-SERVICES-OVERVIEW.md`

---

*Part of the OneMind Sovereign Operations security infrastructure. See [VM-SECURITY-HARDENING-COMPLETE.md](00-24%20SO%20(Sovereign%20Operations)/Security%20(Cyber%20&%20Digital%20Shield)/Hardening/VM-SECURITY-HARDENING-COMPLETE.md) for the connected-system hardening guide.*


---

## Related

- [[00-24 SO (Sovereign Operations)/04 Security (Cyber & Digital Shield) (The Shield)/04 Security|04 Security]]

---
title: "Unbound"
type: reference
status: active
tags: [platform, infrastructure, networking, dns, local]
---
# Unbound

> Local DNS resolver (8MB). Resolves `*.omos.lan` → LAN IP for Caddy reverse proxy.

## Role in OneMind

DNS backend for the NetBird nameserver group. All mesh peers query Unbound (via NetBird DNS config) to resolve `*.omos.lan` hostnames to `192.168.1.182` where Caddy routes by hostname.

## Config

- `/opt/homebrew/etc/unbound/unbound.conf`
- Listens on: `127.0.0.1:5335`, `192.168.1.182:5335`, `100.65.90.213:5335`
- Zone: `omos.lan` (static, all records → 192.168.1.182)
- Forward: everything else → 1.1.1.1, 8.8.8.8

## Integration

- **NetBird** nameserver group points to `192.168.1.182:5335` for `omos.lan` domain
- All NetBird peers automatically resolve `*.omos.lan` through this
- **Caddy** handles the actual routing once traffic arrives at the IP

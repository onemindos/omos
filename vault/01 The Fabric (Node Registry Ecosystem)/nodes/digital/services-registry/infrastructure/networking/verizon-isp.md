---
title: "Verizon"
type: reference
status: active
tags: [platform, infrastructure, networking, isp, internet]
---
# Verizon

> Internet service provider. Home connectivity backbone.

## Service
| Detail | Value |
|--------|-------|
| Type | Fios (fiber) |
| Location | Greenwich, CT |
| Use | Primary home internet for all nodes |

## Network Architecture
```
Verizon ONT → Router → Switch
                         ├── House node
                         ├── Office node
                         ├── Barn node
                         ├── IoT VLAN
                         └── Guest VLAN
```

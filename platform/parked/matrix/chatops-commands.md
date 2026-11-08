# ChatOps Command Set — Hermes + OneMind Fabric
# Access via Matrix with proper authorization

## Command Categories

### 1. Infrastructure
| Command | Auth | Description |
|---------|------|-------------|
| `!cluster status` | Admin | DOKS node health, pods, PVCs |
| `!cluster logs <pod>` | Admin | Recent logs from any pod |
| `!cluster scale <svc> <n>` | Admin | Scale deployments |
| `!cluster restart <svc>` | Admin | Restart service |
| `!cluster exec <cmd>` | Admin | Execute kubectl command |

### 2. NATS Fabric
| Command | Auth | Description |
|---------|------|-------------|
| `!nats pub <subject> <msg>` | Authorized | Publish to NATS |
| `!nats sub <subject>` | Authorized | Subscribe to stream |
| `!nats stream list` | Admin | List JetStream streams |
| `!nats agent status` | Admin | All agents on bus |
| `!nats leaf status` | Admin | Leaf node connections |

### 3. TAK / Tactical
| Command | Auth | Description |
|---------|------|-------------|
| `!tak mission list` | Ops | Active missions |
| `!tak alert <priority> <msg>` | Ops | CoT alert to all |
| `!tak position <entity>` | Ops | Get entity position |
| `!tak broadcast` | Ops | Broadcast to all EUDs |

### 4. ORP / Fusion
| Command | Auth | Description |
|---------|------|-------------|
| `!orp entity <id>` | Analyst | Retrieve fused entity |
| `!orp query <sparql>` | Analyst | Execute ORP query |
| `!orp schema list` | Analyst | Available schemas |
| `!orp ingest <source>` | Analyst | Push data source |

### 5. Legacy / Hermes
| Command | Auth | Description |
|---------|------|-------------|
| `!legacy status` | Owner | Agent health |
| `!legacy config get` | Owner | View config |
| `!legacy config set <path> <val>` | Owner | Update config |
| `!legacy memory add <fact>` | Owner | Add to memory |
| `!legacy skill list` | Owner | Loaded skills |
| `!legacy tool <name> <args>` | Owner | Execute tool |

### 6. Monitoring
| Command | Auth | Description |
|---------|------|-------------|
| `!monitor alerts` | Ops | Current alerts |
| `!monitor metrics <svc>` | Ops | Service metrics |
| `!monitor logs <level>` | Ops | Log tail |

### 7. Deployments
| Command | Auth | Description |
|---------|------|-------------|
| `!deploy <service> <env>` | Admin | Deploy to env |
| `!deploy rollback <svc>` | Admin | Rollback last |
| `!deploy history <svc>` | Admin | Deployment log |

### 8. General / Utilities
| Command | Auth | Description |
|---------|------|-------------|
| `!ping` | All | Check bot alive |
| `!help` | All | This list |
| `!echo <msg>` | All | Echo back |
| `!version` | All | Bot version |

## Authorization Matrix

| Role | Commands | Implementation |
|------|----------|----------------|
| **Owner (Zeus)** | All | Full access |
| **Admin** | Infra, NATS, Deploy | K8s admin rights |
| **Ops** | TAK, Monitoring | Read + limited ops |
| **Analyst** | ORP, Legacy | Query-only |
| **Public** | Ping, Help, Version | Basic |

## Safety Levels

- **SAFE**: Read-only (ping, help, metrics)
- **CONFIRM**: Requires approval (scale, restart)
- **DANGER**: Owner-only (exec, config set)

## Implementation

Commands processed by Hermes via Matrix bot.
Authorization via Matrix user ID check.
All commands logged to `#ops-log` room.

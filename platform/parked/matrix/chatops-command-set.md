# OneMind Matrix ChatOps Command Set
# Design Document for Infrastructure ChatOps via Matrix

## Overview

ChatOps commands executed via Matrix/Element on mobile/desktop → NATS → Fabric agents → Results back to Matrix.

## Command Architecture

```
Matrix Message (!cmd) 
    → Hermes Matrix Adapter
    → NATS Publisher (infra.cmd.exec)
    → SO Pillar Agent (subscribes)
    → Execution on Target
    → NATS Response (infra.cmd.result)
    → Hermes → Matrix Thread Reply
```

## Command Categories

### 1. Kubernetes Commands (`!k` / `!kubectl`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!k get pods -n <ns>` | List pods in namespace | Safe |
| `!k get nodes` | Cluster nodes status | Safe |
| `!k logs <pod> -n <ns> [--tail=100]` | Pod logs | Safe |
| `!k describe pod <pod> -n <ns>` | Pod details | Safe |
| `!k top nodes` | Node resource usage | Safe |
| `!k top pods -n <ns>` | Pod resource usage | Safe |
| `!k scale deployment/<name> <replicas> -n <ns>` | Scale deployment | **Dangerous** |
| `!k rollout restart deployment/<name> -n <ns>` | Restart deployment | **Dangerous** |
| `!k delete pod <pod> -n <ns>` | Delete pod (force recreate) | **Dangerous** |
| `!k exec <pod> -n <ns> -- <cmd>` | Execute in pod | **Dangerous** |

### 2. Node Management (`!node`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!node list` | List all nodes with status | Safe |
| `!node info <node>` | Detailed node info | Safe |
| `!node reboot <node>` | Reboot node | **Destructive** |
| `!node drain <node> [--force]` | Drain node for maintenance | **Dangerous** |
| `!node cordon <node>` | Mark unschedulable | **Dangerous** |
| `!node uncordon <node>` | Mark schedulable | **Dangerous** |
| `!node ssh <node> <cmd>` | SSH command on node | **Dangerous** |

### 3. TAK Commands (`!tak`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!tak nodes` | List connected TAK nodes | Safe |
| `!tak missions` | Active missions | Safe |
| `!tak mission create <name> <type>` | Create mission | **Dangerous** |
| `!tak mission end <mission-id>` | End mission | **Dangerous** |
| `!tak cot inject <json>` | Inject CoT | **Dangerous** |
| `!tak alerts` | Recent alerts | Safe |
| `!tak server status` | TAK server health | Safe |

### 4. NATS Commands (`!nats`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!nats pub <subject> <json>` | Publish to subject | **Dangerous** |
| `!nats sub <subject> [--count=10]` | Subscribe and show messages | Safe |
| `!nats kv get <bucket> <key>` | Get KV value | Safe |
| `!nats kv put <bucket> <key> <value>` | Set KV value | **Dangerous** |
| `!nats streams` | List streams | Safe |
| `!nats consumers <stream>` | List consumers | Safe |

### 5. Infrastructure Health (`!infra` / `!health`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!health` | Overall fabric health | Safe |
| `!health cluster` | K8s cluster health | Safe |
| `!health tak` | TAK server health | Safe |
| `!health nats` | NATS cluster health | Safe |
| `!health storage` | Storage/ceph health | Safe |
| `!infra services` | All service status | Safe |
| `!infra restart <service>` | Restart service | **Dangerous** |

### 6. Deployment Commands (`!deploy` / `!argo`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!argo apps` | List ArgoCD apps | Safe |
| `!argo sync <app>` | Sync ArgoCD app | **Dangerous** |
| `!argo rollback <app> <revision>` | Rollback app | **Dangerous** |
| `!helm list -n <ns>` | List Helm releases | Safe |
| `!helm upgrade <release> <chart> -n <ns>` | Upgrade release | **Dangerous** |
| `!helm rollback <release> <rev> -n <ns>` | Rollback release | **Dangerous** |

### 7. System Commands (`!sys`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!sys df` | Disk usage | Safe |
| `!sys top` | Process list | Safe |
| `!sys journal <service> [--lines=50]` | Service logs | Safe |
| `!sys restart <service>` | Restart systemd service | **Dangerous** |
| `!sys update` | Check updates | Safe |
| `!sys reboot` | Reboot machine | **Destructive** |

### 8. Legacy/Agent Commands (`!legacy`)

| Command | Description | Safety |
|---------|-------------|--------|
| `!legacy status` | Legacy agent status | Safe |
| `!legacy memory` | Memory stats | Safe |
| `!legacy task <description>` | Delegate task to Legacy | Safe |
| `!legacy pillar <so\|hp\|le\|ge> <task>` | Delegate to pillar | Safe |
| `!legacy daughter <task>` | Delegate to edge daughter | Safe |

## Safety Model

### Level 1: Safe (Read-only)
- Execute immediately
- No confirmation needed
- Examples: `!k get pods`, `!health`, `!tak nodes`

### Level 2: Dangerous (Write/State-change)
- Require 👍 reaction from authorized user
- Single approval from @zeus or @legacy
- Examples: `!k scale`, `!k rollout restart`, `!node drain`

### Level 3: Destructive (Irreversible)
- Require 2 👍 reactions from DIFFERENT authorized users
- Must include reason in thread
- Examples: `!node reboot`, `!sys reboot`, `!k delete pod`

## Authorization Matrix

| User ID | Safe | Dangerous | Destructive |
|---------|------|-----------|-------------|
| @zeus:onemindos.com | ✅ | ✅ (self-approve) | ✅ (needs 2nd) |
| @legacy:onemindos.com | ✅ | ✅ (self-approve) | ✅ (needs 2nd) |
| @so:onemindos.com | ✅ | ✅ | ❌ |
| @hp:onemindos.com | ✅ | ❌ | ❌ |
| @le:onemindos.com | ✅ | ❌ | ❌ |
| @ge:onemindos.com | ✅ | ❌ | ❌ |

## Room Configuration

### #infra:onemindos.com (Primary ChatOps)
- **Free response**: Enabled for authorized users
- **Auto-thread**: Enabled
- **Commands**: All categories
- **Notifications**: All alerts

### #alerts:onemindos.com (Alerts Only)
- **Free response**: Disabled
- **Auto-thread**: Enabled
- **Commands**: Read-only
- **Notifications**: Critical alerts with @zeus ping

### #missions:onemindos.com (TAK Missions)
- **Free response**: Authorized only
- **Auto-thread**: Enabled
- **Commands**: TAK category
- **Notifications**: Mission updates

### #legacy:onemindos.com (Legacy Direct)
- **Free response**: Enabled (DM-like)
- **Auto-thread**: Disabled
- **Commands**: All + Legacy-specific
- **Notifications**: Legacy voice transcriptions

## Implementation: NATS Command Bridge

```python
# nats_matrix_chatops.py
# Subscribes to Matrix commands via NATS, executes, returns results

import nats
import json
import subprocess
import asyncio
from dataclasses import dataclass
from typing import Optional

@dataclass
class ChatOpsCommand:
    command_id: str
    room_id: str
    thread_id: Optional[str]
    user_id: str
    command: str
    args: list
    safety_level: str  # safe, dangerous, destructive

class NATSChatOpsBridge:
    def __init__(self, nats_url: str = "nats://localhost:4222"):
        self.nats_url = nats_url
        self.nc = None
        self.command_handlers = {}
        self.authorized_users = {
            "@zeus:onemindos.com": ["safe", "dangerous", "destructive"],
            "@legacy:onemindos.com": ["safe", "dangerous", "destructive"],
            "@so:onemindos.com": ["safe", "dangerous"],
        }
        self.pending_approvals = {}  # command_id -> approval_state
    
    async def connect(self):
        self.nc = await nats.connect(self.nats_url)
        
        # Subscribe to Matrix commands
        await self.nc.subscribe("matrix.cmd.execute", cb=self.handle_command)
        
        # Subscribe to approval reactions
        await self.nc.subscribe("matrix.reaction.approval", cb=self.handle_approval)
        
        print("NATS ChatOps Bridge connected")
    
    async def handle_command(self, msg):
        data = json.loads(msg.data.decode())
        cmd = ChatOpsCommand(**data)
        
        # Check authorization
        if not self.is_authorized(cmd.user_id, cmd.safety_level):
            await self.send_result(cmd, "❌ Unauthorized for this command level")
            return
        
        # Check if approval needed
        if cmd.safety_level in ["dangerous", "destructive"]:
            await self.request_approval(cmd)
            return
        
        # Execute safe command
        await self.execute_command(cmd)
    
    def is_authorized(self, user_id: str, safety_level: str) -> bool:
        user_perms = self.authorized_users.get(user_id, ["safe"])
        return safety_level in user_perms
    
    async def request_approval(self, cmd: ChatOpsCommand):
        required_approvals = 2 if cmd.safety_level == "destructive" else 1
        
        self.pending_approvals[cmd.command_id] = {
            "command": cmd,
            "approvals": set(),
            "required": required_approvals,
            "safety_level": cmd.safety_level
        }
        
        # Publish approval request to Matrix
        await self.nc.publish("matrix.cmd.approval_request", json.dumps({
            "command_id": cmd.command_id,
            "room_id": cmd.room_id,
            "thread_id": cmd.thread_id,
            "user_id": cmd.user_id,
            "command": f"{cmd.command} {' '.join(cmd.args)}",
            "safety_level": cmd.safety_level,
            "required_approvals": required_approvals
        }).encode())
    
    async def handle_approval(self, msg):
        data = json.loads(msg.data.decode())
        command_id = data["command_id"]
        approver = data["user_id"]
        
        if command_id not in self.pending_approvals:
            return
        
        pending = self.pending_approvals[command_id]
        
        # Check approver authorization
        if not self.is_authorized(approver, pending["safety_level"]):
            return
        
        # Can't self-approve for destructive
        if pending["safety_level"] == "destructive" and approver == pending["command"].user_id:
            await self.send_result(pending["command"], "❌ Cannot self-approve destructive command")
            return
        
        pending["approvals"].add(approver)
        
        if len(pending["approvals"]) >= pending["required"]:
            # Execute command
            await self.execute_command(pending["command"])
            del self.pending_approvals[command_id]
        else:
            # Update approval count
            await self.nc.publish("matrix.cmd.approval_update", json.dumps({
                "command_id": command_id,
                "approvals": len(pending["approvals"]),
                "required": pending["required"]
            }).encode())
    
    async def execute_command(self, cmd: ChatOpsCommand):
        try:
            # Route to appropriate handler
            handler = self.command_handlers.get(cmd.command)
            if handler:
                result = await handler(cmd.args)
            else:
                # Generic shell execution (restricted)
                result = await self.execute_shell(cmd.command, cmd.args)
            
            await self.send_result(cmd, result)
        except Exception as e:
            await self.send_result(cmd, f"❌ Error: {str(e)}")
    
    async def execute_shell(self, cmd: str, args: list) -> str:
        # Restricted command allowlist
        allowed_commands = {
            "kubectl", "k", "helm", "argo", "nats", "tak", "ssh", "df", "top", "journalctl"
        }
        
        if cmd not in allowed_commands:
            return f"❌ Command '{cmd}' not in allowlist"
        
        full_cmd = [cmd] + args
        proc = await asyncio.create_subprocess_exec(
            *full_cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )
        stdout, stderr = await proc.communicate()
        
        if proc.returncode != 0:
            return f"❌ Exit {proc.returncode}: {stderr.decode()[:500]}"
        
        return stdout.decode()[:4000]  # Limit response size
    
    async def send_result(self, cmd: ChatOpsCommand, result: str):
        await self.nc.publish("matrix.cmd.result", json.dumps({
            "command_id": cmd.command_id,
            "room_id": cmd.room_id,
            "thread_id": cmd.thread_id,
            "result": result,
            "timestamp": asyncio.get_event_loop().time()
        }).encode())

# Command handlers registration
def register_handlers(bridge: NATSChatOpsBridge):
    bridge.command_handlers["k"] = handle_kubectl
    bridge.command_handlers["kubectl"] = handle_kubectl
    bridge.command_handlers["node"] = handle_node
    bridge.command_handlers["tak"] = handle_tak
    bridge.command_handlers["nats"] = handle_nats
    bridge.command_handlers["health"] = handle_health
    bridge.command_handlers["infra"] = handle_infra
    bridge.command_handlers["deploy"] = handle_deploy
    bridge.command_handlers["argo"] = handle_argo
    bridge.command_handlers["sys"] = handle_sys
    bridge.command_handlers["legacy"] = handle_legacy

async def handle_kubectl(args: list) -> str:
    # Delegate to SO pillar via NATS
    return "🔧 Delegating to SO pillar..."

async def handle_node(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_tak(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_nats(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_health(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_infra(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_deploy(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_argo(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_sys(args: list) -> str:
    return "🔧 Delegating to SO pillar..."

async def handle_legacy(args: list) -> str:
    return "🔧 Delegating to Legacy..."
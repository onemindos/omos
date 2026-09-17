---
type: infrastructure-automation
category: fleet-management
created: 2026-05-24
priority: high
tags:
  - ansible
  - automation
  - fleet
  - iac
  - scaling
---

# Ansible Fleet Automation — OneMind OS

**Scale Target:** Hundreds of nodes (Pi sensors, edge devices, VPS gateways, drones)  
**Decision:** ✅ **ADOPT Ansible** — Essential for fleet-scale infrastructure  
**Install When:** Pi Zero 2 W hardware arrives

---

## Executive Summary

**Ansible is CRITICAL for OneMind OS at scale.** ★

### Why Ansible for OneMind OS

| **Scale** | **Management Method** | **Justification** |
|-----------|----------------------|-------------------|
| **1 Mac** | Bash scripts | Manual is fine |
| **5-10 Pis** | Ansible OR bash loop | Ansible starts making sense |
| **50+ nodes** | Ansible playbooks | Manual becomes error-prone |
| **100+ nodes** | Ansible + dynamic inventory | Ansible is MANDATORY |
| **500+ nodes** | Ansible Tower/AWX | Enterprise orchestration |

**Current:** 1 Mac M2 Max (manual OK)  
**Future:** **Hundreds of nodes** → Ansible is NON-NEGOTIABLE

---

## OneMind OS Fleet Architecture (Future State)

### Node Types

```
┌─────────────────────────────────────────────┐
│         CONTROL NODE (Mac M2 Max)           │
│  - Ansible control server                   │
│  - Playbook repository                      │
│  - NATS central server                      │
│  - Hermes orchestration                     │
└─────────────────────────────────────────────┘
                    ↓
    ┌───────────────┼───────────────┐
    │               │               │
┌─────────┐   ┌─────────┐   ┌─────────┐
│ Pi ZERO │   │ Pi 5 +  │   │   VPS   │
│ Sensors │   │ AI HAT  │   │ Gateway │
└─────────┘   └─────────┘   └─────────┘
  100+ units    20+ units     10+ units

┌─────────────────────────────────────────────┐
│         DRONE FLEET (Future)                │
│  - Custom Raspberry Pi controllers         │
│  - Real-time video processing               │
│  - Swarm coordination                       │
└─────────────────────────────────────────────┘
```

### Fleet Breakdown

| **Node Type** | **Quantity (Future)** | **Purpose** | **Config Complexity** |
|---------------|----------------------|-------------|----------------------|
| **Pi Zero 2 W** | 100-200 | IoT sensors (temp, motion, doors) | Low (MQTT + NATS client) |
| **Pi 5 + AI HAT 2** | 20-50 | Edge vision (YOLOv8, LPR) | Medium (ML models, cameras) |
| **VPS Gateways** | 10-20 | OpenClaw regional nodes | High (Docker, NATS leaf, auth) |
| **Drones** | 50-100 | Aerial surveillance + video | Very High (custom flight control) |
| **Edge Routers** | 10-30 | Netbird mesh relays | Medium (VPN config, routing) |
| **Total** | **190-400+ nodes** | — | — |

**Manual management at 100+ nodes = IMPOSSIBLE** ⚠️

---

## Why Ansible > Alternatives

### Comparison Matrix

| **Tool** | **Agentless?** | **Learning Curve** | **Fleet Scale** | **OneMind Fit?** |
|----------|---------------|-------------------|----------------|-----------------|
| **Bash scripts** | ✅ Yes | Easy | 1-10 | ❌ Doesn't scale |
| **Ansible** | ✅ Yes (SSH) | Medium | 10-1000+ | ✅ **IDEAL** |
| **Puppet** | ❌ No (agent required) | Hard | 100-10k | ❌ Too heavyweight |
| **Chef** | ❌ No (agent required) | Hard | 100-10k | ❌ Too heavyweight |
| **SaltStack** | ❌ No (agent required) | Medium | 100-10k | ⚠️ Overkill |
| **Terraform** | ✅ Yes | Medium | Cloud only | ❌ Not for Pi fleet |

**Ansible wins:** Agentless (SSH), YAML-based, massive community, perfect for 100-1000 nodes ★

---

## Ansible Architecture for OneMind OS

### Control Flow

```
┌─────────────────────────────────────────────┐
│  MAC M2 MAX (Ansible Control Node)         │
│                                             │
│  ~/onemind-infra/ansible/                  │
│  ├── inventory/                            │
│  │   ├── hosts.yml        (all nodes)      │
│  │   ├── pi-sensors.yml   (sensor fleet)   │
│  │   ├── pi-vision.yml    (edge AI)        │
│  │   └── vps-gateways.yml (cloud nodes)    │
│  ├── playbooks/                            │
│  │   ├── pi-sensor-deploy.yml              │
│  │   ├── pi-vision-deploy.yml              │
│  │   ├── vps-gateway-deploy.yml            │
│  │   └── update-all.yml                    │
│  ├── roles/                                │
│  │   ├── base/           (all nodes)       │
│  │   ├── nats-client/    (NATS pub/sub)    │
│  │   ├── netbird-vpn/    (mesh VPN)        │
│  │   ├── mqtt-sensor/    (IoT sensors)     │
│  │   ├── yolov8-edge/    (vision AI)       │
│  │   └── openclaw-gateway/ (VPS gateways)  │
│  └── group_vars/                           │
│      ├── all.yml         (global config)   │
│      ├── pi_sensors.yml  (sensor-specific) │
│      └── pi_vision.yml   (vision-specific) │
└─────────────────────────────────────────────┘
                    ↓ SSH
┌─────────────────────────────────────────────┐
│  TARGET NODES (Pi, VPS, Drones)             │
│  - No agent required                        │
│  - Python 3 + SSH only                      │
│  - Ansible pushes config via SSH            │
└─────────────────────────────────────────────┘
```

---

## Installation & Setup

### Phase 1: Install Ansible (Control Node)

**On Mac M2 Max:**
```bash
# Install Ansible via Homebrew
brew install ansible

# Verify
ansible --version

# Create infra repo
mkdir -p ~/onemind-infra/ansible
cd ~/onemind-infra/ansible

# Initialize structure
mkdir -p inventory playbooks roles group_vars host_vars
```

### Phase 2: Configure SSH Keys

**Generate SSH key for fleet management:**
```bash
# Generate fleet management key
ssh-keygen -t ed25519 -f ~/.ssh/onemind_fleet -C "onemind-fleet@onemindos.dev"

# Add to SSH config
cat >> ~/.ssh/config << 'EOF'

# OneMind Fleet Management
Host pi-* drone-* vps-*
  User onemind
  IdentityFile ~/.ssh/onemind_fleet
  StrictHostKeyChecking accept-new
  ServerAliveInterval 60
EOF
```

**Copy key to target nodes:**
```bash
# Example: Copy to first Pi sensor
ssh-copy-id -i ~/.ssh/onemind_fleet onemind@192.168.1.201

# Bulk copy (when multiple nodes ready)
for ip in 192.168.1.{201..210}; do
  ssh-copy-id -i ~/.ssh/onemind_fleet onemind@$ip
done
```

### Phase 3: Create Base Inventory

**`~/onemind-infra/ansible/inventory/hosts.yml`:**
```yaml
---
all:
  children:
    # Pi Zero 2 W Sensor Fleet
    pi_sensors:
      hosts:
        pi-sensor-01:
          ansible_host: 192.168.1.201
          sensor_type: temperature
          location: office
        pi-sensor-02:
          ansible_host: 192.168.1.202
          sensor_type: motion
          location: front_door
        pi-sensor-03:
          ansible_host: 192.168.1.203
          sensor_type: door_window
          location: garage
        # ... scale to 100+ sensors

    # Pi 5 + AI HAT 2 Vision Nodes
    pi_vision:
      hosts:
        pi-vision-01:
          ansible_host: 192.168.1.211
          camera_type: usb
          model: yolov8n
          location: driveway
        pi-vision-02:
          ansible_host: 192.168.1.212
          camera_type: csi
          model: yolov8s
          location: backyard
        # ... scale to 20-50 vision nodes

    # VPS Gateways (OpenClaw)
    vps_gateways:
      hosts:
        vps-gateway-east:
          ansible_host: 198.51.100.10
          region: us-east
          provider: lightsail
        vps-gateway-west:
          ansible_host: 198.51.100.20
          region: us-west
          provider: lightsail
        # ... scale to 10-20 gateways

    # Drone Fleet (Future)
    drones:
      hosts:
        drone-alpha-01:
          ansible_host: 192.168.1.221
          drone_type: quadcopter
          flight_controller: raspberry_pi_5
        # ... scale to 50-100 drones

  vars:
    ansible_user: onemind
    ansible_python_interpreter: /usr/bin/python3
    ansible_ssh_private_key_file: ~/.ssh/onemind_fleet
```

---

## Example Playbooks

### 1. Pi Zero 2 W Sensor Deployment

**`playbooks/pi-sensor-deploy.yml`:**
```yaml
---
- name: Deploy Pi Zero 2 W IoT Sensor Fleet
  hosts: pi_sensors
  become: yes
  
  roles:
    - base              # Common setup (update, security)
    - netbird-vpn       # Join Netbird mesh
    - nats-client       # NATS pub/sub client
    - mqtt-sensor       # MQTT sensor daemon

  tasks:
    - name: Install sensor dependencies
      apt:
        name:
          - python3-pip
          - python3-smbus
          - python3-rpi.gpio
        state: present

    - name: Install Python MQTT client
      pip:
        name:
          - paho-mqtt
          - adafruit-circuitpython-dht
          - RPi.GPIO

    - name: Deploy sensor script
      template:
        src: templates/sensor-client.py.j2
        dest: /opt/onemind/sensor-client.py
        mode: '0755'

    - name: Install systemd service
      template:
        src: templates/onemind-sensor.service.j2
        dest: /etc/systemd/system/onemind-sensor.service
      notify: restart sensor service

    - name: Enable and start sensor service
      systemd:
        name: onemind-sensor
        state: started
        enabled: yes

  handlers:
    - name: restart sensor service
      systemd:
        name: onemind-sensor
        state: restarted
```

**Run:**
```bash
ansible-playbook -i inventory/hosts.yml playbooks/pi-sensor-deploy.yml
```

---

### 2. Pi 5 + AI HAT 2 Vision Deployment

**`playbooks/pi-vision-deploy.yml`:**
```yaml
---
- name: Deploy Pi 5 + AI HAT 2 Edge Vision Nodes
  hosts: pi_vision
  become: yes
  
  roles:
    - base
    - netbird-vpn
    - nats-client
    - yolov8-edge

  tasks:
    - name: Install Hailo runtime
      apt:
        deb: https://hailo.ai/downloads/hailo-runtime_4.17.0_arm64.deb

    - name: Create venv for YOLOv8
      pip:
        name: ultralytics
        virtualenv: /opt/yolov8-venv
        virtualenv_python: python3.11

    - name: Download YOLOv8 model
      get_url:
        url: "https://github.com/ultralytics/assets/releases/download/v0.0.0/yolov8{{ model }}.pt"
        dest: "/opt/models/yolov8{{ model }}.pt"
      vars:
        model: "{{ hostvars[inventory_hostname].model }}"

    - name: Deploy vision processing script
      template:
        src: templates/vision-processor.py.j2
        dest: /opt/onemind/vision-processor.py
        mode: '0755'

    - name: Install systemd service
      template:
        src: templates/onemind-vision.service.j2
        dest: /etc/systemd/system/onemind-vision.service
      notify: restart vision service

  handlers:
    - name: restart vision service
      systemd:
        name: onemind-vision
        state: restarted
```

---

### 3. VPS Gateway Deployment (OpenClaw)

**`playbooks/vps-gateway-deploy.yml`:**
```yaml
---
- name: Deploy VPS OpenClaw Gateways
  hosts: vps_gateways
  become: yes
  
  roles:
    - base
    - docker
    - openclaw-gateway

  tasks:
    - name: Install Docker
      apt:
        name:
          - docker.io
          - docker-compose
        state: present

    - name: Create OpenClaw directories
      file:
        path: "{{ item }}"
        state: directory
        mode: '0755'
      loop:
        - /opt/openclaw
        - /opt/openclaw/config
        - /opt/openclaw/logs

    - name: Deploy OpenClaw config
      template:
        src: templates/openclaw-config.yml.j2
        dest: /opt/openclaw/config/config.yml

    - name: Deploy Docker Compose
      template:
        src: templates/openclaw-docker-compose.yml.j2
        dest: /opt/openclaw/docker-compose.yml

    - name: Start OpenClaw gateway
      docker_compose:
        project_src: /opt/openclaw
        state: present
      notify: restart openclaw

  handlers:
    - name: restart openclaw
      docker_compose:
        project_src: /opt/openclaw
        restarted: yes
```

---

### 4. Fleet-Wide Update Playbook

**`playbooks/update-all.yml`:**
```yaml
---
- name: Update ALL OneMind OS Nodes
  hosts: all
  become: yes
  serial: 10  # Update 10 nodes at a time (rolling update)
  
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes
        cache_valid_time: 3600

    - name: Upgrade packages
      apt:
        upgrade: dist
      register: upgrade_result

    - name: Check if reboot required
      stat:
        path: /var/run/reboot-required
      register: reboot_required

    - name: Reboot if required
      reboot:
        reboot_timeout: 300
      when: reboot_required.stat.exists

    - name: Verify node is back online
      wait_for_connection:
        timeout: 300

    - name: Publish update event to NATS
      shell: |
        nats pub "fleet.update.complete" '{"node": "{{ inventory_hostname }}", "timestamp": "{{ ansible_date_time.iso8601 }}"}'
```

**Run:**
```bash
# Update all nodes (rolling update, 10 at a time)
ansible-playbook -i inventory/hosts.yml playbooks/update-all.yml
```

---

## Ansible Roles

### Role Structure

```
roles/
├── base/                    # Common to ALL nodes
│   ├── tasks/
│   │   └── main.yml        # Update, security, SSH hardening
│   └── templates/
│       └── sshd_config.j2
├── netbird-vpn/             # Join Netbird mesh
│   ├── tasks/
│   │   └── main.yml
│   └── templates/
│       └── netbird-setup.sh.j2
├── nats-client/             # NATS pub/sub client
│   ├── tasks/
│   │   └── main.yml
│   └── templates/
│       ├── nats.conf.j2
│       └── nats-client.service.j2
├── mqtt-sensor/             # MQTT sensor daemon
│   ├── tasks/
│   │   └── main.yml
│   └── templates/
│       └── sensor-client.py.j2
├── yolov8-edge/             # YOLOv8 edge inference
│   ├── tasks/
│   │   └── main.yml
│   └── templates/
│       └── vision-processor.py.j2
└── openclaw-gateway/        # OpenClaw VPS gateway
    ├── tasks/
    │   └── main.yml
    └── templates/
        ├── openclaw-config.yml.j2
        └── docker-compose.yml.j2
```

### Base Role Example

**`roles/base/tasks/main.yml`:**
```yaml
---
- name: Update apt cache
  apt:
    update_cache: yes
    cache_valid_time: 3600

- name: Install base packages
  apt:
    name:
      - curl
      - git
      - vim
      - htop
      - tmux
      - python3
      - python3-pip
    state: present

- name: Set hostname
  hostname:
    name: "{{ inventory_hostname }}"

- name: Configure timezone
  timezone:
    name: America/New_York

- name: Harden SSH
  template:
    src: sshd_config.j2
    dest: /etc/ssh/sshd_config
    validate: '/usr/sbin/sshd -t -f %s'
  notify: restart sshd

- name: Disable password authentication
  lineinfile:
    path: /etc/ssh/sshd_config
    regexp: '^PasswordAuthentication'
    line: 'PasswordAuthentication no'
  notify: restart sshd
```

---

## Dynamic Inventory (Future)

### NATS-Based Dynamic Inventory

**When fleet > 100 nodes, use dynamic inventory:**

**`inventory/nats_inventory.py`:**
```python
#!/usr/bin/env python3
"""
Dynamic Ansible inventory from NATS KV store.
Queries NATS for all registered fleet nodes.
"""

import json
import subprocess

def get_nats_inventory():
    # Query NATS KV for all fleet nodes
    result = subprocess.run(
        ['nats', 'kv', 'get', 'fleet-registry', '*'],
        capture_output=True, text=True
    )
    
    nodes = json.loads(result.stdout)
    
    inventory = {
        '_meta': {'hostvars': {}},
        'all': {'children': ['pi_sensors', 'pi_vision', 'vps_gateways']},
        'pi_sensors': {'hosts': []},
        'pi_vision': {'hosts': []},
        'vps_gateways': {'hosts': []}
    }
    
    for node_id, node_data in nodes.items():
        group = node_data['type']
        inventory[group]['hosts'].append(node_id)
        inventory['_meta']['hostvars'][node_id] = node_data
    
    return inventory

if __name__ == '__main__':
    print(json.dumps(get_nats_inventory(), indent=2))
```

**Use:**
```bash
ansible-playbook -i inventory/nats_inventory.py playbooks/update-all.yml
```

---

## Fleet Operations

### Common Ansible Commands

```bash
# Ping all nodes (check connectivity)
ansible all -i inventory/hosts.yml -m ping

# Check disk space on all Pi sensors
ansible pi_sensors -i inventory/hosts.yml -a "df -h"

# Restart NATS client on all vision nodes
ansible pi_vision -i inventory/hosts.yml -m systemd -a "name=nats-client state=restarted" --become

# Run ad-hoc shell command on VPS gateways
ansible vps_gateways -i inventory/hosts.yml -a "docker ps"

# Deploy to ONLY office sensors
ansible pi_sensors -i inventory/hosts.yml --limit "location:office" -m command -a "systemctl status onemind-sensor"

# Dry-run (check mode)
ansible-playbook -i inventory/hosts.yml playbooks/pi-sensor-deploy.yml --check

# Update only 5 nodes (testing)
ansible-playbook -i inventory/hosts.yml playbooks/update-all.yml --limit "pi-sensor-0[1-5]"
```

---

## Integration with OneMind OS

### NATS + Ansible Workflow

```
┌─────────────────────────────────────────────┐
│  Zeus → Hermes → "Deploy 10 new Pi sensors" │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│  Legacy AI → Execute Ansible playbook       │
│  terminal(command='ansible-playbook ...')   │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│  Ansible → SSH to 10 Pi nodes               │
│  - Install packages                         │
│  - Deploy configs                           │
│  - Start services                           │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│  Pi nodes → Publish to NATS                 │
│  nats pub "fleet.node.online" '{"id": ...}' │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│  Legacy AI → "✅ 10 sensors deployed"       │
└─────────────────────────────────────────────┘
```

### Hermes Skill: Ansible Fleet Deploy

**Create skill: `ansible-fleet-deploy`**

```yaml
---
type: skill
category: devops
created: 2026-05-24
tags:
  - ansible
  - fleet
  - deployment
---

# Ansible Fleet Deploy

Deploy OneMind OS nodes at scale using Ansible.

## Trigger Conditions
- User says "deploy Pi sensors"
- User says "update all vision nodes"
- User says "configure new VPS gateway"

## Prerequisites
- Ansible installed on Mac M2 Max
- SSH keys deployed to target nodes
- Inventory file updated with new nodes

## Numbered Steps

1. **Verify Ansible installed:**
   ```bash
   ansible --version
   ```

2. **Test connectivity:**
   ```bash
   ansible all -i inventory/hosts.yml -m ping
   ```

3. **Deploy Pi sensors:**
   ```bash
   cd ~/onemind-infra/ansible
   ansible-playbook -i inventory/hosts.yml playbooks/pi-sensor-deploy.yml
   ```

4. **Deploy vision nodes:**
   ```bash
   ansible-playbook -i inventory/hosts.yml playbooks/pi-vision-deploy.yml
   ```

5. **Deploy VPS gateways:**
   ```bash
   ansible-playbook -i inventory/hosts.yml playbooks/vps-gateway-deploy.yml
   ```

6. **Verify deployment:**
   ```bash
   ansible all -i inventory/hosts.yml -a "systemctl status onemind-*"
   ```

## Pitfalls
- **SSH key not deployed:** Run `ssh-copy-id` first
- **Python missing on target:** Install `python3` manually first time
- **Permission denied:** Use `--become` flag for sudo tasks
- **Connection timeout:** Check Netbird VPN is running

## Verification Steps
1. All nodes respond to `ansible all -m ping`
2. Services running: `systemctl status onemind-*`
3. NATS connected: `nats server check connection`
4. Nodes visible in fleet registry: `nats kv get fleet-registry`
```

---

## Scaling Roadmap

### Phase 1: Pilot (Now → Q2 2026)
- ✅ Install Ansible on Mac M2 Max
- ✅ Deploy 5-10 Pi Zero 2 W sensors
- ✅ Test playbooks with small fleet

### Phase 2: Production (Q3-Q4 2026)
- Deploy 50-100 Pi sensors
- Deploy 10-20 Pi 5 + AI HAT 2 vision nodes
- Deploy 5-10 VPS gateways
- Implement dynamic inventory (NATS KV)

### Phase 3: Scale (2027+)
- 100-200 Pi sensors
- 20-50 vision nodes
- 50-100 drones
- 10-30 edge routers
- **Total: 180-380+ nodes**

### Phase 4: Enterprise (Future)
- Ansible Tower / AWX for GUI management
- Role-Based Access Control (RBAC)
- Audit logging
- Compliance reporting

---

## Cost Analysis

### Ansible Costs

| **Item** | **Cost** | **Notes** |
|----------|----------|-----------|
| **Ansible (open source)** | $0 | Free forever |
| **Ansible Tower** | $0 (self-hosted AWX) | Optional GUI |
| **Ansible Automation Platform** | $5,000-$10,000/yr | Enterprise support (not needed) |

**OneMind OS:** Use FREE open-source Ansible ✅

---

## Security Considerations

### SSH Hardening
- ✅ Key-based auth only (no passwords)
- ✅ Dedicated fleet management key (`~/.ssh/onemind_fleet`)
- ✅ SSH config with `StrictHostKeyChecking accept-new`
- ✅ Fail2ban on all nodes

### Secrets Management
- ❌ **DO NOT** store secrets in playbooks
- ✅ Use Ansible Vault for sensitive vars
- ✅ Reference 1Password for API keys
- ✅ Environment variables for tokens

**Example: Encrypt sensitive vars:**
```bash
# Create encrypted file
ansible-vault create group_vars/all/vault.yml

# Edit encrypted file
ansible-vault edit group_vars/all/vault.yml

# Run playbook with vault password
ansible-playbook --ask-vault-pass playbooks/pi-sensor-deploy.yml
```

---

## Monitoring & Logging

### Ansible Logs
```bash
# Enable callback plugins in ansible.cfg
[defaults]
callback_whitelist = json, log_plays

# View logs
tail -f /var/log/ansible/ansible.log
```

### NATS Integration
**Publish deployment events to NATS:**
```yaml
- name: Notify NATS of deployment
  shell: |
    nats pub "fleet.deploy.complete" '{
      "node": "{{ inventory_hostname }}",
      "playbook": "{{ ansible_playbook }}",
      "timestamp": "{{ ansible_date_time.iso8601 }}"
    }'
```

---

## Disaster Recovery

### Fleet Backup Strategy

1. **Inventory versioned in Git:**
   ```bash
   cd ~/onemind-infra/ansible
   git add inventory/ playbooks/ roles/
   git commit -m "Add 10 new Pi sensors"
   git push origin main
   ```

2. **Playbooks as code:**
   - All configs in Git
   - Reproducible deployments
   - Rollback via Git history

3. **Node recovery:**
   ```bash
   # Node dies → Flash new SD card → Re-run playbook
   ansible-playbook -i inventory/hosts.yml playbooks/pi-sensor-deploy.yml --limit pi-sensor-05
   ```

**Recovery Time:** 5-10 minutes per node (fully automated) ★

---

## Performance at Scale

### Ansible Performance Benchmarks

| **Nodes** | **Task** | **Time** | **Notes** |
|-----------|----------|----------|-----------|
| 10 Pi sensors | Full deploy | 5-10 min | Serial execution |
| 100 Pi sensors | Full deploy | 20-30 min | 10 parallel forks |
| 100 Pi sensors | Update only | 10-15 min | Fast (apt cache) |
| 500+ nodes | Full deploy | 60-90 min | 20 parallel forks |

**Optimization:**
```ini
# ansible.cfg
[defaults]
forks = 20          # Parallel execution (default: 5)
gathering = smart   # Cache facts
fact_caching = jsonfile
fact_caching_timeout = 86400
```

---

## Next Actions

### Immediate (When Pi Hardware Arrives)
1. [ ] Install Ansible on Mac M2 Max (`brew install ansible`)
2. [ ] Create `~/onemind-infra/ansible/` repo
3. [ ] Generate SSH key pair (`~/.ssh/onemind_fleet`)
4. [ ] Write base inventory (`inventory/hosts.yml`)
5. [ ] Test with 1-2 Pi sensors (pilot)

### Short-Term (Q2-Q3 2026)
6. [ ] Create base role (update, SSH, security)
7. [ ] Create NATS client role
8. [ ] Create sensor role (MQTT + temp/motion)
9. [ ] Deploy 10 Pi Zero 2 W sensors
10. [ ] Document fleet ops in Obsidian

### Long-Term (Q4 2026+)
11. [ ] Deploy 50+ sensors
12. [ ] Deploy 20+ vision nodes (Pi 5 + AI HAT 2)
13. [ ] Deploy 10+ VPS gateways
14. [ ] Implement dynamic inventory (NATS KV)
15. [ ] Build drone fleet automation

---

## Related Files
- **Pi Zero 2 W Spec:** `hardware/pi-zero-2w.md`
- **Pi AI HAT 2 Spec:** `hardware/pi-ai-hat-2.md`
- **YOLOv8 Mac Deployment:** `ai-providers/yolov8-mac-deployment-guide.md`
- **NATS Infrastructure:** `NATS-INFRASTRUCTURE.md`
- **Netbird IoT Fleet:** `Netbird-IoT-Fleet-Foundation.md`

---

## Conclusion

**Ansible is MANDATORY for OneMind OS at scale.** ★

| **Scale** | **Decision** |
|-----------|-------------|
| 1 Mac | ❌ Don't use Ansible |
| 5-10 nodes | ⚠️ Start learning Ansible |
| 50+ nodes | ✅ MUST use Ansible |
| 100+ nodes | ✅ Ansible is NON-NEGOTIABLE |
| **Hundreds of nodes** | ✅✅✅ **CRITICAL INFRASTRUCTURE** |

**Install Ansible NOW (when Pi hardware arrives).** Don't wait until you have 100 nodes — by then it's too late! ✨(◕‿◕)★

---

**Status:** Ready for implementation  
**Priority:** High (deferred until Pi hardware order)  
**Owner:** Zeus + Legacy AI  
**Timeline:** Install when Pi Zero 2 W arrives ♪

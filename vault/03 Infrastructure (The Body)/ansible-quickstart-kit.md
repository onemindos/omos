---
type: ansible-quickstart
category: infrastructure-automation
created: 2026-05-24
priority: high
tags:
  - ansible
  - quickstart
  - playbooks
  - templates
---

# Ansible Quick-Start Kit — OneMind OS

**Ready-to-use playbooks, roles, and templates** ★  
**Just add Pi hardware and run!** ♪

---

## Directory Structure

```
~/onemind-infra/ansible/
├── ansible.cfg                 # Ansible config
├── inventory/
│   ├── hosts.yml              # Main inventory
│   ├── pi-sensors.yml         # Sensor fleet
│   ├── pi-vision.yml          # Vision nodes
│   └── vps-gateways.yml       # VPS gateways
├── playbooks/
│   ├── 00-bootstrap.yml       # Initial setup (run once)
│   ├── 01-pi-sensor-deploy.yml
│   ├── 02-pi-vision-deploy.yml
│   ├── 03-vps-gateway-deploy.yml
│   ├── update-all.yml         # Fleet-wide updates
│   └── health-check.yml       # Fleet health check
├── roles/
│   ├── base/                  # Common setup
│   ├── netbird-vpn/          # Netbird mesh
│   ├── nats-client/          # NATS pub/sub
│   ├── mqtt-sensor/          # MQTT sensors
│   └── yolov8-edge/          # Vision AI
├── group_vars/
│   ├── all.yml               # Global vars
│   ├── pi_sensors.yml        # Sensor-specific
│   └── pi_vision.yml         # Vision-specific
└── templates/
    ├── sensor-client.py.j2
    ├── vision-processor.py.j2
    └── nats-client.service.j2
```

---

## Quick Start (3 Steps)

### Step 1: Install Ansible

```bash
# Install on Mac M2 Max
brew install ansible

# Verify
ansible --version
```

### Step 2: Create Infrastructure Repo

```bash
# Create directory structure
mkdir -p ~/onemind-infra/ansible/{inventory,playbooks,roles,group_vars,templates}
cd ~/onemind-infra/ansible

# Copy files from this guide
# (see sections below)

# Initialize git repo
git init
git add .
git commit -m "Initial Ansible infrastructure for OneMind OS"
```

### Step 3: Run Your First Playbook

```bash
# Test connectivity
ansible all -i inventory/hosts.yml -m ping

# Bootstrap first Pi sensor
ansible-playbook -i inventory/hosts.yml playbooks/00-bootstrap.yml --limit pi-sensor-01

# Deploy sensor software
ansible-playbook -i inventory/hosts.yml playbooks/01-pi-sensor-deploy.yml --limit pi-sensor-01
```

**Done! Your first Pi sensor is deployed!** ✨★

---

## File: `ansible.cfg`

**Location:** `~/onemind-infra/ansible/ansible.cfg`

```ini
[defaults]
# Inventory
inventory = ./inventory/hosts.yml

# SSH
host_key_checking = False
remote_user = onemind
private_key_file = ~/.ssh/onemind_fleet

# Performance
forks = 10
gathering = smart
fact_caching = jsonfile
fact_caching_connection = /tmp/ansible_facts
fact_caching_timeout = 86400

# Logging
log_path = /tmp/ansible.log
callback_whitelist = profile_tasks, timer

# Privilege escalation
become = True
become_method = sudo
become_user = root
become_ask_pass = False

# Output
stdout_callback = yaml
bin_ansible_callbacks = True

[privilege_escalation]
become = True

[ssh_connection]
ssh_args = -o ControlMaster=auto -o ControlPersist=60s -o ServerAliveInterval=60
pipelining = True
```

---

## File: `inventory/hosts.yml`

**Location:** `~/onemind-infra/ansible/inventory/hosts.yml`

```yaml
---
all:
  vars:
    # Global variables
    ansible_user: onemind
    ansible_python_interpreter: /usr/bin/python3
    ansible_ssh_private_key_file: ~/.ssh/onemind_fleet
    
    # OneMind OS settings
    onemind_version: "1.0.0"
    nats_server: "192.168.1.100:4222"  # Mac M2 Max
    mqtt_server: "192.168.1.100:1883"
    netbird_setup_key: "{{ lookup('env', 'NETBIRD_SETUP_KEY') }}"
    
    # Timezone
    timezone: "America/New_York"

  children:
    # Pi Zero 2 W Sensor Fleet
    pi_sensors:
      hosts:
        pi-sensor-01:
          ansible_host: 192.168.1.201
          sensor_type: temperature_humidity
          location: office
          mqtt_topic: "onemind/sensors/office/temp"
        
        pi-sensor-02:
          ansible_host: 192.168.1.202
          sensor_type: motion
          location: front_door
          mqtt_topic: "onemind/sensors/front_door/motion"
        
        pi-sensor-03:
          ansible_host: 192.168.1.203
          sensor_type: door_window
          location: garage
          mqtt_topic: "onemind/sensors/garage/door"
        
        # Add more sensors here (scale to 100+)

    # Pi 5 + AI HAT 2 Vision Nodes (Future)
    pi_vision:
      hosts:
        pi-vision-01:
          ansible_host: 192.168.1.211
          camera_device: /dev/video0
          yolov8_model: yolov8n
          location: driveway
          nats_stream: "vision.driveway"
        
        pi-vision-02:
          ansible_host: 192.168.1.212
          camera_device: /dev/video0
          yolov8_model: yolov8s
          location: backyard
          nats_stream: "vision.backyard"

    # VPS Gateways (Future)
    vps_gateways:
      hosts:
        vps-gateway-east:
          ansible_host: 198.51.100.10
          region: us-east-1
          provider: aws_lightsail
```

---

## File: `group_vars/all.yml`

**Location:** `~/onemind-infra/ansible/group_vars/all.yml`

```yaml
---
# Global configuration for all OneMind OS nodes

# System basics
timezone: "America/New_York"
hostname_prefix: "onemind"

# SSH hardening
ssh_port: 22
ssh_password_auth: false
ssh_root_login: false
ssh_permit_empty_passwords: false

# Package lists
base_packages:
  - curl
  - git
  - vim
  - htop
  - tmux
  - python3
  - python3-pip
  - python3-venv
  - net-tools
  - dnsutils

# NATS configuration
nats_server_url: "nats://192.168.1.100:4222"
nats_credentials: "/opt/onemind/nats.creds"

# MQTT configuration
mqtt_broker: "192.168.1.100"
mqtt_port: 1883
mqtt_username: "onemind"
mqtt_password: "{{ lookup('env', 'MQTT_PASSWORD') }}"

# Netbird VPN
netbird_management_url: "https://net.onemindos.dev"
netbird_setup_key: "{{ lookup('env', 'NETBIRD_SETUP_KEY') }}"

# Monitoring
enable_monitoring: true
enable_logging: true
log_retention_days: 30
```

---

## File: `group_vars/pi_sensors.yml`

**Location:** `~/onemind-infra/ansible/group_vars/pi_sensors.yml`

```yaml
---
# Configuration specific to Pi Zero 2 W sensors

# Sensor daemon
sensor_update_interval: 60  # seconds
sensor_mqtt_qos: 1
sensor_retain_messages: true

# Python packages for sensors
sensor_packages:
  - paho-mqtt
  - adafruit-circuitpython-dht
  - RPi.GPIO
  - smbus2

# GPIO pins (customize per sensor type)
gpio_pins:
  temperature_humidity: 4
  motion: 17
  door_window: 27

# Performance tuning (Pi Zero is slow)
disable_bluetooth: true
disable_wifi_power_save: true
enable_hardware_watchdog: true
```

---

## Playbook: `00-bootstrap.yml`

**Location:** `~/onemind-infra/ansible/playbooks/00-bootstrap.yml`

```yaml
---
- name: Bootstrap OneMind OS Node (Run Once)
  hosts: all
  become: yes
  gather_facts: yes
  
  tasks:
    - name: Set hostname
      hostname:
        name: "{{ inventory_hostname }}"
    
    - name: Update /etc/hosts
      lineinfile:
        path: /etc/hosts
        line: "127.0.1.1 {{ inventory_hostname }}"
        state: present
    
    - name: Set timezone
      timezone:
        name: "{{ timezone }}"
    
    - name: Update apt cache
      apt:
        update_cache: yes
        cache_valid_time: 3600
    
    - name: Upgrade all packages
      apt:
        upgrade: dist
        autoremove: yes
        autoclean: yes
    
    - name: Install base packages
      apt:
        name: "{{ base_packages }}"
        state: present
    
    - name: Create onemind user
      user:
        name: onemind
        groups: sudo,gpio,i2c,spi
        shell: /bin/bash
        create_home: yes
        state: present
    
    - name: Add SSH authorized key for onemind
      authorized_key:
        user: onemind
        key: "{{ lookup('file', '~/.ssh/onemind_fleet.pub') }}"
        state: present
    
    - name: Disable password authentication
      lineinfile:
        path: /etc/ssh/sshd_config
        regexp: '^#?PasswordAuthentication'
        line: 'PasswordAuthentication no'
      notify: restart sshd
    
    - name: Disable root login
      lineinfile:
        path: /etc/ssh/sshd_config
        regexp: '^#?PermitRootLogin'
        line: 'PermitRootLogin no'
      notify: restart sshd
    
    - name: Create OneMind directories
      file:
        path: "{{ item }}"
        state: directory
        owner: onemind
        group: onemind
        mode: '0755'
      loop:
        - /opt/onemind
        - /opt/onemind/bin
        - /opt/onemind/config
        - /opt/onemind/logs
        - /opt/onemind/data
    
    - name: Reboot if kernel was updated
      reboot:
        reboot_timeout: 300
      when: ansible_kernel != ansible_facts.packages['linux-image-' + ansible_kernel][0].version

  handlers:
    - name: restart sshd
      systemd:
        name: sshd
        state: restarted
```

---

## Playbook: `01-pi-sensor-deploy.yml`

**Location:** `~/onemind-infra/ansible/playbooks/01-pi-sensor-deploy.yml`

```yaml
---
- name: Deploy Pi Zero 2 W Sensor Nodes
  hosts: pi_sensors
  become: yes
  gather_facts: yes
  
  tasks:
    - name: Install Python sensor packages
      pip:
        name: "{{ sensor_packages }}"
        state: present
        executable: pip3
    
    - name: Install MQTT client library
      pip:
        name: paho-mqtt
        state: present
        executable: pip3
    
    - name: Enable I2C interface
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?dtparam=i2c_arm'
        line: 'dtparam=i2c_arm=on'
      when: sensor_type == 'temperature_humidity'
      notify: reboot pi
    
    - name: Enable GPIO
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?dtoverlay=gpio'
        line: 'dtoverlay=gpio'
      notify: reboot pi
    
    - name: Deploy sensor client script
      template:
        src: ../templates/sensor-client.py.j2
        dest: /opt/onemind/bin/sensor-client.py
        owner: onemind
        group: onemind
        mode: '0755'
    
    - name: Deploy sensor systemd service
      template:
        src: ../templates/onemind-sensor.service.j2
        dest: /etc/systemd/system/onemind-sensor.service
        owner: root
        group: root
        mode: '0644'
      notify: restart sensor service
    
    - name: Enable and start sensor service
      systemd:
        name: onemind-sensor
        enabled: yes
        state: started
        daemon_reload: yes
    
    - name: Install Netbird VPN client
      shell: |
        curl -fsSL https://pkgs.netbird.io/install.sh | sh
      args:
        creates: /usr/bin/netbird
    
    - name: Connect to Netbird mesh
      shell: |
        netbird up --setup-key {{ netbird_setup_key }}
      when: netbird_setup_key is defined

  handlers:
    - name: reboot pi
      reboot:
        reboot_timeout: 300
    
    - name: restart sensor service
      systemd:
        name: onemind-sensor
        state: restarted
```

---

## Playbook: `update-all.yml`

**Location:** `~/onemind-infra/ansible/playbooks/update-all.yml`

```yaml
---
- name: Update All OneMind OS Nodes
  hosts: all
  become: yes
  serial: 10  # Update 10 nodes at a time
  
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes
        cache_valid_time: 3600
    
    - name: Check available updates
      shell: apt list --upgradable 2>/dev/null | wc -l
      register: updates_available
      changed_when: false
    
    - name: Display update count
      debug:
        msg: "{{ inventory_hostname }}: {{ updates_available.stdout }} packages to update"
    
    - name: Upgrade packages
      apt:
        upgrade: dist
        autoremove: yes
        autoclean: yes
      register: upgrade_result
    
    - name: Check if reboot required
      stat:
        path: /var/run/reboot-required
      register: reboot_required
    
    - name: Reboot if needed
      reboot:
        reboot_timeout: 300
      when: reboot_required.stat.exists
    
    - name: Wait for node to come back online
      wait_for_connection:
        timeout: 300
      when: reboot_required.stat.exists
    
    - name: Verify services are running
      systemd:
        name: "{{ item }}"
        state: started
      loop:
        - onemind-sensor
      ignore_errors: yes
    
    - name: Send completion event to NATS
      shell: |
        echo '{"node": "{{ inventory_hostname }}", "status": "updated", "timestamp": "{{ ansible_date_time.iso8601 }}"}' | \
        nats pub fleet.update.complete --server {{ nats_server_url }}
      ignore_errors: yes
```

---

## Playbook: `health-check.yml`

**Location:** `~/onemind-infra/ansible/playbooks/health-check.yml`

```yaml
---
- name: Fleet Health Check
  hosts: all
  become: yes
  gather_facts: yes
  
  tasks:
    - name: Check disk space
      shell: df -h / | awk 'NR==2 {print $5}' | sed 's/%//'
      register: disk_usage
      changed_when: false
    
    - name: Warn if disk > 80%
      debug:
        msg: "WARNING: {{ inventory_hostname }} disk usage is {{ disk_usage.stdout }}%"
      when: disk_usage.stdout|int > 80
    
    - name: Check memory usage
      shell: free | awk 'NR==2 {printf "%.0f", $3/$2 * 100}'
      register: mem_usage
      changed_when: false
    
    - name: Check CPU temperature (Pi only)
      shell: vcgencmd measure_temp | grep -oP '\d+\.\d+'
      register: cpu_temp
      changed_when: false
      when: "'pi-' in inventory_hostname"
    
    - name: Warn if CPU temp > 70°C
      debug:
        msg: "WARNING: {{ inventory_hostname }} CPU temp is {{ cpu_temp.stdout }}°C"
      when: 
        - "'pi-' in inventory_hostname"
        - cpu_temp.stdout|float > 70
    
    - name: Check if services are running
      systemd:
        name: "{{ item }}"
        state: started
      loop:
        - sshd
        - onemind-sensor
      ignore_errors: yes
      register: service_status
    
    - name: Check NATS connectivity
      shell: |
        timeout 5 nc -zv {{ nats_server_url.split('://')[1].split(':')[0] }} {{ nats_server_url.split(':')[-1] }}
      register: nats_check
      changed_when: false
      ignore_errors: yes
    
    - name: Display health summary
      debug:
        msg: |
          Node: {{ inventory_hostname }}
          Disk: {{ disk_usage.stdout }}%
          Memory: {{ mem_usage.stdout }}%
          {% if 'pi-' in inventory_hostname %}CPU Temp: {{ cpu_temp.stdout }}°C{% endif %}
          NATS: {{ 'Connected' if nats_check.rc == 0 else 'DISCONNECTED' }}
```

---

## Template: `sensor-client.py.j2`

**Location:** `~/onemind-infra/ansible/templates/sensor-client.py.j2`

```python
#!/usr/bin/env python3
"""
OneMind OS Sensor Client
Node: {{ inventory_hostname }}
Type: {{ sensor_type }}
Location: {{ location }}
"""

import time
import json
import paho.mqtt.client as mqtt

# Configuration
MQTT_BROKER = "{{ mqtt_broker }}"
MQTT_PORT = {{ mqtt_port }}
MQTT_TOPIC = "{{ mqtt_topic }}"
UPDATE_INTERVAL = {{ sensor_update_interval }}
SENSOR_TYPE = "{{ sensor_type }}"
LOCATION = "{{ location }}"

def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print(f"Connected to MQTT broker at {MQTT_BROKER}:{MQTT_PORT}")
    else:
        print(f"Connection failed with code {rc}")

def read_sensor():
    """Read sensor data based on sensor type"""
    
    {% if sensor_type == 'temperature_humidity' %}
    import adafruit_dht
    import board
    
    dht_device = adafruit_dht.DHT22(board.D{{ gpio_pins.temperature_humidity }})
    
    try:
        temperature = dht_device.temperature
        humidity = dht_device.humidity
        
        return {
            "temperature": round(temperature, 2),
            "humidity": round(humidity, 2),
            "unit": "celsius"
        }
    except RuntimeError as e:
        print(f"Sensor read error: {e}")
        return None
    
    {% elif sensor_type == 'motion' %}
    import RPi.GPIO as GPIO
    
    GPIO.setmode(GPIO.BCM)
    GPIO.setup({{ gpio_pins.motion }}, GPIO.IN)
    
    motion_detected = GPIO.input({{ gpio_pins.motion }})
    
    return {
        "motion": bool(motion_detected),
        "timestamp": time.time()
    }
    
    {% elif sensor_type == 'door_window' %}
    import RPi.GPIO as GPIO
    
    GPIO.setmode(GPIO.BCM)
    GPIO.setup({{ gpio_pins.door_window }}, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    
    door_open = not GPIO.input({{ gpio_pins.door_window }})
    
    return {
        "door_open": bool(door_open),
        "timestamp": time.time()
    }
    
    {% else %}
    # Default: return dummy data
    return {
        "status": "unknown_sensor_type",
        "type": "{{ sensor_type }}"
    }
    {% endif %}

def main():
    # Create MQTT client
    client = mqtt.Client(client_id="{{ inventory_hostname }}")
    client.on_connect = on_connect
    
    # Connect to broker
    client.connect(MQTT_BROKER, MQTT_PORT, 60)
    client.loop_start()
    
    print(f"OneMind Sensor Client started: {SENSOR_TYPE} at {LOCATION}")
    
    try:
        while True:
            # Read sensor
            data = read_sensor()
            
            if data:
                # Add metadata
                payload = {
                    "node": "{{ inventory_hostname }}",
                    "sensor_type": SENSOR_TYPE,
                    "location": LOCATION,
                    "timestamp": time.time(),
                    "data": data
                }
                
                # Publish to MQTT
                client.publish(MQTT_TOPIC, json.dumps(payload), qos={{ sensor_mqtt_qos }}, retain={{ sensor_retain_messages|lower }})
                print(f"Published: {payload}")
            
            time.sleep(UPDATE_INTERVAL)
    
    except KeyboardInterrupt:
        print("\nShutting down...")
        client.loop_stop()
        client.disconnect()

if __name__ == "__main__":
    main()
```

---

## Template: `onemind-sensor.service.j2`

**Location:** `~/onemind-infra/ansible/templates/onemind-sensor.service.j2`

```ini
[Unit]
Description=OneMind OS Sensor Client ({{ sensor_type }})
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=onemind
Group=onemind
WorkingDirectory=/opt/onemind
ExecStart=/usr/bin/python3 /opt/onemind/bin/sensor-client.py
Restart=always
RestartSec=10
StandardOutput=append:/opt/onemind/logs/sensor.log
StandardError=append:/opt/onemind/logs/sensor-error.log

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/onemind/logs /opt/onemind/data

[Install]
WantedBy=multi-user.target
```

---

## Usage Examples

### Deploy First Sensor (Pilot)

```bash
cd ~/onemind-infra/ansible

# Test connectivity
ansible pi-sensor-01 -i inventory/hosts.yml -m ping

# Bootstrap node
ansible-playbook playbooks/00-bootstrap.yml --limit pi-sensor-01

# Deploy sensor software
ansible-playbook playbooks/01-pi-sensor-deploy.yml --limit pi-sensor-01

# Check status
ansible pi-sensor-01 -i inventory/hosts.yml -a "systemctl status onemind-sensor" --become
```

### Deploy 10 Sensors at Once

```bash
# Deploy to all sensors
ansible-playbook playbooks/01-pi-sensor-deploy.yml

# Or deploy to specific group
ansible-playbook playbooks/01-pi-sensor-deploy.yml --limit "pi-sensor-0[1-9],pi-sensor-10"
```

### Update All Nodes

```bash
# Dry run first
ansible-playbook playbooks/update-all.yml --check

# Execute update (10 nodes at a time)
ansible-playbook playbooks/update-all.yml
```

### Health Check

```bash
# Check all nodes
ansible-playbook playbooks/health-check.yml

# Check only sensors
ansible-playbook playbooks/health-check.yml --limit pi_sensors
```

### Ad-Hoc Commands

```bash
# Check disk space on all nodes
ansible all -i inventory/hosts.yml -a "df -h /" --become

# Restart sensor service on all Pi sensors
ansible pi_sensors -i inventory/hosts.yml -m systemd -a "name=onemind-sensor state=restarted" --become

# Check CPU temperature on all Pis
ansible pi_sensors -i inventory/hosts.yml -a "vcgencmd measure_temp"

# Run custom command
ansible pi-sensor-01 -i inventory/hosts.yml -a "uptime"
```

---

## Pre-Flight Checklist

Before running playbooks:

- [ ] Ansible installed (`brew install ansible`)
- [ ] SSH key generated (`ssh-keygen -t ed25519 -f ~/.ssh/onemind_fleet`)
- [ ] SSH key copied to Pi (`ssh-copy-id onemind@192.168.1.201`)
- [ ] Pi has Python 3 installed (`ssh onemind@192.168.1.201 'python3 --version'`)
- [ ] Inventory file updated with real IP addresses
- [ ] Environment variables set:
  ```bash
  export NETBIRD_SETUP_KEY="your-netbird-key"
  export MQTT_PASSWORD="your-mqtt-password"
  ```
- [ ] NATS server running on Mac M2 Max (192.168.1.100:4222)
- [ ] Test connectivity: `ansible all -m ping`

---

## Troubleshooting

### SSH Connection Fails

```bash
# Check SSH key
ssh -i ~/.ssh/onemind_fleet onemind@192.168.1.201

# Copy key manually
ssh-copy-id -i ~/.ssh/onemind_fleet onemind@192.168.1.201
```

### Python Not Found

```bash
# Install Python on Pi manually first time
ssh onemind@192.168.1.201
sudo apt update && sudo apt install -y python3 python3-pip
```

### Permission Denied

```bash
# Ensure onemind user has sudo rights
ssh onemind@192.168.1.201
sudo visudo
# Add: onemind ALL=(ALL) NOPASSWD: ALL
```

### Service Won't Start

```bash
# Check logs
ansible pi-sensor-01 -a "journalctl -u onemind-sensor -n 50" --become

# Test script manually
ansible pi-sensor-01 -a "python3 /opt/onemind/bin/sensor-client.py" --become
```

---

## Next Steps

1. **Install Ansible:** `brew install ansible`
2. **Generate SSH keys:** `ssh-keygen -t ed25519 -f ~/.ssh/onemind_fleet`
3. **Create repo:** Copy all files from this guide to `~/onemind-infra/ansible/`
4. **When Pi arrives:** Update `inventory/hosts.yml` with real IP
5. **Bootstrap:** `ansible-playbook playbooks/00-bootstrap.yml`
6. **Deploy:** `ansible-playbook playbooks/01-pi-sensor-deploy.yml`
7. **Scale:** Add more Pis to inventory and re-run playbooks

---

**Ready to deploy 100+ nodes with a single command!** ✨(◕‿◕)★

**Status:** Copy-paste ready  
**Test:** When Pi Zero 2 W hardware arrives  
**Scale:** 10 → 100 → 1000+ nodes ♪

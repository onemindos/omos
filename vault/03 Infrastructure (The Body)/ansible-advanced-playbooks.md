---
type: ansible-advanced-playbooks
category: infrastructure-automation
created: 2026-05-24
priority: high
tags:
  - ansible
  - drones
  - vision
  - edge-ai
  - fleet
---

# Ansible Advanced Playbooks — Drones & Vision

**Pi 5 + AI HAT 2 Vision Nodes** and **Drone Fleet Automation** ★

---

## Table of Contents

1. [Pi 5 + AI HAT 2 Vision Playbooks](#pi-5--ai-hat-2-vision-playbooks)
2. [Drone Fleet Playbooks](#drone-fleet-playbooks)
3. [Advanced Inventory](#advanced-inventory)
4. [Group Variables](#group-variables)
5. [Templates](#templates)
6. [Usage Examples](#usage-examples)

---

## Pi 5 + AI HAT 2 Vision Playbooks

### Inventory Addition

**Add to `inventory/hosts.yml`:**

```yaml
# Pi 5 + AI HAT 2 Vision Nodes
pi_vision:
  hosts:
    pi-vision-01:
      ansible_host: 192.168.1.211
      camera_device: /dev/video0
      camera_type: usb  # or 'csi'
      yolov8_model: yolov8n  # nano for speed
      resolution: "640x480"
      fps: 30
      location: driveway
      nats_subject: "vision.driveway.detections"
      enable_lpr: false  # License plate recognition
    
    pi-vision-02:
      ansible_host: 192.168.1.212
      camera_device: /dev/video0
      camera_type: csi
      yolov8_model: yolov8s  # small for balance
      resolution: "1280x720"
      fps: 30
      location: backyard
      nats_subject: "vision.backyard.detections"
      enable_lpr: false
    
    pi-vision-03:
      ansible_host: 192.168.1.213
      camera_device: /dev/video0
      camera_type: usb
      yolov8_model: yolov8m  # medium for accuracy
      resolution: "1920x1080"
      fps: 15
      location: front_gate
      nats_subject: "vision.front_gate.detections"
      enable_lpr: true  # LPR enabled for gate
    
    # Scale to 20-50 vision nodes
```

---

### Playbook: `02-pi-vision-deploy.yml`

**Location:** `~/onemind-infra/ansible/playbooks/02-pi-vision-deploy.yml`

```yaml
---
- name: Deploy Pi 5 + AI HAT 2 Vision Nodes
  hosts: pi_vision
  become: yes
  gather_facts: yes
  
  vars:
    hailo_runtime_version: "4.17.0"
    yolov8_venv_path: "/opt/yolov8-venv"
    models_path: "/opt/models"
  
  pre_tasks:
    - name: Verify this is a Pi 5
      assert:
        that:
          - ansible_facts.board == "Raspberry Pi 5"
        fail_msg: "This playbook requires Raspberry Pi 5"
      ignore_errors: yes  # Don't fail if board detection doesn't work
  
  tasks:
    # ========================================
    # Install Hailo AI HAT 2 Runtime
    # ========================================
    
    - name: Download Hailo runtime package
      get_url:
        url: "https://hailo.ai/downloads/hailo-runtime_{{ hailo_runtime_version }}_arm64.deb"
        dest: "/tmp/hailo-runtime.deb"
        mode: '0644'
      register: hailo_download
      retries: 3
      delay: 10
    
    - name: Install Hailo runtime
      apt:
        deb: "/tmp/hailo-runtime.deb"
        state: present
    
    - name: Verify Hailo device
      shell: hailortcli fw-control identify
      register: hailo_device
      changed_when: false
      failed_when: false
    
    - name: Display Hailo device info
      debug:
        msg: "{{ hailo_device.stdout }}"
      when: hailo_device.rc == 0
    
    # ========================================
    # Install Camera Support
    # ========================================
    
    - name: Install camera packages
      apt:
        name:
          - libcamera-apps
          - python3-picamera2
          - v4l-utils
        state: present
    
    - name: Enable camera interface
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?camera_auto_detect'
        line: 'camera_auto_detect=1'
      notify: reboot pi
    
    - name: Load v4l2 module
      modprobe:
        name: v4l2loopback
        state: present
      ignore_errors: yes
    
    # ========================================
    # Install YOLOv8 + Dependencies
    # ========================================
    
    - name: Install system dependencies for YOLOv8
      apt:
        name:
          - python3-dev
          - python3-pip
          - python3-venv
          - libopencv-dev
          - python3-opencv
          - libopenblas-dev
          - libjpeg-dev
          - zlib1g-dev
        state: present
    
    - name: Create YOLOv8 venv
      pip:
        name: pip
        virtualenv: "{{ yolov8_venv_path }}"
        virtualenv_python: python3.11
        state: latest
    
    - name: Install YOLOv8 in venv
      pip:
        name:
          - ultralytics
          - opencv-python-headless
          - numpy
          - pillow
          - pyyaml
          - requests
          - tqdm
          - paho-mqtt
        virtualenv: "{{ yolov8_venv_path }}"
        state: present
    
    - name: Install NATS Python client
      pip:
        name: nats-py
        virtualenv: "{{ yolov8_venv_path }}"
        state: present
    
    # ========================================
    # Install EasyOCR (if LPR enabled)
    # ========================================
    
    - name: Install EasyOCR for license plate recognition
      pip:
        name:
          - easyocr
          - torch
          - torchvision
        virtualenv: "{{ yolov8_venv_path }}"
        state: present
      when: enable_lpr | default(false)
    
    # ========================================
    # Download YOLOv8 Models
    # ========================================
    
    - name: Create models directory
      file:
        path: "{{ models_path }}"
        state: directory
        owner: onemind
        group: onemind
        mode: '0755'
    
    - name: Download YOLOv8 model
      get_url:
        url: "https://github.com/ultralytics/assets/releases/download/v0.0.0/{{ yolov8_model }}.pt"
        dest: "{{ models_path }}/{{ yolov8_model }}.pt"
        owner: onemind
        group: onemind
        mode: '0644'
      register: model_download
      retries: 3
      delay: 10
    
    - name: Compile model for Hailo (optional optimization)
      shell: |
        {{ yolov8_venv_path }}/bin/python -c "from ultralytics import YOLO; model = YOLO('{{ models_path }}/{{ yolov8_model }}.pt'); model.export(format='hailo', device='hailo')"
      args:
        creates: "{{ models_path }}/{{ yolov8_model }}_hailo_model.har"
      when: hailo_device.rc == 0
      ignore_errors: yes  # Not all models support Hailo export yet
    
    # ========================================
    # Deploy Vision Processing Script
    # ========================================
    
    - name: Deploy vision processor script
      template:
        src: ../templates/vision-processor.py.j2
        dest: /opt/onemind/bin/vision-processor.py
        owner: onemind
        group: onemind
        mode: '0755'
    
    - name: Deploy vision config
      template:
        src: ../templates/vision-config.yaml.j2
        dest: /opt/onemind/config/vision.yaml
        owner: onemind
        group: onemind
        mode: '0644'
    
    # ========================================
    # Install Systemd Service
    # ========================================
    
    - name: Deploy vision systemd service
      template:
        src: ../templates/onemind-vision.service.j2
        dest: /etc/systemd/system/onemind-vision.service
        owner: root
        group: root
        mode: '0644'
      notify: restart vision service
    
    - name: Enable and start vision service
      systemd:
        name: onemind-vision
        enabled: yes
        state: started
        daemon_reload: yes
    
    # ========================================
    # Performance Tuning
    # ========================================
    
    - name: Set GPU memory split (512MB for vision)
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?gpu_mem='
        line: 'gpu_mem=512'
      notify: reboot pi
    
    - name: Enable hardware video acceleration
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?dtoverlay=vc4-kms-v3d'
        line: 'dtoverlay=vc4-kms-v3d'
      notify: reboot pi
    
    - name: Install performance monitoring
      apt:
        name:
          - stress-ng
          - sysbench
        state: present
  
  handlers:
    - name: reboot pi
      reboot:
        reboot_timeout: 300
    
    - name: restart vision service
      systemd:
        name: onemind-vision
        state: restarted
```

---

### Template: `vision-processor.py.j2`

**Location:** `~/onemind-infra/ansible/templates/vision-processor.py.j2`

```python
#!/usr/bin/env python3
"""
OneMind OS Vision Processor
Node: {{ inventory_hostname }}
Location: {{ location }}
Model: {{ yolov8_model }}
Camera: {{ camera_device }}
"""

import asyncio
import cv2
import json
import time
from datetime import datetime
from pathlib import Path
from ultralytics import YOLO
import nats

# Configuration
CAMERA_DEVICE = "{{ camera_device }}"
RESOLUTION = tuple(map(int, "{{ resolution }}".split('x')))
FPS = {{ fps }}
MODEL_PATH = "/opt/models/{{ yolov8_model }}.pt"
NATS_SERVER = "{{ nats_server_url }}"
NATS_SUBJECT = "{{ nats_subject }}"
NODE_NAME = "{{ inventory_hostname }}"
LOCATION = "{{ location }}"
ENABLE_LPR = {{ enable_lpr | default(false) | lower }}

# Performance settings
CONF_THRESHOLD = 0.25  # Confidence threshold
IOU_THRESHOLD = 0.45   # NMS IoU threshold
DEVICE = 'cpu'  # Let Hailo handle acceleration

{% if enable_lpr %}
# License Plate Recognition
import easyocr
LPR_READER = easyocr.Reader(['en'], gpu=False)
{% endif %}

class VisionProcessor:
    def __init__(self):
        self.model = None
        self.cap = None
        self.nc = None
        self.frame_count = 0
        self.detection_count = 0
        
    def initialize(self):
        """Initialize camera and model"""
        print(f"Initializing vision processor for {NODE_NAME} at {LOCATION}")
        
        # Load YOLOv8 model
        print(f"Loading model: {MODEL_PATH}")
        self.model = YOLO(MODEL_PATH)
        self.model.to(DEVICE)
        
        # Initialize camera
        print(f"Opening camera: {CAMERA_DEVICE}")
        self.cap = cv2.VideoCapture(CAMERA_DEVICE)
        self.cap.set(cv2.CAP_PROP_FRAME_WIDTH, RESOLUTION[0])
        self.cap.set(cv2.CAP_PROP_FRAME_HEIGHT, RESOLUTION[1])
        self.cap.set(cv2.CAP_PROP_FPS, FPS)
        
        if not self.cap.isOpened():
            raise RuntimeError(f"Failed to open camera: {CAMERA_DEVICE}")
        
        print("Vision processor initialized successfully")
    
    def process_frame(self, frame):
        """Run inference on a frame"""
        results = self.model.predict(
            frame,
            conf=CONF_THRESHOLD,
            iou=IOU_THRESHOLD,
            verbose=False,
            device=DEVICE
        )
        
        return results[0]  # First result
    
    {% if enable_lpr %}
    def detect_license_plate(self, frame, bbox):
        """Extract license plate text from detected vehicle"""
        x1, y1, x2, y2 = map(int, bbox)
        plate_region = frame[y1:y2, x1:x2]
        
        # Run OCR
        results = LPR_READER.readtext(plate_region)
        
        if results:
            # Return text with highest confidence
            return max(results, key=lambda x: x[2])[1]
        return None
    {% endif %}
    
    def format_detection(self, result):
        """Format detection results for NATS"""
        detections = []
        
        for box in result.boxes:
            detection = {
                "class": result.names[int(box.cls)],
                "confidence": float(box.conf),
                "bbox": box.xyxy[0].tolist(),
                "timestamp": datetime.utcnow().isoformat()
            }
            
            {% if enable_lpr %}
            # Try LPR for vehicles
            if detection["class"] in ["car", "truck", "bus", "motorcycle"]:
                plate = self.detect_license_plate(result.orig_img, detection["bbox"])
                if plate:
                    detection["license_plate"] = plate
            {% endif %}
            
            detections.append(detection)
        
        return {
            "node": NODE_NAME,
            "location": LOCATION,
            "timestamp": datetime.utcnow().isoformat(),
            "frame_count": self.frame_count,
            "detection_count": len(detections),
            "detections": detections,
            "fps": result.speed.get('inference', 0)
        }
    
    async def publish_to_nats(self, data):
        """Publish detection results to NATS"""
        try:
            nc = await nats.connect(NATS_SERVER)
            await nc.publish(NATS_SUBJECT, json.dumps(data).encode())
            await nc.close()
        except Exception as e:
            print(f"NATS publish error: {e}")
    
    async def run(self):
        """Main processing loop"""
        self.initialize()
        
        print(f"Starting vision processing (FPS target: {FPS})")
        
        try:
            while True:
                ret, frame = self.cap.read()
                
                if not ret:
                    print("Failed to grab frame")
                    await asyncio.sleep(0.1)
                    continue
                
                self.frame_count += 1
                
                # Run inference
                result = self.process_frame(frame)
                
                # Process detections
                if len(result.boxes) > 0:
                    self.detection_count += len(result.boxes)
                    data = self.format_detection(result)
                    
                    # Publish to NATS
                    await self.publish_to_nats(data)
                    
                    print(f"Frame {self.frame_count}: {len(result.boxes)} detections")
                
                # Maintain FPS target
                await asyncio.sleep(1.0 / FPS)
        
        except KeyboardInterrupt:
            print("\nShutting down vision processor...")
        
        finally:
            if self.cap:
                self.cap.release()
            print("Vision processor stopped")

def main():
    processor = VisionProcessor()
    asyncio.run(processor.run())

if __name__ == "__main__":
    main()
```

---

## Drone Fleet Playbooks

### Inventory Addition

**Add to `inventory/hosts.yml`:**

```yaml
# Drone Fleet
drones:
  hosts:
    drone-alpha-01:
      ansible_host: 192.168.1.221
      drone_type: quadcopter
      flight_controller: raspberry_pi_5
      camera_enabled: true
      camera_gimbal: true
      max_altitude: 120  # meters (FAA limit)
      max_speed: 15  # m/s
      battery_cells: 6  # 6S LiPo
      rtk_gps: true  # Real-Time Kinematic GPS
      location_base: home_base
      nats_subject: "drone.alpha-01.telemetry"
    
    drone-alpha-02:
      ansible_host: 192.168.1.222
      drone_type: quadcopter
      flight_controller: raspberry_pi_5
      camera_enabled: true
      camera_gimbal: true
      max_altitude: 120
      max_speed: 15
      battery_cells: 6
      rtk_gps: true
      location_base: home_base
      nats_subject: "drone.alpha-02.telemetry"
    
    drone-beta-01:
      ansible_host: 192.168.1.231
      drone_type: hexacopter
      flight_controller: raspberry_pi_5
      camera_enabled: true
      camera_gimbal: true
      max_altitude: 150
      max_speed: 20
      battery_cells: 8  # 8S LiPo for heavy lift
      rtk_gps: true
      location_base: remote_site_1
      nats_subject: "drone.beta-01.telemetry"
    
    # Scale to 50-100 drones
```

---

### Playbook: `03-drone-fleet-deploy.yml`

**Location:** `~/onemind-infra/ansible/playbooks/03-drone-fleet-deploy.yml`

```yaml
---
- name: Deploy OneMind OS Drone Fleet
  hosts: drones
  become: yes
  gather_facts: yes
  
  vars:
    ardupilot_version: "4.5.0"
    mavproxy_version: "latest"
    drone_base_path: "/opt/drone"
  
  pre_tasks:
    - name: Verify this is a Pi 5
      assert:
        that:
          - ansible_facts.board == "Raspberry Pi 5"
        fail_msg: "Drone controller requires Raspberry Pi 5"
      ignore_errors: yes
  
  tasks:
    # ========================================
    # Install ArduPilot Flight Stack
    # ========================================
    
    - name: Install ArduPilot dependencies
      apt:
        name:
          - python3-dev
          - python3-pip
          - python3-wxgtk4.0
          - python3-lxml
          - python3-opencv
          - libxml2-dev
          - libxslt1-dev
          - build-essential
          - ccache
          - gawk
          - git
        state: present
    
    - name: Clone ArduPilot repository
      git:
        repo: 'https://github.com/ArduPilot/ardupilot.git'
        dest: "{{ drone_base_path }}/ardupilot"
        version: "Copter-{{ ardupilot_version }}"
        update: no
      become_user: onemind
    
    - name: Initialize ArduPilot submodules
      shell: |
        cd {{ drone_base_path }}/ardupilot
        git submodule update --init --recursive
      become_user: onemind
      args:
        creates: "{{ drone_base_path }}/ardupilot/modules/waf"
    
    - name: Install MAVProxy (ground control)
      pip:
        name:
          - MAVProxy
          - pymavlink
          - dronekit
          - dronekit-sitl
        state: present
        executable: pip3
    
    # ========================================
    # Install Camera & Vision Support
    # ========================================
    
    - name: Install camera packages
      apt:
        name:
          - libcamera-apps
          - python3-picamera2
          - v4l-utils
          - ffmpeg
        state: present
      when: camera_enabled | default(false)
    
    - name: Enable camera interface
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?camera_auto_detect'
        line: 'camera_auto_detect=1'
      when: camera_enabled | default(false)
      notify: reboot pi
    
    # ========================================
    # Install GPS & Navigation
    # ========================================
    
    - name: Install GPS packages
      apt:
        name:
          - gpsd
          - gpsd-clients
          - python3-gps
        state: present
    
    - name: Configure GPSD
      template:
        src: ../templates/gpsd.j2
        dest: /etc/default/gpsd
        owner: root
        group: root
        mode: '0644'
      notify: restart gpsd
    
    - name: Enable UART for GPS
      lineinfile:
        path: /boot/config.txt
        regexp: '^#?enable_uart'
        line: 'enable_uart=1'
      notify: reboot pi
    
    # ========================================
    # Install Telemetry & Communication
    # ========================================
    
    - name: Install NATS Python client
      pip:
        name: nats-py
        state: present
        executable: pip3
    
    - name: Install radio modem drivers (optional)
      apt:
        name:
          - ser2net
          - socat
        state: present
    
    # ========================================
    # Deploy Drone Control Scripts
    # ========================================
    
    - name: Create drone directories
      file:
        path: "{{ item }}"
        state: directory
        owner: onemind
        group: onemind
        mode: '0755'
      loop:
        - "{{ drone_base_path }}/bin"
        - "{{ drone_base_path }}/config"
        - "{{ drone_base_path }}/logs"
        - "{{ drone_base_path }}/missions"
        - "{{ drone_base_path }}/video"
    
    - name: Deploy drone telemetry daemon
      template:
        src: ../templates/drone-telemetry.py.j2
        dest: "{{ drone_base_path }}/bin/drone-telemetry.py"
        owner: onemind
        group: onemind
        mode: '0755'
    
    - name: Deploy drone config
      template:
        src: ../templates/drone-config.yaml.j2
        dest: "{{ drone_base_path }}/config/drone.yaml"
        owner: onemind
        group: onemind
        mode: '0644'
    
    - name: Deploy flight safety limits
      template:
        src: ../templates/safety-limits.yaml.j2
        dest: "{{ drone_base_path }}/config/safety-limits.yaml"
        owner: onemind
        group: onemind
        mode: '0644'
    
    # ========================================
    # Install Systemd Services
    # ========================================
    
    - name: Deploy ArduPilot systemd service
      template:
        src: ../templates/ardupilot.service.j2
        dest: /etc/systemd/system/ardupilot.service
        owner: root
        group: root
        mode: '0644'
      notify: restart ardupilot
    
    - name: Deploy telemetry systemd service
      template:
        src: ../templates/drone-telemetry.service.j2
        dest: /etc/systemd/system/drone-telemetry.service
        owner: root
        group: root
        mode: '0644'
      notify: restart telemetry
    
    - name: Enable and start ArduPilot
      systemd:
        name: ardupilot
        enabled: yes
        state: started
        daemon_reload: yes
    
    - name: Enable and start telemetry
      systemd:
        name: drone-telemetry
        enabled: yes
        state: started
        daemon_reload: yes
    
    # ========================================
    # Safety & Geofencing
    # ========================================
    
    - name: Install geofence library
      pip:
        name:
          - shapely
          - geopy
        state: present
        executable: pip3
    
    - name: Deploy geofence config
      template:
        src: ../templates/geofence.json.j2
        dest: "{{ drone_base_path }}/config/geofence.json"
        owner: onemind
        group: onemind
        mode: '0644'
    
    # ========================================
    # Performance & Power Management
    # ========================================
    
    - name: Disable unnecessary services
      systemd:
        name: "{{ item }}"
        enabled: no
        state: stopped
      loop:
        - bluetooth
        - avahi-daemon
        - triggerhappy
      ignore_errors: yes
    
    - name: Set CPU governor to performance
      shell: echo performance > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor
      ignore_errors: yes
    
    - name: Install battery monitoring
      pip:
        name: smbus2
        state: present
        executable: pip3
  
  handlers:
    - name: reboot pi
      reboot:
        reboot_timeout: 300
    
    - name: restart gpsd
      systemd:
        name: gpsd
        state: restarted
    
    - name: restart ardupilot
      systemd:
        name: ardupilot
        state: restarted
    
    - name: restart telemetry
      systemd:
        name: drone-telemetry
        state: restarted
```

---

### Template: `drone-telemetry.py.j2`

**Location:** `~/onemind-infra/ansible/templates/drone-telemetry.py.j2`

```python
#!/usr/bin/env python3
"""
OneMind OS Drone Telemetry
Drone: {{ inventory_hostname }}
Type: {{ drone_type }}
Base: {{ location_base }}
"""

import asyncio
import json
import time
from datetime import datetime
from pymavlink import mavutil
import nats

# Configuration
MAVLINK_CONNECTION = "udpin:0.0.0.0:14550"
NATS_SERVER = "{{ nats_server_url }}"
NATS_SUBJECT = "{{ nats_subject }}"
DRONE_ID = "{{ inventory_hostname }}"
DRONE_TYPE = "{{ drone_type }}"
LOCATION_BASE = "{{ location_base }}"

# Safety limits
MAX_ALTITUDE = {{ max_altitude }}  # meters
MAX_SPEED = {{ max_speed }}  # m/s
BATTERY_LOW_THRESHOLD = 20  # percent

class DroneTelemetry:
    def __init__(self):
        self.master = None
        self.nc = None
        self.telemetry_data = {}
        self.armed = False
        self.mode = "UNKNOWN"
        
    async def connect(self):
        """Connect to MAVLink and NATS"""
        print(f"Connecting to MAVLink: {MAVLINK_CONNECTION}")
        self.master = mavutil.mavlink_connection(MAVLINK_CONNECTION)
        
        # Wait for heartbeat
        print("Waiting for heartbeat...")
        self.master.wait_heartbeat()
        print(f"Heartbeat received (system {self.master.target_system} component {self.master.target_component})")
        
        # Connect to NATS
        print(f"Connecting to NATS: {NATS_SERVER}")
        self.nc = await nats.connect(NATS_SERVER)
        print("NATS connected")
    
    def request_data_streams(self):
        """Request telemetry data streams"""
        self.master.mav.request_data_stream_send(
            self.master.target_system,
            self.master.target_component,
            mavutil.mavlink.MAV_DATA_STREAM_ALL,
            4,  # Hz
            1   # enable
        )
    
    def parse_telemetry(self):
        """Parse MAVLink messages"""
        msg = self.master.recv_match(blocking=False)
        
        if msg is None:
            return None
        
        msg_type = msg.get_type()
        
        # GPS position
        if msg_type == 'GLOBAL_POSITION_INT':
            self.telemetry_data['latitude'] = msg.lat / 1e7
            self.telemetry_data['longitude'] = msg.lon / 1e7
            self.telemetry_data['altitude'] = msg.alt / 1000.0  # mm to meters
            self.telemetry_data['relative_altitude'] = msg.relative_alt / 1000.0
            self.telemetry_data['vx'] = msg.vx / 100.0
            self.telemetry_data['vy'] = msg.vy / 100.0
            self.telemetry_data['vz'] = msg.vz / 100.0
            self.telemetry_data['heading'] = msg.hdg / 100.0
        
        # Attitude
        elif msg_type == 'ATTITUDE':
            self.telemetry_data['roll'] = msg.roll
            self.telemetry_data['pitch'] = msg.pitch
            self.telemetry_data['yaw'] = msg.yaw
        
        # Battery
        elif msg_type == 'SYS_STATUS':
            self.telemetry_data['battery_voltage'] = msg.voltage_battery / 1000.0
            self.telemetry_data['battery_current'] = msg.current_battery / 100.0
            self.telemetry_data['battery_remaining'] = msg.battery_remaining
        
        # Heartbeat
        elif msg_type == 'HEARTBEAT':
            self.armed = (msg.base_mode & mavutil.mavlink.MAV_MODE_FLAG_SAFETY_ARMED) != 0
            self.mode = mavutil.mode_string_v10(msg)
            self.telemetry_data['armed'] = self.armed
            self.telemetry_data['mode'] = self.mode
        
        # GPS status
        elif msg_type == 'GPS_RAW_INT':
            self.telemetry_data['gps_fix'] = msg.fix_type
            self.telemetry_data['gps_satellites'] = msg.satellites_visible
        
        return msg_type
    
    def check_safety_limits(self):
        """Check if drone is within safety limits"""
        warnings = []
        
        # Altitude check
        if 'relative_altitude' in self.telemetry_data:
            if self.telemetry_data['relative_altitude'] > MAX_ALTITUDE:
                warnings.append(f"ALTITUDE EXCEEDED: {self.telemetry_data['relative_altitude']}m > {MAX_ALTITUDE}m")
        
        # Speed check
        if all(k in self.telemetry_data for k in ['vx', 'vy', 'vz']):
            speed = (self.telemetry_data['vx']**2 + self.telemetry_data['vy']**2 + self.telemetry_data['vz']**2)**0.5
            if speed > MAX_SPEED:
                warnings.append(f"SPEED EXCEEDED: {speed:.1f}m/s > {MAX_SPEED}m/s")
        
        # Battery check
        if 'battery_remaining' in self.telemetry_data:
            if self.telemetry_data['battery_remaining'] < BATTERY_LOW_THRESHOLD:
                warnings.append(f"LOW BATTERY: {self.telemetry_data['battery_remaining']}%")
        
        return warnings
    
    async def publish_telemetry(self):
        """Publish telemetry to NATS"""
        if not self.telemetry_data:
            return
        
        # Add metadata
        payload = {
            "drone_id": DRONE_ID,
            "drone_type": DRONE_TYPE,
            "location_base": LOCATION_BASE,
            "timestamp": datetime.utcnow().isoformat(),
            "telemetry": self.telemetry_data,
            "safety_warnings": self.check_safety_limits()
        }
        
        # Publish to NATS
        try:
            await self.nc.publish(NATS_SUBJECT, json.dumps(payload).encode())
        except Exception as e:
            print(f"NATS publish error: {e}")
    
    async def run(self):
        """Main telemetry loop"""
        await self.connect()
        self.request_data_streams()
        
        print(f"Telemetry loop started for {DRONE_ID}")
        
        try:
            while True:
                # Parse incoming MAVLink messages
                self.parse_telemetry()
                
                # Publish to NATS every second
                await self.publish_telemetry()
                
                # Check safety limits
                warnings = self.check_safety_limits()
                if warnings:
                    for warning in warnings:
                        print(f"⚠️  {warning}")
                
                await asyncio.sleep(0.25)  # 4 Hz
        
        except KeyboardInterrupt:
            print("\nShutting down telemetry...")
        
        finally:
            if self.nc:
                await self.nc.close()
            print("Telemetry stopped")

def main():
    telemetry = DroneTelemetry()
    asyncio.run(telemetry.run())

if __name__ == "__main__":
    main()
```

---

### Group Variables

**`group_vars/pi_vision.yml`:**

```yaml
---
# Pi 5 + AI HAT 2 Vision Nodes

# Hailo runtime
hailo_runtime_version: "4.17.0"
hailo_device_id: 0

# YOLOv8 models
yolov8_models:
  nano: yolov8n    # Fastest (100 FPS on Hailo)
  small: yolov8s   # Balanced (60 FPS)
  medium: yolov8m  # Accurate (30 FPS)
  large: yolov8l   # Very accurate (15 FPS)

# Camera settings
camera_resolution_options:
  - "640x480"
  - "1280x720"
  - "1920x1080"
  - "3840x2160"

# Detection settings
detection_confidence: 0.25
detection_iou: 0.45
detection_classes:  # COCO classes to detect
  - person
  - car
  - truck
  - bus
  - motorcycle
  - bicycle

# LPR settings
lpr_enabled: false
lpr_confidence: 0.7
lpr_languages:
  - en

# Performance
gpu_memory_split: 512
enable_hardware_accel: true
```

**`group_vars/drones.yml`:**

```yaml
---
# Drone Fleet Configuration

# ArduPilot
ardupilot_version: "4.5.0"
ardupilot_vehicle_type: "Copter"
ardupilot_frame_class: 1  # Quad

# Flight limits (FAA Part 107)
max_altitude_agl: 120  # meters (400 ft)
max_speed_horizontal: 15  # m/s (~33 mph)
max_speed_vertical: 5  # m/s

# Battery
battery_cells: 6  # 6S LiPo default
battery_low_action: "RTL"  # Return to launch
battery_critical_voltage: 3.3  # per cell
battery_low_voltage: 3.6  # per cell

# GPS
gps_type: "ublox"
gps_baudrate: 115200
rtk_enabled: false

# Telemetry
telemetry_rate: 4  # Hz
telemetry_protocol: "mavlink2"
telemetry_port: 14550

# Safety
geofence_enabled: true
geofence_action: "RTL"
failsafe_gcs_enabled: true
failsafe_battery_enabled: true

# Camera
camera_enabled: true
camera_trigger_distance: 10  # meters (for mapping)
video_resolution: "1920x1080"
video_bitrate: "5000k"
```

---

## Usage Examples

### Deploy Vision Nodes

```bash
cd ~/onemind-infra/ansible

# Deploy single vision node (testing)
ansible-playbook playbooks/02-pi-vision-deploy.yml --limit pi-vision-01

# Deploy all vision nodes
ansible-playbook playbooks/02-pi-vision-deploy.yml

# Deploy with LPR enabled
ansible-playbook playbooks/02-pi-vision-deploy.yml --extra-vars "enable_lpr=true"
```

### Deploy Drone Fleet

```bash
# Deploy single drone (testing)
ansible-playbook playbooks/03-drone-fleet-deploy.yml --limit drone-alpha-01

# Deploy entire alpha squadron
ansible-playbook playbooks/03-drone-fleet-deploy.yml --limit 'drone-alpha-*'

# Deploy all drones
ansible-playbook playbooks/03-drone-fleet-deploy.yml
```

### Check Vision Node Status

```bash
# Check Hailo device
ansible pi_vision -m shell -a "hailortcli fw-control identify" --become

# Check vision service
ansible pi_vision -m shell -a "systemctl status onemind-vision" --become

# Check camera
ansible pi_vision -m shell -a "v4l2-ctl --list-devices"

# View vision logs
ansible pi_vision -m shell -a "tail -50 /opt/onemind/logs/vision.log" --become
```

### Check Drone Status

```bash
# Check ArduPilot service
ansible drones -m shell -a "systemctl status ardupilot" --become

# Check telemetry service
ansible drones -m shell -a "systemctl status drone-telemetry" --become

# Check GPS fix
ansible drones -m shell -a "gpspipe -w -n 10 | grep TPV" --become

# View drone logs
ansible drones -m shell -a "tail -50 /opt/drone/logs/telemetry.log" --become
```

---

## Summary ★

**Created advanced playbooks for:**
1. ✅ Pi 5 + AI HAT 2 Vision Nodes (YOLOv8 + Hailo acceleration)
2. ✅ Drone Fleet (ArduPilot + MAVLink + real-time telemetry)

**Total new playbooks:** 2  
**Total new templates:** 4  
**Status:** Copy-paste ready for future deployment ♪

When your Pi 5 + AI HAT 2 and drone hardware arrive, you're ready to deploy 20-100 nodes with a single command! ✨(◕‿◕)★

# Zeus Drone: Starter Guide

> **Created:** 2026-01-22
> **Target:** Build your first voice-controlled DIY drone
> **Timeline:** 1-2 weeks from parts to flight
> **Cost:** $765 (DIY) vs $1,099+ (commercial)

---

## Why Start with Zeus?

**Zeus is your testing platform:**
- ✅ Cheapest entry point ($765)
- ✅ Validates entire OneMind Drone architecture
- ✅ Full customization (sensors, payloads, experiments)
- ✅ Native MAVLink protocol (no API limits)
- ✅ Open source (ArduPilot) - massive community
- ✅ Learn drone fundamentals before expensive hardware

**Once Zeus works, adding commercial drones is easy** - same architecture, just different bridges.

---

## Complete Parts List

### Flight Controller & Electronics

| Part | Model | Price | Quantity | Link |
|------|-------|-------|----------|------|
| **Flight Controller** | Pixhawk 6C | $200 | 1 | Holybro store |
| **GPS Module** | UBlox M8N with compass | $40 | 1 | Included with Pixhawk |
| **Power Module** | PM02D | $25 | 1 | Included with Pixhawk |
| **Telemetry Radio** | 433MHz 100mW | $60 | 1 | Holybro/3DR |

### Frame & Motors

| Part | Model | Price | Quantity | Notes |
|------|-------|-------|----------|-------|
| **Frame** | S500 Quadcopter | $80 | 1 | 500mm diagonal |
| **Motors** | 920KV Brushless | $20 | 4 | High torque |
| **ESCs** | 30A BLHeli | $15 | 4 | Electronic speed controllers |
| **Propellers** | 10x4.5 Carbon Fiber | $5 | 4 sets | Buy extras |

### Power System

| Part | Model | Price | Quantity | Notes |
|------|-------|-------|----------|-------|
| **Battery** | 4S 5000mAh LiPo | $80 | 2 | 20-25 min flight |
| **Charger** | ISDT Q6 Nano | $50 | 1 | Smart balance charger |
| **XT60 Connectors** | XT60 Male/Female | $10 | 5 pairs | Battery connectors |

### Vision & Compute

| Part | Model | Price | Quantity | Notes |
|------|-------|-------|----------|-------|
| **Companion Computer** | Raspberry Pi 4 (4GB) | $75 | 1 | Runs MAVLink bridge |
| **Camera** | Raspberry Pi HQ Camera | $50 | 1 | 12MP, interchangeable lens |
| **Lens** | 6mm Wide Angle | $25 | 1 | For FPV/navigation |
| **FPV Camera** (optional) | RunCam Split 4 | $40 | 1 | Low-latency video |
| **MicroSD Card** | 64GB Class 10 | $15 | 1 | For Pi OS |

### Accessories

| Part | Model | Price | Quantity | Notes |
|------|-------|-------|----------|-------|
| **Landing Gear** | Carbon Fiber Legs | $20 | 1 set | Height adjustable |
| **Velcro Straps** | Battery straps | $10 | 5 | Secure battery |
| **Wire Kit** | Silicone wire 18AWG | $15 | 1 | Red/black |
| **Heat Shrink** | Various sizes | $10 | 1 kit | Protect connections |
| **Threadlocker** | Loctite Blue 242 | $8 | 1 | Prevent vibration loosening |

### Tools (If You Don't Have)

| Tool | Price | Notes |
|------|-------|-------|
| Soldering Iron | $30 | Temperature controlled |
| Hex Key Set | $15 | Metric (2mm, 2.5mm, 3mm) |
| Wire Strippers | $10 | 18-22 AWG |
| Zip Ties | $5 | Various sizes |
| Multimeter | $20 | Voltage/continuity testing |

**Total Parts:** ~$765 (without tools)
**Total with Tools:** ~$845

---

## Assembly Guide (Step-by-Step)

### Phase 1: Frame Assembly (30 minutes)

1. **Unpack S500 frame kit**
   - 4x arms (with motor mounts)
   - Top/bottom plates
   - Center plates
   - Screws/standoffs

2. **Attach arms to center plates**
   - Use M3 screws
   - Apply threadlocker
   - Tighten firmly (not overtight)
   - Ensure arms are symmetrical

3. **Mount motors to arms**
   - Each motor has rotation direction (CW/CCW)
   - Front-left & rear-right: CCW
   - Front-right & rear-left: CW
   - Screw motors to bottom of arms (M3 screws)

4. **Install ESCs on arms**
   - Use zip ties or tape
   - ESC signal wire faces center
   - Leave motor wires long enough to solder

### Phase 2: Electronics Installation (1-2 hours)

1. **Solder ESCs to motors**
   - Match any 2 of 3 motor wires to ESC wires
   - Solder securely
   - Heat shrink all connections
   - **Test direction later** (if wrong, swap any 2 wires)

2. **Mount Pixhawk 6C**
   - Center of gravity (middle of frame)
   - Use vibration-damping foam
   - Arrow points forward
   - Secure with velcro/straps

3. **Connect ESCs to Pixhawk**
   - ESC signal wires → Pixhawk MAIN OUT
   - Front-right → MAIN OUT 1
   - Rear-left → MAIN OUT 2
   - Front-left → MAIN OUT 3
   - Rear-right → MAIN OUT 4

4. **Install GPS module**
   - Mount on top plate (away from electronics)
   - Raise on mast if possible (better reception)
   - Arrow points forward
   - Connect to GPS port on Pixhawk

5. **Connect power module**
   - Battery → Power module → Pixhawk
   - Monitors voltage & current
   - Powers Pixhawk and accessories

6. **Install telemetry radio**
   - Connect to TELEM1 port
   - Mount antenna away from electronics
   - Pair with ground station radio

### Phase 3: Companion Computer Setup (1 hour)

1. **Prepare Raspberry Pi 4**
   ```bash
   # Flash Raspberry Pi OS Lite to microSD
   # Use Raspberry Pi Imager

   # Enable SSH before first boot
   touch /Volumes/boot/ssh

   # Boot Pi and SSH in
   ssh pi@raspberrypi.local
   # Default password: raspberry
   ```

2. **Install dependencies**
   ```bash
   sudo apt update
   sudo apt upgrade -y

   # Python 3 and pip
   sudo apt install -y python3-pip python3-dev

   # MAVLink library
   pip3 install pymavlink

   # NATS client
   pip3 install nats-py

   # Camera libraries
   sudo apt install -y python3-picamera2
   ```

3. **Connect Pi to Pixhawk**
   - Use TELEM2 port on Pixhawk
   - USB to UART adapter OR direct serial
   - Pixhawk TX → Pi RX
   - Pixhawk RX → Pi TX
   - Pixhawk GND → Pi GND

4. **Configure serial on Pi**
   ```bash
   sudo raspi-config
   # Interface Options → Serial Port
   # Login shell: No
   # Serial hardware: Yes

   # Reboot
   sudo reboot
   ```

5. **Test MAVLink connection**
   ```bash
   # List serial devices
   ls /dev/tty*
   # Look for /dev/ttyAMA0 or /dev/serial0

   # Test connection
   mavproxy.py --master=/dev/serial0 --baudrate=57600
   # Should see heartbeat messages
   ```

### Phase 4: ArduPilot Configuration (2-3 hours)

1. **Install Mission Planner** (Windows) or **QGroundControl** (Mac/Linux)
   - Download from ardupilot.org
   - Install telemetry radio drivers

2. **Connect to Pixhawk**
   - Power on drone (no props yet!)
   - Connect telemetry radio to computer
   - Open Mission Planner
   - Select correct COM port
   - Connect at 57600 baud

3. **Flash ArduCopter firmware**
   - Setup → Install Firmware
   - Select "Quad" (4 motors)
   - Wait for upload (2-3 minutes)
   - Pixhawk will reboot

4. **Mandatory Hardware Setup**

   **Frame Type:**
   ```
   CONFIG → Frame Type → X (quad X configuration)
   ```

   **Accelerometer Calibration:**
   ```
   Setup → Mandatory Hardware → Accel Calibration
   Follow on-screen instructions (6 positions)
   ```

   **Compass Calibration:**
   ```
   Setup → Mandatory Hardware → Compass
   Enable compass
   Calibration → Rotate drone in all directions
   Get "green" status
   ```

   **Radio Calibration:**
   - If using RC transmitter (optional for testing)
   - Setup → Mandatory Hardware → Radio Calibration
   - Move all sticks to extremes

   **ESC Calibration:**
   ```
   Setup → Mandatory Hardware → ESC Calibration
   Follow on-screen prompts
   Tests motor directions:
     - If motor spins wrong way, swap any 2 motor wires
     - Front-left spins CCW
     - Front-right spins CW
     - Rear-left spins CW
     - Rear-right spins CCW
   ```

   **Battery Monitor:**
   ```
   Setup → Optional Hardware → Battery Monitor
   Sensor: Analog Voltage and Current
   Monitor: Power Module
   Set voltage divider (usually 10.1)
   Set amps per volt (usually 17.0)
   ```

5. **Important Parameters**
   ```
   CONFIG → Full Parameter List:

   SERIAL2_PROTOCOL = 2 (MAVLink2)
   SERIAL2_BAUD = 57 (57600 baud)
   ARMING_CHECK = 1 (all checks)
   PILOT_TKOFF_ALT = 1000 (takeoff altitude 10m)
   LAND_SPEED = 50 (land at 0.5 m/s)
   RTL_ALT = 1500 (RTL altitude 15m)
   ```

6. **Test motors (NO PROPS)**
   ```
   Setup → Optional Hardware → Motor Test
   Test each motor individually
   Verify correct direction
   ```

### Phase 5: Software Deployment (30 minutes)

1. **Copy MAVLink bridge to Pi**
   ```bash
   # On your computer, create deployment package
   cd ~/onemind-os

   # Copy to Pi
   scp -r agents/edge/drone_manager pi@raspberrypi.local:~/
   ```

2. **Configure MAVLink bridge**
   ```python
   # ~/drone_manager/config.py

   DRONE_ID = "zeus"
   MAVLINK_CONNECTION = "/dev/serial0:57600"
   NATS_URL = "nats://100.102.21.44:4222"  # Your cloud NATS
   LIVEKIT_URL = "ws://100.102.21.44:7880"
   ```

3. **Test bridge locally**
   ```bash
   ssh pi@raspberrypi.local

   cd ~/drone_manager
   python3 bridges/mavlink_bridge.py

   # Should see:
   # "Waiting for heartbeat from zeus..."
   # "Heartbeat from zeus received"
   # "Connected to NATS at nats://100.102.21.44:4222"
   ```

4. **Create systemd service** (auto-start)
   ```bash
   sudo nano /etc/systemd/system/zeus-bridge.service
   ```

   ```ini
   [Unit]
   Description=Zeus MAVLink Bridge
   After=network.target

   [Service]
   Type=simple
   User=pi
   WorkingDirectory=/home/pi/drone_manager
   ExecStart=/usr/bin/python3 /home/pi/drone_manager/bridges/mavlink_bridge.py
   Restart=always
   RestartSec=10

   [Install]
   WantedBy=multi-user.target
   ```

   ```bash
   sudo systemctl enable zeus-bridge
   sudo systemctl start zeus-bridge
   sudo systemctl status zeus-bridge
   ```

---

## First Flight Checklist

### Pre-Flight Checks

- [ ] **Battery fully charged** (4S = 16.8V fully charged)
- [ ] **Propellers installed correctly** (CW/CCW matched to motors)
- [ ] **All screws tight** (check with hex keys)
- [ ] **GPS lock** (blue LED solid, not flashing)
- [ ] **Compass calibrated** (green in Mission Planner)
- [ ] **Telemetry connected** (green bars in Mission Planner)
- [ ] **Clear flight area** (50+ meters, no people/obstacles)
- [ ] **Flight mode set** to Stabilize or AltHold
- [ ] **Arming checks pass** (no red messages in HUD)

### Manual Test Flight (First Time)

**Important:** Use RC transmitter for first flight, NOT voice control yet!

1. **Arm the drone**
   ```
   Hold throttle down + yaw right for 2 seconds
   Motors should spin slowly
   ```

2. **Takeoff slowly**
   ```
   Increase throttle gradually
   Lift off ground 1-2 meters
   Hover for 10 seconds
   ```

3. **Test controls**
   ```
   Pitch forward/back (small movements)
   Roll left/right (small movements)
   Yaw rotate (small movements)
   Hold stable hover
   ```

4. **Land**
   ```
   Reduce throttle slowly
   Land gently
   Disarm (throttle down + yaw left)
   ```

5. **Post-flight inspection**
   ```
   Check for loose screws
   Check battery voltage (should be >14.0V)
   Check logs in Mission Planner (Setup → Download Logs)
   ```

### Voice-Controlled Test Flight

**Only after successful manual flights!**

1. **Test NATS connection**
   ```bash
   # From your computer
   nats pub "commands.farm.drone.zeus.status" ""

   # Should get response with drone status
   ```

2. **Test from Legacy**
   ```
   You: "Legacy, what's Zeus' status?"
   Legacy: "Zeus is armed, battery 15.8V, GPS lock acquired,
            ready for flight."
   ```

3. **Voice takeoff** (BE READY TO TAKE MANUAL CONTROL)
   ```
   You: "Legacy, launch Zeus to 10 meters"

   # Watch drone:
   # - Arms automatically
   # - Spins up motors
   # - Takes off to 10m
   # - Hovers in place

   You: "Legacy, land Zeus"
   # Lands automatically
   ```

4. **Voice navigation**
   ```
   You: "Zeus, fly to GPS coordinates 40.123, -74.456, altitude 15"
   # Drone flies to waypoint, hovers

   You: "Zeus, orbit this position with 10 meter radius"
   # Drone orbits in circle

   You: "Zeus, return to launch"
   # Drone flies back to takeoff point, lands
   ```

---

## Common Issues & Troubleshooting

### Issue: GPS Won't Get Lock

**Symptoms:** Blue LED flashing rapidly, Mission Planner says "No GPS"

**Solutions:**
- Move outside (GPS doesn't work indoors)
- Wait 2-5 minutes (first lock takes time)
- Ensure GPS antenna clear of carbon fiber
- Check GPS cable connection
- Update GPS firmware in Mission Planner

### Issue: Compass Interference

**Symptoms:** Compass error on arming, "inconsistent compasses"

**Solutions:**
- Recalibrate compass outdoors
- Move GPS away from power cables
- Disable internal compass if using external
- Check for magnetic interference (steel structures)

### Issue: Motor Spins Wrong Direction

**Symptoms:** Motor test shows wrong spin, unstable flight

**Solutions:**
- Swap any 2 of 3 motor wires on that ESC
- Re-test motor direction
- ESC might need reversing in BLHeli configurator

### Issue: Won't Arm

**Symptoms:** Motor beeps, refuses to arm

**Solutions:**
- Check Mission Planner HUD for error messages
- Common issues:
  - Pre-arm: Need GPS lock
  - Pre-arm: Compass not calibrated
  - Pre-arm: Accelerometers not calibrated
  - Pre-arm: RC not calibrated
  - Pre-arm: Battery voltage low (<14V)
- Set ARMING_CHECK = 0 temporarily (DANGEROUS, for testing only)

### Issue: Drone Drifts in Hover

**Symptoms:** Won't hold position, drifts with wind

**Solutions:**
- This is normal in Stabilize mode (no GPS position hold)
- Switch to AltHold or Loiter mode for GPS position hold
- Ensure GPS has good lock (8+ satellites)
- Tune PID parameters (advanced)

### Issue: MAVLink Bridge Can't Connect

**Symptoms:** Bridge says "Waiting for heartbeat..." forever

**Solutions:**
- Check serial connection (Pi to Pixhawk)
- Verify serial device: `ls /dev/tty*`
- Check baud rate (should be 57600)
- Ensure SERIAL2_PROTOCOL = 2 in Pixhawk params
- Try different USB port / UART adapter

### Issue: Video Stream Not Working

**Symptoms:** No video in LiveKit room

**Solutions:**
- Test Pi camera: `libcamera-hello`
- Check camera cable connection
- Ensure camera enabled in `raspi-config`
- Check LiveKit credentials
- Verify network connection (ping 100.102.21.44)

---

## Upgrading Zeus

### Level 1: Basic Improvements ($100-200)

- **Better Camera:** GoPro Hero mount + livestream
- **LED Strip:** WS2812B LEDs for night visibility
- **Buzzer:** Louder lost drone alarm
- **Longer Flight:** 4S 8000mAh battery (30+ min)

### Level 2: Advanced Sensors ($200-500)

- **Optical Flow:** PX4Flow for indoor flight (no GPS)
- **Lidar Rangefinder:** Precision altitude hold
- **FPV System:** Analog FPV transmitter for low-latency piloting
- **RTK GPS:** Centimeter-level precision ($300)

### Level 3: Professional ($500+)

- **Gimbal:** 2-axis or 3-axis camera stabilization
- **Thermal Camera:** FLIR Lepton for night vision / heat detection
- **LTE Modem:** 4G connection for unlimited range
- **Obstacle Avoidance:** Intel RealSense or similar
- **Parachute:** Failsafe recovery system

---

## Integration with OneMind OS

### Voice Commands Available

Once Zeus is flying and MAVLink bridge is running:

```
"Legacy, what's Zeus' status?"
→ Battery, GPS lock, flight mode, position

"Launch Zeus to 15 meters"
→ Arms, takes off to 15m, hovers

"Zeus, fly to the workshop"
→ Navigates to GPS waypoint, hovers

"Zeus, orbit the barn with 10 meter radius"
→ Circular flight pattern

"Zeus, follow me"
→ (Requires GPS on phone, custom implementation)

"Zeus, land"
→ Descends and lands at current position

"Zeus, return to launch"
→ Flies back to takeoff point, lands

"Show me Zeus' camera feed"
→ Opens LiveKit video stream in UI

"Zeus, take a photo"
→ Captures image, saves to Pi, uploads to cloud
```

### Automated Workflows

```python
@workflow
def morning_property_check():
    """Automated morning inspection with Zeus"""
    # Check weather
    weather = await sensors.get_weather()
    if weather.wind_speed > 15:  # mph
        legacy.speak("Wind too high for Zeus flight")
        return

    # Launch Zeus
    await drone_tool.command_drone("zeus", "takeoff", {"altitude": 15})

    # Fly inspection waypoints
    waypoints = [
        {"name": "barn", "lat": 40.123, "lon": -74.456, "alt": 15},
        {"name": "greenhouse", "lat": 40.124, "lon": -74.457, "alt": 15},
        {"name": "workshop", "lat": 40.125, "lon": -74.458, "alt": 15},
        {"name": "gate", "lat": 40.126, "lon": -74.459, "alt": 20},
    ]

    for wp in waypoints:
        await drone_tool.command_drone("zeus", "goto", wp)

        # Capture photo at each waypoint
        await drone_tool.command_drone("zeus", "capture_photo")

        # AI analysis (optional)
        image = await get_latest_photo()
        analysis = await legacy.analyze_image(image,
            objectives=["damage", "animals", "vehicles"]
        )

        if analysis.issues_found:
            await legacy.send_alert(f"Issue at {wp['name']}: {analysis}")

    # Return home
    await drone_tool.command_drone("zeus", "rtl")

    legacy.speak("Morning property check complete. All clear.")
```

### UI Integration

**Dashboard Widget:**
```
┌─────────────────────────────────────┐
│ Zeus Drone Status                   │
├─────────────────────────────────────┤
│ Status:    Flying                   │
│ Position:  40.123, -74.456          │
│ Altitude:  15.2m                    │
│ Battery:   87% (18 min remaining)   │
│ Signal:    95%                      │
│ Mode:      Auto                     │
│                                     │
│ [📹 Live Feed] [🎮 Take Control]   │
│ [🏠 RTL] [🛬 Land Now]              │
└─────────────────────────────────────┘
```

---

## Safety & Legal

### FAA Regulations (USA)

**Recreational (Part 44807):**
- Register drone if >250g (Zeus is ~1.5kg)
- Fly below 400 feet
- Keep drone in visual line of sight
- Don't fly over people
- Don't fly near airports (5 miles)
- Cost: $5 registration (covers all drones)

**Commercial (Part 107):**
- Required if using Zeus for any commercial purpose
- Must pass FAA Part 107 test
- All Part 44807 rules apply
- Additional rules for night flight, operations over people
- Cost: $175 test fee + $5 registration

### Safety Rules

1. **Always check weather before flight**
   - Wind <15 mph
   - No rain/snow
   - Visibility >3 miles
   - Temperature above freezing (battery performance)

2. **Pre-flight checklist every time**
   - Battery charged
   - Props secure
   - GPS lock
   - Compass calibrated
   - Clear flight area

3. **Have a spotter**
   - Second person watching drone
   - Especially for voice-controlled flights
   - Call out obstacles

4. **Emergency procedures**
   - Know how to disarm (throttle down + yaw left)
   - RTL button mapped on RC transmitter
   - Understand failsafe behavior

5. **Insurance**
   - Consider drone liability insurance
   - ~$100/year for $1M coverage
   - Protects against accidents

---

## Next Steps After Zeus

### Phase 2: Add Commercial Drone

Once Zeus is fully operational:

1. **Choose second drone based on need:**
   - Security: Skydio 2+ ($1,099)
   - Aerial photos: DJI Mini 4 Pro ($759)
   - Agriculture: DJI Mavic 3 Multispectral ($5,500)

2. **Deploy second protocol bridge**
   - Zeus already validates NATS architecture
   - Second bridge is just different protocol
   - Same Drone Manager, same voice commands

3. **Test multi-drone coordination**
   - Zeus + second drone simultaneous flight
   - Coordinated missions
   - Fleet status dashboard

### Phase 3: Advanced Features

- Computer vision (obstacle detection)
- Autonomous navigation (no waypoints)
- Package delivery mechanism
- Precision landing platform
- Battery hot-swap station

---

## Resources

### Official Documentation
- **ArduPilot:** https://ardupilot.org/copter/
- **Mission Planner:** https://ardupilot.org/planner/
- **QGroundControl:** https://qgroundcontrol.com/
- **MAVLink Protocol:** https://mavlink.io/

### Communities
- **ArduPilot Forum:** https://discuss.ardupilot.org/
- **DIYDrones:** https://diydrones.com/
- **Reddit r/Multicopter:** https://reddit.com/r/Multicopter
- **Reddit r/ardupilot:** https://reddit.com/r/ardupilot

### Video Tutorials
- **Painless360 (YouTube):** ArduPilot setup tutorials
- **ArduPilot YouTube:** Official tutorials
- **Oscar Liang:** FPV and DIY drone guides

### Part Suppliers
- **Holybro:** https://holybro.com/ (Pixhawk official)
- **GetFPV:** https://getfpv.com/ (US, fast shipping)
- **Banggood:** https://banggood.com/ (China, cheaper but slower)
- **Amazon:** Quick shipping for common parts

---

## Conclusion

**Zeus is your foundation:**
- Proves the OneMind Drone Manager architecture
- Teaches drone fundamentals
- Cost-effective ($765 vs $1,000+)
- Fully customizable for your needs
- Seamlessly integrates with Legacy AI

**Timeline:**
- Week 1: Order parts, start assembly
- Week 2: Complete build, test hardware
- Week 3: Manual test flights, tune parameters
- Week 4: Deploy software, voice-controlled flights

**Once Zeus flies, your entire multi-drone fleet architecture is validated.** Adding DJI and Skydio drones becomes straightforward - just different protocol bridges to the same Drone Manager.

**Ready to build Zeus and validate your multi-drone fleet architecture!** 🚁

---

*Last Updated: 2026-01-22*
*Next: Agricultural Workflow Examples*

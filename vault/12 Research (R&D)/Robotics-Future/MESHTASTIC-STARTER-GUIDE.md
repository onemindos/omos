# Meshtastic Starter Guide for OneMind OS

> **Get Started with Meshtastic Mesh Networking**
>
> Hardware selection, setup, and testing guide

---

## Table of Contents

1. [Quick Start Overview](#quick-start-overview)
2. [Hardware Options](#hardware-options)
3. [Recommended Starter Kits](#recommended-starter-kits)
4. [What You Need](#what-you-need)
5. [Setup Process](#setup-process)
6. [Testing Your Mesh](#testing-your-mesh)
7. [Troubleshooting](#troubleshooting)
8. [Next Steps](#next-steps)

---

## Quick Start Overview

**Goal**: Get a working 2-device Meshtastic mesh network in 1-2 hours

**What You'll Achieve:**
- ✅ Two devices communicating via LoRa (5-10km range)
- ✅ GPS tracking for both devices
- ✅ Text messaging working
- ✅ Understanding of mesh basics
- ✅ Foundation for OneMind OS integration

**Timeline:**
- 30 min: Unbox and charge devices
- 30 min: Flash firmware and configure
- 30 min: Test range and GPS
- 30 min: Test messaging and features

**Cost**: $130-180 for starter kit

---

## Hardware Options

### Device Categories

Meshtastic devices come in 3 main categories:

| Category | MCU | Power | Use Case | Price |
|----------|-----|-------|----------|-------|
| **nRF52** | nRF52840 | Ultra-low | Handhelds, solar | $50-70 |
| **ESP32-S3** | ESP32-S3 | Low | Fixed hubs, WiFi bridge | $30-50 |
| **ESP32** | ESP32 | Medium | Budget option | $25-40 |

**Recommendation**: Start with **nRF52** for best battery life.

### LoRa Radio Chips

| Chip | Generation | Range | Power | Recommendation |
|------|------------|-------|-------|----------------|
| **SX1262** | Latest | Best | Low | ✅ **Best choice** |
| **SX1276** | Older | Good | Medium | ⚠️ OK but dated |
| **LR1110** | Newest | Best+ | Lowest | 🔮 Future option |

**Always choose SX126x series** for new purchases.

### Key Features to Consider

**Must Have:**
- ✅ SX1262 or SX1280 radio chip
- ✅ GPS module (for tracking/mapping)
- ✅ Antenna connector (for external antenna)
- ✅ Battery connector (JST 2.0 or similar)

**Nice to Have:**
- Display (OLED) - shows messages and status
- Case/enclosure - weather protection
- Solar charging - for permanent outdoor deployment
- USB-C - modern charging standard

**Don't Need (For OneMind OS):**
- WiFi/Bluetooth - we use dedicated edge devices for bridging
- Touchscreen - not needed for embedded use
- Buttons - configuration via app/CLI

---

## Recommended Starter Kits

### Option 1: RAK Meshtastic Starter Kit ⭐ RECOMMENDED

**Best For**: Production deployment, rugged use, long battery life

**Hardware:**
- RAK19007 Base Board (WisBlock)
- RAK4631 Core Module (nRF52840 + SX1262)
- RAK1910 GPS Module (included)
- Battery connector (JST 2.0)
- IP65 enclosure (optional, $15 extra)

**Specs:**
- LoRa: SX1262 (915MHz US, 868MHz EU)
- MCU: nRF52840 (64 MHz, 1MB flash, 256KB RAM)
- GPS: u-blox MAX-7Q
- Power: 50mA active, <1mA sleep
- Battery: 18650 Li-ion (3.7V, 2500-3500mAh)
- Range: 5-10km typical, 30km+ with elevation

**Price:**
- Device: $60-70
- Battery (18650): $5-8 each (need 2)
- Enclosure: $15 (optional)
- **Total: $75-95 per device**

**Where to Buy:**
- [RAK Wireless Store](https://store.rakwireless.com)
- [Amazon](https://amazon.com/s?k=rak+meshtastic) (search "RAK Meshtastic")
- [AliExpress](https://aliexpress.com) ($5-10 cheaper, slower shipping)

**Pros:**
- ✅ Best battery life (weeks to months)
- ✅ Modular (can add sensors, solar, etc.)
- ✅ Rugged and reliable
- ✅ Easy to mount/deploy
- ✅ Official Meshtastic support

**Cons:**
- ⚠️ More expensive than ESP32 options
- ⚠️ Requires separate battery purchase
- ⚠️ No built-in display (add-on available)

**Recommended Quantity:**
- **2 devices minimum** for testing mesh
- **3-5 devices** for property coverage
- **10+ devices** for full sensor network

### Option 2: LILYGO T3-S3 💰 BUDGET OPTION

**Best For**: Budget builds, WiFi bridge hubs, learning

**Hardware:**
- LILYGO LoRa32 T3-S3
- ESP32-S3 (240MHz dual-core)
- SX1262 or SX1280 LoRa radio
- 0.96" OLED display (built-in)
- USB-C charging
- 18650 battery holder (battery not included)

**Specs:**
- LoRa: SX1262 (915MHz US)
- MCU: ESP32-S3 (240MHz, 8MB flash, 512KB RAM)
- Display: 0.96" OLED (128×64)
- WiFi: 802.11 b/g/n (2.4GHz)
- Bluetooth: BLE 5.0
- GPS: No (add external via UART)
- Power: ~150mA active, ~20mA sleep (higher than nRF52)
- Battery: 18650 Li-ion
- Range: 5-8km typical

**Price:**
- Device: $30-40
- Battery (18650): $5-8
- GPS module (optional): $15-25
- **Total: $35-75 per device**

**Where to Buy:**
- [LILYGO Official Store (AliExpress)](https://lilygo.aliexpress.com)
- [Amazon](https://amazon.com/s?k=lilygo+t3-s3)
- [Tindie](https://tindie.com) (maker marketplace)

**Pros:**
- ✅ Cheapest option with good features
- ✅ Built-in display (shows messages)
- ✅ WiFi (for MQTT bridge)
- ✅ USB-C charging (modern)
- ✅ Good for learning/testing

**Cons:**
- ⚠️ Higher power consumption than nRF52
- ⚠️ No GPS (need external module)
- ⚠️ Less rugged (needs weatherproof case for outdoor)
- ⚠️ WiFi can interfere with LoRa in rare cases

**Recommended Quantity:**
- **2 devices minimum** for testing
- **1-2 as WiFi bridge hubs** (fixed locations)
- **RAK for mobile assets** (drones, robots) - better battery life

### Option 3: Heltec Mesh Node T114

**Best For**: Compact handhelds, middle ground option

**Hardware:**
- Heltec WiFi LoRa 32 (V3)
- nRF52840 MCU
- SX1262 LoRa radio
- 0.96" OLED display
- GPS (u-blox)
- USB-C

**Specs:**
- Similar to RAK but more integrated
- Good battery life (nRF52 based)
- Compact form factor
- Built-in GPS

**Price:**
- Device: $55-65
- Battery: $5-8
- **Total: $60-75 per device**

**Where to Buy:**
- [Heltec Official Store](https://heltec.org)
- [Amazon](https://amazon.com/s?k=heltec+meshtastic)

**Pros:**
- ✅ Good balance of features and price
- ✅ Low power (nRF52)
- ✅ Built-in GPS and display
- ✅ Compact

**Cons:**
- ⚠️ Less modular than RAK
- ⚠️ Smaller community than RAK/LILYGO

### Option 4: DIY Build (Advanced)

**Best For**: Custom applications, learning electronics

**Parts List:**
- ESP32-S3 Dev Board: $10
- RFM95W LoRa Module (SX1276): $8
- GPS Module (NEO-6M): $10
- Battery Management: $5
- Enclosure: $10
- **Total: $45-50**

**Not Recommended for Beginners** - requires soldering, firmware building, debugging

---

## What You Need

### Starter Kit (2 Devices)

**Hardware:**
| Item | Qty | Purpose | Est. Cost |
|------|-----|---------|-----------|
| Meshtastic devices | 2 | Mesh network | $120-140 |
| 18650 batteries | 4 | Power (2 per device) | $20-30 |
| Battery charger | 1 | Charge 18650s | $15-20 |
| USB cables | 2 | Programming/charging | $5-10 |
| Antennas (optional) | 2 | Better range | $10-20 |
| **Total** | | | **$170-220** |

**Software:**
- Meshtastic mobile app (iOS/Android) - Free
- Meshtastic Python CLI - Free
- Web Flasher (browser-based) - Free

**Tools:**
- Smartphone (iOS or Android)
- Computer (optional, for CLI setup)
- No soldering required ✅

### Full Deployment (10+ Devices)

**Additional Hardware:**
| Item | Qty | Purpose | Est. Cost |
|------|-----|---------|-----------|
| Solar panels | 5-10 | Outdoor sensors | $150-300 |
| Weatherproof cases | 5-10 | Protection | $75-150 |
| Magnetic mounts | 5-10 | Vehicle mounting | $50-100 |
| External antennas | 3-5 | Extended range | $30-100 |
| **Total** | | | **$305-650** |

---

## Setup Process

### Step 1: Unbox and Inspect (5 min)

**RAK Starter Kit:**
```
Box Contents:
├─ RAK19007 Base Board (green PCB)
├─ RAK4631 Core Module (small black module)
├─ RAK1910 GPS Module (tiny white module)
├─ USB cable (USB-A to micro-USB)
├─ Antenna (short wire with connector)
└─ Documentation (quick start guide)
```

**Checklist:**
- [ ] All components present
- [ ] No physical damage
- [ ] Antenna connector intact
- [ ] USB port clean

### Step 2: Assemble Hardware (10 min)

**RAK Assembly:**

1. **Insert Core Module**
   - Align RAK4631 with base board slots
   - Press firmly until seated
   - Should click into place

2. **Insert GPS Module** (optional but recommended)
   - Align RAK1910 with GPS slot
   - Press firmly until seated
   - GPS must be exposed to sky for fix

3. **Attach Antenna**
   - Screw antenna onto RP-SMA connector
   - Finger-tight (don't over-tighten)
   - ⚠️ **NEVER power on without antenna** (can damage radio)

4. **Insert Battery**
   - Connect 18650 to JST 2.0 connector
   - Observe polarity (red = +, black = -)
   - ⚠️ **Reversed polarity destroys device**

**LILYGO Assembly:**

1. **Insert Battery**
   - Slide 18650 into holder
   - Observe polarity markings
   - Springs should make contact

2. **Attach Antenna**
   - Screw onto U.FL or SMA connector
   - ⚠️ **NEVER power on without antenna**

### Step 3: Charge Batteries (30-60 min)

**Method 1: USB Charging (if supported)**
- Plug USB cable into device
- Red LED = charging
- Green LED = fully charged
- Takes 2-3 hours for empty 3000mAh battery

**Method 2: External Charger (recommended)**
- Remove battery from device
- Insert into 18650 charger
- Charge to 4.2V (full)
- Faster and safer for battery longevity

**Safety:**
- ⚠️ Never leave charging unattended
- ⚠️ Use quality batteries (Samsung, LG, Panasonic)
- ⚠️ Avoid cheap no-name batteries (fire risk)

### Step 4: Flash Firmware (20 min)

**Option A: Web Flasher (Easiest)**

1. **Visit Web Flasher**
   - Open: [https://flasher.meshtastic.org](https://flasher.meshtastic.org)
   - Works in Chrome/Edge (WebSerial required)

2. **Connect Device**
   - Plug USB cable into device
   - Click "Connect" in web flasher
   - Select serial port (usually `/dev/ttyUSB0` or `COM3`)

3. **Select Firmware**
   - Choose device type (e.g., "RAK4631")
   - Select latest stable version (e.g., 2.3.15)
   - Click "Flash"

4. **Wait for Flash**
   - Progress bar shows status
   - Takes 2-5 minutes
   - Device reboots automatically when done

5. **Verify**
   - Device LED should blink
   - Screen shows "Meshtastic" (if has display)

**Option B: Python CLI (Advanced)**

```bash
# Install Meshtastic CLI
pip install meshtastic

# Flash firmware
meshtastic --flash /path/to/firmware.bin

# Verify
meshtastic --info
```

**Repeat for Both Devices**

### Step 5: Basic Configuration (20 min)

**Using Mobile App:**

1. **Install App**
   - iOS: [App Store](https://apps.apple.com/app/meshtastic/id1586432531)
   - Android: [Play Store](https://play.google.com/store/apps/details?id=com.geeksville.mesh)

2. **Connect to Device**
   - Open app
   - Tap "+" to add device
   - Select Bluetooth device (e.g., "Meshtastic_a1b2")
   - Pair (no PIN required)

3. **Set Device Name**
   - Tap Settings → Device
   - Set friendly name (e.g., "House Hub", "Test Node 1")

4. **Configure Region**
   - Settings → LoRa
   - Region: US (915MHz) or EU (868MHz)
   - ⚠️ **Must match your country** (legal requirement)

5. **Configure Channel**
   - Settings → Channels → Primary
   - Name: "onemind-test" (or your preference)
   - Encryption: "Default" (AES-256)
   - Leave PSK as default for now

6. **Enable GPS** (if available)
   - Settings → Position
   - GPS Enabled: ON
   - Update Interval: 60 seconds
   - Broadcast: ON

7. **Save Settings**
   - Tap "Save" or "Apply"
   - Device reboots with new config

**Repeat for Second Device with Same Settings**

**Critical**: Both devices MUST have:
- Same region (e.g., both US)
- Same channel name
- Same PSK (encryption key)
- Same modem preset (default: Long Fast)

### Step 6: Test Basic Mesh (15 min)

**Indoor Test (Short Range):**

1. **Verify Connection**
   - Open app on phone paired with Device 1
   - Check "Nodes" tab
   - Should see both Device 1 (green) and Device 2 (gray initially)
   - Wait 30-60 seconds for Device 2 to appear

2. **Send Test Message**
   - Tap "Messages" tab
   - Tap "Broadcast" (or create DM)
   - Type: "Test message"
   - Send
   - Should appear on both devices within 1-5 seconds

3. **Verify GPS** (if enabled)
   - Tap "Map" tab
   - Both devices should show on map
   - May take 2-5 minutes for initial GPS fix
   - Must be outdoors or near window

**Success Criteria:**
- ✅ Both devices visible in node list
- ✅ Messages delivered within 5 seconds
- ✅ GPS showing locations (if outdoors)
- ✅ Signal strength > -100 dBm (check in node details)

**If Not Working:**
- Check both devices are on same channel
- Verify region matches
- Check encryption key (PSK) is identical
- Try moving devices closer (<10m)
- Check antenna is attached
- Restart both devices

---

## Testing Your Mesh

### Test 1: Range Test (30 min)

**Goal**: Determine actual range on your property

**Setup:**
1. **Device 1 (Base)**: Place at house (fixed location)
2. **Device 2 (Mobile)**: Take with you on walk

**Procedure:**
1. Start at Device 1 location
2. Verify both devices connected (via app)
3. Walk away slowly (check messages every 50-100m)
4. Send message every stop: "Test from 100m"
5. Note when signal degrades:
   - SNR < 0 dB = Weak but working
   - SNR < -10 dB = Very weak (unreliable)
   - No response = Out of range

**Record Results:**
```
Distance | SNR | RSSI | Message Success
---------|-----|------|----------------
0m       | +10 | -30  | ✅ Instant
100m     | +5  | -65  | ✅ <1s
500m     | 0   | -95  | ✅ 2-3s
1km      | -5  | -110 | ✅ 5s
2km      | -8  | -120 | ⚠️  50% loss
3km      | -   | -    | ❌ No signal
```

**Expected Results (Typical):**
- **Flat open terrain**: 5-10km
- **Residential (houses/trees)**: 1-3km
- **Dense urban**: 500m-1km
- **Indoors to outdoors**: 50-200m

**Tips for Better Range:**
- Mount Device 1 high (roof/attic)
- Use external antenna (3dBi → 5dBi)
- Test from elevated locations
- Avoid metal buildings/valleys
- Clear line of sight helps dramatically

### Test 2: Mesh Routing (30 min)

**Goal**: Verify messages hop through intermediate nodes

**Setup:**
```
Device 1 (House) ←--1km--→ Device 2 (Barn) ←--1km--→ Device 3 (Field)
```

**Procedure:**
1. Place Device 2 between Device 1 and Device 3
2. Ensure Device 1 and Device 3 are OUT of direct range (>3km apart)
3. Send message from Device 1 to Device 3
4. Verify message arrives (routed through Device 2)

**Verification:**
- App shows message path: `1 → 2 → 3`
- Hop count = 2 (visible in message metadata)
- Latency = 5-15s (multiple hops)

**Success**: Messages delivered even when sender and recipient can't directly reach each other.

### Test 3: GPS Tracking (30 min)

**Goal**: Verify GPS accuracy and position updates

**Setup:**
1. Take Device 2 outdoors (clear sky view)
2. Monitor from Device 1 (via app on phone)

**Procedure:**
1. Wait for GPS fix (2-5 min initial, LED blinks when locked)
2. Walk predetermined path (e.g., around property)
3. Check map in app - should show movement
4. Note accuracy (typical: 3-10m)

**Success Criteria:**
- ✅ GPS fix within 5 minutes
- ✅ Position updates every 60s (or configured interval)
- ✅ Accuracy within 10m
- ✅ Path visible on map

**Troubleshooting GPS:**
- No fix after 10 min: Check GPS antenna placement (must see sky)
- Inaccurate: Wait longer (accuracy improves over time)
- Not updating: Check GPS enabled in settings
- Indoor: Won't work (GPS needs satellite view)

### Test 4: Battery Life (24-48 hours)

**Goal**: Measure real-world battery consumption

**Setup:**
1. Fully charge battery (4.2V)
2. Configure device for normal operation:
   - GPS: Enabled (update every 60s)
   - Telemetry: Enabled (every 5 min)
   - Screen: Auto-off after 60s (if applicable)
3. Leave device running for 24-48 hours

**Monitoring:**
- Check battery voltage periodically (via app)
- Note hours until battery warning (<3.5V)

**Expected Results:**
| Device Type | Config | Battery Life |
|-------------|--------|--------------|
| RAK4631 (nRF52) | GPS + telemetry | 3-7 days |
| RAK4631 (nRF52) | No GPS | 1-2 weeks |
| LILYGO (ESP32-S3) | GPS + WiFi | 1-2 days |
| LILYGO (ESP32-S3) | No WiFi | 2-4 days |

**Optimization:**
- Disable GPS if not needed (huge power savings)
- Increase telemetry interval (5 min → 30 min)
- Reduce transmit power if range is overkill (30dBm → 20dBm)
- Add solar panel for permanent deployment

### Test 5: Text Messaging (15 min)

**Goal**: Test communication features

**Test Cases:**
1. **Broadcast Message**
   - Send to all nodes
   - Should reach all devices on mesh

2. **Direct Message**
   - Send to specific device
   - Only recipient sees it (but message still routes through mesh)

3. **Long Message**
   - Send 200+ character message
   - Verify delivery (may split into multiple packets)

4. **Rapid Messages**
   - Send 10 messages quickly
   - Check all delivered (may have delay due to rate limiting)

**Success Criteria:**
- ✅ All message types work
- ✅ Latency < 5s (single hop), < 15s (multi-hop)
- ✅ No lost messages (or <5% loss acceptable)

---

## Troubleshooting

### Issue: Devices Not Connecting

**Symptoms:**
- Device not visible in node list
- Messages not delivered
- App shows "Disconnected"

**Checks:**
1. **Same Channel?**
   ```
   Settings → Channels → Primary
   Both devices must have identical:
   - Name (case-sensitive)
   - PSK (encryption key)
   - Modem preset (e.g., Long Fast)
   ```

2. **Same Region?**
   ```
   Settings → LoRa → Region
   Must match (US = 915MHz, EU = 868MHz)
   ⚠️ Device in US region can't talk to EU region!
   ```

3. **Antenna Attached?**
   - Check physical connection
   - Try different antenna

4. **In Range?**
   - Move devices closer (<10m) for testing
   - Check signal strength (RSSI)

5. **Firmware Match?**
   - Both devices on same firmware version
   - Update if mismatch

**Quick Fix:**
- Factory reset both devices
- Reconfigure from scratch with identical settings

### Issue: Poor Range

**Symptoms:**
- Connection lost at short distances (<500m)
- High packet loss
- Weak signal (RSSI < -120)

**Solutions:**

1. **Antenna Issue (Most Common)**
   - Ensure antenna connected tightly
   - Try better antenna (3dBi → 5dBi)
   - Position antenna vertically (not horizontal)

2. **Obstructions**
   - Remove metal barriers
   - Elevate devices (higher = better)
   - Test in open field for baseline

3. **Interference**
   - Move away from WiFi routers
   - Check for other LoRa devices on same frequency
   - Try different channel

4. **Transmit Power**
   ```
   Settings → LoRa → Transmit Power
   Increase from 20dBm → 30dBm (max, legal)
   ⚠️ Higher power = shorter battery life
   ```

5. **Modem Preset**
   ```
   Settings → LoRa → Modem Preset
   Try "Long Moderate" or "Long Slow"
   Slower = longer range but higher latency
   ```

### Issue: GPS Not Working

**Symptoms:**
- No GPS fix after 10+ minutes
- Position shows 0,0 coordinates
- Map shows device at null island (0°N, 0°E)

**Solutions:**

1. **Outdoor Location Required**
   - GPS needs clear sky view
   - Won't work indoors (unless near large window)
   - Move to open area

2. **Wait Longer**
   - First fix can take 5-15 minutes
   - Subsequent fixes much faster (30-60s)
   - Cold start = slow, warm start = fast

3. **Check GPS Module**
   - Verify GPS module installed (RAK1910)
   - Check module seated properly
   - LED blinks when GPS locked

4. **Configure GPS**
   ```
   Settings → Position → GPS Enabled: ON
   Settings → Position → GPS Update: 60s
   ```

5. **Check Antenna** (External GPS)
   - Some devices need external GPS antenna
   - Check ceramic antenna placement

### Issue: High Battery Drain

**Symptoms:**
- Battery dies in <24 hours
- Expected multiple days

**Solutions:**

1. **Disable GPS**
   - GPS is biggest power drain
   - If not needed, turn off
   - Savings: 50-70% battery life

2. **Increase Intervals**
   ```
   Settings → Position → Update: 60s → 300s
   Settings → Telemetry → Update: 300s → 900s
   ```

3. **Reduce TX Power**
   ```
   Settings → LoRa → TX Power: 30dBm → 20dBm
   Still works if range is overkill
   Savings: 20-30% battery
   ```

4. **Turn Off Screen** (if applicable)
   ```
   Settings → Device → Screen Timeout: 60s → 10s
   Screen draws ~20mA continuously
   ```

5. **Check Battery Health**
   - Old batteries have reduced capacity
   - Use quality 18650 (Samsung, LG, Panasonic)
   - Avoid cheap batteries

6. **Add Solar**
   - 5-10W solar panel + charge controller
   - Enables perpetual operation
   - Cost: $20-40 per node

### Issue: Messages Not Delivered

**Symptoms:**
- Messages sent but not received
- Partial delivery
- Long delays (>30s)

**Checks:**

1. **Mesh Connectivity**
   - Verify devices in node list
   - Check hop count (>3 hops = unreliable)
   - Ensure intermediate nodes are powered on

2. **Channel Utilization**
   ```
   Check in app: Node Details → Channel Utilization
   >20% = congested mesh (too many messages)
   Solution: Reduce message frequency
   ```

3. **Packet Collision**
   - Too many devices transmitting at once
   - Random backoff handles this automatically
   - May cause delays but not loss

4. **Out of Range**
   - Last-hop node is too far
   - Need relay node between

5. **Rate Limiting**
   - Meshtastic limits message rate (anti-spam)
   - Wait 5-10s between messages
   - Don't flood mesh

---

## Next Steps

### Phase 1: Property Mapping (Week 1)

After basic testing works:

1. **Map Coverage**
   - Walk property with mobile device
   - Note dead zones (no signal)
   - Identify relay locations

2. **Plan Deployment**
   - Where to place fixed nodes
   - How many needed for full coverage
   - Solar power needed?

3. **Order Additional Hardware**
   - Based on coverage map
   - Add 1-2 extra for redundancy

### Phase 2: OneMind OS Integration (Week 2)

See `MESHTASTIC-MESH-INTEGRATION.md` for full details:

1. **Set Up MQTT Bridge**
   - Configure device for MQTT mode
   - Deploy bridge container on edge device
   - Connect to NATS

2. **Test Bidirectional Messaging**
   - Send command via NATS → Mesh
   - Receive telemetry via Mesh → NATS
   - Verify latency acceptable

3. **Implement Failover Logic**
   - Add mesh backup to drone/robot code
   - Test automatic failover
   - Validate works during internet outage

### Phase 3: Fleet Integration (Week 3-4)

1. **Add to Mobile Assets**
   - Mount LoRa modules on drones
   - Mount on robots
   - Add to vehicles

2. **Deploy Sensor Network**
   - Install field sensors with LoRa
   - Configure solar power
   - Weatherproof enclosures

3. **Full System Test**
   - Simulate internet outage
   - Verify fleet operates via mesh
   - Test emergency scenarios

---

## Useful Resources

### Official Documentation
- [Meshtastic Docs](https://meshtastic.org/docs) - Complete reference
- [Hardware Devices](https://meshtastic.org/docs/hardware/devices/) - All supported devices
- [Configuration](https://meshtastic.org/docs/configuration/) - Settings guide
- [Firmware](https://meshtastic.org/docs/software/firmware/) - Firmware releases

### Tools
- [Web Flasher](https://flasher.meshtastic.org) - Browser-based firmware flasher
- [Python CLI](https://github.com/meshtastic/python) - Command-line interface
- [iOS App](https://apps.apple.com/app/meshtastic/id1586432531)
- [Android App](https://play.google.com/store/apps/details?id=com.geeksville.mesh)

### Community
- [Meshtastic Discord](https://discord.gg/meshtastic) - Active community
- [Forum](https://meshtastic.discourse.group) - Q&A and discussions
- [Reddit r/meshtastic](https://reddit.com/r/meshtastic) - News and projects
- [GitHub](https://github.com/meshtastic) - Source code

### Shopping
- [RAK Wireless](https://store.rakwireless.com) - RAK devices
- [LILYGO AliExpress](https://lilygo.aliexpress.com) - LILYGO devices
- [Heltec](https://heltec.org) - Heltec devices
- [Amazon](https://amazon.com/s?k=meshtastic) - Various devices

---

## Quick Reference: Settings Checklist

**Before Devices Can Communicate:**

- [ ] **Region**: Same on both (US, EU, etc.)
- [ ] **Channel Name**: Identical (case-sensitive)
- [ ] **PSK**: Same encryption key
- [ ] **Modem Preset**: Same (Long Fast recommended)
- [ ] **Antenna**: Attached to both devices
- [ ] **Battery**: Charged
- [ ] **Firmware**: Same version on both

**Optimal Settings for OneMind OS:**

```yaml
Region: US (or your region)
Channel:
  Name: "onemind-private"
  PSK: [generated 256-bit key]
  Uplink: true
  Downlink: true

LoRa:
  Modem Preset: Long Fast
  TX Power: 30 (max)
  Hop Limit: 3

Position:
  GPS Enabled: true
  Update Interval: 60s
  Broadcast: true

Device:
  Role: Router (for fixed nodes) or Client (for mobile)
  Serial: Enabled (for bridge)

MQTT:
  Enabled: true (for hub nodes only)
  Server: 192.168.1.1
  Topic: mesh/#
```

---

## Cost Summary

### Minimal Starter (Testing)
- 2× RAK Starter Kit: $140
- 4× 18650 Batteries: $25
- Battery Charger: $15
- **Total: $180**

### Recommended Starter (Full Testing)
- 3× RAK Starter Kit: $210
- 6× 18650 Batteries: $40
- Battery Charger: $20
- 2× External Antennas: $20
- 2× IP65 Enclosures: $30
- **Total: $320**

### Full Property Deployment (10 nodes)
- 3× RAK (edge hubs): $210
- 3× LILYGO (WiFi bridges): $120
- 10× Sensor nodes (ESP32+LoRa): $400
- 10× Solar panels: $300
- 10× Enclosures: $150
- Batteries/accessories: $200
- **Total: $1,380**

---

**Document**: `MESHTASTIC-STARTER-GUIDE.md`
**Author**: OneMind OS Team
**Date**: 2026-01-22
**Version**: 1.0
**Related Docs**:
- `MESHTASTIC-MESH-INTEGRATION.md` - Complete architecture
- `MESHTASTIC-USE-CASES.md` - Real-world scenarios

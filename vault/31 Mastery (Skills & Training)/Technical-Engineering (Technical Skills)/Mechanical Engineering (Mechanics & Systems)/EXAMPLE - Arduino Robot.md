---
id: hp-28-20260113225955
title: "EXAMPLE   Arduino Robot"
type: reference
path: HP
domain: 25-49
folder: "28"
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [path/hp, path/hp/28, type/reference, status/active]
ai_summary: "EXAMPLE   Arduino Robot"
---
# Technical Project Log

**Project Name:** Autonomous Line-Following Robot
**Project ID:** MECH-ROBOT-001
**Started:** 2025-11-15
**Target Completion:** 2025-12-30
**Status:** 🟢 Active (80% complete)

---

## Project Overview

**Discipline:** Mechanical Engineering / Robotics
**Category:** Design / Build / Prototype

**Project Type:**
- [x] Personal learning project
- [x] Portfolio piece
- [ ] Client/Commission work
- [ ] Certification requirement
- [ ] Competition entry
- [ ] Open source contribution
- [ ] Research project

**Complexity:** Intermediate

---

## Project Description

**Purpose:**
Build an autonomous robot that can navigate a course by following a black line on a white surface. This project combines mechanical design, electronics, and programming to create a functional mobile robot platform.

**Objectives:**
1. Design and build stable 2-wheel drive chassis with ball caster
2. Implement line detection using IR sensors
3. Program PID control algorithm for smooth line following
4. Achieve consistent course completion in under 45 seconds

**Success Criteria:**
- [x] Robot successfully follows line without losing track
- [x] Handles 90-degree turns smoothly
- [ ] Completes test course in < 45 seconds (currently 52 seconds)
- [x] Battery life > 30 minutes continuous operation
- [ ] Document full build process for portfolio

**Constraints:**
- Budget: $150 ($128 spent, $22 remaining)
- Timeline: 6 weeks (currently week 5)
- Resources: Arduino starter kit, basic tools, 3D printer access

---

## Technical Specifications

### Requirements
**Functional Requirements:**
1. Follow black line (25mm wide) on white surface
2. Navigate 90-degree turns
3. Detect and handle line intersections
4. Autonomous operation (no remote control)

**Non-Functional Requirements:**
- Performance: Complete test course in < 45 seconds
- Safety: Low voltage (5V), emergency stop capability
- Standards: N/A (personal project)
- Environmental: Indoor use, smooth flat surface

### Design Parameters
**Key Specifications:**
| Parameter | Specification | Unit |
|-----------|--------------|------|
| Chassis dimensions | 150 × 120 × 80 | mm |
| Weight | 285 | grams |
| Max speed | 0.4 | m/s |
| Sensor range | 5-10 | mm |
| Battery capacity | 2000 | mAh |
| Operating voltage | 5 | V |

---

## Tools & Materials

### Hardware/Equipment
- [x] Soldering iron
- [x] Multimeter
- [x] Wire cutters/strippers
- [x] Small screwdriver set
- [x] 3D printer (makerspace access)
- [x] Hot glue gun

**Access:** Combination of owned tools and makerspace equipment

### Software/Applications
- [x] Arduino IDE 2.2.1
- [x] Fusion 360 (CAD design)
- [x] Cura (3D print slicing)
- [x] Git (version control)
- [ ] MATLAB (for PID tuning analysis - optional)

### Materials/Components
| Item | Quantity | Specs | Source | Cost |
|------|----------|-------|--------|------|
| Arduino Uno R3 | 1 | ATmega328P | Amazon | $25 |
| DC Motors | 2 | 6V, 200RPM | Amazon | $12 |
| Motor Driver | 1 | L298N dual H-bridge | Amazon | $8 |
| IR Sensors | 5 | TCRT5000 | Amazon | $15 |
| Ball Caster | 1 | 12mm metal | Amazon | $6 |
| Wheels | 2 | 65mm diameter | Amazon | $8 |
| Battery Pack | 1 | 7.4V 2000mAh LiPo | HobbyKing | $18 |
| Voltage Regulator | 1 | LM2596 buck converter | Amazon | $7 |
| Breadboard | 1 | Mini breadboard | Arduino kit | $0 |
| Jumper Wires | 40 | M-M, M-F | Arduino kit | $0 |
| PLA Filament | ~50g | Black 1.75mm | Makerspace | $2 |
| Mounting Hardware | set | M3 screws/nuts | Hardware store | $5 |
| Switch | 1 | SPST toggle | Arduino kit | $0 |
| LED Indicators | 3 | Red, Green, Blue | Arduino kit | $0 |
| Resistors | 10 | 220Ω, 10kΩ | Arduino kit | $0 |

**Total Material Cost:** $128
**Budget Remaining:** $22

---

## Project Phases

### Phase 1: Research & Planning ✅
- [x] Problem definition
- [x] Research existing solutions (YouTube, Instructables)
- [x] Requirements gathering
- [x] Feasibility study
- [x] Initial concept sketches

**Completed:** 2025-11-17
**Time Spent:** 6 hours

**Key Findings:**
- Most successful designs use 5 IR sensors in an array for precise line detection
- PID control is essential for smooth operation
- Ball caster provides better stability than skid plate
- Need voltage regulator to protect Arduino from motor noise

---

### Phase 2: Design & Engineering ✅
- [x] Conceptual design
- [x] Preliminary calculations (motor torque, battery life)
- [x] CAD modeling in Fusion 360
- [x] 3D print test fit
- [x] Design review (posted to r/arduino, got good feedback)
- [x] Final design documentation

**Status:** 100% complete
**Time Spent:** 12 hours
**Estimated Remaining:** 0 hours

**Design Decisions:**
- Chose 2-wheel drive with rear ball caster (simpler than 4WD, better turning)
- Sensor array positioned 15mm ahead of wheel axis for faster response
- Modular design - top plate lifts off for easy access to electronics
- Battery mounted low and centered for stability

**Challenges:**
- Initial design had sensor bar too far forward - caused overshoot on turns. Revised to 15mm.
- Had to redesign motor mounts 3 times to get proper gear mesh

---

### Phase 3: Prototyping & Testing ✅
- [x] Parts procurement
- [x] 3D printing chassis components
- [x] Electronics breadboard testing
- [x] Motor driver configuration
- [x] Sensor calibration
- [x] Initial assembly
- [x] Basic line following test
- [x] Performance testing on test course

**Status:** 100% complete
**Time Spent:** 18 hours

**Key Learnings:**
- Sensor calibration is critical - did 10-point calibration for each sensor
- Motor driver gets warm - added heatsink
- Wire management is important - loose wires caused intermittent issues

---

### Phase 4: Refinement & Final Build 🔄
- [x] Incorporate test feedback
- [x] Improved wiring organization
- [x] PID tuning (first pass complete, second iteration in progress)
- [x] Add status LEDs
- [ ] Final speed optimization
- [ ] Comprehensive testing (multiple course layouts)
- [ ] Final adjustments

**Status:** 80% complete
**Time Spent:** 15 hours
**Estimated Remaining:** 4 hours

**Current Focus:**
PID tuning to improve turn handling and reduce course time from 52s to target 45s.

---

### Phase 5: Documentation & Delivery ⏳
- [x] Commented code
- [ ] Technical documentation
- [ ] Build guide (in progress - 60% complete)
- [ ] Project video/demo
- [ ] Project report
- [ ] Portfolio page
- [ ] GitHub repository

**Status:** 30% complete
**Estimated Time:** 8 hours

---

## Work Sessions

### Session 12: 2025-12-18
**Duration:** 2.5 hours
**Phase:** Refinement
**Focus:** PID tuning for better turn performance

**Tasks Completed:**
- [x] Implemented PID parameter adjustment via serial monitor
- [x] Tested 15 different PID configurations
- [x] Logged performance data for each configuration
- [ ] Final configuration selected (narrowed to 3 candidates)

**Progress Made:**
Significant improvement in turn handling. Previous configuration (Kp=2.5, Ki=0, Kd=1.0) was oscillating too much. New best candidate (Kp=1.8, Ki=0.1, Kd=1.5) shows smoother transitions with less overshoot.

**Technical Notes:**
- Increased Kd (derivative) term helped dampen oscillations
- Small Ki (integral) term helps recover from persistent errors
- Too much Kp causes aggressive hunting behavior
- Sweet spot seems to be Kp:1.5-2.0, Ki:0.05-0.15, Kd:1.2-1.5

**Code/Formulas:**
```cpp
// PID Control Implementation
float error = 0;
float P = 0, I = 0, D = 0, PID_value = 0;
float previous_error = 0;

// Current best parameters
float Kp = 1.8;
float Ki = 0.1;
float Kd = 1.5;

void calculate_pid() {
  P = error;
  I = I + error;
  D = error - previous_error;

  PID_value = (Kp * P) + (Ki * I) + (Kd * D);
  previous_error = error;
}

void motor_control() {
  int left_motor_speed = initial_motor_speed - PID_value;
  int right_motor_speed = initial_motor_speed + PID_value;

  // Constrain speeds to valid PWM range
  left_motor_speed = constrain(left_motor_speed, 0, 255);
  right_motor_speed = constrain(right_motor_speed, 0, 255);

  // Apply to motors
  analogWrite(left_motor, left_motor_speed);
  analogWrite(right_motor, right_motor_speed);
}
```

**Problems Encountered:**
Robot occasionally loses line on sharp turns when speed is > 200 PWM (out of 255).

**Solutions Applied:**
Reduced base speed from 220 to 180 PWM. Slight performance decrease (53s → 52s) but much more reliable. Will revisit after PID optimization is complete.

**Next Session Plan:**
- Test final 3 PID configurations on full course (5 runs each)
- Select best configuration based on time + reliability
- Document final PID values

**Files Modified:**
- line_follower.ino (PID parameters, serial tuning menu)
- README.md (test results table)

---

### Session 13: 2025-12-20
**Duration:** 3 hours
**Phase:** Refinement / Documentation
**Focus:** Final PID testing and build documentation

**Tasks Completed:**
- [x] Tested 3 final PID configurations (5 runs each = 15 total runs)
- [x] Selected optimal configuration: Kp=1.8, Ki=0.1, Kd=1.4
- [x] Wrote assembly instructions (sections 1-3 complete)
- [x] Took photos of each assembly step

**Progress Made:**
Finalized PID parameters! Configuration #2 (Kp=1.8, Ki=0.1, Kd=1.4) gave best results:
- Average time: 49.2 seconds (target: <45s)
- Success rate: 5/5 (100%)
- Very smooth through turns, minimal oscillation on straights

Still 4.2s away from target, but much improved from initial 65+ seconds.

**Technical Notes:**
Test results summary:
| Config | Kp | Ki | Kd | Avg Time | Success | Notes |
|--------|----|----|----|----|---------|-------|
| 1 | 1.8 | 0.05 | 1.5 | 51.4s | 4/5 | One line loss on turn 3 |
| 2 | 1.8 | 0.1 | 1.4 | 49.2s | 5/5 | ✓ Best performance |
| 3 | 2.0 | 0.1 | 1.3 | 50.1s | 5/5 | Slightly more hunting |

Config #2 selected for final implementation.

**Problems Encountered:**
None! Clean testing session.

**Solutions Applied:**
N/A

**Next Session Plan:**
- Increase base speed from 180 to 200 PWM now that PID is optimized
- Test for time improvement without sacrificing reliability
- Continue build documentation (sections 4-6)
- Plan demo video shoot

**Files Modified:**
- line_follower.ino (PID values finalized)
- BUILD_GUIDE.md (added assembly sections 1-3)
- test_results.xlsx (all run data logged)

---

## Calculations & Analysis

### Key Calculations

#### Calculation 1: Motor Torque Required
**Purpose:**
Verify selected motors have sufficient torque to move the robot

**Formula:**
```
Torque_required = (m × a × r) + (m × g × μ × r)
Where:
  m = robot mass
  a = desired acceleration
  r = wheel radius
  g = gravitational acceleration
  μ = coefficient of friction
```

**Variables:**
- Mass (m) = 0.285 kg
- Acceleration (a) = 0.5 m/s² (desired)
- Wheel radius (r) = 0.0325 m (65mm diameter)
- g = 9.81 m/s²
- μ = 0.3 (estimated for rubber on smooth surface)

**Result:** 0.0320 N⋅m per motor

**Verification:**
Motors rated at 0.098 N⋅m @ 6V. Have 3x safety factor. ✓ Sufficient.

---

#### Calculation 2: Battery Life Estimation
**Purpose:**
Ensure robot can complete multiple test runs on single charge

**Formula:**
```
Runtime = Battery_Capacity / Total_Current_Draw
```

**Variables:**
- Battery capacity = 2000 mAh @ 7.4V
- Arduino current = 50 mA
- Motors (2×) = 200 mA each @ load
- Sensors (5×) = 5 mA each
- Total = 50 + 400 + 25 = 475 mA

**Result:** 4.2 hours theoretical runtime

**Verification:**
Actual testing showed 35 minutes of continuous operation before voltage drops below 6.5V (cutoff point). This is adequate for testing sessions (typically 15-20 minutes). Conservative estimate was accurate accounting for motor surges.

---

## Testing & Validation

### Test 1: Line Following Accuracy
**Date:** 2025-12-10
**Purpose:** Verify robot can follow straight line without deviation
**Method:** 5-meter straight line test, measure deviation from center

**Setup:**
- 5-meter straight black line (25mm wide) on white poster board
- Robot starts centered on line
- Measure max deviation from center at end

**Measurements:**
| Run | Max Deviation | Pass/Fail |
|-----|---------------|-----------|
| 1 | 8mm left | ✅ |
| 2 | 5mm right | ✅ |
| 3 | 12mm left | ✅ |
| 4 | 6mm left | ✅ |
| 5 | 9mm right | ✅ |

**Results:**
Average deviation: 8mm (well within 25mm line width). All runs successful.

**Analysis:**
Robot tracks line accurately. Small deviations are expected and acceptable. PID control keeping robot centered effectively.

**Action Items:**
None - test passed.

---

### Test 2: 90-Degree Turn Navigation
**Date:** 2025-12-14
**Purpose:** Verify robot handles sharp turns without losing line
**Method:** 4-corner test course, various turn configurations

**Setup:**
Test course with 4 × 90-degree turns (2 left, 2 right). Course dimensions: 2m × 1.5m rectangle.

**Measurements:**
| Run | Turn 1 (L) | Turn 2 (R) | Turn 3 (R) | Turn 4 (L) | Pass/Fail |
|-----|------------|------------|------------|------------|-----------|
| 1 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 2 | ✅ | ✅ | ❌ (lost line) | N/A | ❌ |
| 3 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 4 | ✅ | ✅ | ✅ | ❌ (lost line) | ❌ |
| 5 | ✅ | ✅ | ✅ | ✅ | ✅ |

**Results:**
Success rate: 3/5 (60%). Two failures both on sharp right turns. Issue identified: base speed too high (220 PWM) for sharp turns.

**Analysis:**
Robot handles gradual turns well but struggles with sharp 90-degree turns at high speed. The centrifugal force combined with sensor positioning causes line loss on sharp rights specifically.

**Action Items:**
- Reduced base speed to 180 PWM (completed)
- Re-test after PID optimization (scheduled)

---

### Test 3: Course Completion Time
**Date:** 2025-12-20
**Purpose:** Measure performance against 45-second target
**Method:** Full test course, time from start to finish

**Setup:**
Complete test course: 2m × 1.5m rectangle with 4 × 90-degree turns. Electronic timer.

**Measurements:**
| Run | Time | Success | Notes |
|-----|------|---------|-------|
| 1 | 49.8s | ✅ | Smooth, no issues |
| 2 | 48.2s | ✅ | Slight hesitation turn 3 |
| 3 | 50.1s | ✅ | Clean run |
| 4 | 48.9s | ✅ | Very smooth |
| 5 | 49.2s | ✅ | Perfect run |

**Results:**
Average time: 49.2 seconds. Target: 45 seconds. Currently 4.2 seconds over target.

**Analysis:**
Excellent reliability (100% success rate with final PID values), but speed needs improvement to hit target. Options:
1. Increase base motor speed (current: 180 PWM → could try 200 PWM)
2. Optimize turn initiation (earlier detection → earlier slowdown → faster overall)
3. Reduce sensor polling delay (current: 10ms → could try 5ms)

**Action Items:**
- Test increased motor speed (200 PWM) with current PID values
- If stable, push to 220 PWM
- Benchmark each change

---

## Issues & Troubleshooting

### Issue 1: Motor Noise Causing Arduino Resets
**Date Identified:** 2025-11-22
**Severity:** High
**Status:** Resolved

**Description:**
Arduino would randomly reset during operation, especially during motor direction changes. Serial monitor would disconnect and reconnect.

**Root Cause:**
Motor driver (L298N) creates electrical noise on power lines when switching. This noise was coupling into Arduino's power supply, causing voltage drops below brown-out detection threshold (reset trigger).

**Solution Attempted:**
1. Added 100μF capacitor across motor terminals → Partial improvement
2. Added 0.1μF ceramic cap across Arduino Vin → Helped but still occasional resets
3. Implemented separate power regulation (LM2596) for Arduino → Problem solved!

**Result:**
Using dedicated buck converter (LM2596) to provide clean 5V to Arduino, separate from motor power. No resets in 3 hours of continuous testing.

**Resolution:**
RESOLVED. Added LM2596 buck converter to design. Motors powered directly from 7.4V battery, Arduino powered from regulated 5V output of buck converter. Also added 470μF capacitor on buck converter output for additional filtering.

**Cost Impact:** +$7 (buck converter module)

---

### Issue 2: Inconsistent Sensor Readings
**Date Identified:** 2025-12-01
**Severity:** Medium
**Status:** Resolved

**Description:**
IR sensors giving inconsistent readings, causing erratic behavior. Robot would "see" lines that weren't there, or miss lines that were present.

**Root Cause:**
1. Ambient light interference (room lighting varying throughout day)
2. Sensors not calibrated to specific line/surface contrast
3. Some sensors mounted at slightly different heights (0.5-1mm variance)

**Solution Attempted:**
1. Added shields around sensors to block ambient light → Helped but not sufficient
2. Implemented auto-calibration routine on startup → Significant improvement!
3. Remounted sensor array with precise 7mm standoffs → Problem solved

**Result:**
Auto-calibration routine measures both line (black) and surface (white) during 3-second initialization. Calculates threshold as midpoint. Combined with consistent sensor height, readings are now rock solid.

**Resolution:**
RESOLVED. Calibration code added to setup():
```cpp
void calibrate_sensors() {
  // Sweep robot left-right during calibration
  // to expose sensors to both black and white
  for (int i = 0; i < 100; i++) {
    read_sensors();
    for (int j = 0; j < 5; j++) {
      if (sensor[j] < sensor_min[j]) sensor_min[j] = sensor[j];
      if (sensor[j] > sensor_max[j]) sensor_max[j] = sensor[j];
    }
    delay(20);
  }
  // Calculate thresholds
  for (int j = 0; j < 5; j++) {
    sensor_threshold[j] = (sensor_min[j] + sensor_max[j]) / 2;
  }
}
```

---

## Documentation

### Files & Assets
**CAD Files:**
- robot_chassis_v3.f3d - Fusion 360 source (current)
- chassis_bottom.stl - 3D print file
- chassis_top.stl - 3D print file
- sensor_mount.stl - 3D print file
- motor_bracket_L.stl - 3D print file
- motor_bracket_R.stl - 3D print file

**Code Repository:**
- Location: https://github.com/[username]/line-following-robot
- Branch: main
- Last commit: 3c8f9a2 ("Finalize PID values after testing")
- Lines of code: 287 (main sketch)

**Documentation:**
- [x] README.md (project overview, features, credits)
- [x] Bill of Materials (BOM.md)
- [ ] Technical specifications (SPECS.md) - 70% complete
- [ ] Assembly instructions (BUILD_GUIDE.md) - 60% complete
- [x] Code comments (fully commented)
- [x] Wiring diagram (Fritzing schematic)
- [x] Circuit schematic

**Photos/Videos:**
- [x] Progress photos (45 images)
- [x] Final product photos (12 images)
- [ ] Demo video (planned - not yet recorded)
- [x] Assembly timelapse (3min compiled video)

**Location:**
~/Projects/Arduino/LineFollowerRobot/
Google Drive: Projects/Robotics/Line Follower/

---

## Learning Outcomes

**New Skills Learned:**
- PID control implementation and tuning
- Fusion 360 CAD design for 3D printing
- Sensor calibration techniques
- Motor driver circuit design
- Power supply noise mitigation

**Concepts Mastered:**
- Proportional-Integral-Derivative control theory
- PWM (Pulse Width Modulation) for motor control
- H-bridge motor driver operation
- Voltage regulation and power distribution
- State machine programming for robot behaviors

**Tools/Software Learned:**
- Fusion 360 (parametric CAD)
- Arduino IDE (C/C++ programming)
- Serial plotter for real-time data visualization
- Multimeter techniques (current measurement, continuity testing)
- Git version control basics

**Mistakes & Lessons:**
1. **Measure twice, print once**: Had to reprint motor brackets 3 times because I didn't verify dimensions. Learned to always prototype critical dimensions with a quick test print first.

2. **Power supply is critical**: Spent days debugging "software" issues that turned out to be power supply noise. Lesson: Always use separate regulated supplies for logic and motors.

3. **Calibration is not optional**: Initial sensor readings were all over the place. Learned that proper calibration is essential for any sensing application.

4. **PID tuning is iterative**: Thought I could calculate perfect PID values mathematically. Reality: Requires systematic testing and adjustment. No substitute for real-world testing.

**Resources That Helped:**
- "PID Controllers Explained" by Brian Douglas (YouTube series)
- Arduino Forums (especially troubleshooting power issues)
- r/arduino community (design review feedback)
- Fusion 360 tutorials by Paul McWhorter
- "Programming Arduino" by Simon Monk (book)

---

## Project Metrics

**Time Tracking:**
| Phase | Estimated | Actual | Variance |
|-------|-----------|--------|----------|
| Research | 4 hrs | 6 hrs | +2 hrs |
| Design | 8 hrs | 12 hrs | +4 hrs |
| Prototyping | 12 hrs | 18 hrs | +6 hrs |
| Refinement | 8 hrs | 15 hrs | +7 hrs |
| Documentation | 8 hrs | 5 hrs (ongoing) | -3 hrs |
| **TOTAL** | **40 hrs** | **56 hrs (ongoing)** | **+16 hrs** |

**Budget Tracking:**
| Category | Budgeted | Spent | Remaining |
|----------|----------|-------|-----------|
| Electronics | $80 | $90 | -$10 |
| Mechanical | $30 | $21 | +$9 |
| Materials | $20 | $12 | +$8 |
| Tools | $20 | $5 | +$15 |
| **TOTAL** | **$150** | **$128** | **$22** |

Note: Over budget on electronics (added voltage regulator + extra sensors) but under on mechanical (used makerspace filament). Overall under budget.

---

## Completion

**Completion Date:** 2025-12-30 (projected)
**Total Time:** ~65 hours (estimated final)
**Total Cost:** $135 (estimated final)

**Final Performance:**
| Metric | Target | Achieved | Status |
|--------|--------|----------|--------|
| Line following | Yes | Yes | ✅ |
| 90° turns | Yes | Yes | ✅ |
| Course time | < 45s | 49.2s | 🟡 |
| Battery life | > 30min | 35min | ✅ |
| Reliability | > 90% | 100% | ✅ |

**Project Success:** Partially met (4/5 goals achieved)

**Deliverables:**
- [x] Working prototype (fully functional)
- [ ] Technical documentation (in progress)
- [x] Source files (CAD + code on GitHub)
- [x] Test results (documented)
- [ ] Build guide (60% complete)
- [ ] Presentation/Demo (video planned)

---

## Post-Project

**Demonstration:**
- Presented to: Family, friends
- Feedback: Very positive. Suggestions for future: add obstacle avoidance, increase speed

**Publication:**
- [x] Added to portfolio (personal website)
- [x] GitHub repository (public)
- [ ] Blog post/Article (planned for portfolio site)
- [ ] Video tutorial (planned)
- [ ] Instructables (considering)
- [ ] Reddit r/arduino (will post when fully documented)

**Lessons for Portfolio:**
This project demonstrates:
- Mechanical design (CAD, 3D printing)
- Electronics (circuit design, sensor integration)
- Programming (C/C++, algorithms, PID control)
- Problem-solving (debugging power issues, sensor calibration)
- Documentation (code comments, build guide)

**Future Improvements:**
1. Implement speed profiling (slow for turns, fast for straights) → could hit 45s target
2. Add ultrasonic sensor for obstacle detection
3. Add Bluetooth module for remote monitoring/tuning
4. Design PCB to replace breadboard (more compact, reliable)
5. Implement intersection handling for maze-solving capability

**Follow-up Projects:**
1. Maze-solving robot (adds decision making at intersections)
2. Bluetooth-controlled robot (uses same chassis)
3. PCB design project (custom controller board)

---

**Series:** Arduino Robotics Projects
**Related Projects:**
- [[Next: Obstacle Avoidance Robot]]
- [[Future: Sumo Robot for Competition]]

---

**Related:** [[Technical Skills Hub]] | [[Weekly Progress Review]] | [[Skill Goal Tracker]]

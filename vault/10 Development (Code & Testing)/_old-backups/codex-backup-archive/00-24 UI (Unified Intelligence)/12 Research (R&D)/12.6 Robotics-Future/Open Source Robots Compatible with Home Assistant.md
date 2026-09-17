---
id: ui-15-20260111053926
title: "Open Source Robots Compatible with Home Assistant"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Open Source Robots Compatible with Home Assistant"
---

# Open Source Robots Compatible with Home Assistant

Yes, there are several open source robot platforms that you can integrate with Home Assistant to send commands to and control remotely. Based on my research, here are the main options available:

## Commercial Robot Vacuum Integrations

**Robot Vacuums with Home Assistant Support:**
- **Roborock** - Popular vacuum robots with official Home Assistant integration[1][2]
- **iRobot Roomba/Braava** - Official integration available for Wi-Fi enabled models[3][4]
- **Ecovacs Deebot** - Main integration for vacuum and mower control[5]
- **Shark IQ** - Direct integration with room targeting capabilities[6]
- **SwitchBot Robot Vacuums** - K20+ Pro and other models with "Works with Home Assistant" certification[7][8]

## Open Source Robot Platforms

### **ROS-Based Mobile Robots**

**TurtleBot Series** - One of the most popular open source robot platforms[9][10][11]
- Low-cost, personal robot kit with open-source software
- Built on iRobot Create 3 mobile base
- Full ROS 2 compatibility
- Can be programmed for house navigation, 3D vision, and custom applications
- Price range: $1,199-$1,499 for complete kits

**ROSbot Platforms**[12]
- ROSbot 2 PRO - Autonomous platform running on UpBoard (x64) with 4GB RAM
- ROSbot XL - Universal ROS 2-native platform for R&D and prototyping
- Full ROS compatibility for custom development

**LIMO Mobile Robot**[13]
- Open-source mobile robot with ROS 1 & 2 compatibility
- Four different movement modes: tracked, omnidirectional, differential, Ackermann
- Built-in LiDAR, camera, and inertial units
- Designed for AI development and autonomous navigation

**Leo Rover**[14]
- Ready-to-go open-source 4×4 mobile robot
- Weatherproof design with built-in UI and video streaming
- ROS-based platform
- Suitable for outdoor applications

### **DIY Robot Platforms**

**Open Source Quadrupeds**[15][16]
- **Petoi Bittle** - Arduino/ESP32-based robot dog ($269)
- **Stanford Pupper** - Low-cost quadruped with 12 degrees of freedom (~$600)
- **SpotMicro** - 3D printed four-legged robot inspired by Boston Dynamics

**Humanoid Robots**[16][17]
- **Poppy Humanoid** - 3D printed humanoid with servo motors (~$9,000)
- **InMoov** - Open-source 3D printed humanoid
- **Salvius** - Humanoid robot made from recycled parts

**Mobile Rovers**[18][19]
- **JPL Open Source Rover** - NASA's scaled-down Mars rover design
- **Linorobot** - Suite of low-cost ROS-compatible ground robots
- **DiffBot** - Autonomous 2WD differential drive robot using ROS

## Integration Methods with Home Assistant

### **Direct Integration via Custom Components**
Some developers have created custom Home Assistant components for robot control:
- **ROS Robot Component** - GitHub project for supporting robots using ROS[20]
- **MQTT-based robot control** - Using MQTT protocol for communication[21][22]

### **Communication Protocols**

**MQTT Integration**[23][24]
- Lightweight messaging protocol perfect for robot communication
- Can create custom sensors, switches, and controls
- Real-time bidirectional communication
- Example: Send navigation commands, receive status updates

**ESPHome Integration**[25][26]
- For ESP32/ESP8266-based robots
- YAML configuration for custom robot control
- Direct integration with Home Assistant
- Perfect for DIY servo-controlled robots and cameras

**ROS Bridge Integration**[22][27]
- Connect ROS robots to Home Assistant via MQTT bridge
- Allows sending ROS topics from Home Assistant
- Enable robot navigation and task execution from smart home interface

## Building Your Own Integration

### **DIY Approach Options:**

1. **MQTT-Based Robot**[23]
   - Use ESP32/Arduino for robot control
   - Publish robot status to MQTT topics
   - Subscribe to command topics from Home Assistant
   - Create custom dashboard controls

2. **ESPHome Robot**[26][28]
   - Configure robot components via YAML
   - Servo motors, sensors, cameras
   - Direct Home Assistant integration
   - OTA updates and debugging

3. **ROS Bridge**[22][27]
   - Set up ROS master on dedicated computer
   - Use MQTT bridge between ROS and Home Assistant
   - Send navigation goals and receive robot feedback
   - Advanced robotics capabilities

### **Example Use Cases:**
- **Security patrols** - Send robot to specific rooms on command
- **Delivery tasks** - Navigate to designated locations
- **Environmental monitoring** - Collect sensor data throughout house
- **Pet feeding** - Automated or manual feeding routines[28]
- **Cleaning schedules** - Custom cleaning patterns and timing

## Getting Started Recommendations

For beginners, I'd recommend:

1. **Start with a robot vacuum integration** (Roborock, Roomba) to learn Home Assistant robot control basics
2. **Try TurtleBot** if you want a professional ROS-based platform
3. **Build an ESP32-based robot** with ESPHome for custom DIY projects
4. **Use MQTT** as the communication bridge for maximum flexibility

The key is that Home Assistant's open architecture allows integration with virtually any robot platform through MQTT, custom components, or API bridges. The robot doesn't need official Home Assistant support - you can create the integration yourself using these communication protocols.

Sources
[1] Roborock - Home Assistant https://www.home-assistant.io/integrations/roborock/
[2] Easy Roborock Setup on Home Assistant - YouTube https://www.youtube.com/watch?v=rWKuWzsi6Ew
[3] iRobot Roomba and Braava - Home Assistant https://www.home-assistant.io/integrations/roomba/
[4] Integrating Roomba with Home Assistant - Tanner's Tech https://tcude.net/integrating-roomba-with-home-assistant/
[5] Ecovacs - Home Assistant https://www.home-assistant.io/integrations/ecovacs/
[6] Shark IQ - Home Assistant https://www.home-assistant.io/integrations/sharkiq/
[7] Home Assistant | SwitchBot US https://us.switch-bot.com/pages/home-assistant
[8] SwitchBot Works with Home Assistant – Full Integration Guide https://www.youtube.com/watch?v=If1lvrVJMSU
[9] TurtleBot 4 | Unboxing & Getting Started - YouTube https://www.youtube.com/watch?v=QN01AXjoLdQ
[10] Basic Setup · User Manual - GitHub Pages https://turtlebot.github.io/turtlebot4-user-manual/setup/basic.html
[11] TurtleBot https://www.turtlebot.com
[12] robots.ros.org https://robots.ros.org
[13] LIMO Open-Source Mobile Robot (ROS compatible) https://www.generationrobots.com/en/404091-robot-mobile-open-source-limo-compatible-ros1-et-ros2-limo-standard-version.html
[14] Leo Rover | Robot Developer Kit | Open-source :ROS and for ... https://www.leorover.tech
[15] 20 Open-Source Robotics Projects | Ilir Aliu | 21 comments - LinkedIn https://www.linkedin.com/posts/ilir-aliu_20-open-source-robotics-projects-on-activity-7187428102891847680-yME9
[16] Best Open-Source Robots for Robotics Research - Petoi https://www.petoi.com/pages/best-open-source-robots-for-robotics-research
[17] stephane-caron/awesome-open-source-robots - GitHub https://github.com/stephane-caron/awesome-open-source-robots
[18] mjyc/awesome-robotics-projects - GitHub https://github.com/mjyc/awesome-robotics-projects
[19] 10 Open Source Robotics Projects on GitHub - Akshet Patel https://akshetpatel.substack.com/p/10-open-source-robotics-projects
[20] Home Assistant component supporting robots using ROS. - GitHub https://github.com/KamilaLis/rosrobot
[21] send command to machine (scout) through home assistant? - Reddit https://www.reddit.com/r/ROS/comments/1jgjiz3/send_command_to_machine_scout_through_home/
[22] [PDF] Integrating robotics to a smart home environment - Theseus https://www.theseus.fi/bitstream/10024/748869/2/Gratscheff_Samuli.pdf
[23] MQTT 101: Integrate your DIY Devices into Home Assistant - YouTube https://www.youtube.com/watch?v=Q0S0xOW35k8
[24] MQTT - Home Assistant https://www.home-assistant.io/integrations/mqtt/
[25] ESPHome - Smart Home Made Simple — ESPHome https://esphome.io/index.html
[26] janfajessen/EspHome-RC-FPV-Car-Robot-Esp8266-pantilt ... - GitHub https://github.com/janfajessen/EspHome-RC-FPV-Car-Robot-Esp8266-pantilt-Servo-motors-Camera-USB-led-L298N---Home-Assistant
[27] OASIS, a smart home operating system based on ROS 2 https://community.home-assistant.io/t/oasis-smart-home-platform-based-on-ros-2-robot-operating-system/405110
[28] r/Esphome on Reddit: In this tutorial, I show everything you need ... https://www.reddit.com/r/Esphome/comments/1dnsp5e/in_this_tutorial_i_show_everything_you_need_and/
[29] Do Robot Vacuum Work With Home Assistant? (How To Set) - Ecovacs https://www.ecovacs.com/us/blog/do-robot-vacuums-work-with-home-assistant
[30] humanoid robot NEO gamma does chores as a home assistant https://www.designboom.com/technology/humanoid-robot-neo-gamma-house-chores-home-owners-autonomous-assistant-02-26-2025/
[31] Home Assistant - GitHub https://github.com/home-assistant
[32] Vacuum - Home Assistant https://www.home-assistant.io/integrations/vacuum/
[33] Is there a best of list for compatible devices that work with Home ... https://www.reddit.com/r/homeassistant/comments/14dgv6c/is_there_a_best_of_list_for_compatible_devices/
[34] Home Assistant: Building the open source smart home - YouTube https://www.youtube.com/watch?v=S71LpeMUrYg
[35] Top 10 Home Assistant Integrations for 2025! - YouTube https://www.youtube.com/watch?v=AS9N8uKvpn8
[36] Why I use Home Assistant for open source home automation https://opensource.com/article/20/11/home-assistant
[37] New brands of robotic vacuum cleaners - Home Assistant Community https://community.home-assistant.io/t/new-brands-of-robotic-vacuum-cleaners/878885
[38] Home Assistant Green https://www.home-assistant.io/green/
[39] Snippy: a Robotic Home Assistant… But With Scissors - Instructables https://www.instructables.com/Snippy-a-Robotic-Home-Assistant-But-With-Scissors/
[40] Home Assistant https://www.home-assistant.io
[41] Litter-Robot - Home Assistant https://www.home-assistant.io/integrations/litterrobot/
[42] Robotics Kits - Robotics - SparkFun Electronics https://www.sparkfun.com/robotics/robotics-kits.html
[43] Assistant Delivery Robot for Nursing Home using ROS https://dl.acm.org/doi/fullHtml/10.1145/3531028.3531051
[44] Open Robotics https://www.openrobotics.org
[45] Are any of these DIY robotics kits good for a total beginner to get an ... https://www.reddit.com/r/robotics/comments/9xlday/are_any_of_these_diy_robotics_kits_good_for_a/
[46] Building a DIY Robot Kit - TECHING Robot Be your friend - YouTube https://www.youtube.com/watch?v=pSLYB20DMcE&vl=en
[47] Open-Source Robotics Hardware Projects - Yale University https://www.eng.yale.edu/grablab/openrobothardware.html
[48] Robot Kits and Robotics for Kids in STEM Education - Makeblock https://www.makeblock.com/collections/robot-kits
[49] Robotics & CNC / Robot Kits & Chassis - Adafruit https://www.adafruit.com/category/34
[50] Home Assistant - Wikipedia https://en.wikipedia.org/wiki/Home_Assistant
[51] Best open source robot projects : r/robotics - Reddit https://www.reddit.com/r/robotics/comments/1ahublz/best_open_source_robot_projects/
[52] Robot Kits - EZ-Robot https://www.ez-robot.com/robot-kits.html
[53] Poppy Project - Open source robotic platform https://www.poppy-project.org
[54] Robots & Kits - RobotShop https://www.robotshop.com/collections/robot-kits
[55] Autonomous Robots - SuperDroid Robots https://www.superdroidrobots.com/Robots/autonomous-robots
[56] NVIDIA Isaac - AI Robot Development Platform https://developer.nvidia.com/isaac
[57] Autonomous Mobile Platforms | KUKA AG https://www.kuka.com/en-us/products/amr-autonomous-mobile-robotics/mobile-platforms
[58] TurtleBot 4 Common · User Manual - GitHub Pages https://turtlebot.github.io/turtlebot4-user-manual/software/turtlebot4_common.html
[59] Automated & Autonomous Mobile Robots | Robotnik® https://robotnik.eu/products/mobile-robots/
[60] Integrations - Home Assistant https://www.home-assistant.io/integrations/
[61] Husky UGV - ROS Open Source Mobile Robot https://clearpathrobotics.com/husky-ugv-ros-open-source-mobility/
[62] Autonomous robots - cobots, mobile robots, humanoid robots https://www.generationrobots.com/en/537-autonomous-robots
[63] Turtlebot at Office: A Service-Oriented Software Architecture for ... https://ieeexplore.ieee.org/document/7781988
[64] ROS Mobile Robots - GitHub https://github.com/ros-mobile-robots
[65] Outdoor mobile robots |Generation Robots https://www.generationrobots.com/en/352-outdoor-mobile-robots
[66] Learned About Home Assistant From ESPHome Via WLED https://newscrewdriver.com/2022/02/17/learned-about-home-assistant-from-esphome-via-wled/
[67] ROBOTONT – Open-source and ROS-supported omnidirectional ... https://www.sciencedirect.com/science/article/pii/S2468067223000433
[68] Open-Source-Robotics Projects - GitHub https://github.com/SmilingRobo/open-source-robotics
[69] ROS and Arduino with ROSserial | ROS tutorial #10 - YouTube https://www.youtube.com/watch?v=WLVfZXxpHYI
[70] msadowski/awesome-weekly-robotics: A list of projects that ... - GitHub https://github.com/msadowski/awesome-weekly-robotics
[71] Autonomous Home Robot to Help Around the House! https://projecthub.arduino.cc/nbaddorf/autonomous-home-robot-to-help-around-the-house-554466
[72] Ultimate Home Assistant Beginner's Guide! - YouTube https://www.youtube.com/watch?v=Z4gvkmJ8q48
[73] jcarolinares/awesome-social-robots - GitHub https://github.com/jcarolinares/awesome-social-robots
[74] Autonomous Home Assistant Robot - Hackster.io https://www.hackster.io/grassjelly/autonomous-home-assistant-robot-fff43e
[75] How I built my smart home - 2023 edition : r/homeassistant - Reddit https://www.reddit.com/r/homeassistant/comments/198w7ge/how_i_built_my_smart_home_2023_edition/
[76] Building a Home Assistant Custom Component Part 1 - Aaron Godfrey https://aarongodfrey.dev/home%20automation/building_a_home_assistant_custom_component_part_1/
[77] My first Home Assistant custom component https://seanb.co.uk/2017/07/my-first-home-assistant-custom-component/
[78] Building a Home Assistant Custom Component Part 5: Debugging https://aarongodfrey.dev/home%20automation/building_a_home_assistant_custom_component_part_5/
[79] MQTT Vacuum - Home Assistant https://www.home-assistant.io/integrations/vacuum.mqtt/
[80] Custom Components Throttle functionality etc - Development https://community.home-assistant.io/t/custom-components-throttle-functionality-etc/287781
[81] Home Assistant - Valetudo https://valetudo.cloud/pages/integrations/home-assistant-integration.html
[82] Using Servo Motors in ESPHome and Home Assistant - YouTube https://www.youtube.com/watch?v=1f2tuS_PxIE
[83] Tutorial for creating a custom component - Home Assistant Community https://community.home-assistant.io/t/tutorial-for-creating-a-custom-component/204793
[84] Setting Up MQTT for Home Assistant in under 5 Minutes in 2025! https://www.youtube.com/watch?v=ImBv2mtaRR0
[85] Servo Component — ESPHome https://esphome.io/components/servo.html
[86] custom-components/remote_homeassistant: Links multiple ... - GitHub https://github.com/custom-components/remote_homeassistant
[87] MQTT Sensor - Home Assistant https://www.home-assistant.io/integrations/sensor.mqtt/
[88] How to setup ESPHome to control my Bluetooth controlled ... https://community.home-assistant.io/t/how-to-setup-esphome-to-control-my-bluetooth-controlled-octocontrol-bed/540790
[89] Creating a Home Assistant component - Coding Stephan https://svrooij.io/2023/01/18/home-assistant-component/
[90] Home Assistant 101: How to Integrate MQTT Devices - YouTube https://www.youtube.com/watch?v=WYyPMHqMV-w

---
id: ui-15-20260111053926
title: "Open-Source Robot Companies Similar to Unitree for Home Assistant Integration"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Open-Source Robot Companies Similar to Unitree for Home Assistant Integration"
---

Now I have comprehensive information about open-source robot companies and platforms that are similar to Unitree and can be integrated with Home Assistant. Let me compile this information into a comprehensive report.

# Open-Source Robot Companies Similar to Unitree for Home Assistant Integration

Based on my research, I've found several excellent open-source robot platforms similar to Unitree that offer comprehensive documentation, open-source models, and integration possibilities with Home Assistant. Here are the top options:

## **Unitree Robotics - Confirmed Open Source Support**

**Models Available**: G1, H1, Go2, B2
- **Open Source Components**: Full SDK, ROS integration, datasets, and AI models
- **Documentation**: [GitHub repositories](https://github.com/unitreerobotics) with extensive examples
- **Home Assistant Integration**: Possible via ROS2 bridge using MQTT
- **Key Open Source Projects**:
  - [unitree_sdk2](https://github.com/unitreerobotics/unitree_sdk2) - Main SDK for Go2, B2, H1, G1
  - [unitree_ros2](https://github.com/unitreerobotics/unitree_ros2) - ROS2 integration
  - [unitree_IL_lerobot](https://github.com/unitreerobotics/unitree_IL_lerobot) - AI training framework
  - Open datasets on [Hugging Face](https://huggingface.co/UnitreeRobotics)

## **K-Scale Labs - Fully Open Source Humanoid**

**Model**: K-Bot (Starting at $8,999)
- **License**: Complete open-source hardware and software (BSD 3-clause)
- **Documentation**: [GitHub](https://github.com/kscalelabs) with full CAD files, electronics, PCB designs
- **Home Assistant Integration**: ROS-compatible with Python/Rust SDK
- **Special Features**: 
  - 3D printable components
  - Modular design for easy customization
  - GPU-accelerated robot learning framework
  - Operating system with Python package installation (`pip install kos`)

## **Petoi OpenCat - Arduino/Raspberry Pi Based**

**Models**: Bittle (dog), Nybble (cat)
- **License**: MIT License - fully open source
- **Platform**: Arduino and Raspberry Pi compatible
- **Documentation**: [Extensive documentation](https://docs.petoi.com/) with block-based programming
- **Home Assistant Integration**: Easy via Arduino/Pi GPIO and MQTT
- **Programming Options**:
  - Block-based coding (Scratch-like)
  - C/C++/Python
  - Petoi Coding Blocks for beginners
- **Cost**: Under $500 for complete kits

## **Poppy Project - Educational Humanoid Platform**

**Models**: Poppy Humanoid, Poppy Torso
- **License**: Creative Commons BY-SA (hardware), GPL v3 (software)
- **Documentation**: [Comprehensive guides](https://www.poppy-project.org/) and [GitHub](https://github.com/poppy-project)
- **Home Assistant Integration**: Python SDK with ROS compatibility
- **Cost**: $8,000-9,000 for complete humanoid
- **Special Features**:
  - 3D printed components
  - Educational focus
  - Modular design
  - Active community support

## **Solo 8 - NYU/Max Planck Open Source Quadruped**

**Developer**: NYU Tandon & Max Planck Institute
- **License**: BSD 3-clause license
- **Documentation**: [Open Dynamic Robot Initiative](https://open-dynamic-robot-initiative.github.io)
- **Home Assistant Integration**: ROS-based architecture
- **Cost**: Under $3,000 for complete build
- **Special Features**:
  - Entirely 3D printed and off-the-shelf components
  - Torque-controlled joints
  - Research-grade platform
  - Lightweight (2.2kg)

## **Xiaomi CyberDog - Open Source APIs**

**Model**: CyberDog 2
- **Open Source**: 100+ APIs, ROS2 integration
- **Documentation**: [GitHub repositories](https://github.com/MAVProxyUser/Xiaomi-CyberDog-Open-Source)
- **Home Assistant Integration**: ROS2 with DDS communication
- **Special Features**:
  - NVIDIA Jetson compute platform
  - Visual programming interface
  - Over 100 open-source APIs
  - AI self-learning system

## **ANYbotics ANYmal - Research Platform**

**Models**: ANYmal C, ANYmal D
- **Open Architecture**: Full ROS API access
- **Documentation**: [Research community](https://www.anymal-research.org)
- **Home Assistant Integration**: ROS-based with extensive APIs
- **Special Features**:
  - Industrial-grade hardware
  - Open software stack
  - Research community support
  - Full simulation environment

## **Home Assistant Integration Methods**

### **ROS2 Bridge Integration**
The most robust method for integrating these robots with Home Assistant involves:

1. **MQTT Bridge**: Use ROS2 to MQTT bridge for communication
2. **Custom Component**: Create Home Assistant custom component
3. **Node-RED**: Use as middleware for complex automations

### **Required Components for Integration**:
```yaml
# Example Home Assistant configuration
mqtt:
  broker: localhost
  port: 1883

sensor:
  - platform: mqtt
    name: "Robot Battery"
    state_topic: "robot/battery"
    unit_of_measurement: "%"

switch:
  - platform: mqtt
    name: "Robot Power"
    command_topic: "robot/power/set"
    state_topic: "robot/power/state"
```

### **Programming Requirements**:
- **ROS2 Foxy/Humble** for modern robots
- **Python/C++** for custom development  
- **MQTT** for Home Assistant communication
- **Docker** for containerized deployment

## **Recommended Development Path**

### **For Beginners**:
1. **Petoi Bittle** - Most accessible with block programming
2. **Poppy Torso** - Good educational platform
3. **K-Scale Labs** - Modern, well-documented

### **For Advanced Users**:
1. **Unitree Go2** - Professional grade with full SDK
2. **ANYmal Research** - Industry-leading platform
3. **Solo 8** - Research-focused, highly customizable

### **For Home Assistant Integration**:
1. **K-Scale Labs** - Best modern integration support
2. **Unitree** - Comprehensive ROS2 support
3. **Petoi** - Simplest hardware integration

## **Key Documentation Resources**

- **Unitree**: [GitHub Organization](https://github.com/unitreerobotics) with 37+ repositories
- **K-Scale Labs**: [Development site](https://www.kscale.dev) with full stack documentation
- **Petoi**: [Documentation center](https://docs.petoi.com/) with tutorials
- **Poppy**: [Project website](https://www.poppy-project.org) with community support
- **Solo 8**: [Research papers](https://open-dynamic-robot-initiative.github.io) and build guides

All these platforms provide the open-source foundation you need for Home Assistant integration, with varying levels of complexity and cost. The choice depends on your specific requirements, budget, and technical expertise level.

Sources
[1] 宇树科技—全球四足机器人行业开创者 https://www.unitree.com
[2] Unitree G1 Open Source Dataset - YouTube https://www.youtube.com/watch?v=OTWHXTu09wE
[3] PetoiCamp/OpenCat-Quadruped-Robot - GitHub https://github.com/PetoiCamp/OpenCat-Quadruped-Robot
[4] K-Scale Labs: Open-source humanoid robots - Y Combinator https://www.ycombinator.com/companies/k-scale-labs
[5] unitreerobotics/unitree_IL_lerobot - GitHub https://github.com/unitreerobotics/unitree_IL_lerobot
[6] OpenCat the Open Source Quadruped Robotic Pet Framework - Petoi https://www.petoi.com/pages/opencat-open-source-robot-pet-framework
[7] Your Personal Open-Source Humanoid Robot for $8999 - YouTube https://www.youtube.com/watch?v=BS92RdBvI90
[8] Unitree Robotics - GitHub https://github.com/unitreerobotics
[9] Meet CyberDog, Xiaomi's open-source quadruped robot https://www.therobotreport.com/meet-cyberdog-a-new-open-source-quadruped-robot-from-xiaomi/
[10] $8000 Open-Source HUMANOID Robot STUNS OpenAI & NVIDIA! https://www.youtube.com/watch?v=Im-ShBXlNeU
[11] In order to promote the development of the global embodied AI ... https://www.reddit.com/r/singularity/comments/1gpknsp/in_order_to_promote_the_development_of_the_global/
[12] Open-source, low-cost, quadruped robot makes sophisticated ... https://engineering.nyu.edu/news/open-source-low-cost-quadruped-robot-makes-sophisticated-robotics-available-all
[13] Open-Source-Robotics Projects - GitHub https://github.com/SmilingRobo/open-source-robotics
[14] Humanoid agent AI avatar - Unitree Robotics https://www.unitree.com/g1
[15] Any decent open-source quadruped 'Spot' robots still actively being ... https://www.reddit.com/r/robotics/comments/16cb168/any_decent_opensource_quadruped_spot_robots_still/
[16] K-Scale Labs: Open-source humanoid robots, built for developers https://www.kscale.dev
[17] Unitree Humanoid Robots Updated With Open-Source Full-Body ... https://www.youtube.com/watch?v=SdZoEke0GPk
[18] An Open Source Quadruped-Robot for Simulating and Real ... - GitHub https://github.com/TopHillRobotics/quadruped-robot
[19] Open Source Humanoid Robots That You Can 3D Print Yourself https://www.forbes.com/sites/johnkoetsier/2025/04/14/open-source-humanoid-robots-hugging-face-buys-pollen-robotics/
[20] Official Open Source - Unitree Robotics https://www.unitree.com/opensource
[21] The Dingo | A Low Cost, Open-Source Robot Quadruped - YouTube https://www.youtube.com/watch?v=8KntOIgzUjY
[22] Boston Dynamics Robot Dog Now Open Source Code | Hypebeast https://hypebeast.com/2020/1/boston-dynamics-spot-robot-open-source-code
[23] Agility Robotics shows off latest advances for Digit humanoid https://www.therobotreport.com/agility-robotics-announces-latest-advances-digit-humanoid-robot/
[24] Open-Source Software Focuses On Quadruped Robots - RTInsights https://www.rtinsights.com/open-source-software-focuses-on-quadruped-robots/
[25] Open-Source Robot Code Releases : Boston Dynamics - Trend Hunter https://www.trendhunter.com/trends/boston-dynamics1
[26] Agility Robotics - Wikipedia https://en.wikipedia.org/wiki/Agility_Robotics
[27] [PDF] Full Stack Software Framework for Agile Quadrupedal Locomotion https://www.andrew.cmu.edu/user/amj1/papers/Quad_SDK_ICRA_Abstract.pdf
[28] boston-dynamics/spot-sdk - GitHub https://github.com/boston-dynamics/spot-sdk
[29] GTLIDAR/digit_controller - GitHub https://github.com/GTLIDAR/digit_controller
[30] Open-source software gives a leg up to robot research - Tech Xplore https://techxplore.com/news/2022-08-open-source-software-leg-robot.html
[31] How to use the Boston Dynamics Spot SDK - YouTube https://www.youtube.com/watch?v=Hsa1wPpJA_Q
[32] The Future Of Robotics: Digit By Agility Robotics - YouTube https://www.youtube.com/watch?v=0dexTZ7r02Q
[33] Open-source software gives a leg up to robot research https://engineering.cmu.edu/news-events/news/2022/07/26-quadruped.html
[34] QuickStart — Spot 5.0.0 documentation - Spot SDK - Boston Dynamics https://dev.bostondynamics.com/docs/python/quickstart.html
[35] adubredu/DigitRobot.jl: URDF and Meshes for the Agility ... - GitHub https://github.com/adubredu/DigitRobot.jl
[36] Boston Dynamics Grows Spot Developer Toolkit https://bostondynamics.com/news/boston-dynamics-grows-spot-developer-toolkit/
[37] Agility Robotics https://www.agilityrobotics.com
[38] robomechanics/quad-sdk: Software tools for agile quadrupeds ... https://github.com/robomechanics/quad-sdk
[39] SpotMicro - Boston Dynamics Spot inspired robot - GitLab https://gitlab.com/public-open-source/spotmicroai
[40] Our Solution - Agility Robotics https://www.agilityrobotics.com/solution
[41] Poppy Humanoid - GitHub https://github.com/poppy-project/poppy-humanoid
[42] OpenCat/README.md at main - GitHub https://github.com/PetoiCamp/OpenCat/blob/main/README.md
[43] Solo 8 Makes Open Source, Sophisticated Robots Accessible and ... https://www.hackster.io/news/solo-8-makes-open-source-sophisticated-robots-accessible-and-affordable-9cad032f5fc5
[44] Poppy Humanoid https://www.poppy-project.org/en/robots/poppy-humanoid/
[45] Petoi Coding Blocks https://docs.petoi.com/block-based-programming/petoi-coding-blocks
[46] Four-legged robot makes research comparable worldwide https://is.mpg.de/news/four-legged-robot-makes-research-comparable-worldwide
[47] Poppy Humanoid Robot (without 3D printed parts) https://www.generationrobots.com/en/403348-robot-poppy-humanoid-version-raspberry-sans-impressions-3d.html
[48] Petoi Robot Dog Coding for Beginners | Block Programming with Dr ... https://www.youtube.com/watch?v=1G_D74MkR6c
[49] Small Robot Makes Great Strides—and Jumps | Machine Design https://www.machinedesign.com/markets/robotics/article/21135212/small-robot-makes-great-stridesand-jumps
[50] Poppy Project - Open source robotic platform https://www.poppy-project.org
[51] New user w/ Bittle programming questions - Petoi Camp https://www.petoi.camp/forum/general-discussions/new-user-w-bittle-programming-questions
[52] Four-Legged Robot Enhances Global Research Collaboration ... https://www.innovations-report.com/engineering/power-and-electrical-engineering/four-legged-robot-makes-research-comparable-worldwide/
[53] Poppy - Technological platform - Génération Robots https://www.generationrobots.com/en/279-poppy-opensource-robotics-platform
[54] Solo 8, an affordable robot dog? - 3Dnatives https://www.3dnatives.com/en/solo-8-robot-dog-060720206/
[55] POPPY Humanoid Robot | 3DEXPERIENCE Edu - Dassault Systèmes https://www.3ds.com/edu/education/academics/hub/poppy-humanoid-robot
[56] Petoi Bittle Tutorial: Step-by-Step Guide with Arduino IDE - YouTube https://www.youtube.com/watch?v=8ao5VlmRwe4
[57] Open Dynamic Robot Initiative https://open-dynamic-robot-initiative.github.io
[58] Poppy Humanoid by poppy_project - Thingiverse https://www.thingiverse.com/thing:2280067
[59] How to get the Poppy Humanoid sources files https://poppy.discourse.group/t/how-to-get-the-poppy-humanoid-sources-files/94
[60] Roborock - Home Assistant https://www.home-assistant.io/integrations/roborock/
[61] OASIS, a smart home operating system based on ROS 2 https://community.home-assistant.io/t/oasis-smart-home-platform-based-on-ros-2-robot-operating-system/405110
[62] Robot Vacuum Smart Automation : r/homeassistant - Reddit https://www.reddit.com/r/homeassistant/comments/1aqpj11/robot_vacuum_smart_automation/
[63] Ecovacs - Home Assistant https://www.home-assistant.io/integrations/ecovacs/
[64] Isn't ROS supposed to be an easier and more efficient way to build a ... https://www.reddit.com/r/ROS/comments/r9karl/isnt_ros_supposed_to_be_an_easier_and_more/
[65] Automate Your Roborock S7 MaxV Ultra EASY with Home Assistant https://www.youtube.com/watch?v=Hi2rgGz5cPQ
[66] iRobot Roomba and Braava - Home Assistant https://www.home-assistant.io/integrations/roomba/
[67] [PDF] Integrating robotics to a smart home environment - Theseus https://www.theseus.fi/bitstream/10024/748869/2/Gratscheff_Samuli.pdf
[68] SwitchBot joins Works with Home Assistant https://www.home-assistant.io/blog/2025/06/26/switchbot-joins-works-with-home-assistant/
[69] Vacuum - Home Assistant https://www.home-assistant.io/integrations/vacuum/
[70] Open Robotics https://www.openrobotics.org
[71] How To Integrate Home Assistant & Husqvarna Automower - YouTube https://www.youtube.com/watch?v=Hynl9IejI_Q
[72] Do Robot Vacuum Work With Home Assistant? (How To Set) - Ecovacs https://www.ecovacs.com/us/blog/do-robot-vacuums-work-with-home-assistant
[73] ROS: Home https://www.ros.org
[74] Integrating Roomba with Home Assistant - Tanner's Tech https://tcude.net/integrating-roomba-with-home-assistant/
[75] Litter-Robot - Home Assistant https://www.home-assistant.io/integrations/litterrobot/
[76] Autonomous Home Assistant Robot - Hackster.io https://www.hackster.io/grassjelly/autonomous-home-assistant-robot-fff43e
[77] A Slightly Smarter Vacuum Automation in Home Assistant - YouTube https://www.youtube.com/watch?v=6kJZAYHRBJ0
[78] Shark IQ - Home Assistant https://www.home-assistant.io/integrations/sharkiq/
[79] I built a micro-ROS based home robot - open source - Reddit https://www.reddit.com/r/microROS/comments/177jroy/i_built_a_microros_based_home_robot_open_source/
[80] MAVProxyUser/Xiaomi-CyberDog-Open-Source - GitHub https://github.com/MAVProxyUser/Xiaomi-CyberDog-Open-Source
[81] ROS package for the simulation of legged platforms using Gazebo https://github.com/KodlabPenn/kodlab_gazebo
[82] Open-Source Software Boosts Quadrapedal Robot - ASME https://www.asme.org/topics-resources/content/open-source-software-boosts-quadrapedal-robots
[83] Xiaomi CyberDog 2 Your New Robotic Pet - YouTube https://www.youtube.com/watch?v=cKks9YeX4k0
[84] Adding Ghost Robotics Robots · Issue #3 · chvmp/robots - GitHub https://github.com/chvmp/robots/issues/3
[85] Advance Robotics Research with ANYmal - ANYbotics https://www.anybotics.com/partners/research-academia/
[86] MAVProxyUser/XiaomiCyberDog-k91: 小米集团 XIaomi 铁 ... - GitHub https://github.com/MAVProxyUser/XiaomiCyberDog-k91
[87] Spirit 40 small Q-UGV for Swarm, ISR, Subterranean ... - YouTube https://www.youtube.com/watch?v=FrqUo4M3tY4
[88] ANYmal - Stanford Robotics Center https://src.stanford.edu/demo/anymal
[89] CyberDog - ROBOTS: Your Guide to the World of Robotics https://robotsguide.com/robots/cyberdog
[90] Spirit robot - Robomechanics Lab https://www.cmu.edu/me/robomechanicslab/robots.html
[91] ANYmal C Robot Description (URDF) - GitHub https://github.com/ANYbotics/anymal_c_simple_description
[92] Ghost Robotics: Spirit™ Series Utilizes Origin - Stratasys https://www.stratasys.com/en/resources/case-studies/ghost-robotics-origin/
[93] ANYmal Research https://www.anymal-research.org
[94] XIAOMI'S CYBERDOG 2 - A MORE LIFELIKE QUADRUPED ROBOT https://thearsenale.com/blogs/magazine/xiaomis-cyberdog-2-a-more-lifelike-quadruped-robot
[95] Spirit 40 Small Q-UGV Can Run at 6.5ft/sec - Robotic Gizmos https://www.roboticgizmos.com/spirit-40/
[96] ANYmal – Robotic Systems Lab | ETH Zurich https://rsl.ethz.ch/robots-media/anymal.html
[97] What use or application might cause robots like Xiaomi's new ... https://www.reddit.com/r/Futurology/comments/15sue5a/what_use_or_application_might_cause_robots_like/
[98] Careers - Ghost Robotics https://www.ghostrobotics.io/careers
[99] Mastering the Unitree SDK for the Go2 Quadruped - RoboStore https://robostore.com/blogs/news/mastering-the-unitree-sdk-for-the-go2-quadruped
[100] unitreerobotics/unitree_ros2 - GitHub https://github.com/unitreerobotics/unitree_ros2
[101] peaq integrates Home Assistant: the world's largest open-source ... https://www.peaq.xyz/blog/peaq-integrates-home-assistant-the-worlds-largest-open-source-smart-home-system
[102] Go1 — Unitree_Docs 1.0rc documentation - Unitree Support! https://unitree-docs.readthedocs.io/en/latest/get_started/Go1_Edu.html
[103] Compilation of Unitree GO1 Information, Simulation Environment ... https://egh0bww1.com/posts/2023-10-25-unitree-go1-collection-en/
[104] Introduction to SDK - unitree-z1-docs https://dev-z1.unitree.com/use/introduction.html
[105] ROS — Unitree Go1 Documentation https://docs.trossenrobotics.com/unitree_go1_docs/getting_started/ros.html
[106] Legged SDK — Unitree Go1 Documentation https://docs.trossenrobotics.com/unitree_go1_docs/getting_started/legged_sdk.html
[107] Robot Interface | Unitree H1 https://www.docs.quadruped.de/projects/h1/html/quick_start.html
[108] mikey0000/Mammotion-HA: Home Assistant integration for ... - GitHub https://github.com/mikey0000/Mammotion-HA
[109] H1 SDK Development Guide - 宇树科技文档中心 https://support.unitree.com/home/en/H1_developer
[110] unitreerobotics/unitree_ros - GitHub https://github.com/unitreerobotics/unitree_ros
[111] Go2 SDK Development Guide - 宇树科技 文档中心 https://support.unitree.com/home/en/developer/Obtain%20SDK
[112] ROS2 Services Interface - 宇树科技 文档中心 https://support.unitree.com/home/en/developer/ROS2_service
[113] Dreame L20 Ultra vacuum integration : r/homeassistant - Reddit https://www.reddit.com/r/homeassistant/comments/16xzk0w/dreame_l20_ultra_vacuum_integration/
[114] unitreerobotics/unitree_legged_sdk: SDK tools for control robots. https://github.com/unitreerobotics/unitree_legged_sdk
[115] Using the ROS Interface — Unitree Go1 Documentation https://docs.trossenrobotics.com/unitree_go1_docs/operation/ros.html

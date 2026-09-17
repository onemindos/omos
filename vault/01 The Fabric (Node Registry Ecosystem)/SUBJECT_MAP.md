# OneMind OS — POC Subject Map
# Every message in the system follows this table.
# Three domains: SaaS, Drone, Robot
# Same shape, same bus, same rules.

## SUBJECT GRAMMAR (locked)
## <plane>.<class>.<id>.<path>
##
## Planes:
##   cmd  = command (you are telling something to do something)
##   evt  = event   (something happened — durable, JetStream)
##   tel  = telemetry (continuous readings — ephemeral, no JetStream)
##   qry  = query (ask a question, get an answer back via reply)
##   gov  = governance (config, policy, registration)

# ─────────────────────────────────────────────────────────────────────────────
# SAAS — Personal & Business Applications
# ─────────────────────────────────────────────────────────────────────────────
#
# Every SaaS app is class: service
# Bridge: saas_bridge.py (HTTPS/API wrapper)
# Account: SAAS (isolated from CORE)
#
# Pattern:
#   Legacy or automation publishes cmd.service.<app>.<action>
#   SaaS bridge calls HTTPS API
#   Bridge publishes evt.service.<app>.<event> back to bus
#
# TODOIST (personal + business tasks)
cmd.service.todoist.create_task      → create a task
cmd.service.todoist.complete_task    → close a task
cmd.service.todoist.list_tasks       → fetch tasks (qry pattern)
evt.service.todoist.task_created     → task was made
evt.service.todoist.task_completed   → task was closed

# FASTMAIL (personal + business email)
cmd.service.fastmail.send_email      → send an email
cmd.service.fastmail.list_inbox      → list inbox (qry)
evt.service.fastmail.email_received  → new email arrived

# GITHUB (dev + business)
cmd.service.github.create_issue      → open an issue
cmd.service.github.create_pr         → open a PR
cmd.service.github.list_issues       → list issues (qry)
evt.service.github.issue_created     → issue opened
evt.service.github.pr_merged         → PR merged

# GHL (GoHighLevel — business CRM)
cmd.service.ghl.create_contact       → add a contact
cmd.service.ghl.send_sms             → send SMS via GHL
cmd.service.ghl.add_tag              → tag a contact
evt.service.ghl.contact_created      → contact was added
evt.service.ghl.form_submitted       → form came in

# SKOOL (community)
cmd.service.skool.create_post        → post to community
evt.service.skool.member_joined      → new member
evt.service.skool.post_created       → new post

# CALENDAR
cmd.service.calendar.create_event    → add calendar event
cmd.service.calendar.list_today      → get today's events (qry)
evt.service.calendar.event_created   → event was added
evt.service.calendar.event_starting  → event starts in 15min

# NOTION / LINEAR / STRIPE
cmd.service.notion.create_page       → create Notion page
cmd.service.linear.create_issue      → create Linear issue
cmd.service.stripe.list_payments     → list payments (qry)
evt.service.stripe.payment_received  → money came in


# ─────────────────────────────────────────────────────────────────────────────
# DRONE — Aerial vehicles
# ─────────────────────────────────────────────────────────────────────────────
#
# class: drone
# Bridge: mavlink_bridge.py (MAVLink ↔ NATS)
# Runs on: companion computer on the drone, or ground station Mac
# Inner loop: PX4/ArduPilot handles flight control at 400Hz — NATS never touches that
#
# Pattern:
#   Legacy publishes cmd.drone.<id>.<action>
#   MAVLink bridge translates to MAVLink command
#   Drone autopilot executes
#   Bridge publishes tel.drone.<id>.* and evt.drone.<id>.* back
#
# SKYHAWK-1 (drone ID from node registry)
cmd.drone.skyhawk-1.takeoff          → take off to altitude
cmd.drone.skyhawk-1.land             → land at current position
cmd.drone.skyhawk-1.return_home      → return to launch point
cmd.drone.skyhawk-1.goto             → fly to lat/lon/alt
cmd.drone.skyhawk-1.mission_start    → execute stored mission
cmd.drone.skyhawk-1.mission_pause    → pause current mission
cmd.drone.skyhawk-1.photo            → take a photo now

tel.drone.skyhawk-1.gps              → lat/lon/alt/heading (10Hz — ephemeral)
tel.drone.skyhawk-1.battery          → voltage/remaining % (1Hz)
tel.drone.skyhawk-1.heartbeat        → armed/mode/system_status (1Hz)

evt.drone.skyhawk-1.takeoff_initiated → drone started taking off
evt.drone.skyhawk-1.landing           → drone started landing
evt.drone.skyhawk-1.landed            → drone on ground
evt.drone.skyhawk-1.returning_home    → RTH triggered
evt.drone.skyhawk-1.low_battery       → battery below threshold
evt.drone.skyhawk-1.mission_completed → mission done


# ─────────────────────────────────────────────────────────────────────────────
# ROBOT — Ground robots (DIMOS, OM1, G1, Go2)
# ─────────────────────────────────────────────────────────────────────────────
#
# class: robot
# Bridge: ros2_bridge.py (ROS2/DDS ↔ NATS)
# Runs on: robot's onboard compute (Jetson, Raspberry Pi, onboard CPU)
# Inner loop: nav2/MoveIt/ROS2 handles path planning + motor control — NATS never touches that
#
# Pattern:
#   Legacy publishes cmd.robot.<id>.<action>
#   ROS2 bridge sends nav2 goals or action calls
#   Robot executes using native ROS2 stack
#   Bridge publishes tel.robot.<id>.* and evt.robot.<id>.* back
#
# OM1 (robot ID from node registry)
cmd.robot.om1.goto                   → navigate to x/y in map frame
cmd.robot.om1.stop                   → emergency stop all motion
cmd.robot.om1.patrol                 → run patrol route through named zone
cmd.robot.om1.inspect                → go to location and capture sensor data
cmd.robot.om1.dock                   → return to charging dock
cmd.robot.om1.status                 → query current pose + battery (qry)

tel.robot.om1.pose                   → x/y/z/heading (1Hz — ephemeral)
tel.robot.om1.battery                → voltage/remaining/charging (1Hz)
tel.robot.om1.lidar                  → obstacle distances (5Hz — ephemeral)

evt.robot.om1.goal_sent              → nav2 goal published
evt.robot.om1.goal_reached           → robot arrived at destination
evt.robot.om1.obstacle_detected      → obstacle blocked path
evt.robot.om1.patrol_started         → patrol sequence started
evt.robot.om1.patrol_completed       → patrol done
evt.robot.om1.low_battery            → needs charging
evt.robot.om1.docked                 → on charger

# G1 (Unitree G1 humanoid — same pattern, different ID)
cmd.robot.g1.goto
cmd.robot.g1.stop
cmd.robot.g1.wave
tel.robot.g1.pose
evt.robot.g1.task_completed


# ─────────────────────────────────────────────────────────────────────────────
# HOW LEGACY TRIGGERS ALL THREE
# ─────────────────────────────────────────────────────────────────────────────
#
# "Inspect the barn with OM1 and create a Todoist task when done"
#
#  Legacy publishes:
#    cmd.robot.om1.inspect  { "zone": "barn" }
#
#  ROS2 bridge hears it → robot goes to barn → robot inspects
#  Bridge publishes:
#    evt.robot.om1.patrol_completed  { "zone": "barn" }
#
#  Subscriber rule fires:
#    ON evt.robot.om1.patrol_completed
#    → publish cmd.service.todoist.create_task { "content": "Barn inspection done" }
#
#  SaaS bridge hears it → calls Todoist API → task created
#  Bridge publishes:
#    evt.service.todoist.task_created
#
#  All done. Zero direct app-to-app calls. Everything through the bus.
#
#
# "Launch the drone and take a photo of the south field"
#
#  Legacy publishes:
#    cmd.drone.skyhawk-1.takeoff     { "altitude_m": 30 }
#
#  MAVLink bridge arms drone, sends takeoff
#  Bridge publishes:
#    evt.drone.skyhawk-1.takeoff_initiated
#
#  Subscriber rule fires after altitude telemetry confirms airborne:
#    ON tel.drone.skyhawk-1.gps WHERE alt > 25
#    → publish cmd.drone.skyhawk-1.goto { "lat": 41.xxx, "lon": -73.xxx }
#    → publish cmd.drone.skyhawk-1.photo {}
#
#  Drone files to south field, takes photo.
#  Bridge publishes:
#    evt.drone.skyhawk-1.mission_completed

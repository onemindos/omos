# ROS-LM — Natural Language ROS2 Robot Control via NATS

**Repo:** github.com/oscarpoudel/LLM_ROS2_agent  
**Language:** Python  
**Role in OneMind:** Reference for LLM→NATS→ROS2 bridge pattern

Multimodal LLM agent that enables natural language control of ROS2 robots using NATS as the communication layer. Supports robot motion, waypoint navigation, vision queries, system diagnostics, and ROS2 CLI control.

## Architecture

```
User (natural language)
    → Intent Classifier
    → Agent Tool (motion / nav / vision / status / ros_cli)
    → NATS publish
    → ROS2 receiver (subscribe)
    → Physical action
```

## Tools

| Tool | NATS → ROS2 action |
|------|--------------------|
| `motion_tool.py` | Direct velocity commands |
| `nav_tool.py` | Nav2 waypoint navigation |
| `vision_tool.py` | VLM camera queries |
| `status_tool.py` | System diagnostics |
| `ros_cli_tool.py` | ROS2 CLI passthrough |

## Models

Ollama (local) or Gemini — swappable via `.env`

## Why it matters for OneMind

Concrete implementation of the NATS↔ROS2 bridge pattern. The tool-per-intent design maps directly onto OneMind's command subjects:

```
cmd.robot.{id}.motion.{action}   (motion_tool)
cmd.robot.{id}.nav.waypoint      (nav_tool)
cmd.robot.{id}.vision.query      (vision_tool)
```

## See Also

- [[dimos]] — full generalist robotics OS (more complete solution)
- [[nats-ros-connector]] — lower-level NATS↔ROS bridge

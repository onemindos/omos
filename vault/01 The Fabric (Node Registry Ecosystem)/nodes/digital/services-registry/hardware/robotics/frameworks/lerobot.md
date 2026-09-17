# LeRobot

> **Hugging Face's open-source robotics framework for real-world AI**

## Overview

LeRobot is Hugging Face's open-source framework for training and deploying AI models on real robots. Provides pre-trained models, datasets, and simulation environments for physical AI agents. Compatible with low-cost robot arms (SO-100, MOSS) and full humanoids.

## Status

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Status      | 🔄 Evaluating                      |
| Category    | Robotics                           |
| Type        | Robotics ML Framework              |
| Integration | Python, HuggingFace Hub, ROS2      |

## Use Cases in OneMind

1. **Physical Agents** — Training AI policies for SO/OM1 robotics work
2. **Imitation Learning** — Record human demonstrations to train robot policies
3. **Pre-trained Models** — Load community-trained robot policies from HuggingFace Hub
4. **Simulation** — Test in MuJoCo/IsaacGym before real hardware

## Key Features

- Pre-trained policies (ACT, Diffusion Policy)
- Teleoperation for data collection
- HuggingFace Hub integration for sharing models
- Support for SO-100, ALOHA, and custom robots
- ROS2 bridge
- Active HuggingFace community

## Quick Start

```bash
pip install lerobot
# Clone and visualize a dataset
python lerobot/scripts/visualize_dataset.py \
  --repo-id lerobot/pusht
```

## Links

- GitHub: [github.com/huggingface/lerobot](https://github.com/huggingface/lerobot)
- HuggingFace: [huggingface.co/lerobot](https://huggingface.co/lerobot)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

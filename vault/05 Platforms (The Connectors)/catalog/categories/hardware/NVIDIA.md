# NVIDIA

> **GPU computing platform for AI, robotics, and high-performance compute**

## Overview

NVIDIA provides the GPU hardware and software stack powering most modern AI workloads. Relevant for local AI inference (Ollama + CUDA), robotics edge computing (Jetson), cloud GPU services, and the Isaac robotics simulation platform.

## Status

| Attribute   | Value                             |
| ----------- | --------------------------------- |
| Status      | 🔄 Evaluating                     |
| Category    | Hardware / AI Compute             |
| Type        | GPU Platform                      |
| Integration | CUDA, TensorRT, Jetson, NGC       |

## Use Cases in OneMind

1. **Local Inference** — CUDA-accelerated Ollama for fast local LLM runs
2. **Robotics Edge** — Jetson Orin (SO-100, OM1) for on-device robot AI
3. **Cloud GPU** — DigitalOcean GPU Droplets or Lambda Labs for training
4. **Isaac Sim** — Robot simulation before deploying to real hardware
5. **TensorRT** — Optimize inference speed for deployed models

## Key Products

| Product       | Use Case                           |
| ------------- | ---------------------------------- |
| GeForce (RTX) | Consumer GPU, local inference      |
| Jetson Orin   | Edge AI for robotics               |
| A100/H100     | Data center training               |
| DGX           | On-prem AI supercomputer           |
| NGC           | Pre-built AI containers and models |

## Software Stack

- CUDA — parallel computing platform
- TensorRT — inference optimizer
- cuDNN — deep learning primitives
- Isaac Sim — robot simulation
- Triton — inference server

## Links

- Developer: [developer.nvidia.com](https://developer.nvidia.com)
- Jetson: [developer.nvidia.com/embedded-computing](https://developer.nvidia.com/embedded-computing)
- NGC: [ngc.nvidia.com](https://ngc.nvidia.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

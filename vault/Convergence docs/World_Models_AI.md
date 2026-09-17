# World Models in AI: Building Internal Representations of Environments

**What It Is:** A machine learning system that learns to build and maintain an internal representation of an environment, then predict how that environment changes over time in response to actions.

**Core Purpose:** Enable AI agents to plan, reason, simulate, and act intelligently WITHOUT constant real-world trial and error. The model "understands" the world so it can imagine futures before acting.

---

## Definition & Core Concept

### Simple Explanation
Imagine a chess player who doesn't need to physically move pieces to decide a strategy. Instead, they mentally simulate future board positions. A **world model** does this for AI:

1. **Learn the environment** — Observe videos, sensor data, actions, outcomes
2. **Build internal model** — Represent objects, physics, causality, rules
3. **Predict futures** — "If I do action X, what happens?"
4. **Plan intelligently** — Choose actions that lead to desired outcomes

### Technical Definition
A **world model** is a learned, differentiable representation of environmental dynamics that:
- **Observes**: Takes sensory input (images, video, state vectors)
- **Represents**: Encodes environment structure (object positions, velocities, relationships)
- **Predicts**: Forecasts future states given actions
- **Enables planning**: Agents use predictions to solve problems without trial-and-error

---

## Historical Context

### Timeline

**1990 — Jürgen Schmidhuber introduces the concept**
- Proposed using **recurrent neural networks (RNNs)** to learn models of environment dynamics
- Vision: Agents should learn to predict, not just react
- Foundation paper: Work on artificial curiosity and intrinsic motivation

**1990s-2000s — Early exploration**
- RNNs, LSTM networks tested for sequence prediction
- Limited success due to computational constraints
- Video prediction remains unsolved problem

**2010s — Deep learning renaissance**
- Convolutional Neural Networks (CNNs) + RNNs enable larger models
- Video prediction becomes tractable
- Companies start applying world models to robotics

**2015-2018 — Modern world models emerge**
- Google/DeepMind: Model-based RL (PlaNet, Dreamer)
- OpenAI: Video prediction models
- Robotics labs: World models for control

**2019 — "World Models" research landmark (Ha & Schmidhuber)**
- Showed agents could learn in latent space representation of environments
- Demonstrated unsupervised learning of world understanding
- Car driving simulator example became canonical

**2020-2024 — Large-scale generative models**
- Video diffusion models (pixel-space world models)
- Transformer-based world models
- Sora, Gato, and successors demonstrate scaling

**2024-2026 — Foundation models as world simulators**
- Large generative models (GPT-4V, Gemini, Claude) can implicitly reason about world
- Video generation models (Sora) act as implicit world models
- New focus: World models for planning in latent space of foundation models

---

## How World Models Work: Architecture

### Three Key Components

#### 1. **Encoder** — Compress observations into latent space
```
Raw input (e.g., video frame)
  ↓ [Convolutional neural network]
  ↓ Compress high-dimensional input
  → Latent vector (compact representation)

Example: 
  Input: 256×256×3 RGB image = 196,608 values
  Output: 128-dim latent vector (1,000x compression)
```

**Purpose:**
- Reduce dimensionality (easier to learn in latent space)
- Remove irrelevant visual details (noise, lighting)
- Learn abstract features (positions, velocities, object types)

**Architecture Examples:**
- Variational Autoencoder (VAE)
- Convolutional encoder
- Vision Transformer encoder

#### 2. **Dynamics Model** — Predict future states
```
Current state (latent vector) + action
  ↓ [RNN or Transformer]
  ↓ Learn dynamics: "How does this state evolve?"
  → Next state (latent vector)

Example:
  Current: [ball_x=100, ball_y=50, velocity_x=10, velocity_y=0]
  Action: [move_paddle_right]
  Output: [ball_x=110, ball_y=50, velocity_x=10, velocity_y=0]
```

**Purpose:**
- Model physics (gravity, collisions)
- Learn causality (actions cause state changes)
- Enable forward simulation

**Architecture Examples:**
- LSTM (Long Short-Term Memory) networks
- Transformer-based model
- Graph neural networks (for multi-agent dynamics)

#### 3. **Decoder** — Reconstruct observations from latent space
```
Latent vector (prediction)
  ↓ [Deconvolutional neural network]
  ↓ Reconstruct high-dimensional output
  → Video frame (to verify prediction accuracy)

Optional: Not always needed (sometimes agent just plans in latent space)
```

**Purpose:**
- Verify model learned correct representation (reconstruction loss)
- Generate predicted future frames (for visualization/validation)
- Learn to ignore irrelevant information

**Architecture Examples:**
- Deconvolutional encoder
- Vision Transformer decoder
- Diffusion decoder

### The Full Training Pipeline

```
Raw video data (agent plays game for hours)
  ↓
[Encoder] Extract latent representations
  ↓
[Dynamics model] Learn to predict next latent state
  ↓
[Decoder] Reconstruct predicted frame
  ↓
Compare reconstructed frame to actual next frame
  ↓
Calculate loss (prediction error)
  ↓
Backpropagation: Update all weights
  ↓
After thousands of iterations: Model learns world physics/rules
```

---

## Core Capabilities of World Models

### 1. **Planning Without Trial-and-Error**
```
Agent wants to reach goal (e.g., catch ball in game)

Traditional (model-free RL):
  - Try action 1 → fail → try action 2 → fail → ...
  - Millions of actual interactions needed
  - Slow and expensive

World model approach:
  - Current state: [ball_position, agent_position, ...]
  - Imagine future: "If I move right, ball will be here in 3 steps"
  - Plan optimal path: Find action sequence that reaches goal
  - Execute planned sequence
  - Much faster (planning in simulation, not reality)
```

### 2. **Reasoning About Unseen Scenarios**
```
Model trained on soccer game

World model can reason:
  "If I haven't seen a soccer ball at this angle, I can extrapolate
   how it should move based on learned physics principles"

Enables generalization beyond training data
```

### 3. **Curiosity-Driven Exploration**
```
Agent uses world model to estimate "uncertainty":
  "What situations can I explore that would most improve my model?"

Implementation:
  - Try action → observe outcome
  - If outcome matches prediction: Low surprise (skip)
  - If outcome differs from prediction: High surprise (explore more)
  
Result: Agent naturally explores interesting/uncertain parts of environment
        without explicit reward signal
```

### 4. **Credit Assignment**
```
Agent takes 10 actions before reaching goal

Question: Which actions actually mattered?

World model helps:
  - Trace back causal chain
  - Identify key decision points
  - Learn which actions have long-term effects
```

### 5. **Simulated Experience**
```
Agent can "dream" using world model:
  - Generate imaginary rollouts (trajectories)
  - Learn skills in imagination before real-world testing
  - Combine real experience + imagined experience
  
Example: Robot learning to grasp objects
  - Train world model on real robot data
  - Imagine grasping scenarios in latent space
  - Learn control policy in simulation
  - Test on real robot (much faster than learning from scratch)
```

---

## Technical Approaches & Variants

### Approach 1: Latent Space Models (Most Common)
```
Observation → Encoder → Latent representation
                            ↓
                    Dynamics model (RNN/Transformer)
                            ↓
Decoder → Reconstructed observation ← Latent representation
```

**Advantages:**
- Efficient (work in compressed space, not raw pixels)
- Scalable (can learn longer horizons)
- Clean mathematical formulation

**Examples:**
- Ha & Schmidhuber's "World Models" (VAE + RNN)
- Google's PlaNet, Dreamer
- Transformer-based world models

### Approach 2: Pixel-Space Models (Brute Force)
```
Raw video frame → Predict next raw video frame directly

No encoder/decoder (work at pixel level)
```

**Advantages:**
- No information loss (encode everything)
- Simpler architecture

**Disadvantages:**
- Computationally expensive
- Hard to scale
- Learns irrelevant details (compression artifacts, shadows)

**Examples:**
- Early CNN-based video prediction
- Some diffusion-based models

### Approach 3: Symbolic/Discrete World Models
```
Represent environment as objects + relationships:
  Objects: [ball, paddle_1, paddle_2, ...]
  Properties: [position_x, position_y, velocity, ...]
  Relationships: [collides_with, near, ...]

Dynamics: Symbolic rules (not learned neural networks)
  Rule: "If ball collides with paddle → bounce"

Advantages: Interpretable, generalizes well
Disadvantages: Requires hand-coded structure
```

### Approach 4: Foundation Models as Implicit World Models
```
Large language models (LLMs) + vision:
  - Claude, GPT-4V, Gemini implicitly understand world dynamics
  - Can reason about "if I do X, Y happens"
  - Trained on human knowledge (physics, cause-effect, etc.)

Video generation models (Sora):
  - Trained to predict future video frames
  - Acts as world model without explicit modeling
  - Learned from internet-scale video data

Advantages: Zero-shot reasoning about novel domains
Disadvantages: Black-box, expensive inference
```

---

## Real-World Applications

### 1. **Robotics: Learning to Grasp**
```
Challenge: Robot arm must learn to pick up objects

Traditional: 10,000 real robot trials
  Cost: $$$, time: weeks, hardware wear

With world models:
  1. Collect small amount of real grasp data
  2. Train world model on real experience
  3. Generate imaginary grasp scenarios in model
  4. Learn grasping policy in simulation
  5. Test on real robot (validation only)

Result: 10x faster learning, 10x less hardware cost

Research: Google Robotics, OpenAI, DeepMind
```

### 2. **Autonomous Driving**
```
Challenge: Self-driving car must predict pedestrian behavior

Traditional: Build explicit behavior models
  "Pedestrians cross at crosswalks"
  "Cyclists stay near curb"
  (Hand-crafted rules for every scenario)

With world models:
  1. Train model on real driving video (millions of hours)
  2. Predict: "If I see person walking at 5 mph toward crossing..."
  3. Model predicts likely futures (crosses? stops? continues?)
  4. Planner chooses safe action

Real-world: Tesla, Waymo, Cruise use variants of this

Companies: Tesla Autopilot, Waymo, Cruise, Zoox
```

### 3. **Game Playing & Planning**
```
AlphaGo Zero (world model variant):
  - Model of Go game: Position encoding
  - Dynamics: "If I play move X, what's the board state?"
  - Planning: Monte Carlo tree search (search using world model)
  - Result: Superhuman Go performance

Similar: AlphaZero (Chess, Shogi), MuZero (Atari games)

Key insight: World model + planning = superhuman performance
without needing vast amounts of training data
```

### 4. **Video Generation & Prediction**
```
OpenAI Sora (2024):
  - Trained to predict next video frames
  - Acts as world model: Understands physics, lighting, object interactions
  - Can generate realistic multi-second videos
  - Demonstrates emergent understanding of 3D, physics, causality

Implications:
  - World models are becoming mainstream
  - Scaling to real-world complexity (high-res video)
  - Used for planning, design, content creation
```

### 5. **Planning & Decision-Making**
```
Business process optimization:
  1. Model company operations (workflows, dependencies)
  2. Predict: "If we change process X, how does it affect Y?"
  3. Plan optimal changes without trial-and-error

Example: Supply chain optimization
  - World model: Relationships between suppliers, factories, demand
  - Prediction: "If supplier A delays, how does it cascade?"
  - Planning: Adjust orders/routing before problem occurs
```

### 6. **Scientific Discovery**
```
Physics simulation:
  - Train world model on particle simulation data
  - Model learns physics laws from data (no explicit rules)
  - Use model to predict systems, design experiments

Example: Protein folding
  - World model: Molecular dynamics
  - Prediction: How does protein fold given initial structure?
  - Application: Drug design, protein engineering
```

---

## Key Differences: World Models vs. Other AI Paradigms

### World Models vs. Large Language Models (LLMs)

| Aspect | World Model | LLM |
|--------|-------------|-----|
| **Training data** | Videos, sensor data, interactions | Text (internet scale) |
| **What it learns** | Physics, causality, environment dynamics | Language patterns, world knowledge |
| **Input** | Observations/actions | Text prompts |
| **Output** | Predicted next state/frame | Next token (text) |
| **Planning ability** | Explicit (forward simulate) | Implicit (pattern matching) |
| **Real-world grounding** | Strong (directly learns from reality) | Weak (only text, no embodiment) |
| **Hallucination** | Visible (wrong predictions) | Often subtle (plausible but false text) |
| **Speed** | Fast (inference in latent space) | Slower (token-by-token generation) |

**Key insight:** World models are grounded in physical reality; LLMs are grounded in human text patterns.

### World Models vs. Model-Free RL

| Aspect | World Model | Model-Free RL (DQN, PPO) |
|--------|-------------|-------------------------|
| **Planning** | Uses model to simulate | No simulation, trial-error |
| **Sample efficiency** | High (learns from model predictions) | Low (needs millions of real trials) |
| **Training time** | Faster (learn compact model) | Slow (interact with real world) |
| **Generalization** | Better (understands dynamics) | Poor (memorizes specific states) |
| **Interpretability** | Higher (can visualize predictions) | Lower (black-box policy) |
| **Practical cost** | Lower (fewer real interactions) | Higher (millions of interactions) |

**Example:**
- Model-free: "Try 1M robot grasps, learn which works"
- World model: "Learn physics in simulation, test few grasps on real robot"

---

## Training & Learning Dynamics

### The Learning Process

**Phase 1: Collect experience**
```
Agent or human plays environment (game, drives car, controls robot)
Record: observations + actions + next observations

Data collected: Video (millions of frames) or sensor logs
```

**Phase 2: Train encoder (optional)**
```
Compress observations into latent space
Loss: Reconstruction error (can we reconstruct from latent?)
Regularization: Keep latent space smooth/interpretable
Result: Encoder that captures essential environment structure
```

**Phase 3: Train dynamics model**
```
Input: Current latent state + action
Output: Predicted next latent state
Loss: Prediction error (does prediction match actual next state?)
Optimization: Backpropagation through time

Key challenge: Long-term error accumulation
  - Short-term predictions (1-2 steps): Easy
  - Long-term predictions (50+ steps): Hard
  Solution: Curriculum learning (start short, gradually extend horizon)
```

**Phase 4: Evaluate & iterate**
```
Test on held-out data:
  - How well do predictions match reality?
  - How far ahead can we predict accurately?
  - On which scenarios does model fail?

Iterate:
  - Collect more data on failing scenarios
  - Add architectural improvements
  - Increase model capacity
```

### Training Challenges

1. **Compounding prediction error** — Errors accumulate over long horizons
   - Solution: Latent space training (fewer pixels to error), stochastic dynamics

2. **Multimodal futures** — Many possible futures exist
   - Example: Ball could go left or right after hitting wall
   - Solution: Model distribution over futures (not just mean), use VAE

3. **Learning without labels** — How do we know what's important to predict?
   - Solution: Unsupervised learning + curiosity (predict what's uncertain)

4. **Sim-to-real gap** — Simulation ≠ reality
   - Example: Robot trained in simulation crashes on real robot
   - Solution: Randomization, domain adaptation, test on reality

---

## Key Research Milestones & Papers

### Canonical References

**1. "World Models" (Ha & Schmidhuber, 2018)**
- **Link:** https://en.wikipedia.org/wiki/World_model_(artificial_intelligence)
- **What it showed:** 
  - Agent can learn compact representation of game environment
  - Controller trained in latent space can play game
  - Model can generate imaginative rollouts (dreams)
- **Impact:** Defined modern world models research direction

**2. PlaNet: Deep Planning Network (Hafner et al., 2018)**
- **Research:** Google Brain
- **Key contribution:** Model-based RL using world models
- **Result:** Learned controller without trial-and-error; worked on image-based control

**3. Dreamer: Scalable Belief-Space Planning (Hafner et al., 2019-2021)**
- **Research:** Google Brain
- **Key contribution:** Improved planning with learned world models
- **Result:** State-of-the-art on Atari games using learned model

**4. Sora: Video generation as world modeling (OpenAI, 2024)**
- **Link:** https://openai.com/ (search "Sora")
- **What it shows:** 
  - Large-scale video diffusion model acts as world model
  - Can generate coherent multi-minute videos
  - Implicitly learned physics, lighting, perspective
- **Impact:** Demonstrates scaling to real-world complexity

**5. Generative model (Foundation models)**
- **Link:** https://en.wikipedia.org/wiki/Generative_model
- **What it shows:**
  - LLMs, diffusion models, GANs all learn implicit world models
  - Foundation models contain embedded world understanding
  - Emergent capability: Planning/reasoning without explicit training

---

## Practical Implementation: Getting Started

### Tools & Frameworks

**1. PyTorch** (Most common)
- **Link:** https://pytorch.org/
- **Use:** Build custom world model architectures
- **Example:**
  ```python
  import torch
  import torch.nn as nn
  
  # Encoder: compress image to latent
  class Encoder(nn.Module):
      def __init__(self):
          super().__init__()
          self.conv = nn.Sequential(
              nn.Conv2d(3, 32, 4, stride=2),
              nn.ReLU(),
              nn.Conv2d(32, 64, 4, stride=2),
              nn.ReLU(),
          )
          self.fc = nn.Linear(64*56*56, 128)  # latent_dim=128
      
      def forward(self, x):
          return self.fc(self.conv(x).view(x.size(0), -1))
  ```

**2. TensorFlow / Keras**
- **Use:** Pre-built layers, easier for beginners
- **Advantage:** Good tutorials, Colab notebooks

**3. JAX**
- **Use:** Research, custom gradient computation
- **Advantage:** Functional programming, JIT compilation

### Datasets for Training

1. **Atari games** (Benchmark)
   - Standard test domain
   - Open-source (Gym, Atari)

2. **Real robotics** (If you have hardware)
   - Robot learning datasets
   - Example: MIT Robotics, DeepMind robotics repos

3. **Simulated environments**
   - Gym environments (OpenAI)
   - MuJoCo (physics simulation)
   - PyBullet (robotics simulation)

4. **Video datasets** (For video prediction)
   - YouTube-8M
   - Kinetics
   - UCF101

### Beginner Code Example: Simple World Model

```python
# Pseudocode: Simple VAE-based world model
import torch
import torch.nn as nn

class WorldModel(nn.Module):
    def __init__(self, latent_dim=32):
        super().__init__()
        
        # Encoder: Image → Latent
        self.encoder = nn.Sequential(
            nn.Conv2d(3, 32, 4, stride=2), nn.ReLU(),
            nn.Conv2d(32, 64, 4, stride=2), nn.ReLU(),
            nn.Flatten(),
            nn.Linear(64*56*56, latent_dim)
        )
        
        # Dynamics: Latent + Action → Next Latent
        self.dynamics = nn.LSTM(latent_dim + action_dim, latent_dim, batch_first=True)
        
        # Decoder: Latent → Image
        self.decoder = nn.Sequential(
            nn.Linear(latent_dim, 64*56*56),
            nn.ReLU(),
            nn.ConvTranspose2d(64, 32, 4, stride=2), nn.ReLU(),
            nn.ConvTranspose2d(32, 3, 4, stride=2),
        )
    
    def forward(self, observations, actions):
        # Encode observations to latent
        latents = self.encoder(observations)
        
        # Predict next latents
        latent_actions = torch.cat([latents, actions], dim=-1)
        next_latents, _ = self.dynamics(latent_actions)
        
        # Decode predicted latents
        predictions = self.decoder(next_latents)
        
        return predictions, latents

# Training loop:
# for batch in dataloader:
#     obs, actions, next_obs = batch
#     pred, latents = model(obs, actions)
#     loss = MSE(pred, next_obs)  # Prediction error
#     loss.backward()
#     optimizer.step()
```

---

## Current Frontiers & Open Problems

### 1. **Long-horizon prediction**
- Challenge: Predictions get worse after 10+ steps
- Active research: Hierarchical models, stochastic dynamics, uncertainty quantification

### 2. **Scaling to real-world complexity**
- Challenge: Real environments are high-dimensional, partially observable
- Progress: Sora demonstrates scaling to realistic video
- Next: World models for real robotics at scale

### 3. **Multimodal futures**
- Challenge: Multiple futures are possible (agent can choose)
- Solution: Probabilistic models, ensemble methods
- Active research: Diffusion-based world models

### 4. **Sim-to-real transfer**
- Challenge: Learned models in simulation don't work in reality
- Approaches: Domain randomization, world model adaptation, real-time updates

### 5. **Interpretability**
- Challenge: What does the world model actually learn?
- Research: Visualization, abstraction discovery, concept learning

### 6. **Integration with LLMs**
- Challenge: How to combine explicit world models with implicit knowledge in LLMs?
- Active research: Neurosymbolic approaches, retrieval-augmented planning

---

## World Models vs. Alternatives: When to Use

### Use World Models If:
✓ You need to **plan** (not just react)  
✓ You have **low trial budgets** (expensive real-world interactions)  
✓ You need **generalization** to new scenarios  
✓ You want **interpretability** (can visualize predictions)  
✓ Your environment has **consistent physics/rules**  

### Use Traditional RL If:
- Environment is simple (few states/actions)
- Trial-and-error is cheap (simulation)
- Model learning overhead isn't worth it

### Use LLMs If:
- Knowledge is primarily linguistic/conceptual
- You need common-sense reasoning
- Domain is text-based

### Use Symbolic Methods If:
- Environment structure is known/simple
- You need hard guarantees
- Efficiency is critical

---

## Real Examples in the Wild

### 1. **Tesla Autopilot**
- Uses models of road, pedestrians, physics
- Predicts: "Will pedestrian reach road by time I arrive?"
- Plans: Adjust speed/steering accordingly

### 2. **OpenAI Sora (Video Generation)**
- World model demonstrates physics understanding
- Can generate realistic physics-consistent video

### 3. **DeepMind AlphaGo**
- Uses world model of Go board state
- Plans by searching game tree using model

### 4. **Boston Dynamics Robots**
- Learn models of their bodies, environment
- Use for complex locomotion planning

### 5. **Google's Gato (2022)**
- Multimodal model with implicit world understanding
- Can control Atari games, real robots, play Minecraft

---

## Learning Resources

### Official Documentation
- **PyTorch:** https://pytorch.org/
- **Wikipedia - World Models:** https://en.wikipedia.org/wiki/World_model_(artificial_intelligence)
- **Wikipedia - Generative Models:** https://en.wikipedia.org/wiki/Generative_model

### Landmark Papers (Accessible)
1. Ha & Schmidhuber, "World Models" (2018) — Best starting point
2. Hafner et al., "PlaNet" (2018) — Deep planning networks
3. Hafner et al., "Dreamer" (2019) — Scalable belief planning

### Implementations
- **Gym:** https://github.com/openai/gym (Atari environments)
- **MuJoCo:** https://mujoco.org/ (Physics simulation)
- **PyBullet:** https://pybullet.org/ (Robotics simulation)

### Courses & Tutorials
- **MIT 6.S191** (Intro to Deep Learning) — Covers RL and world models
- **DeepMind RL course** (YouTube) — Includes model-based RL
- **Stanford CS231N** — Includes video prediction lecture

---

## Integration with OneMind OS / Tactical Operations

### Hypothetical Application: Tactical Planning

```
Scenario: Combat operations in contested area

World model application:
  1. Train model on historical/simulated battle data
     - Input: Enemy positions, terrain, weather, friendly forces
     - Dynamics: How do forces move/interact?
     - Output: Predicted future state of battlefield

  2. Use for rapid planning
     - Commander asks: "If I move Delta squad here, what happens?"
     - World model predicts 10 likely futures
     - Commander picks best scenario
     - Faster decision-making (seconds, not hours)

  3. Continuous adaptation
     - Real intel arrives (captured positions, intelligence)
     - Update world model beliefs
     - Re-predict battlefield evolution
     - Adjust plan in real-time

Integration points:
  - NATS: Receive position updates (agents.tactical.position)
  - DDF: Query historical intelligence patterns
  - TAK Server: Display predicted futures on COP
  - Mumble: Broadcast plan updates
```

### Hypothetical Application: Autonomous Tactical Systems

```
Scenario: Unmanned vehicle operating in GPS-denied environment

World model approach:
  1. Vehicle learns model of terrain/obstacles (from drone footage)
  2. Plans path without GPS (uses learned model)
  3. Predicts sensor readings based on model
  4. Detects anomalies (mismatch between prediction and reality)
  5. Alerts human operator to unexpected changes

Benefits:
  - GPS-denied autonomy (no signal needed)
  - Anomaly detection (difference from expected)
  - Graceful degradation (works without comms)
```

---

## Summary: What World Models Enable

| Capability | Impact |
|------------|--------|
| **Planning without trial-and-error** | 10-100x faster learning for embodied agents |
| **Reasoning about unseen scenarios** | Generalization beyond training data |
| **Simulation before action** | Lower cost, higher safety |
| **Credit assignment** | Learn which decisions matter |
| **Interpretability** | Visualize agent's understanding |
| **Anomaly detection** | Detect when reality differs from expectations |

---

**Captured:** 2026-07-09  
**Status:** Active research area + production deployment  
**Relevance:** Foundation for AI planning, robotics, autonomous systems, tactical decision support  
**Learning curve:** Moderate (requires ML background, strong for grad-level study)  
**Practical impact:** High (enables sample-efficient learning, planning, simulation)  
**Key players:** Google DeepMind, OpenAI, Meta, academia  
**Next horizon:** Scaling to real-world robotics, multimodal world understanding, integration with foundation models

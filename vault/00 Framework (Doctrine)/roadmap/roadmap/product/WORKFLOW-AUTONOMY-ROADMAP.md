# Workflow Autonomy Roadmap

## Vision

Transform workflows from **"smart rails"** (predefined steps with decision points) to **"autonomous agents"** (self-directing, learning, adaptive systems).

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         WORKFLOW EVOLUTION                                   │
│                                                                              │
│   CURRENT              PHASE 2              PHASE 3              PHASE 4    │
│   Smart Rails    →    Adaptive Rails   →   Dynamic Steps    →   Autonomous  │
│                                                                              │
│   ┌─────────┐         ┌─────────┐          ┌─────────┐         ┌─────────┐  │
│   │ Step 1  │         │ Step 1  │          │ Step ?  │         │ Agent   │  │
│   │ Step 2  │         │ Step 2  │←retry    │ (LLM    │         │ decides │  │
│   │ Step 3  │         │ Step 3  │←skip     │ decides)│         │ & learns│  │
│   │ (fixed) │         │ (adapt) │          │ Step ?  │         │ forever │  │
│   └─────────┘         └─────────┘          └─────────┘         └─────────┘  │
│                                                                              │
│   Predefined          Retry/skip           Generate steps       Self-modify │
│   structure           based on             at runtime           & remember  │
│                       quality                                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Current State: Smart Rails ✅

**What exists:**
- 8 step types: `step`, `parallel`, `condition`, `router`, `loop`, `approval`, `function`, `team`
- Event triggers: Schedule (cron), Event (NATS), HTTP (webhook)
- Smart approvals: LLM assesses risk, auto-approves low-risk steps
- Context threading: `WorkflowContext` carries state across steps
- Output piping: `{step_name.field}` syntax
- Checkpointing: Recovery from failures
- Caching: Identical steps skip re-execution

**Limitations:**
- Steps are predefined at workflow creation
- Cannot dynamically spawn new agents
- No learning from outcomes
- No self-modification

---

## Phase 2: Adaptive Rails (1-2 Sprints)

### Goal
Workflows can **retry with different strategies** and **skip steps** based on outcome quality.

### Features

#### 2.1 Outcome Feedback System
Track quality/performance after each step execution.

```python
# New model: StepOutcome
class StepOutcome(BaseModel):
    workflow_run_id: str
    step_name: str
    agent_id: Optional[str]

    # Quality metrics
    quality_score: float  # 0.0 - 1.0, LLM-assessed
    execution_time_ms: int
    token_usage: int
    error_occurred: bool

    # Metadata
    input_hash: str  # For caching decisions
    output_summary: str  # Short summary of result
    timestamp: datetime
```

#### 2.2 Quality Assessment
LLM evaluates step output quality.

```python
async def assess_step_quality(
    step: WorkflowStep,
    result: StepResult,
    context: WorkflowContext,
) -> float:
    """
    Use Claude to assess if the step output is high-quality.

    Returns:
        Score 0.0 - 1.0
    """
    prompt = f"""
    Assess the quality of this workflow step output.

    Step: {step.name}
    Task: {step.description}
    Expected: {step.expected_output or "Not specified"}

    Actual Output:
    {result.content[:2000]}

    Rate quality from 0.0 (unusable) to 1.0 (excellent).
    Consider: completeness, accuracy, relevance, formatting.

    Return ONLY a number between 0.0 and 1.0.
    """
    # Quick assessment with Haiku
    score = await assess_with_llm(prompt, model="haiku")
    return float(score)
```

#### 2.3 Adaptive Retry Logic
If quality is low, retry with different agent or modified prompt.

```python
class AdaptiveExecutor(WorkflowExecutor):
    async def execute_step_with_retry(
        self,
        step: WorkflowStep,
        context: WorkflowContext,
        max_retries: int = 2,
        quality_threshold: float = 0.6,
    ) -> StepResult:
        """Execute step, retry if quality is below threshold."""

        for attempt in range(max_retries + 1):
            result = await self.execute_step(step, context)

            # Assess quality
            quality = await assess_step_quality(step, result, context)

            # Record outcome
            await self.record_outcome(step, result, quality, attempt)

            if quality >= quality_threshold:
                return result

            if attempt < max_retries:
                # Try different strategy
                step = await self.adapt_step(step, result, quality, attempt)
                logger.info(f"Retrying {step.name} with adapted strategy (attempt {attempt + 2})")

        # Return best attempt even if below threshold
        return result

    async def adapt_step(
        self,
        step: WorkflowStep,
        result: StepResult,
        quality: float,
        attempt: int,
    ) -> WorkflowStep:
        """Modify step for retry attempt."""

        strategies = [
            self._try_different_agent,
            self._enhance_prompt,
            self._increase_context,
        ]

        strategy = strategies[attempt % len(strategies)]
        return await strategy(step, result)
```

#### 2.4 Skip Logic
Skip steps that aren't needed based on context.

```python
async def should_skip_step(
    step: WorkflowStep,
    context: WorkflowContext,
) -> tuple[bool, str]:
    """
    Determine if step should be skipped.

    Returns:
        (should_skip, reason)
    """
    # Check if output already exists in context
    if step.output_key and step.output_key in context.data:
        return True, "Output already exists in context"

    # Check if step is relevant to current task
    if step.skip_condition:
        should_skip = eval_condition(step.skip_condition, context)
        if should_skip:
            return True, f"Skip condition met: {step.skip_condition}"

    # LLM decides if step is necessary
    if step.allow_smart_skip:
        is_necessary = await assess_step_necessity(step, context)
        if not is_necessary:
            return True, "LLM determined step is not necessary"

    return False, ""
```

### Deliverables
- [ ] `StepOutcome` model and storage
- [ ] `OutcomeStore` for persisting feedback
- [ ] `assess_step_quality()` function
- [ ] `AdaptiveExecutor` class
- [ ] `should_skip_step()` function
- [ ] Outcome dashboard in frontend (optional)

### Success Criteria
- Steps retry when quality < 0.6
- Different agent selected on retry
- Outcome history persisted
- 20% fewer failed workflows

---

## Phase 3: Dynamic Step Generation (3-4 Sprints)

### Goal
LLM can **generate new steps at runtime** based on context and goals.

### Features

#### 3.1 Step Generator
LLM analyzes context and decides what step to execute next.

```python
class DynamicStepGenerator:
    async def generate_next_step(
        self,
        goal: str,
        context: WorkflowContext,
        completed_steps: List[StepResult],
        available_agents: List[str],
        available_tools: List[str],
    ) -> Optional[WorkflowStep]:
        """
        Generate the next step to achieve the goal.

        Returns:
            WorkflowStep if more work needed, None if goal achieved.
        """
        prompt = f"""
        You are a workflow orchestrator. Analyze the current state and decide the next step.

        GOAL: {goal}

        COMPLETED STEPS:
        {self._format_completed_steps(completed_steps)}

        CURRENT CONTEXT:
        {json.dumps(context.data, indent=2)}

        AVAILABLE AGENTS:
        {', '.join(available_agents)}

        AVAILABLE TOOLS:
        {', '.join(available_tools)}

        Decide:
        1. Is the goal achieved? If yes, return {"done": true}
        2. If not, what step should execute next?

        Return JSON:
        {{
            "done": false,
            "step": {{
                "name": "step_name",
                "type": "step|parallel|condition|team",
                "agent_id": "agent-to-use",
                "prompt": "What the agent should do",
                "output_key": "where_to_store_result"
            }},
            "reasoning": "Why this step is needed"
        }}
        """

        response = await self.llm.generate(prompt)
        parsed = json.loads(response)

        if parsed.get("done"):
            return None

        return self._build_step(parsed["step"])
```

#### 3.2 Goal-Driven Workflow
Workflows defined by goal, not steps.

```python
class GoalDrivenWorkflow:
    def __init__(
        self,
        goal: str,
        constraints: List[str] = None,
        max_steps: int = 10,
        available_agents: List[str] = None,
    ):
        self.goal = goal
        self.constraints = constraints or []
        self.max_steps = max_steps
        self.available_agents = available_agents or ["legacy-mother-ai"]

    async def run(self, initial_input: Dict[str, Any]) -> WorkflowResult:
        """Execute workflow by generating steps until goal achieved."""

        context = WorkflowContext(data=initial_input)
        completed_steps = []
        generator = DynamicStepGenerator()

        for i in range(self.max_steps):
            # Generate next step
            next_step = await generator.generate_next_step(
                goal=self.goal,
                context=context,
                completed_steps=completed_steps,
                available_agents=self.available_agents,
            )

            if next_step is None:
                # Goal achieved
                return WorkflowResult(
                    status="completed",
                    output=context.data,
                    steps_executed=len(completed_steps),
                )

            # Execute step
            result = await self.executor.execute_step(next_step, context)
            completed_steps.append(result)

            # Update context
            context.data[next_step.output_key] = result.output

        return WorkflowResult(
            status="max_steps_reached",
            output=context.data,
            steps_executed=len(completed_steps),
        )
```

#### 3.3 Nested Workflow Calls
Workflows can spawn sub-workflows.

```python
class WorkflowCallStep(WorkflowStep):
    """Step that executes another workflow."""

    type: Literal["workflow"] = "workflow"
    workflow_id: str  # ID of workflow to call
    input_mapping: Dict[str, str]  # Map context keys to workflow input
    output_key: str  # Where to store workflow result

    async def execute(self, context: WorkflowContext) -> StepResult:
        # Get workflow definition
        workflow = registry.get(self.workflow_id)

        # Map input from context
        workflow_input = {
            k: context.data.get(v)
            for k, v in self.input_mapping.items()
        }

        # Execute nested workflow
        result = await executor.run(
            workflow_id=self.workflow_id,
            steps=workflow.steps,
            input=workflow_input,
        )

        return StepResult(
            step_name=self.name,
            output=result.output,
            status="completed",
        )
```

### Deliverables
- [ ] `DynamicStepGenerator` class
- [ ] `GoalDrivenWorkflow` class
- [ ] `WorkflowCallStep` for nested workflows
- [ ] Step validation before execution
- [ ] Depth limit for nested calls (prevent infinite loops)
- [ ] Goal achievement assessment

### Success Criteria
- Can execute workflow with just a goal (no predefined steps)
- LLM generates appropriate steps
- Nested workflows execute correctly
- Max depth enforced

---

## Phase 4: Full Autonomy (6-8 Sprints)

### Goal
Workflows become **self-improving agents** that learn, adapt, and optimize over time.

### Features

#### 4.1 Learning Store
Persistent storage of outcomes, patterns, and learned knowledge.

```python
class LearningStore:
    """Long-term memory for workflow learning."""

    async def record_workflow_outcome(
        self,
        workflow_id: str,
        goal: str,
        steps_executed: List[StepOutcome],
        final_quality: float,
        total_time_ms: int,
    ):
        """Record complete workflow execution for learning."""
        pass

    async def get_best_strategy_for_goal(
        self,
        goal: str,
        similar_goals: int = 5,
    ) -> Optional[WorkflowStrategy]:
        """Find best-performing strategy for similar goals."""
        # Semantic search for similar past goals
        # Return strategy that worked best
        pass

    async def get_agent_performance(
        self,
        agent_id: str,
        task_type: str,
    ) -> AgentPerformanceMetrics:
        """Get historical performance for agent on task type."""
        pass

    async def suggest_improvements(
        self,
        workflow_id: str,
    ) -> List[Improvement]:
        """Analyze workflow history, suggest optimizations."""
        pass
```

#### 4.2 Self-Modification
Workflows modify their own structure based on outcomes.

```python
class SelfModifyingWorkflow:
    async def optimize(self) -> List[Modification]:
        """Analyze past runs, modify workflow for better performance."""

        modifications = []

        # Analyze step outcomes
        outcomes = await self.learning_store.get_step_outcomes(self.workflow_id)

        for step_name, step_outcomes in outcomes.items():
            avg_quality = mean([o.quality_score for o in step_outcomes])
            avg_time = mean([o.execution_time_ms for o in step_outcomes])

            # Low quality step - try different agent
            if avg_quality < 0.6:
                better_agent = await self.find_better_agent(step_name)
                if better_agent:
                    modifications.append(Modification(
                        type="change_agent",
                        step_name=step_name,
                        old_value=self.steps[step_name].agent_id,
                        new_value=better_agent,
                    ))

            # Slow step - consider parallelization
            if avg_time > 30000:  # 30 seconds
                can_parallelize = await self.check_parallelization(step_name)
                if can_parallelize:
                    modifications.append(Modification(
                        type="parallelize",
                        step_name=step_name,
                    ))

        return modifications

    async def apply_modifications(self, modifications: List[Modification]):
        """Apply learned modifications to workflow definition."""
        for mod in modifications:
            if mod.type == "change_agent":
                self.steps[mod.step_name].agent_id = mod.new_value
            elif mod.type == "parallelize":
                self.parallelize_step(mod.step_name)
            # ... other modification types

        # Persist updated workflow
        await self.save()
```

#### 4.3 Continuous Learning Loop
Background process that continuously improves workflows.

```python
class WorkflowLearningLoop:
    """Background service that optimizes workflows over time."""

    async def run_forever(self):
        while True:
            # Get workflows with recent runs
            workflows = await self.get_active_workflows()

            for workflow in workflows:
                # Check if enough data for learning
                if await self.has_enough_data(workflow.id):
                    # Analyze and suggest improvements
                    improvements = await self.analyze_workflow(workflow)

                    # Apply low-risk improvements automatically
                    auto_apply = [i for i in improvements if i.risk == "low"]
                    if auto_apply:
                        await workflow.apply_modifications(auto_apply)
                        logger.info(f"Auto-applied {len(auto_apply)} improvements to {workflow.id}")

                    # Queue high-risk improvements for human review
                    human_review = [i for i in improvements if i.risk != "low"]
                    if human_review:
                        await self.queue_for_review(workflow.id, human_review)

            # Run every hour
            await asyncio.sleep(3600)
```

#### 4.4 Pattern Discovery
Discover reusable patterns from successful workflows.

```python
class PatternDiscovery:
    async def discover_patterns(self) -> List[WorkflowPattern]:
        """Analyze successful workflows, extract reusable patterns."""

        # Get high-quality workflow runs
        successful_runs = await self.learning_store.get_successful_runs(
            min_quality=0.8,
            min_runs=10,
        )

        # Cluster similar step sequences
        patterns = await self.cluster_step_sequences(successful_runs)

        # Extract common patterns
        discovered = []
        for pattern in patterns:
            if pattern.frequency > 5:  # Used in 5+ workflows
                discovered.append(WorkflowPattern(
                    name=f"pattern_{pattern.id}",
                    steps=pattern.steps,
                    typical_goal=pattern.common_goal,
                    success_rate=pattern.success_rate,
                ))

        return discovered
```

### Deliverables
- [ ] `LearningStore` with vector search
- [ ] `SelfModifyingWorkflow` class
- [ ] `WorkflowLearningLoop` background service
- [ ] `PatternDiscovery` for reusable patterns
- [ ] Human review queue for high-risk changes
- [ ] Learning dashboard
- [ ] Pattern library

### Success Criteria
- Workflows improve over time without human intervention
- Agent selection optimized based on performance
- Reusable patterns discovered and applied
- Human oversight for risky changes

---

## Architecture Evolution

### Current
```
WorkflowDefinition (static)
    ↓
WorkflowExecutor (sequential)
    ↓
StepResult (output only)
```

### Phase 4
```
LearningStore (memory)
    ↑↓
GoalDrivenWorkflow (dynamic)
    ↓
AdaptiveExecutor (retry/skip)
    ↓
DynamicStepGenerator (LLM)
    ↓
StepOutcome (quality + metrics)
    ↓
SelfModifyingWorkflow (optimize)
    ↓
PatternLibrary (reuse)
```

---

## Timeline

| Phase | Duration | Key Milestone |
|-------|----------|---------------|
| **Current** | Done | Smart rails + approvals working |
| **Phase 2** | 1-2 sprints | Adaptive retry + outcome tracking |
| **Phase 3** | 3-4 sprints | Goal-driven workflows + dynamic steps |
| **Phase 4** | 6-8 sprints | Self-improving autonomous workflows |

**Total: 10-14 sprints** for full autonomy

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Infinite loops | Max step limit (10), depth limit for nested workflows |
| Runaway costs | Token budgets per workflow, quality gates |
| Bad modifications | Human review for high-risk changes |
| Over-optimization | A/B testing before applying changes |
| Data loss | Versioned workflow definitions, rollback capability |

---

## Success Metrics

| Metric | Current | Phase 2 | Phase 3 | Phase 4 |
|--------|---------|---------|---------|---------|
| Manual step definition | 100% | 100% | 50% | 10% |
| Retry success rate | 0% | 60% | 70% | 80% |
| Self-optimization | 0% | 0% | 0% | 50% |
| Pattern reuse | 0% | 0% | 20% | 60% |
| Human intervention | High | Medium | Low | Minimal |

---

## Next Steps

1. **Immediate**: Complete Phase 2 design review
2. **Sprint 1**: Implement `StepOutcome` + `OutcomeStore`
3. **Sprint 2**: Implement `AdaptiveExecutor` with retry logic
4. **Sprint 3**: Add quality assessment
5. **Review**: Assess Phase 2, plan Phase 3 details

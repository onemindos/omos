# Agent Evaluations Guide

Evaluations (evals) measure the quality of AI agents across multiple dimensions: accuracy, performance, reliability, and custom criteria.

## Why Evaluate Agents?

| Reason | Description |
|--------|-------------|
| **Quality Assurance** | Ensure agents give correct answers |
| **Regression Testing** | Catch issues when changing models/prompts |
| **Model Comparison** | Compare GPT-4 vs Claude vs Gemini |
| **Tool Verification** | Verify agents use tools correctly |
| **Performance Tuning** | Measure latency and resource usage |

## Agno Native Evals

Use Agno's built-in evaluation framework:

```python
from agno.eval.accuracy import AccuracyEval
from agno.eval.performance import PerformanceEval
from agno.eval.reliability import ReliabilityEval
from agno.eval.agent_as_judge import AgentAsJudgeEval
```

## Evaluation Types

### 1. Accuracy Eval

Tests if agent responses match expected outputs.

```python
from agno.agent import Agent
from agno.eval.accuracy import AccuracyEval, AccuracyResult
from agno.models.openai import OpenAIChat
from agno.tools.calculator import CalculatorTools

# Create agent to evaluate
agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    tools=[CalculatorTools()],
)

# Create evaluation
evaluation = AccuracyEval(
    model=OpenAIChat(id="o4-mini"),  # Judge model
    agent=agent,
    input="What is 10*5 then to the power of 2?",
    expected_output="2500",
    additional_guidelines="Output should include steps and final answer.",
)

# Run evaluation
result: AccuracyResult = evaluation.run(print_results=True)

# Assert quality
assert result is not None and result.avg_score >= 8
```

### 2. Performance Eval

Measures latency and resource usage.

```python
from agno.eval.performance import PerformanceEval

evaluation = PerformanceEval(
    agent=agent,
    input="What is the capital of France?",
    num_iterations=10,  # Run multiple times for averages
)

result = evaluation.run(print_results=True)

print(f"Avg latency: {result.avg_latency_ms}ms")
print(f"Memory usage: {result.memory_mb}MB")
print(f"Tokens/second: {result.tokens_per_second}")
```

### 3. Reliability Eval

Verifies tool calls and error handling.

```python
from agno.eval.reliability import ReliabilityEval

evaluation = ReliabilityEval(
    agent=agent,
    input="Calculate 5 * 10",
    expected_tool_calls=["multiply"],  # Should use calculator
)

result = evaluation.run(print_results=True)

print(f"Tool call success rate: {result.tool_call_success_rate}")
print(f"Error rate: {result.error_rate}")
```

### 4. Agent-as-Judge Eval

Uses an LLM to judge response quality with custom criteria.

```python
from agno.eval.agent_as_judge import AgentAsJudgeEval

evaluation = AgentAsJudgeEval(
    model=OpenAIChat(id="gpt-4o"),  # Judge model
    agent=agent,
    input="Explain quantum computing to a 5-year-old",
    criteria=[
        "Uses simple language appropriate for children",
        "Avoids technical jargon",
        "Uses relatable analogies",
        "Is engaging and fun",
    ],
)

result = evaluation.run(print_results=True)

for criterion, score in result.scores.items():
    print(f"{criterion}: {score}/10")
```

## Running Multiple Evaluations

### Test Suite Pattern

```python
from agno.agent import Agent
from agno.eval.accuracy import AccuracyEval
from agno.models.openai import OpenAIChat

# Define test cases
test_cases = [
    {"input": "What is 2+2?", "expected": "4"},
    {"input": "Capital of France?", "expected": "Paris"},
    {"input": "What year was Python created?", "expected": "1991"},
]

# Create agent
agent = Agent(model=OpenAIChat(id="gpt-4o"))

# Run all tests
results = []
for test in test_cases:
    eval = AccuracyEval(
        model=OpenAIChat(id="o4-mini"),
        agent=agent,
        input=test["input"],
        expected_output=test["expected"],
    )
    result = eval.run()
    results.append({
        "input": test["input"],
        "passed": result.avg_score >= 8,
        "score": result.avg_score,
    })

# Summary
passed = sum(1 for r in results if r["passed"])
print(f"Passed: {passed}/{len(results)}")
```

### Comparing Models

```python
from agno.agent import Agent
from agno.eval.accuracy import AccuracyEval
from agno.models.openai import OpenAIChat
from agno.models.anthropic import Claude

# Same test, different models
models = [
    ("GPT-4o", OpenAIChat(id="gpt-4o")),
    ("GPT-4o-mini", OpenAIChat(id="gpt-4o-mini")),
    ("Claude Sonnet", Claude(id="claude-3-5-sonnet-20241022")),
]

test_input = "What is the derivative of x^2?"
expected = "2x"

for name, model in models:
    agent = Agent(model=model)
    eval = AccuracyEval(
        model=OpenAIChat(id="o4-mini"),
        agent=agent,
        input=test_input,
        expected_output=expected,
    )
    result = eval.run()
    print(f"{name}: {result.avg_score}/10")
```

## Best Practices

### 1. Use Appropriate Judge Models

```python
# For accuracy evals, use a capable judge model
AccuracyEval(
    model=OpenAIChat(id="o4-mini"),  # Good: reasoning model as judge
    # model=OpenAIChat(id="gpt-4o-mini"),  # OK but less reliable
    agent=agent,
    ...
)
```

### 2. Run Multiple Iterations

```python
# Performance evals need multiple runs for reliable averages
PerformanceEval(
    agent=agent,
    input="...",
    num_iterations=20,  # At least 10-20 for stable metrics
)
```

### 3. Test Edge Cases

```python
edge_cases = [
    # Empty input
    {"input": "", "expected": "clarification request"},
    # Very long input
    {"input": "..." * 1000, "expected": "handles gracefully"},
    # Ambiguous question
    {"input": "What is the best?", "expected": "asks for clarification"},
    # Out of scope
    {"input": "Write malware", "expected": "refuses appropriately"},
]
```

### 4. CI/CD Integration

```python
# pytest test file
import pytest
from agno.agent import Agent
from agno.eval.accuracy import AccuracyEval
from agno.models.openai import OpenAIChat

@pytest.fixture
def agent():
    return Agent(model=OpenAIChat(id="gpt-4o"))

def test_math_accuracy(agent):
    eval = AccuracyEval(
        model=OpenAIChat(id="o4-mini"),
        agent=agent,
        input="What is 15% of 200?",
        expected_output="30",
    )
    result = eval.run()
    assert result.avg_score >= 8, f"Score too low: {result.avg_score}"

def test_knowledge_accuracy(agent):
    eval = AccuracyEval(
        model=OpenAIChat(id="o4-mini"),
        agent=agent,
        input="Who wrote Romeo and Juliet?",
        expected_output="Shakespeare",
    )
    result = eval.run()
    assert result.avg_score >= 8
```

## Metrics Reference

### AccuracyResult

| Field | Description |
|-------|-------------|
| `avg_score` | Average score (0-10) |
| `passed` | Whether threshold met |
| `reasoning` | Judge's explanation |
| `latency_ms` | Time to generate response |

### PerformanceResult

| Field | Description |
|-------|-------------|
| `avg_latency_ms` | Average response time |
| `p50_latency_ms` | 50th percentile latency |
| `p95_latency_ms` | 95th percentile latency |
| `memory_mb` | Memory usage |
| `tokens_per_second` | Throughput |

### ReliabilityResult

| Field | Description |
|-------|-------------|
| `tool_call_success_rate` | % of correct tool calls |
| `error_rate` | % of runs with errors |
| `tool_calls` | List of tool calls made |

## References

- [Agno Evals Documentation](https://docs.agno.com/basics/evals/overview)
- [AccuracyEval Reference](https://docs.agno.com/basics/evals/accuracy/overview)
- [PerformanceEval Reference](https://docs.agno.com/basics/evals/performance)

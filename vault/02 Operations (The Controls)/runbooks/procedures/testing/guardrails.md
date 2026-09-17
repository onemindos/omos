# Guardrails Guide

Guardrails protect your AI agents by validating inputs and outputs, blocking harmful content before it causes issues.

## Why Use Guardrails?

| Risk | Guardrail Solution |
|------|-------------------|
| Violent content | Block violence/graphic categories |
| Hate speech | Block hate/threatening categories |
| Self-harm | Block self-harm instructions |
| Inappropriate content | Block sexual/harassment categories |
| Prompt injection | Validate input patterns |
| PII leakage | Filter sensitive data from outputs |

## Agno Native Guardrails

Use Agno's built-in guardrails:

```python
from agno.guardrails import OpenAIModerationGuardrail
from agno.agent import Agent
from agno.models.openai import OpenAIChat
```

## OpenAI Moderation Guardrail

Detects content violating OpenAI's content policy before processing.

### Basic Usage

```python
from agno.guardrails import OpenAIModerationGuardrail
from agno.agent import Agent
from agno.models.openai import OpenAIChat

# Create guardrail
guardrail = OpenAIModerationGuardrail()

# Attach to agent via pre_hooks
agent = Agent(
    name="Safe Agent",
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[guardrail],  # Runs before processing input
)

# Safe content passes through
agent.print_response("Explain quantum computing")  # OK

# Harmful content is blocked
agent.print_response("How to make weapons")  # Raises InputCheckError
```

### Custom Categories

Block only specific content categories:

```python
from agno.guardrails import OpenAIModerationGuardrail

# Only block violence and hate speech
guardrail = OpenAIModerationGuardrail(
    raise_for_categories=[
        "violence",
        "violence/graphic",
        "hate",
        "hate/threatening",
    ]
)

agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[guardrail],
)
```

### Available Categories

| Category | Description |
|----------|-------------|
| `violence` | Violent content |
| `violence/graphic` | Graphic violence |
| `hate` | Hate speech |
| `hate/threatening` | Threatening hate speech |
| `self-harm` | Self-harm content |
| `self-harm/intent` | Self-harm intent |
| `self-harm/instructions` | Self-harm instructions |
| `sexual` | Sexual content |
| `sexual/minors` | Child exploitation (always blocked) |
| `harassment` | Harassment |
| `harassment/threatening` | Threatening harassment |

### Handling Blocked Content

```python
from agno.exceptions import InputCheckError

try:
    response = agent.run("potentially harmful content")
except InputCheckError as e:
    print(f"Content blocked: {e.message}")
    print(f"Triggered by: {e.check_trigger}")
    # Handle gracefully - maybe ask user to rephrase
```

## Pre-hooks vs Post-hooks

| Hook Type | When it Runs | Use Case |
|-----------|-------------|----------|
| `pre_hooks` | Before agent processes input | Block harmful inputs |
| `post_hooks` | After agent generates output | Filter harmful outputs |

```python
agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[input_guardrail],   # Check input
    post_hooks=[output_guardrail], # Check output
)
```

## Custom Guardrails

Create your own guardrails for specific needs:

```python
from agno.exceptions import InputCheckError

def custom_guardrail(input_text: str) -> str:
    """Block requests containing specific keywords."""
    blocked_words = ["password", "api_key", "secret"]

    for word in blocked_words:
        if word.lower() in input_text.lower():
            raise InputCheckError(
                message=f"Input contains blocked keyword: {word}",
                check_trigger="custom_keyword_filter",
            )

    return input_text  # Pass through if OK

agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[custom_guardrail],
)
```

## PII Filter Example

```python
import re
from agno.exceptions import InputCheckError

def pii_filter(text: str) -> str:
    """Block inputs containing PII patterns."""
    patterns = {
        "SSN": r"\b\d{3}-\d{2}-\d{4}\b",
        "Credit Card": r"\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b",
        "Email": r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b",
    }

    for name, pattern in patterns.items():
        if re.search(pattern, text):
            raise InputCheckError(
                message=f"Input contains {name} - please remove sensitive data",
                check_trigger="pii_filter",
            )

    return text

agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[pii_filter],
)
```

## Multiple Guardrails

Chain multiple guardrails together:

```python
from agno.guardrails import OpenAIModerationGuardrail

agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[
        OpenAIModerationGuardrail(),  # Content moderation
        pii_filter,                    # PII detection
        custom_guardrail,              # Custom rules
    ],
)
```

Guardrails run in order - first one to fail stops processing.

## Async Guardrails

For async agents:

```python
import asyncio
from agno.guardrails import OpenAIModerationGuardrail
from agno.exceptions import InputCheckError

async def main():
    agent = Agent(
        model=OpenAIChat(id="gpt-4o"),
        pre_hooks=[OpenAIModerationGuardrail()],
    )

    try:
        await agent.aprint_response("Safe question about AI")
    except InputCheckError as e:
        print(f"Blocked: {e.message}")

asyncio.run(main())
```

## Best Practices

### 1. Always Use Moderation for Public-Facing Agents

```python
# Any agent exposed to users should have moderation
public_agent = Agent(
    model=OpenAIChat(id="gpt-4o"),
    pre_hooks=[OpenAIModerationGuardrail()],
)
```

### 2. Log Blocked Content for Review

```python
import logging

logger = logging.getLogger("guardrails")

def logging_guardrail(text: str) -> str:
    # Your validation logic
    if is_blocked(text):
        logger.warning(f"Blocked content: {text[:100]}...")
        raise InputCheckError(message="Content blocked")
    return text
```

### 3. Provide Helpful Error Messages

```python
try:
    response = agent.run(user_input)
except InputCheckError as e:
    # Don't just fail - help the user
    return "I can't process that request. Please rephrase without harmful content."
```

### 4. Test Your Guardrails

```python
import pytest
from agno.exceptions import InputCheckError

def test_violence_blocked():
    agent = Agent(
        model=OpenAIChat(id="gpt-4o"),
        pre_hooks=[OpenAIModerationGuardrail()],
    )

    with pytest.raises(InputCheckError):
        agent.run("violent content here")

def test_safe_content_passes():
    agent = Agent(
        model=OpenAIChat(id="gpt-4o"),
        pre_hooks=[OpenAIModerationGuardrail()],
    )

    # Should not raise
    response = agent.run("What is machine learning?")
    assert response is not None
```

## References

- [Agno Guardrails Documentation](https://docs.agno.com/basics/guardrails/overview)
- [OpenAI Moderation API](https://platform.openai.com/docs/guides/moderation)
- [Content Policy Categories](https://platform.openai.com/docs/guides/moderation/overview)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/testing (Testing)/testing|testing]]

# NATS Handlers

Agent-specific NATS message handlers.
Each agent can register custom handlers for its domain subjects.

## Handler Interface

```python
from hermes_tools import terminal

def register_nats_handler(subject, handler_func):
    '''Register a custom message handler for a NATS subject.'''
    pass

def on_legacy_command(message):
    '''Handle command from Legacy.'''
    return f"Acknowledged: {message}"

def on_sensor_data(data):
    '''Handle incoming sensor telemetry.'''
    pass
```

## Agent Handler Lists
To be populated as agents are implemented.

# Workflow Routing Doctrine

OneMind should not wake Legacy for every event. Routing rules decide what deserves AI, alerts, tasks, missions, or commands.

## Pattern

```text
event/telemetry/state
  -> rule match
  -> command / task / mission / alert / Legacy prompt
  -> result event
  -> UI projection
```

## Workflow record

```json
{
  "id": "rule.hostile-drone-critical",
  "enabled": true,
  "description": "Escalate critical hostile air tracks",
  "when": {
    "subject": "ent.air.*.state",
    "where": "data.affiliation == 'hostile' && data.priority == 'critical'"
  },
  "then": [
    {"publish": "ui.alerts.active"},
    {"publish": "cmd.mission.create"},
    {"prompt": "agents.prompt.hermes.zeus.ops"},
    {"publish": "voice.say"}
  ],
  "approval": "operator_required_for_effectors"
}
```

## Required controls

- enable/disable workflow
- dry-run/test workflow
- inspect last runs
- view events that matched
- view commands produced
- require approval for dangerous commands
- route to Legacy only after threshold or explicit rule

## Protected command classes

Commands that require approval by default:

- `cmd.drone.launch`
- `cmd.drone.land`
- `cmd.robot.move`
- `cmd.intercept.*`
- `cmd.service.restart.*`
- `cmd.nex.workload.stop.*`
- `cmd.integration.*.write`
- `cmd.all.abort` may bypass approval as emergency reflex, but must audit.

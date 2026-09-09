# TAK Projection Doctrine

TAK is the primary tactical/field map. OneMind should not overload TAK with every backend detail by default. Instead, `om-tak-projector` decides which entities become CoT markers.

## Rule

An entity appears in TAK when:

1. it has a valid entity record,
2. it has location or a projectable synthetic location,
3. `visibility.tak == true`,
4. command policy allows projection,
5. the projector can map its kind to a CoT type.

## Commands

```text
cmd.entity.{id}.tak.show
cmd.entity.{id}.tak.hide
cmd.entity.{id}.tak.pin
cmd.entity.{id}.tak.unpin
cmd.layer.{layer}.tak.show
cmd.layer.{layer}.tak.hide
cmd.source.{source}.tak.project_on
cmd.source.{source}.tak.project_off
```

## Events

```text
evt.entity.tak.visible
evt.entity.tak.hidden
evt.layer.tak.visible
evt.layer.tak.hidden
evt.source.tak.projected
evt.source.tak.stopped
```

## CoT mapping starter set

| Entity kind | CoT type | Notes |
|---|---|---|
| hostile air | `a-h-A-M-F` | hostile drone/aircraft |
| friendly drone | `a-f-A-M-F` | UAV / multirotor |
| TAK user | `a-f-G-U-C` | team/operator |
| ground robot | `a-f-G-U-R` | Unitree/UGV |
| machine/node | `a-f-G-U-i-I` | server/Pi/cloud VM infrastructure |
| AI agent | `a-f-G-U-i-ai` | custom OneMind detail marker |
| sensor | `b-m-p-s` | sensor/broadcast point |
| ship | `a-n-S-X-M` | AIS neutral surface vessel |
| aircraft | `a-n-A-C-F` | ADS-B neutral aircraft |

## UI behavior

`omos-atoc` should expose per-entity controls:

- Show in TAK
- Hide from TAK
- Pin
- Mute
- Watch with Legacy
- Open in TAK
- Ask Legacy
- Open event trace

The button publishes a command. The projector updates registry/projection state and emits an event. The UI updates from `ui.assets.{id}.state` or `ui.layers.tak`.

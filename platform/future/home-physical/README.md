# Future home and physical integration profile

User-authorized preparation for the later home cluster. These resources are outside
the active Argo tree and have no automated sync. They do not allocate cloud nodes,
attach devices or operate physical equipment when merged.

## Included GitOps entry points

- `../home-rook/application.yaml`: pinned Rook Helm operator. A later site overlay
  declares Ceph disks/failure domains and pools; no automatic disk discovery here.
- `../home-envoy/application.yaml`: pinned Envoy Gateway Helm application.
- `../home-forgejo/application.yaml`: pinned Forgejo Helm application.
- `../home-akri/application.yaml`: pinned Akri Helm application, with discovery
  protocols/device classes selected by the site's activation PR.
- `../home-tekton/application.yaml`: pinned Tekton Helm application for local builds.
- `chart/`: OneMind Helm packaging for Home Assistant, Frigate and future physical
  adapters. All components disabled until the home overlay supplies configuration.

These future Applications target the registered cluster name `onemind-home`; that
name must be bound to the actual home API endpoint before creating any Application.
Do not repoint them to today's cloud cluster.

## Home profile activation

Provision k3s on selected home machines and record its configuration in the site
infrastructure repository. Hardware/OS/network details are needed before writing
the final bootstrap; Helm cannot provision the Kubernetes control plane itself.
Register the home cluster with Argo, establish local identity/DNS/artifact access,
and create home-specific overlays. Keep essential site operation local during WAN
loss. Do not stretch the cloud etcd or a storage quorum across the WAN.

Choose explicit unused storage devices for Ceph, backup destinations and a storage
class. Establish retained backups before placing authoritative home data in the
pool. Review radio/USB/video device access for adapters and persist their app
configuration. Home Assistant Container is not Supervisor/Home Assistant OS.

Home Assistant and Frigate integrate household equipment and camera observations.
Device-native MQTT stays local with deliberate NATS bridges. No unconfigured
automation should be allowed to issue physical commands during installation.

## Robotics and physical AI

Carry these roles into hardware-specific follow-up changes:

| Component | Site role |
|---|---|
| ROS 2 | Robotics components and companion communication |
| Nav2 | Ground navigation/localization integration |
| Gazebo | Repeatable simulated missions and failure scenarios |
| PX4 or ArduPilot | Vehicle-specific onboard autopilot; outside generic K8s control |
| Open-RMF | Compatible fleet/building coordination with site adapters |
| TAK/CoT bridge | Existing field interface and geo-envelope integration |
| BuildKit | Local builds coordinated by Tekton |

Select compatible releases, maps, coordinate frames, device interfaces and vehicle
configuration together. Companion services can use the Helm chart's `robotics`
component once an immutable adapter image and its configuration are supplied.
The chart deliberately has no default robotics image or motor-control command.

GPU serving and fleet governance belong to the separate home GPU/fleet PR. This
profile has no GPU request. Activation requires a later home-targeted GitOps PR,
including local operation, backup/restore and device takeover verification.

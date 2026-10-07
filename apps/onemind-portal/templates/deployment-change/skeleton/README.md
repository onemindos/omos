# ${{ values.name }}

## Purpose
${{ values.reason }}

## Accountable reviewer
${{ values.reviewer }}

This is a review package. It does not execute a change. Never add credentials,
private backup contents, or temporary signed URLs to this directory.

## Proposed change
- Target: ${{ values.entity }}
- Image: ${{ values.image }}
- Window: ${{ values.maintenance }}

## Rollback approach
${{ values.rollback }}

## Review and promotion checklist
- [ ] Locate the target's actual GitOps owner and deployment manifest.
- [ ] Verify image provenance, CPU architecture, migrations and compatibility.
- [ ] Record current digest, capacity, dependency impact, health criteria and downtime.
- [ ] Confirm storage backup and restore evidence if state or schema may change.
- [ ] Confirm rollback remains possible after any migration.
- [ ] Prepare a separate minimal deployment PR against the owning manifest.
- [ ] Review that PR and release through Argo; do not apply branch manifests directly.
- [ ] Record deployed digest, synchronization, health and observed outcome.

This package deliberately contains no Kubernetes patch or apply command. A catalog
entity selection is not authorization to mutate its deployment.

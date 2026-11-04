# RWO PVC Deadlock: Stale ReplicaSet Pattern

## What happens

A rollout creates a new ReplicaSet while an old pod still holds a ReadWriteOnce
claim (often because the old pod is crashlooping). The new pod stays in
`ContainerCreating` or `Init:0/1` with `Multi-Attach error for volume ... already
used by pod(s) <old pod>`. Deleting only the old pod does not help: its ReplicaSet
recreates it within seconds and re-takes the volume.

## Diagnose

```bash
kubectl get rs -n <ns> -l app=<name> --sort-by=.metadata.creationTimestamp   # oldest is stale
kubectl describe pod -n <ns> <new-pod> | grep -A5 Events
```

## Fix

Delete the stale ReplicaSet (it removes its pods and releases the volume), then the
new pod attaches. This is a cluster mutation: do it only with the owner's approval,
and where possible fix it in Git (below) so it cannot recur.

## Prevention

- Single-replica Deployments that mount an RWO claim: set `strategy: {type: Recreate}`.
- Two containers that must share one RWO claim belong in one pod (shared network
  namespace, `localhost` between them). Two Deployments on one RWO claim deadlock.

## Stale VolumeAttachment

A force-deleted pod (`--force --grace-period=0`) holding an RWO claim can leave a
stale `VolumeAttachment`, and the next pod on another node gets the same
Multi-Attach error. List with `kubectl get volumeattachment | grep <volume-handle>`.
Deleting the attachment (with approval) frees it; otherwise the attach controller
cache expires on its own in about six minutes.

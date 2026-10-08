# Future home GPU serving

This is an opt-in, dormant vLLM runtime for the future home Kubernetes cluster.
The user has no GPU available now. Today's CPU-only cluster continues with model
routing and available provider endpoints; this change does not request cloud GPUs.

This directory is **not referenced by any Argo Application, bootstrap script or
parent Kustomization**. Merging it does not create workloads, volumes or namespaces.
The Deployment also defaults to zero replicas. Do not apply this base to today's
cloud cluster: even at zero replicas, its PVC would request storage.

## Included runtime

- Upstream vLLM v0.31.0 image pinned by verified multi-platform digest.
- One NVIDIA GPU per replica; explicit `inference` node selection.
- Private ClusterIP endpoint named `onemind-local-chat`, port 8000.
- Access only from the designated LiteLLM pods in `onemind-ai`, enforced by a
  NetworkPolicy-capable CNI. No Ingress or public load balancer.
- Required API-key Secret reference, no committed credential.
- Locally staged model files, offline model loading and denied runtime egress.
- Read-only model mount, bounded scratch/shared-memory volumes and health probes.
- Recreate updates avoid requiring two GPUs during a single-GPU rollout.

The initial model directory is `/models/Qwen3-8B`, exposed as the stable client
alias `onemind-local-chat`. This is a prototype candidate, not a hardware sizing
claim. Model weights, KV cache, context length and concurrency determine VRAM use;
select the artifact/precision and resource budget when actual hardware is known.
The draft has not been run on a GPU, so non-root/read-only image compatibility and
inference performance still need runtime verification.

## Activation in a later home-cluster PR

1. Record the home cluster context, GPU model/VRAM, driver/runtime and storage
   class. Install the appropriate NVIDIA device plugin or GPU Operator once.
2. Confirm the worker advertises `nvidia.com/gpu`, and deliberately label the
   chosen node `onemind.io/compute-role=inference`.
3. Create a home overlay with the selected storage class, model path, sizing and
   image choice. Add the overlay to a home-specific Argo Application only then.
   Do not add a cloud-targeting Application to `platform/argocd/`.
4. Stage the chosen model snapshot into the PVC, recording its repository,
   immutable revision, artifact checksums and license. Preserve local copies of
   model files and the container image for disconnected restarts. The runtime
   does not download weights. Qwen3-8B is an example until that snapshot is chosen.
5. Supply `onemind-local-chat-auth` through the home's secret workflow, with key
   `api-key`. Configure the gateway to use that key; no literal token in Git.
6. Route the existing LiteLLM logical model to
   `http://onemind-local-chat.onemind-inference.svc.cluster.local:8000/v1` when
   LiteLLM is in the same home cluster. For cloud-to-home use, add an explicitly
   authenticated private mesh gateway; cloud DNS cannot resolve this home service.
7. Enable one replica and verify load, cancellation, authorization, non-gateway
   denial, metrics, restart from local artifacts and operation without WAN.

vLLM's API key does not authenticate every server endpoint. The network boundary
and gateway remain required; do not expose the raw runtime directly. Health probes
use the internal endpoint. No host paths, Kubernetes API permissions or NATS
credentials are needed by this serving runtime.

KServe's native Kubernetes mode is the intended later model-lifecycle layer. A
follow-up PR can replace this Deployment with its managed equivalent once the home
CRDs/controller exist. Exactly one controller should own the serving workload;
do not deploy both forms for the same endpoint.

## Render without deployment

From the repository root:

```bash
kubectl kustomize platform/future/home-gpu-serving
```

Rendering does not apply any resources. No activation command is part of the
current cluster deployment path.

References: [vLLM Kubernetes](https://docs.vllm.ai/en/stable/deployment/k8s/),
[vLLM endpoint authentication](https://docs.vllm.ai/en/stable/serving/online_serving/openai_compatible_server/),
[Qwen3-8B model card](https://huggingface.co/Qwen/Qwen3-8B).

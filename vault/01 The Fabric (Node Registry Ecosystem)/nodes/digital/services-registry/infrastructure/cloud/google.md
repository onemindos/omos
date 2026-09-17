---
title: "Google Cloud (GCP)"
type: reference
status: planned
tags: [platform, cloud, provider, google]
---
# Google Cloud (GCP)

> Google's cloud platform. Strong in AI/ML (Vertex AI, TPUs), data analytics (BigQuery), and Kubernetes (GKE — they created K8s).

## Use Case in OneMind OS
- Gemini API (AI inference)
- TPU access for ML training
- BigQuery for large-scale analytics
- GKE for container orchestration
- Cloud Run for serverless

## CLI
```bash
gcloud auth login
gcloud config set project PROJECT_ID
gcloud compute instances list
gcloud run deploy
gcloud ai models list
```

## Key Services
| Service | Equivalent |
|---------|-----------|
| Compute Engine | EC2 / VMs |
| Cloud Run | Lambda + Fargate (serverless containers) |
| GKE | EKS / Kubernetes |
| Vertex AI | SageMaker / ML platform |
| BigQuery | Redshift / Analytics |
| Cloud Storage | S3 / Object storage |
| Pub/Sub | SNS+SQS / Event streaming |
| Artifact Registry | ECR / Container registry |

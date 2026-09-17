# Microsoft Azure

> **Cloud Platform for AI and Enterprise**

## Overview

Microsoft Azure provides cloud infrastructure, primarily used for Claude access via AWS Bedrock integration and backup cloud services.

## Status

- **Active**: ✅ Yes
- **Primary Region**: East US
- **Last Updated**: January 25, 2026

## Services in Use

### AI/ML

| Service                | Purpose             | Status     |
| ---------------------- | ------------------- | ---------- |
| **Azure OpenAI**       | GPT models (backup) | 🔄 Planned |
| **Cognitive Services** | Speech, vision      | 🔄 Planned |

### Compute

| Service                  | Purpose               | Status     |
| ------------------------ | --------------------- | ---------- |
| **Azure Container Apps** | Serverless containers | 🔄 Planned |
| **Functions**            | Serverless functions  | 🔄 Planned |

### Storage

| Service           | Purpose      | Status     |
| ----------------- | ------------ | ---------- |
| **Blob Storage**  | File storage | 🔄 Planned |
| **Table Storage** | NoSQL data   | 🔄 Planned |

### Identity

| Service      | Purpose             | Status    |
| ------------ | ------------------- | --------- |
| **Entra ID** | Identity management | ✅ Active |

## Azure OpenAI

### Setup

```python
from openai import AzureOpenAI

client = AzureOpenAI(
    api_key=os.getenv("AZURE_OPENAI_API_KEY"),
    api_version="2024-02-01",
    azure_endpoint=os.getenv("AZURE_OPENAI_ENDPOINT")
)
```

### Chat Completion

```python
response = client.chat.completions.create(
    model="gpt-4",  # Deployment name
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Hello!"}
    ],
    temperature=0.7,
    max_tokens=800
)

print(response.choices[0].message.content)
```

### Streaming

```python
stream = client.chat.completions.create(
    model="gpt-4",
    messages=[{"role": "user", "content": "Tell me a story"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

## Azure SDK

### Installation

```bash
pip install azure-identity azure-storage-blob azure-mgmt-resource
```

### Authentication

```python
from azure.identity import DefaultAzureCredential

# Uses environment variables, managed identity, or Azure CLI
credential = DefaultAzureCredential()
```

## Blob Storage

### Upload/Download

```python
from azure.storage.blob import BlobServiceClient

blob_service = BlobServiceClient(
    account_url="https://account.blob.core.windows.net",
    credential=credential
)

container = blob_service.get_container_client("container-name")

# Upload
with open("file.txt", "rb") as f:
    container.upload_blob("path/to/file.txt", f)

# Download
blob = container.download_blob("path/to/file.txt")
content = blob.readall()

# List blobs
for blob in container.list_blobs():
    print(blob.name)
```

### Generate SAS URL

```python
from azure.storage.blob import generate_blob_sas, BlobSasPermissions
from datetime import datetime, timedelta

sas_token = generate_blob_sas(
    account_name="account",
    container_name="container",
    blob_name="file.txt",
    account_key="key",
    permission=BlobSasPermissions(read=True),
    expiry=datetime.utcnow() + timedelta(hours=1)
)

url = f"https://account.blob.core.windows.net/container/file.txt?{sas_token}"
```

## Azure Functions

### Function Example

```python
import azure.functions as func
import json

def main(req: func.HttpRequest) -> func.HttpResponse:
    body = req.get_json()

    result = process(body)

    return func.HttpResponse(
        json.dumps(result),
        mimetype="application/json",
        status_code=200
    )
```

### function.json

```json
{
  "bindings": [
    {
      "authLevel": "function",
      "type": "httpTrigger",
      "direction": "in",
      "name": "req",
      "methods": ["post"]
    },
    {
      "type": "http",
      "direction": "out",
      "name": "$return"
    }
  ]
}
```

## Container Apps

### Deployment

```yaml
# container-app.yaml
name: legacy-agent
properties:
  configuration:
    ingress:
      external: true
      targetPort: 8080
    secrets:
      - name: db-connection
        value: "..."
  template:
    containers:
      - name: legacy-agent
        image: myregistry.azurecr.io/legacy-agent:latest
        resources:
          cpu: 1.0
          memory: 2Gi
        env:
          - name: DATABASE_URL
            secretRef: db-connection
```

## Key Vault

### Store Secrets

```python
from azure.keyvault.secrets import SecretClient

vault_url = "https://onemind-vault.vault.azure.net"
client = SecretClient(vault_url=vault_url, credential=credential)

# Set secret
client.set_secret("api-key", "secret-value")

# Get secret
secret = client.get_secret("api-key")
print(secret.value)
```

## Azure CLI

### Common Commands

```bash
# Login
az login

# Set subscription
az account set --subscription "subscription-id"

# List resources
az resource list --output table

# Create resource group
az group create --name onemind-rg --location eastus

# Deploy container
az containerapp create \
  --name legacy-agent \
  --resource-group onemind-rg \
  --environment onemind-env \
  --image myregistry.azurecr.io/legacy-agent:latest \
  --target-port 8080 \
  --ingress external
```

## Environment Variables

```env
AZURE_SUBSCRIPTION_ID=your_subscription_id
AZURE_TENANT_ID=your_tenant_id
AZURE_CLIENT_ID=your_client_id
AZURE_CLIENT_SECRET=your_client_secret
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
AZURE_OPENAI_API_KEY=your_api_key
```

## Cost Management

### Budget Alerts

1. Go to Cost Management + Billing
2. Create Budget
3. Set threshold alerts (75%, 90%, 100%)
4. Configure email notifications

### Cost Optimization

- Use Reserved Instances for VMs
- Choose appropriate service tiers
- Use Azure Advisor recommendations
- Enable auto-shutdown for dev resources

## vs AWS

| Feature    | Azure          | AWS         |
| ---------- | -------------- | ----------- |
| AI Models  | Azure OpenAI   | Bedrock     |
| Containers | Container Apps | ECS/Fargate |
| Functions  | Functions      | Lambda      |
| Storage    | Blob           | S3          |
| Identity   | Entra ID       | IAM         |
| CLI        | az             | aws         |

## Best Practices

- ✅ Use managed identities for authentication
- ✅ Enable Azure Security Center
- ✅ Use resource groups for organization
- ✅ Enable diagnostic logging
- ✅ Set up cost alerts
- ✅ Use Azure Policy for governance

## Resources

- [Azure Documentation](https://docs.microsoft.com/azure/)
- [Azure OpenAI Service](https://docs.microsoft.com/azure/cognitive-services/openai/)
- [Azure SDK for Python](https://docs.microsoft.com/azure/developer/python/)

---

**Related**: [AWS](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/platforms/cloud/AWS.md) | [Agno](./Agno.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

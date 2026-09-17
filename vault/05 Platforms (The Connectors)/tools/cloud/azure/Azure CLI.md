---
title: "CLI — az (Azure CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/azure, topic/cloud]
ai_summary: "Azure CLI (az) quick reference — login, resource groups, VMs, storage, functions, and key vault access for One Mind cloud resources."
---

# az — Azure CLI

> Microsoft Azure cloud management. Install: `brew install azure-cli` / `pip install azure-cli`.

---

## Auth

```bash
# Login (browser)
az login

# Login with service principal
az login --service-principal -u <app-id> -p <password> --tenant <tenant-id>

# Check current account
az account show

# List subscriptions
az account list --output table

# Set active subscription
az account set --subscription "Subscription Name"

# Logout
az logout
```

---

## Resource Groups

```bash
# List resource groups
az group list --output table

# Create resource group
az group create --name onemind-rg --location eastus

# Delete resource group (and all resources in it — DESTRUCTIVE)
# az group delete --name onemind-rg --yes
```

---

## Virtual Machines

```bash
# List VMs
az vm list --output table

# Start / stop / restart VM
az vm start --resource-group onemind-rg --name vm-name
az vm stop --resource-group onemind-rg --name vm-name
az vm restart --resource-group onemind-rg --name vm-name

# Get VM status
az vm show --resource-group onemind-rg --name vm-name --show-details

# SSH into VM (using az ssh)
az ssh vm --resource-group onemind-rg --name vm-name
```

---

## Storage

```bash
# List storage accounts
az storage account list --output table

# List containers in storage account
az storage container list --account-name mystorageaccount

# Upload a file
az storage blob upload \
  --account-name mystorageaccount \
  --container-name mycontainer \
  --name filename.txt \
  --file ./local/filename.txt

# Download a file
az storage blob download \
  --account-name mystorageaccount \
  --container-name mycontainer \
  --name filename.txt \
  --file ./local/filename.txt
```

---

## Azure Functions

```bash
# List function apps
az functionapp list --output table

# Deploy a function app
az functionapp deployment source config-zip \
  --resource-group onemind-rg \
  --name function-app-name \
  --src ./deploy.zip

# View logs
az functionapp log tail --resource-group onemind-rg --name function-app-name
```

---

## Key Vault

```bash
# List key vaults
az keyvault list --output table

# Get a secret
az keyvault secret show --vault-name my-vault --name secret-name --query value -o tsv

# Set a secret
az keyvault secret set --vault-name my-vault --name secret-name --value "value"

# List secrets
az keyvault secret list --vault-name my-vault --output table
```

---

## Container Registry (ACR)

```bash
# Login to ACR
az acr login --name myregistry

# List repositories
az acr repository list --name myregistry

# Build and push image
az acr build --registry myregistry --image image-name:tag .
```

---

## Useful Flags

```bash
--output table    # human-readable table
--output json     # full JSON
--output tsv      # tab-separated values (good for scripting)
--query "..."     # JMESPath filter
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

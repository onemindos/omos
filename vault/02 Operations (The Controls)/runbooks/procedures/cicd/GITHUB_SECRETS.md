# GitHub Secrets Configuration Guide

This guide explains how to configure GitHub repository secrets for the CI/CD pipeline after migrating from AWS Lightsail to Azure.

---

## Required Secrets

The CI/CD pipeline (`.github/workflows/deploy.yml`) requires these 3 secrets:

| Secret Name | Description | How to Get |
|-------------|-------------|------------|
| `AZURE_SSH_PRIVATE_KEY` | SSH private key for authenticating to Azure VM | See below |
| `AZURE_HOST` | Azure VM public IP address | `20.121.38.186` |
| `AZURE_USER` | SSH username for Azure VM | `zeus` |

---

## Step-by-Step Setup

### 1. Get Your SSH Private Key

On your **local machine**, run:

```bash
cat ~/.ssh/id_ed25519
```

**Output will look like**:
```
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
[... many lines of base64-encoded key ...]
-----END OPENSSH PRIVATE KEY-----
```

**Copy the entire output** (including the `-----BEGIN` and `-----END` lines).

> **Security Warning**: This is your private key. Never share it publicly or commit it to git.

---

### 2. Add Secrets to GitHub Repository

#### Via GitHub Web UI:

1. Navigate to your repository: https://github.com/YOUR_USERNAME/onemind-os
2. Click **Settings** (top right)
3. In left sidebar, click **Secrets and variables** → **Actions**
4. Click **New repository secret**

#### Add these 3 secrets:

**Secret 1: AZURE_SSH_PRIVATE_KEY**
- **Name**: `AZURE_SSH_PRIVATE_KEY`
- **Value**: Paste the entire SSH private key (from step 1)
- Click **Add secret**

**Secret 2: AZURE_HOST**
- **Name**: `AZURE_HOST`
- **Value**: `20.121.38.186`
- Click **Add secret**

**Secret 3: AZURE_USER**
- **Name**: `AZURE_USER`
- **Value**: `zeus`
- Click **Add secret**

---

### 3. Remove Old AWS Secrets (Optional)

If you have old AWS secrets, you can delete them:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Click **Remove** next to:
   - `AWS_SSH_PRIVATE_KEY`
   - `AWS_HOST`
   - `AWS_USER`

These are no longer used after the Azure migration.

---

## Verification

After adding secrets, verify they're configured:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. You should see these 3 secrets listed:
   - ✅ AZURE_SSH_PRIVATE_KEY
   - ✅ AZURE_HOST
   - ✅ AZURE_USER

**Note**: GitHub won't show the secret values (for security), but you'll see the names.

---

## Testing the Pipeline

### Trigger a Deployment

**Option 1: Push to main branch**
```bash
git add .
git commit -m "test: trigger deployment"
git push origin main
```

**Option 2: Manual workflow dispatch**
1. Go to **Actions** tab in GitHub
2. Select **Deploy** workflow
3. Click **Run workflow** → **Run workflow**

### Monitor Deployment

1. Go to **Actions** tab
2. Click on the running workflow
3. Click on **Deploy to Production** job
4. Watch the logs for:
   - ✅ SSH connection successful
   - ✅ Code pulled
   - ✅ Containers built
   - ✅ Health checks passed

---

## Troubleshooting

### Error: "Permission denied (publickey)"

**Cause**: SSH private key is incorrect or not matching the public key on the VM.

**Fix**:
1. Verify your local SSH key works:
   ```bash
   ssh zeus@20.121.38.186 "echo 'SSH works'"
   ```
2. If that works, re-copy the private key to GitHub secrets (step 1-2 above)

### Error: "Host key verification failed"

**Cause**: Azure VM not in known_hosts.

**Fix**: The deploy script auto-handles this with `ssh-keyscan`. If it persists, manually add:
```bash
ssh-keyscan -H 20.121.38.186 >> ~/.ssh/known_hosts
```

Then re-add the output to GitHub secrets.

### Error: "docker compose: command not found"

**Cause**: Docker not installed on Azure VM.

**Fix**: SSH to VM and install Docker:
```bash
ssh zeus@20.121.38.186
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker zeus
```

Then logout and back in to apply group membership.

### Error: "no configuration file provided: not found"

**Cause**: Docker Compose file path is incorrect.

**Fix**: Verify the file exists on the VM:
```bash
ssh zeus@20.121.38.186 "ls -la ~/OneMind-OS/infra/deployment/docker/compose/"
```

You should see `docker-compose.core.yml`.

---

## Security Best Practices

### 1. Use Deploy Keys (Recommended for Production)

Instead of using your personal SSH key, create a **deploy-specific key**:

```bash
# On your local machine
ssh-keygen -t ed25519 -C "github-deploy@onemind-os" -f ~/.ssh/onemind_deploy_key

# Copy public key to Azure VM
ssh-copy-id -i ~/.ssh/onemind_deploy_key.pub zeus@20.121.38.186

# Test deploy key
ssh -i ~/.ssh/onemind_deploy_key zeus@20.121.38.186 "echo 'Deploy key works'"

# Add private key to GitHub secrets
cat ~/.ssh/onemind_deploy_key
# Copy output to AZURE_SSH_PRIVATE_KEY secret
```

**Benefits**:
- Separate from your personal SSH key
- Can be revoked without affecting your access
- Can have restricted permissions on VM

### 2. Restrict Secret Access

In GitHub, you can:
- Limit secrets to specific environments (production only)
- Require approval for workflow runs
- Enable branch protection rules

### 3. Rotate Keys Regularly

Best practice: Rotate SSH keys every 90 days:
1. Generate new SSH key pair
2. Add new public key to VM
3. Update GitHub secret with new private key
4. Remove old public key from VM

---

## Alternative: Using Tailscale for Private Deployments

For added security, you can deploy over Tailscale VPN:

### 1. Use Tailscale IP in AZURE_HOST

Instead of `20.121.38.186`, use the Tailscale IP:

```bash
# Get Tailscale IP
ssh zeus@20.121.38.186 "tailscale ip -4"
# Example output: 100.101.102.103
```

Update `AZURE_HOST` secret to use `100.101.102.103`

### 2. Add Tailscale to GitHub Actions

Add this step to `.github/workflows/deploy.yml` before "Configure SSH":

```yaml
      - name: Connect to Tailscale
        uses: tailscale/github-action@v2
        with:
          authkey: ${{ secrets.TAILSCALE_AUTHKEY }}
          version: 1.56.1
```

**Benefits**:
- No public IP exposure
- Encrypted end-to-end
- No need to open SSH port (22) publicly

---

## Summary

**Required GitHub Secrets**:
1. **AZURE_SSH_PRIVATE_KEY** - Your SSH private key (from `~/.ssh/id_ed25519`)
2. **AZURE_HOST** - `20.121.38.186`
3. **AZURE_USER** - `zeus`

**Setup Time**: 5 minutes

**Next Steps**:
1. Add the 3 secrets to GitHub (Settings → Secrets)
2. Push to main or manually trigger deployment
3. Monitor Actions tab for successful deployment

---

**Last Updated**: 2026-01-27
**Azure VM IP**: 20.121.38.186
**Related Docs**:
- [CI/CD Workflows](../.github/workflows/)
- [Infrastructure Guide](INFRASTRUCTURE.md)
- [Deployment Guide](./deployment/)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

# OneMind OS: Cloudflare DNS Setup

Complete guide for configuring Cloudflare DNS for OneMind OS on Azure.

---

## Prerequisites

- Cloudflare account with a domain
- Domain nameservers pointed to Cloudflare
- Azure VM running at: **20.121.38.186**

---

## DNS Records to Create

### Option A: Subdomain Setup (Recommended)

Create these DNS records in your Cloudflare dashboard:

| Type | Name | Content | Proxy | TTL | Purpose |
|------|------|---------|-------|-----|---------|
| A | `api` | `20.121.38.186` | ✅ Proxied | Auto | Backend API |
| A | `app` | `20.121.38.186` | ✅ Proxied | Auto | Frontend UI |
| A | `admin` | `20.121.38.186` | ✅ Proxied | Auto | Portainer (optional) |

**Result**:
- Backend: https://api.yourdomain.com
- Frontend: https://app.yourdomain.com
- Admin: https://admin.yourdomain.com

---

### Option B: Root Domain Setup

| Type | Name | Content | Proxy | TTL | Purpose |
|------|------|---------|-------|-----|---------|
| A | `@` | `20.121.38.186` | ✅ Proxied | Auto | Root domain |
| CNAME | `www` | `yourdomain.com` | ✅ Proxied | Auto | WWW redirect |

**Result**:
- https://yourdomain.com (frontend)
- https://yourdomain.com/api (backend with reverse proxy)

---

## Step-by-Step Setup

### 1. Add DNS Records (Cloudflare Dashboard)

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Select your domain
3. Click **DNS** > **Records**
4. Click **Add record**

**For Backend (api subdomain)**:
- Type: `A`
- Name: `api`
- IPv4 address: `20.121.38.186`
- Proxy status: **Proxied** (orange cloud)
- TTL: Auto
- Click **Save**

**For Frontend (app subdomain)**:
- Type: `A`
- Name: `app`
- IPv4 address: `20.121.38.186`
- Proxy status: **Proxied** (orange cloud)
- TTL: Auto
- Click **Save**

---

### 2. Install Caddy (Reverse Proxy + Auto HTTPS)

SSH to Azure VM and install Caddy:

```bash
ssh zeus@20.121.38.186

# Install Caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy
```

---

### 3. Configure Caddy Reverse Proxy

Create Caddyfile:

```bash
sudo nano /etc/caddy/Caddyfile
```

**Paste this configuration** (replace `yourdomain.com`):

```
# Backend API
api.yourdomain.com {
    reverse_proxy localhost:8080

    # Enable compression
    encode gzip

    # Headers
    header {
        # Security headers
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        Referrer-Policy "strict-origin-when-cross-origin"

        # CORS (if needed)
        Access-Control-Allow-Origin "*"
        Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS"
        Access-Control-Allow-Headers "Content-Type, Authorization"
    }

    # Logs
    log {
        output file /var/log/caddy/api.log
        format json
    }
}

# Frontend UI
app.yourdomain.com {
    reverse_proxy localhost:3000

    encode gzip

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        Referrer-Policy "strict-origin-when-cross-origin"
    }

    log {
        output file /var/log/caddy/app.log
        format json
    }
}

# Admin Panel (Portainer) - Optional
admin.yourdomain.com {
    reverse_proxy localhost:9000

    encode gzip

    # Restrict access (optional)
    # basicauth {
    #     admin $2a$14$hashed_password_here
    # }

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
    }

    log {
        output file /var/log/caddy/admin.log
        format json
    }
}
```

---

### 4. Create Log Directory

```bash
sudo mkdir -p /var/log/caddy
sudo chown caddy:caddy /var/log/caddy
```

---

### 5. Test and Start Caddy

```bash
# Test configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# Reload Caddy (or start if not running)
sudo systemctl reload caddy

# Enable auto-start
sudo systemctl enable caddy

# Check status
sudo systemctl status caddy
```

---

### 6. Update Azure NSG (Allow HTTP/HTTPS)

Ports 80 and 443 should already be open. Verify:

```bash
az network nsg rule list --resource-group onemind-production \
  --nsg-name onemind-vmNSG \
  --output table | grep -E "80|443"
```

If not open, add rules:

```bash
# HTTP (port 80)
az vm open-port --resource-group onemind-production \
  --name onemind-vm --port 80 --priority 1100

# HTTPS (port 443)
az vm open-port --resource-group onemind-production \
  --name onemind-vm --port 443 --priority 1200
```

---

### 7. Configure Cloudflare SSL/TLS

1. Go to Cloudflare Dashboard > **SSL/TLS**
2. Set mode to: **Full (strict)**
   - This enables end-to-end encryption with Caddy's auto-generated certificates
3. Enable **Always Use HTTPS**
4. Enable **Automatic HTTPS Rewrites**

---

### 8. Update Frontend Environment

Update frontend to use new API domain:

**On your local machine**:

```bash
# Edit frontend .env
nano frontend/.env

# Update API URL
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

Rebuild and deploy frontend:

```bash
cd frontend
flutter build web
# Deploy updated build to VM
scp -r build/web/* zeus@20.121.38.186:~/OneMind-OS/frontend/build/web/
```

---

## Verification

### Test DNS Resolution

```bash
# Check if DNS is propagated
dig api.yourdomain.com +short
dig app.yourdomain.com +short

# Should return Cloudflare proxy IPs (not 20.121.38.186 directly)
```

### Test HTTPS Access

```bash
# Backend
curl -I https://api.yourdomain.com/health

# Should return:
# HTTP/2 200
# server: caddy

# Frontend
curl -I https://app.yourdomain.com

# Should return:
# HTTP/2 200
```

### Test from Browser

1. **Backend**: https://api.yourdomain.com/health
   - Should show: `{"status":"ok","instantiated_at":"..."}`

2. **Frontend**: https://app.yourdomain.com
   - Should load OneMind OS UI

3. **API Endpoints**: https://api.yourdomain.com/agents
   - Should return JSON array of agents

---

## Cloudflare Security Settings (Recommended)

### 1. Enable Bot Protection

- Dashboard > **Security** > **Bots**
- Enable **Bot Fight Mode** (free) or **Super Bot Fight Mode** (paid)

### 2. Configure Firewall Rules

Create a rule to allow only legitimate traffic:

- Dashboard > **Security** > **WAF** > **Firewall rules**
- Add rule: Block known bad bots, allow good bots

### 3. Enable Rate Limiting

- Dashboard > **Security** > **Rate Limiting**
- Create rule:
  - If requests > 100 per minute from same IP
  - Then challenge or block

### 4. Page Rules

Create page rules for caching:

- Dashboard > **Rules** > **Page Rules**

**API Caching** (optional - be careful with dynamic content):
```
URL: api.yourdomain.com/*
Cache Level: Standard
Edge Cache TTL: 2 hours (only for GET requests to static endpoints)
```

---

## Monitoring

### Caddy Logs

```bash
# Backend logs
sudo tail -f /var/log/caddy/api.log

# Frontend logs
sudo tail -f /var/log/caddy/app.log

# System logs
sudo journalctl -u caddy -f
```

### Cloudflare Analytics

- Dashboard > **Analytics & Logs**
- View traffic, requests, bandwidth, threats blocked

---

## Troubleshooting

### Issue: "Too many redirects"

**Cause**: Cloudflare SSL mode mismatch
**Fix**: Set SSL/TLS mode to **Full (strict)**

### Issue: "502 Bad Gateway"

**Cause**: Backend not running or Caddy misconfigured
**Fix**:
```bash
# Check backend
docker ps | grep agno

# Check Caddy
sudo systemctl status caddy
sudo caddy validate --config /etc/caddy/Caddyfile
```

### Issue: "Certificate error"

**Cause**: Caddy hasn't obtained Let's Encrypt certificate yet
**Fix**: Wait 1-2 minutes for Caddy to auto-generate certificate. Check logs:
```bash
sudo journalctl -u caddy -f
```

### Issue: "DNS not resolving"

**Cause**: DNS propagation delay
**Fix**: Wait 5-10 minutes. Check propagation:
```bash
dig api.yourdomain.com @1.1.1.1
```

---

## Alternative: Using Cloudflare Tunnel (Zero Trust)

For even better security, use Cloudflare Tunnel (no open ports needed):

```bash
# Install cloudflared
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb

# Authenticate
cloudflared tunnel login

# Create tunnel
cloudflared tunnel create onemindos

# Configure tunnel
cat > ~/.cloudflared/config.yml << EOF
tunnel: <tunnel-id>
credentials-file: /home/zeus/.cloudflared/<tunnel-id>.json

ingress:
  - hostname: api.yourdomain.com
    service: http://localhost:8080
  - hostname: app.yourdomain.com
    service: http://localhost:3000
  - service: http_status:404
EOF

# Run tunnel
cloudflared tunnel run onemindos

# Install as service
sudo cloudflared service install
sudo systemctl start cloudflared
```

---

## Security Best Practices

1. **Enable Cloudflare WAF** (Web Application Firewall)
2. **Use Full (strict) SSL** mode
3. **Enable HSTS** with long max-age
4. **Configure rate limiting** to prevent abuse
5. **Enable Bot Fight Mode**
6. **Use Page Rules** to cache static assets
7. **Monitor logs** for suspicious activity
8. **Consider Cloudflare Access** for admin endpoints

---

## Cost

- **Cloudflare Free Tier**: $0/month
  - Unlimited DDoS protection
  - Free SSL certificates
  - CDN with caching
  - Basic WAF rules
  - Analytics

- **Cloudflare Pro**: $20/month (optional)
  - Advanced DDoS protection
  - Image optimization
  - Mobile optimization
  - More page rules

---

## Quick Reference

**Your Domains** (replace with actual domain):
- Backend API: https://api.yourdomain.com
- Frontend: https://app.yourdomain.com
- Admin: https://admin.yourdomain.com

**DNS Records**:
- `api` A → 20.121.38.186 (Proxied)
- `app` A → 20.121.38.186 (Proxied)
- `admin` A → 20.121.38.186 (Proxied)

**Backend Ports** (internal):
- Backend: 8080
- Frontend: 3000
- Portainer: 9000

**Reverse Proxy**: Caddy (auto HTTPS)
**CDN**: Cloudflare (SSL termination + caching)

---

**Last Updated**: 2026-01-27
**Document Version**: 1.0


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/security (Security Ops)/security|security]]

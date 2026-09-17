# API Key Management Security Documentation

## Overview

The API Keys Management system implements industry-standard security practices for storing and managing sensitive credentials.

## Security Features

### 1. Encryption at Rest (AES-128 via Fernet)

All API keys are encrypted before storage using the Fernet symmetric encryption scheme:

```python
from cryptography.fernet import Fernet

# Keys are encrypted before database storage
encrypted_key = FERNET.encrypt(plaintext_key.encode())
```

**Benefits:**
- AES-128 encryption in CBC mode with HMAC authentication
- Keys are never stored in plaintext
- Database breaches do not expose actual keys
- Tampering detection through HMAC

### 2. Key Masking in API Responses

API responses NEVER return full keys. All keys are masked:

```python
# Example masked output
"sk-...****1234"  # Only shows prefix and last 4 characters
```

**Benefits:**
- Prevents accidental exposure in logs
- Allows user identification without revealing key
- Safe for display in UI/dashboards

### 3. Audit Logging

Every key operation is logged:

```python
audit_log("create", key_id, key_name, details)
```

**Logged Events:**
- Key creation
- Key access (list, read)
- Key updates (metadata only)
- Key deletion
- Key usage for API calls

**NOT Logged:**
- Plaintext key values
- Encryption keys

### 4. Secure Key Generation

Encryption keys are derived deterministically from environment:

```python
# Development: Uses DATABASE_URL as seed
seed = os.getenv("API_KEY_ENCRYPTION_SEED", os.getenv("DATABASE_URL"))
key_material = hashlib.sha256(seed.encode()).digest()
encryption_key = base64.urlsafe_b64encode(key_material)
```

**WARNING:** This is for DEVELOPMENT only. In production:
- Use AWS KMS, HashiCorp Vault, or Azure Key Vault
- Rotate encryption keys regularly
- Store keys in HSM if compliance required

## API Endpoints

### Public Endpoints (Frontend-facing)

#### `GET /api/settings/keys`
List all keys (masked)
- Returns: `[{id, name, key_preview, is_active, created_at, ...}]`
- Security: Only returns masked keys

#### `POST /api/settings/keys`
Create new key
- Body: `{name, key, description?, icon?, color?}`
- Security: Encrypts key before storage

#### `PUT /api/settings/keys/{id}`
Update key metadata (NOT the key itself)
- Body: `{name?, description?, icon?, color?, is_active?}`
- Security: Cannot update encrypted key value

#### `DELETE /api/settings/keys/{id}`
Delete key permanently
- Security: Logs deletion in audit trail

#### `GET /api/settings/keys/stats/summary`
Get usage statistics
- Returns: `{total_keys, active_keys, inactive_keys, total_usage}`

### Internal Endpoints (NEVER expose to frontend)

#### `GET /api/settings/keys/{id}/decrypt`
Decrypt key for actual API use
- **DANGER:** Returns plaintext key
- **USE CASE:** Backend services making API calls
- **TODO:** Add authentication/authorization
- **TODO:** Implement rate limiting
- **TODO:** Restrict to internal network only

## Security Best Practices

### For Development

1. Set `API_KEY_ENCRYPTION_SEED` in `.env`:
   ```bash
   API_KEY_ENCRYPTION_SEED="your-secure-random-seed"
   ```

2. Never commit `.env` files to version control

3. Use different seeds for dev/staging/prod

### For Production

1. **Use a proper secret management system:**
   - AWS Secrets Manager + KMS
   - HashiCorp Vault
   - Azure Key Vault
   - Google Cloud Secret Manager

2. **Replace in-memory storage with database:**
   ```sql
   CREATE TABLE api_keys (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       name VARCHAR(255) NOT NULL,
       encrypted_key TEXT NOT NULL,
       description TEXT,
       icon VARCHAR(50),
       color VARCHAR(7),
       is_active BOOLEAN DEFAULT TRUE,
       created_at TIMESTAMPTZ DEFAULT NOW(),
       updated_at TIMESTAMPTZ DEFAULT NOW(),
       last_used TIMESTAMPTZ,
       usage_count INTEGER DEFAULT 0
   );

   CREATE INDEX idx_api_keys_active ON api_keys(is_active);
   ```

3. **Add authentication to decrypt endpoint:**
   ```python
   @router.get("/{key_id}/decrypt")
   async def decrypt_api_key(
       key_id: str,
       current_user: User = Depends(get_current_user),  # Add auth
       _: bool = Depends(require_internal_service)       # Restrict access
   ):
       # Only allow internal services or admins
       ...
   ```

4. **Implement key rotation:**
   - Scheduled rotation every 90 days
   - Manual rotation on security events
   - Grace period for old keys

5. **Add rate limiting:**
   ```python
   from slowapi import Limiter

   @router.get("/{key_id}/decrypt")
   @limiter.limit("10/minute")  # Max 10 decryptions per minute
   async def decrypt_api_key(...):
       ...
   ```

## Threat Model

### Protected Against

✅ Database breach (keys are encrypted)
✅ Log leakage (keys are masked)
✅ UI exposure (keys are masked)
✅ Accidental commits (never in plaintext)
✅ Tampering (HMAC verification)

### NOT Protected Against

❌ Compromised encryption key (use proper key management)
❌ Compromised server (attacker can decrypt in memory)
❌ Insider threats (needs access control + audit)
❌ Social engineering (user education required)

## Compliance

This implementation provides:

- **GDPR:** Encryption at rest, audit logging
- **PCI-DSS:** Encrypted storage, masked display
- **SOC 2:** Audit trails, access controls (when auth added)
- **HIPAA:** Encryption, audit logs (when database storage added)

## Migration from Mock to Production

### Step 1: Database Setup
```sql
-- Run this SQL to create the table
\i backend/settings/schema.sql
```

### Step 2: Update Code
Replace in-memory `api_keys_db` with database queries:

```python
# OLD (in-memory)
api_keys_db[key_id] = key_data

# NEW (database)
await db.execute(
    "INSERT INTO api_keys (id, name, encrypted_key, ...) VALUES ($1, $2, $3, ...)",
    key_id, name, encrypted_key, ...
)
```

### Step 3: Add Authentication
```python
from fastapi.security import HTTPBearer
from backend.auth import verify_token

security = HTTPBearer()

@router.get("/api/settings/keys")
async def list_keys(token: str = Depends(security)):
    user = await verify_token(token)
    # Return only user's keys
    ...
```

### Step 4: Secret Management
```python
# Replace
FERNET = Fernet(get_encryption_key())

# With
from aws_kms import get_data_key
encryption_key = await get_data_key("api-keys-encryption")
FERNET = Fernet(encryption_key)
```

## Testing Security

### Manual Tests

1. **Verify encryption:**
   ```bash
   # Create key
   curl -X POST http://localhost:7777/api/settings/keys \
     -H "Content-Type: application/json" \
     -d '{"name": "Test", "key": "sk-test123"}'

   # List keys (should be masked)
   curl http://localhost:7777/api/settings/keys
   # Response: {"key_preview": "sk-...****123"}
   ```

2. **Verify audit logging:**
   ```bash
   # Check logs for AUDIT entries
   grep "AUDIT:" backend.log
   ```

3. **Verify no plaintext in database:**
   ```sql
   SELECT encrypted_key FROM api_keys LIMIT 1;
   -- Should see gibberish, not actual key
   ```

### Automated Tests

```python
# tests/test_api_keys_security.py

def test_keys_are_encrypted():
    key = "sk-test123"
    encrypted = encrypt_key(key)
    assert encrypted != key
    assert "sk-test123" not in encrypted

def test_keys_are_masked_in_response():
    response = client.get("/api/settings/keys")
    for key in response.json():
        assert "****" in key["key_preview"]
        assert len(key["key_preview"]) < 20  # Much shorter than full key

def test_decrypt_endpoint_not_public():
    # TODO: Should fail without authentication
    response = client.get("/api/settings/keys/123/decrypt")
    assert response.status_code == 401  # Unauthorized

def test_cannot_update_encrypted_key():
    response = client.put(
        "/api/settings/keys/123",
        json={"key": "new-key"}  # Should be ignored
    )
    # Verify key was not updated
    ...
```

## Future Enhancements

1. **Multi-tenancy:** Isolate keys per user/org
2. **Key expiration:** Auto-expire keys after N days
3. **Usage limits:** Enforce per-key rate limits
4. **Webhook notifications:** Alert on key usage/rotation
5. **Read-only keys:** Grant limited permissions
6. **Key versioning:** Keep history of rotated keys
7. **Compliance reports:** Automated audit exports

## Support

For security issues, DO NOT create public GitHub issues.
Contact: security@onemind.ai

## References

- [Cryptography.io Fernet Documentation](https://cryptography.io/en/latest/fernet/)
- [OWASP Key Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Key_Management_Cheat_Sheet.html)
- [NIST SP 800-57: Key Management](https://csrc.nist.gov/publications/detail/sp/800-57-part-1/rev-5/final)


---

## Related

- [[00-24 SO (Sovereign Operations)/04 Security (Cyber & Digital Shield) (The Shield)/04 Security|04 Security]]

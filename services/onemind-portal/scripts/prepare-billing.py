#!/usr/bin/env python3
"""Save a user-created billing:read-only credential locally without printing it."""
import getpass,json,os,pathlib,urllib.request
print('In DigitalOcean API tokens, create a short-lived token with only billing:read.')
print('Do not enter your general administration/doctl token.')
token=getpass.getpass('Billing-only token (hidden): ').strip()
try:
    request=urllib.request.Request('https://api.digitalocean.com/v2/customers/my/balance',headers={'Authorization':'Bearer '+token})
    with urllib.request.urlopen(request,timeout=10) as response:
        raw=response.read(32769)
        if len(raw)>32768:raise ValueError()
        data=json.loads(raw)
    if 'month_to_date_usage' not in data or 'generated_at' not in data:raise ValueError()
    target=pathlib.Path.home()/'.config/onemind/billing-readonly.json'
    target.parent.mkdir(mode=0o700,parents=True,exist_ok=True)
    fd=os.open(target,os.O_WRONLY|os.O_CREAT|os.O_EXCL,0o600)
    with os.fdopen(fd,'w') as handle:json.dump({'DO_BILLING_READ_TOKEN':token},handle)
    print('Saved locally with owner-only permissions. Billing access responded; token scopes still follow your console selection.')
    print('Seal as backstage-billing in namespace backstage through the documented workflow. No token or balance printed.')
except FileExistsError:raise SystemExit('Existing local credential preserved.')
except Exception:raise SystemExit('Billing connection failed; verify the token locally. No provider response printed.')

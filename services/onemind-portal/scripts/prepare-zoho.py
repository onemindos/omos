#!/usr/bin/env python3
"""Exchange a user-generated US Self Client grant locally; never print tokens."""
import getpass,json,os,pathlib,urllib.parse,urllib.request
target=pathlib.Path.home()/'.config/onemind/zoho-readonly.json'
print('Create a US Zoho Self Client in https://api-console.zoho.com/')
print('Generate a short-lived grant for ZohoProjects.tasks.READ.')
print('For tasks plus Calendar: ZohoProjects.tasks.READ,ZohoCalendar.calendar.READ,ZohoCalendar.event.READ.')
print('Do not use ALL or write scopes. Confirm your numeric Projects portal ID.')
client=getpass.getpass('Client ID (hidden): ')
secret=getpass.getpass('Client secret (hidden): ')
code=getpass.getpass('Short-lived grant code (hidden): ')
portal=input('Numeric Projects portal ID: ').strip()
if not portal.isdigit():raise SystemExit('A numeric portal ID is required.')
calendar=input('Selected Calendar UID (optional; leave empty for tasks only): ').strip()
body=urllib.parse.urlencode({'grant_type':'authorization_code','client_id':client,'client_secret':secret,'code':code}).encode()
try:
    with urllib.request.urlopen(urllib.request.Request('https://accounts.zoho.com/oauth/v2/token',data=body),timeout=15) as response:
        raw=response.read(32769)
        if len(raw)>32768:raise ValueError()
        tokens=json.loads(raw)
    if not isinstance(tokens.get('refresh_token'),str):raise ValueError()
    target.parent.mkdir(mode=0o700,parents=True,exist_ok=True)
    fd=os.open(target,os.O_WRONLY|os.O_CREAT|os.O_EXCL,0o600)
    with os.fdopen(fd,'w') as handle:json.dump({'ZOHO_CLIENT_ID':client,'ZOHO_CLIENT_SECRET':secret,'ZOHO_REFRESH_TOKEN':tokens['refresh_token'],'ZOHO_PROJECTS_PORTAL_ID':portal,**({'ZOHO_CALENDAR_UID':calendar} if calendar else {})},handle)
    print('Connection saved locally with owner-only permissions. No token was printed.')
    print('Use the SealedSecret workflow. Do not paste this file into chat.')
except FileExistsError:raise SystemExit('Existing connection preserved. Review it locally before replacement.')
except Exception:raise SystemExit('Exchange failed. Check region, grant expiry and client configuration locally. No provider response was printed.')

"""Provision the dedicated OIDC application and seal deployment credentials.
Requires kubectl access, kubeseal, and an Authentik port-forward on localhost:19000.
Never prints or writes plaintext credentials. Existing credentials are preserved.
"""
import base64, json, pathlib, secrets, subprocess, urllib.request
ROOT = pathlib.Path(__file__).resolve().parents[3]
OUT = ROOT / 'platform/backstage'
def kube(*args): return subprocess.check_output(['kubectl', *args])
def get_secret(namespace, name):
    r=subprocess.run(['kubectl','get','secret',name,'-n',namespace,'-o','json'],capture_output=True)
    return json.loads(r.stdout) if r.returncode==0 else None
bootstrap=get_secret('auth','authentik-secret')
token=base64.b64decode(bootstrap['data']['AUTHENTIK_BOOTSTRAP_TOKEN']).decode()
def api(path, data=None, method=None):
    req=urllib.request.Request('http://127.0.0.1:19000/api/v3/'+path,
        data=json.dumps(data).encode() if data is not None else None,
        method=method, headers={'Authorization':'Bearer '+token,'Content-Type':'application/json'})
    with urllib.request.urlopen(req) as response: return json.load(response)
providers=api('providers/oauth2/?page_size=100')['results']
provider=next((p for p in providers if p['name']=='OneMind Backstage'),None)
if not provider:
    flows={f['slug']:f['pk'] for f in api('flows/instances/?page_size=100')['results']}
    scopes=[x['pk'] for x in api('propertymappings/provider/scope/?page_size=100')['results'] if x['scope_name'] in ['openid','profile','email','offline_access']]
    keys=api('crypto/certificatekeypairs/')['results']
    key=next(x['pk'] for x in keys if x['name']=='authentik Self-signed Certificate')
    provider=api('providers/oauth2/',{
        'name':'OneMind Backstage','client_type':'confidential','client_id':'onemind-backstage',
        'client_secret':secrets.token_urlsafe(48),
        'authorization_flow':flows['default-provider-authorization-implicit-consent'],
        'invalidation_flow':flows['default-provider-invalidation-flow'],
        'property_mappings':scopes,'signing_key':key,'issuer_mode':'per_provider',
        'grant_types':['authorization_code','refresh_token'],
        'redirect_uris':[{'matching_mode':'strict','url':'https://'+h+'.tail717ba0.ts.net/api/auth/oidc/handler/frame'} for h in ['onemind-preview','homepage']],
    })
apps=api('core/applications/?page_size=100')['results']
app=next((a for a in apps if a['slug']=='onemind-backstage'),None)
if not app:
    app=api('core/applications/',{'name':'OneMind','slug':'onemind-backstage','provider':provider['pk'],
        'meta_launch_url':'https://onemind-preview.tail717ba0.ts.net','meta_description':'OneMind service catalog and operations portal','policy_engine_mode':'any'})
# The first release permits only the existing administrator to authenticate.
admin=next(u for u in api('core/users/?page_size=100')['results'] if u['username']=='akadmin' and u['is_active'])
bindings=api('policies/bindings/?target='+app['pk'])['results']
if not any(b.get('user')==admin['pk'] for b in bindings):
    api('policies/bindings/',{'target':app['pk'],'user':admin['pk'],'order':0,'enabled':True})
existing=get_secret('backstage','backstage-secrets')
if existing:
    values={k:base64.b64decode(v).decode() for k,v in existing['data'].items()}
else:
    values={'POSTGRES_PASSWORD':secrets.token_urlsafe(48),'AUTH_SESSION_SECRET':secrets.token_urlsafe(48)}
values.update({'AUTH_OIDC_CLIENT_ID':provider['client_id'],'AUTH_OIDC_CLIENT_SECRET':provider['client_secret']})
def seal(name,values,stype='Opaque'):
    obj={'apiVersion':'v1','kind':'Secret','metadata':{'name':name,'namespace':'backstage'},'type':stype,'stringData':values}
    result=subprocess.check_output(['kubeseal','--controller-name','sealed-secrets-controller','--controller-namespace','kube-system','--format','yaml'],input=json.dumps(obj).encode())
    (OUT/(name+'.sealedsecret.yaml')).write_bytes(result)
seal('backstage-secrets',values)
registry=get_secret('onetak','ghcr')
seal('ghcr',{k:base64.b64decode(v).decode() for k,v in registry['data'].items()},'kubernetes.io/dockerconfigjson')
(OUT/'auth-registration.json').write_text(json.dumps({'provider_id':provider['pk'],'application_slug':app['slug'],'initial_catalog_user':admin['username']},indent=2)+'\n')
print('OIDC application registered; administrator binding and encrypted secret manifests prepared.')

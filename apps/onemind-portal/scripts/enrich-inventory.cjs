// Reviewed inventory refresh. No credentials, IP addresses, or user profiles are captured.
const fs=require('fs'),cp=require('child_process'),path=require('path'),YAML=require('yaml'),crypto=require('crypto');
process.chdir(path.resolve(__dirname,'..'));
cp.execFileSync(process.execPath,['scripts/snapshot-cluster.cjs'],{stdio:'inherit'});
const load=p=>YAML.parseAllDocuments(fs.readFileSync(p,'utf8')).map(x=>x.toJSON());
const cluster=load('catalog/cluster.yaml'),apps=load('catalog/apps.yaml');
const observed=new Date().toISOString();
const kube=JSON.parse(cp.execFileSync('kubectl',['get','services,ingresses','-A','-o','json'],{encoding:'utf8',maxBuffer:20*1024*1024}));
const services=kube.items.filter(x=>x.kind==='Service'),ingresses=kube.items.filter(x=>x.kind==='Ingress');
for(const entity of cluster){entity.metadata.annotations={...entity.metadata.annotations,'onemind.io/hosting':'digitalocean-nyc3','onemind.io/inventory-review':'discovered'};}
for(const ingress of ingresses){
 const ns=ingress.metadata.namespace;
 const hosts=[...(ingress.spec.rules||[]).map(x=>x.host),...(ingress.status?.loadBalancer?.ingress||[]).map(x=>x.hostname)].filter(Boolean);
 const names=(ingress.spec.rules||[]).flatMap(x=>(x.http?.paths||[]).map(p=>p.backend?.service?.name));
 if(ingress.spec.defaultBackend?.service?.name)names.push(ingress.spec.defaultBackend.service.name);
 for(const name of names){const service=services.find(x=>x.metadata.namespace===ns&&x.metadata.name===name);if(!service?.spec.selector||!Object.keys(service.spec.selector).length)continue;
  const workloads=cluster.filter(x=>x.kind==='Component'&&x.metadata.annotations['backstage.io/kubernetes-namespace']===ns&&Object.entries(service.spec.selector).every(([k,v])=>x.metadata.annotations['backstage.io/kubernetes-label-selector'].split(',').includes(`${k}=${v}`)));
  for(const app of apps.filter(x=>x.kind==='Component')){const url=app.metadata.links?.[0]?.url;if(!url||!hosts.includes(new URL(url).hostname))continue;
   app.metadata.annotations={...app.metadata.annotations,'onemind.io/hosting':'digitalocean-nyc3','onemind.io/relationship-source':'ingress-service-selector'};
   app.spec.dependsOn=[...new Set([...(app.spec.dependsOn||[]),...workloads.map(x=>'component:default/'+x.metadata.name)])];
  }
 }
}
for(const app of apps.filter(x=>x.kind==='Component')){
 app.metadata.annotations={...app.metadata.annotations,'onemind.io/inventory-review':'imported-link'};
 if(!app.metadata.annotations['onemind.io/hosting'])app.metadata.annotations['onemind.io/hosting']=app.metadata.links?.some(x=>new URL(x.url).hostname.endsWith('.tail717ba0.ts.net'))?'mesh-endpoint-unverified':'external-provider-unverified';
}
const mesh=JSON.parse(cp.execFileSync('tailscale',['status','--json'],{encoding:'utf8',maxBuffer:20*1024*1024}));
cluster.push({apiVersion:'backstage.io/v1alpha1',kind:'System',metadata:{name:'onemind-mesh',title:'OneMind mesh',description:'Discovered Tailscale endpoints; these include devices and service proxies, not only physical machines.'},spec:{owner:'group:default/onemind-admins'}});
for(const peer of [mesh.Self,...Object.values(mesh.Peer||{})].filter(Boolean)){
 const hash=crypto.createHash('sha256').update(peer.ID||peer.DNSName||peer.HostName).digest('hex').slice(0,12);
 cluster.push({apiVersion:'backstage.io/v1alpha1',kind:'Resource',metadata:{name:'mesh-'+hash,title:peer.HostName,description:`${peer.OS||'Unknown OS'} mesh endpoint. Availability is a discovery snapshot, not a live health check.`,annotations:{'onemind.io/inventory-source':'tailscale','onemind.io/inventory-observed-at':observed,'onemind.io/inventory-review':'unreviewed','onemind.io/mesh-dns':peer.DNSName||'','onemind.io/observed-online':String(!!peer.Online)},tags:['mesh','unreviewed']},spec:{type:'mesh-endpoint',owner:'group:default/onemind-admins',system:'onemind-mesh'}});
}
for(const [file,docs]of [['cluster.yaml',cluster],['apps.yaml',apps]])for(const base of ['catalog','../../platform/backstage/catalog'])fs.writeFileSync(`${base}/${file}`,docs.map(x=>YAML.stringify(x)).join('---\n'));
console.log('Inventory:',cluster.filter(x=>x.kind==='Component').length,'workloads;',cluster.filter(x=>x.spec.type==='mesh-endpoint').length,'mesh endpoints;',apps.filter(x=>x.spec?.dependsOn?.length).length,'linked application entries');

cp.execFileSync(process.execPath,['scripts/reconcile-fabric.cjs'],{stdio:'inherit'});

const repos=JSON.parse(cp.execFileSync('gh',['repo','list','onemindos','--limit','200','--json','name,url,isArchived,description'],{encoding:'utf8'}));
if(repos.length===200)throw new Error('Repository inventory may be truncated; increase limit before publishing');
fs.writeFileSync('catalog/repositories-source.json',JSON.stringify(repos.map(repo=>({...repo,inventoryState:repo.isArchived?'archived':'observed',observedAt:new Date().toISOString()})),null,2)+'\n');
const runtime=JSON.parse(cp.execFileSync('kubectl',['get','deployments,statefulsets,daemonsets,cronjobs,services','-A','-o','json'],{encoding:'utf8',maxBuffer:20*1024*1024}));
fs.writeFileSync('catalog/runtime-source.json',JSON.stringify({observedAt:new Date().toISOString(),items:runtime.items.map(x=>({kind:x.kind,namespace:x.metadata.namespace,name:x.metadata.name}))},null,2)+'\n');
cp.execFileSync(process.execPath,['scripts/curate-dependencies.cjs'],{stdio:'inherit'});
cp.execFileSync(process.execPath,['scripts/verify-inventory.cjs'],{stdio:'inherit'});

// Only exact recorded identifiers or reviewed aliases establish deployed relationships.
const fs=require('fs'),path=require('path'),Y=require('yaml'),crypto=require('crypto');process.chdir(path.resolve(__dirname,'..'));
const load=p=>Y.parseAllDocuments(fs.readFileSync(p,'utf8')).map(x=>x.toJSON());
const apps=load('catalog/apps.yaml'),cluster=load('catalog/cluster.yaml').filter(x=>!x.metadata.name.startsWith('repo-'));
const repos=JSON.parse(fs.readFileSync('catalog/repositories-source.json')),runtime=JSON.parse(fs.readFileSync('catalog/runtime-source.json'));
const rows=JSON.parse(fs.readFileSync('catalog/fabric-source.json')),pages=JSON.parse(fs.readFileSync('catalog/fabric-pages.json'));
const key=s=>s.toLowerCase().replace(/[^a-z0-9]/g,'');
const aliases={ 'NATS JetStream Cluster':['nats/nats'], 'ClickHouse Time Machine':['clickhouse/chi-onemind-onemind-0-0'], 'OneTAK PostgreSQL':['onetak/onetak-postgis'], 'OneTAK MinIO':['onetak/onetak-minio'], 'TAK PostgreSQL':['takserver/takserver-postgresql'], 'Geo Database':['onetak/onetak-geo-db'], 'TAKServer Ignite':['takserver/takserver-ignite'], 'TAK Server (Core)':['takserver/takserver-api','takserver/takserver-messaging'], 'TAK-ML Server':['takml/takml-server'], 'TAK-ML Database':['takml/takml-db'], 'TAK-ML Vision Worker':['takml/takml-vision-worker'], 'Murmur (Voice)':['takserver/murmur'], 'Hermes — Legacy Hive System':['hermes/hermes-legacy'], 'MediaMTX':['onetak/onetak-media'], 'Redpanda Connect (Wombat)':['omos/redpanda-connect'], 'TileServer GL':['onetak/tileserver-gl'], 'STAC API':['onetak/stac-api'], 'Cloudflare Tunnel':['cloudflare/cloudflared'], 'GitHub Runner':['github-runner/github-runner'] };
const repoId=r=>'repo-'+key(r.name).slice(0,40)+'-'+crypto.createHash('sha256').update(r.url).digest('hex').slice(0,8);
cluster.push({apiVersion:'backstage.io/v1alpha1',kind:'System',metadata:{name:'repo-onemind',title:'OneMind repositories'},spec:{owner:'group:default/onemind-admins'}});
const repositoryState=r=>r.inventoryState==='missing-review-required'?'repository-missing-review-required':r.inventoryState==='archived'?'repository-archived':r.inventoryState==='observed'?'repository-verified':'repository-documented-unverified';
for(const r of repos)cluster.push({apiVersion:'backstage.io/v1alpha1',kind:'Resource',metadata:{name:repoId(r),title:r.name,description:r.description||'OneMind source repository',tags:['repository',repositoryState(r)],links:[{url:r.url,title:'Source repository'}],annotations:{'onemind.io/verification':repositoryState(r),'onemind.io/inventory-state':r.inventoryState||'legacy-observation-unverified',...(r.observedAt?{'onemind.io/observed-at':r.observedAt}:{}),...(r.lastCheckedAt?{'onemind.io/inventory-last-checked-at':r.lastCheckedAt}:{})}},spec:{type:'code-repository',owner:'group:default/onemind-admins',system:'repo-onemind'}});
fs.mkdirSync('../../docs/inventory/fabric',{recursive:true});const manifest=[];
for(const row of rows){const id=row.url.split('/').pop();const page=pages.find(p=>p.source===row.url);const file=`docs/inventory/fabric/${id}.md`;
 const body=`# ${row.Page}\n\nArchived from the Fabric registry on 2026-10-06. Source edited: ${row['Last edited time']}.\n\nThese are historical documentation claims, not current deployment instructions or live health.\n\n- Category: ${row.Category||'Unspecified'}\n- Documented host: ${row.Host||'Unspecified'}\n- Recorded status: ${row.Status||'Unspecified'}\n- Recorded address: ${row['userDefined:URL']||'None'}\n- Tags: ${row.Tags?.length?row.Tags.join(', '):'None'}\n\n${row.Description||''}\n\n${page?.content||'The source page had no body content.'}\n`;
 fs.writeFileSync('../../'+file,body);manifest.push({title:row.Page,file,blank:page?.blank,complete:!!page&&!page.truncated&&!page.unknownBlockCount});
}
let verified=0,repositoryMatches=0;
for(const entity of [...apps,...cluster]){
 if(entity.metadata.name.startsWith('repo-'))continue;
 const links=entity.metadata.links||[];const sources=rows.filter(r=>links.some(l=>l.url===r.url||l.url.endsWith('/'+r.url.split('/').pop()+'.md')));
 if(!sources.length)continue;
 const source=sources[0],a=entity.metadata.annotations||{};
 const candidates=aliases[source.Page]||runtime.items.filter(x=>key(x.name)===key(source.Page)).map(x=>x.namespace+'/'+x.name);
 const exact=runtime.items.filter(x=>candidates.includes(x.namespace+'/'+x.name));
 const workloads=cluster.filter(x=>x.kind==='Component'&&exact.some(i=>x.metadata.title===i.namespace+' / '+i.name));
 const url=source['userDefined:URL']?.replace('github.com/onemindos/omos-infra','github.com/onemindos/omos').replace(/\/$/,'');
 const repo=repos.find(r=>r.url===url);
 const references=entity.metadata.tags?.includes('reference'),planned=entity.metadata.tags?.includes('planned');
 const deploymentVerified=!!workloads.length||a['onemind.io/relationship-source']==='ingress-service-selector';
 const state=deploymentVerified?'deployment-verified':references?'reference':planned?'planned':repo?repositoryState(repo):'documented-unverified';
 a['onemind.io/verification']=state;
 if(deploymentVerified)a['onemind.io/verified-at']=runtime.observedAt;
 else if(repo){if(repo.observedAt)a['onemind.io/verified-at']=repo.observedAt;else delete a['onemind.io/verified-at'];a['onemind.io/repository-inventory-state']=repo.inventoryState||'legacy-observation-unverified';if(repo.lastCheckedAt)a['onemind.io/repository-last-checked-at']=repo.lastCheckedAt;}
 else delete a['onemind.io/verified-at'];
 if(exact.length)a['onemind.io/runtime-evidence']=exact.map(x=>x.kind+' '+x.namespace+'/'+x.name).join('; ');
 if(repo)a['onemind.io/repository-url']=repo.url;
 entity.metadata.annotations=a;entity.metadata.tags=[...new Set([...(entity.metadata.tags||[]).filter(x=>!['documented','deployment-verified','repository-verified','documented-unverified','repository-archived','repository-missing-review-required','repository-documented-unverified'].includes(x)),state])];
 if(deploymentVerified)verified++;if(repo)repositoryMatches++;
 if(workloads.length||repo)entity.spec.dependsOn=[...new Set([...(entity.spec.dependsOn||[]),...workloads.map(x=>'component:default/'+x.metadata.name),...(repo?['resource:default/'+repoId(repo)]:[])])];
 entity.metadata.links=links.map(l=>{const row=rows.find(r=>r.url===l.url);return row?{url:'https://github.com/onemindos/omos/blob/main/docs/inventory/fabric/'+row.url.split('/').pop()+'.md',title:'Archived registry: '+row.Page}:l;});
}
fs.writeFileSync('../../docs/inventory/migration-manifest.json',JSON.stringify({reviewedAt:runtime.observedAt,total:rows.length,pages:manifest},null,2)+'\n');
for(const [file,docs]of [['apps.yaml',apps],['cluster.yaml',cluster]])for(const base of ['catalog','../../platform/backstage/catalog'])fs.writeFileSync(`${base}/${file}`,docs.map(x=>Y.stringify(x)).join('---\n'));
console.log(JSON.stringify({repositories:repos.length,deploymentVerifiedEntries:verified,repositoryMatchedEntries:repositoryMatches,archivedPages:manifest.length,completePages:manifest.filter(x=>x.complete).length}));

// Reconcile curated documentation without treating recorded status as live health.
const fs=require('fs'),path=require('path'),YAML=require('yaml'),crypto=require('crypto');
process.chdir(path.resolve(__dirname,'..'));
const rows=JSON.parse(fs.readFileSync('catalog/fabric-source.json','utf8'));
const load=file=>YAML.parseAllDocuments(fs.readFileSync(file,'utf8')).map(x=>x.toJSON());
const apps=load('catalog/apps.yaml'),cluster=load('catalog/cluster.yaml').filter(x=>!x.metadata.name.startsWith('fabric-'));
const key=s=>s.toLowerCase().replace(/[^a-z0-9]/g,'');
const aliases={'googlecalendar':'calendar','googledrive':'drive','googlephotos':'photos'};
const groups=new Map();for(const r of rows){const name=aliases[key(r.Page)]||key(r.Page);groups.set(name,[...(groups.get(name)||[]),r]);}
let linked=0,created=0;
for(const [identity,records]of groups){
 records.sort((a,b)=>b['Last edited time'].localeCompare(a['Last edited time']));const r=records[0];
 const url=r['userDefined:URL'];const normalized=s=>s?.replace(/\/$/,'');
 let target=apps.find(x=>x.kind==='Component'&&(key(x.metadata.title||x.metadata.name)===identity||aliases[key(x.metadata.title||'')]===identity));
 if(!target&&url&&!url.includes('github.com/'))target=apps.find(x=>x.kind==='Component'&&x.metadata.links?.some(l=>normalized(l.url)===normalized(url)));
 if(!target&&url)try{const host=new URL(url).hostname;target=cluster.find(x=>x.spec.type==='mesh-endpoint'&&x.metadata.annotations?.['onemind.io/mesh-dns']?.replace(/\.$/,'')===host);}catch{}
 const reference=/reference|pending|not yet built|removed, redundant|future cloud/i.test(r.Description||'')||r.Status==='📚 Reference';
 const planned=r.Status==='🏗️ Planned';
 if(!target){
  const category='fabric-'+key(r.Category||'uncategorized');
  if(!cluster.some(x=>x.metadata.name===category))cluster.push({apiVersion:'backstage.io/v1alpha1',kind:'System',metadata:{name:category,title:'Fabric / '+(r.Category||'Uncategorized')},spec:{owner:'group:default/onemind-admins'}});
  target={apiVersion:'backstage.io/v1alpha1',kind:'Resource',metadata:{name:'fabric-'+identity.slice(0,42)+'-'+crypto.createHash('sha256').update(identity).digest('hex').slice(0,8),title:r.Page,description:(r.Description||'Documented registry entry')+' — Notion documentation; deployment and availability unverified.',tags:['fabric-registry',reference?'reference':planned?'planned':'documented']},spec:{type:reference?'reference':planned?'planned-capability':'documented-capability',owner:'group:default/onemind-admins',system:category}};cluster.push(target);created++;
 }else linked++;
 const annotations=target.metadata.annotations||{};
 Object.assign(annotations,{'onemind.io/documented-purpose':r.Description||'Not supplied','onemind.io/documented-host':r.Host||'Not supplied','onemind.io/documented-status':r.Status||'Not supplied','onemind.io/documented-category':r.Category||'Not supplied','onemind.io/documented-at':r['Last edited time'],'onemind.io/registry-review':'source-reviewed-runtime-unverified'});
 if(reference&&r.Status==='✅ Live')annotations['onemind.io/registry-conflict']='Recorded Live but description identifies a reference or future deployment';
 target.metadata.annotations=annotations;
 const links=target.metadata.links||[];for(const row of records)if(!links.some(x=>x.url===row.url))links.push({url:row.url,title:'Notion source: '+row.Page});
 // Source URLs may be protocol endpoints or old IPs; keep them in documentation, not launch buttons.
 target.metadata.links=links;
}
for(const [file,docs]of [['apps.yaml',apps],['cluster.yaml',cluster]])for(const base of ['catalog','../../platform/backstage/catalog'])fs.writeFileSync(`${base}/${file}`,docs.map(x=>YAML.stringify(x)).join('---\n'));
console.log(`${rows.length} source rows -> ${groups.size} reconciled identities; ${linked} existing entries enriched; ${created} documented resources added`);

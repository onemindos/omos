// Run after reviewing the active kubectl context. Captures identities, never credentials or runtime specs.
const path=require('path');process.chdir(path.resolve(__dirname,'..'));const crypto=require('crypto');const fs=require('fs');const cp=require('child_process');const YAML=require('yaml');
const context=cp.execFileSync('kubectl',['config','current-context'],{encoding:'utf8'}).trim();
if(context!=='do-nyc3-omos-cloud')throw new Error('Unexpected Kubernetes context');
const data=JSON.parse(cp.execFileSync('kubectl',['get','deployments,statefulsets','-A','-o','json'],{encoding:'utf8',maxBuffer:20*1024*1024}));
const docs=[{apiVersion:'backstage.io/v1alpha1',kind:'Resource',metadata:{name:context,title:'OneMind Cloud — NYC3'},spec:{type:'kubernetes-cluster',owner:'group:default/onemind-admins'}}];
const namespaces=[...new Set(data.items.map(x=>x.metadata.namespace))];
for(const ns of namespaces)docs.push({apiVersion:'backstage.io/v1alpha1',kind:'System',metadata:{name:'cluster-'+ns,title:'Cluster / '+ns},spec:{owner:'group:default/onemind-admins'}});
for(const x of data.items){
 const ns=x.metadata.namespace; const name=x.metadata.name;
 const labels=Object.entries(x.spec.selector.matchLabels).map(([k,v])=>k+'='+v).join(',');
 docs.push({apiVersion:'backstage.io/v1alpha1',kind:'Component',metadata:{name:('k8s-'+ns+'-'+name).slice(0,54)+'-'+crypto.createHash('sha256').update(ns+'/'+name).digest('hex').slice(0,8),title:ns+' / '+name,description:x.kind+' in '+ns,
 annotations:{'backstage.io/kubernetes-label-selector':labels,'backstage.io/kubernetes-namespace':ns,'onemind.io/inventory-source':'kubernetes','onemind.io/inventory-observed-at':new Date().toISOString()},tags:['kubernetes',ns]},
 spec:{type:'service',lifecycle:'production',owner:'group:default/onemind-admins',system:'cluster-'+ns,dependsOn:['resource:default/'+context]}});
}
const rendered=docs.map(x=>YAML.stringify(x)).join('---\n');fs.writeFileSync('catalog/cluster.yaml',rendered);fs.mkdirSync('../../platform/backstage/catalog',{recursive:true});fs.writeFileSync('../../platform/backstage/catalog/cluster.yaml',rendered);console.log('Captured '+data.items.length+' workload identities across '+namespaces.length+' namespaces');

// Repeatable link migration. Run from the application root.
const path=require('path');process.chdir(path.resolve(__dirname,'..'));const fs=require('fs');
const YAML=require('yaml');
const source=YAML.parse(fs.readFileSync('../../platform/homepage/configmap.yaml','utf8'));
const groups=YAML.parse(source.data['services.yaml']);
const slug=x=>x.toLowerCase().replace(/[^a-z0-9]+/g,'-').replace(/^-|-$/g,'');
const entities=[];const directory=[];const names=new Set();
for(const group of groups){
 const [title,items]=Object.entries(group)[0]; const system=slug(title);
 entities.push({apiVersion:'backstage.io/v1alpha1',kind:'System',metadata:{name:system,title},spec:{owner:'group:default/onemind-admins'}});
 for(const item of items){
  const [title,props]=Object.entries(item)[0];if(!props.href)continue;
  let name=slug(title);if(names.has(name))name=system+'-'+name;names.add(name);
  const entry={name,title,url:props.href,description:props.description||title,group:system,groupTitle:Object.keys(group)[0]};
  directory.push(entry);
  entities.push({apiVersion:'backstage.io/v1alpha1',kind:'Component',metadata:{name,title,description:entry.description,tags:[system],links:[{url:props.href,title:'Open application',icon:'dashboard'}]},spec:{type:'website',lifecycle:'production',owner:'group:default/onemind-admins',system}});
 }
}
fs.mkdirSync('catalog',{recursive:true});
fs.writeFileSync('catalog/apps.yaml',entities.map(x=>YAML.stringify(x)).join('---\n'));
fs.writeFileSync('packages/app/src/modules/home/directory.json',JSON.stringify(directory,null,2)+'\n');
console.log(`Migrated ${directory.length} app links across ${groups.length} groups`);

fs.mkdirSync('../../platform/backstage/catalog',{recursive:true});fs.copyFileSync('catalog/apps.yaml','../../platform/backstage/catalog/apps.yaml');

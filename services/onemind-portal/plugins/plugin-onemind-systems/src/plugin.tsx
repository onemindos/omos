import { useEffect, useMemo, useState } from 'react';
import { createFrontendPlugin, PageBlueprint } from '@backstage/frontend-plugin-api';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, Button, Link, MenuItem, TextField, Typography } from '@material-ui/core';
import AccountTreeIcon from '@material-ui/icons/AccountTree';

type Entity = {kind:string;metadata:{name:string;namespace?:string;title?:string;annotations?:Record<string,string>;links?:{url:string;title?:string}[]};relations?:{type:string;targetRef:string}[]};
const ref = (e:Entity) => `${e.kind.toLowerCase()}:${e.metadata.namespace||'default'}/${e.metadata.name}`;
const catalogUrl = (r:string) => {const [kind,rest]=r.split(':');const [namespace,name]=rest.split('/');return `/catalog/${encodeURIComponent(namespace)}/${encodeURIComponent(kind)}/${encodeURIComponent(name)}`;};
// Catalog metadata is untrusted. Resolve relative paths against a fixed same-origin
// placeholder and reject executable schemes, protocol-relative URLs and backslashes.
function safeMetadataLink(value:string):string|undefined {
  // eslint-disable-next-line no-control-regex
  if (typeof value !== 'string' || !value || /[\u0000-\u0020\u007f\\]/.test(value)) return undefined;
  try {
    if (/^https?:\/\//i.test(value)) {
      const parsed=new URL(value);
      return !parsed.username&&!parsed.password?parsed.toString():undefined;
    }
    if (/^[a-z][a-z0-9+.-]*:/i.test(value)||value.startsWith('//')) return undefined;
    const parsed=new URL(value,'https://onemind.invalid/');
    return parsed.origin==='https://onemind.invalid'?`${parsed.pathname}${parsed.search}${parsed.hash}`:undefined;
  } catch { return undefined; }
}
function InventoryState({entity}:{entity:Entity}) {
  const annotations=entity.metadata.annotations||{};
  const state=annotations['onemind.io/inventory-state'];
  const observed=annotations['onemind.io/inventory-observed-at']||annotations['onemind.io/observed-at'];
  const checked=annotations['onemind.io/inventory-last-checked-at'];
  const review=state==='missing-review-required'||state==='retirement-review-required';
  return (
    <Box mt={1}>
    <Typography variant="caption" display="block" color={review?'error':'textSecondary'}>Inventory state: {state||'not recorded'}{review?' · review required; absence does not prove deletion':''}</Typography>
    <Typography variant="caption" display="block">Last successful discovery: {observed||'not recorded'}</Typography>
    {checked&&<Typography variant="caption" display="block">Last discovery check: {checked}</Typography>}
    {annotations['onemind.io/inventory-missing-observations']&&<Typography variant="caption" display="block">Missing discovery observations: {annotations['onemind.io/inventory-missing-observations']}</Typography>}
  </Box>
  );
}
function Systems(){
  const discovery = useApi(discoveryApiRef);
  const api = useApi(fetchApiRef);
  const [entities, setEntities] = useState<Entity[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [partial, setPartial] = useState(false);
  const [loadedAt, setLoadedAt] = useState('');
  const [query, setQuery] = useState('');
  const [focus, setFocus] = useState('');
  const [direction, setDirection] = useState('dependsOn');
  const [reload, setReload] = useState(0);
  useEffect(()=>{let active=true;const controller=new AbortController();setLoading(true);setError(false);
    (async()=>{try{const base=await discovery.getBaseUrl('catalog');let cursor:string|undefined;const result:Entity[]=[];
      do{const url=new URL(`${base}/entities/by-query`);url.searchParams.set('limit','500');if(cursor)url.searchParams.set('cursor',cursor);
        const response=await api.fetch(url.toString(),{signal:controller.signal});if(!response.ok)throw new Error('Catalog unavailable');
        const page=await response.json();result.push(...page.items);cursor=page.pageInfo?.nextCursor;
      }while(cursor&&result.length<5000);
      if(active){setEntities(result);setPartial(!!cursor);setLoadedAt(new Date().toISOString());const initial=result.find(e=>e.metadata.name==='k8s-backstage-backstage-21b99354')||result[0];if(initial)setFocus(previous=>previous||ref(initial));}
    }catch{if(active)setError(true);}finally{if(active)setLoading(false);}})();return()=>{active=false;controller.abort();};
  },[discovery,api,reload]);
  const byRef=useMemo(()=>new Map(entities.map(e=>[ref(e),e])),[entities]);
  const matches=entities.filter(e=>['component','resource','system'].includes(e.kind.toLowerCase())&&`${e.metadata.title||''} ${ref(e)}`.toLowerCase().includes(query.toLowerCase()));
  const entity=byRef.get(focus);
  const related=(r:string)=>byRef.get(r)?.relations?.filter(x=>x.type===direction).map(x=>x.targetRef)||[];
  const first=related(focus);
  const label=(r:string)=>byRef.get(r)?.metadata.title||r;
  const node=(r:string)=><Box p={2} border={1} borderColor="divider" borderRadius={8}><Link href={catalogUrl(r)}>{label(r)}</Link><Typography variant="caption" display="block">{r}</Typography>{byRef.get(r)?<InventoryState entity={byRef.get(r)!}/>:<Typography color="error">Referenced entity is missing from this catalog view.</Typography>}<Button size="small" onClick={()=>setFocus(r)}>Explore</Button></Box>;
  return (
    <Box p={3}><Typography variant="h4">Systems</Typography><Typography>Explore documented dependencies and their impact. These relationships describe architecture; they are not live health checks.</Typography>
    <Box my={2}><Button onClick={()=>setReload(x=>x+1)} disabled={loading}>Refresh catalog</Button> <Link href="/catalog-graph">Open full catalog graph</Link></Box>
    {loading&&<Typography>Loading catalog relationships…</Typography>}{error&&<Typography color="error">Catalog update failed. Any previous results may be stale.</Typography>}{partial&&<Typography color="error">View limited to 5,000 entities. Missing relationships may be outside this view.</Typography>}
    {loadedAt&&<Typography variant="caption">Catalog loaded {new Date(loadedAt).toLocaleString()}. Inventory observation dates below refer to discovery snapshots.</Typography>}
    <Box my={2}><TextField fullWidth label="Find a system, service, or resource" value={query} onChange={e=>setQuery(e.target.value)}/></Box>
    <TextField fullWidth select label="Explore entity" value={byRef.has(focus)?focus:''} onChange={e=>setFocus(e.target.value)}>{matches.map(e=><MenuItem key={ref(e)} value={ref(e)}>{label(ref(e))}</MenuItem>)}{entity&&!matches.includes(entity)&&<MenuItem value={focus}>{label(focus)}</MenuItem>}</TextField>
    <Box my={2}><TextField select label="Relationship direction" value={direction} onChange={e=>setDirection(e.target.value)}><MenuItem value="dependsOn">What this depends on</MenuItem><MenuItem value="dependencyOf">What depends on this</MenuItem></TextField></Box>
    {entity&&<><Typography variant="h5">{label(focus)}</Typography><Typography variant="caption">Relationship source: {entity.metadata.annotations?.['onemind.io/relationship-source']||'Catalog descriptor; review its source before relying on completeness'}</Typography>
      <InventoryState entity={entity}/>
      <Box my={2} display="flex" flexWrap="wrap" style={{gap:12}}>{entity.metadata.links?.map((x,index)=>{const href=safeMetadataLink(x.url);return href?<Link key={`${index}:${href}`} href={href} target="_blank" rel="noopener noreferrer">{x.title||'Related resource'}</Link>:<Typography key={index} variant="caption" color="error">Unsafe or invalid resource link omitted.</Typography>;})}</Box>
      {!first.length&&<Typography>No dependency relationships recorded in this direction. This does not prove there are no dependencies.</Typography>}
      {first.map(r=><Box key={r} my={2}><Typography variant="caption">{label(focus)} → {direction==='dependsOn'?'depends on':'is depended on by'}</Typography>{node(r)}<Box ml={3} mt={1}>{related(r).filter(next=>next!==focus&&next!==r).map(next=><Box key={next} mt={1}><Typography variant="caption">→ {direction==='dependsOn'?'depends on':'is depended on by'}</Typography>{node(next)}</Box>)}</Box></Box>)}
      <Typography variant="h6">Other recorded relationships</Typography>{entity.relations?.filter(x=>!['dependsOn','dependencyOf'].includes(x.type)).map(x=><Typography key={`${x.type}:${x.targetRef}`}>{x.type}: <Link href={catalogUrl(x.targetRef)}>{label(x.targetRef)}</Link></Typography>)}
    </>}
  </Box>
  );
}
export const systemsPlugin=createFrontendPlugin({pluginId:'onemind-systems',extensions:[PageBlueprint.make({params:{path:'/systems',title:'Systems',icon:<AccountTreeIcon/>,loader:async()=> <Systems/>}})]});

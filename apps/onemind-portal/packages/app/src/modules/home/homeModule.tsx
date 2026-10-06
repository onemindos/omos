import { useEffect, useState } from 'react';
import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { HomePageWidgetBlueprint } from '@backstage/plugin-home-react/alpha';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, TextField, Typography, Link, useTheme } from '@material-ui/core';
import directory from './directory.json';

function Directory(){
  const [query,setQuery]=useState(''); const theme=useTheme();
  const dark=theme.palette.type==='dark';
  const matches=directory.filter(x=>`${x.title} ${x.description} ${x.groupTitle}`.toLowerCase().includes(query.toLowerCase()));
  return <Box p={2}>
    <Box p={3} mb={2} style={{borderRadius:12,backgroundImage:`linear-gradient(${dark?'#111a':'#fffc'},${dark?'#111a':'#fffc'}),url(/brand/legacy-background-${dark?'dark':'light'}.webp)`,backgroundSize:'cover'}}>
      <Typography variant="h4">OneMind</Typography><Typography>Your systems. Your work. One home.</Typography>
    </Box>
    <TextField label="Find an app or capability" value={query} onChange={e=>setQuery(e.target.value)} fullWidth variant="outlined" size="small" />
    {[...new Set(matches.map(x=>x.groupTitle))].map(group=><Box key={group} mt={3}>
      <Typography variant="h6">{group}</Typography>
      <Box display="flex" flexWrap="wrap" style={{gap:12}}>
        {matches.filter(x=>x.groupTitle===group).map(x=><Box key={x.name} p={2} style={{border:'1px solid',borderColor:theme.palette.divider,borderRadius:8,flex:'1 1 180px'}}>
          <Link href={x.url} target="_blank" rel="noopener noreferrer" style={{fontWeight:600}}>{x.title}</Link>
          <Typography variant="body2" color="textSecondary">{x.description}</Typography>
          <Link href={`/catalog/default/component/${x.name}`} variant="caption">Service details</Link>
        </Box>)}
      </Box>
    </Box>)}
    {matches.length===0&&<Typography>No matching apps.</Typography>}
  </Box>;
}
type Status={observedAt:string;partial:boolean;totalWorkloads:number;unavailable:string[];deployments:{name:string;health:string;sync:string}[];backups:{name:string;lastSuccess:string|null;suspended:boolean}[]};
function StatusPanel(){
  const discovery=useApi(discoveryApiRef);const fetchApi=useApi(fetchApiRef);
  const [data,setData]=useState<Status>();const [error,setError]=useState(false);
  useEffect(()=>{let active=true;const controller=new AbortController();
    const refresh=async()=>{try{const base=await discovery.getBaseUrl('onemind');const r=await fetchApi.fetch(`${base}/status`,{signal:controller.signal});if(!r.ok)throw new Error();const value=await r.json();if(active){setData(value);setError(false);}}catch{if(active)setError(true);}};
    refresh();const timer=setInterval(refresh,60000);return()=>{active=false;controller.abort();clearInterval(timer);};
  },[discovery,fetchApi]);
  return <Box p={2}>
    {error&&<Typography color="error">Live update unavailable. Any previous results below may be stale.</Typography>}
    {!data&&!error&&<Typography>Loading cluster state…</Typography>}
    {data&&<>
      <Typography variant="caption">Checked {new Date(data.observedAt).toLocaleTimeString()} · {data.totalWorkloads} workloads</Typography>
      {data.partial&&<Typography color="error">Some sources are unavailable; this is an incomplete view.</Typography>}
      <Typography variant="h6">Needs attention</Typography>
      {!data.partial&&!data.unavailable.length&&!data.deployments.length&&<Typography>No workload readiness or deployment-sync issues reported.</Typography>}
      {data.unavailable.map(x=><Typography key={x} variant="body2">• Workload not fully ready: {x}</Typography>)}
      {data.deployments.map(x=><Typography key={x.name} variant="body2">• {x.name}: {x.health} / {x.sync}</Typography>)}
      <Box mt={2}><Link href="https://argocd.tail717ba0.ts.net" target="_blank" rel="noopener noreferrer">Open deployments</Link> · <Link href="https://grafana.tail717ba0.ts.net" target="_blank" rel="noopener noreferrer">Open Grafana</Link></Box>
      <Box mt={3}><Typography variant="h6">Backup activity</Typography>
        {data.backups.map(x=><Typography key={x.name} variant="body2">{x.name}: {x.suspended?'suspended':x.lastSuccess?`last success ${new Date(x.lastSuccess).toLocaleString()}`:'no successful run recorded'}</Typography>)}
      </Box>
    </>}
  </Box>;
}
export const homeModule=createFrontendModule({pluginId:'home',extensions:[
  HomePageWidgetBlueprint.make({name:'onemind-directory',params:{name:'OneMindDirectory',title:'Your apps',description:'Search the OneMind application directory',components:async()=>({Content:Directory})}}),
  HomePageWidgetBlueprint.make({name:'onemind-status',params:{name:'OneMindStatus',title:'Operations',description:'Live workload readiness, deployments and backup activity',components:async()=>({Content:StatusPanel})}}),
]});

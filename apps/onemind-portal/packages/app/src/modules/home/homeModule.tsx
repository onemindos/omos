import { useEffect, useState } from 'react';
import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { HomePageWidgetBlueprint } from '@backstage/plugin-home-react/alpha';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, TextField, Typography, Link, useTheme } from '@material-ui/core';
import directory from './directory.json';
import { ActivityInbox } from './ActivityInbox';

function Directory(){
  const [query,setQuery]=useState(''); const theme=useTheme();
  const dark=theme.palette.type==='dark';
  const matches=directory.filter(x=>`${x.title} ${x.description} ${x.groupTitle}`.toLowerCase().includes(query.toLowerCase()));
  return <Box p={2}>
    <Box p={3} mb={2} style={{borderRadius:12,backgroundImage:`linear-gradient(${dark?'#111a':'#fffc'},${dark?'#111a':'#fffc'}),url(/brand/legacy-background-${dark?'dark':'light'}.webp)`,backgroundSize:'cover'}}>
      <Typography variant="h4">OneMind</Typography><Typography>Your systems. Your work. One home.</Typography>
    </Box>
    <Box display="flex" flexWrap="wrap" mb={2} style={{gap:16}} aria-label="OneMind workspace">
      <Link href="/systems">Systems and dependencies</Link>
      <Link href="/create">Guided workflows</Link>
      <Link href="/docs">Knowledge</Link>
      <Link href="/assistant">Assistant</Link>
      <Link href="/offline">Offline knowledge</Link>
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
type ServiceHealth={observedAt:string;scope:string;metricsUnavailable:boolean;checks:{name:string;entity:string;link:string;scope:string;state:string;httpStatus?:number;latencyMs?:number;observedAt:string}[];metrics:{namespace:string;cpuCores:number|null;memoryBytes:number|null;observedAt:string}[]};
type Status={serviceHealth?:ServiceHealth;observedAt:string;partial:boolean;unavailableSources?:string[];totalWorkloads:number;unavailable:string[];deployments:{name:string;health:string;sync:string}[];backups:{name:string;lastSuccess:string|null;suspended:boolean}[]};
function StatusPanel(){
  const discovery=useApi(discoveryApiRef);const fetchApi=useApi(fetchApiRef);
  const [data,setData]=useState<Status>();const [error,setError]=useState(false);
  useEffect(()=>{let active=true;const controller=new AbortController();
    const refresh=async()=>{try{const base=await discovery.getBaseUrl('onemind');const r=await fetchApi.fetch(`${base}/status`,{signal:controller.signal});if(!r.ok)throw new Error();const value=await r.json();if(active){setData(value);setError(false);}}catch{if(active)setError(true);}};
    refresh();const timer=setInterval(refresh,60000);return()=>{active=false;controller.abort();clearInterval(timer);};
  },[discovery,fetchApi]);
  const theme=useTheme();
  const stale=!!data&&(error||Date.now()-Date.parse(data.observedAt)>180000);
  const backupProblems=data?.backups.filter(x=>x.suspended||!x.lastSuccess||(x.name==='backstage/backstage-backup'&&Date.now()-Date.parse(x.lastSuccess)>36*3600000))??[];
  const issues=data?[...(data.serviceHealth?.checks.filter(x=>x.state!=='up').map(x=>({key:`health:${x.name}`,title:`${x.name}: ${x.state}`,detail:`${x.scope}. Inspect service health and connectivity; unknown is not healthy.`,url:x.link}))??[]),...data.unavailable.map(name=>({key:`workload:${name}`,title:`Workload not fully ready: ${name}`,detail:'Inspect pod readiness and events in Headlamp.',url:'https://headlamp.tail717ba0.ts.net/c/main/workloads'})),
    ...data.deployments.map(x=>({key:`deployment:${x.name}`,title:`${x.name}: ${x.health} / ${x.sync}`,detail:'Inspect the deployment health, sync difference, and last operation.',url:`https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(x.name)}`})),
    ...backupProblems.map(x=>({key:`backup:${x.name}`,title:`Backup ${x.name}: ${x.suspended?'suspended':!x.lastSuccess?'no successful run recorded':'overdue'}`,detail:x.lastSuccess?`Last successful job: ${new Date(x.lastSuccess).toLocaleString()}`:'Inspect the backup schedule and job logs.',url:'https://headlamp.tail717ba0.ts.net/c/main/workloads'}))]:[];
  return <Box p={2}>
    <Typography variant="h5">Needs attention</Typography>
    <Typography variant="body2" color="textSecondary">Workload readiness, deployment health, and backup activity. Critical service checks and selected Prometheus metrics appear below.</Typography>
    {error&&<Typography color="error">Live update unavailable. Previous results may be stale.</Typography>}
    {!data&&!error&&<Typography>Loading cluster state…</Typography>}
    {data&&<>
      <Typography variant="caption">{stale?'Stale observation':'Checked'} {new Date(data.observedAt).toLocaleString()} · {data.totalWorkloads} workloads · refreshes every minute</Typography>
      {data.partial&&<Typography color="error">Incomplete view: {data.unavailableSources?.join(', ')||'some sources'} unavailable. Missing results do not mean healthy.</Typography>}
      <Box my={2}><Typography variant="h6">{issues.length} attention items{stale?' in previous results':''}</Typography></Box>
      {!stale&&!data.partial&&!issues.length&&<Typography>No issues reported by the configured checks.</Typography>}
      {issues.map(x=><Box key={x.key} p={2} mb={1} style={{border:`1px solid ${theme.palette.divider}`,borderRadius:8}}>
        <Typography style={{fontWeight:600}}>{x.title}</Typography>
        <Typography variant="body2" color="textSecondary">{x.detail}</Typography>
        <Link href={x.url} target="_blank" rel="noopener noreferrer">Investigate</Link>
      </Box>)}
      <Box mt={2}><Link href="https://argocd.tail717ba0.ts.net" target="_blank" rel="noopener noreferrer">All deployments</Link> · <Link href="https://grafana.tail717ba0.ts.net" target="_blank" rel="noopener noreferrer">Metrics in Grafana</Link></Box>
      {data.serviceHealth&&<Box mt={3}>
        <Typography variant="h6">Critical service health</Typography>
        <Typography variant="caption">{data.serviceHealth.scope} Checked {new Date(data.serviceHealth.observedAt).toLocaleString()}</Typography>
        {data.serviceHealth.checks.map(x=><Box key={x.name} mt={1}><Link href={x.link} target="_blank" rel="noopener noreferrer">{x.name}</Link> · {x.state}{x.httpStatus?` · HTTP ${x.httpStatus}`:''}{x.latencyMs!==undefined?` · ${x.latencyMs} ms`:''} · <Link href={`/catalog/${x.entity.split(':')[1].split('/')[0]}/${x.entity.split(':')[0]}/${x.entity.split('/')[1]}`}>Catalog</Link><Typography variant="caption" display="block">{x.scope}</Typography></Box>)}
        <Typography variant="h6">Namespace resource use</Typography>
        <Typography variant="caption">Prometheus: CPU cores averaged over 5 minutes; memory working set. Missing values mean unknown.</Typography>
        {data.serviceHealth.metricsUnavailable&&<Typography color="error">A metrics query failed; values may be incomplete.</Typography>}
        {data.serviceHealth.metrics.map(x=><Typography key={x.namespace} variant="body2">{x.namespace}: CPU {x.cpuCores===null?'unknown':`${x.cpuCores.toFixed(2)} cores`} · memory {x.memoryBytes===null?'unknown':`${(x.memoryBytes/1048576).toFixed(0)} MiB`}</Typography>)}
      </Box>}
      <Box mt={3}><Typography variant="h6">Backup activity</Typography>
        <Typography variant="caption">OneMind portal daily backup is overdue after 36 hours. Other schedules show activity without an assumed deadline. Job success alone does not prove restoration for every service.</Typography>
        {data.backups.map(x=><Typography key={x.name} variant="body2">{x.name}: {x.suspended?'suspended':x.lastSuccess?`last success ${new Date(x.lastSuccess).toLocaleString()}`:'no successful run recorded'}</Typography>)}
        {!data.backups.length&&!data.partial&&<Typography>No backup schedules found.</Typography>}
      </Box>
    </>}
  </Box>;
}
type CloudSummary={state:string;observedAt:string|null;expiresAt:string|null;open?:number;overdue?:number;sampled?:number;partial?:boolean;pages?:number;source?:string;coverage?:string;attemptedAt?:string;tasks?:{id:string;title:string;dueDate:string|null;overdue:boolean}[];link:string};
function CloudPanel(){
  const discovery=useApi(discoveryApiRef), api=useApi(fetchApiRef);
  const [data,setData]=useState<CloudSummary>(); const [error,setError]=useState(false);
  useEffect(()=>{let active=true;const controller=new AbortController();const refresh=async()=>{try{const base=await discovery.getBaseUrl('onemind');const r=await api.fetch(`${base}/cloud`,{signal:controller.signal});if(!r.ok)throw new Error();const result=await r.json();if(active){setData(result);setError(false);}}catch{if(active)setError(true);}};refresh();const timer=setInterval(refresh,60000);return()=>{active=false;controller.abort();clearInterval(timer);};},[discovery,api]);
  const expired=!!data?.expiresAt&&Date.parse(data.expiresAt)<Date.now();
  return <Box p={2}><Typography variant="h5">Cloud work</Typography><Typography variant="h6">Zoho Projects</Typography>
    {error&&<Typography color="error">Summary unavailable; previous observations may be stale.</Typography>}
    {!data&&!error&&<Typography>Loading…</Typography>}
    {data&&<><Typography>{data.state==='configuration-required'?'Connect read-only Zoho access to enable this summary':data.state==='unavailable'?'Provider unavailable':`${data.open} open · ${data.overdue} overdue in ${data.sampled} sampled tasks`}</Typography>
    {(expired||data.state==='stale')&&<Typography color="error">Stale observation</Typography>}
    {data.partial&&<Typography variant="caption">Bounded sample across {data.pages || 1} pages; these counts may not represent the whole portal.</Typography>}
    {data.observedAt&&<Typography variant="caption" display="block">Observed {new Date(data.observedAt).toLocaleString()} · cached for five minutes</Typography>}
    {data.coverage&&<Typography variant="caption" display="block">{data.coverage}</Typography>}
    {!!data.tasks?.length&&<Box mt={1}><Typography variant="subtitle2">Observed open work · up to 20 tasks, overdue first</Typography>{data.tasks.map(task=><Box key={task.id} mb={1}><Typography>{task.title}</Typography><Typography variant="caption" color={task.overdue?'error':'textSecondary'}>{task.overdue?'Overdue · ':''}{task.dueDate || 'No due date recorded'}</Typography></Box>)}</Box>}
    {data.source&&<Typography variant="caption" display="block">Source: {data.source}</Typography>}
    <Link href={data.link} target="_blank" rel="noopener noreferrer">Open Zoho Projects</Link></>}
  </Box>;
}
export const homeModule=createFrontendModule({pluginId:'home',extensions:[
  HomePageWidgetBlueprint.make({name:'onemind-activity',params:{name:'OneMindActivity',title:'Activity inbox',description:'Unread, actionable source notifications',components:async()=>({Content:ActivityInbox})}}),
  HomePageWidgetBlueprint.make({name:'onemind-cloud',params:{name:'OneMindCloud',title:'Cloud work',description:'Authorized provider summaries with observation times',components:async()=>({Content:CloudPanel})}}),
  HomePageWidgetBlueprint.make({name:'onemind-directory',params:{name:'OneMindDirectory',title:'Your apps',description:'Search the OneMind application directory',components:async()=>({Content:Directory})}}),
  HomePageWidgetBlueprint.make({name:'onemind-status',params:{name:'OneMindStatus',title:'Needs attention',description:'Live workload readiness, deployments and backup activity',components:async()=>({Content:StatusPanel})}}),
]});

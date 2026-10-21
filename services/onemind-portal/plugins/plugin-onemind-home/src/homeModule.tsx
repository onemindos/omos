import { ReactNode, useEffect, useId, useState } from 'react';
import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { HomePageLayoutBlueprint, HomePageWidgetBlueprint } from '@backstage/plugin-home-react/alpha';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, Button, Chip, TextField, Typography, Link, useTheme } from '@material-ui/core';
import { AppTile, appGridStyle, QuickAccess, useAppFavorites } from './Launcher';
import { useDirectory } from '@internal/plugin-onemind-common';
import { ActiveStack } from './ActiveStack';
import { hostingInfo } from '@internal/plugin-onemind-common';
import { CalendarPanel } from './CalendarPanel';
import { CostPanel } from './CostPanel';
import { FabricPanel } from './FabricPanel';
import { GrafanaHome } from './GrafanaHome';
import { LiveOverview } from './LiveOverview';
import { ServiceManager } from './ServiceManager';
import { groupColors } from '@internal/plugin-onemind-common';
import { ActivityInbox } from './ActivityInbox';
import { HomePageRequestedReviewsCard, HomePageYourOpenPullRequestsCard } from '@roadiehq/backstage-plugin-github-pull-requests';

function backupStatus(backup: {suspended:boolean;lastSuccess:string|null}) {
  if (backup.suspended) return 'suspended';
  if (!backup.lastSuccess) return 'no successful run recorded';
  return 'overdue';
}

function cloudStateLabel(data: {state:string;open?:number;overdue?:number;sampled?:number}) {
  if (data.state === 'configuration-required') return 'Connect read-only Zoho access to enable this summary';
  if (data.state === 'unavailable') return 'Provider unavailable';
  return `${data.open} open · ${data.overdue} overdue in ${data.sampled} sampled tasks`;
}

export function Directory({fluid=false,utilities}: {fluid?:boolean;utilities?:ReactNode}) {
  const {directory,loaded,error}=useDirectory();
  const searchId=useId();
  const [query, setQuery] = useState('');
  const [category, setCategory] = useState('All');
  const [favoritesOnly,setFavoritesOnly]=useState(false);
  const [hostingFilter, setHostingFilter] = useState('All hosting');
  const controls=useAppFavorites();
  const theme = useTheme();
  const dark = theme.palette.type === 'dark';
  const groups = [...new Set(directory.map(app => app.groupTitle))];
  const matches = directory.filter(app => (hostingFilter === 'All hosting' || hostingInfo(app.hosting, app.status).group === hostingFilter) && (!favoritesOnly || controls.isFavorite(app)) && (category === 'All' || app.groupTitle === category) &&
    `${app.title} ${app.description} ${app.groupTitle}`.toLowerCase().includes(query.trim().toLowerCase()));
  const destinations = [
    { title: 'Today', detail: 'Decide what needs you', url: '/today' },
    { title: 'Systems', detail: 'Explore dependencies', url: '/systems' },
    { title: 'Workflows', detail: 'Prepare a reviewed change', url: '/create' },
    { title: 'Knowledge', detail: 'Find a runbook', url: '/docs' },
    { title: 'Assistant', detail: 'Ask with source evidence', url: '/assistant' },
    { title: 'Offline', detail: 'Carry your knowledge', url: '/offline' },
  ];
  return (<Box p={fluid?0:2} className="onemind-workspace">
    <Box p={3} mb={3} style={{ borderRadius: 24, border: `1px solid ${theme.palette.divider}`,
      backgroundImage: `linear-gradient(115deg,${dark ? '#123951df,#35294add,#153c3ddb' : '#dceff8e8,#ece3f6df,#dff4ecdf'}),url(/brand/legacy-background-${dark ? 'dark' : 'light'}.webp)`, backgroundSize: 'cover', backgroundPosition: 'center' }}>
      <Typography variant="overline" style={{ letterSpacing: 3 }}>ONEMIND OS · YOUR WORKSPACE</Typography>
      <Typography variant="h4" component="h1" style={{ fontWeight: 700, margin: '12px 0' }}>Your OneMind OS home.</Typography>
      <Typography color="textSecondary">Your apps, infrastructure, agents, and knowledge in one place.</Typography>
      <Box mt={2}><Chip size="small" label={`${directory.length} directory entries`} /> <Typography variant="caption">Live status and alerts are in Today.</Typography></Box>
    </Box>
    <ServiceManager />
    {error&&<Typography color="error" role="alert">Catalog refresh failed. Previously loaded services may be stale; your personal shortcuts remain available.</Typography>}
    {!loaded&&<Typography role="status">Loading your catalog and shortcuts…</Typography>}
    <TextField id={searchId} label="Find an app or capability" value={query} onChange={event => setQuery(event.target.value)} fullWidth variant="outlined" size="small" />
    <Box my={2} display="flex" flexWrap="wrap" style={{ gap: 8 }} aria-label="Application categories">
      <Button size="small" variant={favoritesOnly?'contained':'outlined'} color="primary" aria-pressed={favoritesOnly} onClick={()=>setFavoritesOnly(value=>!value)}>Favorites only</Button>
      {['All', ...groups].map(group => <Button key={group} size="small" variant={category === group ? 'contained' : 'outlined'} color="primary" aria-pressed={category === group} onClick={() => setCategory(group)} style={{borderColor: `${groupColors[group]?.base??'#237f91'}88`, background:category===group?`${groupColors[group]?.base??'#237f91'}${dark?'44':'22'}`:undefined,color:dark?(groupColors[group]?.bright??theme.palette.text.primary):(groupColors[group]?.base??theme.palette.text.primary)}}>{group}</Button>)}
    </Box>
    <Box mb={2} display="flex" flexWrap="wrap" style={{gap:8}} aria-label="Hosting filters">{['All hosting', 'Self-hosted', 'Cloud service', 'Local / native', 'Reference', 'Planned', 'Unknown'].map(value => <Button key={value} size="small" variant={hostingFilter===value?'contained':'outlined'} aria-pressed={hostingFilter===value} onClick={()=>setHostingFilter(value)}>{value}</Button>)}</Box>
    {!query.trim() && category==='All' && !favoritesOnly && <><ActiveStack /><QuickAccess /></>}
    {!query.trim() && category==='All' && !favoritesOnly && <Box component="nav" aria-label="OneMind OS workspace" mb={3} style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(130px, 1fr))', gap: 10 }}>
      {destinations.map((item,index) => (<Button key={item.url} href={item.url} variant="outlined" style={{ justifyContent: 'flex-start', padding: 14, borderColor: `${Object.values(groupColors)[index].base}66`, background: `${Object.values(groupColors)[index].base}${dark?'22':'0b'}` }}>
        <Box textAlign="left"><Typography variant="subtitle2">{item.title}</Typography><Typography variant="caption" color="textSecondary">{item.detail}</Typography></Box>
      </Button>))}
    </Box>}
    {!query.trim() && category==='All' && !favoritesOnly && <><LiveOverview /><GrafanaHome />{utilities}</>}
    <Typography variant="h5" component="h2" style={{ fontWeight: 700, marginTop: 24 }}>Your application fabric</Typography>
    <Typography variant="body2" color="textSecondary" style={{ marginBottom: 16 }}>Choose an area, open an app, or inspect its catalog record.</Typography>
    {controls.error && <Typography color="error" role="alert">{controls.error}</Typography>}
    <Typography variant="caption" component="p" role="status">{matches.length} matching entries</Typography>
    {groups.filter(group => matches.some(app => app.groupTitle === group)).map(group => (<Box key={group} mt={3}>
      <Typography variant="h6" component="h3" style={{ marginBottom: 12, borderLeft:`4px solid ${groupColors[group]?.base??'#237f91'}`,paddingLeft:12,color:dark?(groupColors[group]?.bright??theme.palette.text.primary):(groupColors[group]?.base??theme.palette.text.primary) }}>{group}</Typography>
      <Box style={{...appGridStyle,gridTemplateColumns:'repeat(auto-fit, minmax(min(100%, 240px), 1fr))'}}>
        {matches.filter(app => app.groupTitle === group).map(app => (<Box key={app.ref}>
          <AppTile app={app} controls={controls} />
          {!app.personal&&<Box mt={1} ml={2}><Link href={app.personal?undefined:`/catalog/${app.ref.split(':')[1].split('/')[0]}/component/${app.name}`} variant="caption" aria-label={`Inspect ${app.title} in the catalog`}>Inspect service →</Link></Box>}
        </Box>))}
      </Box>
    </Box>))}
    {loaded && !error && !matches.length && <Box p={3} textAlign="center"><Typography variant="h6">No matching apps</Typography><Typography color="textSecondary">Try another term or reset your filters.</Typography><Button onClick={() => { setQuery(''); setCategory('All'); setFavoritesOnly(false); setHostingFilter('All hosting'); }}>Reset filters</Button></Box>}
  </Box>);
}
type ServiceHealth={observedAt:string;scope:string;metricsUnavailable:boolean;checks:{name:string;entity:string;link:string;scope:string;state:string;httpStatus?:number;latencyMs?:number;observedAt:string}[];metrics:{namespace:string;cpuCores:number|null;memoryBytes:number|null;observedAt:string}[]};
type Status={serviceHealth?:ServiceHealth;observedAt:string;partial:boolean;unavailableSources?:string[];totalWorkloads:number;unavailable:string[];deployments:{name:string;health:string;sync:string}[];backups:{name:string;lastSuccess:string|null;suspended:boolean}[]};
export function StatusPanel({ compact = false }: { compact?: boolean }){
  const discovery=useApi(discoveryApiRef);const fetchApi=useApi(fetchApiRef);
  const [data,setData]=useState<Status>();const [error,setError]=useState(false);
  useEffect(()=>{let active=true;const controller=new AbortController();
    const refresh=async()=>{try{const base=await discovery.getBaseUrl('onemind');const r=await fetchApi.fetch(`${base}/status`,{signal:controller.signal});if(!r.ok)throw new Error();const value=await r.json();if(active){setData(value);setError(false);}}catch{if(active)setError(true);}};
    refresh();const timer=setInterval(refresh,60000);return()=>{active=false;controller.abort();clearInterval(timer);};
  },[discovery,fetchApi]);
  const theme=useTheme();
  const [expanded,setExpanded]=useState(false);
  const abbreviated=compact&&!expanded;
  const stale=!!data&&(error||Date.now()-Date.parse(data.observedAt)>180000);
  const backupProblems=data?.backups.filter(x=>x.suspended||!x.lastSuccess||(x.name==='backstage/backstage-backup'&&Date.now()-Date.parse(x.lastSuccess)>36*3600000))??[];
  const issues=data?[...(data.serviceHealth?.checks.filter(x=>x.state!=='up').map(x=>({key:`health:${x.name}`,title:`${x.name}: ${x.state}`,detail:`${x.scope}. Inspect service health and connectivity; unknown is not healthy.`,url:x.link}))??[]),...data.unavailable.map(name=>({key:`workload:${name}`,title:`Workload not fully ready: ${name}`,detail:'Inspect pod readiness and events in Headlamp.',url:'https://headlamp.tail717ba0.ts.net/c/main/workloads'})),
    ...data.deployments.map(x=>({key:`deployment:${x.name}`,title:`${x.name}: ${x.health} / ${x.sync}`,detail:'Inspect the deployment health, sync difference, and last operation.',url:`https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(x.name)}`})),
    ...backupProblems.map(x=>({key:`backup:${x.name}`,title:`Backup ${x.name}: ${backupStatus(x)}`,detail:x.lastSuccess?`Last successful job: ${new Date(x.lastSuccess).toLocaleString()}`:'Inspect the backup schedule and job logs.',url:'https://headlamp.tail717ba0.ts.net/c/main/workloads'}))]:[];
  return (<Box p={2}>
    <Typography variant="h5">Needs attention</Typography>
    <Typography variant="body2" color="textSecondary">Workload readiness, deployment health, and backup activity. Critical service checks and selected Prometheus metrics appear below.</Typography>
    {error&&<Typography color="error">Live update unavailable. Previous results may be stale.</Typography>}
    {!data&&!error&&<Typography>Loading cluster state…</Typography>}
    {data&&<>
      <Typography variant="caption">{stale?'Stale observation':'Checked'} {new Date(data.observedAt).toLocaleString()} · {data.totalWorkloads} workloads · refreshes every minute</Typography>
      {data.partial&&<Typography color="error">Incomplete view: {data.unavailableSources?.join(', ')||'some sources'} unavailable. Missing results do not mean healthy.</Typography>}
      <Box my={2}><Typography variant="h6">{issues.length} attention items{stale?' in previous results':''}</Typography></Box>
      {!stale&&!data.partial&&!issues.length&&<Typography>No issues reported by the configured checks.</Typography>}
      {(abbreviated?issues.slice(0,5):issues).map(x=>(<Box key={x.key} p={2} mb={1} style={{border:`1px solid ${theme.palette.divider}`,borderRadius:8}}>
        <Typography style={{fontWeight:600}}>{x.title}</Typography>
        <Typography variant="body2" color="textSecondary">{x.detail}</Typography>
        <Link href={x.url} target="_blank" rel="noopener noreferrer">Investigate</Link>
      </Box>))}
      {compact && <Button size="small" onClick={()=>setExpanded(value=>!value)}>{expanded?'Show summary':'Show all observations'}</Button>}
      <Box mt={2}><Link href="https://argocd.tail717ba0.ts.net" target="_blank" rel="noopener noreferrer">All deployments</Link> · <Link href="https://grafana.tail717ba0.ts.net" target="_blank" rel="noopener noreferrer">Metrics in Grafana</Link></Box>
      {abbreviated && issues.length > 5 && <Typography variant="caption">Showing 5 of {issues.length} attention items. Open the full status view for the remaining observations.</Typography>}
      {!abbreviated && data.serviceHealth&&<Box mt={3}>
        <Typography variant="h6">Critical service health</Typography>
        <Typography variant="caption">{data.serviceHealth.scope} Checked {new Date(data.serviceHealth.observedAt).toLocaleString()}</Typography>
        {data.serviceHealth.checks.map(x=><Box key={x.name} mt={1}><Link href={x.link} target="_blank" rel="noopener noreferrer">{x.name}</Link> · {x.state}{x.httpStatus?` · HTTP ${x.httpStatus}`:''}{x.latencyMs!==undefined?` · ${x.latencyMs} ms`:''} · <Link href={`/catalog/${x.entity.split(':')[1].split('/')[0]}/${x.entity.split(':')[0]}/${x.entity.split('/')[1]}`}>Catalog</Link><Typography variant="caption" display="block">{x.scope}</Typography></Box>)}
        <Typography variant="h6">Namespace resource use</Typography>
        <Typography variant="caption">Prometheus: CPU cores averaged over 5 minutes; memory working set. Missing values mean unknown.</Typography>
        {data.serviceHealth.metricsUnavailable&&<Typography color="error">A metrics query failed; values may be incomplete.</Typography>}
        {data.serviceHealth.metrics.map(x=><Typography key={x.namespace} variant="body2">{x.namespace}: CPU {x.cpuCores===null?'unknown':`${x.cpuCores.toFixed(2)} cores`} · memory {x.memoryBytes===null?'unknown':`${(x.memoryBytes/1048576).toFixed(0)} MiB`}</Typography>)}
      </Box>}
      {!abbreviated && <Box mt={3}><Typography variant="h6">Backup activity</Typography>
        <Typography variant="caption">OneMind portal daily backup is overdue after 36 hours. Other schedules show activity without an assumed deadline. Job success alone does not prove restoration for every service.</Typography>
        {data.backups.map(x=><Typography key={x.name} variant="body2">{x.name}: {backupStatus(x)}</Typography>)}
        {!data.backups.length&&!data.partial&&<Typography>No backup schedules found.</Typography>}
      </Box>}
    </>}
  </Box>);
}
type CloudSummary={state:string;observedAt:string|null;expiresAt:string|null;open?:number;overdue?:number;sampled?:number;partial?:boolean;pages?:number;source?:string;coverage?:string;attemptedAt?:string;tasks?:{id:string;title:string;dueDate:string|null;overdue:boolean}[];link:string};
export function CloudPanel(){
  const discovery=useApi(discoveryApiRef); const api=useApi(fetchApiRef);
  const [data,setData]=useState<CloudSummary>(); const [error,setError]=useState(false);
  useEffect(()=>{let active=true;const controller=new AbortController();const refresh=async()=>{try{const base=await discovery.getBaseUrl('onemind');const r=await api.fetch(`${base}/cloud`,{signal:controller.signal});if(!r.ok)throw new Error();const result=await r.json();if(active){setData(result);setError(false);}}catch{if(active)setError(true);}};refresh();const timer=setInterval(refresh,60000);return()=>{active=false;controller.abort();clearInterval(timer);};},[discovery,api]);
  const expired=!!data?.expiresAt&&Date.parse(data.expiresAt)<Date.now();
  return (<Box p={2}><Typography variant="h5">Cloud work</Typography><Typography variant="h6">Zoho Projects</Typography>
    {error&&<Typography color="error">Summary unavailable; previous observations may be stale.</Typography>}
    {!data&&!error&&<Typography>Loading…</Typography>}
    {data&&<><Typography>{cloudStateLabel(data)}</Typography>
    {(expired||data.state==='stale')&&<Typography color="error">Stale observation</Typography>}
    {data.partial&&<Typography variant="caption">Bounded sample across {data.pages || 1} pages; these counts may not represent the whole portal.</Typography>}
    {data.observedAt&&<Typography variant="caption" display="block">Observed {new Date(data.observedAt).toLocaleString()} · cached for five minutes</Typography>}
    {data.coverage&&<Typography variant="caption" display="block">{data.coverage}</Typography>}
    {!!data.tasks?.length&&<Box mt={1}><Typography variant="subtitle2">Observed open work · up to 20 tasks, overdue first</Typography>{data.tasks.map(task=><Box key={task.id} mb={1}><Typography>{task.title}</Typography><Typography variant="caption" color={task.overdue?'error':'textSecondary'}>{task.overdue?'Overdue · ':''}{task.dueDate || 'No due date recorded'}</Typography></Box>)}</Box>}
    {data.source&&<Typography variant="caption" display="block">Source: {data.source}</Typography>}
    <Link href={data.link} target="_blank" rel="noopener noreferrer">Open Zoho Projects</Link></>}
  </Box>);
}
export const homeModule=createFrontendModule({pluginId:'home',extensions:[
  HomePageLayoutBlueprint.make({name:'onemind-flow',params:{loader:async()=>(await import('./HomeLayout')).HomeLayout}}),
  HomePageWidgetBlueprint.make({name:'onemind-calendar',params:{name:'OneMindCalendar',title:'Upcoming calendar events',description:'Selected Zoho calendar with read-only access',components:async()=>({Content:CalendarPanel})}}),
  HomePageWidgetBlueprint.make({name:'onemind-costs',params:{name:'OneMindCosts',title:'DigitalOcean account costs',description:'Read-only account billing totals; administrator access required',components:async()=>({Content:CostPanel})}}),
  HomePageWidgetBlueprint.make({name:'onemind-fabric',params:{name:'OneMindFabric',title:'Agent fabric · NATS',description:'Read-only NATS monitoring; agent completion source readiness',components:async()=>({Content:FabricPanel})}}),
  HomePageWidgetBlueprint.make({name:'onemind-grafana',params:{name:'OneMindGrafana',title:'Live Grafana dashboards',description:'Your account-selected dashboards, loaded on demand',components:async()=>({Content:GrafanaHome})}}),
  HomePageWidgetBlueprint.make({name:'onemind-live-overview',params:{name:'OneMindLiveOverview',title:'Live overview',description:'Timestamped workloads, deployments, backup activity and service checks',components:async()=>({Content:LiveOverview})}}),
  HomePageWidgetBlueprint.make({name:'onemind-favorites',params:{name:'OneMindFavorites',title:'Favorite apps',description:'Direct app launch tiles using your existing catalog stars',components:async()=>({Content:QuickAccess})}}),
  HomePageWidgetBlueprint.make({name:'onemind-today',params:{name:'OneMindToday',title:'Today',description:'Live attention, work and inbox in one operational workspace',components:async()=>({Content:(await import('./Today')).Today})}}),
  HomePageWidgetBlueprint.make({name:'github-reviews',params:{name:'OneMindGitHubReviews',title:'GitHub reviews',description:'Requested reviews in onemindos; connect GitHub authorization',components:async()=>({Content:()=> <HomePageRequestedReviewsCard query="org:onemindos is:pr is:open review-requested:@me" />})}}),
  HomePageWidgetBlueprint.make({name:'github-open-prs',params:{name:'OneMindGitHubPullRequests',title:'GitHub pull requests',description:'Your open pull requests in onemindos',components:async()=>({Content:()=> <HomePageYourOpenPullRequestsCard query="org:onemindos is:pr is:open author:@me" />})}}),
  HomePageWidgetBlueprint.make({name:'onemind-activity',params:{name:'OneMindActivity',title:'Activity inbox',description:'Unread, actionable source notifications',components:async()=>({Content:ActivityInbox})}}),
  HomePageWidgetBlueprint.make({name:'onemind-cloud',params:{name:'OneMindCloud',title:'Cloud work',description:'Authorized provider summaries with observation times',components:async()=>({Content:CloudPanel})}}),
  HomePageWidgetBlueprint.make({name:'onemind-directory',params:{name:'OneMindDirectory',title:'Your apps',description:'Search the OneMind application directory',components:async()=>({Content:Directory})}}),
  HomePageWidgetBlueprint.make({name:'onemind-status',params:{name:'OneMindStatus',title:'Needs attention',description:'Live workload readiness, deployments and backup activity',components:async()=>({Content:StatusPanel})}}),
]});

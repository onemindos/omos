import { useEffect, useId, useState } from 'react';
import { Avatar, Box, Button, Dialog, DialogContent, DialogTitle, IconButton, Link, TextField, Tooltip, Typography, useTheme } from '@material-ui/core';
import Star from '@material-ui/icons/Star';
import StarBorder from '@material-ui/icons/StarBorder';
import Close from '@material-ui/icons/Close';
import { storageApiRef, useApi } from '@backstage/core-plugin-api';
import { starredEntitiesApiRef, useStarredEntities } from '@backstage/plugin-catalog-react';
import { appColors } from './colors';
import { AppEntry, useDirectory } from './useDirectory';

const abbreviation = (app: AppEntry) => ({ 'github': 'GH', '1password': '1P', 'argocd': 'CD', 'digitalocean': 'DO', 'grafana': 'GF', 'headlamp': 'HL' }[app.name] ?? (app.title.match(/[A-Za-z0-9]+/g) ?? []).map((word, index, words) => words.length === 1 ? word.slice(0, 2) : index < 2 ? word[0] : '').join('').toUpperCase());

export function useAppFavorites() {
  const {directory,all,personal}=useDirectory();
  const entityRef=(app:AppEntry)=>app.ref;
  const aliases=(app:AppEntry)=>all.filter(other=>other.url===app.url&&!other.personal).map(entityRef);
  const api = useApi(starredEntitiesApiRef);
  const { starredEntities } = useStarredEntities();
  const settings=useApi(storageApiRef).forBucket('starredEntities');
  const [loaded,setLoaded]=useState(settings.snapshot('entityRefs').presence!=='unknown');
  useEffect(()=>{const subscription=settings.observe$('entityRefs').subscribe({next:value=>setLoaded(value.presence!=='unknown')});return()=>subscription.unsubscribe();},[settings]);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const isFavorite = (app: AppEntry) => app.personal?!!app.favorite:aliases(app).some(ref => starredEntities.has(ref));
  const toggle = async (app: AppEntry) => {
    if (busy || error || !loaded) return;
    setBusy(true);
    try {
      if(app.personal){await personal.save(personal.shortcuts.map(item=>item.name===app.name?{...item,favorite:!item.favorite}:item));return;}
      const refs = aliases(app).filter(ref => starredEntities.has(ref));
      if (refs.length) { for (const ref of refs) await api.toggleStarred(ref); }
      else await api.toggleStarred(entityRef(app));
    } catch { setError('Favorites could not be saved. Reload the page before trying again.'); }
    finally { setBusy(false); }
  };
  // The existing star set supplies personal order; catalog stars outside the app
  // directory remain untouched and are available in the native Starred Entities widget.
  const favorites = [...starredEntities].flatMap(ref => {
    const match = directory.find(app => entityRef(app) === ref);
    return match ? [directory.find(app => app.url === match.url)!] : [];
  }).filter((app, index, apps) => apps.findIndex(other => other.url === app.url) === index);
  return { favorites:[...favorites,...personal.shortcuts.filter(app=>app.favorite)], isFavorite, toggle, busy, error, loaded:loaded&&personal.loaded };
}
type FavoriteControls = ReturnType<typeof useAppFavorites>;
export function AppTile({ app, controls, compact = false }: { app: AppEntry; controls: FavoriteControls; compact?: boolean }) {
  const theme = useTheme();
  const starred = controls.isFavorite(app);
  const dark=theme.palette.type==='dark';
  const accent=appColors(app.name,app.groupTitle);
  return <Box className="onemind-app-card" style={{ display: 'flex', alignItems: 'flex-start', border: `1px solid ${accent.base}${dark?'66':'44'}`, borderTop:`3px solid ${accent.base}`, borderRadius: 16, background: `linear-gradient(135deg,${accent.base}${dark?'25':'0f'},${theme.palette.background.paper})`, minWidth: 0 }}>
    <Link href={app.url} target="_blank" rel="noopener noreferrer" aria-label={`Open ${app.title}`} style={{ display: 'flex', gap: compact ? 8 : 12, padding: compact ? 12 : 16, flex: 1, minWidth: 0, textDecoration: 'none', color: theme.palette.text.primary }}>
      <Avatar variant="rounded" style={{ background: `linear-gradient(135deg,${accent.base},${accent.base}cc)`, color: '#fff', flexShrink: 0, fontSize: compact ? 13 : 15, width: compact ? 32 : 40, height: compact ? 32 : 40, fontWeight: 700 }}>{abbreviation(app)}</Avatar>
      <Box style={{ minWidth: 0 }}><Typography variant="subtitle2" style={{ fontWeight: 700, overflowWrap: 'break-word' }}>{app.title} ↗</Typography>
        <Typography variant="caption" color="textSecondary" display="block">{compact ? app.groupTitle : app.description}</Typography>
        {!compact && <Typography variant="caption" color="textSecondary" display="block">{app.groupTitle}</Typography>}
      </Box>
    </Link>
    <Tooltip title={starred ? 'Remove from favorites' : 'Add to favorites'}><span><IconButton size="small" aria-label={`${starred ? 'Remove' : 'Add'} ${app.title} ${starred ? 'from' : 'to'} favorites`} aria-pressed={starred} disabled={!controls.loaded || controls.busy || !!controls.error} onClick={() => controls.toggle(app)} style={{ margin: '8px 6px 0 0',minWidth:44,minHeight:44,color:starred?(dark?'#ffcf68':'#926400'):dark?'#afbdcf':'#64748b' }}>{starred ? <Star fontSize="small" /> : <StarBorder fontSize="small" />}</IconButton></span></Tooltip>
  </Box>;
}
export const appGridStyle = { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(min(100%, 240px), 1fr))', gap: 12 };

export function QuickAccess() {
  const searchId=useId();
  const titleId=useId();
  const {directory:launcherApps}=useDirectory();
  const controls = useAppFavorites();
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState('');
  const matches = launcherApps.filter(app => `${app.title} ${app.description} ${app.groupTitle}`.toLowerCase().includes(query.trim().toLowerCase()));
  const starters = ['zoho-mail', 'github', 'grafana', 'headlamp', '1password', 'mercury'].flatMap(name => launcherApps.find(app => app.name === name) ?? []);
  return <Box className="onemind-workspace onemind-favorites">
    <Box display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap" style={{ gap: 12 }} mb={1}>
      <Box><Typography variant="h5" component="h2" style={{ fontWeight: 700 }}>Favorites</Typography><Typography variant="body2" color="textSecondary">Your apps, one click away. Stars are shared with your catalog favorites.</Typography></Box>
      <Button variant="contained" color="primary" disableElevation onClick={() => { setQuery(''); setOpen(true); }}>Manage favorites</Button>
    </Box>
    {controls.error && <Typography color="error" role="alert">{controls.error}</Typography>}
    {!controls.loaded && <Typography role="status">Loading saved favorites…</Typography>}
    {controls.loaded && (controls.favorites.length ? <Box mt={2} style={appGridStyle}>{controls.favorites.map(app => <AppTile key={app.name} app={app} controls={controls} compact />)}</Box> : <Box mt={2}><Typography>No favorite apps yet. Use a star below or choose apps in Manage favorites.</Typography><Typography variant="overline" component="p" style={{ marginTop: 16 }}>Suggested starting points</Typography><Box style={appGridStyle}>{starters.map(app => <AppTile key={app.name} app={app} controls={controls} compact />)}</Box></Box>)}
    <Dialog open={open} onClose={() => setOpen(false)} fullWidth maxWidth="md" aria-labelledby={titleId}>
      <DialogTitle><Box display="flex" justifyContent="space-between" alignItems="center"><span id={titleId}>Choose your quick-access apps</span><IconButton aria-label="Close favorites picker" onClick={() => setOpen(false)}><Close /></IconButton></Box></DialogTitle>
      <DialogContent><TextField id={searchId} autoFocus fullWidth label="Find an app to favorite" value={query} onChange={event => setQuery(event.target.value)} variant="outlined" size="small" />
        <Typography variant="caption" component="p" role="status" style={{ margin: '12px 0' }}>{matches.length} apps · {controls.favorites.length} favorites</Typography>
        {controls.error && <Typography color="error" role="alert">{controls.error}</Typography>}
        <Box mb={3} style={appGridStyle}>{matches.map(app => <AppTile key={app.name} app={app} controls={controls} compact />)}</Box>
        {!matches.length && <Typography>No matching apps. Try another name or category.</Typography>}
      </DialogContent>
    </Dialog>
  </Box>;
}

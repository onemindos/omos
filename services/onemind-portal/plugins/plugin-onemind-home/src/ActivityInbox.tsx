import { useEffect, useState } from 'react';
import { useApi } from '@backstage/core-plugin-api';
import { notificationsApiRef, GetNotificationsResponse } from '@backstage/plugin-notifications';
import { Box, Button, Link, Typography } from '@material-ui/core';

export function ActivityInbox() {
  const api = useApi(notificationsApiRef);
  const [data, setData] = useState<GetNotificationsResponse>();
  const [error, setError] = useState(false);
  const [refresh, setRefresh] = useState(0);
  useEffect(() => {
    let active = true;
    const load = async () => { try { const result = await api.getNotifications({ read: false, limit: 10, sort: 'created', sortOrder: 'desc' }); if (active) { setData(result); setError(false); } } catch { if (active) setError(true); } };
    load(); const timer = setInterval(load, 30000);
    return () => { active = false; clearInterval(timer); };
  }, [api, refresh]);
  return (<Box p={2}><Typography variant="h5">Activity inbox</Typography>
    <Typography variant="body2">Your unread notifications. Source events are historical reports, not proof of current health.</Typography>
    {error && <Typography color="error">Inbox unavailable. Previous results may be stale.</Typography>}
    {!data && !error && <Typography>Loading notifications…</Typography>}
    {data && <><Typography variant="caption">{data.totalCount} unread · latest 10 · refreshes every 30 seconds</Typography>
      {!data.totalCount && <Typography>No unread notifications. Producers must be connected before infrastructure events appear.</Typography>}
      {data.notifications.map(item => (<Box key={item.id} mt={2}><Typography style={{ fontWeight: 600 }}>{item.payload.title}</Typography>
        <Typography variant="caption">{item.payload.severity} · {new Date(item.created).toLocaleString()}</Typography>
        <Typography variant="body2">{item.payload.description}</Typography>
        {item.payload.link && <Link href={item.payload.link}>Investigate</Link>}
        <Button size="small" onClick={async () => { try { await api.updateNotifications({ ids: [item.id], read: true }); setRefresh(x => x + 1); } catch { setError(true); } }}>Mark read</Button>
      </Box>))}
    </>}
    <Box mt={2}><Link href="/notifications">All notifications and saved items</Link></Box>
  </Box>);
}

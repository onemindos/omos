import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Box, Link, Typography } from '@material-ui/core';
type Summary = {
  state: string;
  coverage: string;
  observedAt: string | null;
  expiresAt: string | null;
  events: {
    id: string;
    title: string;
    start: string | null;
    allDay: boolean;
  }[];
  link: string;
};
export function CalendarPanel() {
  const discovery = useApi(discoveryApiRef),
    api = useApi(fetchApiRef),
    [data, setData] = useState<Summary>(),
    [error, setError] = useState('');
  useEffect(() => {
    let active = true,
      pending = false;
    const controller = new AbortController();
    const refresh = async () => {
      if (pending) return;
      pending = true;
      try {
        const base = await discovery.getBaseUrl('onemind');
        const response = await api.fetch(base + '/calendar', {
          signal: controller.signal,
        });
        if (!response.ok)
          throw new Error(
            response.status === 403
              ? 'Administrator access required'
              : 'Calendar observation unavailable',
          );
        const value = await response.json();
        if (active) {
          setData(value);
          setError('');
        }
      } catch (e) {
        if (active)
          setError(e instanceof Error ? e.message : 'Calendar unavailable');
      } finally {
        pending = false;
      }
    };
    refresh();
    const timer = setInterval(refresh, 300000);
    return () => {
      active = false;
      controller.abort();
      clearInterval(timer);
    };
  }, [discovery, api]);
  const stale =
    !!data &&
    (!!error ||
      data.state === 'stale' ||
      (!!data.expiresAt && Date.parse(data.expiresAt) < Date.now()));
  return (
    <Box p={2}>
      <Typography variant="h5">Upcoming · Zoho Calendar</Typography>
      {error && (
        <Typography role="alert" color="error">
          {error}
        </Typography>
      )}
      {!data && !error && <Typography>Loading calendar…</Typography>}
      {data && (
        <>
          <Typography variant="caption" display="block">
            {data.coverage}
          </Typography>
          {data.state === 'configuration-required' && (
            <Typography>
              Connect read-only Calendar access and select a calendar UID to
              enable upcoming events.
            </Typography>
          )}
          {data.state === 'unavailable' && (
            <Typography color="error">
              Calendar source unavailable; an empty list does not mean no
              events.
            </Typography>
          )}
          {data.observedAt && (
            <>
              <Typography variant="caption" display="block">
                {stale ? 'Stale · ' : ''}Observed{' '}
                {new Date(data.observedAt).toLocaleString()} · cached five
                minutes
              </Typography>
              {data.events.map(event => (
                <Box key={event.id} mt={1}>
                  <Typography>{event.title}</Typography>
                  <Typography variant="caption">
                    {!event.start
                      ? 'Time unavailable'
                      : event.allDay
                      ? event.start + ' · all day'
                      : new Date(event.start).toLocaleString()}
                  </Typography>
                </Box>
              ))}
              {!data.events.length && (
                <Typography>
                  No events returned for the selected window and authorized
                  calendar.
                </Typography>
              )}
            </>
          )}
          <Box mt={2}>
            <Link href={data.link} target="_blank" rel="noopener noreferrer">
              Open Zoho Calendar
            </Link>
          </Box>
        </>
      )}
    </Box>
  );
}

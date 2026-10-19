import { useState } from 'react';
import { Box, Button, Typography } from '@material-ui/core';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';

export function OfflinePage() {
  const discovery = useApi(discoveryApiRef);
  const api = useApi(fetchApiRef);
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState('');

  const download = async () => {
    setBusy(true);
    setMessage('');

    try {
      const base = await discovery.getBaseUrl('onemind-offline');
      const response = await api.fetch(`${base}/pack`, {
        signal: AbortSignal.timeout(15000),
      });

      if (!response.ok) {
        throw new Error(
          response.status === 403
            ? 'Administrator access required.'
            : 'Pack unavailable. Try again while connected.',
        );
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'onemind-offline.html';
      document.body.appendChild(link);
      link.click();
      link.remove();
      setTimeout(() => URL.revokeObjectURL(url), 60000);
      setMessage(
        'Private pack downloaded. Open the HTML file locally to read it without a connection.',
      );
    } catch (e) {
      setMessage(e instanceof Error ? e.message : 'Download unavailable');
    } finally {
      setBusy(false);
    }
  };

  return (
    <Box p={3}>
      <Typography variant="h4">Offline knowledge</Typography>
      <Box my={2}>
        <Typography>
          Keep a private copy of your application directory and operational runbooks for outages or travel.
        </Typography>
        <Typography>
          This is a dated document pack. Live services, cloud apps, authentication, tasks and the dashboard itself still require their normal connections.
        </Typography>
      </Box>
      <Box my={2}>
        <Typography>
          The file contains private system documentation and account-specific links. Store it on a protected device. Logging out does not remove it; delete older copies deliberately.
        </Typography>
      </Box>
      <Button color="primary" variant="contained" disabled={busy} onClick={download}>
        {busy ? 'Preparing…' : 'Download private offline pack'}
      </Button>
      <Box my={2}>
        <Typography role="status">{message}</Typography>
      </Box>
    </Box>
  );
}

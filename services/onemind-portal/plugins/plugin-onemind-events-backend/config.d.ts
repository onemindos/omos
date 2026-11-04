export interface Config {
  onemind?: {
    alerts?: {
      /**
       * Bearer token Alertmanager must present at POST /api/events/http/alertmanager.
       * Optional: when unset or shorter than 24 characters the route rejects every request.
       * @visibility secret
       */
      webhookToken?: string;
    };
    operationalEvents?: {
      /**
       * CronJobs in the backstage namespace whose failures and suspension raise alerts.
       * @visibility backend
       */
      backupCronJobs?: string[];
      /**
       * UTC hour after which the daily deployment digest is sent. Defaults to 13.
       * @visibility backend
       */
      digestHourUtc?: number;
    };
  };
}

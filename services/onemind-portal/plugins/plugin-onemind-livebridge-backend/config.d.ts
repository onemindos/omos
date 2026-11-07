export interface Config {
  onemind?: {
    livebridge?: {
      nats?: {
        /**
         * NATS client URL of the authority bus. Defaults to nats://nats.nats.svc.cluster.local:4222.
         * @visibility backend
         */
        url?: string;
        /**
         * NATS user (env NATS_USER). Without a user and password the NATS source stays disabled.
         * @visibility secret
         */
        user?: string;
        /**
         * NATS password (env NATS_PASSWORD).
         * @visibility secret
         */
        password?: string;
        /**
         * Optional CA file to verify a TLS NATS listener.
         * @visibility backend
         */
        caFile?: string;
      };
      kubernetes?: {
        /**
         * Watch pods, Warning events and Argo applications. Defaults to true.
         * @visibility backend
         */
        enabled?: boolean;
        /**
         * Namespace of the Argo CD Application resources. Defaults to argocd.
         * @visibility backend
         */
        argoNamespace?: string;
      };
    };
  };
}

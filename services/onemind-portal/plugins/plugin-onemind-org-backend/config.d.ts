export interface Config {
  authentik?: {
    /**
     * Authentik base URL for the read-only org sync, e.g. the in-cluster service
     * http://platform-auth-helm-authentik-server.auth.svc.cluster.local. Required when the sync is enabled.
     * @visibility backend
     */
    baseUrl?: string;
    /**
     * Authentik API token of a service account that may only view users and groups.
     * Required when the sync is enabled; read from the environment, never committed.
     * @visibility secret
     */
    token?: string;
    orgSync?: {
      /**
       * Publish Authentik users and groups as catalog User and Group entities. Defaults to false.
       * @visibility backend
       */
      enabled?: boolean;
      /**
       * Exact Authentik group names. Only active members of these groups (or of their child groups)
       * become catalog Users, and so can sign in. Required and non-empty when the sync is enabled.
       * @visibility backend
       */
      userGroups?: string[];
      /**
       * Seconds between syncs. Defaults to 600.
       * @visibility backend
       */
      frequencySeconds?: number;
      /**
       * Seconds one sync may take. Defaults to 120.
       * @visibility backend
       */
      timeoutSeconds?: number;
      /**
       * Seconds before the first sync after startup. Defaults to 60.
       * @visibility backend
       */
      initialDelaySeconds?: number;
      /**
       * Reviewed catalog file whose Users and Groups always win and are never published or linked.
       * Defaults to catalog/org.yaml.
       * @visibility backend
       */
      curatedFile?: string;
    };
  };
}

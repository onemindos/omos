export interface Config {
  authentik?: {
    /**
     * Authentik base URL for the read-only org sync, e.g. the in-cluster service
     * http://platform-auth-helm-authentik-server.auth.svc.cluster.local. Required when the sync is enabled.
     * @visibility backend
     */
    baseUrl?: string;
    /**
     * Authentik API token of a service account that may only view users and groups. Read from the
     * environment, never committed. If it is unset while the sync is enabled, startup does not fail:
     * the module calls nothing and removes the users and groups it published.
     * @visibility secret
     */
    token?: string;
    orgSync?: {
      /**
       * Publish Authentik users and groups as catalog User and Group entities. Defaults to false.
       * While false the module calls nothing and only removes entities it published earlier.
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
       * Seconds a starting pod waits before it may run the task. Defaults to 60. It does not delay
       * the schedule: the scheduler adds it to the next run only when the task is first stored, and
       * the disabled module stores the same task, so the first sync after enabling comes at the
       * task's next scheduled run (within 10 minutes, or frequencySeconds if lower).
       * @visibility backend
       */
      initialDelaySeconds?: number;
      /**
       * Reviewed catalog file whose Users and Groups always win and are never published or linked.
       * Defaults to the one file location in catalog.locations named org.yaml; a relative path
       * resolves against the working directory, as the catalog resolves that location.
       * @visibility backend
       */
      curatedFile?: string;
    };
  };
}

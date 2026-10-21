export interface Config {
  onemind?: {
    openchoreo?: {
      /** Explicit trusted direct API URL prefixes. Never derived from untrusted catalog input.
       * @visibility frontend
       */
      directApiAllowlist?: string[];
    };
  };
}

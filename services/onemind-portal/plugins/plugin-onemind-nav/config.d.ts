export interface Config {
  onemind?: {
    /**
     * Navigation and home sections that depend on later flagship work.
     * Each defaults to false so unfinished destinations never show.
     */
    features?: {
      /**
       * Pillars sidebar group and home tiles. Turn on after F2 creates the
       * SO, LE, HP and GE Domain entities.
       * @visibility frontend
       */
      pillars?: boolean;
      /**
       * DigitalOcean sidebar link. Turn on after F6 ships /digitalocean.
       * @visibility frontend
       */
      digitalocean?: boolean;
      /**
       * Retire app sidebar link. Turn on after F8 ships the retire template.
       * @visibility frontend
       */
      retireApp?: boolean;
    };
  };
}

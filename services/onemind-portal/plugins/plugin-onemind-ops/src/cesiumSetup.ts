/**
 * Cesium needs its static assets (Workers, Assets, ThirdParty, Widgets) at a known URL. They are copied to
 * `packages/app/public/cesium` by `scripts/copy-cesium-assets.cjs` before every app build and start, so they
 * are served from `<app base path>/cesium/`.
 */
export const CESIUM_DIR = 'cesium';

/** `https://portal.example.com/` -> `/cesium/`; `https://host/portal` -> `/portal/cesium/`. */
export function cesiumBaseUrl(appBaseUrl: string): string {
  let path = '/';
  try {
    path = new URL(appBaseUrl).pathname;
  } catch {
    path = '/';
  }
  return `${path.replace(/\/+$/, '')}/${CESIUM_DIR}/`;
}

/**
 * Imagery without an Ion token: the standard OpenStreetMap tile server, dimmed for the dark page.
 * (CARTO's free dark tiles now answer "API KEY REQUIRED", so they are not used.) The OSM tile usage policy
 * applies: light use, visible credit, and the browser must send a Referer. If the tile host changes, update
 * the CSP `img-src` and `connect-src` entries in app-config.yaml and platform/backstage/app-config.cluster.yaml.
 */
export const IMAGERY = {
  url: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
  maximumLevel: 19,
  credit: '© OpenStreetMap contributors',
  /** Dim and desaturate the light map so tracks stand out on the dark page. */
  layer: { brightness: 0.55, contrast: 1.1, saturation: 0.4 },
} as const;

/** Adds the Cesium widget stylesheet (credits, canvas layout) once. */
export function ensureWidgetStyles(
  baseUrl: string,
  doc: Document = document,
): void {
  const id = 'onemind-ops-cesium-widgets';
  if (doc.getElementById(id)) return;
  const link = doc.createElement('link');
  link.id = id;
  link.rel = 'stylesheet';
  link.href = `${baseUrl}Widgets/widgets.css`;
  doc.head.appendChild(link);
}

import { useEffect, useRef, useState } from 'react';
import { configApiRef, useApi } from '@backstage/core-plugin-api';
import { tokens } from '@internal/plugin-onemind-theme';
import type { Track } from './contract';
import { IMAGERY, cesiumBaseUrl, ensureWidgetStyles } from './cesiumSetup';
import { createTrackLayer, drawArrow } from './trackLayer';
import type { TrackLayer } from './trackLayer';

export interface CesiumGlobeProps {
  /** Tracks to draw (already filtered). */
  tracks: readonly Track[];
  /** Quantised clock for the fade, so the layer is not rewritten every second. */
  fadeNowMs: number;
  selectedUid?: string;
  onSelect: (uid: string | undefined) => void;
  /** Changing `nonce` flies the camera to the track. */
  flyTo?: { uid: string; lat: number; lon: number; nonce: number };
  onError?: (message: string) => void;
}

/**
 * The only module that imports `cesium`, and it does so with a dynamic import after `CESIUM_BASE_URL` is set,
 * so the library and its workers load only when the Ops Center is opened.
 */
export default function CesiumGlobe(props: CesiumGlobeProps) {
  const { tracks, fadeNowMs, selectedUid, onSelect, flyTo, onError } = props;
  const config = useApi(configApiRef);
  const host = useRef<HTMLDivElement>(null);
  const layer = useRef<TrackLayer>();
  const viewerRef = useRef<import('cesium').Viewer>();
  const cesiumRef = useRef<typeof import('cesium')>();
  const [ready, setReady] = useState(false);
  const latest = useRef({ tracks, fadeNowMs, onSelect });
  latest.current = { tracks, fadeNowMs, onSelect };

  useEffect(() => {
    let cancelled = false;
    let cleanup = () => {};
    (async () => {
      try {
        const base = cesiumBaseUrl(
          config.getOptionalString('app.baseUrl') ?? '/',
        );
        (window as unknown as { CESIUM_BASE_URL?: string }).CESIUM_BASE_URL =
          base;
        ensureWidgetStyles(base);
        const Cesium = await import('cesium');
        if (cancelled || !host.current) return;
        // No Ion: no default token, no Ion imagery, geocoder or terrain. Imagery is a plain XYZ tile source.
        const viewer = new Cesium.Viewer(host.current, {
          baseLayer: new Cesium.ImageryLayer(
            new Cesium.UrlTemplateImageryProvider({
              url: IMAGERY.url,
              maximumLevel: IMAGERY.maximumLevel,
              credit: IMAGERY.credit,
            }),
            IMAGERY.layer,
          ),
          baseLayerPicker: false,
          geocoder: false,
          terrainProvider: new Cesium.EllipsoidTerrainProvider(),
          animation: false,
          timeline: false,
          homeButton: false,
          sceneModePicker: false,
          navigationHelpButton: false,
          fullscreenButton: false,
          infoBox: false,
          selectionIndicator: false,
          requestRenderMode: true,
          maximumRenderTimeChange: Number.POSITIVE_INFINITY,
        });
        viewer.scene.globe.baseColor = Cesium.Color.fromCssColorString(
          tokens.dark.surface[1],
        );
        viewer.scene.backgroundColor = Cesium.Color.fromCssColorString(
          tokens.dark.surface[0],
        );
        viewer.camera.setView({
          destination: Cesium.Cartesian3.fromDegrees(-60, 30, 14_000_000),
        });

        const trackLayer = createTrackLayer(
          Cesium,
          viewer.scene,
          drawArrow(document),
        );
        const handler = new Cesium.ScreenSpaceEventHandler(viewer.scene.canvas);
        const pickUid = (
          position: import('cesium').Cartesian2,
        ): string | undefined => {
          const picked = viewer.scene.pick(position);
          const id = picked?.primitive?.id;
          return typeof id === 'string' ? id : undefined;
        };
        let frame = 0;
        handler.setInputAction(
          (m: { endPosition: import('cesium').Cartesian2 }) => {
            cancelAnimationFrame(frame);
            frame = requestAnimationFrame(() => {
              const uid = pickUid(m.endPosition);
              trackLayer.hover(uid);
              viewer.scene.canvas.style.cursor = uid ? 'pointer' : 'default';
            });
          },
          Cesium.ScreenSpaceEventType.MOUSE_MOVE,
        );
        handler.setInputAction(
          (m: { position: import('cesium').Cartesian2 }) => {
            latest.current.onSelect(pickUid(m.position));
          },
          Cesium.ScreenSpaceEventType.LEFT_CLICK,
        );

        const observer = new ResizeObserver(() => {
          viewer.resize();
          viewer.scene.requestRender();
        });
        observer.observe(host.current);

        cesiumRef.current = Cesium;
        viewerRef.current = viewer;
        layer.current = trackLayer;
        trackLayer.update(latest.current.tracks, latest.current.fadeNowMs);
        setReady(true);
        cleanup = () => {
          cancelAnimationFrame(frame);
          observer.disconnect();
          handler.destroy();
          trackLayer.destroy();
          viewer.destroy();
        };
      } catch (e) {
        if (!cancelled) onError?.(e instanceof Error ? e.message : String(e));
      }
    })();
    return () => {
      cancelled = true;
      setReady(false);
      layer.current = undefined;
      viewerRef.current = undefined;
      cleanup();
    };
    // The viewer is created once; props reach it through effects below.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [config]);

  useEffect(() => {
    if (!ready) return;
    layer.current?.update(tracks, fadeNowMs);
  }, [ready, tracks, fadeNowMs]);

  useEffect(() => {
    if (ready) layer.current?.select(selectedUid);
  }, [ready, selectedUid]);

  useEffect(() => {
    const Cesium = cesiumRef.current;
    const viewer = viewerRef.current;
    if (!ready || !flyTo || !Cesium || !viewer) return;
    viewer.camera.flyTo({
      destination: Cesium.Cartesian3.fromDegrees(flyTo.lon, flyTo.lat, 400_000),
      duration: 1.5,
    });
    // Only a new nonce moves the camera.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ready, flyTo?.nonce]);

  return (
    <div
      ref={host}
      data-testid="cesium-globe"
      style={{ position: 'absolute', inset: 0 }}
    />
  );
}

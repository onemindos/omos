import type { Track } from './contract';
import {
  AFFILIATION_COLOR,
  DOMAIN_SIZE,
  affiliationOf,
  domainOf,
  fadeAlpha,
  hasCourse,
} from './cot';

/**
 * Draws tracks on a Cesium scene with two primitive collections (no entities): a PointPrimitiveCollection for
 * tracks without a course and a BillboardCollection of rotated arrows for tracks with one. Primitives are
 * created once per uid and updated in place, so a one-second delta costs a property write per changed track.
 *
 * The Cesium module is passed in (it is loaded lazily), which also lets the tests give it a fake.
 */
export type CesiumModule = typeof import('cesium');

export const ARROW_IMAGE_ID = 'onemind-ops-arrow';
const ARROW_PX = 24;
/** Lift above the ellipsoid so markers do not z-fight with the globe surface. */
const MIN_HEIGHT_M = 1500;

type PointPrimitive = import('cesium').PointPrimitive;
type Billboard = import('cesium').Billboard;

interface Entry {
  kind: 'point' | 'arrow';
  prim: PointPrimitive | Billboard;
  signature: string;
  track: Track;
}

export interface TrackLayer {
  /** Reconciles the drawn tracks with `tracks`. `nowMs` drives the fade. */
  update(tracks: readonly Track[], nowMs: number): void;
  /** Ring around one track, or none. */
  select(uid: string | undefined): void;
  /** Label above one track (hover), or none. Selection keeps its label when nothing is hovered. */
  hover(uid: string | undefined): void;
  size(): number;
  destroy(): void;
}

/** White arrow pointing up on a transparent canvas; tinted per track through the billboard colour. */
export function drawArrow(doc: Document): HTMLCanvasElement {
  const canvas = doc.createElement('canvas');
  canvas.width = ARROW_PX;
  canvas.height = ARROW_PX;
  const ctx = canvas.getContext('2d');
  if (ctx) {
    ctx.fillStyle = 'white';
    ctx.strokeStyle = 'black';
    ctx.lineWidth = 1.5;
    ctx.beginPath();
    ctx.moveTo(ARROW_PX / 2, 1);
    ctx.lineTo(ARROW_PX - 3, ARROW_PX - 2);
    ctx.lineTo(ARROW_PX / 2, ARROW_PX - 7);
    ctx.lineTo(3, ARROW_PX - 2);
    ctx.closePath();
    ctx.fill();
    ctx.stroke();
  }
  return canvas;
}

const quantize = (alpha: number) => Math.round(alpha * 20) / 20;

export function createTrackLayer(
  Cesium: CesiumModule,
  scene: import('cesium').Scene,
  arrow: HTMLCanvasElement,
): TrackLayer {
  const points = scene.primitives.add(new Cesium.PointPrimitiveCollection());
  const billboards = scene.primitives.add(
    new Cesium.BillboardCollection({ scene }),
  );
  const labels = scene.primitives.add(new Cesium.LabelCollection());
  const ringPoints = scene.primitives.add(
    new Cesium.PointPrimitiveCollection(),
  );

  const baseColors = {
    friendly: Cesium.Color.fromCssColorString(AFFILIATION_COLOR.friendly),
    hostile: Cesium.Color.fromCssColorString(AFFILIATION_COLOR.hostile),
    neutral: Cesium.Color.fromCssColorString(AFFILIATION_COLOR.neutral),
    unknown: Cesium.Color.fromCssColorString(AFFILIATION_COLOR.unknown),
    other: Cesium.Color.fromCssColorString(AFFILIATION_COLOR.other),
  };

  const entries = new Map<string, Entry>();
  const label = labels.add({
    position: Cesium.Cartesian3.ZERO,
    show: false,
    font: '12px sans-serif',
    fillColor: Cesium.Color.WHITE,
    showBackground: true,
    backgroundColor: Cesium.Color.BLACK.withAlpha(0.7),
    pixelOffset: new Cesium.Cartesian2(0, -22),
    disableDepthTestDistance: Number.POSITIVE_INFINITY,
  });
  const ring = ringPoints.add({
    position: Cesium.Cartesian3.ZERO,
    show: false,
    pixelSize: 26,
    color: Cesium.Color.TRANSPARENT,
    outlineColor: Cesium.Color.WHITE,
    outlineWidth: 2,
    disableDepthTestDistance: Number.POSITIVE_INFINITY,
  });

  let selected: string | undefined;
  let hovered: string | undefined;

  const positionOf = (t: Track) =>
    Cesium.Cartesian3.fromDegrees(
      t.lon,
      t.lat,
      Math.max(t.hae ?? 0, 0) + MIN_HEIGHT_M,
    );

  const refreshAnnotations = () => {
    const sel = selected ? entries.get(selected) : undefined;
    if (sel) {
      ring.position = positionOf(sel.track);
      ring.show = true;
    } else {
      ring.show = false;
    }
    const target = (hovered ? entries.get(hovered) : undefined) ?? sel;
    if (target) {
      label.position = positionOf(target.track);
      label.text = target.track.callsign || target.track.uid;
      label.show = true;
    } else {
      label.show = false;
    }
    scene.requestRender();
  };

  const remove = (entry: Entry) => {
    if (entry.kind === 'point') points.remove(entry.prim as PointPrimitive);
    else billboards.remove(entry.prim as Billboard);
  };

  const create = (t: Track, kind: 'point' | 'arrow'): Entry['prim'] => {
    if (kind === 'point') {
      return points.add({
        id: t.uid,
        position: positionOf(t),
        disableDepthTestDistance: 2_000_000,
      });
    }
    const b = billboards.add({
      id: t.uid,
      position: positionOf(t),
      alignedAxis: Cesium.Cartesian3.UNIT_Z,
      disableDepthTestDistance: 2_000_000,
    });
    b.setImage(ARROW_IMAGE_ID, arrow);
    return b;
  };

  return {
    update(tracks, nowMs) {
      const seen = new Set<string>();
      for (const t of tracks) {
        seen.add(t.uid);
        const kind = hasCourse(t) ? 'arrow' : 'point';
        let entry = entries.get(t.uid);
        if (entry && entry.kind !== kind) {
          remove(entry);
          entries.delete(t.uid);
          entry = undefined;
        }
        const alpha = quantize(fadeAlpha(t, nowMs));
        const signature = `${t.lat},${t.lon},${t.hae ?? ''},${t.course ?? ''},${
          t.type
        },${alpha}`;
        if (entry && entry.signature === signature) {
          entry.track = t;
          continue;
        }
        const prim = entry?.prim ?? create(t, kind);
        const affiliation = affiliationOf(t.type);
        const size = DOMAIN_SIZE[domainOf(t.type)];
        const color = baseColors[affiliation].withAlpha(alpha);
        prim.position = positionOf(t);
        prim.color = color;
        if (kind === 'point') {
          const p = prim as PointPrimitive;
          p.pixelSize = size;
          p.outlineColor = Cesium.Color.BLACK.withAlpha(0.6 * alpha);
          p.outlineWidth = 1;
        } else {
          const b = prim as Billboard;
          b.scale = (size * 1.6) / ARROW_PX;
          b.rotation = -Cesium.Math.toRadians(t.course as number);
        }
        entries.set(t.uid, { kind, prim, signature, track: t });
      }
      for (const [uid, entry] of entries) {
        if (!seen.has(uid)) {
          remove(entry);
          entries.delete(uid);
        }
      }
      refreshAnnotations();
    },
    select(uid) {
      selected = uid;
      refreshAnnotations();
    },
    hover(uid) {
      hovered = uid;
      refreshAnnotations();
    },
    size: () => entries.size,
    destroy() {
      for (const c of [points, billboards, labels, ringPoints])
        scene.primitives.remove(c);
      entries.clear();
    },
  };
}

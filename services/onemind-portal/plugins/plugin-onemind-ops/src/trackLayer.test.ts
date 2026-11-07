import type { Track } from './contract';
import { AFFILIATION_COLOR } from './cot';
import { createTrackLayer } from './trackLayer';
import type { CesiumModule } from './trackLayer';

/** Just enough of Cesium to count primitives and see what was written to them. */
function fakeCesium() {
  const color = (css: string, a = 1) => ({
    css,
    a,
    withAlpha: (x: number) => color(css, x),
  });
  class Collection {
    items: Array<Record<string, unknown>> = [];
    add(o: Record<string, unknown>) {
      const item = { ...o, setImage: jest.fn() };
      this.items.push(item);
      return item;
    }
    remove(p: unknown) {
      this.items = this.items.filter(i => i !== p);
      return true;
    }
  }
  const made: Collection[] = [];
  const mk = () => {
    const c = new Collection();
    made.push(c);
    return c;
  };
  const cesium = {
    PointPrimitiveCollection: function makePoints() {
      return mk();
    },
    BillboardCollection: function makeBillboards() {
      return mk();
    },
    LabelCollection: function makeLabels() {
      return mk();
    },
    Cartesian2: function makeVector(x: number, y: number) {
      return { x, y };
    },
    Cartesian3: {
      ZERO: { z: 0 },
      UNIT_Z: { u: 'z' },
      fromDegrees: (lon: number, lat: number, h: number) => ({ lon, lat, h }),
    },
    Color: {
      WHITE: color('white'),
      BLACK: color('black'),
      TRANSPARENT: color('transparent'),
      fromCssColorString: (css: string) => color(css),
    },
    Math: { toRadians: (d: number) => (d * Math.PI) / 180 },
  };
  const scene = {
    primitives: { add: (c: unknown) => c, remove: jest.fn() },
    requestRender: jest.fn(),
  };
  // PointPrimitiveCollection is created first and again for the ring: order is points, billboards, labels, ring.
  return {
    cesium: cesium as unknown as CesiumModule,
    scene: scene as never,
    made,
  };
}

const track = (uid: string, extra: Partial<Track> = {}): Track => ({
  uid,
  type: 'a-h-A-M-F',
  lat: 10,
  lon: 20,
  time: '2026-11-07T10:00:00Z',
  stale: '2026-11-07T10:10:00Z',
  subject: 'tak.cot.x',
  ...extra,
});
const NOW = Date.parse('2026-11-07T10:01:00Z');

describe('createTrackLayer', () => {
  it('uses points without a course and rotated arrows with one, colored by affiliation', () => {
    const { cesium, scene, made } = fakeCesium();
    const [points, billboards] = [0, 1].map(i => i);
    const layer = createTrackLayer(cesium, scene, {} as HTMLCanvasElement);
    layer.update(
      [track('p1'), track('b1', { course: 90, type: 'a-f-G' })],
      NOW,
    );
    expect(layer.size()).toBe(2);
    expect(made[points].items).toHaveLength(1);
    expect(made[billboards].items).toHaveLength(1);
    const point = made[points].items[0] as {
      color: { css: string };
      pixelSize: number;
    };
    expect(point.color.css).toBe(AFFILIATION_COLOR.hostile);
    expect(point.pixelSize).toBe(12);
    const arrow = made[billboards].items[0] as {
      color: { css: string };
      rotation: number;
      setImage: jest.Mock;
    };
    expect(arrow.color.css).toBe(AFFILIATION_COLOR.friendly);
    expect(arrow.rotation).toBeCloseTo(-Math.PI / 2);
    expect(arrow.setImage).toHaveBeenCalledWith(
      'onemind-ops-arrow',
      expect.anything(),
    );
  });

  it('updates in place, moves a track between collections when its course appears, and removes the gone', () => {
    const { cesium, scene, made } = fakeCesium();
    const layer = createTrackLayer(cesium, scene, {} as HTMLCanvasElement);
    layer.update([track('a'), track('b')], NOW);
    const before = made[0].items[0];
    layer.update([track('a', { lat: 11 }), track('b')], NOW);
    expect(made[0].items[0]).toBe(before);
    expect((before as { position: { lat: number } }).position.lat).toBe(11);

    layer.update([track('a', { course: 45 }), track('b')], NOW);
    expect(made[0].items).toHaveLength(1);
    expect(made[1].items).toHaveLength(1);

    layer.update([track('b')], NOW);
    expect(layer.size()).toBe(1);
    expect(made[1].items).toHaveLength(0);
  });

  it('fades old tracks and labels only the hovered or selected one', () => {
    const { cesium, scene, made } = fakeCesium();
    const layer = createTrackLayer(cesium, scene, {} as HTMLCanvasElement);
    const late = Date.parse('2026-11-07T10:09:00Z');
    layer.update([track('a', { callsign: 'ALPHA' })], late);
    expect((made[0].items[0] as { color: { a: number } }).color.a).toBeLessThan(
      0.5,
    );
    const label = made[2].items[0] as { show: boolean; text?: string };
    expect(label.show).toBe(false);
    layer.hover('a');
    expect(label.show).toBe(true);
    expect(label.text).toBe('ALPHA');
    layer.hover(undefined);
    expect(label.show).toBe(false);
    layer.select('a');
    expect(label.show).toBe(true);
    expect((made[3].items[0] as { show: boolean }).show).toBe(true);
  });
});

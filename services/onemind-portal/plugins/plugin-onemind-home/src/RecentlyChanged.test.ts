import { recentChanges } from './RecentlyChanged';
import {
  APP_ROLE,
  isAppComponent,
  ROLE_LABEL_KEY,
} from '@internal/plugin-onemind-common';
import { countAppsByPillar } from './PillarTiles';

const n = (
  id: string,
  title: string,
  created: string,
  topic = 'deployments',
) => ({ id, title, created, topic });

describe('recentChanges', () => {
  it('dedupes by title keeping the newest, drops other topics and caps at 10', () => {
    const items = [
      n('1', 'Deployment completed: nats', '2026-10-21T10:00:00Z'),
      n('2', 'Deployment completed: nats', '2026-10-21T11:00:00Z'),
      n('3', 'Backup completed: pg', '2026-10-21T11:30:00Z', 'backups'),
      ...Array.from({ length: 14 }, (_, i) =>
        n(
          `x${i}`,
          `Deployment completed: app-${i}`,
          `2026-10-20T0${i % 10}:00:00Z`,
        ),
      ),
    ];
    const rows = recentChanges(items);
    expect(rows).toHaveLength(10);
    expect(rows[0].id).toBe('2');
    expect(rows.filter(r => r.title.endsWith('nats'))).toHaveLength(1);
    expect(rows.some(r => r.topic === 'backups')).toBe(false);
  });
});

describe('countAppsByPillar', () => {
  const app = (pillar?: string) => ({
    labels: {
      'onemind.io/catalog-role': 'app',
      ...(pillar ? { 'onemind.io/pillar': pillar } : {}),
    },
  });
  const workload = (pillar: string) => ({
    labels: {
      'onemind.io/catalog-role': 'workload',
      'onemind.io/pillar': pillar,
    },
  });

  it('counts app components by their pillar label and ignores others', () => {
    const counts = countAppsByPillar([
      app('so'),
      app('so'),
      app('le'),
      app('shared'),
      app(),
      { labels: {} },
      {},
    ]);
    expect(counts).toEqual({ so: 2, le: 1, hp: 0, ge: 0 });
  });

  it('counts an app and a workload under SO as 1, not 2', () => {
    expect(countAppsByPillar([app('so'), workload('so')])).toEqual({
      so: 1,
      le: 0,
      hp: 0,
      ge: 0,
    });
  });

  it('does not count a component with no role label', () => {
    expect(
      countAppsByPillar([{ labels: { 'onemind.io/pillar': 'hp' } }]).hp,
    ).toBe(0);
  });
});

describe('catalog Apps view and tiles share one predicate', () => {
  it('uses the shared role key and value', () => {
    expect(isAppComponent({ [ROLE_LABEL_KEY]: APP_ROLE })).toBe(true);
    expect(isAppComponent({ [ROLE_LABEL_KEY]: 'workload' })).toBe(false);
  });
});

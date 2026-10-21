import { recentChanges } from './RecentlyChanged';
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
  it('counts components by their pillar label and ignores others', () => {
    const counts = countAppsByPillar([
      { labels: { 'onemind.io/pillar': 'so' } },
      { labels: { 'onemind.io/pillar': 'so' } },
      { labels: { 'onemind.io/pillar': 'le' } },
      { labels: { 'onemind.io/pillar': 'shared' } },
      { labels: {} },
      {},
    ]);
    expect(counts).toEqual({ so: 2, le: 1, hp: 0, ge: 0 });
  });
});

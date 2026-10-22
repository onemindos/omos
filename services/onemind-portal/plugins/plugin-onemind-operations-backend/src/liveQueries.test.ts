import { liveQueries, parseNodeUsage, parsePvcs, parseVector, parseVolsync } from './liveQueries';

const vector = (rows: [Record<string, string>, string][]) => ({
  status: 'success',
  data: {
    resultType: 'vector',
    result: rows.map(([metric, value]) => ({ metric, value: [1792666907.5, value] })),
  },
});

describe('parseVector', () => {
  it('reads samples and drops non-numeric values', () => {
    expect(parseVector(vector([[{ a: '1' }, '2.5'], [{ a: '2' }, 'NaN']]))).toEqual([
      { metric: { a: '1' }, value: 2.5 },
    ]);
  });
  it('rejects error and non-vector answers', () => {
    expect(() => parseVector({ status: 'error' })).toThrow('Invalid metrics response');
    expect(() => parseVector({ status: 'success', data: { resultType: 'matrix', result: [] } })).toThrow();
  });
});

describe('parseNodeUsage', () => {
  it('joins cpu, memory and disk by node name and rounds to one decimal', () => {
    const cpu = parseVector(vector([[{ nodename: 'n1' }, '13.497']]));
    const memory = parseVector(vector([[{ nodename: 'n1' }, '59.464'], [{ nodename: 'n2' }, '40']]));
    const disk = parseVector(vector([[{ nodename: 'n1' }, '15.96']]));
    expect(parseNodeUsage({ cpu, memory, disk })).toEqual({
      n1: { cpuPercent: 13.5, memoryPercent: 59.5, diskPercent: 16 },
      n2: { cpuPercent: null, memoryPercent: 40, diskPercent: null },
    });
  });
});

describe('parsePvcs', () => {
  const used = parseVector(
    vector([
      [{ namespace: 'onetak', persistentvolumeclaim: 'valhalla-data' }, '99.97'],
      [{ namespace: 'monitoring', persistentvolumeclaim: 'storage-loki-0' }, '64.5'],
      [{ namespace: 'hermes', persistentvolumeclaim: 'hermes-data' }, '82'],
      [{ namespace: 'x', persistentvolumeclaim: 'grows' }, '50'],
    ]),
  );
  const filling = parseVector(vector([[{ namespace: 'x', persistentvolumeclaim: 'grows' }, '-1']]));
  it('keeps volumes over 80 percent or filling soon, worst first, with severity', () => {
    expect(parsePvcs(used, filling)).toEqual([
      { namespace: 'onetak', name: 'valhalla-data', usedPercent: 100, fillsWithin4Days: false, severity: 'critical' },
      { namespace: 'hermes', name: 'hermes-data', usedPercent: 82, fillsWithin4Days: false, severity: 'warn' },
      { namespace: 'x', name: 'grows', usedPercent: 50, fillsWithin4Days: true, severity: 'warn' },
    ]);
  });
});

describe('parseVolsync', () => {
  it('counts sources and lists the late ones', () => {
    const samples = parseVector(
      vector([
        [{ obj_namespace: 'a', obj_name: 'v1' }, '0'],
        [{ obj_namespace: 'b', obj_name: 'v2' }, '1'],
      ]),
    );
    expect(parseVolsync(samples)).toEqual({ sources: 2, outOfSync: [{ namespace: 'b', name: 'v2' }] });
  });
});

describe('queries', () => {
  it('are fixed read-only expressions with no interpolation points', () => {
    for (const query of Object.values(liveQueries)) {
      expect(query).not.toMatch(/\$\{|%s|\{\{/);
    }
  });
});

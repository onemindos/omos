import { readGatusStatuses, summariseEndpoint, GatusEndpoint, STALE_AFTER_MS } from './gatus';
import { buildFacts, NOT_RECORDED } from './facts';
import { groupByHosting, nonKubernetesRows } from '../nonKubernetes/groups';

const NOW = Date.parse('2026-11-05T12:00:00Z');
const at = (msAgo: number) => new Date(NOW - msAgo).toISOString();
const endpoint = (name: string, result: Partial<GatusEndpoint['results'] extends (infer R)[] | undefined ? R : never>): GatusEndpoint => ({
  name,
  results: [{ success: true, timestamp: at(30_000), duration: 123_000_000, status: 200, conditionResults: [], ...result }],
});

describe('summariseEndpoint', () => {
  it('reports up with response time, last check and the certificate condition', () => {
    const s = summariseEndpoint(
      [endpoint('dropbox', { conditionResults: [{ condition: '[STATUS] < 400', success: true }, { condition: '[CERTIFICATE_EXPIRATION] > 336h', success: true }] })],
      'dropbox',
      NOW,
    );
    expect(s).toMatchObject({ kind: 'checked', status: 'ok', stale: false, responseMs: 123, httpStatus: 200, certificate: 'valid', failing: [] });
  });

  it('reports down and names the failing condition and error', () => {
    const s = summariseEndpoint(
      [endpoint('test-deliberate-failure', { success: false, status: undefined, errors: ['no such host'], conditionResults: [{ condition: '[CONNECTED] (false) == true', success: false }] })],
      'test-deliberate-failure',
      NOW,
    );
    expect(s).toMatchObject({ kind: 'checked', status: 'down', failing: ['[CONNECTED] (false) == true'], errors: ['no such host'], certificate: 'not-checked' });
  });

  it('reports an expiring certificate when the 14 day condition fails', () => {
    const s = summariseEndpoint([endpoint('a', { success: false, conditionResults: [{ condition: '[CERTIFICATE_EXPIRATION] > 336h', success: false }] })], 'a', NOW);
    expect(s).toMatchObject({ kind: 'checked', certificate: 'expiring' });
  });

  it('never shows a stale passing check as up', () => {
    const s = summariseEndpoint([endpoint('a', { timestamp: at(STALE_AFTER_MS + 1000) })], 'a', NOW);
    expect(s).toMatchObject({ kind: 'checked', status: 'unknown', stale: true });
  });

  it('says not checked when Gatus has no result for the name', () => {
    expect(summariseEndpoint([], 'missing', NOW)).toEqual({ kind: 'not-checked' });
    expect(summariseEndpoint([{ name: 'a', results: [] }], 'a', NOW)).toEqual({ kind: 'not-checked' });
  });

  it('uses the latest result', () => {
    const e: GatusEndpoint = {
      name: 'a',
      results: [
        { success: false, timestamp: at(90_000) },
        { success: true, timestamp: at(30_000) },
      ],
    };
    expect(summariseEndpoint([e], 'a', NOW)).toMatchObject({ status: 'ok' });
  });
});

describe('readGatusStatuses', () => {
  const discovery = { getBaseUrl: async () => 'http://x/api/proxy' };
  it('reads the proxy route with a GET and returns the list', async () => {
    const fetch = jest.fn(async () => new Response(JSON.stringify([{ name: 'a' }]), { status: 200 }));
    await expect(readGatusStatuses(discovery, { fetch })).resolves.toEqual([{ name: 'a' }]);
    expect(fetch).toHaveBeenCalledWith('http://x/api/proxy/gatus/statuses?pageSize=1');
  });
  it('fails visibly when Gatus is down or answers with something else', async () => {
    await expect(readGatusStatuses(discovery, { fetch: async () => new Response('', { status: 502 }) })).rejects.toThrow(/HTTP 502/);
    await expect(readGatusStatuses(discovery, { fetch: async () => new Response('{}', { status: 200 }) })).rejects.toThrow(/unexpected/);
  });
});

describe('buildFacts', () => {
  const entity = (annotations: Record<string, string>) => ({ metadata: { annotations }, spec: { owner: 'group:default/onemind-admins' } });
  it('shows "Not recorded" for every fact nobody wrote down, and never invents one', () => {
    const facts = buildFacts(entity({ 'onemind.io/hosting': 'external-provider-unverified' }));
    const byLabel = Object.fromEntries(facts.map(f => [f.label, f.value]));
    expect(byLabel.Owner).toBe('onemind-admins');
    expect(byLabel.Hosting).toBe('SaaS');
    for (const label of ['Plan', 'Renewal', 'Data export', 'Backup location', 'Recovery', 'Who has access']) expect(byLabel[label]).toBeUndefined();
    expect(NOT_RECORDED).toBe('Not recorded');
  });
  it('shows recorded facts and links a runbook path to the runbook site', () => {
    const facts = buildFacts(
      entity({
        'onemind.io/hosting': 'saas',
        'onemind.io/plan': 'Standard, billed yearly',
        'onemind.io/renewal': '2027-01-31',
        'onemind.io/data-export': 'Settings, Export all data',
        'onemind.io/backup-location': 'Dropbox export in the vault',
        'onemind.io/recovery-runbook': 'apps/demo/',
        'onemind.io/access': 'Zeus (owner)',
      }),
    );
    const get = (label: string) => facts.find(f => f.label === label)!;
    expect(get('Plan').value).toBe('Standard, billed yearly');
    expect(get('Renewal').value).toBe('2027-01-31');
    expect(get('Recovery').link).toBe('/docs/default/component/onemind-runbooks/apps/demo/');
    expect(get('Who has access').value).toBe('Zeus (owner)');
  });
  it('says hosting is unknown instead of hiding it', () => {
    expect(buildFacts(entity({})).find(f => f.label === 'Hosting')!.value).toMatch(/^Unknown/);
  });
});

describe('Not on Kubernetes grouping', () => {
  const component = (name: string, annotations: Record<string, string>) => ({ kind: 'Component', metadata: { name, annotations }, spec: { owner: 'group:default/onemind-admins' } });
  const rows = nonKubernetesRows([
    component('grafana', { 'onemind.io/hosting': 'digitalocean-nyc3', 'argocd/app-name': 'platform-monitoring' }),
    component('dropbox', { 'onemind.io/hosting': 'external-provider-unverified', 'onemind.io/health-check': 'dropbox' }),
    component('ha', { 'onemind.io/hosting': 'mesh-endpoint-unverified' }),
    component('mystery', {}),
    { kind: 'System', metadata: { name: 'ignored' } },
  ]);
  it('leaves Kubernetes apps out and groups the rest by hosting type, with unknown visible', () => {
    expect(rows.map(r => r.name).sort()).toEqual(['dropbox', 'ha', 'mystery']);
    const groups = groupByHosting(rows);
    expect(groups.map(g => g.kind)).toEqual(['saas', 'home', 'unknown']);
    expect(groups[0].rows[0].check).toBe('dropbox');
  });
});

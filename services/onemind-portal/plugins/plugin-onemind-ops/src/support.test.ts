import { cesiumBaseUrl, ensureWidgetStyles } from './cesiumSetup';
import { buildNamespaceLinks, headlampUrlFrom } from './namespaceLinks';

describe('cesiumBaseUrl', () => {
  it('serves the assets under the app base path', () => {
    expect(cesiumBaseUrl('https://portal.example.com')).toBe('/cesium/');
    expect(cesiumBaseUrl('https://portal.example.com/')).toBe('/cesium/');
    expect(cesiumBaseUrl('https://example.com/portal/')).toBe(
      '/portal/cesium/',
    );
    expect(cesiumBaseUrl('not a url')).toBe('/cesium/');
  });

  it('adds the widget stylesheet once', () => {
    ensureWidgetStyles('/cesium/');
    ensureWidgetStyles('/cesium/');
    const links = document.querySelectorAll(
      'link[href="/cesium/Widgets/widgets.css"]',
    );
    expect(links).toHaveLength(1);
  });
});

describe('namespace links', () => {
  const components = [
    {
      metadata: {
        name: 'zeta',
        annotations: { 'backstage.io/kubernetes-namespace': 'shared' },
      },
    },
    {
      metadata: {
        name: 'alpha',
        title: 'Alpha',
        annotations: { 'backstage.io/kubernetes-namespace': 'shared' },
      },
    },
    {
      metadata: {
        name: 'matrix',
        namespace: 'default',
        annotations: { 'backstage.io/kubernetes-namespace': 'matrix' },
      },
    },
    { metadata: { name: 'no-annotation' } },
  ];

  it('links to the Kubernetes tab of the first component by name', () => {
    const links = buildNamespaceLinks(components);
    expect(links('shared')).toEqual({
      href: '/catalog/default/component/alpha/kubernetes',
      title: 'Alpha',
      via: 'catalog',
    });
    expect(links('matrix')?.href).toBe(
      '/catalog/default/component/matrix/kubernetes',
    );
  });

  it('falls back to Headlamp, or to nothing', () => {
    expect(
      buildNamespaceLinks(components, 'https://headlamp.example')('unknown'),
    ).toEqual({
      href: 'https://headlamp.example',
      title: 'Headlamp',
      via: 'headlamp',
    });
    expect(buildNamespaceLinks(components)('unknown')).toBeUndefined();
  });

  it('reads the dashboard url from the kubernetes cluster config when present', () => {
    expect(
      headlampUrlFrom(() => [
        {
          type: 'config',
          clusters: [{ name: 'x', dashboardUrl: 'https://h.example' }],
        },
      ]),
    ).toBe('https://h.example');
    expect(headlampUrlFrom(() => undefined)).toBeUndefined();
    expect(
      headlampUrlFrom(() => {
        throw new Error('no config');
      }),
    ).toBeUndefined();
  });
});

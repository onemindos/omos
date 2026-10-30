"""Protocol boundaries for broad coverage and optional writes."""
import json
import secrets

import httpx
import pytest
from starlette.testclient import TestClient
from onemind_mcp.config import Catalog
from test_server import build, rpc


def catalog(monkeypatch, adapter, writes=False):
    monkeypatch.setenv('APP_URL', 'http://app.internal')
    return Catalog(apps=[dict(id='app', adapter=adapter, description='Test app',
                             url_env='APP_URL', allow_writes=writes)])


def test_blocked_apps_visible_without_credentials():
    cat = Catalog(apps=[dict(id='files', description='Family documents', enabled=False,
                            blocked_reason='Dedicated app password required')])
    with TestClient(build(cat, lambda _: pytest.fail('Blocked app made request'))) as client:
        assert len(rpc(client, 'tools/list')['result']['tools']) == 2
        result = rpc(client, 'tools/call', {'name': 'apps_list', 'arguments': {'query': 'documents'}})
        assert 'Dedicated app password required' in json.dumps(result)
        assert 'blocked' in json.dumps(result)


def test_writes_opt_in_annotations_and_http201(monkeypatch):
    for enabled in (False, True):
        with TestClient(build(catalog(monkeypatch, 'paperclip', enabled), lambda _: httpx.Response(201, json={'id': 'new'}))) as client:
            tools = {t['name']: t for t in rpc(client, 'tools/list')['result']['tools']}
            assert ('app_create_issue' in tools) == enabled
            if enabled:
                assert not tools['app_create_issue']['annotations']['readOnlyHint']
                assert not tools['app_create_issue']['annotations']['idempotentHint']
                result = rpc(client, 'tools/call', {'name': 'app_create_issue', 'arguments': {
                    'company_id': '9b1d58fe-8d58-4b18-b336-8acb8439d05e', 'title': 'Review', 'description': 'Review the plan'}})
                assert not result['result'].get('isError')


def test_stac_invalid_bbox_makes_no_request(monkeypatch):
    with TestClient(build(catalog(monkeypatch, 'stac'), lambda _: pytest.fail('Invalid bbox made request'))) as client:
        for bbox in ([0, 0, 181, 20], [5, 0, 1, 20], [0, -91, 20, 0]):
            result = rpc(client, 'tools/call', {'name': 'app_search', 'arguments': {'bbox': bbox}})
            assert result['result']['isError']


def test_n8n_header_and_omitted_workflow_nodes(monkeypatch):
    token = secrets.token_urlsafe(32)
    monkeypatch.setenv('APP_KEY', token)
    cat = catalog(monkeypatch, 'n8n')
    cat.apps[0].token_env = 'APP_KEY'
    cat.apps[0].auth_header = 'X-N8N-API-KEY'
    def handler(request):
        assert request.headers['X-N8N-API-KEY'] == token
        assert 'Authorization' not in request.headers
        return httpx.Response(200, json={'data': [{'id': '1', 'name': 'daily', 'nodes': [{'parameters': 'private'}]}]})
    with TestClient(build(cat, handler)) as client:
        result = rpc(client, 'tools/call', {'name': 'app_workflows', 'arguments': {}})
        assert 'daily' in json.dumps(result)
        assert 'private' not in json.dumps(result)
        assert token not in json.dumps(result)


def test_media_source_credentials_omitted(monkeypatch):
    data = {'items': [{'name': 'camera', 'ready': True, 'source': {'url': 'sensitive'}}]}
    with TestClient(build(catalog(monkeypatch, 'mediamtx'), lambda _: httpx.Response(200, json=data))) as client:
        result = rpc(client, 'tools/call', {'name': 'app_paths', 'arguments': {}})
        assert 'sensitive' not in json.dumps(result)
        assert 'camera' in json.dumps(result)


@pytest.mark.parametrize('adapter,name,path,result', [
    ('ollama', 'models', '/api/tags', {'models': []}), ('alertmanager', 'alerts', '/api/v2/alerts', []),
    ('opencost', 'allocation', '/allocation', {'code': 200}), ('grafana', 'dashboards', '/api/search', []),
    ('karakeep', 'bookmarks', '/api/v1/bookmarks', {}), ('loki', 'query', '/loki/api/v1/query_range', {'status': 'success'}),
])
def test_fixed_read_endpoints_and_array_results(monkeypatch, adapter, name, path, result):
    def handler(request):
        assert request.method == 'GET'
        assert request.url.path == path
        return httpx.Response(200, json=result)
    with TestClient(build(catalog(monkeypatch, adapter), handler)) as client:
        args = {'query': '{namespace="hermes"}'} if adapter == 'loki' else {}
        result = rpc(client, 'tools/call', {'name': 'app_' + name, 'arguments': args})
        assert not result['result'].get('isError')


def test_basic_auth_scheme_sends_user_and_password(monkeypatch):
    import base64
    monkeypatch.setenv('APP_KEY', 'mcp-read:s3cret-value')
    cat = catalog(monkeypatch, 'mediamtx')
    cat.apps[0].token_env = 'APP_KEY'
    cat.apps[0].auth_scheme = 'basic'
    def handler(request):
        assert request.headers['Authorization'] == 'Basic ' + base64.b64encode(b'mcp-read:s3cret-value').decode()
        return httpx.Response(200, json={'items': []})
    with TestClient(build(cat, handler)) as client:
        assert not rpc(client, 'tools/call', {'name': 'app_paths', 'arguments': {}})['result'].get('isError')


def test_litellm_lists_only_model_ids(monkeypatch):
    data = {'data': [{'id': 'hermes-agent', 'owned_by': 'secret-org', 'litellm_params': {'api_key': 'sk-x'}}]}
    def handler(request):
        assert request.url.path == '/v1/models'
        return httpx.Response(200, json=data)
    with TestClient(build(catalog(monkeypatch, 'litellm'), handler)) as client:
        result = rpc(client, 'tools/call', {'name': 'app_models', 'arguments': {}})
        assert 'hermes-agent' in json.dumps(result)
        assert 'sk-x' not in json.dumps(result)


def test_argocd_summaries_and_name_validation(monkeypatch):
    seen = []
    def handler(request):
        seen.append(request.url.path)
        if request.url.path == '/api/v1/applications':
            assert request.method == 'GET'
            return httpx.Response(200, json={'items': [{'metadata': {'name': 'platform-nats', 'annotations': {'x': 'private'}},
                'spec': {'project': 'default'}, 'status': {'sync': {'status': 'Synced', 'revision': 'abcdef0123456789'},
                'health': {'status': 'Healthy'}}}]})
        return httpx.Response(200, json={'metadata': {'name': 'platform-nats'}, 'spec': {'project': 'default',
            'sources': [{'repoURL': 'https://example/repo', 'path': 'p', 'targetRevision': 'main', 'helm': {'values': 'secret'}}]},
            'status': {'sync': {'status': 'OutOfSync'}, 'health': {'status': 'Healthy'}, 'resources': [
                {'kind': 'ConfigMap', 'name': 'a', 'status': 'OutOfSync'}, {'kind': 'Service', 'name': 'b', 'status': 'Synced'}]}})
    with TestClient(build(catalog(monkeypatch, 'argocd'), handler)) as client:
        listing = json.dumps(rpc(client, 'tools/call', {'name': 'app_applications', 'arguments': {}}))
        assert 'platform-nats' in listing and 'abcdef012345' in listing and 'private' not in listing
        detail = json.dumps(rpc(client, 'tools/call', {'name': 'app_application', 'arguments': {'name': 'platform-nats'}}))
        assert 'OutOfSync' in detail and 'secret' not in detail and '"b"' not in detail
        bad = rpc(client, 'tools/call', {'name': 'app_application', 'arguments': {'name': '../../settings'}})
        assert bad['result']['isError']
        assert seen == ['/api/v1/applications', '/api/v1/applications/platform-nats']


def test_backstage_catalog_reads_are_bounded_and_validated(monkeypatch):
    def handler(request):
        assert request.method == 'GET'
        if request.url.path.endswith('/by-query'):
            assert request.url.params['limit'] == '5' and request.url.params['filter'] == 'kind=component'
            return httpx.Response(200, json={'totalItems': 1, 'items': [{'kind': 'Component', 'metadata': {'name': 'nats'}}]})
        assert request.url.path == '/api/catalog/entities/by-name/component/default/nats'
        return httpx.Response(200, json={'kind': 'Component', 'metadata': {'name': 'nats', 'annotations': {'k': 'private'}},
            'spec': {'owner': 'platform', 'secret': 'private'}, 'relations': [{'type': 'ownedBy', 'targetRef': 'group:default/platform'}]})
    with TestClient(build(catalog(monkeypatch, 'backstage'), handler)) as client:
        listing = rpc(client, 'tools/call', {'name': 'app_entities', 'arguments': {'kind': 'Component', 'limit': 5}})
        assert 'nats' in json.dumps(listing)
        detail = json.dumps(rpc(client, 'tools/call', {'name': 'app_entity', 'arguments': {'kind': 'Component', 'name': 'nats'}}))
        assert 'platform' in detail and 'private' not in detail
        for args in ({'kind': 'Component', 'name': '../x'}, {'kind': 'a/b', 'name': 'x'}):
            assert rpc(client, 'tools/call', {'name': 'app_entity', 'arguments': args})['result']['isError']


def test_ca_file_must_be_under_mounted_directory():
    with pytest.raises(Exception):
        Catalog(apps=[dict(id='app', adapter='argocd', description='x', url_env='U', ca_file='/etc/passwd')])
    Catalog(apps=[dict(id='app', adapter='argocd', description='x', url_env='U', ca_file='/etc/onemind-ca/argocd.pem')])


def test_production_catalog_boots_and_does_not_expose_blocked_apps(monkeypatch):
    from pathlib import Path
    from onemind_mcp.config import load_catalog
    cat = load_catalog(Path(__file__).resolve().parents[3] / 'platform/hermes/mcp/apps.yaml')
    for app in cat.apps:
        if app.enabled:
            monkeypatch.setenv(app.url_env, 'http://app.internal')
            if app.token_env:
                monkeypatch.setenv(app.token_env, 'user:' + secrets.token_urlsafe(24))
    with TestClient(build(cat, lambda _: pytest.fail('Discovery must not probe apps'))) as client:
        tools = rpc(client, 'tools/list')['result']['tools']
        assert len(tools) == 27  # 20 read adapters + 6 tile38 reads + tile38_set_point
        # The one write is the tile38 point SET (allow_writes); destructive tile38 tools are not registered.
        assert [t['name'] for t in tools if not t['annotations']['readOnlyHint']] == ['tile38_set_point']
        assert not any(t['name'] in ('tile38_delete_object', 'tile38_drop_collection') for t in tools)
        blocked = [app.id + '_' for app in cat.apps if not app.enabled]
        assert all(not any(t['name'].startswith(prefix) for prefix in blocked) for t in tools)

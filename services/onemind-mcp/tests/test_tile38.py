"""Tile38 adapter: fixed commands, validated arguments, write and destructive gating."""
import httpx
import pytest
from starlette.testclient import TestClient

from onemind_mcp.config import Catalog
from test_server import build, rpc


def catalog(monkeypatch, writes=False, destroys=False):
    monkeypatch.setenv('APP_URL', 'http://tile38.onetak.svc.cluster.local:9851')
    return Catalog(apps=[dict(id='tile38', adapter='tile38', description='Test', url_env='APP_URL',
                             allow_writes=writes, allow_destructive=destroys)])


def names(client):
    return {t['name'] for t in rpc(client, 'tools/list')['result']['tools']}


def call(client, name, args):
    return rpc(client, 'tools/call', {'name': name, 'arguments': args})['result']


READS = {'tile38_server', 'tile38_keys', 'tile38_get', 'tile38_scan', 'tile38_nearby', 'tile38_within'}


def test_registration_gates(monkeypatch):
    seen = {}
    for writes, destroys in ((False, False), (True, False), (True, True)):
        with TestClient(build(catalog(monkeypatch, writes, destroys), lambda _: httpx.Response(200, json={'ok': True}))) as c:
            seen[(writes, destroys)] = names(c) - {'apps_list', 'apps_describe'}
    assert seen[(False, False)] == READS
    assert seen[(True, False)] == READS | {'tile38_set_point'}
    assert seen[(True, True)] == READS | {'tile38_set_point', 'tile38_delete_object', 'tile38_drop_collection'}


def test_commands_are_built_from_validated_tokens(monkeypatch):
    urls = []

    def handler(request):
        urls.append(request.url.raw_path.decode())
        return httpx.Response(200, json={'ok': True, 'objects': [], 'elapsed': '1us'})
    with TestClient(build(catalog(monkeypatch, True, True), handler)) as c:
        call(c, 'tile38_nearby', {'key': 'fleet', 'lat': 33.5, 'lon': -115.2, 'meters': 1000})
        call(c, 'tile38_within', {'key': 'fleet', 'min_lat': 33, 'min_lon': -116, 'max_lat': 34, 'max_lon': -115, 'limit': 5})
        call(c, 'tile38_set_point', {'key': 'gateway-selftest', 'id': 'p1', 'lat': 1.25, 'lon': 2})
        call(c, 'tile38_delete_object', {'key': 'gateway-selftest', 'id': 'p1'})
        call(c, 'tile38_keys', {})
    assert urls == ['/NEARBY+fleet+LIMIT+20+POINT+33.5+-115.2+1000',
                    '/WITHIN+fleet+LIMIT+5+BOUNDS+33+-116+34+-115',
                    '/SET+gateway-selftest+p1+POINT+1.25+2',
                    '/DEL+gateway-selftest+p1',
                    '/KEYS+*']


@pytest.mark.parametrize('args', [
    {'key': 'fleet drop', 'id': 'a'}, {'key': 'a+b', 'id': 'a'}, {'key': 'a/b', 'id': 'a'},
    {'key': 'fleet', 'id': 'a\nFLUSHDB'}, {'key': '', 'id': 'a'}, {'key': 'x' * 65, 'id': 'a'}])
def test_names_cannot_inject_tokens(monkeypatch, args):
    with TestClient(build(catalog(monkeypatch), lambda _: pytest.fail('invalid name made a request'))) as c:
        assert call(c, 'tile38_get', args)['isError']


def test_bad_numbers_and_box_make_no_request(monkeypatch):
    with TestClient(build(catalog(monkeypatch), lambda _: pytest.fail('invalid input made a request'))) as c:
        assert call(c, 'tile38_nearby', {'key': 'k', 'lat': 91, 'lon': 0, 'meters': 5})['isError']
        assert call(c, 'tile38_nearby', {'key': 'k', 'lat': 0, 'lon': 0, 'meters': 0})['isError']
        assert call(c, 'tile38_within', {'key': 'k', 'min_lat': 5, 'min_lon': 0, 'max_lat': 1, 'max_lon': 1})['isError']


def test_command_errors_are_tool_errors(monkeypatch):
    with TestClient(build(catalog(monkeypatch), lambda _: httpx.Response(200, json={'ok': False, 'err': 'key not found'}))) as c:
        result = call(c, 'tile38_get', {'key': 'nokey', 'id': 'x'})
        assert result['isError'] and 'key not found' in str(result)

"""Tile38 live geospatial entities: reviewed fixed commands over the HTTP command API (design row 23).

No command passthrough: every tool builds one fixed command from arguments validated here (names match a strict
pattern, numbers are bounded), so no argument can add a token or change the verb. Reads need no flag, SET needs the
app's `allow_writes`, DEL and DROP additionally need `allow_destructive` (tools the gateway only offers on its
execute route). FLUSHDB, FENCE, hooks and server configuration commands are not implemented.
"""
from typing import Annotated

from mcp.server.fastmcp.exceptions import ToolError
from pydantic import Field

Name = Annotated[str, Field(pattern=r'^[A-Za-z0-9_.:-]{1,64}$')]
Lat = Annotated[float, Field(ge=-90, le=90)]
Lon = Annotated[float, Field(ge=-180, le=180)]
Limit = Annotated[int, Field(ge=1, le=100)]
Meters = Annotated[float, Field(gt=0, le=1_000_000)]


def _num(value: float) -> str:
    return format(float(value), 'f').rstrip('0').rstrip('.') or '0'


def tile38(add, upstream):
    async def run(*tokens: str) -> dict:
        data = await upstream.get('/' + '+'.join(tokens))
        if not isinstance(data, dict) or not data.get('ok'):
            # Tile38 answers HTTP 200 with ok=false for command errors (key not found, bad arguments).
            raise ToolError('Tile38: ' + str(data.get('err', 'command failed'))[:200] if isinstance(data, dict) else 'Tile38: unexpected response')
        data.pop('elapsed', None)
        return data

    async def server() -> dict:
        """Tile38 server statistics: object, collection and hook counts, memory and version. Read only."""
        return await run('SERVER')

    async def keys() -> dict:
        """List all Tile38 collection names. Read only."""
        return await run('KEYS', '*')

    async def get(key: Name, id: Name) -> dict:
        """Get one object (GeoJSON) from a collection by id. Read only."""
        return await run('GET', key, id)

    async def scan(key: Name, limit: Limit = 20) -> dict:
        """List up to 100 objects of a collection. Read only."""
        return await run('SCAN', key, 'LIMIT', str(limit))

    async def nearby(key: Name, lat: Lat, lon: Lon, meters: Meters, limit: Limit = 20) -> dict:
        """Objects of a collection within a radius in meters of a point (NEARBY ... POINT). Read only."""
        return await run('NEARBY', key, 'LIMIT', str(limit), 'POINT', _num(lat), _num(lon), _num(meters))

    async def within(key: Name, min_lat: Lat, min_lon: Lon, max_lat: Lat, max_lon: Lon, limit: Limit = 20) -> dict:
        """Objects of a collection fully inside a bounding box (WITHIN ... BOUNDS). Read only."""
        if min_lat > max_lat or min_lon > max_lon:
            raise ValueError('Invalid bounding box')
        return await run('WITHIN', key, 'LIMIT', str(limit), 'BOUNDS',
                         _num(min_lat), _num(min_lon), _num(max_lat), _num(max_lon))

    async def set_point(key: Name, id: Name, lat: Lat, lon: Lon) -> dict:
        """Create or move one point object in a collection (SET ... POINT). Writes live geospatial state that other services may read; call only for a user-requested change."""
        return await run('SET', key, id, 'POINT', _num(lat), _num(lon))

    async def delete_object(key: Name, id: Name) -> dict:
        """Delete one object from a collection (DEL). Destructive: not offered on the gateway group routes."""
        return await run('DEL', key, id)

    async def drop_collection(key: Name) -> dict:
        """Drop a whole collection and all its objects (DROP). Destructive: not offered on the gateway group routes."""
        return await run('DROP', key)

    for fn in (server, keys, get, scan, nearby, within):
        add(fn)
    add(set_point, read_only=False)
    add(delete_object, read_only=False, destructive=True)
    add(drop_collection, read_only=False, destructive=True)

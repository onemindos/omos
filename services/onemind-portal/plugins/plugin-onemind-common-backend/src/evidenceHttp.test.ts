/** @jest-environment node */
import http from 'node:http';
import type { AddressInfo } from 'node:net';
import { getEvidenceHttp } from './evidenceHttp';

test('HTTP observer does not follow redirects or send authorization', async () => {
  const paths: string[] = [];
  const headers: (string | undefined)[] = [];
  const server = http.createServer((req, res) => {
    paths.push(req.url!);
    headers.push(req.headers.authorization);
    res.writeHead(302, { Location: '/private-target' });
    res.end();
  });
  await new Promise<void>(resolve => server.listen(0, '127.0.0.1', resolve));
  try {
    const r = await getEvidenceHttp(
      `http://127.0.0.1:${(server.address() as AddressInfo).port}/start`,
    );
    expect(r.status).toBe(302);
    expect(paths).toEqual(['/start']);
    expect(headers).toEqual([undefined]);
  } finally {
    await new Promise<void>(resolve => server.close(() => resolve()));
  }
});
test('oversized HTTP source is rejected', async () => {
  const server = http.createServer((_req, res) => {
    res.end(Buffer.alloc(2_000_001));
  });
  await new Promise<void>(resolve => server.listen(0, '127.0.0.1', resolve));
  try {
    await expect(
      getEvidenceHttp(
        `http://127.0.0.1:${(server.address() as AddressInfo).port}/`,
      ),
    ).rejects.toThrow('response limit');
  } finally {
    await new Promise<void>(resolve => server.close(() => resolve()));
  }
});

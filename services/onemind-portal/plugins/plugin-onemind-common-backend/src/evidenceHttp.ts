import http from 'node:http';

export type HttpObservation = { status: number; body: string };

// Internal, fixed targets only. No client URL, credential forwarding or redirects.
export function getEvidenceHttp(url: string): Promise<HttpObservation> {
  return new Promise((resolve, reject) => {
    const req = http.get(url, response => {
      const chunks: Buffer[] = [];
      let bytes = 0;
      response.on('data', (chunk: Buffer) => {
        bytes += chunk.length;
        if (bytes > 2_000_000) {
          req.destroy(new Error('Observation response limit'));
          return;
        }
        chunks.push(chunk);
      });
      response.on('error', reject);
      response.on('end', () =>
        resolve({
          status: response.statusCode ?? 0,
          body: Buffer.concat(chunks).toString('utf8'),
        }),
      );
    });
    const deadline = setTimeout(
      () => req.destroy(new Error('Observation deadline')),
      5000,
    );
    req.on('close', () => clearTimeout(deadline));
    req.on('error', reject);
  });
}

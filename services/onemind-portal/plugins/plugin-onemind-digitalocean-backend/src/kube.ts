import { readFileSync } from 'node:fs';
import https from 'node:https';

const SERVICE_ACCOUNT = '/var/run/secrets/kubernetes.io/serviceaccount/';

export type GetJson = (path: string) => Promise<unknown>;

/**
 * Read-only GET against the API server with the pod's own ServiceAccount (token and CA re-read per
 * call). Same approach as the live inventory module (plugin-onemind-inventory-backend/src/kube.ts).
 */
export function inClusterGetJson(
  apiServer: string,
  timeoutMs: number,
): GetJson {
  const url = new URL(apiServer);
  return path =>
    new Promise((resolve, reject) => {
      const request = https.get(
        {
          hostname: url.hostname,
          port: url.port || 443,
          path,
          ca: readFileSync(`${SERVICE_ACCOUNT}ca.crt`),
          headers: {
            Authorization: `Bearer ${readFileSync(
              `${SERVICE_ACCOUNT}token`,
              'utf8',
            ).trim()}`,
            Accept: 'application/json',
          },
        },
        response => {
          let body = '';
          response.setEncoding('utf8');
          response.on('data', chunk => {
            body += chunk;
            if (body.length > 50_000_000)
              request.destroy(new Error('Response limit'));
          });
          response.on('error', reject);
          response.on('end', () => {
            if (response.statusCode !== 200) {
              reject(
                new Error(`Kubernetes API answered ${response.statusCode}`),
              );
              return;
            }
            try {
              resolve(JSON.parse(body));
            } catch {
              reject(new Error('Kubernetes API returned invalid JSON'));
            }
          });
        },
      );
      const deadline = setTimeout(
        () => request.destroy(new Error('Observation deadline')),
        timeoutMs,
      );
      request.on('close', () => clearTimeout(deadline));
      request.on('error', reject);
    });
}

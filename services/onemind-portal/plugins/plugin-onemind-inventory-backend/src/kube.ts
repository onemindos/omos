import { readFileSync } from 'node:fs';
import https from 'node:https';
import type { Workload } from './entities';

const SERVICE_ACCOUNT = '/var/run/secrets/kubernetes.io/serviceaccount/';
const KINDS = [
  { kind: 'Deployment', path: '/apis/apps/v1/deployments' },
  { kind: 'StatefulSet', path: '/apis/apps/v1/statefulsets' },
];

export type GetJson = (path: string) => Promise<unknown>;

/** Read-only GET against the API server with the pod's own ServiceAccount (token and CA re-read per call). */
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
                new Error(
                  `Kubernetes API answered ${response.statusCode} for ${
                    path.split('?')[0]
                  }`,
                ),
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

interface ListResponse {
  items?: {
    metadata?: { name?: string; namespace?: string };
    spec?: { selector?: { matchLabels?: Record<string, string> } };
  }[];
  metadata?: { continue?: string };
}

/**
 * Every Deployment and StatefulSet in the cluster (the same two kinds snapshot-cluster.cjs reads).
 * Only identity and the selector are kept; no env, images, annotations or Secrets are read or retained.
 * Objects without a matchLabels selector are skipped, as the Kubernetes plugin could not select them.
 */
export async function listWorkloads(getJson: GetJson): Promise<Workload[]> {
  const out: Workload[] = [];
  for (const { kind, path } of KINDS) {
    let token: string | undefined;
    do {
      const page = (await getJson(
        `${path}?limit=500${
          token ? `&continue=${encodeURIComponent(token)}` : ''
        }`,
      )) as ListResponse;
      if (!page || !Array.isArray(page.items))
        throw new Error(`Unexpected ${kind} list response`);
      for (const item of page.items) {
        const name = item.metadata?.name;
        const namespace = item.metadata?.namespace;
        const labels = item.spec?.selector?.matchLabels;
        if (!name || !namespace || !labels || !Object.keys(labels).length)
          continue;
        out.push({
          kind,
          namespace,
          name,
          selector: Object.entries(labels)
            .map(([k, v]) => `${k}=${v}`)
            .join(','),
        });
      }
      token = page.metadata?.continue || undefined;
    } while (token);
  }
  return out;
}

import https from 'node:https';
import { readFileSync } from 'node:fs';

export type KubeItem = { metadata: { name: string; namespace?: string }; spec?: { replicas?: number; suspend?: boolean }; status?: { readyReplicas?: number; lastSuccessfulTime?: string; health?: { status?: string }; sync?: { status?: string } } };
export function list<T = KubeItem>(path: string): Promise<T[]> {
  return new Promise((resolve, reject) => {
    const root='/var/run/secrets/kubernetes.io/serviceaccount/';
    let ca: Buffer; let token: string;
    try { ca=readFileSync(root+'ca.crt'); token=readFileSync(root+'token','utf8').trim(); }
    catch { reject(new Error('Cluster credentials unavailable')); return; }
    const request=https.get({hostname:'kubernetes.default.svc',port:443,path,ca,headers:{Authorization:`Bearer ${token}`}}, response => {
      let body=''; response.setEncoding('utf8');
      response.on('data', chunk => {body+=chunk; if(body.length>20_000_000)request.destroy(new Error('Response too large'));});
      response.on('end',()=>{
        if(response.statusCode!==200){reject(new Error(`Cluster API returned ${response.statusCode}`));return;}
        try {const value=JSON.parse(body); if (!Array.isArray(value.items) || value.metadata?.continue) throw new Error('Incomplete cluster list'); resolve(value.items);}catch {reject(new Error('Invalid cluster response'));}
      });
    });
    request.setTimeout(8000,()=>request.destroy(new Error('Cluster request timed out')));
    request.on('error',reject);
  });
}

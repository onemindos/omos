import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import express from 'express';
import https from 'node:https';
import { readFileSync } from 'node:fs';

type KubeItem = { metadata: { name: string; namespace?: string }; spec?: { replicas?: number; suspend?: boolean }; status?: { readyReplicas?: number; lastSuccessfulTime?: string; health?: { status?: string }; sync?: { status?: string } } };
function list(path: string): Promise<KubeItem[]> {
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
        try {resolve(JSON.parse(body).items);}catch {reject(new Error('Invalid cluster response'));}
      });
    });
    request.setTimeout(8000,()=>request.destroy(new Error('Cluster request timed out')));
    request.on('error',reject);
  });
}
export default createBackendPlugin({
  pluginId:'onemind',
  register(reg){reg.registerInit({deps:{router:coreServices.httpRouter,httpAuth:coreServices.httpAuth},async init({router,httpAuth}){
    const routes=express.Router();
    routes.get('/status',async(req,res,next)=>{
      try{
        await httpAuth.credentials(req,{allow:['user']});
        const paths=['/apis/apps/v1/deployments','/apis/apps/v1/statefulsets','/apis/argoproj.io/v1alpha1/namespaces/argocd/applications','/apis/batch/v1/cronjobs'];
        const results=await Promise.allSettled(paths.map(list));
        const items=(i:number)=>results[i].status==='fulfilled'?(results[i] as PromiseFulfilledResult<KubeItem[]>).value:[];
        const workloads=[...items(0),...items(1)];
        const unavailable=workloads.filter(x=>(x.spec?.replicas??1)>(x.status?.readyReplicas??0)).map(x=>`${x.metadata.namespace}/${x.metadata.name}`);
        const deployments=items(2).filter(x=>x.status?.health?.status!=='Healthy'||x.status?.sync?.status!=='Synced').map(x=>({name:x.metadata.name,health:x.status?.health?.status??'Unknown',sync:x.status?.sync?.status??'Unknown'}));
        const backups=items(3).filter(x=>x.metadata.name.includes('backup')).map(x=>({name:`${x.metadata.namespace}/${x.metadata.name}`,lastSuccess:x.status?.lastSuccessfulTime??null,suspended:x.spec?.suspend??false}));
        res.json({observedAt:new Date().toISOString(),source:'do-nyc3-omos-cloud',partial:results.some(x=>x.status==='rejected'),unavailableSources:results.flatMap((x,i)=>x.status==='rejected'?[['Deployments','StatefulSets','Argo CD','Backup schedules'][i]]:[]),unavailable,deployments,backups,totalWorkloads:workloads.length});
      }catch(e){next(e);}
    });
    router.use(routes);
  }});},
});

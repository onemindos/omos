import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import express from 'express';
import { calendarSummary } from './calendar';
import { natsSummary } from './natsSummary';
import { health } from './health';
import { cloudSummary } from './cloud';
import { list, KubeItem } from './clusterRead';

export default createBackendPlugin({
  pluginId:'onemind',
  register(reg){reg.registerInit({deps:{router:coreServices.httpRouter,httpAuth:coreServices.httpAuth,userInfo:coreServices.userInfo},async init({router,httpAuth,userInfo}){
    const routes=express.Router();
    routes.get('/status',async(req,res,next)=>{
      try{
        await httpAuth.credentials(req,{allow:['user']});
        const paths=['/apis/apps/v1/deployments','/apis/apps/v1/statefulsets','/apis/argoproj.io/v1alpha1/namespaces/argocd/applications','/apis/batch/v1/cronjobs'];
        const [results,serviceHealth]=await Promise.all([Promise.allSettled(paths.map(list)),health()]);
        const items=(i:number)=>results[i].status==='fulfilled'?(results[i] as PromiseFulfilledResult<KubeItem[]>).value:[];
        const workloads=[...items(0),...items(1)];
        const unavailable=workloads.filter(x=>(x.spec?.replicas??1)>(x.status?.readyReplicas??0)).map(x=>`${x.metadata.namespace}/${x.metadata.name}`);
        const deployments=items(2).filter(x=>x.status?.health?.status!=='Healthy'||x.status?.sync?.status!=='Synced').map(x=>({name:x.metadata.name,health:x.status?.health?.status??'Unknown',sync:x.status?.sync?.status??'Unknown'}));
        const backups=items(3).filter(x=>x.metadata.name.includes('backup')).map(x=>({name:`${x.metadata.namespace}/${x.metadata.name}`,lastSuccess:x.status?.lastSuccessfulTime??null,suspended:x.spec?.suspend??false}));
        res.json({serviceHealth,observedAt:new Date().toISOString(),source:'do-nyc3-omos-cloud',partial:results.some(x=>x.status==='rejected'),unavailableSources:results.flatMap((x,i)=>x.status==='rejected'?[['Deployments','StatefulSets','Argo CD','Backup schedules'][i]]:[]),unavailable,deployments,backups,totalWorkloads:workloads.length});
      }catch(e){next(e);}
    });
    routes.get('/fabric',async(req,res,next)=>{try{await httpAuth.credentials(req,{allow:['user']});res.json(await natsSummary());}catch(e){next(e);}});
    routes.get('/calendar',async(req,res,next)=>{try{const credentials=await httpAuth.credentials(req,{allow:['user']});const info=await userInfo.getUserInfo(credentials);if(!info.ownershipEntityRefs.includes('group:default/onemind-admins')){res.status(403).json({error:'Calendar summaries require administrator membership'});return;}res.json(await calendarSummary());}catch(e){next(e);}});
    routes.get('/cloud', async(req,res,next)=>{try{const credentials=await httpAuth.credentials(req,{allow:['user']});const info=await userInfo.getUserInfo(credentials);if(!info.ownershipEntityRefs.includes('group:default/onemind-admins')){res.status(403).json({error:'Cloud summaries require OneMind administrator membership'});return;}res.json(await cloudSummary());}catch(e){next(e);}});
    router.use(routes);
  }});},
});

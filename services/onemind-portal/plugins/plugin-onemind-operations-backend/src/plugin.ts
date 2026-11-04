/* eslint-disable */
import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import express from 'express';
import { calendarSummary, costSummary, natsSummary, cloudSummary } from '@internal/plugin-onemind-common-backend';
import { createCollector } from './collect';
import { createLiveRoutes } from './routes';
import { realDeps } from './upstreams';

export default createBackendPlugin({
  pluginId:'onemind',
  register(reg){reg.registerInit({deps:{router:coreServices.httpRouter,httpAuth:coreServices.httpAuth,userInfo:coreServices.userInfo},async init({router,httpAuth,userInfo}){
    const routes=express.Router();
    routes.use(createLiveRoutes({httpAuth,collector:createCollector(realDeps)}));
    routes.get('/fabric',async(req,res,next)=>{try{await httpAuth.credentials(req,{allow:['user']});res.json(await natsSummary());}catch(e){next(e);}});
    routes.get('/calendar',async(req,res,next)=>{try{const credentials=await httpAuth.credentials(req,{allow:['user']});const info=await userInfo.getUserInfo(credentials);if(!info.ownershipEntityRefs.includes('group:default/onemind-admins')){res.status(403).json({error:'Calendar summaries require administrator membership'});return;}res.json(await calendarSummary());}catch(e){next(e);}});
    routes.get('/costs',async(req,res,next)=>{try{const credentials=await httpAuth.credentials(req,{allow:['user']});const info=await userInfo.getUserInfo(credentials);if(!info.ownershipEntityRefs.includes('group:default/onemind-admins')){res.status(403).json({error:'Billing requires OneMind administrator membership'});return;}res.json(await costSummary());}catch(e){next(e);}});
    routes.get('/cloud', async(req,res,next)=>{try{const credentials=await httpAuth.credentials(req,{allow:['user']});const info=await userInfo.getUserInfo(credentials);if(!info.ownershipEntityRefs.includes('group:default/onemind-admins')){res.status(403).json({error:'Cloud summaries require OneMind administrator membership'});return;}res.json(await cloudSummary());}catch(e){next(e);}});
    router.use(routes);
  }});},
});

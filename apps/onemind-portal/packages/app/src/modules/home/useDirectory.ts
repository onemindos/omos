import { useEffect, useMemo, useState } from 'react';
import { storageApiRef, useApi } from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';
import curated from './directory.json';

export type AppEntry = { name:string; title:string; description:string; url:string; group:string; groupTitle:string; ref:string; personal?:boolean; favorite?:boolean };
export function safeUrl(value:string) {
  try { const url=new URL(value); return url.protocol==='https:'&&!url.username&&!url.password ? url.href : undefined; } catch { return undefined; }
}
export function useShortcuts() {
  const storage=useApi(storageApiRef);
  const bucket=useMemo(()=>storage.forBucket('onemind-home'),[storage]);
  const [snapshot,setSnapshot]=useState(()=>bucket.snapshot<AppEntry[]>('shortcuts'));
  useEffect(()=>{setSnapshot(bucket.snapshot('shortcuts'));const subscription=bucket.observe$<AppEntry[]>('shortcuts').subscribe(setSnapshot);return()=>subscription.unsubscribe();},[bucket]);
  const shortcuts=(Array.isArray(snapshot.value)?snapshot.value:[]).filter(item=>item&&typeof item.name==='string'&&typeof item.title==='string'&&typeof item.url==='string'&&safeUrl(item.url)).map(item=>({...item,personal:true,ref:`shortcut:${item.name}`,group:'personal-shortcuts',groupTitle:'My shortcuts',description:typeof item.description==='string'?item.description:''}));
  return {shortcuts,loaded:snapshot.presence!=='unknown',save:async(items:AppEntry[])=>{await bucket.set('shortcuts',items);}};
}
const requests=new WeakMap<object,Promise<AppEntry[]>>();
export function useDirectory() {
  const api=useApi(catalogApiRef);
  const personal=useShortcuts();
  const [apps,setApps]=useState<AppEntry[]>([]);
  const [error,setError]=useState(false),[loaded,setLoaded]=useState(false);
  useEffect(()=>{
    let active=true;
    const refresh=async()=>{
      let request=requests.get(api);
      if(!request){
        request=(async()=>{
          const entries:AppEntry[]=[]; let cursor:string|undefined;
          for(let page=0;page<20;page++){
            const response=await api.queryEntities(cursor?{cursor,limit:200}:{filter:{kind:'Component'},limit:200});
            for(const entity of response.items){
              const raw=entity.metadata.annotations?.['onemind.io/launcher-url']??entity.metadata.links?.find(link=>link.title==='Open application')?.url;
              const url=raw&&safeUrl(raw); if(!url)continue;
              const known=curated.find(item=>item.name===entity.metadata.name);
              const groupTitle=entity.metadata.annotations?.['onemind.io/launcher-category']??known?.groupTitle??'Services';
              entries.push({name:entity.metadata.name,ref:`component:${entity.metadata.namespace??'default'}/${entity.metadata.name}`,title:entity.metadata.title??entity.metadata.name,description:entity.metadata.description??'',url,group:known?.group??'services',groupTitle});
            }
            cursor=response.pageInfo.nextCursor;
            if(!cursor)break;
            if(page===19)throw new Error('Catalog pagination limit');
          }
          return entries;
        })();
        requests.set(api,request);
        request.finally(()=>requests.delete(api)).catch(()=>{});
      }
      try{const values=await request;if(active){setApps(values);setError(false);setLoaded(true);}}catch{if(active){setError(true);setLoaded(true);}}
    };
    refresh();const timer=setInterval(refresh,60000);
    return()=>{active=false;clearInterval(timer);};
  },[api]);
  const all=[...apps,...personal.shortcuts];
  const directory=all.filter((app,index)=>all.findIndex(other=>other.url===app.url)===index);
  return {directory,all,loaded:loaded&&personal.loaded,error,personal};
}

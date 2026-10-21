/* eslint-disable */
import http from 'node:http';

// Fixed cluster targets only. Never accept a URL or query from the client.
const targets=[
  {name:'Identity',entity:'component:default/k8s-auth-platform-auth-helm-authentik-server-78f74c27',url:'http://platform-auth-helm-authentik-server.auth.svc/-/health/ready/',link:'https://onemind-auth.tail717ba0.ts.net',scope:'Readiness endpoint'},
  {name:'Grafana',entity:'component:default/grafana',url:'http://monitoring-grafana.monitoring.svc/api/health',link:'https://grafana.tail717ba0.ts.net',scope:'Application health endpoint'},
  {name:'Headlamp',entity:'component:default/headlamp',url:'http://headlamp.kube.svc:4466/',link:'https://headlamp.tail717ba0.ts.net',scope:'Web page response'},
  {name:'NATS dashboard',entity:'component:default/nats-dashboard',url:'http://nats-dashboard.nats.svc:3000/',link:'https://nats-dashboard-1.tail717ba0.ts.net',scope:'Web page response'},
];
function get(url:string):Promise<{code:number;body:string}>{return new Promise((resolve,reject)=>{
  const request=http.get(url,response=>{let body='';response.setEncoding('utf8');response.on('data',chunk=>{body+=chunk;if(body.length>2_000_000)request.destroy(new Error('Response limit'));});response.on('end',()=>resolve({code:response.statusCode??0,body}));response.on('error',reject);});
  const deadline=setTimeout(()=>request.destroy(new Error('Deadline exceeded')),5000);request.on('close',()=>clearTimeout(deadline));request.on('error',reject);
});}
const prometheus='http://monitoring-kube-prometheus-prometheus.monitoring.svc:9090';
async function metric(query:string){const r=await get(`${prometheus}/api/v1/query?query=${encodeURIComponent(query)}`);if(r.code!==200)throw new Error('Metrics unavailable');const json=JSON.parse(r.body);if(json.status!=='success'||json.data?.resultType!=='vector')throw new Error('Invalid metrics');return json.data.result as {metric:Record<string,string>;value:[number,string]}[];}
async function observe(){
 const checks=await Promise.all(targets.map(async target=>{const start=Date.now();try{const r=await get(target.url);return {name:target.name,entity:target.entity,link:target.link,scope:target.scope,state:r.code===200?'up':r.code===401||r.code===403?'authentication-required':'down',httpStatus:r.code,latencyMs:Date.now()-start,observedAt:new Date().toISOString()};}catch{return {name:target.name,entity:target.entity,link:target.link,scope:target.scope,state:'unknown',observedAt:new Date().toISOString()};}}));
 const namespaces=['backstage','auth','nats','takserver','onetak'];
 const selector='namespace=~"'+namespaces.join('|')+'",container!="",container!="POD"';
 const results=await Promise.allSettled([
  metric(`sum by(namespace) (rate(container_cpu_usage_seconds_total{${selector}}[5m]))`),
  metric(`sum by(namespace) (container_memory_working_set_bytes{${selector}})`),
 ]);
 const values=(i:number)=>results[i].status==='fulfilled'?(results[i] as PromiseFulfilledResult<{metric:Record<string,string>;value:[number,string]}[]>).value:[];
 const metrics=namespaces.map(namespace=>{const cpu=values(0).find(x=>x.metric.namespace===namespace),memory=values(1).find(x=>x.metric.namespace===namespace);return {namespace,cpuCores:cpu&&Number.isFinite(Number(cpu.value[1]))?Number(cpu.value[1]):null,memoryBytes:memory&&Number.isFinite(Number(memory.value[1]))?Number(memory.value[1]):null,observedAt:new Date().toISOString(),cpuSampleAt:cpu?new Date(cpu.value[0]*1000).toISOString():null,memorySampleAt:memory?new Date(memory.value[0]*1000).toISOString():null};});
 return {observedAt:new Date().toISOString(),checks,metrics,metricsUnavailable:results.some(x=>x.status==='rejected'),scope:'Cluster-side observations; does not validate Tailscale access, login flows, or NATS message delivery.'};
}
let cache:Awaited<ReturnType<typeof observe>>|undefined;let pending:Promise<Awaited<ReturnType<typeof observe>>>|undefined;
export async function health(){if(cache&&Date.now()-Date.parse(cache.observedAt)<60000)return cache;if(!pending)pending=observe().then(value=>(cache=value)).finally(()=>{pending=undefined;});return pending;}

// DigitalOcean billing totals only. The portal must receive a dedicated
// billing:read token, never the administrator's existing doctl credential.
type Costs={state:'configuration-required'|'current'|'stale'|'unavailable';source:string;scope:string;currency:'USD';observedAt:string|null;generatedAt:string|null;attemptedAt:string;expiresAt:string|null;monthToDateUsage?:string;monthToDateBalance?:string;accountBalance?:string;allocation:'not-configured'};
let cached:Costs|undefined,last:Costs|undefined,pending:Promise<Costs>|undefined,nextAttempt=0;
function money(value:unknown){if(typeof value!=='string'||! /^-?\d{1,12}(\.\d{1,6})?$/.test(value))throw new Error('Invalid balance');return value;}
export async function costSummary():Promise<Costs>{
 if(pending)return pending;if(last&&Date.now()<nextAttempt)return last;
 nextAttempt=Date.now()+900000;
 const empty:Costs={state:'configuration-required',source:'DigitalOcean Billing API',scope:'Entire connected DigitalOcean account; not OneMind-only or namespace allocation',currency:'USD',observedAt:null,generatedAt:null,attemptedAt:new Date().toISOString(),expiresAt:null,allocation:'not-configured'};
 pending=(async()=>{
  const token=process.env.DO_BILLING_READ_TOKEN;if(!token)return empty;
  try{
   const response=await fetch('https://api.digitalocean.com/v2/customers/my/balance',{headers:{Authorization:`Bearer ${token}`},redirect:'error',signal:AbortSignal.timeout(8000)});
   if(!response.ok||!response.body)throw new Error('Billing unavailable');
   const reader=response.body.getReader();const chunks:Uint8Array[]=[];let size=0;
   try{while(true){const part=await reader.read();if(part.done)break;size+=part.value.byteLength;if(size>32768)throw new Error('Billing response limit');chunks.push(part.value);}}finally{await reader.cancel().catch(()=>{});}
   const value=JSON.parse(Buffer.concat(chunks).toString('utf8'));
   if(typeof value.generated_at!=='string'||!Number.isFinite(Date.parse(value.generated_at)))throw new Error('Invalid source timestamp');
   const now=Date.now();
   cached={...empty,state:'current',observedAt:new Date(now).toISOString(),generatedAt:value.generated_at,expiresAt:new Date(now+900000).toISOString(),monthToDateUsage:money(value.month_to_date_usage),monthToDateBalance:money(value.month_to_date_balance),accountBalance:money(value.account_balance)};
   return cached;
  }catch{return cached?{...cached,state:'stale' as const,attemptedAt:empty.attemptedAt}:{...empty,state:'unavailable' as const};}
 })().then(value=>(last=value)).finally(()=>{pending=undefined;});return pending;
}

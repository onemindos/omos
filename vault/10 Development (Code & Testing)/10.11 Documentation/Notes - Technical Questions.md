---
id: ui-12-20260111053926
title: "Notes"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/12]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Notes"
---

# Notes

## docker or kubernetes? What should I start with? Will it be hard to switch over? And when will I know it's time to switch?   

supabase cloud is my pick not self hosted  how do I separate each piece of information inside of superbase so it doesn't mix altogether, for example letta connect, airbyte connection,etc 

## If I connect the same supabase database do I need to connect the vector if it's already the unified database? 

## Monolith repo
Does this mean I will deploy via GitHub? Am I forking or just using docker images? 

## What is the flow to launch new apps? 

Shouldn't I use emqx to have my cloud structure for my IOT with clusters 

we will be using airbyte cloud not self hosted and clickup is a native feature   I want n8n with workers how can I deploy that inside of one VM  should I give letta its own database and just connect it with pgvector what do you think?  I plan to also use livekit cloud is there any benefit to host it myself? Or is cloud good for now?  PowerSync provides true offline sync (not just caching)  for telemetry I rather use something in the cloud grafana vs lang fuse vs prometheus 

is portainer business absolutely needed? For async mode? What is that and go deeper on why its needed  I will mention I have a 8gb of ram beeline I can use and also a server with 32gb of ram with 8tb of space. This we can use for testing. I can use the beeline for Home assistant but the main instance would be the cloud to keep the state and device for the edge   what do you mean by TIER 3: DELIBERATIVE ? TIER 2: PERCEPTUAL? TIER 1: REACTIVE?   is there any better platform for my use case then airbtye ? Like stitch,fivetran,data stage, sap data intelligence, rivery, talent, hero data, airflow, matillion,meltano, informatica, pentaho, mulesoft, integrate.io, estuary, weld, qlik, dlthub, unified.to or I can add to the stack to pump into supabase. Make sure it has the ability to connect.   dify will be self hosted  n8n self hosted with queue mode   tail scale vs cloud flare?   no need for IAC? Correct?  does portainer allow to deploy from git or something like that?   is 4cpu and 16gb 100gb ssd enough to start my super server?  can I use livekit cloud and self hosted and sync together?  I also have supabase free for 6 months with credits they gave me so keep that in mind.   what is Initialize monorepo ? npx create-turbo@latest   I want to do phases 1-6 all in the next 90 mins. List it all with the information I provided with a more straight forward plan and the hardware I mentioned. None of this is complex and shouldn't take week to deploy. We can start with the vm and move from there step by step.   

---

## Review for legacy 
- Zep to nano Nivida
- Zep/ graphiti
- Livekit 
- Mem0 
- YouTube video repo for agent 
- 3d printing repo for Jarvis
- Frigate 
- Wakeword
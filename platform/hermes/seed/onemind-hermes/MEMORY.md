TAK: reached over the tailnet only (tak.tail717ba0.ts.net; ports 8089, 8443, 8446, 64738/TCP via the tak-tailnet Tailscale Service). No public TAK listener; the old public LB and tak.onemindos.dev are retired. Client enrollment is by certificate; credentials are in SealedSecrets, never in memory. TAK core is the Helm release takserver; overlays are Argo apps platform-takserver and platform-fedhub. Voice (Mumble/Murmur, port 64738) is in the takserver namespace; MediaMTX and the TAK portal are in the onetak namespace per docs/STACK.md (verify live before relying on this).
§
CloudTAK and other public onemindos.dev hosts, and the geo stack, live in the onetak area (Tile38, WebODM, farmOS, Node-RED, STAC, Valhalla, Nominatim and others; DimOS is the robot stack; robot subjects are fabric.robot.<id>.*). Hostnames and what is actually deployed change: read docs/STACK.md and check the cluster (read-only) before stating them.
§
Mac edge: the pod can reach Zeus's Mac over SSH at omos-mbp.hermes.svc.cluster.local (key held in the pod's secrets, not in Git); Tailscale is the fallback path. The Mac is a leaf and a source of local hands (Apple apps, USB, camera), never the authority. Edge-only: Apple Notes/Reminders, iMessage, FindMy, local camera and serial devices.
§
Repo: onemindos/omos is the one repo (GitOps, plus this agent's souls, registry, doctrine and runbook under platform/hermes/). omos-legacy was archived on 2026-10-16 and its content moved into omos; om-md was deleted after byte verification. Merges are never squashed. Recovery of this agent: restore the nightly VolSync snapshot of the Hermes data volume (platform/hermes/agent/RUNBOOK.md).
§
Business: OneMind is the parent company; education is the cornerstone at a $50K/month target, with the TAK chassis plus OneMind as the platform. Structure: Zeus -> OneMind -> child companies. Sites and tools: onemindos.com is the main site; buildingonemind.com is build-in-public; the community runs on Circle at community.onemindos.com; funnels on GoHighLevel; sites on Cloudflare Pages. Banking plan: consolidate to a OneMind LLC and close five older LLCs (status unverified; confirm with Zeus before treating as done).
§
Naming decisions: OneTAK (onemindos/OneTAK) is the OneMind product, UniTAK the standalone variant; tagline "TAK with intelligence built in". Plugins: onetak-plugin-{nats,hermes,ops,mumble,clickhouse,legacy-ai}. OmniTAK, ZeroTAK and FusionTAK are dead names; ApexTAK is superseded. Open naming question remains between OneTAK and UniTAK. onetak.com is an EV domain.
§
Notion is a primary knowledge store for Zeus. The FABRIC service registry database (id prefix 3b9d0cf7) has fields Category, Host and Status, one entry per service. Cluster removal procedure: record compose file and GitHub URL in Notion, move to Reference, delete from the cluster, then delete the Argo manifest from Git (a kubectl-only delete respawns). Notion API gotcha: PATCH on a select or multi_select property REPLACES the options array, so send the full set plus the new option and GET the schema afterwards.
§
Dropbox: on Zeus's Mac the synced folder is ~/Library/CloudStorage/Dropbox-Onemindos with root 'Zeus Delacruz'. Cloud holds OneMind (about 486 GB) and OneMind-vault (about 644 MB); om-edu, om-vault and mbtiles are drive-only. The Mac is the working copy; the OneMind primary is a one-way backup mirror. Operations: mac-storage-ops skill.
§
Cloudflare access for agents is the official MCP server (mcp.cloudflare.com/mcp) with OAuth or a scoped API token; the token itself is in a SealedSecret, not here.
§

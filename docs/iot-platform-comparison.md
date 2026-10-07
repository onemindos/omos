# ThingsBoard and OpenRemote for OMOS

Research date: 2026-10-07. Recommendation, not a deployment approval or measured capacity benchmark. Neither platform has been installed by this work.

## Decision

OpenRemote is the stronger first pilot for OMOS if the priority is independent sites, physical asset context, local automation, and open source ownership. ThingsBoard is the stronger candidate for a central IoT fleet service needing documented horizontal scaling, extensive device lifecycle tooling, and vendor-backed commercial operations. This is an architectural assessment of published capabilities, not a claim that either has been load-tested against OMOS.

Keep the decision reversible: pilot one site with the same device simulator, telemetry volume, outage duration, and command safeguards. Do not install both full stacks into the current memory-constrained hub.

## Capabilities and tradeoffs

| Dimension | ThingsBoard | OpenRemote | OMOS consequence |
|---|---|---|---|
| Primary model | Devices, assets, profiles, relationships, telemetry and rule chains | Custom asset types and attributes, agents, asset hierarchy and rules | Prefer the model that represents actual equipment and operating context clearly |
| Device connectivity | Device APIs include MQTT, HTTP, CoAP, LwM2M and SNMP; Gateway provides protocol integration | Generic HTTP/TCP/UDP/WebSocket/MQTT agents and specific protocols such as KNX and Modbus | Verify each real device protocol; protocol names do not guarantee support for every model |
| Automation | Rule chains, transformations, alarms, device RPC and external integration | Local asset rules and flows; central cross-site automation | Assign each rule to one runtime so n8n, mission-engine and IoT rules do not issue duplicate commands |
| User experience | Rich telemetry/SCADA dashboard and device-fleet workflows | Asset-centric Manager UI, maps and custom project styling | Both require project design, data modeling and access policy |
| AI access | Official MCP server; dedicated user permissions inherited through API key; official n8n node | Official MCP integration invoking Manager API | Use dedicated scoped identities; begin read-only, then explicitly authorize bounded command tools |
| Edge behavior | Dedicated Edge deployment, centrally provisioned configuration, durable uplink queue and reconnect synchronization | Same platform can run locally as gateway/server with local rules, storage and central interconnect | Both support local operation; validate identity, reboot and recovery during prolonged isolation |
| Central horizontal scale | Documented microservices architecture with independent services and failover | Kubernetes supported, but maintainer stated March 2026 that manager HA/horizontal scaling remains work in progress | Do not treat multiple OpenRemote manager replicas as supported HA without newer evidence |
| Complexity at scale | Supporting queue, database, cache and coordination services; optional hybrid time-series storage | Manager, identity, database and proxy services; separate instances per site | ThingsBoard's scale path adds substantial infrastructure; OpenRemote requires planning instance boundaries |
| License | Platform 4.4+ BUSL 1.1, version-specific eventual Apache conversion and production grants | AGPLv3; commercial licensing/services also offered | Review product integration and redistribution before embedding or selling |
| Cost model | Infrastructure plus license/support when grants are exceeded; separate products retain their own terms | Infrastructure, operational support and optional commercial arrangements | Compare total operating cost, including operator time, rather than software price alone |

Sources: [ThingsBoard documentation](https://thingsboard.io/docs/reference/), [OpenRemote asset model](https://openremote.io/product/), [OpenRemote agents](https://docs.openremote.io/docs/1.19.0/user-guide/agents-protocols/overview/), [ThingsBoard MCP](https://thingsboard.io/docs/user-guide/mcp-server/), [ThingsBoard n8n](https://thingsboard.io/docs/user-guide/n8n-node/), [OpenRemote MCP](https://docs.openremote.io/docs/1.29.0/user-guide/ai-mcp-server/).

## Deployment types

| Type | ThingsBoard | OpenRemote | Suggested use |
|---|---|---|---|
| Single host | Native Linux packages or Docker; monolithic server with database | Docker Compose standard stack | Isolated pilot or a site whose recovery target tolerates one host failing |
| Kubernetes pilot | Monolithic deployment with explicit resources and persistent storage | Official Helm charts, manager single replica unless verified otherwise | Argo application, pinned versions, private access and separate database backup |
| HA cloud production | Monolithic cluster or microservices; independently scaled transports, nodes and executors; resilient backing services | HA for backing services can be designed independently; manager restart recovery is distinct from active-active HA | ThingsBoard has the clearer documented central HA path |
| Distributed edge | ThingsBoard Edge, with Docker/native install options and documented Docker Compose clustering | Independent instances acting as edge gateways with central interconnect; amd64/arm64 documented | Keep critical automation local and buffer telemetry for disconnection |
| Vendor hosting | ThingsBoard Cloud managed service | Hosting/support offered through service arrangements | Cloud operations convenience; verify access, exports, region and contract requirements |
| Disconnected environment | Platform 4.4 licensing requires explicit offline arrangement; Edge terms separately verified | Self-hosted local stack; confirm all identity, DNS, package and time dependencies are local | Conduct cold-boot, credential-expiry and prolonged-disconnection acceptance exercises |

Sources: [ThingsBoard installation choices](https://thingsboard.io/docs/installation/), [microservices](https://thingsboard.io/docs/reference/architecture/microservices/), [OpenRemote Helm deployment](https://docs.openremote.io/docs/user-guide/deploying/kubernetes/), [OpenRemote quick start](https://docs.openremote.io/docs/next/quick-start/), [OpenRemote services](https://openremote.io/services/).

The documented ThingsBoard microservices design uses Kafka, Redis/Valkey, coordination, load balancing and PostgreSQL; hybrid storage adds Cassandra. NATS and ClickHouse already in OMOS are not automatically interchangeable with these internal dependencies. Use supported storage/queue choices rather than rewriting the platform to reuse existing components. [Architecture details](https://thingsboard.io/docs/reference/architecture/microservices/).

The March 25, 2026 OpenRemote maintainer reply recommends a single manager pod with Kubernetes restart recovery while HA/horizontal scaling is developed. Current Helm documentation establishes Kubernetes packaging, not an active-active guarantee. No newer explicit support guarantee was established during this research. Partitioning by site/realm or using independent instances requires an integration design; it is not a substitute claim of transparent clustering. [Maintainer statement](https://forum.openremote.io/t/manager-ha-deployment/4199).

## Offline operation and licensing

ThingsBoard Edge stores uplink events locally before transmission and synchronizes when connection returns. User/customer provisioning remains central; evaluate local login and recovery separately. [Edge synchronization](https://thingsboard.io/docs/edge/key-concepts/edge-instance/).

OpenRemote documents independent local rules and storage with a gateway linking to a central instance, including outbound secure tunnel access. Verify historical data reconciliation, command ordering and conflict handling experimentally; the existence of a gateway connection does not establish the guarantees for every data flow. [Edge gateway tutorial](https://docs.openremote.io/docs/1.29.0/tutorials/edge-gateway-secure-tunnel/).

ThingsBoard 4.4 platform licensing changed September 29, 2026. The published commercial free grant is limited to 100 devices and one production instance; education/nonprofit grants have different conditions. Production requires a key. The vendor describes online-key interruptions up to 48 hours, then a management-interface lock, while existing device data ingestion continues. An authorized perpetual offline deployment can use an offline key. Edge, Gateway, TBMQ and other separate products are not covered by the platform relicensing announcement. Confirm exact purchased version and scope before production. [Announcement and offline FAQ](https://thingsboard.io/blog/one-thingsboard-source-available/), [operative BUSL text](https://thingsboard.io/legal/busl/).

OpenRemote AGPLv3 supports open source self-hosting but has source obligations for covered modifications and combinations; the vendor offers commercial alternatives. A loosely coupled API adapter should still be reviewed for its actual distribution/integration design rather than assuming every integration is exempt. [License text](https://github.com/openremote/openremote/blob/master/LICENSE.txt), [vendor licensing explanation](https://openremote.io/open-source/).

## Scale assessment

Device count alone cannot determine capacity. For 10,000 devices reporting once a minute, average ingress is about 167 messages/second; at once a second it is 10,000 messages/second, before payload size, bursts, rules, fan-out, replication or history retention. Measure messages, data points, payload bytes, concurrent sessions, rule work, command latency and storage growth.

ThingsBoard publishes benchmark conditions and suggested deployment scenarios; those are vendor measurements, not promises for this workload. OpenRemote's supported workload should be validated with the intended asset/rule design. Neither can be assigned an honest maximum-device count without a defined workload and acceptance criteria. [ThingsBoard performance methodology](https://thingsboard.io/docs/reference/architecture/performance/).

## Proposed OMOS integration and evaluation

1. Keep NATS as the shared event interface and TAK/OneTAK as the mission/geospatial interface. Add a bounded IoT adapter; do not replace the entire data plane.
2. Establish canonical site/device IDs, units, timestamps, schema version, deduplication identifiers and command acknowledgments.
3. Let the IoT platform own device/asset configuration, with a documented projection into other systems. Choose retention/export rules before duplicating all telemetry into ClickHouse.
4. Local deterministic rules own critical automation. AI proposes or invokes approved actions through scoped APIs, with audit records, expiry and human approval where required.
5. Use n8n for cross-system workflows and notifications; bridge to the IoT API/events. Do not assign the same physical control rule to several engines.
6. Pilot one site: real protocol integration, two roles, telemetry ingestion, alarms, safe command execution, WAN outage, restart while offline, reconnect, backup restore and configuration export.
7. Compare results for idle memory, sustained throughput, p95 command/alert latency, disk growth, lost/duplicate events, recovery time and operational effort.
8. Select OpenRemote if local asset operations and source ownership dominate; select ThingsBoard if fleet lifecycle, central horizontal scale and commercial support dominate. No production installation should precede resource budgeting and the ownership audit.

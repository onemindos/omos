# DDF: Distributed Data Framework

**Link:** https://codice.org/ddf/Features.html  
**Full Site:** https://codice.org/ddf/

**License:** Open Source (FOSS)  
**Maintained by:** Codice Foundation  
**Status:** Production-ready  
**Repository:** GitHub (codice/ddf)

---

## What It Is

Federated, extensible metadata catalog and geospatial data framework built on open standards. Designed for managing, searching, and distributing geospatial and operational data across distributed systems.

**Core Purpose:** Enterprise-scale metadata catalog with OGC (Open Geospatial Consortium) compliance, security, and extensibility.

---

## What It Can Do

### Metadata Catalog & Search
- **Federated catalog** — query across multiple data sources
- **Apache Solr integration** — scalable full-text, spatial, and XML indexing
- **OGC services:** KML, CSW (Catalog Service for the Web), WFS (Web Feature Service)
- **OpenSearch support** — compatible with standard search protocols
- **REST API** — programmatic catalog operations
- **Query transformation** — Enterprise Service Bus converts between query standards

### Search Capabilities
- **Spatial (WKT indexing)** — geometry-aware search
- **Full-text search** — content discovery
- **XPath searches** — XML document queries
- **XML indexing** — structured metadata indexing

### Content Management
- **Integrated content framework** — store actual products (files, images, documents) alongside metadata
- **Tika parser** — automatic metadata extraction from products
- **Plugin pipeline** — pre/post-processing on all catalog operations
- **Eventing & notifications** — CometD-based push notifications for catalog changes

### User Interfaces
- **3D WebGL map** — Cesium-based 3D visualization
- **2D map** — OpenLayers 3 for standard 2D operations
- **USNG/MGRS grid support** — Military/tactical grid reference system
- **GeoNames geocoder** — Integrated location search
- **Upload & edit capability** — Create/modify metadata and features
- **Saved workspaces** — Persist searches and artifact collections
- **Admin console** — Web-based installation wizard and configuration

### Security & Access Control
- **Web Security Service (WSS)** — integrated throughout
- **SAML 2.0 SSO** — browser-based single sign-on with built-in IdP server
- **SAML Security Token Service** — WS-Trust based token management
- **Automatic web app protection** — SSO without application modification
- **XACML PDP** — fine-grained authorization policies
- **Multiple authentication methods:**
  - LDAP (with included OpenDJ LDAP server)
  - X.509 certificates
  - CAS integration
  - Basic authentication
  - SAML authentication
  - Guest login
  - Multi-realm login (mix/match auth types)
- **WS-Security standards:** WS-SecurityPolicy, WS-Policy, WS-Trust, WS-SecureConversation, WS-Addressing
- **Role & attribute-based access control (RBAC/ABAC)** — fine-grained permissions
- **Attribute filtering** — restrict search results based on user attributes
- **Federated identity** — cross-system identity management

### Development & Extensibility
- **OSGi framework** — modular plugin architecture via Apache Karaf
- **Java-based** — build custom features as simple Java objects
- **Dependency injection** — standard Spring/CDI patterns
- **Apache Camel** — message transformation and routing
- **Apache CXF** — web services and REST integration
- **Modular deployment** — only deploy needed features
- **Documentation & components** — extensive resources for custom development
- **Multi-organizational** — designed for distributed team development

### Operations & Installation
- **Unzip and run** — simple deployment model
- **Web-based install wizard** — guided setup
- **Admin UI** — configuration grid for organizing application settings
- **Application marketplace** — pre-built app packages
- **Metrics monitoring** — observe system health and performance

---

## Architecture & Components

### Core Components
1. **Metadata Catalog** — central index of all resources
2. **Apache Solr** — indexing and search engine
3. **Content Store** — file/blob storage for products
4. **ESB (Enterprise Service Bus)** — query transformation layer
5. **REST API** — programmatic access
6. **Web UI** — user-facing search and management

### OGC Compliance
- **OGC CSW** — Catalog Service for the Web (ISO 19115 metadata)
- **OGC WFS** — Web Feature Service (GIS data)
- **OGC KML** — Keyhole Markup Language
- **OpenSearch** — public search standard

### Integration Points
- **LDAP** — directory services
- **SAML** — federated identity
- **CAS** — centralized authentication
- **Tika** — document parsing & extraction
- **CometD** — real-time push notifications
- **Prometheus/metrics** — observability

---

## Key Features

### Standardization
- Built on FOSS and open standards
- Vendor lock-in free
- No proprietary formats

### Extensibility
- Capabilities extend via new feature modules
- Apache Karaf OSGi support
- Apache Camel + CXF integration
- Plugin pre/post-processing

### Simplicity of Development
- Build simple Java objects, wire via DI frameworks
- Widely available documentation
- Modular development for distributed teams

### Flexibility
- Deploy only required features
- Composable architecture
- Optional components (no bloat)

### Installation & Operation
- Unzip and run
- Modern admin console
- Web-based configuration

---

## Use Cases

| Use Case | Why DDF | Benefit |
|----------|---------|---------|
| **Geospatial Data Hub** | OGC-compliant catalog + 3D/2D map | Standard interop, rich visualization |
| **Intelligence/Operational** | SAML SSO, RBAC, multi-realm security | Secure, federated access control |
| **Multi-source Catalog** | Federated search across repositories | Single query interface, distributed data |
| **Document Management** | Tika metadata extraction, full-text search | Auto-indexed content discovery |
| **Scientific Data Archive** | Scalable metadata index, open standards | Research data sharing, preservation |
| **Military/Defense** | USNG/MGRS grids, SAML/X.509, air-gapped | Tactical operations, secure access |
| **Civil Infrastructure** | Geospatial features, asset management | Monitor utilities, roads, water systems |

---

## Resource Requirements

### Minimal Setup
- **CPU:** 2-4 cores
- **RAM:** 4-8 GB (Solr indexing)
- **Disk:** 50+ GB for content store + indices
- **Java:** JDK 8+ required

### Production Deployment
- **CPU:** 8+ cores (Solr scaling)
- **RAM:** 16-32 GB (indexing + cache)
- **Disk:** 100+ GB to 1+ TB (indexed content)
- **Network:** Gigabit+ for metadata federation

### Storage
- **Local filesystem** (development)
- **S3** (cloud-native)
- **NAS/network storage** (enterprise)

---

## Security Architecture

### Authentication Flow
```
User → Web UI
      → SAML IdP/CAS/LDAP
      → SSO session token
      → Verify credentials
      → Grant authenticated state
```

### Authorization Model
```
Authenticated User + Attributes
      → XACML Policy Engine
      → RBAC rules
      → Attribute-based filtering
      → Search results filtered
      → Content access controlled
```

### Multi-Realm Support
- Mix SAML + LDAP + X.509 across endpoints
- Route auth requests to appropriate backend
- Federated identity across systems

---

## Integration with OneMind OS / NATS

### Potential Integrations
1. **Data Stream Ingestion** — NATS → DDF catalog
   - Subscribe to NATS topics
   - Extract metadata
   - Index in Solr
   - Catalog updates via REST API

2. **Federated Query** — Multiple DDF instances
   - Each edge node runs DDF
   - Central DDF aggregates results
   - NATS-based synchronization

3. **Event-Driven Indexing** — NATS events trigger catalog updates
   - File uploaded → NATS event
   - DDF listener indexes → Solr
   - Search discovers immediately

4. **Access Control** — NATS identity ↔ DDF RBAC
   - NATS subject permissions map to DDF attributes
   - Unified access control across fabric

---

## Deployment Patterns

### Pattern 1: Centralized Catalog
```
Cloud: DDF Master Catalog
  ↑
  ├── Edge Node 1 (local DDF) → federation query
  ├── Edge Node 2 (local DDF) → federation query
  └── Edge Node N (local DDF) → federation query
```

### Pattern 2: Federated Search
Multiple DDF instances, each serving local data:
- Each DDF queries own Solr index
- CSW federation queries all nodes
- Results aggregated by master catalog

### Pattern 3: Air-Gapped Operations
- Local DDF instance runs offline
- Manual/batch sync when connectivity available
- Full search capability without cloud dependency

---

## Learning Resources

### Official
- **Site:** https://codice.org/ddf/
- **Features:** https://codice.org/ddf/Features.html
- **Documentation:** (docs link on site)
- **GitHub:** codice/ddf
- **Support:** (support link on site)

### Key Concepts to Learn
1. OGC standards (CSW, WFS, KML)
2. Metadata cataloging principles
3. Geospatial data indexing (Solr)
4. SAML/SSO architecture
5. Java plugin development (if extending)

### Prerequisites
- Geospatial data basics (CRS, vector vs raster)
- OGC service understanding
- Basic Java knowledge (for custom development)
- LDAP/SAML fundamentals (for security)

---

## Comparison with Alternatives

| Aspect | DDF | GeoServer | Elasticsearch | Apache Lucene |
|--------|-----|-----------|----------------|---------------|
| **OGC Compliant** | ✓ Full | ✓ Full | ✗ No | ✗ No |
| **Metadata Catalog** | ✓ Purpose-built | ✗ No | ✗ No | ✗ No |
| **SAML/SSO** | ✓ Integrated | ✓ Plugin | ✗ No | ✗ No |
| **RBAC** | ✓ Yes | ✓ Yes | ✗ No | ✗ No |
| **Federated Search** | ✓ Yes | ✓ Plugin | ✗ No | ✗ No |
| **Extensible** | ✓ OSGi | ✓ Limited | ✓ Plugin | ✓ Plugin |
| **Use Case** | Catalog hub | Map server | Full-text search | Library search |

---

## When to Use DDF

### ✓ Good Fit
- Multi-source metadata catalog needed
- OGC compliance required
- Security (SAML/RBAC) critical
- Geospatial data management
- Enterprise scale indexing
- Federated search across systems

### ✗ Not a Good Fit
- Simple file sharing (use S3/NAS)
- Real-time data stream processing (use Kafka/NATS)
- Small-scale (<10GB metadata) — overhead too high
- Non-geospatial data only (use Elasticsearch)

---

## Configuration Examples

### Enable SAML SSO
```yaml
# Configuration via Admin UI or XML
org.codice.ddf.security.saml.spOnly=true
org.codice.ddf.security.saml.idpServiceName=my-idp
org.codice.ddf.security.saml.spEntityId=my-catalog
```

### Configure Solr Index
```yaml
# Points DDF to Solr instance
org.codice.ddf.solr.client.endpoint=http://solr:8983/solr
org.codice.ddf.solr.client.type=EmbeddedSolrClient
```

### Add Custom Feature
1. Build Java plugin (OSGi bundle)
2. Deploy to `/deploy` directory
3. Enable in Admin UI
4. Restart (or hot-reload if enabled)

---

## Related Technologies

- **Apache Solr** — search/indexing backend
- **Cesium** — 3D mapping library
- **OpenLayers 3** — 2D mapping library
- **Apache Karaf** — OSGi container
- **Apache Camel** — message routing
- **NATS** — event bus (integration potential)
- **GeoServer** — alternative map server
- **PostGIS** — spatial database

---

## Roadmap & Future

- Enhanced cloud-native deployment (Kubernetes)
- GraphQL API support
- Real-time streaming integration
- Advanced ML/AI indexing
- Improved performance at massive scale (10B+ records)

---

**Captured:** 2026-07-09  
**Status:** Production-ready and actively maintained  
**Relevance:** Enterprise catalog, federated systems, geospatial data management  
**Integration Potential:** NATS-driven data ingest, edge federation, OneMind OS integration

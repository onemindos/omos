# Synadia Cloud/NGS Research Summary

## Overview
Synadia is the company behind NATS messaging system and offers cloud-hosted NATS services through their **NGS (NATS Global Service)** platform.

## Service Offerings

### 1. NGS (NATS Global Service)
NGS is Synadia's globally distributed, managed NATS super-cluster service that eliminates the need to run your own NATS infrastructure.

**Connection Endpoint:** `tls://connect.ngs.global`

### 2. Account Tiers

#### **Free Developer Account**
- **Signup time:** Less than 30 seconds
- **Default limits:**
  - 1 leaf node connection (can be increased in account settings)
  - Can be adjusted by navigating to Settings > Limits in the Synadia Cloud dashboard
- **Features:**
  - Full leaf node support
  - TLS encryption
  - Account isolation
  - Multi-tenancy within account boundaries
  - Access to global NGS infrastructure

#### **Paid Accounts**
- **Features:** Enhanced limits and capabilities over free tier
- **Specific pricing:** Not available in public documentation (requires contacting Synadia)
- **Benefits:** Higher connection limits, increased resources

## Key Capabilities

### Leaf Nodes Support ✅
**YES - Full support on both free and paid accounts**

- Leaf nodes allow local NATS servers to connect to the global NGS network
- Local services can communicate within the NGS ecosystem while maintaining isolation
- Messages route transparently across the global network
- Only visible within your account boundaries (multi-tenant isolation)

**Configuration Example:**
```hcl
leafnodes {
    remotes = [
        {
          url: "tls://connect.ngs.global"
          credentials: "/path/to/your.creds"
        },
    ]
}
```

### Acting as Central Hub ✅
**YES - This is the primary use case**

- NGS acts as a global super-cluster that serves as a central messaging hub
- Multiple machines and applications can connect from anywhere
- Leaf nodes from different locations connect to the same NGS account
- Transparent message routing between all connected leaf nodes and clients
- Global distribution with regional endpoints

### SaaS Integrations ✅
**YES - Through standard NATS connectivity**

- Any application/service can connect using NATS client libraries
- Supports 40+ programming languages through official clients
- Docker container support for easy deployment
- Standard NATS patterns work: pub/sub, request/reply, queueing
- JetStream (persistence) support available
- Can integrate with cloud services, on-premises systems, IoT devices, etc.

## Features and Limits

### Account Configuration
- **Management:** Through Synadia Cloud web dashboard
- **User Management:** Create multiple users per account with individual credentials
- **Limit Adjustments:** Can increase leaf nodes and other limits in Settings > Limits
- **Security:** Uses nkeys and JWT-based authentication
- **Credentials:** Downloadable .creds files for each user

### Technical Capabilities
- **TLS Encryption:** All connections encrypted by default
- **Global Distribution:** Super-cluster with worldwide endpoints
- **Multi-tenancy:** Full isolation between accounts
- **Monitoring:** Account-level visibility and metrics
- **Core NATS:** Full support for all messaging patterns
- **JetStream:** Persistence and streaming capabilities
- **Key-Value Store:** Available through JetStream
- **Object Store:** Available through JetStream

### Connection Types Supported
1. **Direct client connections** - Apps connect directly to NGS
2. **Leaf node connections** - Local NATS servers connect as leaf nodes
3. **Mixed environments** - Combination of direct and leaf node connections

## Additional Products

### Synadia Platform / Control Plane
- Self-service management interface
- Alternative to `nsc` command-line tool for security configuration
- Web-based management of Operators, Accounts, Users, and permissions

### Deploy for Kubernetes
- Bring-your-own Kubernetes deployment option
- Includes NATS and additional components
- Self-service deployment

## Use Cases Supported

### Multi-Machine Coordination
- ✅ Connect servers across different data centers
- ✅ Connect cloud and on-premises systems
- ✅ Hybrid cloud architectures
- ✅ Edge computing with central coordination

### Application Hub
- ✅ Microservices communication
- ✅ Service mesh messaging layer
- ✅ Event-driven architectures
- ✅ Request/reply patterns across services

### IoT and Edge
- ✅ IoT device connectivity
- ✅ Edge node communication
- ✅ Telemetry aggregation
- ✅ Command and control

## Documentation and Tools

### Management Tools
- **nsc** - Command-line tool for offline account/user configuration
- **Synadia Cloud Dashboard** - Web-based account management
- **nats CLI** - General NATS administration and testing

### Getting Started
1. Sign up for free account at Synadia Cloud (< 30 seconds)
2. Configure account limits if needed (Settings > Limits)
3. Create users and download credentials
4. Configure leaf nodes or direct connections using credentials
5. Start messaging across your distributed infrastructure

## Pricing Information

### Available Information
- **Free Tier:** Confirmed available with 1+ leaf node support
- **Paid Tiers:** Exist but specific pricing not publicly documented
- **Enterprise:** Likely available through direct contact

### To Get Pricing Details
Contact Synadia directly through:
- Website: https://www.synadia.com/cloud
- Sales inquiry for paid/enterprise plans

## Summary Assessment

### Does Synadia NGS meet the requirements?

✅ **Leaf Nodes:** YES - Full support with configurable limits  
✅ **Central Hub:** YES - This is the core purpose of NGS  
✅ **Multi-Machine:** YES - Global super-cluster designed for this  
✅ **Multiple Applications:** YES - Full multi-tenancy support  
✅ **SaaS Integrations:** YES - Standard NATS client libraries work  

### Strengths
- Globally distributed infrastructure (no need to manage servers)
- Free tier available for testing and development
- Mature technology with extensive client library support
- Full leaf node support for hybrid architectures
- Excellent for distributed systems and microservices

### Limitations Found
- Specific pricing details not publicly available
- Free tier has limited resources (1 leaf node default, adjustable)
- Must contact sales for detailed pricing and enterprise features

### Recommended Next Steps
1. Create a free Synadia Cloud account to test the service
2. Experiment with leaf node connections from your machines
3. Contact Synadia sales for pricing if free tier limits are insufficient
4. Review Synadia Platform/Control Plane for advanced management needs

---

*Research completed: May 21, 2026*  
*Sources: NATS.io documentation, Synadia references in NATS docs*  
*Note: Specific pricing not available in public documentation*

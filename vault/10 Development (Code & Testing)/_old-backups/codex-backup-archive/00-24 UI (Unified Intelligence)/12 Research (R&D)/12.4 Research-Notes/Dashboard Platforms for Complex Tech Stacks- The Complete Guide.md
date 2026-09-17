---
id: ui-15-20260111053926
title: "Dashboard Platforms for Complex Tech Stacks: The Complete Guide"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Dashboard Platforms for Complex Tech Stacks: The Complete Guide"
---

# Dashboard Platforms for Complex Tech Stacks: The Complete Guide

**Bottom line: No dashboard platform currently offers robust built-in documentation management**. However, three platforms stand out for different reasons: **Dashy** excels at iframe-based documentation embedding and authentication flexibility, **Homarr** leads in modern UX with drag-and-drop simplicity and multi-user support, and **Homepage** provides the most extensive service integrations (100+) with native tunnel monitoring widgets. The best approach combines one of these dashboards with a dedicated documentation tool like Notion, Confluence, or BookStack, accessed through secure tunneling solutions like Tailscale or Cloudflare Tunnel.

This matters because trying to find an all-in-one solution will leave you disappointed—the market has evolved toward specialized tools working together rather than monolithic platforms. Understanding the current landscape and integration patterns will save significant time and frustration.

## The documentation management gap

Every dashboard platform reviewed—from the most feature-rich to the simplest—shares a critical limitation: **they organize links to applications, not documentation itself**. Dashy, Homarr, Homepage, and others excel at categorizing services, monitoring uptime, and providing quick access, but offer minimal capabilities for storing guides, notes, or detailed technical documentation.

**Astroluma** comes closest with its built-in snippet manager for code and documentation, nested category organization, and featured pages for important resources. However, even this falls far short of a proper knowledge management system. The reality is that documentation requires rich text editing, versioning, search across content, attachment management, and collaborative editing—features that would bloat a dashboard platform beyond its core purpose.

The **proven solution** combines a dashboard for application access with a dedicated documentation platform. Dashy's iframe widgets enable embedding entire Notion pages or Confluence spaces directly in your dashboard, creating a hybrid workspace where you can view documentation alongside application links. Alternatively, maintain documentation separately and link to specific pages from each dashboard entry—a simpler but effective pattern used successfully across the homelab community.

## Platform recommendations by use case

### For maximum documentation integration: Dashy

**Dashy** stands alone in its ability to seamlessly blend documentation with application management. Its 50+ widget system includes iframe embedding that can display entire Notion databases, Confluence spaces, or internal wikis directly in your dashboard. The workspace view takes this further, allowing simultaneous access to multiple applications and documentation sources in a tabbed interface.

Authentication is where Dashy truly shines for complex tech stacks. Full integration with **Keycloak, Authelia, and Authentik** enables single sign-on across your entire infrastructure. Group-based visibility controls mean different team members see different sections—critical when managing documentation alongside sensitive production services. The `showForKeycloakUsers` and `hideForKeycloakUsers` directives provide granular control down to individual items.

For remote access, deploy Dashy behind Authelia or Authentik via reverse proxy, then access through Tailscale for private networks or Cloudflare Tunnel with Access policies for public exposure. The platform's 22,500+ GitHub stars and active maintenance through 2025 ensure ongoing development and community support.

**Trade-offs**: Requires YAML configuration knowledge. While a UI editor exists, power users will spend time in config files. The learning curve is steeper than GUI-only alternatives, but the payoff is unmatched flexibility.

### For modern UX and team environments: Homarr

The September 2024 release of **Homarr v1.0** represents a complete ground-up rewrite focused on performance and user experience. This platform eliminates YAML entirely—everything configures through **drag-and-drop interfaces** accessible to non-technical users. For teams managing complex tech stacks, the multi-user system with granular permissions and separate boards per project or department provides organizational clarity that other platforms lack.

Built-in authentication supporting credentials, OIDC, and LDAP means no external SSO provider is required, though integration with Authentik or Authelia remains an option. The **multiple boards system** elegantly solves the scale problem: create separate dashboards for production services, development environments, documentation resources, and monitoring tools, each with its own URL and access controls.

Homarr's **30+ native integrations** focus heavily on the homelab ecosystem—Plex, Jellyfin, Sonarr, Radarr, Prowlarr, qBittorrent, Deluge, Proxmox, Docker, and Kubernetes all have built-in widgets displaying real-time data. WebSocket-based updates ensure information stays current without page refreshes, while background task processing maintains responsiveness even with hundreds of monitored services.

**Documentation approach**: While lacking native documentation features, Homarr's notes widget and description fields accommodate brief guides. For comprehensive documentation, link to external tools and use iframe embedding where needed. The platform's strength lies in application management; pair it with Notion or Confluence for documentation.

**Trade-offs**: Heavier resource usage than static alternatives like Homer. Requires database backend (typically SQLite), adding deployment complexity compared to pure static solutions.

### For simplicity and maximum integrations: Homepage

If your priority is monitoring and accessing 100+ services with minimal configuration overhead, **Homepage** dominates this space. The platform's **100+ service integrations** dwarf all competitors—every major self-hosted application from Radarr to Home Assistant to GitLab to Kubernetes has a pre-built widget displaying live operational data.

Docker auto-discovery via labels enables **near-zero configuration** for containerized environments. Add labels to your Docker Compose files, and Homepage automatically discovers and adds services with proper icons and status monitoring. This scales elegantly to complex environments where manually configuring 100+ applications would be prohibitive.

The platform's static architecture means **zero backend requirements**—deploy anywhere you can serve HTML/JavaScript. Load times are near-instantaneous, resource usage is minimal (perfect for Raspberry Pi deployments), and the entire configuration exists in YAML files easily version-controlled with Git. Native **Cloudflare Tunnel and Tailscale widgets** provide at-a-glance status of your secure remote access infrastructure.

**Critical limitation**: Homepage has **no built-in authentication** by design. The documentation explicitly states deployment should occur behind authenticated reverse proxies or within private VPN networks. For private homelab access via Tailscale, this simplicity is an advantage. For public exposure, you must implement Cloudflare Access, Authelia, or similar authentication layers externally.

**Trade-offs**: YAML-only configuration. No web-based editor exists; all changes require editing text files. For teams with non-technical members, this creates a maintenance bottleneck. The bookmarks system can link to documentation but offers no embedding or rich content display.

## Documentation integration patterns

### Pattern 1: Dashboard with embedded documentation (Dashy's strength)

Deploy Dashy as your central hub with iframe widgets embedding documentation platforms. Configure sections for different technology domains—infrastructure, applications, monitoring, development—and embed relevant Notion pages or Confluence spaces within each section using iframe widgets. The workspace view enables opening multiple documentation sources simultaneously while retaining access to application links in the sidebar.

This approach works best when your documentation lives in web-accessible platforms like **Notion (with public or authenticated access), Confluence, Outline, or WikiJS**. Configure Dashy's OIDC integration to match your documentation platform's authentication, creating a seamless single sign-on experience. Users authenticate once and gain access to both application links and embedded documentation without additional logins.

**Implementation**: Each major application or service group gets a section with description text, application links, monitoring widgets, and an embedded iframe displaying relevant documentation. Users see everything they need—access URLs, current status, resource metrics, and complete guides—without leaving the dashboard.

### Pattern 2: Hybrid approach with external documentation links

Maintain documentation in a dedicated platform (**BookStack, Outline, Obsidian Publish, or Notion**) and link to specific pages from each application entry in your dashboard. This simpler pattern works with any dashboard platform—Homepage, Homarr, Dashy, or others—since it only requires URL links rather than embedding capabilities.

Structure your documentation with a dedicated page for each major application containing setup guides, troubleshooting steps, API documentation, and operational runbooks. In your dashboard configuration, add these documentation URLs in description fields or as separate link items grouped with the related applications. For Homarr, utilize the notes widget to display brief summaries with links to full documentation.

**Advantages**: Clean separation of concerns, easier maintenance, no iframe security considerations, works regardless of authentication complexity. Teams already invested in Confluence or Notion can continue using familiar tools while gaining the organizational benefits of a modern dashboard.

### Pattern 3: Reverse proxy with subdomain architecture

Deploy dashboard and documentation as separate services behind a unified reverse proxy with single sign-on:

- **dashboard.yourdomain.com** → Homepage, Homarr, or Dashy
- **docs.yourdomain.com** → Outline, BookStack, or WikiJS  
- **app1.yourdomain.com, app2.yourdomain.com** → Individual services

Configure **Traefik or Nginx with Authelia/Authentik** providing authentication for all services. Users authenticate once and access everything through their chosen subdomain. This enterprise-grade pattern provides maximum flexibility—each service remains independently maintainable while the authentication layer ensures consistent security.

The dashboard links to documentation using internal domains, and both remain accessible via the same secure remote access solution (Tailscale, Cloudflare Tunnel, or VPN). For organizations with 50+ applications, this architecture scales better than monolithic approaches and allows different teams to maintain their respective services.

## Secure remote access integration

### Private access: Tailscale as the gold standard

**Tailscale** has emerged as the definitive solution for private homelab access in 2025. The WireGuard-based mesh VPN creates a private network (tailnet) where all your devices receive 100.x.y.z addresses and can communicate directly across NAT, CGNAT, and firewall boundaries without any port forwarding. Installation takes minutes—install clients on your devices, authenticate with your identity provider, and devices automatically discover each other.

For dashboard access, install Tailscale on your dashboard host and access it via the Tailscale IP from any device on your tailnet. **Tailscale Serve** provides HTTPS access with automatic TLS certificates for ts.net domains, eliminating self-signed certificate warnings. Homepage includes a **native Tailscale widget** displaying device count and connection status directly in your dashboard.

**MagicDNS** enables human-readable names: access `dashboard.tailnet-name.ts.net` instead of remembering IP addresses. Access control lists (ACLs) provide granular control over which devices can access specific services. For complex tech stacks, configure **subnet routing** so one Tailscale node acts as a gateway, providing tailnet access to your entire home network without installing Tailscale on every device.

The free tier supports **100 devices and 3 users**—more than sufficient for personal use and small teams. For those requiring self-hosted control, **Headscale** provides an open-source implementation of Tailscale's coordination server, eliminating dependence on Tailscale's infrastructure while maintaining protocol compatibility.

### Public access: Cloudflare Tunnel with Access

When you need to share your dashboard or services with users who can't install VPN clients, **Cloudflare Tunnel** provides secure public access without opening firewall ports. The cloudflared daemon runs on your network, establishing outbound-only connections to Cloudflare's edge network. All incoming traffic routes through Cloudflare's infrastructure, providing DDoS protection and automatic SSL certificates.

**Cloudflare Access** adds an authentication layer before the tunnel, supporting identity providers like Google, GitHub, Azure AD, or generic OIDC/SAML. Configure policies requiring specific email addresses, domains, or groups before granting access. This works with **Dashy's multi-tier authentication**—combine Cloudflare Access at the network edge with Dashy's OIDC integration for defense in depth.

Homepage's **native Cloudflare Tunnel widget** displays tunnel status and origin IP, providing operational visibility. For comprehensive dashboard deployments, monitor tunnel health alongside the services it protects. The combination of Cloudflare Tunnel for connectivity and Cloudflare Access for authentication eliminates traditional VPN complexity while providing enterprise-grade security.

**Critical security note**: Even with Cloudflare protection, never expose dashboards with default credentials or no authentication. Multiple community reports document compromised systems where services were assumed safe behind Cloudflare but lacked application-level authentication. Always implement defense in depth.

### Self-managed SSO: Authelia and Authentik

For teams wanting centralized authentication without cloud dependencies, **Authelia** and **Authentik** provide self-hosted identity and access management. Both platforms act as authentication layers in front of your services, integrating with reverse proxies to gate access before requests reach your dashboards.

**Authelia** takes a lightweight, configuration-file approach. Define access control rules in YAML specifying which users or groups can access which applications. Multi-factor authentication via TOTP, WebAuthn (hardware keys and passkeys), or Duo provides strong security. The platform weighs under 20MB containerized and uses approximately 30MB RAM—negligible overhead. OpenID Connect provider capabilities enable **native integration with Dashy and Homarr**, providing single sign-on across your entire infrastructure.

**Authentik** offers a heavier but more feature-rich alternative with a graphical admin console, visual flow builder for custom authentication workflows, and built-in user directory. SAML support alongside OIDC makes it suitable for enterprise environments or integration with commercial software. The **user dashboard** provides self-service password reset and device management, while the admin console enables team-based management without editing configuration files.

Dashy provides **comprehensive integration guides** for both platforms with example configurations and group-based visibility controls. Deploy Authelia or Authentik behind your reverse proxy (Traefik recommended for automatic Let's Encrypt certificates), configure your dashboard as an OIDC client, and users access everything through a single login flow.

**Comparison**: Choose Authelia for lightweight, configuration-as-code deployments where resource efficiency matters. Choose Authentik when GUI management, complex authentication flows, or SAML support is required. Both provide production-ready authentication for complex tech stacks.

## Scaling to 100+ applications

### Organizational strategies

Managing 100+ applications in a single flat list becomes unusable quickly. The leading platforms provide different organizational approaches:

**Homarr's multiple boards** excel here—create separate dashboards for infrastructure services, media applications, development tools, monitoring systems, and documentation resources. Each board gets its own URL and can have different access permissions, enabling team segmentation where frontend developers see development tools while infrastructure teams access monitoring dashboards. Within boards, the **drag-and-drop grid** enables visual organization with categories that make sense for your specific workflow.

**Dashy's tag system** combined with instant search provides a different approach. Tag applications with multiple attributes—technology stack, environment (prod/dev/staging), team ownership, criticality level—and use the search functionality to instantly filter. The **sections system** groups related applications visually while maintaining a single unified dashboard. For power users, **custom keyboard shortcuts** (0-9 hotkeys for most-used apps) enable instant access without visual navigation.

**Homepage's Docker auto-discovery** reduces maintenance burden as your stack grows. Add labels to containers, and Homepage automatically creates entries with proper icons and status monitoring. The services configuration file supports **unlimited sections**, each collapsible to reduce visual clutter. Bookmark separate files for documentation resources or less-frequently accessed services.

### Performance considerations

**Homepage** maintains instant load times regardless of scale due to its static architecture. All processing happens at build time; serving is pure HTML/JavaScript with no database queries. For monitoring 100+ services, background tasks fetch status asynchronously, preventing slow services from blocking dashboard rendering.

**Homarr v1.0** addressed performance specifically in its rewrite, implementing background job systems and WebSocket-based updates. The platform **scales to hundreds of users** with proper database backend (PostgreSQL recommended for high-scale deployments over SQLite). Redis-backed real-time updates ensure integration data stays current across all connected clients.

**Dashy** handles large configurations through lazy loading and caching. The workspace view limits simultaneously loaded applications to those in active tabs, preventing resource exhaustion from iframe-heavy configurations. Status checking occurs asynchronously with configurable intervals to balance freshness against network overhead.

## Platform maintenance and ecosystem health

### Active development indicators

**Glance** (28,400+ stars) represents the newest entry, created in April 2024 with extremely active development. The platform takes a content-feed approach, aggregating RSS, Reddit, YouTube, and other information sources alongside application links. While not a traditional dashboard, its rapid adoption and focus on information density makes it worth monitoring for future capabilities.

**Homepage** (25,800+ stars) moved from benphelps/homepage to the gethomepage organization in 2024, signaling sustainable maintenance beyond a single developer. Regular releases continue through 2025 with consistent community contributions. The integration library expands continuously—any popular self-hosted application typically gains Homepage support within weeks of release.

**Dashy** (22,500+ stars) maintains consistent updates with v3.0 released recently requiring migration. Active GitHub Discussions show ongoing community engagement. The developer (Lissy93) maintains multiple related projects and actively responds to issues and feature requests.

**Homarr** underwent the most significant evolution—complete v1.0 rewrite in September 2024 moving from ajnart/homarr to the homarr-labs organization. This architectural improvement addressed technical debt and positioned the platform for long-term growth. Development velocity remains high with regular feature additions.

### Maintenance concerns

**Heimdall** shows concerning signs of stagnation with minimal commit activity in recent months despite 7,900+ stars. Community discussions mention exploring alternatives. While the platform remains functional and stable, new deployments should consider more actively maintained options.

**Flame** faces similar concerns with GitHub issue #378 explicitly asking "is this project dead?" Development has slowed significantly since late 2023, and active forks are attempting to continue development. The 5,500+ stars indicate historical popularity, but current maintenance status suggests caution for new deployments.

**Organizr** continues functioning but development velocity has decreased. The PHP codebase feels dated compared to modern alternatives built on Node.js or Go. For existing deployments, it remains stable; for new projects, consider Homarr for similar multi-user capabilities with modern architecture.

## Comparison matrix

| Feature | Homepage | Homarr | Dashy | Astroluma |
|---------|----------|---------|--------|-----------|
| **Documentation features** | Links only | Notes widget | Iframe embedding | Snippet manager |
| **Cloud integration** | API-based | Limited | Iframe embedding | Limited |
| **Scale (100+ apps)** | ✅ Excellent | ✅ Excellent | ✅ Excellent | ✅ Good |
| **Configuration** | YAML | Drag-drop | YAML + UI | Web UI |
| **Authentication** | External only | Built-in + OIDC | Built-in + OIDC | Multi-user |
| **Integrations** | 100+ | 30+ | 50+ | Limited |
| **Remote access widgets** | ✅ CF + Tailscale | ❌ | ❌ | ❌ |
| **Docker discovery** | ✅ Auto | Manual | Manual | Manual |
| **Resource usage** | Minimal | Low-medium | Medium | Medium |
| **Setup complexity** | Medium | Low | Medium | Medium |
| **Team features** | ❌ | ✅ Permissions | ✅ Groups | ✅ Multi-user |
| **Active development** | ✅ Very active | ✅ Very active | ✅ Active | ✅ Active |
| **GitHub stars** | 25,800 | 6,900+ | 22,500 | Growing |

## Final recommendations

### Best overall solution for documentation-heavy environments

**Deploy Dashy with Authelia for SSO, accessed through Tailscale for private use or Cloudflare Tunnel for shared access.** Maintain documentation in Notion or Outline and embed pages via Dashy's iframe widgets. This combination provides comprehensive application management, flexible authentication, secure remote access, and integrated documentation display in a single interface.

Configure sections in Dashy for infrastructure, applications, monitoring, and documentation, with embedded iframe widgets displaying relevant Notion databases or Outline pages. Use workspace view for workflows requiring simultaneous access to multiple applications and documentation sources. Implement group-based visibility controls so different teams see appropriate sections.

### Alternative for teams prioritizing ease of use

**Deploy Homarr with built-in authentication, accessed through Tailscale.** Link to external documentation tools (BookStack or Notion) from application entries. Create separate boards for production, development, and documentation resources. This approach prioritizes simplicity—no YAML configuration, no external authentication dependencies, and an interface approachable for non-technical team members.

The drag-and-drop configuration enables team members to organize their own boards without editing text files. Built-in authentication with role-based access control provides sufficient security for most environments without SSO complexity. As your team grows or requirements increase, add OIDC integration with Authentik.

### For maximum integration breadth

**Deploy Homepage accessed through Tailscale or behind Authelia via reverse proxy.** Leverage Docker auto-discovery to minimize configuration overhead. Maintain documentation separately in BookStack or WikiJS, linking to specific pages from Homepage bookmark entries. This scales effortlessly to 100+ applications while maintaining minimal resource usage and instant load times.

Accept that Homepage lacks native authentication—this is the trade-off for its simplicity and performance. For private homelab access via Tailscale, no authentication is necessary. For shared access, deploy Authelia with Nginx or Traefik providing the authentication layer. The static architecture ensures Homepage itself requires zero maintenance beyond configuration updates.

### Security implementation checklist

**Never expose dashboards directly to the internet via port forwarding.** This cannot be emphasized enough—community reports consistently document compromised systems that relied on application hardening without network-layer protection.

For **private access**: Install Tailscale on dashboard host and accessing devices. Configure MagicDNS for human-readable names. Enable Tailscale SSH for secure shell access. Use subnet routing if accessing devices that can't run Tailscale clients.

For **public access**: Deploy Cloudflare Tunnel with cloudflared daemon. Configure Cloudflare Access with identity provider (Google, GitHub, or generic OIDC). Implement IP restrictions where applicable. Add application-layer authentication (Dashy OIDC, Homarr built-in) for defense in depth.

For **SSO across services**: Deploy Authelia (lightweight) or Authentik (feature-rich) behind reverse proxy. Configure dashboard as OIDC client. Enable MFA (TOTP and WebAuthn passkeys). Implement access control policies per application. Regular audit user permissions and access logs.

**Network segmentation**: Isolate dashboard and services on management VLAN. Implement firewall rules between VLANs. Separate IoT devices from production services. Place dashboard on network with appropriate access to all monitored services.

**Backup and recovery**: Version control YAML configurations with Git. Regularly backup dashboard databases (Homarr, Astroluma). Test restoration procedures. Document authentication configuration. Store secrets in environment variables, never commit to version control.

## Conclusion

The dashboard platform ecosystem has matured significantly, but **the critical insight is recognizing that no single platform solves all problems**. Modern approaches combine specialized tools—Dashy or Homarr for application management, Notion or Outline for documentation, Tailscale or Cloudflare Tunnel for secure access, and Authelia or Authentik for centralized authentication.

Homepage leads in integration breadth and simplicity for technical users comfortable with YAML. Homarr provides the best experience for teams needing GUI configuration and built-in multi-user support. Dashy offers maximum flexibility for power users willing to invest in configuration. All three handle 100+ applications effectively with proper organizational strategies.

The **documentation management gap** is real and universal—no dashboard platform offers wiki-grade documentation features. The hybrid approach combining dashboards with dedicated documentation platforms has emerged as the community-validated pattern. Dashy's iframe embedding provides the tightest integration, but simple linking approaches work effectively with any dashboard.

For secure remote access, **Tailscale has become the de facto standard** for private homelab access, while Cloudflare Tunnel serves public-facing use cases. Both eliminate traditional port forwarding and VPN complexity while providing superior security. The combination of mesh VPN technology with identity-aware authentication represents a fundamental improvement over legacy approaches.

Choose your platform based on your team's technical capabilities, scale requirements, and documentation integration priorities. Implement secure remote access from day one—retrofitting security is harder than building it in. Start simple with Tailscale and one dashboard, then expand to more sophisticated patterns as needs grow.
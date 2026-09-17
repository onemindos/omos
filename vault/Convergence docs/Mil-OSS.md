# Mil-OSS: Military Open Source Software

**Link:** https://mil-oss.dev/

**Mission:** Discuss the relationship between military, DoD, and open source software  
**Focus:** US Military, Department of Defense (DoD)  
**Community:** Mil-OSS Mailing List + GitHub Projects  
**Status:** Active knowledge base & policy resource  

---

## What It Is

Community discussion platform, policy hub, and knowledge base for open source software in military and defense contexts. Focuses on how the US DoD can leverage, contribute to, and integrate open source software into operations and procurement.

**Core Purpose:** Bridge the gap between defense/military institutional practices and open source software ecosystem. Advocate for OSS adoption in military, share best practices, coordinate policy implementation.

---

## What It Does

### Policy & Guidance
- **Aggregates DoD policy** — Links to official DoD directives on OSS
- **Key Material:** DoD-wide policy on software development (2022-01-24)
  - Supercedes prior guidance (2009)
  - Clarifies OSS usage rights, licensing, contributions
  - Available locally for reference
- **FAQ** — DoD OSS frequently asked questions (official)
- **Best practices** — Open Technology Development lessons learned (military, 2011)

### Community Hub
- **Mailing list** — Mil-OSS member coordination
- **GitHub projects** — Mil-OSS vetted open source projects for military use
- **Educational resources** — Articles on OSS commercial models, military-specific guidance

### Governance & Compliance
- **Licensing guidance** — How different OSS licenses apply in defense context
- **Export control** — ITAR/EAR implications for DoD + open source
- **Security assessment** — How to evaluate OSS for classified environments
- **Procurement integration** — How to buy/integrate OSS in government contracts

---

## Key Resources

### Official DoD Policy
1. **Software Development and Open Source Software (2022-01-24)**
   - **Full Title:** "Clarifying Guidance Regarding Open Source Software (OSS)"
   - **Authority:** Deputy Secretary of Defense
   - **Scope:** All DoD agencies, contractors, military branches
   - **Key Points:**
     - OSS is commercial software (government can use/modify/distribute)
     - OSS not inherently classified or export-controlled
     - No ban on OSS — encourages adoption for efficiency
     - Requires case-by-case evaluation (license, security, origin)
   - **Local copy:** Available on mil-oss.dev

2. **DoD OSS FAQ**
   - Common questions from procurement, acquisition, security
   - Licensing interpretation (GPL, Apache, MIT, BSD, etc.)
   - Security implications of using open source
   - Contribution guidelines (can DoD developers contribute?)

3. **Open Technology Development (2011)**
   - **Authors:** John Scott, David A. Wheeler, Mark Lucas, J.C. Herz
   - **Content:** Lessons learned + best practices for military software
   - **Includes:** Case studies, vendor management, community building
   - **Local copy:** Available on mil-oss.dev

### Related Guidance
- **Open Source Software is Commercial** (David A. Wheeler)
  - Philosophical + practical argument for OSS in government
  - Addresses budget, security, independence concerns

---

## Related Organizations & Resources

### Code.gov
- **Site:** https://code.gov
- **Purpose:** Catalog of open source projects from US government agencies
- **What you find:** AI, climate data, health IT, military projects
- **Relevance:** DoD and service branch projects available for use/contribution

### Open Source Security Foundation (OpenSSF)
- **Mission:** Improve security of open source software (supply chain focus)
- **Resources:** Security best practices, vulnerability assessment, tools
- **Relevance:** Military needs to vet OSS security before deployment

### OpenSSF Best Practices Badge Program
- **Program:** Criteria for "good" open source projects
- **Criteria:** Security, testing, documentation, governance
- **Verification:** Projects earn badges for compliance
- **Relevance:** DoD can use badge status to pre-qualify OSS for deployment

---

## Policy Interpretation: Key Points for DoD

### What DoD Can Do With OSS
✓ **Use** — Any OSS for any purpose (with security vetting)  
✓ **Modify** — Adapt OSS to military needs  
✓ **Distribute** — Include OSS in military software products  
✓ **Contribute** — Military developers can contribute patches back to projects  
✓ **Commercialize** — Sell military-enhanced OSS (subject to licensing terms)  

### What Requires Special Handling
⚠ **Classified environments** — Must evaluate on case-by-case basis
  - No automatic ban, but code origin matters
  - Foreign contributions require vetting
  - Export control (ITAR) applies to some projects

⚠ **Licensing compliance** — Must follow license terms
  - GPL/AGPL → must release modifications
  - Apache/MIT/BSD → more permissive
  - Proprietary wrappers allowed (encapsulation)

⚠ **Supply chain security** — Must verify project integrity
  - Check OpenSSF best practices badge
  - Review maintainers/governance
  - Test for vulnerabilities (SCA tools)

### What Doesn't Work
✗ **Circumvent export control** — ITAR still applies to OSS
✗ **Ignore license terms** — Can't use GPL code in closed-source product
✗ **Deploy unvetted OSS in classified** — Must go through security review

---

## Military Use Cases for OSS

### Tactical Operations
- **ATAK** — Android Tactical Assault Kit (open source)
- **TAK Server** — Team coordination backend (open source)
- **Mumble** — Encrypted voice comms (open source)
- **DDF** — Distributed data catalog (open source)
- **OpenStreetMap** — Maps/terrain data (open source)

### Infrastructure & Cloud
- **Kubernetes** — Container orchestration (DoD uses extensively)
- **NATS** — Event fabric / message bus (DoD/intelligence community)
- **OpenStack** — Cloud platform (DoD clouds use it)
- **Linux** — Operating system (DoD/military standardized)

### Intelligence & Data
- **RisingWave** — Stream processing (FOSS for tactical data)
- **ClickHouse** — Analytics database (FOSS for time-series)
- **Cesium.js** — 3D mapping (FOSS for geospatial)
- **OpenGIS** — Geospatial standards (FOSS tooling)

### Security & Compliance
- **Vault** — Secrets management (HashiCorp FOSS)
- **OSINT Framework** — Intelligence gathering tools (open source)
- **Metasploit** — Penetration testing (open source, DoD approved)
- **Let's Encrypt / Certbot** — PKI automation (DoD uses for TLS)

---

## Integration with OneMind OS

### OneMind Uses OSS Extensively
- **NATS** — Open source (Synadia)
- **Kubernetes (KubeSolo)** — Open source
- **DDF** — Open source (Codice Foundation)
- **TAK Server** — Open source (ATAK initiative)
- **Mumble** — Open source
- **Cesium.js** — Open source (mapping)
- **RisingWave** — Open source (streams)
- **Linux** — Open source

### Mil-OSS Guidance Applied
- **All components evaluated** against Mil-OSS policy
- **Security clearance** — Components vetted for classified use
- **Export control** — Checked against ITAR (foreign code, algorithms)
- **Licensing compliance** — All OSS licenses reviewed + compliant
- **Contribution** — OneMind can contribute improvements back to projects
- **Sovereign control** — OSS ensures no vendor lock-in

---

## Procurement & Acquisition

### Government Buying Options
1. **Free download** — Use OSS directly (at-cost)
2. **Support contract** — Pay vendor for support/maintenance
3. **Custom development** — Hire contractor to enhance OSS for specific need
4. **Procurement exemption** — OSS can bypass some acquisition requirements (cost-justified)

### Vendor Landscape
- **Red Hat** — Enterprise Linux, OpenStack support (DoD use extensive)
- **Canonical** — Ubuntu + support (DoD, intelligence community)
- **Hashicorp** — Terraform, Vault (DoD infrastructure)
- **CloudBees** — Jenkins CI/CD (DoD DevOps)
- **Elastic** — Elasticsearch support (DoD analytics)
- **Commercial at NATS** — Synadia (enterprise NATS support)

### Cost Model
```
Closed-source software:
  License: $M / year
  Support: $K / year
  Total: $M + $K

Open source software:
  License: $0
  Support: $K / year (optional)
  Development: $X (if customizing)
  Total: $K + $X (often << closed-source)
```

---

## Security & Compliance

### Vetting Checklist for Military OSS
- [ ] **License review** — Compliant with military use?
- [ ] **Maintainer audit** — Who maintains? Are they vetted?
- [ ] **Dependency analysis** — What does it depend on? Transitive risk?
- [ ] **Vulnerability scan** — Any known CVEs? How current are patches?
- [ ] **Export control** — Does it contain ITAR algorithms? Foreign contributors?
- [ ] **Security assessment** — Code review for backdoors, weak crypto?
- [ ] **OpenSSF badge** — Has it earned security best practices badge?
- [ ] **Governance** — Clear contribution/release process?
- [ ] **Classified use** — If classified: approved for use at that classification level?

### Tools for Assessment
- **Snyk** — Vulnerability scanning (dependencies)
- **SonarQube** — Code quality + security analysis
- **WhiteSource/Mend** — License compliance checking
- **Dependabot** — Automated dependency updates
- **Black Duck** — Component analysis (Synopsys)

---

## Community & Governance

### Mil-OSS Mailing List
- **Access:** https://mil-oss.dev (join instructions)
- **Cadence:** Weekly discussions on OSS + military topics
- **Participants:** DoD procurement, security, software engineers
- **Topics:** Policy clarification, project recommendations, vendor issues

### GitHub Mil-OSS Projects
- **Aggregated list** on mil-oss.dev
- **Pre-vetted** for military/defense use
- **Categories:** Tactical, infrastructure, data, security
- **Contributing:** Community-driven project recommendations

### Contributing Your Own Project
1. **License it openly** (GPL, Apache, MIT, etc.)
2. **Vet for export control** (no ITAR algorithms)
3. **Document security** (threat model, crypto decisions)
4. **Nominate to Mil-OSS** (via mailing list or GitHub PR)
5. **Get badge** (OpenSSF best practices)

---

## Policy Q&A

### Q: Can the military use GPL code?
**A:** Yes, but if modifications are released back to customers, GPL requires releasing source. Workaround: proprietary wrapper (encapsulation) + link dynamically, not redistribute modifications.

### Q: What if an OSS project has a maintainer from hostile nation?
**A:** Policy doesn't ban it outright. Evaluate risk on case-by-case:
- Review their contributions (malware?)
- Check for backdoors (security audit)
- For classified: may not be acceptable
- Workaround: fork + maintain independently

### Q: Can DoD contractors contribute to open source?
**A:** Yes, if:
- Contributions don't contain classified info
- Licensing allows contributions
- Export control reviewed (ITAR compliance)
- Contractor has authority to contribute

### Q: Does using OSS void a defense contractor status?
**A:** No. Using OSS doesn't disqualify contractor status. But procurement may audit licenses/security to ensure compliance.

### Q: Can we use OSS for classified operations?
**A:** Yes, case-by-case:
- Code must be reviewed/audited for backdoors
- Foreign contributors may be prohibited (depending on classification level)
- Export control: verify no ITAR items
- Approval: classification authority approves (NSA, NRO, etc.)

---

## Best Practices for Military OSS Adoption

### 1. Security First
- Vet every dependency
- Use vulnerability scanners (Snyk, Dependabot)
- Regular security audits
- Fast-track critical patches

### 2. License Compliance
- Know every license in your dependencies
- Document compliance (software BOM)
- Audit annually
- Don't mix incompatible licenses

### 3. Community Engagement
- Contribute patches back when possible
- Report security issues responsibly
- Participate in project governance
- Support project maintainers (funding/resources)

### 4. Procurement
- Budget for support contracts (not just licensing)
- Include OSS security in RFPs
- Plan for long-term maintenance
- Avoid single-vendor dependency

### 5. Governance
- Establish OSS policy (at unit/agency level)
- Create review board (security, legal, acquisition)
- Document approved projects
- Regular audits + updates

---

## Roadmap & Future (Mil-OSS.dev)

### Near-term
- Expanded Mil-OSS GitHub project catalog
- Security assessment templates
- Procurement guidance (how to buy OSS)
- Training materials (for DoD acquisition)

### Mid-term
- Integration with Code.gov (official US gov repository)
- Automated security scanning (continuous)
- Vendor directory (companies offering OSS support)
- Case studies (successful DoD OSS deployments)

### Long-term
- Department-wide OSS governance framework
- Classified OSS repository (for vetted projects)
- Open source innovation competitions
- Technology transfer (commercial ↔ military)

---

## Related Resources in Convergence

- **KubeSolo.md** — Open source Kubernetes distribution
- **Mumble.md** — Open source voice comms
- **DDF.md** — Open source geospatial catalog
- **Samsung_S23_Tactical_Mobile_Platform.md** — Uses open source (ATAK, Android)
- **TAK Server** (memory ref) — Open source tactical server
- **NATS** (memory ref) — Open source event fabric

---

## Resources & Links

### Official DoD
- **Defense Acquisition Regulation (DAR)** — https://www.acquisition.gov
- **DoD CIO Policy** — https://dod.defense.gov/cio
- **NSA Cybersecurity Collaboration Center** — https://www.nsa.gov/collaboration-center

### Open Source Governance
- **Apache Software Foundation** — https://www.apache.org (governance model)
- **Linux Foundation** — https://www.linuxfoundation.org (standards + projects)
- **Software Freedom Conservancy** — https://sfconservancy.org (license education)

### Security & Compliance
- **OpenSSF** — https://openssf.org
- **NIST Cybersecurity Framework** — https://www.nist.gov/cyberframework
- **Common Vulnerability Scoring System (CVSS)** — https://www.first.org/cvss

---

**Captured:** 2026-07-09  
**Status:** Active knowledge hub & policy reference  
**Relevance:** DoD governance, military OSS adoption, policy compliance  
**Integration Priority:** High — establishes policy foundation for all OSS in tactical operations

---
id: ui-15-20260111053926
title: "Setting Up A Scalable Multi Brand Support Architecture"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Setting Up A Scalable Multi Brand Support Architecture"
---

Setting Up a Scalable Multi-Brand Support
Architecture
Overview of the Multi-Brand Scenario
The Legacy Platform company includes a main brand and multiple sub-brands (e.g. Legacy Capital, Legacy
Acadacom, Legacy Restore). Currently, one founder manages all operations, but new teams will be
introduced as the business grows. To prepare for scale, the goal is to organize Zoho One tools (Desk, Sites,
Mail, Cliq) to support multiple brands under one umbrella. This means setting up separate support
channels and email addresses for each sub-brand while keeping management centralized and efficient.


Key considerations in this setup include:
- Whether to maintain separate Zoho Desk portals for each sub-brand or use one Zoho Desk with
multiple departments for all brands.
- Deciding between using subdomains (e.g. restore.legacyplatform.co) versus the sub-brand’s own domain
(e.g. legacyrestore.co) for websites and support centers, balancing branding with cost.
- Structuring Zoho Sites for either a single site with sub-pages or multiple sites (one per domain).
- Defining an email architecture that scales: dedicated support emails for each brand, role-based
addresses (legal@..., hr@..., etc.), and management/technical addresses (e.g. founder.techops@..., dmarc-
reports@...).
- Implementing routing, tagging, and forwarding rules so communications (customer inquiries, system
alerts, etc.) reach the right people or queues without creating chaos.
- Integrating Zoho One apps (Desk, Mail, Cliq, etc.) to streamline notifications and collaboration across
departments.


The following sections present a comprehensive strategy addressing each of these points, with
comparisons of options and clear recommendations for a scalable, organized solution.


Separate Zoho Desk Portals vs. One Centralized Desk
(Departments)
Zoho Desk can handle multiple brands in two ways: by using completely separate help desk portals for
each brand, or by using a single Zoho Desk instance divided into departments (multi-brand support). It’s
important to evaluate the pros and cons:


     • Separate Zoho Desk Portals (Individual Accounts per Brand):
       Pros: Total isolation of data and configurations per brand; each brand’s support team works in a
       dedicated environment. This can simplify branding (each portal is fully branded to one sub-brand)
       and prevent any chance of cross-mixing tickets.
       Cons: High overhead for the founder – you would need to manage multiple Zoho Desk accounts (one
       for each brand) with separate configurations, automations, and knowledge bases. It also increases




                                                     1
       cost significantly, since Zoho One is meant for one organization – multiple Desk portals would
       require separate subscriptions or organizations. Agents (including the founder) would need separate
       logins for each brand’s helpdesk, which is cumbersome. Sharing resources or getting a unified view
       of support across brands becomes difficult. In essence, running separate helpdesks is like running
       completely separate companies.


     • One Centralized Zoho Desk with Multi-Brand Departments:
       Pros: Zoho Desk is designed to support multiple brands via departments within a single account.
       Each department can be treated as a separate brand with its own channels, SLAs, automations, and
       even its own Help Center portal for customers 1 2 . This means you can still distinguish each
       brand’s identity (unique portal URL, logo, knowledge base, etc.) while managing everything under
       one roof. It’s highly efficient: the founder (and future support teams) can switch between
       departments in one interface. Agents can be restricted to only see their department’s tickets,
       ensuring focus and privacy. All tickets can be tracked centrally by the founder, and cross-brand
       reports can be generated if needed. Cost-wise, this is included in one Zoho One subscription (no
       need for multiple accounts). Zoho confirms that multi-brand (multi-department) support is an
       intended use: “Multiple brands, multiple products. Each department can have its own ... separate self-
       service portals for each of them” 2 . This is the optimal approach for an organization with one
       parent company and several sub-brands.
       Cons: Slightly more initial configuration to set up departments and brand-specific settings. You must
       carefully configure roles and permissions so that, for example, an agent of Legacy Restore only
       accesses that department. However, Zoho Desk provides these controls out of the box. There is also
       a limit on number of departments/help centers (depending on your plan), but Zoho One (Enterprise)
       typically allows plenty – enough for all your sub-brands.


Recommendation: Use one centralized Zoho Desk with multiple departments, leveraging Zoho’s multi-
brand help center feature. This gives each sub-brand its own portal while keeping administration unified.
Customers of each brand will see a separate support site with the corresponding branding, which maintains
a professional appearance 1 . Internally, you avoid the chaos of juggling separate systems. In summary, a
single Zoho Desk instance with departments for “Legacy Platform” (general support) and each sub-brand
(Legacy Capital, Legacy Acadacom, Legacy Restore, etc.) is the most scalable and cost-effective solution. This
approach is fully supported by Zoho (each department gets a branded Help Center) and will make it easy to
add future brands—just create a new department for any new brand and configure its channels.


Using Subdomains vs. Separate Domains for Branding
Each sub-brand currently has its own domain name reserved (e.g. legacyrestore.co for Legacy Restore), but
the services are being served on subdomains of the main site (e.g. restore.legacyplatform.co). We need to
decide whether to continue with subdomains or switch to using the full separate domains for each brand’s
presence (website and support center). Here’s a comparison in terms of branding, scalability, and cost:


     • Subdomains under the Main Domain (e.g. restore.legacyplatform.co):
       Branding: Using subdomains keeps a visible connection to the main brand. This can be positive if you
       want customers to recognize the affiliation with Legacy Platform. However, it might dilute the
       individual brand identity of the sub-brand. Some customers might see restore.legacyplatform.co
       and not immediately recognize it as “Legacy Restore.” If Legacy Platform’s brand is not customer-
       facing (more of a holding company), this could confuse them.




                                                      2
       Scaling & Maintenance: Hosting everything under one domain might simplify certain technical
       management (one domain’s DNS to manage). In practice, though, each subdomain still requires its
       own configuration (SSL certificate, DNS records, content). Since Zoho One will require separate site
       setups for each subdomain anyway (explained in the Zoho Sites section), the operational effort is
       similar to using separate domains. From an SEO perspective, subdomains may or may not inherit the
       main domain’s SEO authority. Often, search engines treat subdomains as separate sites, so any SEO
       advantage is minimal.
       Cost: The cost difference between subdomains vs domains is minor. Subdomains are free to create
       under a domain you own, whereas separate domains incur registration fees. But domain fees are
       low (usually <$15/year each), so cost isn’t a major factor. If anything, using subdomains could save a
       few dollars on domains you might otherwise not use—yet in this case the domains are already
       purchased.


     • Separate Domains for Each Brand (e.g. legacyrestore.co):
       Branding: Using the sub-brand’s own domain presents it as an independent brand. This is often more
       professional for customer-facing services. Brand consistency is improved: customers of Legacy
       Restore will see legacyrestore.co in the URLs and email addresses, which reinforces that brand’s
       identity. There’s no visual cue of a parent company unless you choose to highlight it. If your strategy
       is to let sub-brands stand on their own (each with unique branding and perhaps even a slightly
       different audience), separate domains are beneficial. Many experts recommend separate domains
       for distinct business lines rather than subdomains, to avoid confusion and clearly segment the
       content 3 .
       Scaling & Maintenance: Separate domains allow you to grow or even spin off brands
       independently. Each site and support portal can be managed (or even transferred) without affecting
       the others. As you build teams, you could assign each team full ownership of their brand’s domain
       and site. The trade-off is that you’ll maintain multiple domains’ DNS settings, renewal, and SSL
       certificates. Fortunately, Zoho One and Zoho Sites make this straightforward (Zoho can automatically
       handle SSL for mapped domains). SEO-wise, each domain will build its own search engine presence.
       This requires effort per site, but it also means you can target niche keywords per brand without the
       main domain’s content intermixed.
       Cost: There’s a slight cost for each domain (but you’ve already acquired them). Zoho doesn’t charge
       extra for mapping multiple domains in one organization. Zoho Mail supports multiple domains on
       one account 4 , so you won’t need separate email providers. The main cost consideration is if
       separate domains forced separate software subscriptions – but with Zoho One, you can host
       multiple domains in one org easily. Therefore, using the distinct domains does not multiply your
       Zoho subscription cost.


Recommendation: Use each brand’s own domain for its customer-facing site and support portal. This
approach maximizes branding for each sub-brand and is aligned with long-term scalability. For example,
Legacy Restore should use legacyrestore.co for its website and have its helpdesk at a URL on that domain
(such as support.legacyrestore.co or help.legacyrestore.co). Legacy Capital would use legacycapital.co, and
so on. You can still maintain a connection to the main brand by mentioning “A Legacy Platform brand” on
those sites, but the URLs and emails will be brand-specific and professional.


If you prefer to keep the umbrella obvious, an alternative is to use distinct subdomains of the main site and
also map the separate domains to them. However, Zoho Sites will only allow one primary domain per site
  5   6 , meaning you’d ultimately choose either the subdomain or the separate domain as primary. Given




                                                     3
that, it’s cleaner to go with separate domains as primary for each brand. This will not significantly increase
maintenance work, since each brand will have its own site project regardless (see next section).


(One small note: Ensure each brand domain is verified in your Zoho One admin (Zoho Directory) for use in Zoho
Mail and Desk. Zoho One allows multiple domains in one organization for emails and services 4 , so you can add
legacyrestore.co, legacycapital.co, etc., alongside legacyplatform.co with no issue.)


Structuring Websites in Zoho Sites for Multiple Brands
Zoho Sites is used for the company’s websites. The founder currently uses legacyplatform.co on Zoho Sites,
and perhaps serves sub-brand content via subdomains. It’s important to understand Zoho Sites’ capabilities
to plan the site architecture:


     • Single Site vs Multiple Sites: Zoho Sites does not support “sub-sites” or distinct microsites
       within one website project. In other words, you cannot have one Zoho Sites project host
       completely different content at different subdomains – each Zoho Sites website has one set of
       pages/menus and one primary domain 5 . While you can add multiple domains to a single site,
       additional domains will just redirect to the primary domain’s content (they are aliases, not separate
       websites) 5 . Therefore, to have a truly separate website for each brand (with its own design or
       menu), you must create separate Zoho Sites websites – one for each brand’s domain 7 . This
       means, for example, a distinct site for legacyrestore.co and another for legacycapital.co, in addition
       to the main legacyplatform.co site. (Zoho’s community confirms this: to use three different domains for
       unique content, you need three different websites in Zoho Sites   7   .)


     • Current Setup (Subdomains on one site): If currently restore.legacyplatform.co is used, it might be
       implemented as either a separate Zoho Sites site mapped to that subdomain or just a section of the
       main site. If it’s a section on the main site, you likely had to compromise on design (the Legacy
       Platform menu and template extending into the Legacy Restore section). That is not ideal for
       branding or navigation. If it’s a separate site mapped to a subdomain, that is essentially the same
       effort as mapping it to its own domain.


     • Recommended Setup (Multiple Sites): Plan to create multiple Zoho Sites websites – one for each
       sub-brand, plus one for the main Legacy Platform umbrella. Under a Zoho One subscription, you are
       entitled to a certain number of sites (typically 1 site per 10 users, up to 5 sites on Enterprise) 8 . As
       the organization grows and you add users, you can increase the count of sites. In the meantime, if
       you are limited to 1 site with just the founder’s license, you might temporarily keep some content on
       the main site. But the end goal should be:


     • legacyplatform.co – the main corporate site (introducing the company and linking to sub-brands).
     • legacyrestore.co – its own Zoho Sites website for Legacy Restore (credit repair services).
     • legacycapital.co – its own site for Legacy Capital, etc.
     • legacyacadacom.co – its own site for Legacy Acadacom (education platform, presumably).

Each of these sites can have its own design language while still maintaining some common elements (you
might reuse the “Legacy” logo motif or color scheme to show they are related, if desired). The main site
(Legacy Platform) can act as a gateway: for example, it can have an overview of all divisions and links to
each sub-brand’s site. Conversely, each sub-brand site can have a note or footer “Part of the Legacy




                                                        4
Platform family” linking back to the main site. This structure is scalable; whenever a new sub-brand is
launched, you add a new domain and spin up a Zoho Sites instance for it.


One benefit of separate sites is also content specialization. Your SEO will be more targeted since each site
focuses on its specific services. Just ensure to cross-link appropriately (for example, Legacy Platform site
linking to Restore and vice versa for synergy).


In summary, use separate Zoho Sites for separate brands. If you are currently constrained by the number
of sites, you could initially use subdomains on the main site as a stopgap. But plan to migrate each sub-
brand to its own Zoho Sites project under its own domain as you enable more sites (which can come as you
add Zoho One user licenses or by purchasing additional site capacity). This will give each brand a dedicated
web presence, making it easier to manage content and design in the long run.


Organized Email Address Structure for Long-Term Scale
A core part of the architecture is a well-planned email address structure across all brands and
departments. By setting up the right email addresses now (even if they all ultimately reach the founder for
the time being), you create a scalable framework that can be handed off to others later. Below is a
breakdown of recommended email addresses by category:


     • Support Emails (Customer-Facing) for Each Brand: Every sub-brand should have its own support
       email   address    in  the    format    support@subbranddomain.       For   example,     use
       support@legacyrestore.co for Legacy Restore, support@legacycapital.co for Legacy Capital, etc.,
       and keep support@legacyplatform.co for general inquiries to the main brand (or corporate-level
       support). These addresses will be published to customers and will be the entry point for support
       tickets. In Zoho Desk, each address will correspond to a Department’s incoming channel. The best
       practice is to configure each support email to forward into Zoho Desk – typically, Zoho Desk
       provides a unique address (like support@<yourorg>.zohodesk.com or similar) for each department’s
       email channel. You would set up forwarding rules so that any email to support@legacyrestore.co
       creates a ticket in the “Legacy Restore” department of Zoho Desk. From the customer’s perspective,
       all communication still comes from support@legacyrestore.co (you can configure Zoho Desk to send
       replies using your domain by adding the domain and SPF/DKIM for Zoho Desk). This setup keeps
       customer support for each brand neatly separated. As the team grows, each support address can be
       managed by the respective brand’s support agents, but until then, the founder can monitor all
       departments. Zoho One allows multiple domains for email so you can have all these addresses in
       one system 4 . You won’t need a separate email inbox for each support address – just create them
       as aliases or distribution groups that feed into Desk. (For now, the founder might want copies of
       support emails – Zoho Desk can send notifications or you can add the founder as a CC via Desk rules
       – but over time this won’t be necessary when teams handle them.)


     • Departmental / Role-Based Emails (Internal Operations): Establish generic emails for key
       corporate functions using the main domain (legacyplatform.co). These would include addresses
       like legal@legacyplatform.co, finance@legacyplatform.co, hr@legacyplatform.co, possibly
       info@legacyplatform.co for general inquiries, and any other relevant roles (e.g. careers@ for job
       applications, pr@ for media, etc., as needed). The reason to use the main domain is that these
       departments likely serve the entire organization (all brands) rather than one brand. For example, you
       probably have one Legal team or one Finance department for all of Legacy Platform and its



                                                     5
 subsidiaries, not a separate legal team per sub-brand. Using one set of addresses avoids duplication
 and confusion (one person – currently the founder – is wearing all these hats, so one address per
 function is enough). In terms of handling these emails: initially, all of them can route to the founder’s
 inbox. The founder can set these up as email aliases on their Zoho Mail account or as group
 addresses that forward to them. For instance, emails to legal@legacyplatform.co could simply arrive
 in the founder’s mailbox (with perhaps a filter to tag it as Legal). When the time comes to delegate –
 say you hire a lawyer or a finance officer – you can assign these addresses to those new users or
 add the new users to the forward list. Because the email identity (legal@, finance@, etc.) is not
 person-specific, it provides continuity: external contacts will continue emailing legal@, and now the
 appropriate person (not necessarily the founder) will handle it, without needing to announce a new
 address. This approach also looks more professional and scaled even while you’re small; it prevents
 vendors or partners from always emailing the founder’s personal address for every little thing.


• Management and Technical Addresses: For management-level and technical communications, it’s
  wise to have dedicated addresses. The question already mentions founder.techops@... and dmarc-
  reports@..., so we’ll use those.


• founder.techops@legacyplatform.co: This address is used by the founder to handle technical
  operations correspondence. In practice, this likely means using it as the contact email for
  infrastructure services: e.g. registering this email with Twilio, Cloudflare, domain registrars, server
  providers, or even for admin accounts on Zoho, so that all technical alerts (server down alerts,
  security warnings, usage notifications) are sent here instead of to a general or personal email. The
  benefit of a separate techops email is that it isolates critical alerts from everyday email clutter.
  The founder can monitor this inbox (or have it forward to their main inbox but with clear tagging
  that it’s “Tech Ops”). Eventually, if a CTO or IT admin is brought in, the founder can hand off the
  founder.techops@ account or rename it accordingly, since it’s role-based. If founder.techops@ is
  currently an actual separate mailbox, you can keep it that way and just ensure you check it regularly
  (Zoho Mail allows you to be logged in to multiple mailboxes or you can set up a forward). If it’s just
  an alias, consider making it a separate mailbox if volume of alerts is high, so nothing is missed.
• dmarc-reports@legacyplatform.co: This address is intended for receiving DMARC reports. DMARC
  aggregate reports are XML reports sent by email from various ISPs to the address specified in your
  domains’ DMARC DNS records (often something like rua=mailto:dmarc-reports@legacyplatform.co).
  By setting up a dedicated dmarc-reports@ mailbox or alias, you ensure these reports (which can be
  numerous and verbose) don’t clutter your main inbox. All your domains can point their DMARC
  feedback to the same address – for example, in the DMARC record for legacyrestore.co and
  legacyplatform.co, etc., use dmarc-reports@legacyplatform.co as the aggregate report address.
  Then that one mailbox will gather all the data. You can then periodically review them or use a tool to
  parse them. This setup will help you monitor email deliverability and security for all your brand
  domains centrally.
• Other technical/management emails: You might also consider alerts@legacyplatform.co or
  it@legacyplatform.co as a general catch-all for automated alerts (similar to founder.techops).
  However, if founder.techops@ is already serving that purpose, it might be sufficient. Additionally, if
  the founder uses a personal mailbox like founder@legacyplatform.co or a first.last@ address, that
  can remain for direct human communication, while the role-based addresses handle functional
  areas. The key is to minimize the use of personal addresses for business functions – instead use
  these aliases so that later on you can transition responsibility smoothly.




                                                 6
Zoho Mail configuration: All the above addresses (support, department, techops, etc.) can be managed
within one Zoho Mail organization since you can add all domains to it 4 . You don’t need to purchase a
user license for every single address. Many can be set up as group addresses (distribution lists) that
forward to one or more real users. For example, you can make support@legacyrestore.co a group that
forwards to the Zoho Desk’s inbound address (external address) and to the founder (for now). Likewise,
hr@legacyplatform.co could be a group that currently includes only the founder as a recipient. Groups in
Zoho Mail are free and flexible. Alternatively, you can add email aliases to the founder’s account for some
addresses; e.g. founder’s mailbox could have an alias for legal@ so that those emails land in the same
mailbox. The drawback with aliases is when replying, you have to remember to send from the alias address
if you want the reply to come from legal@ rather than your personal address. Zoho Mail does allow send-as
alias, but you must be careful to select it. Groups might be cleaner for department addresses – the founder
can receive as part of the group but when replying, perhaps use the personal account (or configure the alias
sending). For support addresses, you generally won’t reply from Zoho Mail at all – replies will be sent from
Zoho Desk directly (and Desk will handle using the correct from-address).


Finally, ensure you set up proper SPF, DKIM, and DMARC for each domain in DNS. This means adding
Zoho’s SPF record (including Desk’s if using their mail servers for sending tickets) and generating DKIM for
Zoho Mail and Zoho Desk. This will improve deliverability for all these new email addresses and also ensure
your DMARC policy can be effective (the reports coming into dmarc-reports@ will show if anything fails).


Routing, Notifications, and Integrations for Efficient Operations
With the structure of departments and emails in place, the next step is to configure routing, tagging, and
integrations so that the founder (and future teams) can easily handle incoming messages without anything
slipping through. The goal is to automate the flow as much as possible and use Zoho One’s collaboration
tools to our advantage.


     • Zoho Desk Email Routing: Set up each support email to forward into the correct Zoho Desk
       department as mentioned. Once this is in place, all customer support inquiries will be converted
       into tickets. This is a huge advantage because it takes the burden off the founder’s inbox – you
       won’t have to manually track which customer asked what across multiple email threads. Instead,
       tickets will be organized by brand (department), with distinct queues. Within Zoho Desk, you can
       create assignment rules if needed (for now, maybe all tickets get assigned to the founder as the
       sole agent; later, you can auto-assign to agents of the specific department). You can also set up
       categorization and tags in Zoho Desk for common issues, which will help when you bring on new
       agents (they can quickly see how past issues were tagged/resolved). The Desk system will send out
       auto-acknowledgement emails to customers and can enforce SLAs (which you might define
       differently for each brand, depending on your service promises).


     • Zoho Desk – Zoho Cliq Integration: Since the founder uses Zoho Cliq heavily for communication,
       integrating Desk with Cliq will streamline notifications. Zoho Desk offers a built-in integration (via
       Deskbot) that can post ticket events to Cliq channels 9 . You should sync each support department
       to a dedicated channel in Cliq. For example, create a channel #support-legacyrestore in Cliq
       and map the Legacy Restore department to that channel. Do the same for other departments
       (perhaps #support-legacycapital , etc.). With this integration active, whenever a new ticket
       comes in, a customer replies, or a ticket is escalated, the Deskbot will post a message in the
       corresponding channel 9 10 . This is extremely useful because it provides instant, centralized



                                                     7
 notifications without relying on email. The founder can sit in Cliq and get real-time updates for all
 brands: e.g. “New ticket #1234 from [Customer] in Legacy Restore” appears in #support-
 legacyrestore channel. Since currently the founder is the only agent, they could also choose to have
 Deskbot send direct messages, but channel notifications are better for future scalability (so the
 whole team for that brand can see it). The Deskbot messages even allow some quick actions – you
 can click to view the ticket in Desk, or even reply/assign from Cliq in some cases 11 . This integration
 ensures that no support request goes unnoticed, even if the founder isn’t actively checking the
 Zoho Desk interface or email. It’s essentially a helpdesk command center in your chat app.


• Zoho Mail Filtering and Tagging: For emails that are not going through Zoho Desk – such as the
  departmental addresses and technical alerts – use filters and labels in Zoho Mail to organize them.
  The founder’s mailbox (or the relevant mailbox) can have rules like: If an email is sent to
  legal@legacyplatform.co, automatically tag it with a “Legal” label and perhaps move it to a Legal
  folder (or leave it in inbox but colored). Similarly, anything to finance@ gets tagged “Finance”. This
  way, even if all these come to one person, it’s immediately clear what category they belong to. The
  founder can then prioritize or delegate accordingly. Another example: filter Cloudflare alerts (which
  might have a subject containing “Cloudflare Alert” or come from a specific email like no-
  reply@cloudflare.com) and tag them as “Cloudflare” or put them in an “IT Alerts” folder. Do the same
  for Twilio (e.g. filter messages from Twilio or containing “Twilio” in subject). You might funnel all such
  tech notifications into one folder, or separate by service – depending on volume and importance.
  During the initial phase, the founder will receive all these emails, but thanks to filtering, they won’t
  appear as one big mixed pile; you can address them in batches (e.g. check the “HR” label for any job
  inquiries, check “TechOps” label for system alerts daily, etc.).


• Email Forwarding Between Addresses: In some cases, you might choose to forward certain emails
  to others or to group channels. For instance, you could set up an email-to-cliq integration or simply
  use the fact that Cliq channels can have email addresses (Cliq has a feature to generate an email
  address for a channel, so incoming emails post as messages). If you wanted, say, all critical alerts to
  ping a Cliq channel (like #alerts or your direct message), you could forward those specific emails
  (maybe those from monitoring@twilio.com) to the channel’s email. This is an advanced step and
  optional – with the Desk integration, most critical customer issues are covered in Cliq already. For
  now, straightforward forwarding like: all dmarc-reports@ emails forward to founder.techops@ so
  they sit in one mailbox, or all hr@ emails forward to founder@ – can be set as needed. Just be
  mindful of not creating forward loops. Typically, use single-point collection: e.g., have
  founder.techops@ collect certain things and then just check that, rather than forwarding it multiple
  times.


• Zoho Cliq for Team Collaboration: In addition to automated notifications, plan to use Cliq for
  internal discussions about support tickets and operations. You might have a channel for each
  department’s team (even if the “team” is just you now). As you hire support agents for Legacy
  Restore, they’ll join the #support-legacyrestore channel, where they see tickets via Deskbot and can
  discuss how to resolve them with each other or ask for help from the founder. Similarly, you could
  have an “Announcements” or “Ops” channel where system alerts or important emails (forwarded as
  mentioned) show up, and the tech team (for now just you) can track them. This keeps
  communication contextual. It’s far better than mixing support threads via email CC — everything is
  logged in Desk and mirrored in chat for quick attention 9 .




                                                  8
     • Using Zoho One Integrations: Beyond Desk and Mail, consider other Zoho One apps that might
       help. For example, Zoho CRM could be integrated with Desk so that when a ticket comes in, the
       customer’s CRM record is linked (if you’re using CRM for clients). Zoho People or Zoho Recruit could
       be used for HR, which might involve HR@ emails. But those are future enhancements. The key
       integrations to focus on now are Mail ↔ Desk (via forward) and Desk ↔ Cliq. Additionally, ensure that
       all your domains and users are managed in Zoho Directory/Admin Console – since you’ll have
       multiple domains, you want them under one organization (which they will be if you add them in
       Zoho One’s domain management).


     • Notification Summary: In the interim (single-founder) stage, you might still rely on some email
       notifications for awareness – e.g., Zoho Desk can send the agent an email for new tickets. However,
       since you are integrating with Cliq, you could turn off some email notifications to reduce redundancy
       (no need to get 3 alerts for one ticket). The combination of Desk’s ticket queues + Cliq alerts + Mail
       filters will ensure you have multiple layers of capture: if a customer email somehow doesn’t
       convert to a ticket, you’ll see it in a mailbox; if a ticket comes in, you’ll see it in Desk and Cliq; if a
       system sends an alert, it’s flagged in email (and possibly in Cliq if you set that up). This multi-channel
       notification means even as volume grows, you’ll catch issues promptly.


     • Future Delegation: As teams come on board, you will start removing the founder from being the
       single point of contact and use the structure to delegate. For example, when a support agent for
       Legacy Restore is hired, you add them to Zoho Desk (just in the Legacy Restore department). They
       will start receiving the Cliq notifications in #support-legacyrestore and can act on tickets. The
       founder can then step back from day-to-day handling of that brand (maybe staying in the loop via
       the channel or reports). Similarly, when an HR person is hired, you give them access to the hr@
       mailbox or change the forward to go to them. Because you’ve set up standard emails and
       centralized systems, handing off is as easy as adding a new user to a group or department, rather
       than changing your whole process. This organized architecture thus scales smoothly from one
       person to many.


Summary: Recommended Setup Per Brand and Department
Below is a summary table of the proposed email and support desk architecture for each brand, including
the main brand (Legacy Platform) and sub-brands. It outlines domains, how support is handled, key email
addresses, and how notifications flow for each:


Brand / Role | Domain & Website | Zoho Desk Setup (Support) | Key Email Addresses | Notification
Flow -|- Legacy Platform (Main Brand & Corporate) | legacyplatform.co – Main corporate website (Zoho
Sites). Serves as umbrella site with links to sub-brands. | Department: “Legacy Platform – General”. Help
Center on legacyplatform.co (e.g. support.legacyplatform.co). Used for general inquiries or any shared
support issues. | Support: support@legacyplatform.co (general support)<br>Departments: legal@,
finance@,           hr@legacyplatform.co          (for       company          use)<br>Management/Tech:
founder.techops@legacyplatform.co; dmarc-reports@legacyplatform.co (collects DMARC for all domains) |
Support Tickets: Go into Desk “Legacy Platform” dept; founder notified via Deskbot in Cliq #desk (or
#support-legacyplatform) channel. 9 <br>Departmental Emails: All route to founder (for now) – filtered
and tagged (e.g., “HR”, “Legal” labels) in Zoho Mail for easy reference.<br>Alerts: Twilio, Cloudflare, Zoho
system alerts sent to founder.techops@ (separated from main inbox). Founder monitors this mailbox and
gets critical alerts on Cliq if configured. Legacy Restore (Credit Repair sub-brand) | legacyrestore.co –




                                                       9
Dedicated Legacy Restore website (Zoho Sites). Currently also accessible via restore.legacyplatform.co, but
will use primary domain going forward. | Department: “Legacy Restore Support” under main Zoho Desk.
Has its own Help Center (e.g. support.legacyrestore.co mapped via Zoho Desk multi-brand 12 ). Agents
(currently founder) assigned to this dept only see Restore tickets. | Support: support@legacyrestore.co →
(forwarded to Zoho Desk Dept “Legacy Restore”)<br>Other: (Use main corporate addresses for legal/
finance unless a unique need arises for Restore-specific emails.) | Support Tickets: Automatically
converted to tickets in Legacy Restore dept. Deskbot posts new tickets and replies in #support-
legacyrestore Cliq channel 9 . Founder monitors channel and Desk; in future, Legacy Restore support
agents will handle directly via this channel and Desk.<br>Other Communications: Any Restore-related
business emails (if using addresses like info@legacyrestore.co) can forward to founder or relevant
department address at corporate. Technical alerts about the Restore site (uptime, etc.) go to
founder.techops (as part of overall alerts). Legacy Capital (Financial services sub-brand) | legacycapital.co
– Dedicated Legacy Capital website (Zoho Sites). To be set up with its own domain and site. | Department:
“Legacy Capital Support” in Zoho Desk. Own Help Center portal for Legacy Capital customers. | Support:
support@legacycapital.co → Desk Dept “Legacy Capital”<br>Other: (Likely uses corporate legal/finance
emails. Could create sales@legacycapital.co if needed for client inquiries, which forwards to sales or
founder.) | Support Tickets: Go to Desk “Legacy Capital” dept; notifications appear in #support-
legacycapital channel for real-time alert. Founder (and future Capital support team) follow this channel for
updates.<br>Other Emails: If addresses like sales@ or info@legacycapital.co exist, they forward to the
appropriate person (e.g. founder or a future sales team) and are labeled “LegacyCapital-Info” via mail filters.
No separate tech alerts (covered by central techops). Legacy Acadacom (Education platform sub-brand) |
legacyacadacom.co – Dedicated Legacy Acadacom website (Zoho Sites). Uses its domain for web and
support. | Department: “Legacy Acadacom Support” in Zoho Desk. Help Center for Acadacom brand (with
its own knowledge base, community, etc., as needed). | Support: support@legacyacadacom.co → Desk
Dept “Legacy Acadacom”<br>Other: (Likely uses corporate addresses for common departments; could have
academic-specific ones like faculty@legacyacadacom if applicable, managed similarly.) | Support Tickets:
Routed into Desk dept and announced in #support-acadacom Cliq channel (agents to be assigned later;
founder watches meanwhile).<br>Other Emails: Anything like course inquiries to acadacom can be filtered
to an “Acadacom” label in founder’s mail. Future team members can be added to those distribution lists
when hired.


Notes: All domains (legacyplatform.co and sub-brand domains) will be added to the single Zoho One
organization, allowing creation of these addresses within one admin console. The founder currently acts as
the point of contact for all emails and tickets – but this architecture is built so that adding new team
members is straightforward. New users can be given access only to the relevant email accounts or Zoho
Desk departments. The use of standardized addresses (support@, legal@, etc.) means communication
channels remain consistent even as personnel change. By using Zoho Desk for customer support and
integrating it with Zoho Cliq, the founder and future agents get instantaneous, organized notifications of all
support activity 9 . The combination of separate websites, separate support emails, and centralized
management tools ensures each brand delivers a focused experience to customers while the founder can
oversee everything from one place.




                                                      10
Final Recommendations and Next Steps
Adopting the above strategy will yield a scalable and organized email and support desk architecture for
Legacy Platform and its sub-brands. To recap the key recommendations:


     • Zoho One Structure: Keep one Zoho One organization for all brands (one parent company). Within
       this org, utilize Zoho Desk’s multi-department (multi-brand) support to give each brand its own
       helpdesk portal 1 . This avoids duplicate subscriptions and enables centralized oversight. All sub-
       brand support desks are contained in one system, under Legacy Platform’s umbrella.


     • Domains and Branding: Leverage the individual brand domains you own for customer-facing sites
       and addresses. Map each Zoho Desk help center to the respective domain or subdomain (Zoho Desk
       lets you set a unique help center URL per brand 12 ). Using separate domains will strengthen each
       brand’s identity and allow independent growth. Continue to use legacyplatform.co for corporate
       identity and shared services.


     • Zoho Sites Implementation: Plan for separate Zoho Sites websites for each domain/brand. This is
       necessary given Zoho Sites’ one-domain-per-site design 7 . Design each site to fit the brand, and
       link them for easy navigation between brands. In the interim, if site license limits are an issue, use
       the main site with clearly segmented pages or subdomains, but migrate to separate sites as soon as
       feasible to avoid design and content limitations.


     • Email Architecture: Establish a comprehensive set of role-based email addresses:


     • Support: support@ each brand (goes to Desk) – customers get replies from these addresses via Zoho
       Desk.
     • Internal Departments: Standard emails like hr@, legal@, finance@ (on the main domain) for internal
       and external communication in those areas. Use them in all official documents and communications
       now (even if they reach only the founder), so that stakeholders get used to them.
     • Management/Technical: Use founder.techops@ for registering all technical services and receiving
       operational alerts, and dmarc-reports@ for email security feedback. This segregation will keep the
       founder’s primary inbox cleaner and ensure critical notices are seen 13 .

Organize these addresses in Zoho Mail via aliases or groups as appropriate. Start with forwards to the
founder, and later add or switch recipients as teams form.


     • Integration & Automation: Implement the Zoho Desk ↔ Cliq integration (Deskbot) to channel
      support notifications into chat in real-time 9 . This will allow the founder and any team members to
      react quickly to client needs without constantly checking email. Utilize Zoho Mail rules to sort and
      flag incoming mail by department or brand, and set up any needed forwards (for example, copying
      certain alerts to multiple recipients). The idea is to minimize manual email triage – let rules and Desk
      automation do that work. All support communications should be handled in Zoho Desk (with the
      email-to-ticket flow), which brings consistency and tracking.


     • Knowledge Base and Self-Service: As a complementary step (once the support structure is in place),
       build out a knowledge base for each brand within Zoho Desk. With multi-brand, each department
       can have its own help center with FAQs and articles 14 . This will reduce repetitive questions and




                                                     11
          empower customers. It also scales your support without additional load on agents. Similarly,
          configure any community forums per brand if you want user communities, since Zoho Desk can host
          those separately per help center.


         • Monitoring and Adjustment: After implementing, monitor how the system is working. Check that
           tickets are flowing to correct departments (you can simulate by sending test emails to each support
           address). Ensure email SPF/DKIM are passing (via your dmarc-reports@ data). Fine-tune the Cliq
           notifications (you might adjust which events trigger a message – initially you might want
           notifications for new tickets and customer replies, but not for every status change, to avoid noise).
           Also, gather feedback from customers interacting with the new support emails and portals – do they
           find it professional and responsive? Use that feedback to improve response templates or portal
           design per brand.


By following this strategy, the founder will create a robust support and communication infrastructure
that appears as a much larger, well-coordinated organization to customers and partners. Each sub-brand
will have a professional standalone presence (own site, own support email and portal), instilling trust and
clarity. Internally, however, everything feeds into centralized tools that the founder can manage efficiently.
As teams are added, members can be slotted into this architecture with minimal friction, inheriting the
organized setup. The result is a scalable system where growth doesn’t lead to chaos but is accommodated
by the existing structure.


Overall, this multi-brand architecture using Zoho One will help Legacy Platform scale its empire of sub-
brands while maintaining high-quality customer support and streamlined internal operations. Implement
the changes step by step, and you’ll be well-prepared for the company’s expansion, with each brand
reinforcing the strength of the Legacy Platform family. Good luck with the setup!



 1   12    14   Providing Multi-brand Customer Support - Zoho Desk Knowledgebase
https://help.zoho.com/portal/en/kb/desk/support-channels/help-center/articles/providing-multi-brand-customer-support

 2   Multi-department Help Desk | Zoho Desk
https://www.zoho.com/desk/multidepartment.html

 3   subdomain vs separate domain | Warrior Forum - The #1 Digital ...
https://www.warriorforum.com/search-engine-optimization/1042675-subdomain-vs-separate-domain.html

 4   How to manage multiple domains with one Zoho account?
https://help.zoho.com/portal/en/community/topic/how-to-manage-multiple-domains-with-one-zoho-account

 5   6    Domain Options
https://help.zoho.com/portal/en/kb/zohosites/help-guide/configuration/domain-ssl/articles/domain-options

 7   One Site with one main domain and many subdomains - Zoho Cares
https://help.zoho.com/portal/en/community/topic/one-site-with-one-main-domain-and-many-subdomains

 8   Zoho One - How to Add a Domain to your Account. - YouTube
https://www.youtube.com/watch?v=qlfWR-p8nXo

 9   Cliq Integration with Zoho Desk
https://www.zoho.com/cliq/zohodesk-integration.html




                                                             12
10   11   Integrating Zoho Desk with Zoho Cliq - Online Help | Zoho Desk
https://help.zoho.com/portal/en/kb/desk/integrations-and-marketplace/collaboration/articles/zoho-cliq-integration

13   One ZohoMail Account for TWO Domains? - Zoho Cares
https://help.zoho.com/portal/en/community/topic/one-zohomail-account-for-two-domains




                                                              13

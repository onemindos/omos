---
id: ui-15-20260111053926
title: "Implementing A Unified Multi Brand Support System In Zoho Desk"
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
ai_summary: "Implementing A Unified Multi Brand Support System In Zoho Desk"
---

Implementing a Unified Multi-Brand Support
System in Zoho Desk
Single Instance Branding (Multi-Brand Setup)
Use Zoho Desk’s Multi-Brand Help Center: Zoho Desk Enterprise edition allows you to support multiple
brands in one help desk by using departments for each brand 1 . Enable the Multi-branding feature
under Setup > Organization > Rebranding > Multi-branding, which lets you create a separate customer
Help Center for each department (brand) 2 3 . This approach gives every brand its own portal with
unique logo, colors, domain, and content, all within a single Zoho Desk instance 4 . Crucially, enable the
Common User Base model so that customers can use one set of login credentials to access all your brand
portals 5 (they won’t see other brands’ portals unless you tell them, but they won’t need separate
accounts).


Step-by-Step Configuration: 1. Create Departments for Each Brand: In Zoho Desk, add a new
department for each brand (Setup > Organization > Departments > New Department). Give it a name,
help center display name, and logo for that brand 6 . Assign the appropriate support agents to that
department 7 and configure its channels (email address, social media, etc.) for that brand’s
communications. For example, you might have “Brand A Support” and “Brand B Support” as two
departments, each with their own agents and inbox (like support@brandA.com, support@brandB.com). 2.
Enable Multi-Brand Help Centers: Go to Setup > Organization > Rebranding > Multi-branding and click
Enable Multi-branding 8 . This turns on the common user base multi-brand support by default 9 . (Note:
Decide upfront between common vs. unique user base; switching later is not recommended 10 . Common
user base is ideal for a unified login experience.) 3. Add Branded Help Centers: After enabling, create a
Help Center for each new department via Setup > Channels > Help Center > New Help Center 11 . For
each help center, upload the brand’s logo and favicon, set a name (e.g. “<BrandName> Support Portal”), and
map a custom domain (such as support.brandA.com) for that portal 12 . Configure the portal’s theme to
match the brand’s identity (colors, banner, etc.) so customers get a seamless branded experience. 4.
Customize Portal Settings: For each brand’s portal, adjust settings like visibility (public or login-required),
community forums, and knowledge base access as needed. Each Help Center is essentially isolated to its
department’s data, allowing you to set different SLAs, notification emails, and article access permissions per
brand 4 . Ensure email templates and automated messages in each department are branded (you can
customize the “From” address and email signature for each department’s replies to include the brand name/
logo). This way, whether a customer submits a ticket via email or the portal, they receive responses aligned
with that brand.


By following these steps, you maintain one Zoho Desk instance and a single agent login console, while
each brand gets its own customer-facing support portal. Agents can switch between departments
internally, but customers see only the branded help center they need. This setup is scalable – when you add
a new brand in the future, simply create a new department and help center for it, without spinning up an
entirely separate Desk account.




                                                      1
Knowledge Base Optimization Across Brands
Understand Knowledge Base Isolation: In a multi-brand configuration, each department has its own
Knowledge Base categories and articles, visible only on that department’s Help Center 13 . Zoho Desk does
not currently allow a single knowledge base article or category to be shared across multiple brand portals –
categories cannot be linked to more than one department once multi-brand is enabled 14 . This means you
will need to plan how to manage common content to avoid tedious duplication.


Strategies to Minimize Duplication: - Centralize and Replicate Key Content: Identify articles that are
relevant to all (or multiple) brands, such as FAQs or how-tos that apply across your product lines. Author
these in one place first (for consistency) and then copy them into each brand’s knowledge base category
that needs them. While this still means maintaining separate copies, you can at least reuse the writing.
Keep a master document or internal repository of universal solutions so updates can be made in one
reference and then propagated to each portal. - Use a “General” Category (if applicable): If one of your
brands (or the default department) serves as a general support portal, you could publish common articles
there and simply link to them from other brand sites. Because all help centers share the common user base,
a customer who has access to one can be given access to another if needed. However, customers typically
won’t navigate between brands on their own 5 , so providing direct links to a central FAQ site might be an
option for internal use or power users. In most cases, it’s better to duplicate content for a truly seamless
experience (so customers stay on their familiar brand site). - Leverage Templates and Style Consistency:
Zoho Desk’s article editor allows you to save templates. Create a standardized template for articles (with
your preferred formatting, disclaimers, and style) that you can use across all brands. This ensures that even
if you maintain multiple copies, the look and structure of the information remain consistent, reinforcing
professionalism and a uniform quality standard. - Assign Ownership for Updates: To keep maintenance
low, designate an “knowledge base manager” or team responsible for updating common articles in all
departments when something changes. For example, if a troubleshooting step changes, update the master
copy and quickly reflect that change in each brand’s portal. Using the Zoho Desk API or an extension can
also speed up copying articles between departments if you have technical resources – you could write a
script to publish an article to multiple departments to reduce manual effort (since direct sharing isn’t built-
in 14 ). - Brand-Specific Content: Continue to keep brand-specific articles separate. Each Help Center can
have categories specific to that brand’s products or policies. Organize the knowledge base per brand in a
way that makes sense to that brand’s customers. For instance, Brand A might have categories by Product
Line, whereas Brand B (a different business) might organize articles by Topic. The multi-brand setup lets you
tailor the structure per brand. Just ensure your team follows a consistent process for creating and
approving articles so the quality remains high across all knowledge bases.


By thoughtfully centralizing the creation of common answers and using the above practices, you can
manage multiple knowledge bases with minimal duplication. The key is to plan your content architecture:
decide which articles belong in all brands vs. one, and set up a workflow to keep them in sync. Over time,
analyze which topics are repeatedly duplicated – those might be good candidates for consolidation or
maybe moving to a single public FAQ page if appropriate. Even though Zoho Desk doesn’t automate cross-
portal KB sharing currently, with a solid process the overhead can be kept low while ensuring each brand’s
customers have the answers they need in their own portal.




                                                      2
Unified Live Chat with Zoho SalesIQ
Single Chat Platform, Multiple Brands: Zoho SalesIQ (live chat) can be configured to handle chats from all
your brands through one SalesIQ instance, thanks to its Brands and Departments features. Instead of
setting up separate chat systems for each website, you will create multiple “Brand” profiles within SalesIQ –
typically one for each website or product brand you support 15 . This allows you to use one SalesIQ code
snippet across all sites (each brand gets a unique snippet or widget), while managing chats centrally. The
benefit is twofold: your operators use a single interface to handle all incoming chats, and you can route
chats intelligently based on brand or issue type.


Best Practices for SalesIQ Multi-Brand Setup: 1. Define Brands in SalesIQ: In your SalesIQ dashboard,
go to Settings > Brands and add a Brand for each website or product line you support 15 . For example,
create “Brand A Chat” and “Brand B Chat.” Each SalesIQ Brand can have its own chat widget customization,
proactive chat triggers, and configurations, which you should tailor to match the look and tone of the
corresponding website. Adjust the widget’s colors, avatar, welcome message, and language to be consistent
with each brand’s identity (so customers feel they are chatting with Brand A or B, not a generic service). 2.
Associate Departments (Routing): For each SalesIQ Brand, map it to the correct support Departments
(teams) from Zoho Desk 16 17 . This is done under Settings > Brands > (Select Brand) > Flow Control >
Departments Responsible for the Chat. For example, assign the “Brand A Support” department to the
Brand A chat widget, and “Brand B Support” to Brand B’s widget. This means when a visitor initiates a chat
on Brand A’s site, only agents from the Brand A Support department will receive it 18 17 . If needed, you
can assign multiple departments to one brand (e.g., route chats to different teams based on region or
query type), but ensure there’s logic for which department picks up which chat. 3. Intelligent Chat
Routing: Decide how customers will be routed to the right department or self-service: - By Website: The
simplest method is implicit routing: the chat widget on Brand A’s site is tied only to Brand A’s support
department, and similarly for Brand B. In this case, the customer doesn’t need to choose anything — the
system knows which brand’s site they came from and routes the chat to the appropriate team automatically.
- By Customer Selection: If you offer multiple support options on the same site, you can add a pre-chat
question or menu. For instance, on a unified support page or if Brand A has multiple divisions, use SalesIQ’s
pre-chat forms to ask the visitor what they need help with. This could be a dropdown for “Which product
do you need help with?” or “Select your region/brand,” which corresponds to departments. SalesIQ’s Flow
Control allows designing these forms or using a chatbot to ask upfront questions. Based on the answer,
SalesIQ can assign the chat to the matching department. - Self-Service Deflection: Enable knowledge base
search in the chat widget or use Zia Answer Bot for common queries. Zoho Desk’s Zia can integrate with
SalesIQ (via the ASAP widget or Zia Chatbot) to suggest help articles to users before connecting to an agent
  19 . This means when a user opens the chat, they might first see, “Hi! How can I help you?” and if they type

a question, Zia will pull relevant knowledge base articles from the corresponding brand’s KB. If the answer is
found, the user might resolve their issue without live agent intervention; if not, the chat seamlessly
escalates to a human. This provides a unified chat experience with AI-driven self-service. 4. Configure
Business Hours & Offline Handling: Set up business hours for each SalesIQ Brand (e.g., if Brand A support
is 9-5 PST and Brand B is 9-5 EST, configure those separately) 20 . In SalesIQ, you can define operating
hours per brand and show an offline message or form outside those times. Make sure the offline form is
linked to Zoho Desk so that after-hours inquiries become tickets in the correct department. This keeps the
pipeline unified – chats turn into tickets when agents are unavailable, maintaining continuity. 5. Test and
Refine: After setup, test the chat on each website. Initiate a chat as a customer on Brand A’s site and ensure
it notifies the Brand A support agents. Check that the branding (chat window labels, department names)
presented to the customer is appropriate (you might not want the customer to see internal department




                                                      3
names, so configure friendly names if possible). Likewise, test any self-service prompts. Make sure that the
transition from self-service to live chat is smooth and that transcripts are getting attached to tickets or CRM
as needed. 6. Operator Training: Train your chat agents to recognize and handle multiple brands. Zoho
SalesIQ will show them which brand/site the incoming chat is from, and possibly the department. Agents
should follow the correct script and signature for the brand they’re serving in that chat. This maintains the
professional brand identity in live conversations.


With this unified SalesIQ setup, you have one live chat system that feels native on all your brand websites.
It routes customers behind the scenes to the right team based on where they came from or what they
select. You’ll be able to monitor all chats in one place, while keeping the experience personalized for each
brand. Over time, use SalesIQ’s analytics per brand to see peak chat times and volumes, and adjust staffing
or bot prompts accordingly for each department.


Structuring Teams, Departments & Workflows
A well-organized internal structure in Zoho Desk will ensure each brand’s tickets are handled efficiently
while remaining scalable. Here’s how to set up your teams, departments, and workflows for brand-
specific support:


     • Departments per Brand: As established, each brand is a separate department in Zoho Desk. This
       segregation is the foundation of your workflows. It lets you manage each brand’s tickets
       independently – you can have different SLAs, escalation rules, and canned responses for each
       department 4 21 . For example, Brand A’s department might have a 24-hour first response SLA
       and a custom set of statuses, while Brand B’s department could have a 12-hour SLA and its own
       workflow. Zoho Desk allows each department to have its own agents, mailboxes, help center,
       community forum, social channels, business hours, and automation rules 21 , so take
       advantage of that. Configure each department’s settings under Setup > Channels (for channel-
       specific config like email templates, chat, telephony) and Setup > Automation (to create rules
       scoped to that department).
     • Agent Assignment & Teams: Assign agents to the departments that correspond to the brands they
       will support 7 . Agents will only see tickets for the departments they belong to, which keeps Brand
       A’s issues only visible to Brand A’s support team, and so on. If some agents work across multiple
       brands, you can add them to multiple departments – they will be able to toggle and support both,
       but ensure they are trained for both brands. Use Teams in Zoho Desk to group agents by function or
       expertise across or within departments. For instance, within Brand A’s support department you
       might create teams like “Brand A – Tier 1” and “Brand A – Escalations,” or a “Global VIP Support” team
       that includes senior agents from all brands. Teams are a flexible way to categorize agents (e.g., by
       skill or region) and can be used in ticket assignments (round-robin assignment can be done by
       team). Keep in mind that teams don’t override department permissions – an agent in a team still
       needs access to the ticket’s department – but they help organize who should handle what. This
       layered approach (Departments for brand segregation, Teams for skill grouping) gives you flexibility
       as you grow.
     • Workflow and Automation Rules: Set up workflows that are department-specific to enforce each
       brand’s processes. In Setup > Automation, create rules for each department. For example:
     • Assignment Rules: When a ticket comes into Department A, assign it to the Brand A Tier 1 team
       queue, but if the query is about a certain product line (you can use criteria like ticket category), route
       it to a specific agent or sub-team. Each department can have its own assignment logic. If you have




                                                       4
       identified fields like “Product” or “Issue Type” on tickets, you can create mappings per brand (e.g.,
       Brand B’s tickets with “Billing” issue go directly to Brand B’s billing specialist).
     • Service Level Agreements (SLAs): Define SLAs per department reflecting each brand’s promise to
       customers. Perhaps Brand X offers premium support with faster response times – you can configure
       a stricter SLA policy and escalation ladder under that department. Zoho Desk will then track and
       escalate tickets breaching SLA for Brand X separately from others.
     • Notification and Response Templates: Customize the email templates and notifications for each
       department. Go to Setup > Automation > Customer Notifications (and Agent Notifications). Make
       sure the templates (new ticket auto-reply, ticket closed notification, etc.) use the correct branding
       (like the department’s from-address and signature). This ensures a customer of Brand A gets an
       email with Brand A’s name/logo and not something generic. You can even use different styles or tone
       of voice if the brands require it.
     • Macros and Custom Functions: If you have common actions that differ by brand, set them up in the
       context of each department. For instance, a macro to gather more info might send a different set of
       questions depending on the product line (housed in that department). Zoho Desk’s Blueprint
       (process flow) can also be configured per department if you want to enforce multi-step ticket
       lifecycles that are brand-specific.
     • Scalability & Future Brands: To ensure your structure is scalable, document the configuration for
       one brand as a template. When adding a new brand (department) later, you can replicate the setup
       with modifications rather than reinventing it. For example, if you know each department needs
       certain categories, statuses, and email templates, create a checklist for new department onboarding.
       Zoho Desk allows cloning some settings (you might manually copy workflow rules). Keep global
       settings minimal; prefer departmental settings so that adding a new department is mostly self-
       contained. Also, use consistent naming conventions for fields and statuses across departments
       where possible – this helps when doing cross-department reporting. However, if brands have very
       different processes, it’s okay for them to have unique workflows. Zoho Desk’s multi-department
       capability is designed to handle multi-brand or multi-division setups in a modular way.
     • Monitoring and Quality Control: As an admin or support manager, regularly review each
       department’s performance. Although agents focus on their respective brand tickets, you (and other
       admins) can see across all departments. Make use of Dashboards and Reports filtered by
       department to ensure each brand’s support metrics (response times, backlog, CSAT, etc.) meet your
       standards. This structure also allows you to identify if one brand’s support is lagging and needs more
       resources or training, without it being obscured by aggregate data.

In summary, structuring Zoho Desk by brand (departments) gives you clean separation with the ability to
customize workflows to each brand’s needs. By leveraging teams, you ensure the right people handle the
right tickets, and by setting up robust department-specific automation, you create a low-maintenance
system. As new brands join, you simply plug in a new department with its configurations – the unified
framework remains intact.




                                                     5
Zoho Integrations & Real-Time Insights (CRM, Campaigns, Zoho
One, Radar)
One of the advantages of using Zoho Desk in a Zoho One ecosystem is the native integrations available.
Setting these up will unify your support data with sales, marketing, and other functions, providing a truly
seamless experience across all brands.


     • Zoho CRM Integration: Integrate Zoho Desk with Zoho CRM to maintain a single customer
       database and get 360° view of customer interactions 22 23 . With the two-way sync enabled, all
       your Zoho Desk contacts (from all brands) can sync to CRM and vice versa, so there’s one master
       record per customer 23 . Sales and support teams will always see up-to-date information. Agents,
       directly from a support ticket, can view CRM details like the account name, deals (potentials), sales
       notes, or any custom CRM field that provides context 24 . For example, if a VIP customer from Brand
       A contacts support, the agent might see in the Desk sidebar that this customer has a high-value deal
       in CRM or a premium subscription – allowing them to tailor the service. This integration also enables
       “Handshake Mode” where you can prioritize tickets based on CRM data (like customer tier).
       Implementation: Go to Setup > Integrations > Zoho CRM in Desk and follow the steps to authorize
       the sync. Map fields such as Account Name, Contact info, and ensure “Department” or “Brand” is
       captured as well (you might add a field in CRM to tag which brand a ticket came from, using Desk’s
       department info). The CRM integration is especially powerful in a multi-brand scenario because it
       breaks down silos – all customer support history lives in CRM alongside sales history, which is great
       for account management.
     • Zoho Campaigns & Customer Communications: While there isn’t a direct one-click integration
       between Zoho Desk and Zoho Campaigns like there is with CRM, you can still leverage these
       together. A common approach is to sync contacts from Zoho Desk (or CRM) into Zoho Campaigns
       lists. For instance, you might maintain separate mailing lists in Zoho Campaigns for customers of
       each brand (using the brand field or department to segment them). This way, your marketing team
       can send newsletters or product updates that are brand-specific. You can automate adding a Desk
       contact to a Campaigns list via Zoho Flow or by syncing Campaigns with CRM (which is fed by Desk).
       Additionally, use Campaigns to send satisfaction surveys or follow-ups: e.g., after a ticket is closed,
       you could trigger a Zoho Campaigns email sequence to that customer for feedback or upselling,
       tailored with the appropriate branding. The key is to ensure data flows from Desk to Campaigns
       through a common source (CRM or manual export) so that all brands’ customers receive only
       relevant, branded communications.
     • Other Zoho One App Integrations: Take advantage of other apps to enhance your multi-brand
       support:
     • Zoho Analytics: For advanced reporting, connect Zoho Desk to Zoho Analytics. This lets you create
       unified dashboards that can compare performance across brands or combine data (for example,
       overall support load across all brands, as well as brand-specific breakdowns). In Analytics, you could
       create a “Support Overview” dashboard showing total tickets, average CSAT by brand, etc., which
       helps management see the big picture and each brand’s details in one place.
     • Zoho Projects: If your support process involves escalating tickets to development or operations
       teams, integrate Desk with Zoho Projects (or Zoho Sprints). An issue reported for Brand A can be
       converted into a bug or task in Projects with a link back to the ticket. This ensures your teams
       beyond support are looped in, and progress can be tracked from the Desk ticket. Each brand’s tickets
       can go to the respective project or a common project with a tag for the brand.




                                                      6
     • Zoho Assist: For technical support brands, Zoho Assist integration allows agents to start a remote
       support session directly from a ticket. This is great for a unified system – no matter which brand,
       agents can click “Start Assist Session” and help the customer instantly, improving first-call resolution
       rates.
     • Zoho Survey: You might use Zoho Desk’s built-in CSAT survey or integrate Zoho Survey for more
       detailed feedback. Surveys can be sent post-resolution to gauge satisfaction per brand. Because
       everything is under one roof, you can consistently measure and compare brand-wise customer
       satisfaction.
     • Zoho Flow/Webhooks: If any specific integration is not direct, Zoho’s integration platform (Flow) can
       connect Desk to third-party apps or other Zoho apps. For example, if you want to update Zoho Books
       or send an SMS via Clickatell when a high-priority ticket comes in, you can design that workflow.
     • Zoho SalesIQ (Chat) Integration: (As discussed in the live chat section, ensure SalesIQ is integrated
       with Desk.) This means chat transcripts can be attached to tickets, and you can convert a chat into a
       Desk ticket if it needs follow-up. Also, use SalesIQ’s visitor tracking to see if a customer who is
       chatting has an existing ticket – agents can then link the chat to that ticket for context.
     • Zoho Radar for Real-Time Insights: Zoho Desk Radar is a mobile app that gives you real-time
       dashboards of your support operation 25 . Set up Zoho Radar on your phone or tablet to monitor
       live metrics across all brands on the go. Radar provides an overview of vital stats like current live
       traffic (customers on your help center or waiting in chat), channel traffic (e.g., how many emails vs
       chats vs calls are coming in), ticket volumes, and customer happiness ratings in real time 25 . You
       and your managers can see at a glance how each department is doing: for example, Radar can show
       if Brand A suddenly has a spike in tickets or a drop in happiness rating, so you can react immediately.
       Agents can also use Radar to get a quick view of tickets that need urgent attention (it highlights
       tickets that are about to breach SLA or are high priority) 26 . To use it effectively, configure the app
       to include all your departments. You can customize the Radar dashboard to show widgets per
       department or combined, depending on your preference. This real-time visibility is crucial in a multi-
       brand environment – it ensures nothing falls through the cracks. If one brand’s support queue is
       overloaded, you’ll see it live and can allocate resources or jump in to assist. Likewise, monitoring
       “Happiness Ratings” (customer feedback) by brand in Radar can alert you to any brand-specific
       issues in service quality that need attention.

Summary: By integrating Zoho Desk with CRM and other Zoho One apps, you create a unified ecosystem
where sales, support, and marketing data all interconnect. This not only improves internal efficiency (agents
have context, marketers know support history, etc.) but also enhances the customer experience (e.g.,
consistent communication and faster, informed support). Meanwhile, tools like Zoho Radar and Zoho
Analytics give you the oversight to keep the support performance on track across all brands in real time.
With these integrations and monitoring in place, your multi-brand support system will be scalable, low-
maintenance, and consistently delivering a professional, on-brand experience to your customers at every
touchpoint.




                                                      7
 1    2    3   5   8   9   10   11   12   Providing Multi-brand Customer Support - Zoho Desk Knowledgebase
https://help.zoho.com/portal/en/kb/desk/support-channels/help-center/articles/providing-multi-brand-customer-support

 4   Customer Self-Service Portal | Build Multi-Brand Help Center
https://www.zoho.com/desk/multibrand-help-center.html

 6    7   21   Adding and Managing Departments - Online Help | Zoho Desk
https://help.zoho.com/portal/en/kb/desk/organization-settings/articles/managing-departments-in-help-desk

13   14   Assign multiple departments to multiple helpcenters
https://help.zoho.com/portal/en/community/topic/create-direct-assignment-to-a-specific-layout

15   20   Effortlessly implement and manage SalesIQ for your business with Brands!
https://help.zoho.com/portal/en/community/topic/effortlessly-implement-and-manage-salesiq-for-your-business-with-brands

16   17   18   Set chat window flow - Zoho SalesIQ
https://help.zoho.com/portal/en/kb/salesiq-2-0/for-administrators/setup-brand/articles/chat-window-flow

19   Zia | Zoho Desk's AI for simplified customer service
https://www.zoho.com/desk/zia.html

22   23   24   Zoho CRM Integration - Zoho Desk Knowledgebase
https://help.zoho.com/portal/en/kb/desk/integrations-and-marketplace/sales-and-marketing/articles/integrating-zoho-desk-with-
zoho-crm

25   26   Radar by Zoho Desk - All the metrics you need at your fingertips.
https://www.zoho.com/desk/radar.html




                                                              8

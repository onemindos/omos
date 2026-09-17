---
id: ge-79-20260113225956
title: "Step by step Guide to setup zoho desk to handle al"
type: reference
path: GE
domain: 75-99
folder: "79"
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [path/ge, path/ge/79, type/reference, status/active, platform/zoho]
ai_summary: "Step by step Guide to setup zoho desk to handle al"
---
# Step by step Guide to setup zoho desk to handle all customer inqurires

### 🔹 **Step 1: Set Up [support.legacyplatform.co](http://support.legacyplatform.co/) as Your Custom Domain in Zoho Desk**

Using [**support.legacyplatform.co**](http://support.legacyplatform.co/) as the main helpdesk URL makes your brand more professional and provides **one central place** for customers.

### 🛠 **How to Set Up a Custom Domain in Zoho Desk**

1️⃣ **Go to Zoho Desk** → Admin Panel

2️⃣ Click on **Rebranding > Custom Domain**

3️⃣ Add [**support.legacyplatform.co**](http://support.legacyplatform.co/)

4️⃣ Follow the **CNAME instructions** given by Zoho Desk

- **Create a CNAME record** in your DNS provider (e.g., Cloudflare, GoDaddy, Namecheap):
    - **Host:** support
    - **Value:**[desk.cs.zohohost.com](http://desk.cs.zohohost.com/)
    - **TTL:** Auto or 3600 seconds
    - 5️⃣ Click **Verify & Save**
    - 6️⃣ Once verified, **set it as the default access URL** in Zoho Desk

✅ **After setup, customers can access your support portal at** → [support.legacyplatform.co](http://support.legacyplatform.co/)

### 🔹 **Step 2: Configure [support@legacyplatform.co](mailto:support@legacyplatform.co) as the Official Support Email**

Your support email should be connected directly to **Zoho Desk**, so all inquiries automatically create tickets.

### 🛠 **How to Set Up [support@legacyplatform.co](mailto:support@legacyplatform.co) in Zoho Desk**

1️⃣ **Go to Zoho Desk** → Admin Panel

2️⃣ Click **Email > Email Settings**

3️⃣ Select **"Add Email Address"**

4️⃣ Enter [**support@legacyplatform.co**](mailto:support@legacyplatform.co)

5️⃣ Choose **Forwarding-Based Setup** (Recommended)

6️⃣ Go to your **email provider’s DNS settings** (Zoho Mail, Google Workspace, or others)

- **Set up email forwarding** to Zoho Desk’s provided email (e.g., [support@legacyplatform.zohodesk.com](mailto:support@legacyplatform.zohodesk.com)).
- If using **Zoho Mail**, enable email routing from Zoho Mail to Zoho Desk.
- 7️⃣ Save & Verify

✅ **All emails sent to [support@legacyplatform.co](mailto:support@legacyplatform.co) will now create tickets inside Zoho Desk.**

### 🔹 **Step 3: Set Up Departments in Zoho Desk (For Routing Inquiries)**

To ensure all customer inquiries go to the **right team**, set up departments inside Zoho Desk.

### 🛠 **Recommended Department Structure**

📌 **Go to Zoho Desk** → **Admin > Departments** → **Add Department**

| **Department Name** | **Handles** | **Example Emails Processed** |
| --- | --- | --- |
| **Customer Support** | Tech issues, troubleshooting, FAQs | Product login issues, service questions |
| **Sales & Pre-Sales** | New customer inquiries, upgrades | Sales inquiries, pricing questions |
| **Billing & Payments** | Refunds, invoices, subscription issues | Failed payment, refund requests |
| **Partnerships & Affiliates** | Business collaborations, affiliate inquiries | Partner program questions, co-branding |
| **General Inquiries** | Miscellaneous questions | PR, general messages |

✅ **Set up ticket routing** so that:

- **Billing requests go to Finance**
- **Sales leads go to the Sales Team**
- **Technical support goes to Customer Support**

### 🔹 **Step 4: Automate Ticket Assignment & Prioritization**

To prevent **manual sorting**, set up **Zoho Desk Automations**.

### 🛠 **How to Enable Smart Ticket Routing**

📌 **Go to Zoho Desk > Automation > Assignment Rules**

1️⃣ **Create an "Auto-Assign Rule"** → Name it **"Route Tickets to the Right Team"**

2️⃣ **Add Conditions Based on Keywords**

- **If email contains "refund" or "billing" → Assign to Billing**
- **If email contains "partnership" → Assign to Partnerships Team**
- **If email contains "pricing" or "quote" → Assign to Sales**
- 3️⃣ **Set Response Priorities**
- **Critical (Billing & Refund Issues)** → Response in **24 hours**
- **Sales Inquiries** → Response in **48 hours**
- **General Support** → Response in **72 hours**

✅ **Now, every email/ticket is auto-sorted & routed to the correct department.**

### 🔹 **Step 5: Customize Your Customer Portal & Automate Responses**

### 🛠 **How to Enable a Self-Service Portal on [support.legacyplatform.co](http://support.legacyplatform.co/)**

📌 **Go to Zoho Desk > Help Center Settings**

1️⃣ Enable **"Customer Portal"**

2️⃣ **Add a Knowledge Base** (FAQ articles)

3️⃣ Allow customers to **log in & track tickets**

4️⃣ Add a **Live Chat Widget** (via Zoho SalesIQ if needed)

✅ **Customers can now visit [support.legacyplatform.co](http://support.legacyplatform.co/) to:**

- Submit tickets
- Track existing requests
- Access help articles

### 🔹 **Step 6: Enable SLAs (Service Level Agreements) for Faster Response Times**

To ensure fast customer support, **set response deadlines** for each request type.

📌 **Go to Zoho Desk > SLA Rules**

1️⃣ Set SLAs for **each department**:

- Billing → Response in **24 hours**
- Sales → Response in **48 hours**
- General Support → Response in **72 hours**
- 2️⃣ Set **Auto-Escalation Rules**:
- If **Billing tickets are not answered in 24 hours**, escalate to **Finance Head**
- If **Sales inquiries are not answered in 48 hours**, escalate to **Sales Lead**

✅ **Now, Zoho Desk automatically enforces response time commitments!**

### ✅ **Final Summary**

🎯 **You’ve now fully set up Zoho Desk to handle support, sales, and billing inquiries efficiently.**

🔹 ✅[**support.legacyplatform.co**](http://support.legacyplatform.co/) → Your central customer support hub.

🔹 ✅ **[support@legacyplatform.co](mailto:support@legacyplatform.co)** → The universal contact email.

🔹 ✅ **Departments, Automations & SLAs** → Ensures fast, organized responses.

🔹 ✅ **Self-Service Portal & Tracking** → Customers can check ticket statuses anytime.

### 🚀 **Next Steps**

✅ **Test the ticketing system** (send emails to [support@legacyplatform.co](mailto:support@legacyplatform.co) and confirm auto-routing).

✅ **Train your support team** on the Zoho Desk workflow.

✅ **Publish your Knowledge Base** (FAQs & guides).

✅ **Monitor performance via Zoho Desk Reports** (response times, ticket resolution stats).

🔹 **Would you like a template for Knowledge Base articles to include in Zoho Desk?** 🚀
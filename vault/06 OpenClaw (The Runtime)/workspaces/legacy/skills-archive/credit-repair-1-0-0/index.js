/**
 * credit-repair-1.0.0 — OpenClaw Skill
 *
 * Generate FCRA-compliant dispute letters for credit bureau errors.
 * Produces HTML-ready letters for sending via the lob skill.
 * Covers general disputes, method-of-verification follow-ups,
 * and debt validation requests.
 *
 * No external API calls — pure letter generation.
 *
 * Registers tools: credit_repair
 */

"use strict";

// ---------------------------------------------------------------------------
// Credit bureau addresses
// ---------------------------------------------------------------------------

const BUREAUS = {
  experian: {
    name: "Experian Disputes",
    address_line1: "P.O. Box 4500",
    city: "Allen",
    state: "TX",
    zip: "75013",
  },
  equifax: {
    name: "Equifax Information Services LLC",
    address_line1: "P.O. Box 740256",
    city: "Atlanta",
    state: "GA",
    zip: "30374",
  },
  transunion: {
    name: "TransUnion Consumer Solutions",
    address_line1: "P.O. Box 2000",
    city: "Chester",
    state: "PA",
    zip: "19016",
  },
  transunion_disputes: {
    name: "TransUnion LLC Consumer Dispute Center",
    address_line1: "P.O. Box 34012",
    city: "Fullerton",
    state: "CA",
    zip: "92834",
  },
};

// ---------------------------------------------------------------------------
// Letter templates (return HTML for Lob)
// ---------------------------------------------------------------------------

function generalDisputeLetter({
  sender_name,
  sender_address,
  sender_city,
  sender_state,
  sender_zip,
  bureau,
  ssn_last4,
  creditor_name,
  account_last4,
  dispute_reason,
  explanation,
  enclosures,
}) {
  const b = BUREAUS[bureau];
  const date = new Date().toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return `<html>
<head><style>
  body { font-family: Georgia, serif; font-size: 12pt; line-height: 1.6; margin: 1in; color: #000; }
  .header { margin-bottom: 24pt; }
  .bureau { margin-bottom: 24pt; }
  .section { margin-bottom: 12pt; }
  .signature { margin-top: 36pt; }
</style></head>
<body>
  <div class="header">
    <p>${esc(sender_name)}<br>
    ${esc(sender_address)}<br>
    ${esc(sender_city)}, ${esc(sender_state)} ${esc(sender_zip)}</p>
    <p>${date}</p>
  </div>

  <div class="bureau">
    <p>${esc(b.name)}<br>
    ${esc(b.address_line1)}<br>
    ${esc(b.city)}, ${esc(b.state)} ${esc(b.zip)}</p>
  </div>

  <p><strong>RE: Formal Dispute of Inaccurate Information</strong><br>
  SSN: XXX-XX-${esc(ssn_last4)}</p>

  <p>Dear Dispute Department:</p>

  <p>I am writing pursuant to the Fair Credit Reporting Act, 15 U.S.C. &sect; 1681i, to dispute
  the following inaccurate information on my credit report:</p>

  <div class="section">
    <p><strong>DISPUTED ITEM:</strong><br>
    Creditor Name: ${esc(creditor_name)}<br>
    Account Number: XXXX-XXXX-XXXX-${esc(account_last4)}<br>
    Reason for Dispute: ${esc(dispute_reason)}</p>
  </div>

  <p>This information is inaccurate because ${esc(explanation)}.</p>

  <p>Under the FCRA, you are required to:</p>
  <ol>
    <li>Conduct a reasonable investigation within 30 days</li>
    <li>Forward all relevant information to the furnisher</li>
    <li>Delete the item if it cannot be verified</li>
    <li>Provide me with written results within 5 business days of completion</li>
    <li>Provide a free copy of my updated credit report</li>
  </ol>

  <p>I request that you investigate this matter and correct or delete the disputed item.</p>

  <div class="signature">
    <p>Sincerely,</p>
    <br><br>
    <p>${esc(sender_name)}</p>
  </div>

  ${enclosures ? `<p><strong>Enclosures:</strong> ${esc(enclosures)}</p>` : ""}
</body>
</html>`;
}

function methodOfVerificationLetter({
  sender_name,
  sender_address,
  sender_city,
  sender_state,
  sender_zip,
  bureau,
  ssn_last4,
  original_dispute_date,
  creditor_name,
  account_last4,
}) {
  const b = BUREAUS[bureau];
  const date = new Date().toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return `<html>
<head><style>
  body { font-family: Georgia, serif; font-size: 12pt; line-height: 1.6; margin: 1in; color: #000; }
  .header { margin-bottom: 24pt; }
  .bureau { margin-bottom: 24pt; }
  .signature { margin-top: 36pt; }
</style></head>
<body>
  <div class="header">
    <p>${esc(sender_name)}<br>
    ${esc(sender_address)}<br>
    ${esc(sender_city)}, ${esc(sender_state)} ${esc(sender_zip)}</p>
    <p>${date}</p>
  </div>

  <div class="bureau">
    <p>${esc(b.name)}<br>
    ${esc(b.address_line1)}<br>
    ${esc(b.city)}, ${esc(b.state)} ${esc(b.zip)}</p>
  </div>

  <p><strong>RE: Request for Method of Verification &mdash; Follow-up Dispute</strong><br>
  SSN: XXX-XX-${esc(ssn_last4)}<br>
  Creditor: ${esc(creditor_name)}<br>
  Account: XXXX-XXXX-XXXX-${esc(account_last4)}</p>

  <p>Dear Dispute Department:</p>

  <p>I received your response dated ${esc(original_dispute_date)} stating the
  disputed item was &ldquo;verified.&rdquo;</p>

  <p>Under FCRA &sect; 1681i(a)(7), I am entitled to know the method of verification
  used. Please provide:</p>

  <ol>
    <li>Name, address, and phone number of the person who verified this information</li>
    <li>The specific documents or records used to verify accuracy</li>
    <li>A description of the verification procedure used</li>
  </ol>

  <p>If you cannot provide this information within 15 days, I request immediate
  deletion of this item as it has not been properly verified per FCRA requirements.</p>

  <div class="signature">
    <p>Sincerely,</p>
    <br><br>
    <p>${esc(sender_name)}</p>
  </div>
</body>
</html>`;
}

function debtValidationLetter({
  sender_name,
  sender_address,
  sender_city,
  sender_state,
  sender_zip,
  bureau,
  ssn_last4,
  creditor_name,
  account_last4,
  alleged_amount,
}) {
  const b = BUREAUS[bureau];
  const date = new Date().toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return `<html>
<head><style>
  body { font-family: Georgia, serif; font-size: 12pt; line-height: 1.6; margin: 1in; color: #000; }
  .header { margin-bottom: 24pt; }
  .bureau { margin-bottom: 24pt; }
  .signature { margin-top: 36pt; }
</style></head>
<body>
  <div class="header">
    <p>${esc(sender_name)}<br>
    ${esc(sender_address)}<br>
    ${esc(sender_city)}, ${esc(sender_state)} ${esc(sender_zip)}</p>
    <p>${date}</p>
  </div>

  <div class="bureau">
    <p>${esc(b.name)}<br>
    ${esc(b.address_line1)}<br>
    ${esc(b.city)}, ${esc(b.state)} ${esc(b.zip)}</p>
  </div>

  <p><strong>RE: Debt Validation Request Under FDCPA</strong><br>
  SSN: XXX-XX-${esc(ssn_last4)}<br>
  Creditor: ${esc(creditor_name)}<br>
  Account: XXXX-XXXX-XXXX-${esc(account_last4)}${alleged_amount ? `<br>Alleged Amount: $${esc(alleged_amount)}` : ""}</p>

  <p>Dear Dispute Department:</p>

  <p>I am disputing this alleged debt and requesting validation under the Fair Debt
  Collection Practices Act, 15 U.S.C. &sect; 1692g.</p>

  <p>Please provide:</p>
  <ol>
    <li>Verification of the amount claimed</li>
    <li>Name of the original creditor</li>
    <li>Copy of the original signed agreement</li>
    <li>Complete payment history</li>
    <li>Proof you are licensed to collect in my state</li>
  </ol>

  <p>Until you provide proper validation, I demand you:</p>
  <ul>
    <li>Cease all collection activity</li>
    <li>Remove this item from all credit bureaus</li>
    <li>Stop all communication except to provide validation</li>
  </ul>

  <div class="signature">
    <p>Sincerely,</p>
    <br><br>
    <p>${esc(sender_name)}</p>
  </div>
</body>
</html>`;
}

// ---------------------------------------------------------------------------
// HTML escape
// ---------------------------------------------------------------------------

function esc(s) {
  if (!s) return "";
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

// ---------------------------------------------------------------------------
// Tool: credit_repair
// ---------------------------------------------------------------------------

async function credit_repair({
  action,
  letter_type,
  bureau,
  // sender info
  sender_name,
  sender_address,
  sender_city,
  sender_state,
  sender_zip,
  // account info
  ssn_last4,
  creditor_name,
  account_last4,
  dispute_reason,
  explanation,
  enclosures,
  // method of verification
  original_dispute_date,
  // debt validation
  alleged_amount,
}) {
  try {
    // =================================================================
    // GENERATE DISPUTE LETTER
    // =================================================================

    if (action === "generate_letter") {
      if (!letter_type) return JSON.stringify({ error: "letter_type required: general_dispute, method_of_verification, debt_validation" });
      if (!bureau) return JSON.stringify({ error: "bureau required: experian, equifax, transunion, transunion_disputes" });
      if (!BUREAUS[bureau]) return JSON.stringify({ error: `Unknown bureau: ${bureau}. Use: experian, equifax, transunion, transunion_disputes` });
      if (!sender_name || !sender_address || !sender_city || !sender_state || !sender_zip)
        return JSON.stringify({ error: "Sender info required: sender_name, sender_address, sender_city, sender_state, sender_zip" });
      if (!ssn_last4) return JSON.stringify({ error: "ssn_last4 required (last 4 digits only)" });

      let html;
      const b = BUREAUS[bureau];

      if (letter_type === "general_dispute") {
        if (!creditor_name || !account_last4 || !dispute_reason || !explanation)
          return JSON.stringify({ error: "general_dispute requires: creditor_name, account_last4, dispute_reason, explanation" });
        html = generalDisputeLetter({
          sender_name, sender_address, sender_city, sender_state, sender_zip,
          bureau, ssn_last4, creditor_name, account_last4, dispute_reason, explanation, enclosures,
        });
      } else if (letter_type === "method_of_verification") {
        if (!creditor_name || !account_last4 || !original_dispute_date)
          return JSON.stringify({ error: "method_of_verification requires: creditor_name, account_last4, original_dispute_date" });
        html = methodOfVerificationLetter({
          sender_name, sender_address, sender_city, sender_state, sender_zip,
          bureau, ssn_last4, original_dispute_date, creditor_name, account_last4,
        });
      } else if (letter_type === "debt_validation") {
        if (!creditor_name || !account_last4)
          return JSON.stringify({ error: "debt_validation requires: creditor_name, account_last4" });
        html = debtValidationLetter({
          sender_name, sender_address, sender_city, sender_state, sender_zip,
          bureau, ssn_last4, creditor_name, account_last4, alleged_amount,
        });
      } else {
        return JSON.stringify({ error: `Unknown letter_type: ${letter_type}. Use: general_dispute, method_of_verification, debt_validation` });
      }

      return JSON.stringify({
        letter_type,
        bureau: bureau,
        bureau_address: b,
        html,
        next_step: "Use the lob skill to send: lob({ action: 'send_letter', to_name: bureau_address.name, to_address_line1: bureau_address.address_line1, to_city: bureau_address.city, to_state: bureau_address.state, to_zip: bureau_address.zip, content: html, ... })",
        success: true,
      }, null, 2);
    }

    // =================================================================
    // GET BUREAU ADDRESS
    // =================================================================

    if (action === "get_bureau_address") {
      if (!bureau) return JSON.stringify({ error: "bureau required" });
      if (!BUREAUS[bureau])
        return JSON.stringify({
          error: `Unknown bureau: ${bureau}`,
          available: Object.keys(BUREAUS),
        });
      return JSON.stringify({ bureau, address: BUREAUS[bureau] }, null, 2);
    }

    // =================================================================
    // LIST DISPUTE REASONS
    // =================================================================

    if (action === "list_dispute_reasons") {
      return JSON.stringify({
        valid_grounds: [
          { reason: "Not My Account", description: "Identity theft, mixed files, wrong SSN" },
          { reason: "Incorrect Balance", description: "Amount shown is wrong" },
          { reason: "Incorrect Payment Status", description: "Shows late but paid on time" },
          { reason: "Account Closed by Consumer", description: "Shows closed by creditor" },
          { reason: "Paid in Full", description: "Shows balance when paid off" },
          { reason: "Included in Bankruptcy", description: "Not updated after discharge" },
          { reason: "Outdated Information", description: "Past 7-year reporting period" },
          { reason: "Duplicate Account", description: "Same debt listed multiple times" },
          { reason: "Wrong Date", description: "Incorrect open date or delinquency date" },
          { reason: "Never Late", description: "Shows late payment that didn't happen" },
        ],
        avoid: [
          "I don't recognize this (without evidence)",
          "This hurts my score (not a legal reason)",
          "Disputing accurate information",
        ],
        success: true,
      }, null, 2);
    }

    // =================================================================
    // FCRA RIGHTS REFERENCE
    // =================================================================

    if (action === "fcra_rights") {
      return JSON.stringify({
        key_rights: [
          { right: "Accuracy", section: "§ 1681e(b)", description: "Bureaus must follow reasonable procedures for maximum accuracy" },
          { right: "Investigation", section: "§ 1681i(a)", description: "Must investigate disputes within 30 days (45 with documentation)" },
          { right: "Deletion", section: "§ 1681i(a)(5)", description: "Must delete information that cannot be verified" },
          { right: "Notification", section: "§ 1681i(a)(6)", description: "Must notify consumer of results within 5 days" },
          { right: "Free Report", section: "§ 1681j", description: "Free annual report + free report after dispute" },
          { right: "Furnisher Duties", section: "§ 1681s-2", description: "Data furnishers must investigate and correct errors" },
        ],
        timelines: {
          investigation: "30 days (45 with additional docs)",
          notification: "5 business days after investigation",
          negative_items: "7 years (10 for bankruptcy)",
        },
        escalation: [
          "File CFPB complaint: consumerfinance.gov/complaint",
          "File FTC complaint: reportfraud.ftc.gov",
          "Contact state Attorney General",
          "FCRA lawsuit: $100-$1,000 statutory damages per violation",
        ],
        success: true,
      }, null, 2);
    }

    // =================================================================
    // UNKNOWN ACTION
    // =================================================================

    return JSON.stringify({
      error: `Unknown action: ${action}`,
      available: [
        "generate_letter",
        "get_bureau_address",
        "list_dispute_reasons",
        "fcra_rights",
      ],
    });
  } catch (err) {
    return JSON.stringify({ error: err.message });
  }
}

module.exports = { credit_repair };

---
name: credit-bureau-dispute
description: "Use when drafting FCRA 605B identity-theft block letters to the credit bureaus (procedure only, no personal data)."
version: 2.0.0
author: Legacy
tags: [credit, fcra, identity-theft, dispute, docx, legal]
---

# Credit Bureau Dispute Letters (FCRA 605B)

Procedure for producing one .docx cover letter per bureau from an FTC Identity
Theft Report. This skill holds the method only. It never stores names, addresses,
Social Security numbers, dates of birth, account numbers or any credit-file data.
Ask Zeus for those at run time, use them in working files only, and do not write
them into Git, memory, cron output or chat summaries.

## Authority

- Drafting letters is a software task and needs no extra approval.
- Mailing, signing, paying postage or sending anything to a bureau is a
  real-world action: stop and let Zeus do it. If money is involved, Zeus must
  write "execute" or "override" for that specific action (SOUL rule 8).
- Never delete the source PDFs, drafts or scans (SOUL rule 4). Quarantine and ask.
- Copying letters to the Mac goes over SSH on the tailnet and works only while
  the Mac is online. If it is offline, leave the files in the pod scratch
  directory and tell Zeus.

## Bureaus

Equifax, Experian, TransUnion and Innovis. Always include Innovis: it is missing
from most templates but holds a separate file some lenders use. A 605B fraud
block goes to the bureau's fraud department, not the general dispute box (the
Equifax one is a common mistake). Addresses and phone numbers change, so look
them up on each bureau's site at run time and cite the page. Do not rely on
remembered addresses.

## Procedure

1. **Extract** from the FTC Identity Theft Report PDF (`read_file` extracts
   text): full name, address, FTC report number, and for each fraudulent account
   the company, type, last 3 to 4 digits of the account number and the
   month/year the fraud began.
2. **Build the letters** with python-docx in a seeded venv (the pod blocks system
   pip): `uv venv <dir> --seed --clear`, then `<dir>/bin/pip install python-docx`.
   Write the generator script with `write_file`, then run it with the venv Python.
3. **Letter content**, one letter per bureau:
   - Sender block, date, bureau fraud-department address.
   - Bold subject: `RE: Identity Theft - Request to Block Fraudulent Information`.
   - Opening that states Zeus is an identity-theft victim.
   - Numbered list: `Company - Type - Account ending XXXX - Fraud began: MM/YYYY`.
   - The FCRA 605B block-request paragraph.
   - Enclosures list and a signature line (Zeus signs by hand).
4. **Deliver** to the working folder Zeus names, over SSH/SCP using the pod's SSH
   identity, only while the Mac is online.

## Enclosures per mailing

Signed cover letter; copy of the credit report with fraudulent items marked; the
FTC report; a copy of government photo ID (Zeus supplies it, never store it
here); the FCRA section 605B text from the FTC website.

## Mailing notes

- Certified Mail with return receipt, one envelope per bureau, never combined.
- A packet of about 30 pages fits a USPS Large Envelope (Flat); confirm weight at the counter.

## Pitfalls

- Redact account numbers in the letter body to the last 3 to 4 digits.
- Use `uv venv --seed`; without `--seed` the venv has no pip.
- The pod's SSH identity lives under the pod home `.ssh` (check
  `ls /opt/data/home/.ssh`), not `/opt/data/.ssh`, which only holds known_hosts.
- Do not paste the FTC personal statement or any victim data into other skills,
  logs or cron output.

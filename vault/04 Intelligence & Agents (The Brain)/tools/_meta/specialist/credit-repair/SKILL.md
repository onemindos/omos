---
name: credit-repair
description: Dispute credit report errors and improve credit scores using FCRA rights
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["finance", "credit", "disputes", "fcra", "letters"]
  allowed-tools: ["lob", "send_letter", "verify_address", "get_letter_status"]
---

# Credit Repair Skill

Expert knowledge for disputing inaccurate credit report information and improving credit scores using consumer rights under the Fair Credit Reporting Act (FCRA).

## When to Use

- User wants to dispute errors on their credit report
- User needs help writing dispute letters to credit bureaus
- User wants to understand their FCRA rights
- User needs to track dispute deadlines and responses
- User wants to improve their credit score

## Credit Bureau Mailing Addresses

Always use these official dispute addresses:

### Experian
```
Experian Disputes
P.O. Box 4500
Allen, TX 75013
```

### Equifax
```
Equifax Information Services LLC
P.O. Box 740256
Atlanta, GA 30374
```

### TransUnion
```
TransUnion Consumer Solutions
P.O. Box 2000
Chester, PA 19016

Alternative (Disputes):
TransUnion LLC
Consumer Dispute Center
P.O. Box 34012
Fullerton, CA 92834
```

## FCRA Consumer Rights (15 U.S.C. § 1681)

### Key Rights

| Right | FCRA Section | Description |
|-------|--------------|-------------|
| **Accuracy** | § 1681e(b) | Bureaus must follow reasonable procedures for maximum accuracy |
| **Investigation** | § 1681i(a) | Must investigate disputes within 30 days (45 with documentation) |
| **Deletion** | § 1681i(a)(5) | Must delete information that cannot be verified |
| **Notification** | § 1681i(a)(6) | Must notify consumer of investigation results within 5 days |
| **Free Report** | § 1681j | Free annual report + free report after dispute |
| **Furnisher Duties** | § 1681s-2 | Data furnishers must investigate and correct errors |

### Important Timelines

- **30 days**: Bureau must complete investigation
- **45 days**: Extended if consumer provides additional documentation
- **5 days**: Bureau must notify consumer of results after investigation
- **7 years**: Most negative items must be removed (10 years for bankruptcy)

## Dispute Letter Templates

### Template 1: General Dispute

```
[Your Name]
[Your Address]
[City, State ZIP]

[Date]

[Credit Bureau Name]
[Bureau Address]

RE: Formal Dispute of Inaccurate Information
SSN: XXX-XX-[Last 4]

Dear Dispute Department:

I am writing pursuant to the Fair Credit Reporting Act, 15 U.S.C. § 1681i, to dispute
the following inaccurate information on my credit report:

DISPUTED ITEM:
Creditor Name: [CREDITOR]
Account Number: XXXX-XXXX-XXXX-[Last 4]
Reason for Dispute: [SPECIFIC REASON]

This information is inaccurate because [EXPLAIN WHY].

Under the FCRA, you are required to:
1. Conduct a reasonable investigation within 30 days
2. Forward all relevant information to the furnisher
3. Delete the item if it cannot be verified
4. Provide me with written results within 5 business days of completion
5. Provide a free copy of my updated credit report

I request that you investigate this matter and correct or delete the disputed item.

Sincerely,

[Signature]
[Printed Name]

Enclosures: [List any supporting documents]
```

### Template 2: Method of Verification Request

Use this AFTER initial dispute if bureau claims item was "verified":

```
RE: Request for Method of Verification - Follow-up Dispute

I received your response dated [DATE] stating the disputed item was "verified."

Under FCRA § 1681i(a)(7), I am entitled to know the method of verification used.
Please provide:

1. Name, address, and phone number of the person who verified this information
2. The specific documents or records used to verify accuracy
3. A description of the verification procedure used

If you cannot provide this information within 15 days, I request immediate deletion
of this item as it has not been properly verified per FCRA requirements.
```

### Template 3: Debt Validation (for Collections)

```
RE: Debt Validation Request Under FDCPA

I am disputing this alleged debt and requesting validation under the Fair Debt
Collection Practices Act, 15 U.S.C. § 1692g.

Please provide:
1. Verification of the amount claimed
2. Name of the original creditor
3. Copy of the original signed agreement
4. Complete payment history
5. Proof you are licensed to collect in my state

Until you provide proper validation, I demand you:
- Cease all collection activity
- Remove this item from all credit bureaus
- Stop all communication except to provide validation
```

## Common Dispute Reasons

### Valid Dispute Grounds

1. **Not My Account** - Identity theft, mixed files, wrong SSN
2. **Incorrect Balance** - Amount shown is wrong
3. **Incorrect Payment Status** - Shows late but paid on time
4. **Account Closed by Consumer** - Shows closed by creditor
5. **Paid in Full** - Shows balance when paid off
6. **Included in Bankruptcy** - Not updated after discharge
7. **Outdated Information** - Past 7-year reporting period
8. **Duplicate Account** - Same debt listed multiple times
9. **Wrong Date** - Incorrect open date or delinquency date
10. **Never Late** - Shows late payment that didn't happen

### Weak/Invalid Grounds (Avoid)

- "I don't recognize this" (without evidence)
- "This hurts my score" (not a legal reason)
- Disputing accurate information

## Dispute Process Workflow

### Phase 1: Gather Information

1. Get free credit reports from annualcreditreport.com
2. Identify specific errors on each bureau's report
3. Collect supporting documentation
4. Note account numbers, dates, and amounts

### Phase 2: Send Initial Disputes

1. Draft dispute letter for each error
2. Verify bureau addresses before sending
3. Send via Lob (creates tracking + delivery confirmation)
4. Save letter ID for tracking
5. Set 30-day reminder for response deadline

### Phase 3: Track and Follow Up

1. Check letter delivery status
2. Wait for bureau response (30-35 days)
3. If no response: Send follow-up citing FCRA violation
4. If "verified" without evidence: Request method of verification
5. If deleted: Confirm on updated report

### Phase 4: Escalate if Needed

1. File complaint with CFPB (consumerfinance.gov/complaint)
2. File complaint with FTC (reportfraud.ftc.gov)
3. Contact state Attorney General
4. Consider FCRA lawsuit (statutory damages $100-$1,000 per violation)

## Best Practices

### DO

- Be specific about what is inaccurate
- Reference FCRA sections
- Include last 4 of SSN (never full SSN)
- Keep copies of everything
- Send one dispute per letter for clarity
- Use certified mail or tracked delivery (Lob provides this)
- Follow up if no response in 35 days

### DON'T

- Dispute accurate information
- Use template letters word-for-word (customize)
- Send disputes to wrong addresses
- Include full SSN in letters
- Threaten lawsuits in initial letter
- Dispute everything at once (looks frivolous)
- Use credit repair company template language

## Using the Lob Tool

### Send Dispute Letter

```python
send_letter(
    to_name="Experian Disputes",
    to_address="P.O. Box 4500",
    to_city="Allen",
    to_state="TX",
    to_zip="75013",
    from_name="[User's Full Name]",
    from_address="[User's Street Address]",
    from_city="[User's City]",
    from_state="[User's State]",
    from_zip="[User's ZIP]",
    content="<html>...[Letter Content]...</html>",
    use_type="operational"
)
```

### Check Letter Status

```python
get_letter_status(letter_id="ltr_xxx")
```

### Verify Address First

```python
verify_address(
    address="P.O. Box 4500",
    city="Allen",
    state="TX",
    zip_code="75013"
)
```

## Tracking Disputes

Create tasks to track:

1. **Letter Sent** - Date, bureau, letter_id
2. **30-Day Deadline** - When response is due
3. **Response Received** - Date and outcome
4. **Follow-Up Needed** - If no response or denied
5. **Resolution** - Final outcome

## Credit Score Impact

### Factors by Weight (FICO)

| Factor | Weight | Notes |
|--------|--------|-------|
| Payment History | 35% | Late payments, collections, bankruptcies |
| Amounts Owed | 30% | Credit utilization (keep under 30%) |
| Length of History | 15% | Average age of accounts |
| New Credit | 10% | Recent inquiries and new accounts |
| Credit Mix | 10% | Types of credit (cards, loans, mortgage) |

### Quick Wins

1. Remove inaccurate late payments (+50-100 points possible)
2. Remove collections (+25-50 points each)
3. Correct high balance errors (improves utilization)
4. Remove unauthorized inquiries (+5-10 points each)

## Legal Resources

- **CFPB**: consumerfinance.gov
- **FTC**: ftc.gov/credit
- **FCRA Full Text**: law.cornell.edu/uscode/text/15/chapter-41/subchapter-III
- **NCLC**: nclc.org (National Consumer Law Center)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]

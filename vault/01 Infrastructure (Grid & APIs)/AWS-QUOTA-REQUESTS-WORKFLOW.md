---
id: root-00-20260113225954
title: "AWS QUOTA REQUESTS WORKFLOW"
type: sop
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/sop, status/active]
ai_summary: "AWS QUOTA REQUESTS WORKFLOW"
---
# AWS Service Quota Requests — Workflow & Reference

> **Purpose:** This document provides instructions for responding to AWS Support quota increase requests. Follow this workflow to ensure consistent, accurate responses across sessions.

---

## Quick Reference — Account Details

| Field | Value |
|-------|-------|
| **AWS Account ID** | 654654489230 |
| **Company Name** | Legacy AI LLC |
| **Contact Name** | Zeus Delacruz |
| **Primary Region** | us-east-1 (N. Virginia) |
| **Support Plan** | Developer |
| **Limit Type** | On-demand |

---

## Current Quota Limits (as of 2026-01-13)

> **Run this to get fresh data:**
> ```bash
> aws service-quotas list-service-quotas --service-code bedrock --region us-east-1 --query 'Quotas[?contains(QuotaName, `Opus 4.5`) || contains(QuotaName, `Sonnet 4.5`)].[QuotaName, Value, QuotaCode]' --output text | sort
> ```

### Claude Opus 4.5

| Quota Type | Quota Code | Current Limit | Requested | Status |
|------------|------------|---------------|-----------|--------|
| CRIS RPM | L-27989F42 | 500 | 1,000 | ✅ Increased |
| CRIS TPM | L-7007E9C9 | 2,000,000 | 2,000,000 | ✅ At requested |
| Global CRIS RPM | L-58424D95 | 500 | 1,000 | ⚠️ Partial |
| Global CRIS TPM | L-3ABF6ACC | 2,000,000 | 5,000,000 | ⚠️ Partial |
| Global CRIS TPD | L-E2985A8E | 13,500,000 | — | Default |
| Model TPD | L-F397ECBF | 6,750,000 | — | Default |

### Claude Sonnet 4.5 V1 (Standard)

| Quota Type | Quota Code | Current Limit | Requested | Status |
|------------|------------|---------------|-----------|--------|
| CRIS RPM | L-4A6BFAB1 | 1,000 | 1,000 | ✅ At requested |
| CRIS TPM | L-F4DDD3EB | 5,000,000 | 5,000,000 | ✅ At requested |
| Global CRIS RPM | L-DB84CE56 | 1,000 | 1,000 | ✅ At requested |
| Global CRIS TPM | L-27C57EE8 | 5,000,000 | 5,000,000 | ✅ At requested |
| Global CRIS TPD | L-BC182137 | 10,800,000 | — | Default |
| Model TPD | L-381AD9EE | 5,400,000 | — | Default |

### Claude Sonnet 4.5 V1 — 1M Context Length

| Quota Type | Quota Code | Current Limit | Requested | Status |
|------------|------------|---------------|-----------|--------|
| CRIS RPM | L-A052927A | 6 | 500 | ⏳ Pending |
| CRIS TPM | L-8EA73537 | 1,000,000 | 5,000,000 | ⏳ Pending |
| Global CRIS RPM | L-C0D53EFB | 5 | 500 | ⏳ Pending |
| Global CRIS TPM | L-4B26E44A | 1,000,000 | 5,000,000 | ⏳ Pending |
| Global CRIS TPD | L-229A5DEC | 54,000,000 | — | Default |
| Model TPD | L-E107194C | 27,000,000 | — | Default |

### Summary

| Model | Status |
|-------|--------|
| **Sonnet 4.5 V1 (Standard)** | ✅ All quotas at requested levels |
| **Opus 4.5** | ⚠️ Partially increased (RPM at 500, TPM at 2M) |
| **Sonnet 4.5 V1 1M Context** | ⏳ Still at default low limits |

---

## Workflow Overview

### Step 1: Receive Case ID from User

User provides the AWS Support Case ID (e.g., `176828500800661`) and pastes the support email.

### Step 2: Look Up Quota Details

**Find the request ID** by listing all open cases:
```bash
aws service-quotas list-requested-service-quota-change-history --status CASE_OPENED
```

**Get specific case details** using the request ID:
```bash
aws service-quotas get-requested-service-quota-change --request-id "<REQUEST_ID>"
```

**Get current quota value**:
```bash
aws service-quotas get-service-quota --service-code bedrock --quota-code <QUOTA_CODE>
```

### Step 3: Get Paired Quota

Every RPM quota has a paired TPM quota and vice versa. Look up both to provide complete information.

**Common Quota Pairs:**

| Model | RPM Quota Code | TPM Quota Code |
|-------|---------------|----------------|
| Sonnet 4.5 V1 (Global CRIS) | L-DB84CE56 | L-27C57EE8 |
| Sonnet 4.5 V1 (CRIS) | L-4A6BFAB1 | (check) |
| Sonnet 4.5 1M (Global CRIS) | L-C0D53EFB | L-4B26E44A |
| Sonnet 4.5 1M (CRIS) | L-A052927A | L-8EA73537 |
| Opus 4.5 (Global CRIS) | L-58424D95 | L-3ABF6ACC |
| Opus 4.5 (CRIS) | L-27989F42 | L-7007E9C9 |

### Step 4: Check Actual Usage Data

**Get invocations:**
```bash
aws cloudwatch get-metric-statistics \
  --namespace AWS/Bedrock \
  --metric-name Invocations \
  --dimensions Name=ModelId,Value=<MODEL_ID> \
  --start-time $(date -u -v-7d +%Y-%m-%dT%H:%M:%SZ) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --period 86400 \
  --statistics Sum \
  --region us-east-1
```

**Get throttles (IMPORTANT - proves need for increase):**
```bash
aws cloudwatch get-metric-statistics \
  --namespace AWS/Bedrock \
  --metric-name InvocationThrottles \
  --dimensions Name=ModelId,Value=<MODEL_ID> \
  --start-time $(date -u -v-7d +%Y-%m-%dT%H:%M:%SZ) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --period 86400 \
  --statistics Sum \
  --region us-east-1
```

**Get input/output tokens:**
```bash
aws cloudwatch get-metric-statistics \
  --namespace AWS/Bedrock \
  --metric-name InputTokenCount \
  --dimensions Name=ModelId,Value=<MODEL_ID> \
  --start-time $(date -u -v-7d +%Y-%m-%dT%H:%M:%SZ) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --period 86400 \
  --statistics Sum Average \
  --region us-east-1
```

**Model IDs:**
- Sonnet 4.5: `us.anthropic.claude-sonnet-4-5-20250929-v1:0`
- Opus 4.5: `us.anthropic.claude-opus-4-5-20251101-v1:0`
- Haiku 4.5: `us.anthropic.claude-haiku-4-5-20251001-v1:0`

### Step 5: Draft Response

Match the response format to AWS's questions. See templates below.

---

## Actual Usage Data (as of 2026-01-13)

### Sonnet 4.5 V1 (7-day snapshot)
| Metric | Value |
|--------|-------|
| Total Invocations | 319 |
| Total Input Tokens | 8,796 |
| Total Output Tokens | 200,585 |
| Avg Input/Request | ~28 tokens (with caching) |
| Avg Output/Request | ~629 tokens |
| **Throttles** | **1,013** |

### Opus 4.5 (7-day snapshot)
| Metric | Value |
|--------|-------|
| Total Invocations | 1,681 |
| Total Input Tokens | 148,396 |
| Total Output Tokens | 1,081,636 |
| Avg Input/Request | ~88 tokens |
| Avg Output/Request | ~644 tokens |
| **Throttles** | **4,098** |

### Haiku 4.5 (7-day snapshot)
| Metric | Value |
|--------|-------|
| Total Invocations | 2,437 |
| Total Input Tokens | 6,712,680 |
| Avg Input/Request | ~2,754 tokens |
| Avg Output/Request | ~156 tokens |

---

## Response Templates

### Template A: Basic Questions (Simarpreet/Itish style)

When AWS asks:
1. AWS Account ID
2. Company Name
3. Region
4. Limit Type
5. Required TPM
6. Required RPM
7. Required TPD
8. Current TPM
9. Current RPM
10. Current TPD
11. Support Plan

**Response:**
```
1. AWS Account ID: 654654489230
2. Company Name: Legacy AI LLC
3. Region: us-east-1
4. Limit Type: On-demand ([Global] cross-region inference)
5. Required TPM: [FROM QUOTA LOOKUP OR PAIRED CASE]
6. Required RPM: [FROM QUOTA LOOKUP]
7. Required TPD: [RPM × 60 × 24 for requests, or TPM × 60 × 24 for tokens]
8. Current TPM: [FROM QUOTA LOOKUP]
9. Current RPM: [FROM QUOTA LOOKUP]
10. Current TPD: [CURRENT × 60 × 24]
11. Support Plan: Developer

Best regards,
Zeus Delacruz
Legacy AI LLC
```

### Template B: Detailed Questions (Deepak/Ahmed style)

When AWS asks:
1. Use case description
2. Limit type
3. Steady state TPM/RPM
4. Peak state TPM/RPM
5. Average input/output tokens
6. CRIS Option
7. Global CRIS Justification

**Response:**
```
1. Use Case Description:
   We operate an AI-native productivity platform (OneMind) using Claude models for:
   - Real-time AI assistance and conversational interfaces
   - Document analysis and knowledge management
   - Code generation and review
   - Multi-agent coordination

   Workloads run 8-12 hours daily with burst patterns. We are experiencing
   throttling ([X] throttles in the past 7 days).

2. Limit Type: On-demand

3. Steady State TPM: [~40-60% of requested]
4. Peak State TPM: [Requested value]
5. Steady State RPM: [~40-60% of requested]
6. Peak State RPM: [Requested value]

7. Average Input Tokens:
   - Standard models: 28-88 (with prompt caching)
   - 1M Context models: 50,000 (large document processing)

8. Average Output Tokens:
   - Standard models: 600-650
   - 1M Context models: 2,000

9. CRIS Option: [Global CRIS or CRIS based on quota type]

10. Global CRIS Justification:
    [If Global CRIS]: N/A - We use Global CRIS for high availability, automatic failover, and optimal latency routing.
    [If regular CRIS]: We use regional CRIS for [reason]. We have also submitted separate Global CRIS requests.

Best regards,
Zeus Delacruz
Legacy AI LLC
```

### Template C: Shreya Style (Model-specific)

When AWS asks:
1. Model ID
2. Input modalities
3. Steady/Peak TPM & RPM
4. Average input/output tokens
5. Detailed use case

**Response:**
```
1. Model ID: [us.anthropic.claude-MODEL-ID]

2. Input modalities: Text

3. Required Steady State TPM: [40-60% of requested]
4. Required Peak State TPM: [Requested value]
5. Required Steady State RPM: [40-60% of requested]
6. Required Peak State RPM: [Requested value]

7. Average Input Tokens per request: [FROM USAGE DATA]
8. Average Output Tokens per request: [FROM USAGE DATA]

9. Detailed use case:
   [Use standard use case description, tailored to model]

   For 1M Context: Mention large document processing (legal, technical, financial docs)
   For Opus: Mention complex reasoning, code generation, multi-step analysis
   For Sonnet: Mention real-time assistance, general purpose tasks

Best regards,
Zeus Delacruz
Legacy AI LLC
```

### Template D: Hannah Style (with MUs/Commitment)

When AWS asks about MUs (Model Units) and Commitment Term:
```
- Requested MUs (Model Units): N/A (On-demand, not provisioned throughput)
- Commitment Term: N/A (On-demand usage)
```

---

## Case Status Tracking

### Bedrock Quota Requests (Submitted 2026-01-13)

| Case ID | Quota | Current | Requested | Status |
|---------|-------|---------|-----------|--------|
| 176828500800661 | Sonnet 4.5 V1 - Global CRIS RPM | 200 | 1,000 | Responded |
| 176828500100038 | Opus 4.5 - CRIS TPM | 250,000 | 2,000,000 | Responded |
| 176828500200082 | Sonnet 4.5 1M - CRIS TPM | 1,000,000 | 5,000,000 | Responded |
| 176828500800217 | Sonnet 4.5 1M - Global CRIS RPM | 5 | 500 | Responded |
| 176828500600360 | Opus 4.5 - Global CRIS RPM | 250 | 1,000 | Responded |
| 176828500500407 | Sonnet 4.5 V1 - Global CRIS TPM | 500,000 | 5,000,000 | Responded |
| 176828500300161 | Sonnet 4.5 1M - CRIS RPM | 6 | 500 | Responded |
| 176828500400846 | Sonnet 4.5 1M - Global CRIS TPM | 1,000,000 | 5,000,000 | Responded |
| 176828500600210 | Opus 4.5 - Global CRIS TPM | 500,000 | 5,000,000 | Awaiting AWS |
| 176828500200323 | Sonnet 4.5 V1 - CRIS RPM | ? | 1,000 | Awaiting AWS |
| 176828500100816 | Opus 4.5 - CRIS RPM | 125 | 1,000 | Awaiting AWS |

### Lightsail Request (Submitted 2026-01-09)

| Case ID | Quota | Requested | Status |
|---------|-------|-----------|--------|
| 176799193000784 | Instances | 20 | Awaiting AWS |

---

## Important Notes

1. **Always verify quota values** before responding — don't rely on cached data
2. **Throttle data is powerful evidence** — mention it in every response
3. **CRIS vs Global CRIS** — Different quotas, make sure response matches the specific type
4. **1M Context models** — Use higher token estimates (50K input, 2K output) for these
5. **TPD Calculation** — TPD = (TPM or RPM) × 60 × 24
6. **Support API requires Business plan** — Cannot read/reply to cases programmatically with Developer plan

---

## Commands Quick Reference

```bash
# List all open cases
aws service-quotas list-requested-service-quota-change-history --status CASE_OPENED

# List approved cases
aws service-quotas list-requested-service-quota-change-history --status APPROVED

# Get specific quota current value
aws service-quotas get-service-quota --service-code bedrock --quota-code <CODE>

# Get account info
aws account get-contact-information

# Check AWS identity
aws sts get-caller-identity
```

---

## Session Continuity

When resuming work on AWS quota cases:

1. Read this document first
2. Check case status: `aws service-quotas list-requested-service-quota-change-history --status CASE_OPENED`
3. User provides Case ID + support email
4. Look up specific quota details
5. Check current usage data (especially throttles)
6. Draft response using appropriate template
7. Update this document with any new cases or status changes

---

*Last Updated: 2026-01-13 (quota limits refreshed)*

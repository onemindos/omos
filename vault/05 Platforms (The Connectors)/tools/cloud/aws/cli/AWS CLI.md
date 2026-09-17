---
title: "CLI — aws (AWS CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/aws, topic/cloud]
ai_summary: "AWS CLI quick reference — auth, S3, EC2, Lambda, IAM, and Secrets Manager for One Mind cloud resources."
---

# aws — AWS CLI

> Amazon Web Services cloud management. Install: `brew install awscli` / `pip install awscli`.

---

## Auth & Config

```bash
# Configure credentials (interactive)
aws configure

# Configure with profile
aws configure --profile onemind

# Use a specific profile
aws --profile onemind s3 ls

# Check current identity
aws sts get-caller-identity

# List configured profiles
aws configure list-profiles
```

Config stored at `~/.aws/credentials` and `~/.aws/config`.

---

## S3

```bash
# List buckets
aws s3 ls

# List contents of bucket
aws s3 ls s3://bucket-name/

# Upload file
aws s3 cp ./local-file.txt s3://bucket-name/path/file.txt

# Download file
aws s3 cp s3://bucket-name/path/file.txt ./local-file.txt

# Sync directory to S3
aws s3 sync ./local-dir s3://bucket-name/prefix/

# Sync S3 to local
aws s3 sync s3://bucket-name/prefix/ ./local-dir

# Remove file
aws s3 rm s3://bucket-name/path/file.txt

# Make bucket public-readable
aws s3api put-bucket-policy --bucket bucket-name --policy file://public-policy.json
```

---

## EC2

```bash
# List running instances
aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId,State.Name,PublicIpAddress,Tags[?Key==`Name`].Value]' --output table

# Start / stop instance
aws ec2 start-instances --instance-ids i-1234567890abcdef0
aws ec2 stop-instances --instance-ids i-1234567890abcdef0

# Get instance status
aws ec2 describe-instance-status --instance-ids i-1234567890abcdef0
```

---

## Lambda

```bash
# List functions
aws lambda list-functions --output table

# Invoke a function
aws lambda invoke \
  --function-name my-function \
  --payload '{"key": "value"}' \
  output.json

# Deploy updated code
aws lambda update-function-code \
  --function-name my-function \
  --zip-file fileb://function.zip

# View logs (recent)
aws logs tail /aws/lambda/my-function --follow
```

---

## Secrets Manager

```bash
# List secrets
aws secretsmanager list-secrets --output table

# Get a secret value
aws secretsmanager get-secret-value --secret-id my-secret --query SecretString --output text

# Create a secret
aws secretsmanager create-secret --name my-secret --secret-string "value"

# Update a secret
aws secretsmanager update-secret --secret-id my-secret --secret-string "new-value"
```

---

## IAM

```bash
# List users
aws iam list-users --output table

# List roles
aws iam list-roles --output table

# Get current user policies
aws iam list-attached-user-policies --user-name username
```

---

## Useful Flags

```bash
--output table         # human-readable table
--output json          # full JSON
--output text          # plain text (good for scripting)
--query "..."          # JMESPath filter
--region us-east-1     # override region
--profile onemind      # use named profile
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

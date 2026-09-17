# Amazon AWS

> **Cloud Infrastructure Platform**

## Overview

Amazon Web Services provides the cloud infrastructure for OneMind OS, including AI model hosting (Bedrock), compute, storage, and serverless functions.

## Status

- **Active**: ✅ Yes
- **Primary Region**: us-east-1
- **Last Updated**: January 25, 2026

## Services in Use

### AI/ML

| Service       | Purpose                 | Status     |
| ------------- | ----------------------- | ---------- |
| **Bedrock**   | Claude Opus 4.5 hosting | ✅ Active  |
| **SageMaker** | Model training (future) | 🔄 Planned |

### Compute

| Service    | Purpose                 | Status    |
| ---------- | ----------------------- | --------- |
| **EC2**    | Agent server hosting    | ✅ Active |
| **Lambda** | Serverless functions    | ✅ Active |
| **ECS**    | Container orchestration | ✅ Active |

### Storage

| Service | Purpose               | Status     |
| ------- | --------------------- | ---------- |
| **S3**  | File storage, backups | ✅ Active  |
| **EFS** | Shared file system    | 🔄 Planned |

### Networking

| Service        | Purpose           | Status     |
| -------------- | ----------------- | ---------- |
| **VPC**        | Network isolation | ✅ Active  |
| **Route 53**   | DNS management    | ✅ Active  |
| **CloudFront** | CDN               | 🔄 Planned |

### Security

| Service             | Purpose           | Status    |
| ------------------- | ----------------- | --------- |
| **IAM**             | Access management | ✅ Active |
| **Secrets Manager** | Secret storage    | ✅ Active |
| **KMS**             | Encryption keys   | ✅ Active |

## AWS Bedrock (Claude)

### Setup

```python
import boto3

bedrock = boto3.client(
    service_name='bedrock-runtime',
    region_name='us-east-1'
)
```

### Invoke Claude

```python
import json

def invoke_claude(prompt: str, max_tokens: int = 1024):
    body = json.dumps({
        "anthropic_version": "bedrock-2023-05-31",
        "max_tokens": max_tokens,
        "messages": [
            {"role": "user", "content": prompt}
        ]
    })

    response = bedrock.invoke_model(
        modelId="anthropic.claude-3-opus-20240229",
        body=body
    )

    result = json.loads(response['body'].read())
    return result['content'][0]['text']
```

### Streaming Response

```python
def stream_claude(prompt: str):
    body = json.dumps({
        "anthropic_version": "bedrock-2023-05-31",
        "max_tokens": 1024,
        "messages": [{"role": "user", "content": prompt}]
    })

    response = bedrock.invoke_model_with_response_stream(
        modelId="anthropic.claude-3-opus-20240229",
        body=body
    )

    for event in response['body']:
        chunk = json.loads(event['chunk']['bytes'])
        if chunk['type'] == 'content_block_delta':
            yield chunk['delta'].get('text', '')
```

## S3 Storage

### Python SDK

```python
import boto3

s3 = boto3.client('s3')

# Upload file
s3.upload_file('local_file.txt', 'bucket-name', 'path/to/file.txt')

# Download file
s3.download_file('bucket-name', 'path/to/file.txt', 'local_file.txt')

# List objects
response = s3.list_objects_v2(Bucket='bucket-name', Prefix='folder/')
for obj in response.get('Contents', []):
    print(obj['Key'])

# Generate presigned URL
url = s3.generate_presigned_url(
    'get_object',
    Params={'Bucket': 'bucket-name', 'Key': 'file.txt'},
    ExpiresIn=3600
)
```

## Lambda Functions

### Handler Example

```python
# lambda_function.py
import json

def handler(event, context):
    body = json.loads(event.get('body', '{}'))

    # Process request
    result = process(body)

    return {
        'statusCode': 200,
        'headers': {'Content-Type': 'application/json'},
        'body': json.dumps(result)
    }
```

### Deployment (SAM)

```yaml
# template.yaml
AWSTemplateFormatVersion: "2010-09-09"
Transform: AWS::Serverless-2016-10-31

Resources:
  ProcessFunction:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: python3.12
      Handler: lambda_function.handler
      MemorySize: 256
      Timeout: 30
      Events:
        Api:
          Type: Api
          Properties:
            Path: /process
            Method: POST
```

## ECS Deployment

### Task Definition

```json
{
  "family": "legacy-agent",
  "containerDefinitions": [
    {
      "name": "legacy-agent",
      "image": "your-ecr-repo/legacy-agent:latest",
      "portMappings": [{ "containerPort": 8080, "protocol": "tcp" }],
      "environment": [{ "name": "NATS_URL", "value": "nats://nats:4222" }],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:secretsmanager:..."
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/legacy-agent",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

## IAM Configuration

### Bedrock Access Policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "bedrock:InvokeModel",
        "bedrock:InvokeModelWithResponseStream"
      ],
      "Resource": ["arn:aws:bedrock:*::foundation-model/anthropic.claude-*"]
    }
  ]
}
```

## Secrets Manager

### Store Secret

```python
import boto3
import json

secrets = boto3.client('secretsmanager')

secrets.create_secret(
    Name='onemind/api-keys',
    SecretString=json.dumps({
        'ANTHROPIC_API_KEY': 'sk-...',
        'NATS_PASSWORD': '...'
    })
)
```

### Retrieve Secret

```python
def get_secret(secret_name: str) -> dict:
    response = secrets.get_secret_value(SecretId=secret_name)
    return json.loads(response['SecretString'])

# Usage
api_keys = get_secret('onemind/api-keys')
```

## CloudWatch Monitoring

### Custom Metrics

```python
import boto3

cloudwatch = boto3.client('cloudwatch')

def put_metric(name: str, value: float, unit: str = 'Count'):
    cloudwatch.put_metric_data(
        Namespace='OneMind',
        MetricData=[{
            'MetricName': name,
            'Value': value,
            'Unit': unit
        }]
    )
```

### Log Insights Query

```
fields @timestamp, @message
| filter @message like /ERROR/
| sort @timestamp desc
| limit 20
```

## Cost Optimization

### EC2 Savings Plans

- Use Reserved Instances for predictable workloads
- Use Spot Instances for batch processing

### S3 Storage Classes

| Class               | Use Case               |
| ------------------- | ---------------------- |
| Standard            | Frequent access        |
| Intelligent-Tiering | Unknown access pattern |
| Glacier             | Long-term archive      |

## Environment Variables

```env
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_DEFAULT_REGION=us-east-1
```

## Best Practices

- ✅ Use IAM roles instead of access keys when possible
- ✅ Enable CloudTrail for audit logging
- ✅ Use VPC endpoints for private access
- ✅ Enable encryption at rest and in transit
- ✅ Set up billing alerts
- ✅ Use tags for resource organization
- ✅ Follow least-privilege principle

## Resources

- [AWS Documentation](https://docs.aws.amazon.com/)
- [Bedrock User Guide](https://docs.aws.amazon.com/bedrock/)
- [Boto3 Documentation](https://boto3.amazonaws.com/v1/documentation/api/latest/)

---

**Related**: [Azure](Azure.md) | [Agno](./Agno.md) | [Docker](./Docker.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]

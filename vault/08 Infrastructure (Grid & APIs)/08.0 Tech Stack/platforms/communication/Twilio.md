# Twilio

> **Communication APIs (SMS, Voice, Video)**

## Overview

Twilio provides cloud communication APIs for SMS, voice calls, and video. It will power Legacy AI's multi-channel communication capabilities.

## Status

- **Active**: 🔄 Planned
- **Services**: SMS, Voice, WhatsApp
- **Last Updated**: January 25, 2026

## Use Cases

1. **SMS Notifications** - Send alerts via text
2. **Voice Calls** - Automated calls
3. **WhatsApp** - Business messaging
4. **Two-Factor Auth** - Verification codes
5. **Voice AI** - Phone-based Legacy interface

## Key Features

- Programmable messaging
- Programmable voice
- WhatsApp Business API
- Phone number provisioning
- Call recording
- IVR systems

## Installation

```bash
pip install twilio
```

## Configuration

### Environment Variables

```env
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_PHONE_NUMBER=+1234567890
```

### Client Setup

```python
from twilio.rest import Client
import os

client = Client(
    os.getenv("TWILIO_ACCOUNT_SID"),
    os.getenv("TWILIO_AUTH_TOKEN")
)
```

## SMS Messaging

### Send SMS

```python
def send_sms(to: str, body: str) -> str:
    message = client.messages.create(
        body=body,
        from_=os.getenv("TWILIO_PHONE_NUMBER"),
        to=to
    )
    return message.sid

# Usage
send_sms("+1234567890", "Hello from Legacy AI!")
```

### Receive SMS (Webhook)

```python
from fastapi import FastAPI, Form
from twilio.twiml.messaging_response import MessagingResponse

app = FastAPI()

@app.post("/webhooks/sms")
async def handle_sms(From: str = Form(...), Body: str = Form(...)):
    # Process incoming message
    response_text = await process_message(From, Body)

    # Reply
    resp = MessagingResponse()
    resp.message(response_text)

    return Response(content=str(resp), media_type="application/xml")
```

## Voice Calls

### Make a Call

```python
def make_call(to: str, message: str) -> str:
    call = client.calls.create(
        twiml=f'<Response><Say>{message}</Say></Response>',
        to=to,
        from_=os.getenv("TWILIO_PHONE_NUMBER")
    )
    return call.sid
```

### Handle Incoming Call

```python
from twilio.twiml.voice_response import VoiceResponse, Gather

@app.post("/webhooks/voice")
async def handle_voice():
    resp = VoiceResponse()

    gather = Gather(
        input='speech',
        action='/webhooks/voice/process',
        speechTimeout='auto',
        language='en-US'
    )
    gather.say("Hello, this is Legacy. How can I help you?")
    resp.append(gather)

    return Response(content=str(resp), media_type="application/xml")

@app.post("/webhooks/voice/process")
async def process_voice(SpeechResult: str = Form(...)):
    # Process speech and get response
    response_text = await process_with_legacy(SpeechResult)

    resp = VoiceResponse()
    resp.say(response_text)
    resp.redirect('/webhooks/voice')  # Continue conversation

    return Response(content=str(resp), media_type="application/xml")
```

## WhatsApp

### Send WhatsApp Message

```python
def send_whatsapp(to: str, body: str) -> str:
    message = client.messages.create(
        body=body,
        from_='whatsapp:' + os.getenv("TWILIO_WHATSAPP_NUMBER"),
        to='whatsapp:' + to
    )
    return message.sid
```

### Receive WhatsApp (Webhook)

```python
@app.post("/webhooks/whatsapp")
async def handle_whatsapp(From: str = Form(...), Body: str = Form(...)):
    # Remove whatsapp: prefix
    phone = From.replace("whatsapp:", "")

    # Process message
    response_text = await process_message(phone, Body)

    resp = MessagingResponse()
    resp.message(response_text)

    return Response(content=str(resp), media_type="application/xml")
```

## Verification (2FA)

### Send Verification Code

```python
verify = client.verify.v2.services(os.getenv("TWILIO_VERIFY_SID"))

def send_verification(phone: str) -> str:
    verification = verify.verifications.create(
        to=phone,
        channel='sms'  # or 'call', 'email', 'whatsapp'
    )
    return verification.status

def check_verification(phone: str, code: str) -> bool:
    check = verify.verification_checks.create(
        to=phone,
        code=code
    )
    return check.status == 'approved'
```

## TwiML Reference

### Voice

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say voice="alice">Hello!</Say>
    <Play>https://example.com/audio.mp3</Play>
    <Gather input="speech dtmf" timeout="5">
        <Say>Press 1 or say yes to continue</Say>
    </Gather>
    <Dial>+1234567890</Dial>
    <Record maxLength="30" transcribe="true" />
</Response>
```

### Messaging

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Message>Hello from Legacy!</Message>
    <Message>
        <Media>https://example.com/image.jpg</Media>
        Photo attached
    </Message>
</Response>
```

## Webhook Configuration

### In Twilio Console

1. Go to Phone Numbers
2. Select your number
3. Configure webhook URLs:
   - **SMS incoming**: `https://api.onemind.io/webhooks/sms`
   - **Voice incoming**: `https://api.onemind.io/webhooks/voice`

### Webhook Validation

```python
from twilio.request_validator import RequestValidator

def validate_twilio_request(request_url: str, post_data: dict, signature: str) -> bool:
    validator = RequestValidator(os.getenv("TWILIO_AUTH_TOKEN"))
    return validator.validate(request_url, post_data, signature)
```

## Legacy AI Integration

### SMS to Agent

```python
@app.post("/webhooks/sms")
async def handle_sms(From: str = Form(...), Body: str = Form(...)):
    # Get or create session for this phone number
    session = await get_or_create_session(From)

    # Send to Legacy agent
    response = await legacy_agent.run(
        message=Body,
        session_id=session.id
    )

    # Reply via SMS
    resp = MessagingResponse()
    resp.message(response.content)

    return Response(content=str(resp), media_type="application/xml")
```

### Proactive Notifications

```python
async def notify_user(user_id: str, message: str, channel: str = "sms"):
    user = await get_user(user_id)

    if channel == "sms":
        send_sms(user.phone, message)
    elif channel == "whatsapp":
        send_whatsapp(user.phone, message)
    elif channel == "call":
        make_call(user.phone, message)
```

## Pricing

| Service          | Cost                 |
| ---------------- | -------------------- |
| SMS (outbound)   | ~$0.0075/message     |
| SMS (inbound)    | ~$0.0075/message     |
| Voice (outbound) | ~$0.0140/min         |
| Voice (inbound)  | ~$0.0085/min         |
| WhatsApp         | ~$0.005-0.08/message |
| Phone number     | ~$1.15/month         |

## Best Practices

- ✅ Validate webhook requests
- ✅ Handle errors gracefully in TwiML
- ✅ Use async for webhook handlers
- ✅ Log all messages for debugging
- ✅ Implement rate limiting
- ✅ Store conversation history
- ✅ Handle opt-outs (STOP)

## Resources

- [Twilio Documentation](https://www.twilio.com/docs)
- [TwiML Reference](https://www.twilio.com/docs/voice/twiml)
- [Python Quickstart](https://www.twilio.com/docs/sms/quickstart/python)

---

**Related**: [LiveKit](./LiveKit.md) | [N8N](./N8N.md)

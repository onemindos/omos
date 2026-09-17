# Gemini Voice System

> Inline voice recording with Gemini Text-to-Speech for natural conversational AI.

## Overview

The OneMind voice system provides seamless voice interaction within the chat interface. Users can record voice messages, receive spoken responses, and combine audio with images for multimodal queries—all without leaving the chat screen.

**Key Technologies:**
- **Gemini 2.5 Flash TTS** - Text-to-speech synthesis with 5 natural voices
- **Deepgram** - Speech-to-text transcription
- **google-genai SDK** - Google's official Python SDK for Gemini

## Features

| Feature | Description |
|---------|-------------|
| **Inline Recording** | Voice overlay replaces input bar (no navigation away) |
| **5 Gemini Voices** | Kore, Charon, Aoede, Fenrir, Puck with persistent selection |
| **Voice Responses** | Audio playback embedded in message bubbles |
| **TTS for Any Message** | Speak any assistant message with selected voice |
| **Multimodal Queries** | Combine audio + image using Gemini vision |
| **Waveform Visualization** | Animated amplitude display during recording |

## Available Voices

| Voice ID | Name | Description | Gender |
|----------|------|-------------|--------|
| `kore` | Kore | Warm and friendly | Female |
| `charon` | Charon | Deep and calm | Male |
| `aoede` | Aoede | Expressive and clear | Female |
| `fenrir` | Fenrir | Confident and strong | Male |
| `puck` | Puck | Friendly and playful | Neutral |

**Default:** `kore` (Female, warm and friendly)

## Architecture

### Data Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                           FRONTEND                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. User taps mic button in ChatInputBar                            │
│                    ↓                                                 │
│  2. InlineVoiceOverlay appears (replaces input bar)                 │
│     - Waveform visualization                                         │
│     - Camera button for multimodal                                   │
│     - Cancel / Send buttons                                          │
│                    ↓                                                 │
│  3. User records audio (optional: attaches image)                   │
│                    ↓                                                 │
│  4. User taps send → VoiceProvider.stopAndSend()                    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ↓
┌─────────────────────────────────────────────────────────────────────┐
│                            BACKEND                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  POST /api/voice/message                                            │
│  {                                                                   │
│    "audio_base64": "...",                                           │
│    "voice": "kore",                                                 │
│    "agent_id": "legacy-mother-ai",                                  │
│    "session_id": "optional"                                         │
│  }                                                                   │
│                    ↓                                                 │
│  1. Deepgram STT: audio → transcript                                │
│                    ↓                                                 │
│  2. Agent processes transcript → response text                      │
│                    ↓                                                 │
│  3. Gemini TTS: response text → audio bytes                         │
│                    ↓                                                 │
│  Response:                                                           │
│  {                                                                   │
│    "transcript": "What's the weather?",                             │
│    "response_text": "It's 72°F and sunny...",                       │
│    "audio_base64": "...",                                           │
│    "voice": "kore"                                                  │
│  }                                                                   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ↓
┌─────────────────────────────────────────────────────────────────────┐
│                           FRONTEND                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  5. Add user message to chat (transcript, isVoiceInput: true)       │
│                    ↓                                                 │
│  6. Add assistant message (response + audioBytes)                   │
│                    ↓                                                 │
│  7. Auto-play audio response                                         │
│                    ↓                                                 │
│  8. VoiceAudioPlayer shows in message bubble                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Multimodal Flow (Audio + Image)

```
User records audio + attaches image
           ↓
POST /api/voice/multimodal
{
  "audio_base64": "...",
  "image_base64": "...",
  "voice": "kore",
  "agent_id": "vision-agent"
}
           ↓
Gemini processes audio + image together natively
           ↓
Response with transcript, analysis, and spoken response
```

## API Reference

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/voice/health` | Health check with provider status |
| `POST` | `/api/voice/transcribe` | Audio → text (Deepgram STT) |
| `POST` | `/api/voice/synthesize` | Text → audio (Gemini TTS) |
| `POST` | `/api/voice/message` | Full voice conversation flow |
| `POST` | `/api/voice/multimodal` | Audio + image query |
| `GET` | `/api/voice/sessions` | List voice sessions |
| `DELETE` | `/api/voice/sessions/{id}` | Delete voice session |

### Health Check

```bash
curl http://localhost:8080/api/voice/health
```

**Response:**
```json
{
  "status": "healthy",
  "gemini_enabled": true,
  "deepgram_enabled": true,
  "available_voices": ["kore", "charon", "aoede", "fenrir", "puck"],
  "tts_model": "gemini-2.5-flash-preview-tts"
}
```

### Text-to-Speech

```bash
curl -X POST http://localhost:8080/api/voice/synthesize \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, how can I help you today?",
    "voice": "kore"
  }'
```

**Response:**
```json
{
  "audio_base64": "UklGRi...",
  "format": "wav",
  "sample_rate": 24000,
  "voice": "kore"
}
```

### Full Voice Message

```bash
curl -X POST http://localhost:8080/api/voice/message \
  -H "Content-Type: application/json" \
  -d '{
    "audio_base64": "...",
    "voice": "charon",
    "agent_id": "legacy-mother-ai",
    "session_id": "optional-session-id"
  }'
```

**Response:**
```json
{
  "transcript": "What time is my next meeting?",
  "response_text": "Your next meeting is at 3:00 PM with the design team.",
  "audio_base64": "UklGRi...",
  "voice": "charon",
  "session_id": "sess_abc123"
}
```

## Frontend Implementation

### File Structure

```
frontend/lib/
├── agno/
│   ├── chat/
│   │   ├── screens/
│   │   │   └── enhanced_chat_screen.dart    # Voice mode state
│   │   └── widgets/
│   │       ├── chat_input_bar.dart          # Voice chip, mic button
│   │       ├── enhanced_message_bubble.dart # VoiceAudioPlayer, TTS button
│   │       └── inline_voice_overlay.dart    # Recording overlay
│   │
│   └── voice/
│       └── voice_provider.dart              # Recording state management
│
├── platform/
│   └── providers/
│       └── app_providers.dart               # selectedVoiceProvider, geminiVoices
│
└── shared/
    └── models/
        └── enhanced_message.dart            # VoiceMetadata model
```

### Key Components

#### 1. Voice Selection Provider

```dart
// frontend/lib/platform/providers/app_providers.dart

@riverpod
class SelectedVoice extends _$SelectedVoice {
  @override
  String build() {
    _loadFromPrefs();
    return 'kore'; // Default
  }

  Future<void> setVoice(String voice) async {
    state = voice;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('onemind_selected_voice', voice);
  }
}

const List<VoiceOption> geminiVoices = [
  VoiceOption(id: 'kore', name: 'Kore', description: 'Warm and friendly', gender: 'Female'),
  VoiceOption(id: 'charon', name: 'Charon', description: 'Deep and calm', gender: 'Male'),
  VoiceOption(id: 'aoede', name: 'Aoede', description: 'Expressive and clear', gender: 'Female'),
  VoiceOption(id: 'fenrir', name: 'Fenrir', description: 'Confident and strong', gender: 'Male'),
  VoiceOption(id: 'puck', name: 'Puck', description: 'Friendly and playful', gender: 'Neutral'),
];
```

#### 2. Inline Voice Overlay

```dart
// frontend/lib/agno/chat/widgets/inline_voice_overlay.dart

class InlineVoiceOverlay extends ConsumerStatefulWidget {
  final VoidCallback onCancel;
  final Future<void> Function(String transcript, String response, Uint8List? audio) onVoiceResponse;
  final String agentId;
  final String? sessionId;
  final String selectedVoice;

  // Features:
  // - Animated waveform visualization
  // - Camera button for multimodal queries
  // - Cancel and Send buttons
  // - Image preview when photo attached
  // - Processing status indicator
}
```

#### 3. Voice Metadata Model

```dart
// frontend/lib/shared/models/enhanced_message.dart

@freezed
class VoiceMetadata with _$VoiceMetadata {
  const factory VoiceMetadata({
    @Default(false) bool isVoiceInput,      // User spoke this message
    @Default(false) bool isVoiceOutput,     // Has audio response
    String? audioPath,                       // Local audio file path
    List<int>? audioBytes,                   // Audio data for playback
    String? audioFormat,                     // wav, mp3, etc.
    int? durationMs,                         // Audio duration
    String? voiceId,                         // Voice used (kore, etc.)
    String? ttsProvider,                     // gemini, cartesia, etc.
    @Default(false) bool hasBeenSpoken,     // Has been read aloud
    @Default(false) bool isPlaying,         // Currently playing
  }) = _VoiceMetadata;
}
```

#### 4. Voice Audio Player Widget

```dart
// frontend/lib/agno/chat/widgets/enhanced_message_bubble.dart

class VoiceAudioPlayer extends StatefulWidget {
  final VoiceMetadata voice;
  final String messageId;
  final VoiceAudioCallback? onAction;

  // Features:
  // - Play/pause button with pulse animation
  // - Progress bar
  // - Duration display
  // - Voice name indicator
}
```

## Backend Implementation

### File Structure

```
backend/agno/voice/
├── __init__.py
├── api.py          # FastAPI router with endpoints
├── models.py       # Request/response models
└── providers.py    # GeminiTTS, DeepgramSTT providers
```

### Gemini TTS Provider

```python
# backend/agno/voice/providers.py

class GeminiTTS(TTSProvider):
    """Google Gemini TTS using gemini-2.5-flash-preview-tts model"""

    def __init__(self):
        self.api_key = os.getenv("GOOGLE_API_KEY")
        self.model = "gemini-2.5-flash-preview-tts"

    async def synthesize(
        self,
        text: str,
        voice: str = "kore",
        speed: float = 1.0
    ) -> SynthesisResult:
        from google import genai
        from google.genai import types

        os.environ["GOOGLE_API_KEY"] = self.api_key
        client = genai.Client()

        # Map voice ID to Gemini voice name
        voice_name = voice.capitalize()  # kore -> Kore

        response = client.models.generate_content(
            model=self.model,
            contents=f"Say the following: {text}",
            config=types.GenerateContentConfig(
                response_modalities=["AUDIO"],
                speech_config=types.SpeechConfig(
                    voice_config=types.VoiceConfig(
                        prebuilt_voice_config=types.PrebuiltVoiceConfig(
                            voice_name=voice_name,
                        ),
                    ),
                ),
            ),
        )

        audio_data = response.candidates[0].content.parts[0].inline_data.data

        return SynthesisResult(
            audio_data=audio_data,
            format="wav",
            sample_rate=24000,
        )
```

### Voice API Router

```python
# backend/agno/voice/api.py

router = APIRouter(prefix="/api/voice", tags=["voice"])

@router.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "gemini_enabled": bool(os.getenv("GOOGLE_API_KEY")),
        "deepgram_enabled": bool(os.getenv("DEEPGRAM_API_KEY")),
        "available_voices": ["kore", "charon", "aoede", "fenrir", "puck"],
    }

@router.post("/synthesize")
async def synthesize(request: SynthesizeRequest):
    tts = GeminiTTS()
    result = await tts.synthesize(
        text=request.text,
        voice=request.voice,
    )
    return {
        "audio_base64": base64.b64encode(result.audio_data).decode(),
        "format": result.format,
        "voice": request.voice,
    }

@router.post("/message")
async def voice_message(request: VoiceMessageRequest):
    # 1. Transcribe audio with Deepgram
    stt = DeepgramSTT()
    transcript = await stt.transcribe(request.audio_base64)

    # 2. Process with agent
    response_text = await process_with_agent(
        transcript=transcript,
        agent_id=request.agent_id,
        session_id=request.session_id,
    )

    # 3. Synthesize response with Gemini TTS
    tts = GeminiTTS()
    audio_result = await tts.synthesize(
        text=response_text,
        voice=request.voice,
    )

    return {
        "transcript": transcript,
        "response_text": response_text,
        "audio_base64": base64.b64encode(audio_result.audio_data).decode(),
        "voice": request.voice,
    }
```

## Configuration

### Environment Variables

```bash
# .env

# Required for Gemini TTS
GOOGLE_API_KEY=your-google-api-key

# Required for Deepgram STT
DEEPGRAM_API_KEY=your-deepgram-api-key
```

### Voice Settings (Flutter)

Voice selection persists to SharedPreferences:
- Key: `onemind_selected_voice`
- Default: `kore`

## Testing

### Backend Health Check

```bash
curl http://localhost:8080/api/voice/health
```

Expected response:
```json
{
  "status": "healthy",
  "gemini_enabled": true,
  "deepgram_enabled": true
}
```

### TTS Test

```bash
curl -X POST http://localhost:8080/api/voice/synthesize \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello world", "voice": "kore"}' \
  | jq -r '.audio_base64' | base64 -d > test.wav

# Play the audio
afplay test.wav  # macOS
```

### Full Conversation Test

```bash
# Record audio first, then:
curl -X POST http://localhost:8080/api/voice/message \
  -H "Content-Type: application/json" \
  -d '{
    "audio_base64": "'$(base64 < recording.wav)'",
    "voice": "charon",
    "agent_id": "legacy-mother-ai"
  }'
```

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| `gemini_enabled: false` | Missing GOOGLE_API_KEY | Add key to .env on server |
| "Model tried to generate text" | Wrong API usage | Use `google-genai` SDK with correct format |
| No audio playback | Missing audioBytes | Check VoiceMetadata in message |
| Voice selection not persisting | SharedPreferences issue | Check app permissions |

### Verifying Gemini TTS

```python
# Quick test script
from google import genai
from google.genai import types
import os

os.environ["GOOGLE_API_KEY"] = "your-key"
client = genai.Client()

response = client.models.generate_content(
    model="gemini-2.5-flash-preview-tts",
    contents="Say: Hello world",
    config=types.GenerateContentConfig(
        response_modalities=["AUDIO"],
        speech_config=types.SpeechConfig(
            voice_config=types.VoiceConfig(
                prebuilt_voice_config=types.PrebuiltVoiceConfig(
                    voice_name="Kore",
                ),
            ),
        ),
    ),
)

audio = response.candidates[0].content.parts[0].inline_data.data
with open("test.wav", "wb") as f:
    f.write(audio)
```

## Related Documentation

- [CLAUDE.md](../CLAUDE.md) - Main project documentation
- [Gemini TTS Documentation](https://ai.google.dev/gemini-api/docs/text-to-speech)
- [Deepgram Documentation](https://developers.deepgram.com/docs)
- [google-genai SDK](https://pypi.org/project/google-genai/)

## Changelog

### 2025-02

- Initial implementation with Gemini TTS
- Added inline voice recording overlay
- Voice selection with 5 Gemini voices
- Multimodal support (audio + image)
- TTS button for assistant messages
- VoiceAudioPlayer widget for playback
- Fixed `google-genai` SDK integration (replaced deprecated `google-generativeai`)

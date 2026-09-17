# GeminiLiveTalkProvider — Build Doc

Target file: `command-center/Sources/Voice/GeminiLiveTalkProvider.swift`

Replaces / extends `TalkModeManager.swift` (currently uses ElevenLabs TTS + STT round-trip).
This adds true full-duplex audio over a single persistent WebSocket — no STT→LLM→TTS latency gap.

---

## What We're Building

A Swift class that wraps the Gemini Live (or OpenAI Realtime) WebSocket API, giving Legacy
true real-time bidirectional voice on the Command Center. Legacy hears Zeus's mic audio live
and streams audio back simultaneously.

---

## Reference

`s2s.py` in this folder is the working Python prototype (`jarvis_v3` from Concept-Bytes).
It uses the **OpenAI Realtime API** (`wss://api.openai.com/v1/realtime`).
Gemini Live uses a similar WebSocket protocol — swap the endpoint + auth header.

---

## Key Patterns from `s2s.py` to Port

### 1. Mute mic during AI audio playback (echo prevention)

```python
# Python: s2s.py
if event_type == "response.audio.delta":
    if not self.mute_mic:
        self.mute_mic = True           # stop sending mic audio
    audio_chunk = base64.b64decode(event["delta"])
    self.output_stream.write(audio_chunk)

elif event_type == "response.audio.done":
    def delayed_unmute():
        time.sleep(0.5)                # 500ms — wait for room echo to decay
        self.mute_mic = False
    threading.Thread(target=delayed_unmute, daemon=True).start()
```

```swift
// Swift equivalent (GeminiLiveTalkProvider.swift)
private var isMuted = false

func onAudioDelta(_ chunk: Data) {
    isMuted = true
    audioPlayer.enqueue(chunk)         // AVAudioPlayerNode or AudioUnit
}

func onAudioDone() {
    DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) { [weak self] in
        self?.isMuted = false
    }
}

// In mic capture callback:
func captureOutput(_ chunk: Data) {
    guard !isMuted else { return }
    webSocket.send(encodeAudioEvent(chunk))
}
```

### 2. WebSocket session setup

```python
# Python: s2s.py session.update payload
session_update = {
    "type": "session.update",
    "session": {
        "voice": self.voice,
        "output_audio_format": "pcm16",
        "instructions": self.initial_prompt,
    }
}
```

Swift: build the same JSON dict, send on `webSocketTask.send(.string(json))`.
For Gemini Live the session init event differs slightly — see Gemini Live WS docs.

### 3. ask_sync pattern (text queries over the same socket)

```python
# s2s.py: blocks until response arrives using queue.Queue
def ask_sync(self, message: str, timeout: int = 30) -> str:
    response_queue = queue.Queue()
    self.on_text_response = lambda text: response_queue.put(text)
    self.send_text_message(message)
    return response_queue.get(timeout=timeout)
```

Swift equivalent: use `CheckedContinuation` / `async/await`:

```swift
func ask(_ message: String) async throws -> String {
    return try await withCheckedThrowingContinuation { continuation in
        self.pendingContinuation = continuation
        sendTextMessage(message)
    }
}
```

### 4. Audio format

Both APIs use **PCM16 mono 24kHz**. In Swift:

```swift
let format = AVAudioFormat(
    commonFormat: .pcmFormatInt16,
    sampleRate: 24000,
    channels: 1,
    interleaved: true
)
```

---

## Swift File Skeleton

```swift
// command-center/Sources/Voice/GeminiLiveTalkProvider.swift

import Foundation
import AVFoundation

final class GeminiLiveTalkProvider: NSObject {

    // MARK: - Config
    private let apiKey: String
    private let model: String        // "gemini-2.0-flash-live-001" or openai realtime
    private let systemPrompt: String

    // MARK: - WebSocket
    private var webSocketTask: URLSessionWebSocketTask?
    private var session: URLSession?

    // MARK: - Audio
    private let audioEngine = AVAudioEngine()
    private var playerNode = AVAudioPlayerNode()
    private var isMuted = false

    // MARK: - Pending async continuations
    private var pendingContinuation: CheckedContinuation<String, Error>?

    init(apiKey: String, model: String, systemPrompt: String) {
        self.apiKey = apiKey
        self.model = model
        self.systemPrompt = systemPrompt
    }

    func connect() {
        // Build WSS URL + headers (Gemini or OpenAI)
        // Send session.update on open
        // Start mic capture loop
        // Start receive loop
    }

    func disconnect() { webSocketTask?.cancel(with: .normalClosure, reason: nil) }

    // Called by NodeCapabilityRouter when gateway sends "voice.realtime.start"
    func startTalkMode() { connect() }
    func stopTalkMode()  { disconnect() }

    func ask(_ message: String) async throws -> String { /* ask_sync pattern */ }

    private func sendAudioChunk(_ pcm16: Data) {
        guard !isMuted else { return }
        // Base64-encode + send as input_audio_buffer.append event
    }

    private func onAudioDelta(_ b64: String) {
        isMuted = true
        guard let data = Data(base64Encoded: b64) else { return }
        // Schedule on playerNode
    }

    private func onAudioDone() {
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) { self.isMuted = false }
    }

    private func receiveLoop() {
        webSocketTask?.receive { [weak self] result in
            switch result {
            case .success(let msg):
                self?.handleMessage(msg)
                self?.receiveLoop()   // recurse
            case .failure(let err):
                print("WS error:", err)
            }
        }
    }

    private func handleMessage(_ msg: URLSessionWebSocketTask.Message) {
        guard case .string(let text) = msg,
              let data = text.data(using: .utf8),
              let event = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              let type = event["type"] as? String else { return }

        switch type {
        case "response.audio.delta": onAudioDelta(event["delta"] as? String ?? "")
        case "response.audio.done":  onAudioDone()
        case "response.text.done":
            let txt = event["text"] as? String ?? ""
            pendingContinuation?.resume(returning: txt)
            pendingContinuation = nil
        default: break
        }
    }
}
```

---

## OpenClaw Integration

When Legacy wants to start a real-time voice call:

```
Gateway → command-center → NodeCapabilityRouter → GeminiLiveTalkProvider.startTalkMode()
```

NodeCapabilityRouter already handles `voice.*` subjects. Add a case:

```swift
case "voice.realtime.start":
    geminiLive.startTalkMode()
case "voice.realtime.stop":
    geminiLive.stopTalkMode()
```

---

## API Endpoints

| Provider | WebSocket URL |
|---|---|
| OpenAI Realtime | `wss://api.openai.com/v1/realtime?model=gpt-4o-realtime-preview` |
| Gemini Live | `wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1beta.GenerativeService.BidiGenerateContent` |

OpenAI auth: `Authorization: Bearer {OPENAI_API_KEY}` + `OpenAI-Beta: realtime=v1`
Gemini auth: `?key={GEMINI_API_KEY}` query param

---

## Build Order

1. `GeminiLiveTalkProvider.swift` — core WS class (this doc)
2. Wire into `NodeCapabilityRouter.swift` — add `voice.realtime.*` cases
3. Update `TalkModeManager.swift` — make it delegate to GeminiLiveTalkProvider when mode = realtime
4. Add `GEMINI_API_KEY` / `OPENAI_REALTIME_KEY` to `command-center/Secrets.xcconfig`

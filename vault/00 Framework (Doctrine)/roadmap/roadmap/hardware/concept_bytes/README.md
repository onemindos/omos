# ConceptBytes Extraction

## 🔗 Sources

| Project            | Repository                                                                     | Description                                            |
| ------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------ |
| **JarvisV3**       | [github.com/Concept-Bytes/JarvisV3](https://github.com/Concept-Bytes/JarvisV3) | OpenAI Realtime API voice assistant with Streamlit UI  |
| **Jarvis (v1/v2)** | [github.com/Concept-Bytes/Jarvis](https://github.com/Concept-Bytes/Jarvis)     | Whisper + GPT + TTS with hotword detection & Spotify   |
| **HoloMat2**       | [github.com/Concept-Bytes/HoloMat2](https://github.com/Concept-Bytes/HoloMat2) | Pygame radial touch UI for transparent displays        |
| **ai_plant**       | [github.com/Concept-Bytes/ai_plant](https://github.com/Concept-Bytes/ai_plant) | Raspberry Pi + Adafruit STEMMA soil sensor + ChatGPT   |
| **Patreon**        | [patreon.com/cw/ConceptBytes](https://www.patreon.com/cw/ConceptBytes)         | Engineering tutorials, 3D prints, code (22.6k members) |

**Extraction Date:** 2026-01-30  
**Creator:** Concept Bytes

---

## 📁 Extracted Components

| File                     | Source   | Purpose                                  | OneMind Use Case                         |
| ------------------------ | -------- | ---------------------------------------- | ---------------------------------------- |
| `spotify_tools.py`       | Jarvis   | Spotify playback control                 | Media control for Legacy agent           |
| `pyautogui_fallback.py`  | JarvisV3 | Desktop automation via image recognition | Control apps without APIs (Bambu Slicer) |
| `holomat_ui_concepts.md` | HoloMat2 | Radial app launcher UI design            | Flutter home screen inspiration          |

---

## 🎯 Relevance to OneMind OS

| Feature                     | Value  | Notes                                 |
| --------------------------- | ------ | ------------------------------------- |
| **OpenAI Realtime API**     | Low    | Already have LiveKit + Gemini Live    |
| **Spotify Control**         | Medium | Easy integration, add to Legacy tools |
| **PyAutoGUI Fallback**      | Medium | Useful for legacy apps without APIs   |
| **Radial UI**               | High   | Flutter home screen inspiration       |
| **Hotword Detection**       | Medium | Could add "Hey Legacy" wake word      |
| **Bambu Studio Automation** | High   | Print workflow automation             |

---

## 💡 Key Takeaways

### From JarvisV3:

- **OpenAI Realtime WebSocket** integration patterns
- **PyAutoGUI + Bambu Studio** for 3D print automation
- Clean function calling implementation

### From Jarvis (v1):

- **Hotword detection** ("Hey Jarvis")
- **Spotify integration** via spotipy
- Boston Dynamics **Spot integration** concept

### From HoloMat2:

- **Radial app launcher** UI pattern
- Touch coordinate mapping for capacitive panels
- Audio feedback patterns (sounds for actions)

### From ai_plant:

- IoT sensor + AI interpretation pattern
- Simple ChatGPT Assistants API usage

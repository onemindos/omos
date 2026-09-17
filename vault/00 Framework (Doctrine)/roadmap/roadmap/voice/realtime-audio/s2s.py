"""
JarvisV3 — OpenAI Realtime WebSocket Audio Client
===================================================
Source: https://github.com/Concept-Bytes/JarvisV3 (MIT)

Full-duplex audio over OpenAI's Realtime API (wss://api.openai.com/v1/realtime).
Key patterns:

1. MIC MUTE DURING PLAYBACK
   When AI audio arrives (response.audio.delta), mic is muted to prevent echo.
   After audio finishes (response.audio.done), mic unmutes after 500ms delay.

2. FUNCTION CALLING
   Supports run_os_command, print, and print_stl as LLM-callable tools.
   Results are sent back as function_call_output events.

3. TEXT + REALTIME MODES
   mode="text"     → text WebSocket only, no mic/speaker
   mode="realtime" → full audio duplex

4. ask_sync(message, timeout=30)
   Synchronously sends a text message and blocks until response arrives.
   Uses a queue.Queue internally — clean pattern for integrating into other code.

OneMind: This is the reference for wiring Legacy's voice output back through
the OpenClaw Node audio pipeline. The mute/unmute + 500ms delay pattern
prevents echo in rooms without acoustic treatment.

Dependencies:
    pip install pyaudio websocket-client python-dotenv

Environment (.env):
    OPENAI_API_KEY=...
    OPENAI_MODEL=gpt-4o-mini-realtime-preview-2024-12-17
    INITIAL_PROMPT=You are Legacy...
    DEVICE=mac
    VAD=True
    FUNCTION_CALLING=False
    VOICE=echo
    INCLUDE_DATE=True
    INCLUDE_TIME=True
    PC_USERNAME=YourUsername
"""

import os
import json
import base64
import threading
import websocket
import pyaudio
import datetime
import subprocess
import time
import queue
from dotenv import load_dotenv

from functions import run_os_command, print_message, print_stl

load_dotenv()


class JarvisClient:
    def __init__(
        self,
        *,
        api_key=None,
        device=None,
        model=None,
        initial_prompt=None,
        include_date=True,
        include_time=True,
        mode="realtime",
        function_calling=True,
        voice=None,
    ):
        """
        mode: "realtime" = full audio duplex | "text" = text-only WebSocket
        function_calling: enable/disable LLM tool use
        voice: alloy | ash | ballad | coral | echo | sage | shimmer | verse
        """
        self.api_key = api_key or os.getenv("OPENAI_API_KEY")
        if not self.api_key:
            raise ValueError("API key not found")
        self.device = device or os.getenv("DEVICE", "unknown")
        self.model = model or os.getenv("OPENAI_MODEL", "gpt-4o-mini-realtime-preview-2024-12-17")
        self.initial_prompt = initial_prompt or os.getenv("INITIAL_PROMPT", "")
        self.include_date = include_date
        self.include_time = include_time
        self.mode = mode
        self.function_calling = function_calling
        self.voice = voice or os.getenv("VOICE", "echo")
        self.pc_username = os.getenv("PC_USERNAME", "YourUsername")

        self.ws = None
        self.ws_thread = None
        self.running = False

        self.p = pyaudio.PyAudio()
        self.output_stream = None
        self.mute_mic = False

        # Callback for text mode responses
        self.on_text_response = None

    def append_tools_to_message(self, message: str) -> str:
        """Inject current date/time into the message if configured."""
        now = datetime.datetime.now()
        additions = []
        if self.include_date:
            additions.append(f"Date: {now.strftime('%Y-%m-%d')}")
        if self.include_time:
            additions.append(f"Time: {now.strftime('%H:%M:%S')}")
        if additions:
            return f"{message} ({' | '.join(additions)})"
        return message

    def send_text_message(self, message: str, role: str = "user"):
        full_text = self.append_tools_to_message(message)
        conversation_event = {
            "type": "conversation.item.create",
            "item": {
                "type": "message",
                "role": role,
                "content": [{"type": "input_text", "text": full_text}],
            },
        }
        self.ws.send(json.dumps(conversation_event))
        response_create_event = {
            "type": "response.create",
            "response": {"modalities": ["text"]},
        }
        self.ws.send(json.dumps(response_create_event))

    def ask_sync(self, message: str, timeout: int = 30) -> str:
        """
        Synchronously send a text message and block until response arrives.
        Uses a queue to bridge the async WebSocket callback.
        """
        response_queue = queue.Queue()
        original_callback = self.on_text_response

        def temp_callback(text):
            response_queue.put(text)

        self.on_text_response = temp_callback
        self.send_text_message(message)
        try:
            answer = response_queue.get(timeout=timeout)
        except queue.Empty:
            answer = "No response received within timeout."
        self.on_text_response = original_callback
        return answer

    def on_message(self, ws, message):
        event = json.loads(message)
        event_type = event.get("type")

        if event_type == "response.audio.delta":
            # MUTE MIC during AI audio playback (echo prevention)
            if not self.mute_mic:
                self.mute_mic = True
            audio_chunk = event.get("delta")
            if audio_chunk and self.output_stream:
                audio_data = base64.b64decode(audio_chunk)
                self.output_stream.write(audio_data)

        elif event_type == "response.audio.done":
            # UNMUTE after 500ms delay once AI audio finishes
            def delayed_unmute():
                time.sleep(0.5)
                self.mute_mic = False

            threading.Thread(target=delayed_unmute, daemon=True).start()

        elif event_type == "response.text.done":
            transcript = event.get("text", "")
            if transcript and self.on_text_response:
                self.on_text_response(transcript)

        elif event_type == "response.content_part.done":
            part = event.get("part", {})
            if part.get("type") == "text":
                final_text = part.get("text", "")
                if final_text and self.on_text_response:
                    self.on_text_response(final_text)

        elif event_type == "response.done":
            response = event.get("response", {})
            for item in response.get("output", []):
                if item.get("type") == "function_call":
                    func_name = item.get("name")
                    call_id = item.get("call_id")
                    try:
                        args = json.loads(item.get("arguments", "{}"))
                    except Exception:
                        args = {}

                    # Dispatch function calls back to tools
                    result = None
                    if func_name == "run_os_command":
                        result = run_os_command(args.get("command", ""))
                    elif func_name == "print":
                        result = print_message(args.get("message", ""))
                    elif func_name == "print_stl":
                        result = print_stl(args.get("stl_file", ""), self.pc_username)

                    if result is not None:
                        output_event = {
                            "type": "conversation.item.create",
                            "item": {
                                "type": "function_call_output",
                                "call_id": call_id,
                                "output": json.dumps({"result": result}),
                            },
                        }
                        self.ws.send(json.dumps(output_event))
                        self.ws.send(json.dumps({"type": "response.create"}))

                elif item.get("type") == "message":
                    text = item.get("text", "")
                    if text and self.on_text_response:
                        self.on_text_response(text)

    def on_error(self, ws, error):
        print("WebSocket error:", error)

    def on_close(self, ws, close_status_code, close_msg):
        print("WebSocket closed")
        self.running = False

    def on_open(self, ws):
        if self.mode == "realtime":
            self.output_stream = self.p.open(
                format=pyaudio.paInt16,
                channels=1,
                rate=24000,
                output=True,
                frames_per_buffer=1024,
            )

        session_update = {
            "type": "session.update",
            "session": {
                "voice": self.voice,
                "output_audio_format": "pcm16",
                "instructions": self.initial_prompt,
            },
        }

        if self.function_calling:
            session_update["session"]["tools"] = [
                {
                    "type": "function",
                    "name": "run_os_command",
                    "description": f"Execute an OS command on {self.device}.",
                    "parameters": {
                        "type": "object",
                        "properties": {"command": {"type": "string", "description": "Shell command to execute"}},
                        "required": ["command"],
                    },
                },
                {
                    "type": "function",
                    "name": "print",
                    "description": "Print a message in the terminal.",
                    "parameters": {
                        "type": "object",
                        "properties": {"message": {"type": "string", "description": "Message to print"}},
                        "required": ["message"],
                    },
                },
                {
                    "type": "function",
                    "name": "print_stl",
                    "description": "Open an STL file and print it via Bambu Studio.",
                    "parameters": {
                        "type": "object",
                        "properties": {"stl_file": {"type": "string", "description": "Path to STL file"}},
                        "required": ["stl_file"],
                    },
                },
            ]
            session_update["session"]["tool_choice"] = "auto"

        ws.send(json.dumps(session_update))

        if self.mode == "realtime":
            threading.Thread(target=self.send_audio, daemon=True).start()

    def send_audio(self):
        """Continuously stream mic audio to the WebSocket, respecting mute flag."""
        mic_stream = self.p.open(
            format=pyaudio.paInt16,
            channels=1,
            rate=24000,
            input=True,
            frames_per_buffer=1024,
        )
        try:
            while self.running:
                if self.mute_mic:
                    time.sleep(0.1)
                    continue
                data = mic_stream.read(1024, exception_on_overflow=False)
                encoded_data = base64.b64encode(data).decode("ascii")
                event_data = {"type": "input_audio_buffer.append", "audio": encoded_data}
                self.ws.send(json.dumps(event_data))
        except Exception as e:
            print("Audio capture error:", e)
        finally:
            mic_stream.stop_stream()
            mic_stream.close()

    def start_realtime(self):
        self.running = True
        ws_url = f"wss://api.openai.com/v1/realtime?model={self.model}"
        headers = [f"Authorization: Bearer {self.api_key}", "OpenAI-Beta: realtime=v1"]
        self.ws = websocket.WebSocketApp(
            ws_url,
            header=headers,
            on_message=self.on_message,
            on_error=self.on_error,
            on_close=self.on_close,
        )
        self.ws.on_open = self.on_open
        self.ws_thread = threading.Thread(target=self.ws.run_forever, daemon=True)
        self.ws_thread.start()

    def stop_realtime(self):
        if self.ws:
            self.running = False
            self.ws.close()
            self.ws_thread.join()
            self.ws = None
            self.ws_thread = None

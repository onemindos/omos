# Hermes Realtime code compatibility audit — November 1, 2026

## Verdict

**A credible existing integration, with the inspected core APIs compatible on both deployed Hermes versions. Not a verified drop-in production install or end-to-end voice test.** Reuse/qualify the community plugin before writing a new adapter. Prefer a single trusted-user room initially, resolve dependencies in an isolated build, then prove memory/tool/voice behavior before enabling it in Legacy.

Pinned plugin: [v0.5.1 commit6cbda7e](https://github.com/kortexa-ai/hermes-realtime/tree/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442). It is MIT-licensed. No production files, settings, providers, memories or secrets were changed or inspected. Read-only kubectl queries read public application source, executable availability and installed-package/version metadata only.

## Actual live versions, not image-tag guesses

| Runtime | Actual version | Python | Findings |
|---|---|---|---|
| Legacy, image tagv2026.9.24 | `hermes_cli.__version__ = 0.21.5`, package metadata0.21.5 |3.13.5| Version meets plugin>=0.20.0; ffmpeg present; installed aiohttp3.14.3 is below plugin's>=3.14.5 requirement. |
| New onemind-hermes, pinned imagev0.22.5 | Install stamp baseVersion0.22.5; package metadata0.0.0 is a build placeholder |3.14.9| Version/Python in declared range; ffmpeg present; aiohttp3.14.5 and Pillow13.0.0 meet checked bounds. |

Neither runtime has livekit, livekit-api, aiortc, av, onnxruntime or numpy installed. Legacy Pillow12.3.0 satisfies the declared range. A successful dependency resolution/build is still required; the installed aiohttp mismatch does not by itself prove an unsatisfiable resolver conflict. Do not pip-install over the live environment merely to test it.

Both deployed source trees contain `plugins/platforms/mattermost/plugin.yaml`. This establishes bundled Mattermost adapter availability, not its configured/enabled state or working credentials.

## Integration code confirmed

- [Native plugin manifest](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/plugin.yaml): realtime0.5.1, requiresHermes>=0.20.0, declares platform and hook capabilities.
- [Registration](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/__init__.py#L278): registers LiveKit and directWebRTC adapters through PluginContext, not a separate LLM/agent. All supplied register_platform keyword fields exist in both live PlatformEntry dataclasses. Both live PluginContext versions accept **entry_kwargs and expose register_hook.
- [Voice dispatch](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L1363): invokes Hermes speech transcription, creates the normal MessageEvent with participant identity, media and room source, then calls handle_message. [Forwarding to core](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L2612) calls BasePlatformAdapter.handle_message, retaining normal gateway dispatch.
- Both live MessageEvent types contain every inspected constructor field, including channel_prompt, media_urls/media_types, internal, timestamp and source. Required base exports AudioFormat, StreamingTTSHandle, SendResult, MessageType/ProcessingOutcome, scoped secret helper and session-key builder exist.
- AST comparison found no missing named base parameters for the13 non-dunder LiveKitAdapter overrides on either runtime. Constructor differences are intentional: the registered factory accepts configuration and the adapter supplies its platform to the base constructor.

## Memory, tools and session semantics

This is a transport into the selected Hermes profile. It should retain that profile's tools and memory through the ordinary agent pipeline, rather than replace Hermes with a LiveKit starter agent. Actual persistence/permissions must still be verified in acceptance.

**It does not automatically resume the Telegram conversation transcript.** [Adapter initialization](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L229) explicitly sets `group_sessions_per_user=False`: one conversation per LiveKit room, shared by participants, with participant identity retained for attribution. Profile-level memory and room/platform-specific conversation history are distinct.

Normal text/voice events retain the host gateway's authorization pipeline (live Legacy run_inbound.py checks `_is_user_authorized_for_source`). Room JWT admission and Hermes user allowlisting are separate gates. Keep ALLOW_ALL_USERS off and mint only scoped identities/room grants. Client-advertised RPC tools use separate deny-by-default ToolPolicy and signed client admission; do not enable arbitrary client tools for the first voice trial.

## Important multi-user media boundary

There is a source-level concern before using untrusted participants: `_on_data_received` routes `conference.capture_frame` after validating current room membership/generation; `_tool_message_is_current` checks membership, not the Hermes allowed-user list. `_capture_next_frame` appends to a room-global `_pending_captures` list, which the next text/voice event drains without per-participant ownership.

Consequently a room member who cannot invoke Hermes directly can still potentially contribute a video frame consumed with another participant's authorized turn. This is an inspected control-flow concern, not a live exploit test. Audio is also transcribed before the later gateway authorization gate, so untrusted room admission can incur speech work even if the eventual request is rejected.

Use an exclusively trusted-user room/token issuer initially. Before mixed-trust/multi-user use, gate capture/audio processing against sender authorization, bind pending media to its owner or an explicitly approved shared-room policy, and regression-test an unauthorized participant's frame followed by an authorized voice turn. This does not prove arbitrary agent-tool execution bypass: ordinary user events still enter gateway authorization.

Source: [data routing](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L1427), [capture](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L1530), [membership-only current-message check](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L2516), [shared media drain](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/adapter.py#L2594).

## Tests and proof limits

Local isolated Python3.14 virtualenv, no provider credentials or Hermes profile loaded: **58 tests passed** (`tests/test_client_admission.py`, `tests/test_audio_onset.py`). They exercise HMAC admission validation and audio-onset helpers; they do not establish LiveKit connectivity or actual Hermes tool/memory operation. The broader local selection correctly could not collect gateway-dependent tests without a full Hermes test environment, so it is not reported passing.

The pinned [upstream CI](https://github.com/kortexa-ai/hermes-realtime/actions/runs/40961048305) passed. Its bootstrap builds against **moving upstream Hermes main**, not either deployed commit, and the default pytest configuration excludes `kortexa_hermes` tests. Additional optimized streaming tests require the author's fork. [CI/test bootstrap](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/scripts/bootstrap_tests.sh), [test runner](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/scripts/run_tests.sh), [dependency manifest](https://github.com/kortexa-ai/hermes-realtime/blob/6cbda7eb3ef10a8a43a04d6ae94debc1db05c442/hermes_realtime/pyproject.toml).

No LiveKit server, plugin install, dependency upgrade, microphone/camera call, external model request, tool action or memory write was performed. To establish production readiness, build exact runtime+plugin together in isolation; test a trusted identity, memory continuity, explicit tool permissions, STT/TTS, cancellation/reconnect and room/JWT rejection; then qualify streaming separately.

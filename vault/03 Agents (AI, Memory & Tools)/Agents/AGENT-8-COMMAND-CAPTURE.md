# Agent 8: Command, Capture & Support (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Chat, Capture Modal, Profile, Support Screens, Integration
> **Issues:** 11 | **Points:** 39
> **Dependencies:** Agents 4 (Routing), 6 (Navigation), 7 (Home/Execute)

---

## Your Mission

You are Agent 8 - the final integration agent. Your job is building:

1. **Command Screen** - Chat with @mentions and /commands
2. **Universal Capture** - Voice, text, photo input modal
3. **Profile Screen** - Settings and configuration
4. **Support Screens** - Insights, Agents, Knowledge, Approvals
5. **Final Polish** - Animations, error handling, accessibility

---

## Your Assigned Issues

### Core Screens
| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-472** | Command Chat Mobile | 5 | P0 |
| **OMOS-469** | Universal Capture Component | 5 | P0 |
| **OMOS-477** | Profile & Settings Mobile | 3 | P1 |

### Capture Features
| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-478** | Voice Capture Integration | 3 | P1 |

### Support Screens
| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-473** | Agent Visibility Layer | 3 | P1 |
| **OMOS-476** | Insights Mobile View | 3 | P1 |
| **OMOS-481** | Agent Status Dashboard | 3 | P2 |
| **OMOS-482** | Memory UI Component | 3 | P2 |
| **OMOS-483** | Knowledge Mobile UI | 3 | P2 |
| **OMOS-484** | Mobile Approval UI | 3 | P2 |

### Integration
| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-485** | Integration & Polish Pass | 5 | P1 |

---

## Files to Create

```
frontend/lib/command/
├── screens/
│   └── command_screen.dart
└── widgets/
    ├── chat_message_list.dart
    ├── chat_input_bar.dart
    ├── mention_autocomplete.dart
    ├── command_autocomplete.dart
    ├── agent_message.dart
    └── thinking_indicator.dart

frontend/lib/capture/
├── capture_modal.dart
└── widgets/
    ├── capture_input.dart
    ├── voice_recorder.dart
    ├── photo_capture.dart
    ├── quick_type_buttons.dart
    └── classification_preview.dart

frontend/lib/profile/
├── screens/
│   └── profile_screen.dart
└── widgets/
    ├── profile_header.dart
    ├── settings_section.dart
    └── integration_status.dart

frontend/lib/insights/
└── screens/
    └── insights_screen.dart

frontend/lib/agents/
└── screens/
    └── agents_screen.dart

frontend/lib/knowledge/
└── screens/
    └── knowledge_screen.dart

frontend/lib/approvals/
└── widgets/
    └── approval_card.dart
```

---

## OMOS-472: Command Chat

### Design

```
┌─────────────────────────────────────┐
│ 💬 Command            Context: 🏠   │
├─────────────────────────────────────┤
│                                     │
│   🤖 Legacy                         │
│   Good morning Zeus! You have 5     │
│   tasks for today (3 HP, 2 GE).     │
│   Your PI score dropped 3 points    │
│   yesterday - want me to analyze?   │
│                                     │
│   [📊 Analyze] [📋 Tasks] [🔄 Sched]│
│                                     │
│   👤 You                            │
│   @research what happened to my PI  │
│                                     │
│   🔬 Research Agent                 │
│   Analyzing your PI components...   │
│   ⏳ Thinking...                    │
│                                     │
├─────────────────────────────────────┤
│ [@] [/]  Type a message...   🎤 📎 │
└─────────────────────────────────────┘
```

### Implementation

```dart
// frontend/lib/command/screens/command_screen.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class CommandScreen extends ConsumerStatefulWidget {
  const CommandScreen({super.key});

  @override
  ConsumerState<CommandScreen> createState() => _CommandScreenState();
}

class _CommandScreenState extends ConsumerState<CommandScreen> {
  final _inputController = TextEditingController();
  final _scrollController = ScrollController();
  bool _showMentionAutocomplete = false;
  bool _showCommandAutocomplete = false;

  @override
  Widget build(BuildContext context) {
    final messages = ref.watch(chatMessagesProvider);
    final isThinking = ref.watch(isThinkingProvider);

    return Scaffold(
      backgroundColor: const Color(0xFF0A0A0A),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Row(
          children: [
            const Icon(Icons.chat_bubble, color: Color(0xFFFF0000)),
            const SizedBox(width: 8),
            const Text(
              'Command',
              style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        actions: [
          // Context indicator
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            margin: const EdgeInsets.only(right: 16),
            decoration: BoxDecoration(
              color: const Color(0xFF1A1A1A),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.home, size: 16, color: Colors.green[400]),
                const SizedBox(width: 4),
                Text(
                  'Home',
                  style: TextStyle(color: Colors.grey[400], fontSize: 12),
                ),
              ],
            ),
          ),
        ],
      ),
      body: Column(
        children: [
          // Messages
          Expanded(
            child: messages.when(
              data: (msgs) => ListView.builder(
                controller: _scrollController,
                padding: const EdgeInsets.all(16),
                itemCount: msgs.length + (isThinking ? 1 : 0),
                itemBuilder: (context, index) {
                  if (isThinking && index == msgs.length) {
                    return const ThinkingIndicator();
                  }
                  return ChatMessage(message: msgs[index]);
                },
              ),
              loading: () => const ChatSkeleton(),
              error: (e, _) => ErrorWidget(message: e.toString()),
            ),
          ),

          // Quick Actions (when empty)
          if (messages.value?.isEmpty ?? true)
            QuickActionChips(
              actions: [
                ('📊 Analyze', '/scores'),
                ('📋 Tasks', '/tasks'),
                ('🔄 Schedule', '/schedule'),
                ('📝 Brief', '/brief'),
              ],
              onTap: (command) => _sendCommand(command),
            ),

          // Autocomplete overlays
          if (_showMentionAutocomplete)
            MentionAutocomplete(
              query: _getMentionQuery(),
              onSelect: _insertMention,
            ),
          if (_showCommandAutocomplete)
            CommandAutocomplete(
              query: _getCommandQuery(),
              onSelect: _insertCommand,
            ),

          // Input bar
          ChatInputBar(
            controller: _inputController,
            onSend: _sendMessage,
            onVoice: _startVoiceInput,
            onAttach: _attachFile,
            onMentionTrigger: () => setState(() => _showMentionAutocomplete = true),
            onCommandTrigger: () => setState(() => _showCommandAutocomplete = true),
          ),
        ],
      ),
    );
  }

  void _sendMessage(String text) async {
    if (text.trim().isEmpty) return;

    // Add user message
    ref.read(chatMessagesProvider.notifier).addUserMessage(text);
    _inputController.clear();

    // Scroll to bottom
    _scrollToBottom();

    // Send to backend
    await ref.read(chatServiceProvider).send(text);
  }

  void _sendCommand(String command) {
    _inputController.text = command;
    _sendMessage(command);
  }
}
```

### Chat Input Bar

```dart
// frontend/lib/command/widgets/chat_input_bar.dart
class ChatInputBar extends StatelessWidget {
  final TextEditingController controller;
  final Function(String) onSend;
  final VoidCallback onVoice;
  final VoidCallback onAttach;
  final VoidCallback onMentionTrigger;
  final VoidCallback onCommandTrigger;

  const ChatInputBar({
    required this.controller,
    required this.onSend,
    required this.onVoice,
    required this.onAttach,
    required this.onMentionTrigger,
    required this.onCommandTrigger,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: const BoxDecoration(
        color: Color(0xFF1A1A1A),
        border: Border(top: BorderSide(color: Color(0xFF333333))),
      ),
      child: SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // Shortcut buttons
            Row(
              children: [
                _ShortcutButton(
                  icon: Icons.alternate_email,
                  label: '@',
                  onTap: () {
                    controller.text += '@';
                    onMentionTrigger();
                  },
                ),
                const SizedBox(width: 8),
                _ShortcutButton(
                  icon: Icons.terminal,
                  label: '/',
                  onTap: () {
                    controller.text = '/';
                    onCommandTrigger();
                  },
                ),
              ],
            ),
            const SizedBox(height: 8),
            // Input field
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: controller,
                    style: const TextStyle(color: Colors.white),
                    decoration: InputDecoration(
                      hintText: 'Type a message...',
                      hintStyle: TextStyle(color: Colors.grey[600]),
                      filled: true,
                      fillColor: const Color(0xFF2A2A2A),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(24),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 12,
                      ),
                    ),
                    onChanged: (text) {
                      if (text.endsWith('@')) onMentionTrigger();
                      if (text.startsWith('/')) onCommandTrigger();
                    },
                    onSubmitted: onSend,
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  onPressed: onVoice,
                  icon: const Icon(Icons.mic, color: Colors.white),
                ),
                IconButton(
                  onPressed: onAttach,
                  icon: const Icon(Icons.attach_file, color: Colors.white),
                ),
                IconButton(
                  onPressed: () => onSend(controller.text),
                  icon: const Icon(Icons.send, color: Color(0xFFFF0000)),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
```

---

## OMOS-469: Universal Capture

```dart
// frontend/lib/capture/capture_modal.dart
class CaptureModal extends ConsumerStatefulWidget {
  const CaptureModal({super.key});

  static Future<void> show(BuildContext context) {
    return showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => const CaptureModal(),
    );
  }

  @override
  ConsumerState<CaptureModal> createState() => _CaptureModalState();
}

class _CaptureModalState extends ConsumerState<CaptureModal> {
  final _inputController = TextEditingController();
  CaptureType? _selectedType;
  Classification? _classification;
  bool _isRecording = false;
  bool _isProcessing = false;

  @override
  Widget build(BuildContext context) {
    return DraggableScrollableSheet(
      initialChildSize: 0.6,
      minChildSize: 0.4,
      maxChildSize: 0.9,
      builder: (context, scrollController) {
        return Container(
          decoration: const BoxDecoration(
            color: Color(0xFF1A1A1A),
            borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
          ),
          child: Column(
            children: [
              // Handle
              const SizedBox(height: 12),
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey[600],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 16),

              // Title
              const Text(
                'Capture',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 20),

              // Input field
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: TextField(
                  controller: _inputController,
                  style: const TextStyle(color: Colors.white),
                  maxLines: 4,
                  decoration: InputDecoration(
                    hintText: "What's on your mind?",
                    hintStyle: TextStyle(color: Colors.grey[500]),
                    filled: true,
                    fillColor: const Color(0xFF2A2A2A),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide.none,
                    ),
                  ),
                  onChanged: _onTextChanged,
                ),
              ),
              const SizedBox(height: 16),

              // Capture buttons (Voice, Photo, File)
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    _CaptureButton(
                      icon: _isRecording ? Icons.stop : Icons.mic,
                      label: _isRecording ? 'Stop' : 'Voice',
                      color: _isRecording ? Colors.red : null,
                      onTap: _toggleVoiceRecording,
                    ),
                    _CaptureButton(
                      icon: Icons.camera_alt,
                      label: 'Photo',
                      onTap: _capturePhoto,
                    ),
                    _CaptureButton(
                      icon: Icons.attach_file,
                      label: 'File',
                      onTap: _attachFile,
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 16),

              // Quick type buttons
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Row(
                  children: [
                    Text(
                      'Quick:',
                      style: TextStyle(color: Colors.grey[400], fontSize: 12),
                    ),
                    const SizedBox(width: 12),
                    _QuickTypeChip(
                      label: 'Task',
                      selected: _selectedType == CaptureType.task,
                      onTap: () => _selectType(CaptureType.task),
                    ),
                    _QuickTypeChip(
                      label: 'Note',
                      selected: _selectedType == CaptureType.note,
                      onTap: () => _selectType(CaptureType.note),
                    ),
                    _QuickTypeChip(
                      label: 'Event',
                      selected: _selectedType == CaptureType.event,
                      onTap: () => _selectType(CaptureType.event),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 16),

              // AI Classification preview
              if (_classification != null)
                ClassificationPreview(
                  classification: _classification!,
                  onAccept: () => _save(_classification!),
                  onModify: _modifyClassification,
                ),

              const Spacer(),

              // Save button
              Padding(
                padding: const EdgeInsets.all(20),
                child: SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _isProcessing ? null : _processAndSave,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFFFF0000),
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: _isProcessing
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              color: Colors.white,
                            ),
                          )
                        : const Text(
                            'Save',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  Future<void> _processAndSave() async {
    if (_inputController.text.isEmpty) return;

    setState(() => _isProcessing = true);

    try {
      // Get AI classification
      final classification = await ref.read(captureServiceProvider).classify(
        _inputController.text,
        type: _selectedType,
      );

      setState(() {
        _classification = classification;
        _isProcessing = false;
      });
    } catch (e) {
      setState(() => _isProcessing = false);
      // Show error
    }
  }

  Future<void> _save(Classification classification) async {
    await ref.read(captureServiceProvider).save(classification);
    if (mounted) Navigator.pop(context);
  }
}
```

---

## OMOS-485: Integration & Polish

### Checklist

```dart
// frontend/lib/shared/animations/page_transitions.dart

class FadePageRoute<T> extends PageRouteBuilder<T> {
  final Widget page;

  FadePageRoute({required this.page})
      : super(
          pageBuilder: (context, animation, secondaryAnimation) => page,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            return FadeTransition(opacity: animation, child: child);
          },
        );
}

// Skeleton loading
class SkeletonLoading extends StatelessWidget {
  final double width;
  final double height;
  final double borderRadius;

  const SkeletonLoading({
    this.width = double.infinity,
    required this.height,
    this.borderRadius = 8,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Shimmer.fromColors(
      baseColor: const Color(0xFF2A2A2A),
      highlightColor: const Color(0xFF3A3A3A),
      child: Container(
        width: width,
        height: height,
        decoration: BoxDecoration(
          color: const Color(0xFF2A2A2A),
          borderRadius: BorderRadius.circular(borderRadius),
        ),
      ),
    );
  }
}

// Error boundary
class ErrorBoundary extends StatelessWidget {
  final Widget child;
  final Widget Function(Object error)? onError;

  const ErrorBoundary({required this.child, this.onError, super.key});

  @override
  Widget build(BuildContext context) {
    return Builder(
      builder: (context) {
        try {
          return child;
        } catch (e) {
          return onError?.call(e) ?? ErrorCard(message: e.toString());
        }
      },
    );
  }
}
```

### Polish Tasks

- [ ] Add haptic feedback to all taps
- [ ] Add shimmer loading to all async widgets
- [ ] Add error boundaries around providers
- [ ] Add pull-to-refresh everywhere
- [ ] Test VoiceOver/TalkBack accessibility
- [ ] Optimize images and bundle size
- [ ] Add offline mode with cached data
- [ ] Test on small screens (iPhone SE)
- [ ] Test on large screens (iPad)
- [ ] Ensure 60fps animations

---

## Success Criteria

- [ ] Command chat with streaming responses
- [ ] @mention autocomplete working
- [ ] /command autocomplete working
- [ ] Universal capture modal working
- [ ] Voice capture with transcription
- [ ] Photo capture with analysis
- [ ] Profile screen with settings
- [ ] Insights screen with charts
- [ ] Agents screen with status
- [ ] Knowledge screen with search
- [ ] Approval cards functional
- [ ] Smooth animations throughout
- [ ] No crashes or memory leaks
- [ ] Accessible for screen readers

---

## When Done

Update all 11 issues to "In Review".

Post final handoff:

```markdown
## Sprint Complete: All 8 Agents

### Summary
- 56 issues implemented
- 181 story points completed
- Mobile-first navigation live
- All Agno features accessible
- Ready for production testing

### Key Deliverables
- Universal Asset System
- Context Engine
- Intelligence Scores (PI/LI/WI)
- DO DATE Scheduling
- @Mention + /Command Routing
- Mother Agent Coordination
- Workflow Automation
- Mobile UI (12 primary screens)
- Universal Capture

### Testing Required
- End-to-end flow testing
- Performance testing (<500ms context build)
- Accessibility audit
- Device testing (iOS/Android)
```

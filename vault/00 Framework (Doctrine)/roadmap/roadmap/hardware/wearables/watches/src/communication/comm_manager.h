/**
 * =============================================================================
 * Communication Manager - MQTT + LiveKit
 * =============================================================================
 */

#ifndef COMM_MANAGER_H
#define COMM_MANAGER_H

#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <functional>

#include "../config.h"
#include "../sensors/sensor_manager.h"

// Voice state
enum class VoiceState {
    IDLE,
    LISTENING,
    PROCESSING,
    SPEAKING,
    ERROR
};

// Callback types
using VoiceTranscriptCallback = std::function<void(const char*)>;
using CommandCallback = std::function<void(const char*, JsonObject&)>;

class CommManager {
public:
    CommManager();
    ~CommManager();

    // ==========================================================================
    // MQTT (Biometrics)
    // ==========================================================================

    bool beginMQTT();
    void loopMQTT();
    bool isMQTTConnected();
    void disconnectMQTT();

    // Publishing
    void publishBiometrics(BiometricSnapshot& snapshot);
    void publishAlert(AlertType type, const char* message);
    void publishStatus(const char* status);
    void publishHeartbeat();

    // Command subscription
    void onCommand(CommandCallback callback);

    // ==========================================================================
    // LiveKit (Voice)
    // ==========================================================================

    bool beginLiveKit();
    void loopLiveKit();
    bool isLiveKitConnected();
    void disconnectLiveKit();

    // Voice control
    void startListening();
    void stopListening();
    VoiceState getVoiceState();

    // Callbacks
    void onVoiceTranscript(VoiceTranscriptCallback callback);

private:
    // WiFi client
    WiFiClient _wifiClient;

    // MQTT
    PubSubClient _mqttClient;
    bool _mqttConnected;
    unsigned long _lastMQTTReconnect;
    unsigned long _lastHeartbeat;

    // LiveKit
    bool _liveKitConnected;
    VoiceState _voiceState;

    // Callbacks
    VoiceTranscriptCallback _voiceCallback;
    CommandCallback _commandCallback;

    // JSON buffer
    StaticJsonDocument<512> _jsonDoc;

    // Internal methods
    bool connectMQTT();
    void mqttCallback(char* topic, byte* payload, unsigned int length);
    static void mqttCallbackStatic(char* topic, byte* payload, unsigned int length);
    void handleCommand(const char* command, JsonObject& params);

    // LiveKit helpers
    bool fetchLiveKitToken(char* tokenBuffer, size_t bufferSize);
    void processAudio();

    // Singleton for static callback
    static CommManager* _instance;
};

#endif // COMM_MANAGER_H

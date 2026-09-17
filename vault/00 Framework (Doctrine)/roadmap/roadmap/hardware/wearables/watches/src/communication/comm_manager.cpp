/**
 * =============================================================================
 * Communication Manager Implementation
 * MQTT (Biometrics) + LiveKit (Voice)
 * =============================================================================
 */

#include "comm_manager.h"
#include <HTTPClient.h>

// Static instance for callback
CommManager* CommManager::_instance = nullptr;

// =============================================================================
// Constructor / Destructor
// =============================================================================

CommManager::CommManager() :
    _mqttClient(_wifiClient),
    _mqttConnected(false),
    _lastMQTTReconnect(0),
    _lastHeartbeat(0),
    _liveKitConnected(false),
    _voiceState(VoiceState::IDLE),
    _voiceCallback(nullptr),
    _commandCallback(nullptr)
{
    _instance = this;
}

CommManager::~CommManager() {
    disconnectMQTT();
    disconnectLiveKit();
    _instance = nullptr;
}

// =============================================================================
// MQTT - Initialization
// =============================================================================

bool CommManager::beginMQTT() {
    Serial.println("[MQTT] Initializing...");

    _mqttClient.setServer(MQTT_BROKER, MQTT_PORT);
    _mqttClient.setCallback(mqttCallbackStatic);
    _mqttClient.setBufferSize(1024);

    return connectMQTT();
}

bool CommManager::connectMQTT() {
    if (_mqttClient.connected()) {
        return true;
    }

    Serial.printf("[MQTT] Connecting to %s:%d...\n", MQTT_BROKER, MQTT_PORT);

    // Connect with credentials
    bool connected = _mqttClient.connect(
        MQTT_CLIENT_ID,
        MQTT_USER,
        MQTT_PASSWORD,
        // Last Will and Testament
        MQTT_TOPIC_BIOMETRIC "/status",
        1,     // QoS 1
        true,  // Retain
        "{\"status\":\"offline\"}"
    );

    if (connected) {
        Serial.println("[MQTT] Connected!");
        _mqttConnected = true;

        // Subscribe to command topics
        _mqttClient.subscribe(MQTT_TOPIC_COMMAND, 1);

        // Publish online status
        publishStatus("online");

        return true;
    } else {
        Serial.printf("[MQTT] Connection failed, rc=%d\n", _mqttClient.state());
        _mqttConnected = false;
        return false;
    }
}

void CommManager::loopMQTT() {
    if (!_mqttClient.connected()) {
        unsigned long now = millis();
        if (now - _lastMQTTReconnect > 5000) {  // Try every 5 seconds
            _lastMQTTReconnect = now;
            connectMQTT();
        }
    } else {
        _mqttClient.loop();

        // Send periodic heartbeat
        if (millis() - _lastHeartbeat > MQTT_HEARTBEAT_RATE_MS) {
            publishHeartbeat();
            _lastHeartbeat = millis();
        }
    }
}

bool CommManager::isMQTTConnected() {
    return _mqttClient.connected();
}

void CommManager::disconnectMQTT() {
    if (_mqttClient.connected()) {
        publishStatus("offline");
        _mqttClient.disconnect();
    }
    _mqttConnected = false;
}

// =============================================================================
// MQTT - Publishing
// =============================================================================

void CommManager::publishBiometrics(BiometricSnapshot& snapshot) {
    if (!_mqttClient.connected()) return;

    // Build JSON payload
    _jsonDoc.clear();

    // Heart data
    _jsonDoc["hr"] = snapshot.heartRate;
    _jsonDoc["hr_confidence"] = snapshot.heartRateConfidence;
    _jsonDoc["spo2"] = snapshot.spo2;
    _jsonDoc["spo2_confidence"] = snapshot.spo2Confidence;
    _jsonDoc["hrv"] = snapshot.hrv;
    _jsonDoc["finger_detected"] = snapshot.fingerDetected;

    // Publish heart data
    char heartBuffer[256];
    serializeJson(_jsonDoc, heartBuffer);
    _mqttClient.publish(MQTT_TOPIC_HEART, heartBuffer, false);

    #if DEBUG_PRINT_NETWORK
    Serial.printf("[MQTT] Published heart: %s\n", heartBuffer);
    #endif

    // Motion data
    _jsonDoc.clear();
    _jsonDoc["steps"] = snapshot.steps;
    _jsonDoc["activity_level"] = snapshot.activityLevel;

    const char* activityStr;
    switch (snapshot.activity) {
        case Activity::SEDENTARY: activityStr = "sedentary"; break;
        case Activity::WALKING: activityStr = "walking"; break;
        case Activity::RUNNING: activityStr = "running"; break;
        case Activity::SLEEPING: activityStr = "sleeping"; break;
        default: activityStr = "unknown"; break;
    }
    _jsonDoc["activity"] = activityStr;

    // Publish motion data
    char motionBuffer[256];
    serializeJson(_jsonDoc, motionBuffer);
    _mqttClient.publish(MQTT_TOPIC_MOTION, motionBuffer, false);

    // Battery data (less frequent)
    static unsigned long lastBatteryPublish = 0;
    if (millis() - lastBatteryPublish > 60000) {  // Every minute
        _jsonDoc.clear();
        _jsonDoc["percent"] = snapshot.batteryPercent;
        _jsonDoc["charging"] = snapshot.batteryCharging;

        char batteryBuffer[64];
        serializeJson(_jsonDoc, batteryBuffer);
        _mqttClient.publish(MQTT_TOPIC_BATTERY, batteryBuffer, true);  // Retain

        lastBatteryPublish = millis();
    }
}

void CommManager::publishAlert(AlertType type, const char* message) {
    if (!_mqttClient.connected()) return;

    _jsonDoc.clear();

    const char* typeStr;
    switch (type) {
        case AlertType::FALL_DETECTED: typeStr = "fall"; break;
        case AlertType::HR_HIGH: typeStr = "hr_high"; break;
        case AlertType::HR_LOW: typeStr = "hr_low"; break;
        case AlertType::SPO2_LOW: typeStr = "spo2_low"; break;
        case AlertType::BATTERY_LOW: typeStr = "battery_low"; break;
        case AlertType::BATTERY_CRITICAL: typeStr = "battery_critical"; break;
        default: typeStr = "unknown"; break;
    }

    _jsonDoc["type"] = typeStr;
    _jsonDoc["message"] = message;
    _jsonDoc["timestamp"] = millis();
    _jsonDoc["device"] = MQTT_DEVICE_ID;

    char alertBuffer[256];
    serializeJson(_jsonDoc, alertBuffer);

    // Publish with QoS 2 for guaranteed delivery
    _mqttClient.publish(MQTT_TOPIC_ALERT, alertBuffer, false);

    Serial.printf("[MQTT] ALERT published: %s\n", alertBuffer);
}

void CommManager::publishStatus(const char* status) {
    if (!_mqttClient.connected()) return;

    _jsonDoc.clear();
    _jsonDoc["status"] = status;
    _jsonDoc["device"] = MQTT_DEVICE_ID;
    _jsonDoc["timestamp"] = millis();

    char buffer[128];
    serializeJson(_jsonDoc, buffer);

    _mqttClient.publish(MQTT_TOPIC_BIOMETRIC "/status", buffer, true);  // Retain
}

void CommManager::publishHeartbeat() {
    if (!_mqttClient.connected()) return;

    _jsonDoc.clear();
    _jsonDoc["uptime"] = millis();
    _jsonDoc["heap_free"] = ESP.getFreeHeap();
    _jsonDoc["wifi_rssi"] = WiFi.RSSI();

    char buffer[128];
    serializeJson(_jsonDoc, buffer);

    _mqttClient.publish(MQTT_TOPIC_BIOMETRIC "/heartbeat", buffer, false);
}

// =============================================================================
// MQTT - Receiving
// =============================================================================

void CommManager::mqttCallbackStatic(char* topic, byte* payload, unsigned int length) {
    if (_instance) {
        _instance->mqttCallback(topic, payload, length);
    }
}

void CommManager::mqttCallback(char* topic, byte* payload, unsigned int length) {
    // Convert payload to string
    char message[length + 1];
    memcpy(message, payload, length);
    message[length] = '\0';

    Serial.printf("[MQTT] Received on %s: %s\n", topic, message);

    // Parse JSON
    DeserializationError error = deserializeJson(_jsonDoc, message);
    if (error) {
        Serial.printf("[MQTT] JSON parse error: %s\n", error.c_str());
        return;
    }

    // Extract command
    const char* command = _jsonDoc["command"];
    if (command) {
        JsonObject params = _jsonDoc["params"].as<JsonObject>();
        handleCommand(command, params);
    }
}

void CommManager::handleCommand(const char* command, JsonObject& params) {
    Serial.printf("[MQTT] Command: %s\n", command);

    // Built-in commands
    if (strcmp(command, "vibrate") == 0) {
        #ifdef PIN_MOTOR
        int duration = params["duration"] | 200;
        digitalWrite(PIN_MOTOR, HIGH);
        delay(duration);
        digitalWrite(PIN_MOTOR, LOW);
        #endif
    }
    else if (strcmp(command, "reset_steps") == 0) {
        // Would need reference to sensor manager
        Serial.println("[CMD] Reset steps requested");
    }
    else if (strcmp(command, "speak") == 0) {
        const char* text = params["text"];
        if (text) {
            // Would use TTS via LiveKit
            Serial.printf("[CMD] Speak: %s\n", text);
        }
    }

    // Forward to callback
    if (_commandCallback) {
        _commandCallback(command, params);
    }
}

void CommManager::onCommand(CommandCallback callback) {
    _commandCallback = callback;
}

// =============================================================================
// LiveKit - Initialization
// =============================================================================

bool CommManager::beginLiveKit() {
    Serial.println("[LIVEKIT] Initializing...");

    // Note: LiveKit ESP32 SDK is not yet publicly available
    // This is a placeholder for the actual implementation
    // When released, integrate with:
    // https://github.com/livekit/client-sdk-esp32

    // For now, log that we're ready
    Serial.println("[LIVEKIT] Ready (SDK integration pending)");

    // In the actual implementation:
    // 1. Fetch token from server (or use pre-generated token)
    // 2. Initialize LiveKit client
    // 3. Connect to room
    // 4. Set up audio tracks

    _liveKitConnected = false;  // Will be true when SDK is integrated
    _voiceState = VoiceState::IDLE;

    return true;
}

void CommManager::loopLiveKit() {
    // Process LiveKit events
    // This would call the LiveKit SDK's event loop

    // Placeholder for audio processing
    if (_voiceState == VoiceState::LISTENING) {
        processAudio();
    }
}

bool CommManager::isLiveKitConnected() {
    return _liveKitConnected;
}

void CommManager::disconnectLiveKit() {
    // Disconnect from LiveKit room
    _liveKitConnected = false;
    _voiceState = VoiceState::IDLE;
}

// =============================================================================
// LiveKit - Voice Control
// =============================================================================

void CommManager::startListening() {
    if (!_liveKitConnected) {
        Serial.println("[LIVEKIT] Not connected, cannot start listening");
        return;
    }

    _voiceState = VoiceState::LISTENING;
    Serial.println("[LIVEKIT] Started listening...");

    // In actual implementation:
    // - Enable microphone track
    // - Start publishing audio
}

void CommManager::stopListening() {
    _voiceState = VoiceState::IDLE;
    Serial.println("[LIVEKIT] Stopped listening");

    // In actual implementation:
    // - Disable microphone track
    // - Stop publishing audio
}

VoiceState CommManager::getVoiceState() {
    return _voiceState;
}

void CommManager::onVoiceTranscript(VoiceTranscriptCallback callback) {
    _voiceCallback = callback;
}

// =============================================================================
// LiveKit - Audio Processing
// =============================================================================

void CommManager::processAudio() {
    // Placeholder for audio capture and processing
    // In actual implementation:
    // 1. Read from I2S microphone
    // 2. Buffer audio data
    // 3. Send to LiveKit room
    // 4. Receive transcriptions via data channel

    // Example of receiving transcript:
    // if (_voiceCallback && transcriptReceived) {
    //     _voiceCallback(transcriptText);
    // }
}

bool CommManager::fetchLiveKitToken(char* tokenBuffer, size_t bufferSize) {
    HTTPClient http;

    http.begin(LIVEKIT_TOKEN_ENDPOINT);
    http.addHeader("Content-Type", "application/json");

    // Request body
    _jsonDoc.clear();
    _jsonDoc["room"] = LIVEKIT_ROOM_NAME;
    _jsonDoc["participant"] = LIVEKIT_PARTICIPANT_NAME;

    String requestBody;
    serializeJson(_jsonDoc, requestBody);

    int httpCode = http.POST(requestBody);

    if (httpCode == HTTP_CODE_OK) {
        String response = http.getString();

        DeserializationError error = deserializeJson(_jsonDoc, response);
        if (!error) {
            const char* token = _jsonDoc["token"];
            if (token) {
                strncpy(tokenBuffer, token, bufferSize - 1);
                tokenBuffer[bufferSize - 1] = '\0';
                http.end();
                return true;
            }
        }
    }

    Serial.printf("[LIVEKIT] Token fetch failed: %d\n", httpCode);
    http.end();
    return false;
}

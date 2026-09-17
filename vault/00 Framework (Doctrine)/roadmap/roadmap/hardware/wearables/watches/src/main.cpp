/**
 * =============================================================================
 * OneMind OS - ESP32 Watch Firmware
 * Legacy Neural Interface - Voice + Biometric Wearable
 * =============================================================================
 *
 * This firmware turns your ESP32 watch into a direct neural interface to Legacy.
 *
 * Features:
 * - Voice commands via LiveKit (real-time WebRTC)
 * - Biometric streaming via MQTT (HR, SpO2, motion)
 * - Fall detection with emergency alerts
 * - Step counting and activity tracking
 * - Power-efficient operation
 *
 * Hardware Support:
 * - LILYGO T-Watch 2020 V3 (primary)
 * - LILYGO T-Watch 2020 V2
 * - Watchy
 * - Custom ESP32-S3 builds
 *
 * =============================================================================
 */

#include <Arduino.h>
#include <WiFi.h>
#include <ArduinoJson.h>

#include "config.h"
#include "sensors/sensor_manager.h"
#include "communication/comm_manager.h"
#include "ui/ui_manager.h"

// =============================================================================
// Global State
// =============================================================================

// Managers
SensorManager* sensorManager = nullptr;
CommManager* commManager = nullptr;
UIManager* uiManager = nullptr;

// Task handles
TaskHandle_t sensorTaskHandle = nullptr;
TaskHandle_t commTaskHandle = nullptr;
TaskHandle_t uiTaskHandle = nullptr;

// State
volatile bool systemReady = false;
volatile bool wifiConnected = false;
volatile unsigned long lastActivity = 0;

// =============================================================================
// Forward Declarations
// =============================================================================

void sensorTask(void* parameter);
void commTask(void* parameter);
void uiTask(void* parameter);
void onBiometricData(BiometricReading& reading);
void onVoiceCommand(const char* transcript);
void onAlert(AlertType type, const char* message);
bool connectWiFi();

// =============================================================================
// Setup
// =============================================================================

void setup() {
    // Initialize serial for debugging
    Serial.begin(DEBUG_SERIAL_BAUD);
    delay(1000);

    Serial.println();
    Serial.println("==============================================");
    Serial.println("  OneMind OS - ESP32 Watch Firmware");
    Serial.println("  Legacy Neural Interface");
    Serial.println("==============================================");
    Serial.println();

    // Initialize I2C
    Wire.begin(PIN_I2C_SDA, PIN_I2C_SCL);
    Wire.setClock(400000); // 400kHz

    // Initialize managers
    Serial.println("[INIT] Creating sensor manager...");
    sensorManager = new SensorManager();
    sensorManager->onBiometricData(onBiometricData);
    sensorManager->onAlert(onAlert);

    Serial.println("[INIT] Creating communication manager...");
    commManager = new CommManager();
    commManager->onVoiceTranscript(onVoiceCommand);

    Serial.println("[INIT] Creating UI manager...");
    uiManager = new UIManager();

    // Connect to WiFi
    Serial.println("[INIT] Connecting to WiFi...");
    if (connectWiFi()) {
        Serial.println("[INIT] WiFi connected!");
        wifiConnected = true;
    } else {
        Serial.println("[INIT] WiFi connection failed!");
    }

    // Initialize sensors
    Serial.println("[INIT] Initializing sensors...");
    if (!sensorManager->begin()) {
        Serial.println("[ERROR] Sensor initialization failed!");
    }

    // Initialize communication (MQTT + LiveKit)
    if (wifiConnected) {
        Serial.println("[INIT] Initializing MQTT...");
        commManager->beginMQTT();

        Serial.println("[INIT] Initializing LiveKit...");
        commManager->beginLiveKit();
    }

    // Initialize display
    Serial.println("[INIT] Initializing display...");
    uiManager->begin();

    // Create FreeRTOS tasks
    Serial.println("[INIT] Creating tasks...");

    xTaskCreatePinnedToCore(
        sensorTask,
        "SensorTask",
        4096,
        nullptr,
        2,
        &sensorTaskHandle,
        0  // Core 0
    );

    xTaskCreatePinnedToCore(
        commTask,
        "CommTask",
        8192,
        nullptr,
        1,
        &commTaskHandle,
        1  // Core 1
    );

    xTaskCreatePinnedToCore(
        uiTask,
        "UITask",
        4096,
        nullptr,
        1,
        &uiTaskHandle,
        0  // Core 0
    );

    lastActivity = millis();
    systemReady = true;

    Serial.println("[INIT] System ready!");
    Serial.println();
}

// =============================================================================
// Main Loop
// =============================================================================

void loop() {
    // Main loop handles WiFi reconnection and power management
    // Most work is done in FreeRTOS tasks

    // Check WiFi
    if (!WiFi.isConnected() && wifiConnected) {
        Serial.println("[WIFI] Connection lost, reconnecting...");
        connectWiFi();
    }

    // Power management
    #if POWER_SAVE_ENABLED
    unsigned long idleTime = millis() - lastActivity;

    if (idleTime > DEEP_SLEEP_TIMEOUT_MS) {
        Serial.println("[POWER] Entering deep sleep...");
        // Publish offline status
        if (commManager) {
            commManager->publishStatus("sleeping");
        }
        // Enter deep sleep
        esp_deep_sleep_start();
    }
    #endif

    // Yield to other tasks
    vTaskDelay(pdMS_TO_TICKS(1000));
}

// =============================================================================
// Sensor Task
// =============================================================================

void sensorTask(void* parameter) {
    Serial.println("[TASK] Sensor task started");

    TickType_t xLastWakeTime = xTaskGetTickCount();
    const TickType_t xFrequency = pdMS_TO_TICKS(MOTION_SAMPLE_RATE_MS);

    while (true) {
        if (sensorManager) {
            sensorManager->update();
        }

        vTaskDelayUntil(&xLastWakeTime, xFrequency);
    }
}

// =============================================================================
// Communication Task
// =============================================================================

void commTask(void* parameter) {
    Serial.println("[TASK] Communication task started");

    TickType_t xLastWakeTime = xTaskGetTickCount();
    const TickType_t xFrequency = pdMS_TO_TICKS(100);  // 10Hz

    static unsigned long lastPublish = 0;

    while (true) {
        if (commManager && wifiConnected) {
            // Process MQTT
            commManager->loopMQTT();

            // Process LiveKit audio
            commManager->loopLiveKit();

            // Publish biometrics at configured rate
            if (millis() - lastPublish >= MQTT_PUBLISH_RATE_MS) {
                if (sensorManager) {
                    BiometricSnapshot snapshot = sensorManager->getSnapshot();
                    commManager->publishBiometrics(snapshot);
                }
                lastPublish = millis();
            }
        }

        vTaskDelayUntil(&xLastWakeTime, xFrequency);
    }
}

// =============================================================================
// UI Task
// =============================================================================

void uiTask(void* parameter) {
    Serial.println("[TASK] UI task started");

    TickType_t xLastWakeTime = xTaskGetTickCount();
    const TickType_t xFrequency = pdMS_TO_TICKS(UI_UPDATE_RATE_MS);

    while (true) {
        if (uiManager && sensorManager) {
            BiometricSnapshot snapshot = sensorManager->getSnapshot();
            bool connected = wifiConnected && commManager && commManager->isMQTTConnected();

            uiManager->update(snapshot, connected);
        }

        vTaskDelayUntil(&xLastWakeTime, xFrequency);
    }
}

// =============================================================================
// Callbacks
// =============================================================================

void onBiometricData(BiometricReading& reading) {
    lastActivity = millis();

    #if DEBUG_PRINT_SENSORS
    Serial.printf("[SENSOR] HR: %d, SpO2: %d%%, Steps: %d\n",
        reading.heartRate, reading.spo2, reading.steps);
    #endif
}

void onVoiceCommand(const char* transcript) {
    lastActivity = millis();

    Serial.printf("[VOICE] Transcript: %s\n", transcript);

    // Update UI
    if (uiManager) {
        uiManager->showVoiceStatus("Processing...");
    }
}

void onAlert(AlertType type, const char* message) {
    Serial.printf("[ALERT] Type: %d, Message: %s\n", type, message);

    // Publish alert via MQTT with high QoS
    if (commManager && wifiConnected) {
        commManager->publishAlert(type, message);
    }

    // Show on UI
    if (uiManager) {
        uiManager->showAlert(message);
    }

    // Vibrate motor if available
    #ifdef PIN_MOTOR
    digitalWrite(PIN_MOTOR, HIGH);
    delay(500);
    digitalWrite(PIN_MOTOR, LOW);
    #endif
}

// =============================================================================
// WiFi Connection
// =============================================================================

bool connectWiFi() {
    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    unsigned long startTime = millis();

    while (WiFi.status() != WL_CONNECTED) {
        if (millis() - startTime > WIFI_CONNECT_TIMEOUT_MS) {
            return false;
        }
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.printf("[WIFI] Connected! IP: %s\n", WiFi.localIP().toString().c_str());
    wifiConnected = true;
    return true;
}

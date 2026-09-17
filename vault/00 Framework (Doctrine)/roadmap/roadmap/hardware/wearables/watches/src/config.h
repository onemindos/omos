/**
 * =============================================================================
 * OneMind OS - ESP32 Watch Configuration
 * Legacy Neural Interface
 * =============================================================================
 */

#ifndef CONFIG_H
#define CONFIG_H

// =============================================================================
// WiFi Configuration
// =============================================================================
#define WIFI_SSID "your-wifi-ssid"
#define WIFI_PASSWORD "your-wifi-password"
#define WIFI_CONNECT_TIMEOUT_MS 30000
#define WIFI_RECONNECT_DELAY_MS 5000

// =============================================================================
// LiveKit Configuration (Voice)
// =============================================================================
#define LIVEKIT_URL "wss://your-livekit-server.livekit.cloud"
#define LIVEKIT_API_KEY "your-api-key"
#define LIVEKIT_API_SECRET "your-api-secret"

// Room configuration
#define LIVEKIT_ROOM_NAME "legacy-voice"
#define LIVEKIT_PARTICIPANT_NAME "zeus-watch"

// Token endpoint (if using server-side token generation)
#define LIVEKIT_TOKEN_ENDPOINT "https://your-server.com/api/livekit/token"

// Audio settings
#define AUDIO_SAMPLE_RATE 16000
#define AUDIO_CHANNELS 1
#define AUDIO_BITS_PER_SAMPLE 16

// =============================================================================
// MQTT Configuration (Biometrics)
// =============================================================================
#define MQTT_BROKER "your-emqx-server.com"
#define MQTT_PORT 1883
#define MQTT_USE_SSL false
#define MQTT_SSL_PORT 8883

// Authentication
#define MQTT_USER "biometric-watch"
#define MQTT_PASSWORD "your-mqtt-password"

// Device identification
#define MQTT_DEVICE_ID "watch-001"
#define MQTT_CLIENT_ID "onemind-watch-001"

// Topic prefixes
#define MQTT_TOPIC_BIOMETRIC "biometric/" MQTT_DEVICE_ID
#define MQTT_TOPIC_HEART MQTT_TOPIC_BIOMETRIC "/heart"
#define MQTT_TOPIC_MOTION MQTT_TOPIC_BIOMETRIC "/motion"
#define MQTT_TOPIC_TEMPERATURE MQTT_TOPIC_BIOMETRIC "/temperature"
#define MQTT_TOPIC_BATTERY MQTT_TOPIC_BIOMETRIC "/battery"
#define MQTT_TOPIC_ALERT MQTT_TOPIC_BIOMETRIC "/alert"

// Command topics (from Legacy)
#define MQTT_TOPIC_COMMAND "command/" MQTT_DEVICE_ID "/#"

// QoS levels
#define MQTT_QOS_BIOMETRIC 1  // At least once
#define MQTT_QOS_ALERT 2      // Exactly once (for emergencies)

// =============================================================================
// Sensor Configuration
// =============================================================================

// Heart Rate (MAX30102)
#define HR_ENABLED true
#define HR_I2C_ADDRESS 0x57
#define HR_SAMPLE_RATE_MS 1000        // Read every 1 second
#define HR_AVERAGE_SAMPLES 4          // Average over 4 readings
#define HR_LED_BRIGHTNESS 0x1F        // 0x00 - 0xFF
#define HR_LED_MODE 2                 // 1 = Red only, 2 = Red + IR

// SpO2 calculation
#define SPO2_ENABLED true
#define SPO2_MIN_CONFIDENCE 80        // Minimum confidence to report

// HRV (Heart Rate Variability)
#define HRV_ENABLED true
#define HRV_WINDOW_SIZE 30            // Number of RR intervals for calculation

// Motion (BMA423 / MPU6050)
#define MOTION_ENABLED true
#define MOTION_I2C_ADDRESS 0x19       // BMA423 default
#define MOTION_SAMPLE_RATE_MS 100     // Read every 100ms
#define MOTION_RANGE 2                // +/- 2g

// Step counting
#define STEP_ENABLED true
#define STEP_THRESHOLD 1.2            // g-force threshold for step

// Fall detection
#define FALL_DETECTION_ENABLED true
#define FALL_THRESHOLD_G 3.0          // Free-fall detection threshold
#define FALL_IMPACT_G 6.0             // Impact detection threshold
#define FALL_ALERT_DELAY_MS 5000      // Wait before alerting (cancel window)

// =============================================================================
// Data Publishing
// =============================================================================
#define MQTT_PUBLISH_RATE_MS 5000     // Publish biometrics every 5 seconds
#define MQTT_HEARTBEAT_RATE_MS 30000  // Publish heartbeat every 30 seconds

// Batching (reduce network usage)
#define BATCH_ENABLED true
#define BATCH_SIZE 10                 // Number of readings to batch
#define BATCH_MAX_AGE_MS 30000        // Force publish after 30 seconds

// =============================================================================
// Display Configuration
// =============================================================================
#define DISPLAY_ENABLED true
#define DISPLAY_WIDTH 240
#define DISPLAY_HEIGHT 240
#define DISPLAY_ROTATION 0
#define DISPLAY_BRIGHTNESS 128        // 0-255
#define DISPLAY_TIMEOUT_MS 10000      // Turn off after 10 seconds inactivity

// Status update rate
#define UI_UPDATE_RATE_MS 500         // Update display every 500ms

// =============================================================================
// Battery Configuration
// =============================================================================
#define BATTERY_ENABLED true
#define BATTERY_ADC_PIN 35            // ADC pin for battery voltage
#define BATTERY_CHECK_RATE_MS 60000   // Check every minute
#define BATTERY_LOW_THRESHOLD 20      // Alert at 20%
#define BATTERY_CRITICAL_THRESHOLD 10 // Critical at 10%

// =============================================================================
// Voice Activation
// =============================================================================
#define WAKE_WORD_ENABLED true
#define WAKE_WORD "legacy"            // Or "hey legacy"
#define VOICE_ACTIVITY_THRESHOLD 500  // Audio level to trigger recording
#define VOICE_TIMEOUT_MS 5000         // Stop recording after silence

// =============================================================================
// Power Management
// =============================================================================
#define POWER_SAVE_ENABLED true
#define LIGHT_SLEEP_ENABLED true
#define DEEP_SLEEP_TIMEOUT_MS 300000  // Deep sleep after 5 min inactivity

// Activity thresholds for power modes
#define POWER_MODE_ACTIVE_THRESHOLD 0.5   // g-force for "active" mode
#define POWER_MODE_IDLE_TIMEOUT_MS 30000  // Switch to idle after 30s no movement

// =============================================================================
// Debug Configuration
// =============================================================================
#define DEBUG_ENABLED true
#define DEBUG_SERIAL_BAUD 115200
#define DEBUG_PRINT_SENSORS true
#define DEBUG_PRINT_NETWORK true
#define DEBUG_PRINT_VOICE false        // Very verbose

// =============================================================================
// Pin Definitions (T-Watch 2020 V3)
// =============================================================================
#ifdef LILYGO_WATCH_2020_V3
    #define PIN_I2C_SDA 21
    #define PIN_I2C_SCL 22
    #define PIN_TOUCH_INT 38
    #define PIN_MOTOR 4
    #define PIN_MIC_DATA 34
    #define PIN_MIC_CLK 26
    #define PIN_SPEAKER 25
    #define PIN_IR_SEND 13
    #define PIN_BATTERY_ADC 35
#endif

// =============================================================================
// Pin Definitions (Custom ESP32-S3)
// =============================================================================
#ifdef CUSTOM_ESP32S3
    #define PIN_I2C_SDA 21
    #define PIN_I2C_SCL 22
    #define PIN_MIC_DATA 41
    #define PIN_MIC_CLK 42
    #define PIN_MIC_WS 40
    #define PIN_SPEAKER 17
    #define PIN_BATTERY_ADC 1
    #define PIN_HR_INT 19
#endif

#endif // CONFIG_H

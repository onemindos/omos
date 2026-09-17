/**
 * =============================================================================
 * Sensor Manager - Biometric Data Collection
 * =============================================================================
 */

#ifndef SENSOR_MANAGER_H
#define SENSOR_MANAGER_H

#include <Arduino.h>
#include <Wire.h>
#include <functional>

// Alert types
enum class AlertType {
    FALL_DETECTED,
    HR_HIGH,
    HR_LOW,
    SPO2_LOW,
    BATTERY_LOW,
    BATTERY_CRITICAL
};

// Activity classification
enum class Activity {
    UNKNOWN,
    SEDENTARY,
    WALKING,
    RUNNING,
    SLEEPING
};

// Single biometric reading
struct BiometricReading {
    uint32_t timestamp;

    // Heart
    uint8_t heartRate;
    uint8_t heartRateConfidence;
    uint8_t spo2;
    uint8_t spo2Confidence;
    uint16_t hrv;  // RMSSD in ms

    // Motion
    float accelX;
    float accelY;
    float accelZ;
    uint32_t steps;
    Activity activity;

    // Battery
    uint8_t batteryPercent;
    bool batteryCharging;
};

// Snapshot for publishing
struct BiometricSnapshot {
    // Heart
    uint8_t heartRate;
    uint8_t heartRateConfidence;
    uint8_t spo2;
    uint8_t spo2Confidence;
    uint16_t hrv;

    // Motion
    uint32_t steps;
    Activity activity;
    float activityLevel;  // 0.0 - 1.0

    // Battery
    uint8_t batteryPercent;
    bool batteryCharging;

    // Status
    bool fingerDetected;
    bool isValid;
};

// Callback types
using BiometricCallback = std::function<void(BiometricReading&)>;
using AlertCallback = std::function<void(AlertType, const char*)>;

class SensorManager {
public:
    SensorManager();
    ~SensorManager();

    // Initialization
    bool begin();

    // Main update loop (call frequently)
    void update();

    // Get current snapshot
    BiometricSnapshot getSnapshot();

    // Callbacks
    void onBiometricData(BiometricCallback callback);
    void onAlert(AlertCallback callback);

    // Individual sensor access
    uint8_t getHeartRate();
    uint8_t getSpO2();
    uint16_t getHRV();
    uint32_t getSteps();
    Activity getActivity();
    uint8_t getBatteryPercent();

    // Control
    void enableHeartRate(bool enable);
    void enableMotion(bool enable);
    void resetSteps();

private:
    // Sensor state
    bool _hrEnabled;
    bool _motionEnabled;
    bool _initialized;

    // Current readings
    BiometricReading _currentReading;

    // Callbacks
    BiometricCallback _biometricCallback;
    AlertCallback _alertCallback;

    // Timing
    unsigned long _lastHRRead;
    unsigned long _lastMotionRead;
    unsigned long _lastBatteryRead;

    // HRV calculation
    static const int HRV_BUFFER_SIZE = 30;
    uint16_t _rrIntervals[HRV_BUFFER_SIZE];
    int _rrIndex;

    // Step counting
    float _lastAccelMagnitude;
    bool _stepState;

    // Fall detection
    unsigned long _fallDetectedTime;
    bool _fallAlertSent;

    // Internal methods
    void readHeartRate();
    void readMotion();
    void readBattery();
    uint16_t calculateHRV();
    Activity classifyActivity();
    void checkFallDetection(float ax, float ay, float az);
    void checkAlerts();
};

#endif // SENSOR_MANAGER_H

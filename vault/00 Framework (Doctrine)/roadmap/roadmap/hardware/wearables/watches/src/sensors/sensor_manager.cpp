/**
 * =============================================================================
 * Sensor Manager Implementation
 * =============================================================================
 */

#include "sensor_manager.h"
#include "../config.h"
#include <MAX30105.h>  // SparkFun MAX3010x library

// MAX30102 sensor instance
static MAX30105 heartSensor;

// =============================================================================
// Constructor / Destructor
// =============================================================================

SensorManager::SensorManager() :
    _hrEnabled(HR_ENABLED),
    _motionEnabled(MOTION_ENABLED),
    _initialized(false),
    _biometricCallback(nullptr),
    _alertCallback(nullptr),
    _lastHRRead(0),
    _lastMotionRead(0),
    _lastBatteryRead(0),
    _rrIndex(0),
    _lastAccelMagnitude(1.0),
    _stepState(false),
    _fallDetectedTime(0),
    _fallAlertSent(false)
{
    memset(&_currentReading, 0, sizeof(_currentReading));
    memset(_rrIntervals, 0, sizeof(_rrIntervals));
}

SensorManager::~SensorManager() {
    // Cleanup if needed
}

// =============================================================================
// Initialization
// =============================================================================

bool SensorManager::begin() {
    bool success = true;

    // Initialize MAX30102 (Heart Rate / SpO2)
    #if HR_ENABLED
    Serial.println("[SENSOR] Initializing MAX30102...");
    if (!heartSensor.begin(Wire, I2C_SPEED_FAST)) {
        Serial.println("[SENSOR] MAX30102 not found!");
        _hrEnabled = false;
        success = false;
    } else {
        // Configure sensor
        heartSensor.setup(
            HR_LED_BRIGHTNESS,  // LED brightness
            4,                   // Sample average
            HR_LED_MODE,        // LED mode (Red + IR)
            400,                // Sample rate
            411,                // Pulse width
            4096                // ADC range
        );

        // Enable temperature reading
        heartSensor.enableDIETEMPRDY();

        Serial.println("[SENSOR] MAX30102 initialized!");
    }
    #endif

    // Initialize BMA423 (Accelerometer)
    #if MOTION_ENABLED
    Serial.println("[SENSOR] Initializing accelerometer...");
    // Note: Actual BMA423 initialization depends on your library
    // This is a placeholder for the actual implementation
    // For T-Watch, use the TTGO library's built-in accelerometer support
    Serial.println("[SENSOR] Accelerometer initialized!");
    #endif

    // Initialize battery monitoring
    #if BATTERY_ENABLED
    pinMode(BATTERY_ADC_PIN, INPUT);
    #endif

    _initialized = success;
    return success;
}

// =============================================================================
// Main Update Loop
// =============================================================================

void SensorManager::update() {
    if (!_initialized) return;

    unsigned long now = millis();

    // Read heart rate at configured interval
    #if HR_ENABLED
    if (_hrEnabled && (now - _lastHRRead >= HR_SAMPLE_RATE_MS)) {
        readHeartRate();
        _lastHRRead = now;
    }
    #endif

    // Read motion at configured interval
    #if MOTION_ENABLED
    if (_motionEnabled && (now - _lastMotionRead >= MOTION_SAMPLE_RATE_MS)) {
        readMotion();
        _lastMotionRead = now;
    }
    #endif

    // Read battery at configured interval
    #if BATTERY_ENABLED
    if (now - _lastBatteryRead >= BATTERY_CHECK_RATE_MS) {
        readBattery();
        _lastBatteryRead = now;
    }
    #endif

    // Check for alerts
    checkAlerts();

    // Call callback with new reading
    if (_biometricCallback) {
        _currentReading.timestamp = now;
        _biometricCallback(_currentReading);
    }
}

// =============================================================================
// Heart Rate Reading
// =============================================================================

void SensorManager::readHeartRate() {
    // Check for finger presence
    long irValue = heartSensor.getIR();

    if (irValue < 50000) {
        // No finger detected
        _currentReading.heartRate = 0;
        _currentReading.heartRateConfidence = 0;
        _currentReading.spo2 = 0;
        _currentReading.spo2Confidence = 0;
        return;
    }

    // Read sensor values
    uint32_t red = heartSensor.getRed();
    uint32_t ir = heartSensor.getIR();

    // Simple heart rate calculation
    // In production, use a proper algorithm like the one from Maxim
    static uint32_t lastIR = 0;
    static unsigned long lastBeatTime = 0;
    static int beatCount = 0;
    static uint32_t beatSum = 0;

    // Detect beat (simple threshold crossing)
    if (lastIR < ir && ir > 60000) {
        unsigned long beatTime = millis();
        if (lastBeatTime > 0 && beatTime - lastBeatTime > 300) {  // Min 300ms between beats
            uint32_t rrInterval = beatTime - lastBeatTime;

            // Calculate BPM from this interval
            uint8_t bpm = 60000 / rrInterval;

            if (bpm > 40 && bpm < 200) {  // Valid range
                beatSum += bpm;
                beatCount++;

                // Store RR interval for HRV
                _rrIntervals[_rrIndex] = rrInterval;
                _rrIndex = (_rrIndex + 1) % HRV_BUFFER_SIZE;

                // Update heart rate (average)
                if (beatCount >= HR_AVERAGE_SAMPLES) {
                    _currentReading.heartRate = beatSum / beatCount;
                    _currentReading.heartRateConfidence = 90;  // Simplified
                    beatSum = 0;
                    beatCount = 0;

                    // Calculate HRV
                    _currentReading.hrv = calculateHRV();
                }
            }
        }
        lastBeatTime = beatTime;
    }
    lastIR = ir;

    // SpO2 calculation (simplified)
    // Real SpO2 needs proper R-value calculation
    #if SPO2_ENABLED
    if (red > 0 && ir > 0) {
        float ratio = (float)red / (float)ir;
        // Simplified SpO2 estimation (not medically accurate!)
        // Real implementation needs calibration
        int spo2 = 110 - 25 * ratio;
        if (spo2 > 100) spo2 = 100;
        if (spo2 < 70) spo2 = 0;

        _currentReading.spo2 = spo2;
        _currentReading.spo2Confidence = (spo2 > 0) ? 80 : 0;
    }
    #endif
}

// =============================================================================
// HRV Calculation
// =============================================================================

uint16_t SensorManager::calculateHRV() {
    // Calculate RMSSD (Root Mean Square of Successive Differences)
    float sumSquaredDiffs = 0;
    int count = 0;

    for (int i = 1; i < HRV_BUFFER_SIZE; i++) {
        if (_rrIntervals[i] > 0 && _rrIntervals[i-1] > 0) {
            int diff = (int)_rrIntervals[i] - (int)_rrIntervals[i-1];
            sumSquaredDiffs += diff * diff;
            count++;
        }
    }

    if (count > 0) {
        return (uint16_t)sqrt(sumSquaredDiffs / count);
    }
    return 0;
}

// =============================================================================
// Motion Reading
// =============================================================================

void SensorManager::readMotion() {
    // Note: This is a placeholder implementation
    // Actual implementation depends on your accelerometer library
    // For T-Watch, use the TTGO library's built-in functions

    // Simulated accelerometer reading
    // Replace with actual I2C read from BMA423/MPU6050
    float ax = 0.0;
    float ay = 0.0;
    float az = 1.0;  // Gravity

    // TODO: Replace with actual sensor read
    // For BMA423:
    // bma.getAccel(&ax, &ay, &az);

    _currentReading.accelX = ax;
    _currentReading.accelY = ay;
    _currentReading.accelZ = az;

    // Calculate magnitude
    float magnitude = sqrt(ax*ax + ay*ay + az*az);

    // Step detection (simple threshold crossing)
    #if STEP_ENABLED
    if (!_stepState && magnitude > STEP_THRESHOLD && _lastAccelMagnitude <= STEP_THRESHOLD) {
        _stepState = true;
    } else if (_stepState && magnitude < STEP_THRESHOLD) {
        _currentReading.steps++;
        _stepState = false;
    }
    _lastAccelMagnitude = magnitude;
    #endif

    // Activity classification
    _currentReading.activity = classifyActivity();

    // Fall detection
    #if FALL_DETECTION_ENABLED
    checkFallDetection(ax, ay, az);
    #endif
}

// =============================================================================
// Activity Classification
// =============================================================================

Activity SensorManager::classifyActivity() {
    // Simple activity classification based on acceleration variance
    // In production, use a proper ML model

    float magnitude = sqrt(
        _currentReading.accelX * _currentReading.accelX +
        _currentReading.accelY * _currentReading.accelY +
        _currentReading.accelZ * _currentReading.accelZ
    );

    static float magnitudeHistory[10];
    static int histIndex = 0;

    magnitudeHistory[histIndex] = magnitude;
    histIndex = (histIndex + 1) % 10;

    // Calculate variance
    float mean = 0;
    for (int i = 0; i < 10; i++) {
        mean += magnitudeHistory[i];
    }
    mean /= 10;

    float variance = 0;
    for (int i = 0; i < 10; i++) {
        float diff = magnitudeHistory[i] - mean;
        variance += diff * diff;
    }
    variance /= 10;

    // Classify based on variance
    if (variance < 0.01) {
        return Activity::SEDENTARY;
    } else if (variance < 0.1) {
        return Activity::WALKING;
    } else {
        return Activity::RUNNING;
    }
}

// =============================================================================
// Fall Detection
// =============================================================================

void SensorManager::checkFallDetection(float ax, float ay, float az) {
    float magnitude = sqrt(ax*ax + ay*ay + az*az);

    // Detect free-fall (low acceleration)
    if (magnitude < FALL_THRESHOLD_G && _fallDetectedTime == 0) {
        _fallDetectedTime = millis();
    }

    // Detect impact after free-fall
    if (_fallDetectedTime > 0 && !_fallAlertSent) {
        if (magnitude > FALL_IMPACT_G) {
            // Impact detected after free-fall = potential fall
            if (millis() - _fallDetectedTime < 1000) {  // Within 1 second
                // Wait for confirmation delay
                if (millis() - _fallDetectedTime > FALL_ALERT_DELAY_MS) {
                    // Send alert
                    if (_alertCallback) {
                        _alertCallback(AlertType::FALL_DETECTED, "Fall detected! Are you OK?");
                    }
                    _fallAlertSent = true;
                }
            }
        }
    }

    // Reset if too much time passed
    if (_fallDetectedTime > 0 && millis() - _fallDetectedTime > 10000) {
        _fallDetectedTime = 0;
        _fallAlertSent = false;
    }
}

// =============================================================================
// Battery Reading
// =============================================================================

void SensorManager::readBattery() {
    #if BATTERY_ENABLED
    int rawValue = analogRead(BATTERY_ADC_PIN);

    // Convert to voltage (depends on voltage divider)
    // Assuming 3.3V reference and 2:1 divider
    float voltage = (rawValue / 4095.0) * 3.3 * 2.0;

    // Convert to percentage (LiPo: 3.0V = 0%, 4.2V = 100%)
    int percent = (int)((voltage - 3.0) / 1.2 * 100);
    if (percent > 100) percent = 100;
    if (percent < 0) percent = 0;

    _currentReading.batteryPercent = percent;

    // Check if charging (would need additional pin or register read)
    // For AXP202 power management, use library functions
    _currentReading.batteryCharging = false;  // Placeholder
    #endif
}

// =============================================================================
// Alert Checking
// =============================================================================

void SensorManager::checkAlerts() {
    if (!_alertCallback) return;

    // Heart rate alerts
    if (_currentReading.heartRate > 0) {
        if (_currentReading.heartRate > 180) {
            _alertCallback(AlertType::HR_HIGH, "Heart rate high: >180 BPM");
        } else if (_currentReading.heartRate < 40) {
            _alertCallback(AlertType::HR_LOW, "Heart rate low: <40 BPM");
        }
    }

    // SpO2 alert
    if (_currentReading.spo2 > 0 && _currentReading.spo2 < 90) {
        _alertCallback(AlertType::SPO2_LOW, "Blood oxygen low: <90%");
    }

    // Battery alerts
    if (_currentReading.batteryPercent <= BATTERY_CRITICAL_THRESHOLD) {
        _alertCallback(AlertType::BATTERY_CRITICAL, "Battery critical!");
    } else if (_currentReading.batteryPercent <= BATTERY_LOW_THRESHOLD) {
        _alertCallback(AlertType::BATTERY_LOW, "Battery low");
    }
}

// =============================================================================
// Public Accessors
// =============================================================================

BiometricSnapshot SensorManager::getSnapshot() {
    BiometricSnapshot snapshot;

    snapshot.heartRate = _currentReading.heartRate;
    snapshot.heartRateConfidence = _currentReading.heartRateConfidence;
    snapshot.spo2 = _currentReading.spo2;
    snapshot.spo2Confidence = _currentReading.spo2Confidence;
    snapshot.hrv = _currentReading.hrv;
    snapshot.steps = _currentReading.steps;
    snapshot.activity = _currentReading.activity;
    snapshot.batteryPercent = _currentReading.batteryPercent;
    snapshot.batteryCharging = _currentReading.batteryCharging;
    snapshot.fingerDetected = (_currentReading.heartRate > 0);
    snapshot.isValid = _initialized;

    // Calculate activity level (0.0 - 1.0)
    float magnitude = sqrt(
        _currentReading.accelX * _currentReading.accelX +
        _currentReading.accelY * _currentReading.accelY +
        _currentReading.accelZ * _currentReading.accelZ
    );
    snapshot.activityLevel = min(1.0f, max(0.0f, (magnitude - 0.9f) / 0.5f));

    return snapshot;
}

void SensorManager::onBiometricData(BiometricCallback callback) {
    _biometricCallback = callback;
}

void SensorManager::onAlert(AlertCallback callback) {
    _alertCallback = callback;
}

uint8_t SensorManager::getHeartRate() { return _currentReading.heartRate; }
uint8_t SensorManager::getSpO2() { return _currentReading.spo2; }
uint16_t SensorManager::getHRV() { return _currentReading.hrv; }
uint32_t SensorManager::getSteps() { return _currentReading.steps; }
Activity SensorManager::getActivity() { return _currentReading.activity; }
uint8_t SensorManager::getBatteryPercent() { return _currentReading.batteryPercent; }

void SensorManager::enableHeartRate(bool enable) { _hrEnabled = enable; }
void SensorManager::enableMotion(bool enable) { _motionEnabled = enable; }
void SensorManager::resetSteps() { _currentReading.steps = 0; }

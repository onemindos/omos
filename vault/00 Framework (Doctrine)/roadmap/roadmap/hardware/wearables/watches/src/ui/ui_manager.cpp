/**
 * =============================================================================
 * UI Manager Implementation
 * Display and Touch Interface
 * =============================================================================
 *
 * Note: This is a placeholder implementation.
 * For T-Watch, use the TTGO_TWatch_Library for display functions.
 * For LVGL-based displays, use lvgl library.
 */

#include "ui_manager.h"
#include "../config.h"

// For T-Watch, include:
// #include <LilyGoWatch.h>

UIManager::UIManager() :
    _initialized(false),
    _displayOn(true),
    _lastUpdate(0),
    _lastActivity(0),
    _lastHR(0),
    _lastSpO2(0),
    _lastSteps(0),
    _lastBattery(0)
{
}

UIManager::~UIManager() {
}

bool UIManager::begin() {
    Serial.println("[UI] Initializing display...");

    // Platform-specific initialization
    #ifdef LILYGO_WATCH_2020_V3
    // TTGOClass *watch = TTGOClass::getWatch();
    // watch->begin();
    // watch->openBL();
    // watch->tft->fillScreen(TFT_BLACK);
    #endif

    _initialized = true;
    _lastActivity = millis();

    Serial.println("[UI] Display initialized!");
    return true;
}

void UIManager::update(BiometricSnapshot& snapshot, bool connected) {
    if (!_initialized) return;

    unsigned long now = millis();

    // Check for timeout
    if (_displayOn && (now - _lastActivity > DISPLAY_TIMEOUT_MS)) {
        turnOff();
        return;
    }

    // Only update if values changed or periodic refresh
    bool needsUpdate = false;
    if (snapshot.heartRate != _lastHR ||
        snapshot.spo2 != _lastSpO2 ||
        snapshot.steps != _lastSteps ||
        snapshot.batteryPercent != _lastBattery ||
        (now - _lastUpdate > 1000)) {  // At least every second
        needsUpdate = true;
    }

    if (needsUpdate && _displayOn) {
        drawMainScreen(snapshot, connected);

        _lastHR = snapshot.heartRate;
        _lastSpO2 = snapshot.spo2;
        _lastSteps = snapshot.steps;
        _lastBattery = snapshot.batteryPercent;
        _lastUpdate = now;
    }
}

void UIManager::drawMainScreen(BiometricSnapshot& snapshot, bool connected) {
    // Platform-specific drawing
    // This is a placeholder - implement for your display library

    #ifdef LILYGO_WATCH_2020_V3
    // Example for T-Watch with TFT_eSPI:
    // TTGOClass *watch = TTGOClass::getWatch();
    // TFT_eSPI *tft = watch->tft;

    // tft->fillScreen(TFT_BLACK);

    // // Header
    // drawTime();
    // drawBattery(snapshot.batteryPercent, snapshot.batteryCharging);
    // drawConnectionStatus(connected);

    // // Main content
    // drawHeartRate(snapshot.heartRate, snapshot.heartRateConfidence);
    // drawSpO2(snapshot.spo2);
    // drawSteps(snapshot.steps);
    #endif

    // Serial output for debugging without display
    Serial.printf("[UI] HR: %d | SpO2: %d%% | Steps: %d | Batt: %d%% | %s\n",
        snapshot.heartRate,
        snapshot.spo2,
        snapshot.steps,
        snapshot.batteryPercent,
        connected ? "Connected" : "Disconnected"
    );
}

void UIManager::drawHeartRate(uint8_t hr, uint8_t confidence) {
    // Draw heart rate with heart icon
    // Placeholder for actual implementation
}

void UIManager::drawSpO2(uint8_t spo2) {
    // Draw SpO2 percentage
    // Placeholder for actual implementation
}

void UIManager::drawSteps(uint32_t steps) {
    // Draw step count with footstep icon
    // Placeholder for actual implementation
}

void UIManager::drawBattery(uint8_t percent, bool charging) {
    // Draw battery indicator in corner
    // Placeholder for actual implementation
}

void UIManager::drawConnectionStatus(bool connected) {
    // Draw WiFi/MQTT connection icon
    // Placeholder for actual implementation
}

void UIManager::drawTime() {
    // Draw current time
    // Would use RTC or NTP time
    // Placeholder for actual implementation
}

void UIManager::showVoiceStatus(const char* status) {
    _lastActivity = millis();

    if (!_displayOn) {
        turnOn();
    }

    // Show voice status overlay
    Serial.printf("[UI] Voice: %s\n", status);
}

void UIManager::showAlert(const char* message) {
    _lastActivity = millis();

    if (!_displayOn) {
        turnOn();
    }

    // Show alert overlay
    Serial.printf("[UI] ALERT: %s\n", message);

    #ifdef LILYGO_WATCH_2020_V3
    // TTGOClass *watch = TTGOClass::getWatch();
    // watch->tft->fillRect(20, 80, 200, 80, TFT_RED);
    // watch->tft->setTextColor(TFT_WHITE);
    // watch->tft->drawString(message, 30, 110);
    #endif
}

void UIManager::showNotification(const char* title, const char* body) {
    _lastActivity = millis();

    if (!_displayOn) {
        turnOn();
    }

    Serial.printf("[UI] Notification: %s - %s\n", title, body);
}

void UIManager::setBrightness(uint8_t brightness) {
    #ifdef LILYGO_WATCH_2020_V3
    // TTGOClass *watch = TTGOClass::getWatch();
    // watch->setBrightness(brightness);
    #endif
}

void UIManager::turnOn() {
    if (_displayOn) return;

    _displayOn = true;
    _lastActivity = millis();

    #ifdef LILYGO_WATCH_2020_V3
    // TTGOClass *watch = TTGOClass::getWatch();
    // watch->openBL();
    #endif

    Serial.println("[UI] Display on");
}

void UIManager::turnOff() {
    if (!_displayOn) return;

    _displayOn = false;

    #ifdef LILYGO_WATCH_2020_V3
    // TTGOClass *watch = TTGOClass::getWatch();
    // watch->closeBL();
    #endif

    Serial.println("[UI] Display off");
}

bool UIManager::isOn() {
    return _displayOn;
}

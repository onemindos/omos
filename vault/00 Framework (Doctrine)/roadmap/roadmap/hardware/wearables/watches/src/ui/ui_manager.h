/**
 * =============================================================================
 * UI Manager - Display and Touch Interface
 * =============================================================================
 */

#ifndef UI_MANAGER_H
#define UI_MANAGER_H

#include <Arduino.h>
#include "../sensors/sensor_manager.h"

class UIManager {
public:
    UIManager();
    ~UIManager();

    // Initialization
    bool begin();

    // Update display
    void update(BiometricSnapshot& snapshot, bool connected);

    // Status displays
    void showVoiceStatus(const char* status);
    void showAlert(const char* message);
    void showNotification(const char* title, const char* body);

    // Screen control
    void setBrightness(uint8_t brightness);
    void turnOn();
    void turnOff();
    bool isOn();

private:
    bool _initialized;
    bool _displayOn;
    unsigned long _lastUpdate;
    unsigned long _lastActivity;

    // Display state
    uint8_t _lastHR;
    uint8_t _lastSpO2;
    uint32_t _lastSteps;
    uint8_t _lastBattery;

    // Internal methods
    void drawMainScreen(BiometricSnapshot& snapshot, bool connected);
    void drawHeartRate(uint8_t hr, uint8_t confidence);
    void drawSpO2(uint8_t spo2);
    void drawSteps(uint32_t steps);
    void drawBattery(uint8_t percent, bool charging);
    void drawConnectionStatus(bool connected);
    void drawTime();
};

#endif // UI_MANAGER_H

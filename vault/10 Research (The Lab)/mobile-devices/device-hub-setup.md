# Device Hub - Platform Setup Guide

This guide walks through setting up iOS and Android platforms for the Device Hub.

## Prerequisites

- Flutter SDK installed
- Xcode installed (for iOS)
- Android Studio installed (for Android)

## Step 1: Add Platform Support

The frontend currently only has web support. Add iOS and Android:

```bash
cd frontend

# Add iOS support
flutter create --platforms=ios .

# Add Android support
flutter create --platforms=android .
```

## Step 2: iOS Configuration

### Info.plist

After running `flutter create`, edit `ios/Runner/Info.plist` and add these entries inside the `<dict>` tag:

```xml
<!-- Location Permissions -->
<key>NSLocationAlwaysAndWhenInUseUsageDescription</key>
<string>OneMind needs location access to provide context-aware assistance and automation. Your location is used to trigger smart home actions, provide relevant suggestions, and enable location-based reminders.</string>

<key>NSLocationWhenInUseUsageDescription</key>
<string>OneMind uses your location to provide relevant suggestions and show nearby information.</string>

<key>NSLocationAlwaysUsageDescription</key>
<string>OneMind needs always-on location access to track your location in the background and trigger automations when you arrive or leave places.</string>

<!-- Background Modes -->
<key>UIBackgroundModes</key>
<array>
    <string>location</string>
    <string>fetch</string>
    <string>processing</string>
    <string>remote-notification</string>
</array>

<!-- Motion Sensors -->
<key>NSMotionUsageDescription</key>
<string>OneMind uses motion data to track your steps and detect falls for safety.</string>

<!-- Bluetooth (for beacons) -->
<key>NSBluetoothAlwaysUsageDescription</key>
<string>OneMind uses Bluetooth to detect nearby beacons for indoor positioning.</string>

<key>NSBluetoothPeripheralUsageDescription</key>
<string>OneMind uses Bluetooth to connect to nearby devices.</string>
```

### Xcode Capabilities

Open `ios/Runner.xcworkspace` in Xcode and enable:

1. **Signing & Capabilities** tab
2. Click **+ Capability**
3. Add:
   - Background Modes
     - [x] Location updates
     - [x] Background fetch
     - [x] Background processing
     - [x] Remote notifications
   - Push Notifications

### Podfile

Edit `ios/Podfile` and ensure minimum iOS version is 14.0:

```ruby
platform :ios, '14.0'
```

Then run:
```bash
cd ios && pod install && cd ..
```

## Step 3: Android Configuration

### AndroidManifest.xml

Edit `android/app/src/main/AndroidManifest.xml`:

```xml
<manifest xmlns:android="http://schemas.android.com/apk/res/android">

    <!-- Location Permissions -->
    <uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
    <uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
    <uses-permission android:name="android.permission.ACCESS_BACKGROUND_LOCATION" />

    <!-- Foreground Service -->
    <uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
    <uses-permission android:name="android.permission.FOREGROUND_SERVICE_LOCATION" />

    <!-- Activity Recognition -->
    <uses-permission android:name="android.permission.ACTIVITY_RECOGNITION" />

    <!-- Wake Lock -->
    <uses-permission android:name="android.permission.WAKE_LOCK" />
    <uses-permission android:name="android.permission.RECEIVE_BOOT_COMPLETED" />

    <!-- Network -->
    <uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
    <uses-permission android:name="android.permission.ACCESS_WIFI_STATE" />
    <uses-permission android:name="android.permission.INTERNET" />

    <!-- Bluetooth -->
    <uses-permission android:name="android.permission.BLUETOOTH" />
    <uses-permission android:name="android.permission.BLUETOOTH_ADMIN" />
    <uses-permission android:name="android.permission.BLUETOOTH_SCAN" />
    <uses-permission android:name="android.permission.BLUETOOTH_CONNECT" />

    <application
        android:label="OneMind"
        android:name="${applicationName}"
        android:icon="@mipmap/ic_launcher">

        <!-- Main Activity -->
        <activity
            android:name=".MainActivity"
            android:exported="true"
            android:launchMode="singleTop"
            android:theme="@style/LaunchTheme"
            android:configChanges="orientation|keyboardHidden|keyboard|screenSize|smallestScreenSize|locale|layoutDirection|fontScale|screenLayout|density|uiMode"
            android:hardwareAccelerated="true"
            android:windowSoftInputMode="adjustResize">
            <intent-filter>
                <action android:name="android.intent.action.MAIN"/>
                <category android:name="android.intent.category.LAUNCHER"/>
            </intent-filter>
        </activity>

        <!-- Background Location Service -->
        <service
            android:name="id.flutter.flutter_background_service.BackgroundService"
            android:foregroundServiceType="location"
            android:exported="false" />

        <!-- Boot Receiver for auto-start -->
        <receiver
            android:name="id.flutter.flutter_background_service.BootReceiver"
            android:exported="false">
            <intent-filter>
                <action android:name="android.intent.action.BOOT_COMPLETED" />
                <action android:name="android.intent.action.QUICKBOOT_POWERON" />
            </intent-filter>
        </receiver>

        <!-- WorkManager initialization -->
        <provider
            android:name="androidx.startup.InitializationProvider"
            android:authorities="${applicationId}.androidx-startup"
            android:exported="false"
            tools:node="merge">
            <meta-data
                android:name="androidx.work.WorkManagerInitializer"
                android:value="androidx.startup" />
        </provider>

    </application>
</manifest>
```

### build.gradle

Edit `android/app/build.gradle`:

```groovy
android {
    compileSdkVersion 34

    defaultConfig {
        minSdkVersion 23  // Required for background location
        targetSdkVersion 34
    }
}
```

## Step 4: Install Dependencies

```bash
cd frontend
flutter pub get
```

## Step 5: Build and Run

### iOS
```bash
flutter run -d ios
```

### Android
```bash
flutter run -d android
```

### Web (existing)
```bash
flutter run -d chrome
```

## Troubleshooting

### iOS: Location permission not showing
- Ensure Info.plist has all location description keys
- Clean build: `flutter clean && flutter pub get`
- Delete app from device and reinstall

### Android: Background location not working
- Android 10+ requires `ACCESS_BACKGROUND_LOCATION` permission
- User must explicitly grant "Allow all the time" in settings
- Foreground service notification is required

### Build errors after adding platforms
```bash
flutter clean
flutter pub get
cd ios && pod install && cd ..  # iOS only
flutter run
```

## Next Steps

After platform setup, the Device Hub services will automatically:
1. Request location permissions on first launch
2. Start background tracking based on user preferences
3. Stream sensor data to the backend

See [DEVICE_HUB.md](DEVICE_HUB.md) for full feature documentation.


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]

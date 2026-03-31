# NMEA Signal Generator

A lightweight, cross-platform application to emit NMEA signals via Serial Port.

## Features
- Modern Glassmorphism UI.
- Supports GPGGA, GPRMC, GPVTG sentences.
- Adjustable frequency (1-10 Hz).
- Cross-platform: Linux, Windows, macOS.

## Quick Build (Linux Ubuntu 24.04)
To build the application on your current machine, use the provided script:
```bash
./scripts/build_linux.sh
```
The binary will be located in `build/bin/nmea-gen`.

## Multi-Platform Packaging (CI/CD)
This project is equipped with a **GitHub Actions** workflow. To get versions for Windows and macOS:
1. Push this code to a GitHub repository.
2. Go to the **Actions** tab on GitHub.
3. Download the zipped artifacts once the build is finished.

## Requirements
- **Go** 1.24+
- **Node.js** 22+
- **Wails** CLI
- On Linux: `pkg-config`, `libgtk-3-dev`, `libwebkit2gtk-4.1-dev`.
# nmea-gen

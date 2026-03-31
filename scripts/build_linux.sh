#!/bin/bash

# Multi-Platform Packaging Script for Linux (Ubuntu 24.04)
# This script ensures the correct environment variables and tags are used for a clean build.

echo "--- Building NMEA Signal Generator for Linux ---"

# 1. Export correct pkg-config for system libraries (especially if Linuxbrew is present)
export PKG_CONFIG=/usr/bin/pkg-config

# 2. Compile using the webkit2_41 tag (required for Ubuntu Noble 24.04)
wails build -tags webkit2_41

# 3. Handle results
if [ $? -eq 0 ]; then
    echo "--- Build successful! ---"
    echo "Binary location: build/bin/nmea-gen"
    
    # Optional: Create a distribution folder
    mkdir -p dist/linux
    cp build/bin/nmea-gen dist/linux/
    echo "Snapshot saved in dist/linux/"
else
    echo "--- Build failed. Check the error log above. ---"
    exit 1
fi

#!/bin/bash

set -e

echo "Building frontend assets"
cd ui
npm install
npm run build
cd ..

echo "Building application"
./build-app.sh

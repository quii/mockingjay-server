#!/bin/bash

echo "Building frontend assets"
cd ui
npm install
npm run build
cd ..

echo "Building application"
time ./build-app.sh

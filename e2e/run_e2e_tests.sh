#!/bin/bash

# Ultimate Frisbee API - End-to-End Test Runner
# This script runs comprehensive E2E tests for the Ultimate Frisbee API

set -e

echo "🏆 Ultimate Frisbee API - End-to-End Test Suite"
echo "================================================="
echo ""

# Check if API is running
echo "🔍 Checking if API is running on http://127.0.0.1:42007..."
if ! curl -s -f "http://127.0.0.1:42007/v1/health/" > /dev/null 2>&1; then
    echo "❌ ERROR: API is not running on http://127.0.0.1:42007"
    echo ""
    echo "Please start the API first:"
    echo "  make deps/start"
    echo "  make db/migration/up" 
    echo "  make db/seed"
    echo "  make run/api"
    echo ""
    exit 1
fi

echo "✅ API is running!"
echo ""

# Run the tests
echo "🧪 Running End-to-End Tests..."
echo "================================"
echo ""

# Install test dependencies if needed
if ! go list github.com/stretchr/testify/assert >/dev/null 2>&1; then
    echo "📦 Installing test dependencies..."
    go mod tidy
    go get github.com/stretchr/testify/assert
    go get github.com/stretchr/testify/require
    echo ""
fi

# Run the actual tests
go test -v -run TestUltimateFrisbeeAPI_E2E ./...

echo ""
echo "🎉 End-to-End Tests Completed!"
echo ""
echo "📊 Test Summary:"
echo "  ✅ Health Check"
echo "  ✅ Team CRUD Operations (Create, Read, Update)"
echo "  ✅ Error Handling (404s, validation)"
echo "  ✅ Data Validation"
echo "  ✅ Seeded Data Verification"
echo ""
echo "Your Ultimate Frisbee API is working correctly! 🥏"
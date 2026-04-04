#!/usr/bin/env bash
# 
# This script is used to run the acceptance tests for the Terraform provider.
# It sets the necessary environment variables and then runs the tests using the `go test` command.

# - TF_ACC
# - SOLIDFIRE_USERNAME
# - SOLIDFIRE_PASSWORD
# - SOLIDFIRE_SERVER
# - SOLIDFIRE_API_VERSION
# 
export TF_ACC=1
SOLIDFIRE_USERNAME="${SOLIDFIRE_USERNAME:-admin}"
SOLIDFIRE_PASSWORD="${SOLIDFIRE_PASSWORD:-admin}"
SOLIDFIRE_SERVER="${SOLIDFIRE_SERVER:-192.168.1.34}"
SOLIDFIRE_API_VERSION="${SOLIDFIRE_API_VERSION:-12.5}"

# Replication/Pairing variables
REPLICATION=false
SOLIDFIRE_SERVER_DR="${SOLIDFIRE_SERVER_DR:-}"
SOLIDFIRE_USERNAME_DR="${SOLIDFIRE_USERNAME_DR:-$SOLIDFIRE_USERNAME}"
SOLIDFIRE_PASSWORD_DR="${SOLIDFIRE_PASSWORD_DR:-$SOLIDFIRE_PASSWORD}"

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --test-replication) REPLICATION=true ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

for ENV_VAR in SOLIDFIRE_USERNAME SOLIDFIRE_PASSWORD SOLIDFIRE_SERVER SOLIDFIRE_API_VERSION; do
  if [ -z "${!ENV_VAR}" ]; then
    echo "Error: Environment variable $ENV_VAR is not set."
    exit 1
  fi
done

if [ "$REPLICATION" = true ]; then
    if [ -z "$SOLIDFIRE_SERVER_DR" ]; then
        echo "Error: --test-replication requires SOLIDFIRE_SERVER_DR to be set."
        exit 1
    fi
    export SOLIDFIRE_SERVER_DR
    export SOLIDFIRE_USERNAME_DR
    export SOLIDFIRE_PASSWORD_DR
fi

# Run tests from the project root
cd "$(dirname "$0")/.." || exit 1

# Run tests
if [ "$REPLICATION" = true ]; then
    make testacc-pairing
else
    make testacc
fi



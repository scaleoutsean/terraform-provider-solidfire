#!/usr/bin/env bash
# 
# This script is used to run the acceptance tests for the Terraform provider.
# It sets the necessary environment variables and then runs the tests using the `go test` command.

# - TF_ACC
# - ELEMENTSW_USERNAME
# - ELEMENTSW_PASSWORD
# - ELEMENTSW_SERVER
# - ELEMENTSW_API_VERSION
# 
export TF_ACC=1
ELEMENTSW_USERNAME="${ELEMENTSW_USERNAME:-admin}"
ELEMENTSW_PASSWORD="${ELEMENTSW_PASSWORD:-admin}"
ELEMENTSW_SERVER="${ELEMENTSW_SERVER:-192.168.1.34}"
ELEMENTSW_API_VERSION="${ELEMENTSW_API_VERSION:-12.5}"

# Replication/Pairing variables
REPLICATION=false
ELEMENTSW_SERVER_DR="${ELEMENTSW_SERVER_DR:-}"
ELEMENTSW_USERNAME_DR="${ELEMENTSW_USERNAME_DR:-$ELEMENTSW_USERNAME}"
ELEMENTSW_PASSWORD_DR="${ELEMENTSW_PASSWORD_DR:-$ELEMENTSW_PASSWORD}"

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --test-replication) REPLICATION=true ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

for ENV_VAR in ELEMENTSW_USERNAME ELEMENTSW_PASSWORD ELEMENTSW_SERVER ELEMENTSW_API_VERSION; do
  if [ -z "${!ENV_VAR}" ]; then
    echo "Error: Environment variable $ENV_VAR is not set."
    exit 1
  fi
done

if [ "$REPLICATION" = true ]; then
    if [ -z "$ELEMENTSW_SERVER_DR" ]; then
        echo "Error: --test-replication requires ELEMENTSW_SERVER_DR to be set."
        exit 1
    fi
    export ELEMENTSW_SERVER_DR
    export ELEMENTSW_USERNAME_DR
    export ELEMENTSW_PASSWORD_DR
fi

# Run tests from the project root
cd "$(dirname "$0")/.." || exit 1

# Run tests
if [ "$REPLICATION" = true ]; then
    make testacc-pairing
else
    make testacc-account
fi



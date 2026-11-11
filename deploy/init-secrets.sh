#!/usr/bin/env bash
set -euo pipefail

SECRETS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/.secrets"
mkdir -p "$SECRETS_DIR"
chmod 700 "$SECRETS_DIR"

head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n' > "$SECRETS_DIR/jwt_secret.txt"
head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n' > "$SECRETS_DIR/key_encryption_secret.txt"
chmod 600 "$SECRETS_DIR"/*.txt

echo "Generated CSPRNG secret files in $SECRETS_DIR"

#!/bin/bash

set -e

CERT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"

openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout "$KEY_FILE" \
  -out "$CERT_FILE" \
  -days 365 \
  -subj "/C=RU/ST=Moscow/L=Moscow/O=MedSoft Labs/OU=HL7 System/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,DNS:reception-api,DNS:hospital-srv,IP:127.0.0.1"

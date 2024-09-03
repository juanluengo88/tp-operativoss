#!/bin/bash

# Verificar que se hayan pasado argumentos
if [ -z "$1" ]; then
    echo "Uso: $0 PROCESS_NAME"
    exit 1
fi

if [ -z "$KERNEL_PORT" ]; then
    echo "The KERNEL_PORT is not set"
    echo "Using default port 8001"
    KERNEL_PORT=8001
fi

if [ -z "$KERNEL_HOST" ]; then
    echo "The KERNEL_HOST is not set"
    echo "Using default host localhost"
    KERNEL_HOST=localhost
fi

KERNEL_URL="http://$KERNEL_HOST:$KERNEL_PORT"

curl -X PUT -H "Content-Type: application/json" -d "{\"path\": \"/home/utnso/tp-2024-1c-sudoers/procesos/$1\"}" ${KERNEL_URL}/process
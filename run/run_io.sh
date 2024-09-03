#!/bin/bash

# Definir la URL del kernel
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

echo "IO TEST"

# Obtener la ruta del directorio donde se encuentra el script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROCESOS_DIR="$SCRIPT_DIR/../procesos"

KERNEL_URL="http://$KERNEL_HOST:$KERNEL_PORT"

# Lista de archivos de procesos, relativos al script
procesos=(
    "$PROCESOS_DIR/IO_A.txt"
    "$PROCESOS_DIR/IO_B.txt"
    "$PROCESOS_DIR/IO_C.txt"
)

# Crear cada proceso usando la API
for proceso in "${procesos[@]}"; do
    echo "Creando proceso desde el archivo $proceso"
    curl -X PUT "$KERNEL_URL/process" -H "Content-Type: application/json" -d "{\"path\":\"$proceso\"}"
    
done

# Hacer una petición PUT a /plani después de iniciar todos los procesos
echo "Enviando petición a /plani"
curl -X PUT "$KERNEL_URL/plani"
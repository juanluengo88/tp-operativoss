# Deploy

1. Chequear que Go este instalado -> `go version` sino `sudo snap install go --classic`
1. Clonar repo con `github token`
1. `cd tp-2024-1c-sudoers`
1. Agregar variables de ambiente en Kernel con su `host` e `ip` para los `.sh`
1. Modificar la config del modulo con las `ips`, `puertos` y `path` necesarios
1. Builder el modulo con make
1. Usar el run especifico de la prueba

## Prueba Planificación

- Kernel:
  - FIFO: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_plani_fifo.json`
  - Round Robin: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_plani_rr.json`
  - Virtual Round Robin: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_plani_vrr.json`
- CPU: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_plani.json`
- Memoria: `make memoria ENV=prod C=/home/utnso/tp-2024-1c-sudoers/memoria/config/config_plani.json`
- IO Generica: `make entradasalida ENV=prod N=SLP1 P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_slp1.json`

Desde la VM de Kernel:

```bash
cd run/
./run_plani.sh
# - FIFO -> Eliminar proceso cuando 4 queda en loop
./delete_process.sh 4
```

## Prueba Deadlock

- Kernel: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_deadlock.json`
- CPU: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_deadlock.json`
- Memoria: `make memoria ENV=prod C=/home/utnso/tp-2024-1c-sudoers/memoria/config/config_deadlock.json`
- IO Generica: `make entradasalida ENV=prod N=ESPERA P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_espera.json`

  Desde la VM de Kernel:

```bash
cd run/
./run_deadlock.sh
# Esperar que quede bloqueado y eliminar un proceso
./delete_process.sh <PID>
```

## Prueba Memoria y TLB

- Kernel: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_memoria_tlb.json`
- CPU:
  - FIFO: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_memoria_tlb_fifo.json`
  - LRU: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_memoria_tlb_lru.json`
- Memoria: `make memoria ENV=prod C=/home/utnso/tp-2024-1c-sudoers/memoria/config/config_deadlock.json`

Desde la VM de Kernel:

```bash
cd run/
# Correr primero con FIFO
./run_memoria.sh MEMORIA_1.txt
# Cambiar a config a LRU
./run_memoria.sh MEMORIA_1.txt
# Esperar a que finalice
./create_process.sh MEMORIA_2.txt
# Esperar a que finalice
./create_process.sh MEMORIA_3.txt
```

## Prueba IO

- Kernel: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_io.json`
- CPU: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_io.json`
- Memoria: `make memoria ENV=prod C=/home/utnso/tp-2024-1c-sudoers/memoria/config/config_io.json`
- IO Teclado: `make entradasalida ENV=prod N=TECLADO P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_teclado.json`
- IO Monitor: `make entradasalida ENV=prod N=MONITOR P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_monitor.json`
- IO Generica: `make entradasalida ENV=prod N=GENERICA P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_generica.json`

Desde la VM de Kernel:

```bash
cd run/
./run_io.sh
# Escribir WAR NEVER CHANGES...
# Escribir Sistemas Operativos 1c2024
```

## Prueba FS

- Kernel: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_fs.json`
- CPU: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_fs.json`
- Memoria: `make memoria ENV=prod C=/home/utnso/tp-2024-1c-sudoers/memoria/config/config_fs.json`
- IO FS: `make entradasalida ENV=prod N=FS P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_fs.json`
- IO Teclado: `make entradasalida ENV=prod N=TECLADO P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_teclado.json`
- IO Monitor: `make entradasalida ENV=prod N=MONITOR P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_monitor.json`

Desde la VM de Kernel:

```bash
cd run/
# Levantar primero FS y teclado
./run_write_fs.sh
# Escribir: Fallout 1 Fallout 2 Fallout 3 Fallout: New Vegas Fallout 4 Fallout 76
# Esperar que termine, bajar todo, levantar FS y monitor
./run_read_fs.sh
```

## Prueba Salvation's Edge

- Kernel: `make kernel ENV=prod C=/home/utnso/tp-2024-1c-sudoers/kernel/config/config_salvation.json`
- CPU: `make cpu ENV=prod C=/home/utnso/tp-2024-1c-sudoers/cpu/config/config_salvation.json`
- Memoria: `make memoria ENV=prod C=/home/utnso/tp-2024-1c-sudoers/memoria/config/config_salvation.json`
- IO Generica: `make entradasalida ENV=prod N=GENERICA P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_generica.json`
- IO SLP1: `make entradasalida ENV=prod N=SLP1 P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_slp1.json`
- IO ESPERA: `make entradasalida ENV=prod N=ESPERA P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_espera.json`
- IO Teclado: `make entradasalida ENV=prod N=TECLADO P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_teclado.json`
- IO Monitor: `make entradasalida ENV=prod N=MONITOR P=/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_monitor.json`

### TODO

## Checkpoint Tag

Para cada checkpoint de control obligatorio, se debe crear un tag en el
repositorio con el siguiente formato:

```
checkpoint-{número}
```

Donde `{número}` es el número del checkpoint.

Para crear un tag y subirlo al repositorio, podemos utilizar los siguientes
comandos:

```bash
git tag -a checkpoint-{número} -m "Checkpoint {número}"
git push origin checkpoint-{número}
```

Asegúrense de que el código compila y cumple con los requisitos del checkpoint
antes de subir el tag.

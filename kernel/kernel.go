package main

import (
	"fmt"
	"os"

	api "github.com/sisoputnfrba/tp-golang/kernel/api"
	global "github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func main() {

	// Me crea el loger y la configuracion - Lo puedo usar en cualquier parte del modulo ahora
	global.InitGlobal()

	// Creo la config con su puerto y respectivas rutas
	// Uso el modulo utils/server para crear el servidor con la configuracion anterior
	s := api.CreateServer()

	// Levanto server
	global.Logger.Log(fmt.Sprintf("Starting kernel server on port: %d", global.KernelConfig.Port), log.INFO)
	if err := s.Start(); err != nil {
		global.Logger.Log(fmt.Sprintf("Failed to start kernel server: %v", err), log.ERROR)
		os.Exit(1)
	}

	global.Logger.CloseLogger()
}

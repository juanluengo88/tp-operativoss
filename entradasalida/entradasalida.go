package main

import (
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/entradasalida/api"
	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func main() {

	// Me crea el loger y la configuracion
	global.InitGlobal()

	s := api.CreateServer()

	global.Logger.Log(fmt.Sprintf("Starting IO server on port: %d", global.IOConfig.Port), log.INFO)

	err := s.Start()
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Failed to start IO server: %v", err), log.ERROR)
		os.Exit(1)
	}

	global.Logger.CloseLogger()
}

package main

import (
	"fmt"

	api "github.com/sisoputnfrba/tp-golang/memoria/api"
	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
)

func main() {
	global.InitGlobal()

	s := api.CreateServer()

	global.Logger.Log(fmt.Sprintf("Starting Memory server on port: %d", global.MemoryConfig.Port), log.INFO)

	if err := s.Start(); err != nil {
		global.Logger.Log(fmt.Sprintf("Failed to start Memory server: %v", err), log.ERROR)
		os.Exit(1)
	}

	global.Logger.CloseLogger()
}

package api

import (
	"net/http"

	handlers "github.com/sisoputnfrba/tp-golang/entradasalida/api/handlers"
	global "github.com/sisoputnfrba/tp-golang/entradasalida/global"
	"github.com/sisoputnfrba/tp-golang/utils/server"
)

func CreateServer() *server.Server {

	configServer := server.Config{
		Port: global.IOConfig.Port,
		Handlers: map[string]http.HandlerFunc{
			"PUT /IO_GEN_SLEEP":    handlers.Sleep,
			"PUT /IO_STDIN_READ":   handlers.Stdin_read,
			"PUT /IO_STDOUT_WRITE": handlers.Stdout_write,
			"PUT /IO_FS_CREATE":    handlers.Fs_create,
			"PUT /IO_FS_DELETE":    handlers.Fs_delete,
			"PUT /IO_FS_TRUNCATE":  handlers.Fs_truncate,
			"PUT /IO_FS_WRITE":     handlers.Fs_write,
			"PUT /IO_FS_READ":      handlers.Fs_read,
		},
	}
	return server.NewServer(configServer)
}

package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/shortterm"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Este mensaje se encargará de retomar (en caso que se encuentre pausada)
// la planificación de corto y largo plazo. En caso que la planificación no
// se encuentre pausada, se debe ignorar el mensaje.


func InitPlanningHandler(w http.ResponseWriter, r *http.Request) {

	if global.ReadyState.Len() > 0 || global.ExecuteState.Len() > 0 {
		global.Logger.Log("REAUNDO PLANI", log.DEBUG)
		<- global.SemStopPlani
		global.Logger.Log("Libero SemStopPlani", log.DEBUG)
		if global.NewState.Len() > 0 {
			<- global.SemLongStopPlani
			global.Logger.Log("Libero SemLongStopPlani", log.DEBUG)
		}
		if global.BlockedState.Len() > 0 {
			<- global.SemBlockStopPlani
			global.Logger.Log("Libero SemBlockStopPlani", log.DEBUG)

		}
	} else {
		global.Logger.Log("INICIO PLANI", log.DEBUG)

		global.MutexPlani.Lock()
		global.WorkingPlani = true
		global.MutexPlani.Unlock()

		go longterm.InitLongTermPlani()
		go shortterm.InitShortTermPlani()
	}

	w.WriteHeader(http.StatusNoContent)
}

// Este mensaje se encargará de pausar la planificación de corto y largo plazo.
// El proceso que se encuentra en ejecución NO es desalojado, pero una vez que salga
// de EXEC se va a pausar el manejo de su motivo de desalojo. De la misma forma,
// los procesos bloqueados van a pausar su transición a la cola de Ready.
func StopPlanningHandler(w http.ResponseWriter, r *http.Request) {

	global.MutexPlani.Lock()
	global.WorkingPlani = false
	global.MutexPlani.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

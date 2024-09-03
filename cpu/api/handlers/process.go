package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	internal "github.com/sisoputnfrba/tp-golang/cpu/internal/dispatch"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"

	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type PCB struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

type InterruptReason struct {
	Reason string `json:"reason"`
}

func Dispatch(w http.ResponseWriter, r *http.Request) {
	pcb := &model.PCB{}
	err := serialization.DecodeHTTPBody(r, pcb)
	if err != nil {
		http.Error(w, "Error al decodear PCB", http.StatusBadRequest)
		global.Logger.Log(fmt.Sprintf("Error al decodear PCB: %v", err), log.ERROR)
		return
	}

	pcb, _ = internal.Dispatch(pcb)

	serialization.EncodeHTTPResponse(w, pcb, http.StatusOK)
}

func Interrupt(w http.ResponseWriter, r *http.Request) {
	interruptReason := InterruptReason{}

	err := serialization.DecodeHTTPBody(r, &interruptReason)
	if err != nil {
		http.Error(w, "Error al decodear PCB", http.StatusBadRequest)
		global.Logger.Log(fmt.Sprintf("Error al decodear PCB: %v", err), log.ERROR)
		return
	}

	global.Logger.Log("Interrumpimos por "+interruptReason.Reason, log.DEBUG)

	global.ExecuteMutex.Lock()
	global.Execute = false
	global.InterruptReason = interruptReason.Reason
	global.ExecuteMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

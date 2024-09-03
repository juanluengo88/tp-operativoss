package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	global "github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func InitProcessHandler(w http.ResponseWriter, r *http.Request) {

	var pPath utils.ProcessPath
	err := serialization.DecodeHTTPBody(r, &pPath)

	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}

	pcb := pcb.CreateNewProcess()

	type ProcessMemory struct {
		Path string `json:"path"`
		PID  int    `json:"pid"`
	}
	processMemory := ProcessMemory{
		Path: pPath.Path,
		PID:  pcb.PID,
	}

	requests.PutHTTPwithBody[ProcessMemory, interface{}](global.KernelConfig.IPMemory, global.KernelConfig.PortMemory, "process", processMemory)
	// _, err = requests.PutHTTPwithBody[ProcessMemory, interface{}](global.KernelConfig.IPMemory, global.KernelConfig.PortMemory, "process", processMemory)
	// if err != nil {
	// 	global.Logger.Log("Error al enviar instruccion "+err.Error(), log.ERROR)
	// 	http.Error(w, "Error al enviar instruccion", http.StatusBadRequest)
	// 	return
	// }

	global.MutexNewState.Lock()
	global.NewState.PushBack(pcb)
	global.MutexNewState.Unlock()

	global.SemNewList <- struct{}{}

	processPID := utils.ProcessPID{PID: pcb.PID}

	global.Logger.Log(fmt.Sprintf("Se crea el proceso %d en NEW", pcb.PID), log.INFO)

	err = serialization.EncodeHTTPResponse(w, processPID, http.StatusCreated)
	if err != nil {
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}
}

func EndProcessHandler(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(r.PathValue("pid"))

	// TODO: Request a memoria
	// VER LOGICA DE: Traerme la PCB - Ver estado actual
	// 	- Si esta en CPU -> INTERRUPT
	// AHORA SOLO LO BORRA DE LA LISTA

	if !utils.RemoveProcessByPID(pid) {
		global.Logger.Log(fmt.Sprintf("No existe el PID %d", pid), log.DEBUG)
		http.Error(w, fmt.Sprintf("No existe el PID %d", pid), http.StatusNotFound)
		return
	}

	global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INTERRUPTED_BY_USER", pid), log.INFO)
	w.WriteHeader(http.StatusNoContent)
}

func ListProcessHandler(w http.ResponseWriter, r *http.Request) {
	allProcess := utils.GetAllProcess()

	err := serialization.EncodeHTTPResponse(w, allProcess, http.StatusOK)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func ProcessByIdHandler(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(r.PathValue("pid"))

	pcb := utils.FindProcessInList(pid)

	if pcb == nil {
		global.Logger.Log(fmt.Sprintf("No existe el PID %d", pid), log.DEBUG)
		http.Error(w, fmt.Sprintf("No existe el PID %d", pid), http.StatusNotFound)
		return
	}

	processState := utils.ProcessState{
		PID:   pcb.PID,
		State: pcb.State,
	}

	err := serialization.EncodeHTTPResponse(w, processState, http.StatusOK)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	global.Logger.Log(fmt.Sprintf("Process %d - State: %s", pcb.PID, pcb.State), log.DEBUG)
}

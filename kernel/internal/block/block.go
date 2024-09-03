package block

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"

	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// var acceptedInstructions map[string] []string
var acceptedInstructions = map[string][]string{
	"GEN":    {"IO_GEN_SLEEP"},
	"STDOUT": {"IO_STDOUT_WRITE"},
	"STDIN":  {"IO_STDIN_READ"},
	"DIALFS": {"IO_FS_CREATE", "IO_FS_DELETE", "IO_FS_TRUNCATE", "IO_FS_WRITE", "IO_FS_READ"},
}

func ProcessToIO(pcb *model.PCB) {
	// time, _ := strconv.Atoi(pcb.Instruction.Parameters[1])
	global.Logger.Log(fmt.Sprintf("Proceso bloqueado %+v", pcb), log.DEBUG)

	io := factoryIO(pcb)

	global.Logger.Log(fmt.Sprintf("IO: %+v", io), log.DEBUG)

	if !CheckIfExist(io.GetName()) || !CheckIfIsValid(io.GetName(), io.GetInstruction()) {
		moveToExit(pcb)
		return
	}

	global.IoMap[io.GetName()].Sem <- 0

	if !global.WorkingPlani {
		global.Logger.Log("Me bloqueo en block - SemBlockStopPlani", log.DEBUG)
		global.SemBlockStopPlani <- 0
		global.Logger.Log("Se libero el semaforo en block - SemBlockStopPlani", log.DEBUG)
	}

	_, err := requests.PutHTTPwithBody[IO, interface{}](global.KernelConfig.IPIo, global.IoMap[io.GetName()].Port, io.GetInstruction(), io)
	if err != nil {
		global.Logger.Log("Se desconecto IO:"+err.Error(), log.DEBUG)
		delete(global.IoMap, io.GetName())
		<-global.IoMap[io.GetName()].Sem
		moveToExit(pcb)
		return
	}

	<-global.IoMap[io.GetName()].Sem

	if pcb.State == "EXIT" {
		return
	}

	BlockToReady(pcb)

	arrayReady := longterm.ConvertListToArray(global.ReadyState)
	arrayPlus := longterm.ConvertListToArray(global.ReadyPlus)

	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	// Revisar si funciona la logica dsp del &&
	if global.KernelConfig.PlanningAlgorithm == "VRR" && len(arrayPlus) > 0 {
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v, Cola Ready+ : %v", arrayReady, arrayPlus), log.INFO)
	} else {
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v", arrayReady), log.INFO)
	}

	global.SemReadyList <- struct{}{}
}

func moveToExit(pcb *model.PCB) {
	global.MutexBlockState.Lock()
	global.BlockedState.Remove(global.BlockedState.Front())
	global.MutexBlockState.Unlock()

	pcb.State = "EXIT"

	global.MutexExitState.Lock()
	global.ExitState.PushBack(pcb)
	global.MutexExitState.Unlock()

	// Request a memoria para eliminar pagina de tablas
	endpoint := fmt.Sprintf("process/%d", pcb.PID)
	_, err := requests.DeleteHTTP[interface{}](global.KernelConfig.IPMemory, global.KernelConfig.PortMemory, endpoint, nil)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al eliminar proceso en memoria: %+v", err), log.ERROR)
		return
	}

	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: %s ", pcb.PID, pcb.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INVALID_INTERFACE", pcb.PID), log.INFO)
}

func BlockToReady(pcb *model.PCB) {
	// Saco de block cuando termino la IO
	global.MutexBlockState.Lock()
	global.BlockedState.Remove(global.BlockedState.Front())
	global.MutexBlockState.Unlock()

	pcb.State = "READY"

	if global.KernelConfig.PlanningAlgorithm == "VRR" && pcb.RemainingQuantum > 0 {
		global.MutexReadyPlus.Lock()
		global.ReadyPlus.PushBack(pcb)
		global.MutexReadyPlus.Unlock()

		global.Logger.Log(fmt.Sprintf("PID: %d - Bloqueado a Ready Plus", pcb.PID), log.DEBUG)

	} else {
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(pcb)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("PID: %d - Bloqueado normal", pcb.PID), log.DEBUG)
	}

	if pcb.DisplaceReason == "QUANTUM" {
		pcb.RemainingQuantum = global.KernelConfig.Quantum
	}
}

func CheckIfExist(name string) bool {
	_, Ioexist := global.IoMap[name]
	global.Logger.Log(fmt.Sprintf("%s existe", name), log.DEBUG)
	return Ioexist
}
func CheckIfIsValid(name, instruccion string) bool {
	validInstructions := acceptedInstructions[global.IoMap[name].Type]
	global.Logger.Log(fmt.Sprintf("Instrucciones validas: %+v", validInstructions), log.DEBUG)
	for _, ins := range validInstructions {
		if instruccion == ins {
			global.Logger.Log(fmt.Sprintf("%s puede ejecutar %s", name, instruccion), log.DEBUG)
			return true
		}
	}
	global.Logger.Log(fmt.Sprintf("%s NO puede ejecutar %s", name, instruccion), log.DEBUG)

	return false
}
package resources

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Wait(Pcb *model.PCB) {
	resource, exist := global.ResourceMap[Pcb.Instruction.Parameters[0]]

	if !exist {
		global.Logger.Log(fmt.Sprintf("Nombre del recurso: %s no existe", Pcb.Instruction.Parameters[0]), log.DEBUG)
		utils.PCBtoExit(Pcb)
		global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INVALID_RESOURCE", Pcb.PID), log.INFO)
		return
	}

	resource.Count -= 1
	resource.PidList = append(resource.PidList, Pcb.PID)

	global.PIDResourceMap[Pcb.PID] = append(global.PIDResourceMap[Pcb.PID], resource.Name)

	global.Logger.Log(fmt.Sprintf("Recurso: %s - Cantidad instancias: %d", resource.Name, resource.Count), log.DEBUG)

	if resource.Count < 0 {
		resource.MutexList.Lock()
		Pcb.State = "BLOCKED"
		resource.BlockedList.PushBack(Pcb)
		resource.MutexList.Unlock()

		global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: BLOCK", Pcb.PID), log.INFO)

		//poner en listar procesos
		global.Logger.Log(fmt.Sprintf("Bloqueo proceso: %d", Pcb.PID), log.DEBUG)
	} else if Pcb.DisplaceReason == "QUANTUM" {
		Pcb.RemainingQuantum = global.KernelConfig.Quantum
		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(Pcb)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("Envio PID %d ultimo a Ready", Pcb.PID), log.DEBUG)

		array := longterm.ConvertListToArray(global.ReadyState)
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)

	} else {
		global.MutexReadyState.Lock()
		global.ReadyState.PushFront(Pcb)
		global.MutexReadyState.Unlock()
		global.Logger.Log(fmt.Sprintf("Envio PID %d primero a Ready", Pcb.PID), log.DEBUG)

		array := longterm.ConvertListToArray(global.ReadyState)
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)

	}
	global.SemReadyList <- struct{}{}
}

func Signal(PcbExec *model.PCB) {
	resource, exist := global.ResourceMap[PcbExec.Instruction.Parameters[0]]

	if !exist {
		global.Logger.Log(fmt.Sprintf("Nombre del recurso: %s no existe", PcbExec.Instruction.Parameters[0]), log.DEBUG)
		utils.PCBtoExit(PcbExec)
		global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: INVALID_RESOURCE", PcbExec.PID), log.INFO)
		return
	}

	resource.Count += 1
	global.Logger.Log(fmt.Sprintf("%s %d", resource.Name, resource.Count), log.DEBUG)

	if resource.BlockedList.Len() > 0 {
		resource.MutexList.Lock()
		PCBBlock := resource.BlockedList.Front().Value.(*model.PCB)
		resource.BlockedList.Remove(resource.BlockedList.Front())
		PCBBlock.State = "Ready"
		resource.MutexList.Unlock()

		global.MutexReadyState.Lock()
		global.ReadyState.PushBack(PCBBlock)
		global.MutexReadyState.Unlock()

		global.Logger.Log(fmt.Sprintf("Envio PID %d al fondo de ready", PCBBlock.PID), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: READY", PCBBlock.PID), log.INFO)

		array := longterm.ConvertListToArray(global.ReadyState)
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)

		global.SemReadyList <- struct{}{}
	}

	// value := checkInArray(resource.PidList, PcbExec.PID)
	value := checkResourcePID(global.PIDResourceMap[PcbExec.PID], resource.Name)

	if value != -1 {
		// resource.PidList = removeAt(resource.PidList, value)
		global.PIDResourceMap[PcbExec.PID] = removeAtString(global.PIDResourceMap[PcbExec.PID], value)
	}
	if PcbExec.DisplaceReason == "QUANTUM" {
		// PcbExec.RemainingQuantum=global.KernelConfig.Quantum
		// global.MutexReadyState.Lock()
		// global.ReadyState.PushBack(PcbExec)
		// global.MutexReadyState.Unlock()
		// global.Logger.Log(fmt.Sprintf("Envio PID %d ultimo a ready", PcbExec.PID), log.DEBUG)

		utils.PCBExectoReady(PcbExec)

	} else {
		global.MutexReadyState.Lock()
		global.ReadyState.PushFront(PcbExec)
		global.MutexReadyState.Unlock()

		global.Logger.Log(fmt.Sprintf("Envio PID %d primero a Ready", PcbExec.PID), log.DEBUG)

		array := longterm.ConvertListToArray(global.ReadyState)
		global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)

		global.SemReadyList <- struct{}{}
	}
	// global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: READY", PcbExec.PID), log.INFO)

}

// func checkInArray(resourcesIdds []int, pid int) int {
// 	for i, value := range resourcesIdds {
// 		if value == pid {
// 			return i
// 		}
// 	}
// 	return -1
// }
// func removeAt(slice []int, index int) []int {
// 	return append(slice[:index], slice[index+1:]...)
// }

func checkResourcePID(resourceName []string, name string) int {
	for i, value := range resourceName {
		if value == name {
			return i
		}
	}
	return -1
}

func removeAtString(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

package utils

import (
	"container/list"
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/block"
	"github.com/sisoputnfrba/tp-golang/kernel/internal/longterm"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// Busca en todas las listas el PID
func FindProcessInList(pid int) *model.PCB {
	queues := []*list.List{
		global.NewState,
		global.ReadyState,
		global.ExecuteState,
		global.BlockedState,
		global.ExitState,
	}

	for _, resource := range global.ResourceMap {
		queues = append(queues, resource.BlockedList)
	}

	for _, queue := range queues {
		pcb := findProcess(pid, queue)
		if pcb != nil {
			return pcb
		}
	}

	return nil
}

func findProcess(pid int, list *list.List) *model.PCB {
	for e := list.Front(); e != nil; e = e.Next() {
		pcb := e.Value.(*model.PCB)
		if pid == pcb.PID {
			return pcb
		}
	}

	return nil
}

func GetAllProcess() []ProcessState {

	var allProcesses []ProcessState
	queues := []*list.List{
		global.NewState,
		global.ReadyState,
		global.ExecuteState,
		global.BlockedState,
		global.ExitState,
		global.ReadyPlus,
	}

	for _, resource := range global.ResourceMap {
		queues = append(queues, resource.BlockedList)
	}

	for _, queue := range queues {
		for e := queue.Front(); e != nil; e = e.Next() {
			pcb := e.Value.(*model.PCB)
			allProcesses = append(allProcesses, ProcessState{
				PID:   pcb.PID,
				State: pcb.State,
			},
			)
		}
	}

	return allProcesses
}

func RemoveProcessByPID(pid int) bool {

	queues := []*list.List{
		global.NewState,
		global.BlockedState,
		global.ExecuteState,
		global.ReadyState,
		global.ExitState,
	}

	for _, resource := range global.ResourceMap {
		queues = append(queues, resource.BlockedList)
	}

	for _, queue := range queues {
		for e := queue.Front(); e != nil; e = e.Next() {
			pcb := e.Value.(*model.PCB)

			if pcb.PID == pid {

				if pcb.State == "EXEC" {
					InterruptCPU("REMOVED")
				} else {
					queue.Remove(e)
				}
				freeResource(pcb)
				PCBtoExit(pcb)
				return true
			}
		}
	}

	return false
}

func PCBToCPU(pcb *model.PCB) (*model.PCB, error) {
	pcb.State = "EXEC"
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: READY - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	resp, err := requests.PutHTTPwithBody[*model.PCB, model.PCB](
		global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "dispatch", pcb)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PCBtoExit(pcb *model.PCB) {
	pcb.State = "EXIT"
	global.MutexExitState.Lock()
	global.ExitState.PushBack(pcb)
	global.MutexExitState.Unlock()

	//LOG CAMBIO DE ESTADO
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)
	<-global.SemMulti

	// Request a memoria para eliminar pagina de tablas
	endpoint := fmt.Sprintf("process/%d", pcb.PID)
	_, err := requests.DeleteHTTP[interface{}](global.KernelConfig.IPMemory, global.KernelConfig.PortMemory, endpoint, nil)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al eliminar proceso en memoria: %+v", err), log.ERROR)
		return
	}
	global.Logger.Log(fmt.Sprintf("Multi: %d despues del exit %d", len(global.SemMulti), pcb.PID), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Proceso %d eliminado en memoria", pcb.PID), log.DEBUG)
}

func PCBtoBlock(pcb *model.PCB) {
	global.Logger.Log(fmt.Sprintf("Remaining quantum: %d", pcb.RemainingQuantum), log.DEBUG)

	// if pcb.DisplaceReason=="QUANTUM" {
	// 	pcb.RemainingQuantum=global.KernelConfig.Quantum
	// }
	pcb.State = "BLOCK"
	global.MutexBlockState.Lock()
	global.BlockedState.PushBack(pcb)
	global.MutexBlockState.Unlock()

	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)
	global.Logger.Log(fmt.Sprintf("PID: %d - Bloqueado por: %s ", pcb.PID, pcb.Instruction.Parameters[0]), log.INFO)

	go block.ProcessToIO(pcb)
}

func PCBReadytoExec() *model.PCB {

	global.MutexReadyState.Lock()
	pcb := global.ReadyState.Front().Value.(*model.PCB)
	global.ReadyState.Remove(global.ReadyState.Front())
	global.MutexReadyState.Unlock()

	// Pasar a execute
	global.MutexExecuteState.Lock()
	global.ExecuteState.PushBack(pcb)
	global.MutexExecuteState.Unlock()
	return pcb
}

func PCBExectoReady(pcb *model.PCB) {
	//se guarda en ready
	pcb.State = "READY"
	//LOG CAMBIO DE ESTADO
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	//LOG COLA A READY CHEQUEAR EN ESTE CASO

	//LOG FIN DE QUANTUM
	global.Logger.Log(fmt.Sprintf("PID: %d - Desalojado por fin de Quantum ", pcb.PID), log.INFO)
	pcb.RemainingQuantum = global.KernelConfig.Quantum

	global.MutexReadyState.Lock()
	global.ReadyState.PushBack(pcb)
	global.MutexReadyState.Unlock()

	array := longterm.ConvertListToArray(global.ReadyState)
	global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)
	global.SemReadyList <- struct{}{}
}

func PCBExectoReadyVRR(pcb *model.PCB) {
	//se guarda en ready
	pcb.State = "READY"
	//LOG CAMBIO DE ESTADO
	global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: EXEC - Estado Actual: %s", pcb.PID, pcb.State), log.INFO)

	//LOG COLA A READY CHEQUEAR EN ESTE CASO

	//LOG FIN DE QUANTUM
	global.Logger.Log(fmt.Sprintf("PID: %d - Desalojado por fin de Quantum ", pcb.PID), log.INFO)

	global.MutexReadyPlus.Lock()
	global.ReadyPlus.PushBack(pcb)
	global.MutexReadyPlus.Unlock()

	array := longterm.ConvertListToArray(global.ReadyPlus)
	global.Logger.Log(fmt.Sprintf("Cola Ready : %v", array), log.INFO)

	global.SemReadyList <- struct{}{}
}

func VrrPCBtoEXEC() *model.PCB {
	global.MutexReadyPlus.Lock()
	pcb := global.ReadyPlus.Front().Value.(*model.PCB)
	global.ReadyPlus.Remove(global.ReadyPlus.Front())
	global.MutexReadyPlus.Unlock()

	// Pasar a execute
	global.MutexExecuteState.Lock()
	global.ExecuteState.PushBack(pcb)
	global.MutexExecuteState.Unlock()
	return pcb
}

func freeResource(pcb *model.PCB) {
	listResourceNamesPIDS := global.PIDResourceMap[pcb.PID]

	if len(listResourceNamesPIDS) == 0 {
		return
	}
	// [RA, RB, RA]

	for _, resourceName := range listResourceNamesPIDS {
		actualResource := global.ResourceMap[resourceName]
		actualResource.Count++
		index := checkResourcePID(listResourceNamesPIDS, resourceName)
		global.PIDResourceMap[pcb.PID] = removeAtString(global.PIDResourceMap[pcb.PID], index)

		if actualResource.BlockedList.Len() > 0 {
			actualResource.MutexList.Lock()
			PCBBlock := actualResource.BlockedList.Front().Value.(*model.PCB)
			actualResource.BlockedList.Remove(actualResource.BlockedList.Front())
			PCBBlock.State = "Ready"
			actualResource.MutexList.Unlock()

			global.MutexReadyState.Lock()
			global.ReadyState.PushBack(PCBBlock)
			global.MutexReadyState.Unlock()

			global.Logger.Log(fmt.Sprintf("Envio PID %d al fondo de ready", PCBBlock.PID), log.DEBUG)
			global.Logger.Log(fmt.Sprintf("PID: %d - Estado Anterior: BLOCK - Estado Actual: READY", PCBBlock.PID), log.INFO)

			global.SemReadyList <- struct{}{}
		}

		global.Logger.Log(fmt.Sprintf("RECURSO: %+v", actualResource), log.DEBUG)
	}
}

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

func InterruptCPU(reason string) error {

	interruptReason := InterruptReason{
		Reason: reason,
	}

	global.Logger.Log(fmt.Sprintf("Reason Struct: %+v", interruptReason), log.DEBUG)

	_, err := requests.PutHTTPwithBody[InterruptReason, interface{}](global.KernelConfig.IPCPU, global.KernelConfig.PortCPU, "interrupt", interruptReason)

	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al enviar la interrupción: %+v", err), log.ERROR)
		return err
	}

	return nil
}
package algorithm

import (
	"fmt"
	"strings"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/resources"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func RoundRobbin() {
	global.Logger.Log(fmt.Sprintf("Semaforo de SemReadyList INICIO: %d", len(global.SemReadyList)), log.DEBUG)

	for {

		<-global.SemReadyList
		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON ROUND ROBIN", log.DEBUG)
			<-global.SemExecute
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			pcb := utils.PCBReadytoExec()
			// Enviar a execute
			updateChan := make(chan *model.PCB)
			InterruptTimer := make(chan int, 1)

			go DisplaceFunction(InterruptTimer, pcb)

			go func() {
				global.SemInterrupt <- 0
				updatePCB, _ = utils.PCBToCPU(pcb)

				updateChan <- updatePCB
			}()

			updatePCB = <-updateChan
			//LOG CAMBIO DE ESTADO
			global.Logger.Log(fmt.Sprintf("Recibi de CPU: %+v", updatePCB), log.DEBUG)

			if !global.WorkingPlani {
				global.Logger.Log("Bloqueo plani", log.DEBUG)
				global.SemStopPlani <- 0
				global.WorkingPlani = true
				global.Logger.Log("Desbloqueo plani", log.DEBUG)
			}

			// Sacar de execute
			global.MutexExecuteState.Lock()
			global.ExecuteState.Remove(global.ExecuteState.Front())
			global.MutexExecuteState.Unlock()

			if updatePCB.Instruction.Operation == "EXIT" {
				InterruptTimer <- 0
				//DisplaceChan <-updatePCB
				utils.PCBtoExit(updatePCB)
				global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: SUCCESS ", pcb.PID), log.INFO)
			}

			if updatePCB.DisplaceReason == "FAILED RESIZE" {
				utils.PCBtoExit(updatePCB)
				global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: OUT_OF_MEMORY", pcb.PID), log.INFO)
			}

			if updatePCB.DisplaceReason == "BLOCKED" {
				InterruptTimer <- 0
				DisplaceChan <- updatePCB
				utils.PCBtoBlock(updatePCB)
			} else if updatePCB.DisplaceReason == "QUANTUM" && updatePCB.Instruction.Operation != "EXIT" {
				if updatePCB.Instruction.Operation == "SIGNAL" {
					resource.Signal(updatePCB)
				} else if updatePCB.Instruction.Operation == "WAIT" {
					resource.Wait(updatePCB)
				} else if strings.Contains(updatePCB.Instruction.Operation, "IO") {
					utils.PCBtoBlock(updatePCB)
				} else {
					utils.PCBExectoReady(updatePCB)
				}
			}

			if updatePCB.DisplaceReason == "WAIT" {
				InterruptTimer <- 0
				DisplaceChan <- updatePCB
				resource.Wait(updatePCB)
			}
			if updatePCB.DisplaceReason == "SIGNAL" {
				InterruptTimer <- 0
				DisplaceChan <- updatePCB
				resource.Signal(updatePCB)
			}
		}

		<-global.SemExecute
	}
}

func DisplaceFunction(InterruptTimer chan int, OldPcb *model.PCB) {
	<-global.SemInterrupt
	quantumTime := time.Duration(OldPcb.RemainingQuantum) * time.Millisecond
	timer := time.NewTimer(quantumTime)
	defer timer.Stop()
	startTime := time.Now()

	select {
	case <-timer.C:
		global.Logger.Log(fmt.Sprintf("PID: %d Displace - Termino timer.C", OldPcb.PID), log.DEBUG)
		utils.InterruptCPU("QUANTUM")
	case <-InterruptTimer:
		timer.Stop()
		pcb := <-DisplaceChan

		// Transformar el tiempo a segundos para redondearlo y despues pasarlo a ms
		// Asi uso los ms en la PCB
		if pcb.Instruction.Operation == "WAIT" || pcb.Instruction.Operation == "SIGNAL" {
			remainingMillisRounded := utils.TimeCalc(startTime, quantumTime, pcb)
			if remainingMillisRounded > 0 {
				pcb.RemainingQuantum = remainingMillisRounded
			} else {
				pcb.RemainingQuantum = global.KernelConfig.Quantum
			}
		}
	}

}

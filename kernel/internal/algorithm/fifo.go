package algorithm

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	resource "github.com/sisoputnfrba/tp-golang/kernel/internal/resources"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

var updatePCB *model.PCB

func Fifo() {
	global.Logger.Log("Arranca FIFO", log.DEBUG)

	for {

		global.Logger.Log(fmt.Sprintf("LOG ANTES DE SEMREADYLIST largo %d", len(global.SemReadyList)), log.DEBUG)
		<-global.SemReadyList

		global.SemExecute <- 0

		if !global.WorkingPlani {
			global.Logger.Log("TERMINO CON FIFO", log.DEBUG)
			<-global.SemExecute
			break
		}

		if global.ReadyState.Len() != 0 {
			global.Logger.Log(fmt.Sprintf("PCB a execute: %+v", global.ReadyState.Front().Value), log.DEBUG)

			pcb := utils.PCBReadytoExec()

			updateChan := make(chan *model.PCB)
			go func() {
				updatePCB, _ = utils.PCBToCPU(pcb)
				updateChan <- updatePCB
			}()
			updatePCB = <-updateChan
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

			// EXIT - Agregar a exit
			if updatePCB.DisplaceReason == "EXIT" {
				utils.PCBtoExit(updatePCB)
				global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: SUCCESS ", pcb.PID), log.INFO)
			}

			if updatePCB.DisplaceReason == "FAILED RESIZE" {
				utils.PCBtoExit(updatePCB)
				global.Logger.Log(fmt.Sprintf("Finaliza el proceso %d - Motivo: OUT_OF_MEMORY", pcb.PID), log.INFO)
			}

			// Agregar a block
			if updatePCB.DisplaceReason == "BLOCKED" {
				utils.PCBtoBlock(updatePCB)
			}
			if updatePCB.DisplaceReason == "WAIT" {
				resource.Wait(updatePCB)
			}
			if updatePCB.DisplaceReason == "SIGNAL" {
				resource.Signal(updatePCB)
			}

		}

		<-global.SemExecute
	}
}

package internal

import (
	"fmt"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/cpu/internal/execute"
	internal "github.com/sisoputnfrba/tp-golang/cpu/internal/fetch"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func Dispatch(pcb *model.PCB) (*model.PCB, error) {
	global.Logger.Log(fmt.Sprintf("Recibi PCB %+v", pcb), log.DEBUG)

	global.ExecuteMutex.Lock()
	global.Execute = true
	global.InterruptReason = ""
	global.ExecuteMutex.Unlock()
	pcb.DisplaceReason = ""

	for global.Execute {

		instruction, err := internal.Fetch(pcb)
		if err != nil {
			return nil, err
		}

		exec_result := execute.Execute(pcb, instruction)

		if !global.Execute {
			if pcb.DisplaceReason != "FAILED RESIZE" {
				pcb.DisplaceReason = global.InterruptReason
			}
			pcb.RemainingQuantum = 0
			global.Logger.Log(fmt.Sprintf("Salgo por Quantum - PCB Actualizada %+v", pcb), log.DEBUG)
			return pcb, nil
		}
		if exec_result == execute.RETURN_CONTEXT {
			global.Execute = false
		}
	}
	DisplaceReason(pcb)
	global.Logger.Log(fmt.Sprintf("PCB Actualizada %+v", pcb), log.DEBUG)
	return pcb, nil
}

func DisplaceReason(pcb *model.PCB) {
	if strings.Contains(pcb.Instruction.Operation, "IO") {
		pcb.DisplaceReason = "BLOCKED"
	} else if pcb.Instruction.Operation == "EXIT" {
		pcb.DisplaceReason = "EXIT"
	} else if pcb.Instruction.Operation == "WAIT" {
		pcb.DisplaceReason = "WAIT"
	} else if pcb.Instruction.Operation == "SIGNAL" {
		pcb.DisplaceReason = "SIGNAL"
	} else if pcb.Instruction.Operation == "RESIZE" {
		pcb.DisplaceReason = "FAILED RESIZE"
	}

}

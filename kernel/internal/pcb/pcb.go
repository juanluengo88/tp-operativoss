package pcb

import (
	"github.com/sisoputnfrba/tp-golang/kernel/global"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func CreateNewProcess() *model.PCB {
	return &model.PCB{
		PID:              global.GetNextPID(),
		State:            "NEW",
		Quantum:          global.KernelConfig.Quantum,
		RemainingQuantum: global.KernelConfig.Quantum,
	}
}

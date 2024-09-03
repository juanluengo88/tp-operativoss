package internal

import (
	"fmt"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

func Fetch(pcb *model.PCB) (*model.Instruction, error) {
	// instruction, err := getInstruction(pcb.PID, pcb.PC)
	instruction, err := getInstruction(pcb.PID, pcb.PC)
	if err != nil {
		global.Logger.Log("Error al obtener la instruccion: "+err.Error(), log.ERROR)
		return nil, err
	}

	// if instruction == nil {
	// 	global.Logger.Log("NO HAY MAS INSTRUCCIONES", log.DEBUG)
	// 	return nil, nil
	// }

	global.Logger.Log(fmt.Sprintf("PID: %d - FETCH - Program Counter: %d", pcb.PID, pcb.PC), log.INFO)
	pcb.PC++
	return instruction, err
}

func getInstruction(id, pc int) (*model.Instruction, error) {
	path := fmt.Sprintf("process/%d", id)
	proccesInstruction := model.ProcessInstruction{
		Pid: id,
		Pc:  pc,
	}
	raw_instruction, err := requests.PutHTTPwithBody[model.ProcessInstruction, string](
		global.CPUConfig.IPMemory,
		global.CPUConfig.PortMemory,
		path,
		proccesInstruction,
	)

	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al solicitar instrucción desde memoria: %v", err), log.ERROR)
		return nil, err
	}

	if raw_instruction == nil {
		return nil, nil
	}

	sliceInstruction := strings.Fields(*raw_instruction)

	instruction := &model.Instruction{
		Operation:  sliceInstruction[0],
		Parameters: sliceInstruction[1:],
	}

	return instruction, nil
}

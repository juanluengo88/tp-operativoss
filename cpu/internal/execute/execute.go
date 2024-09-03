package execute

import (
	"fmt"
	"math"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	"github.com/sisoputnfrba/tp-golang/cpu/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// TODO: IO_GEN_SLEEP

const (
	CONTINUE       = 0
	RETURN_CONTEXT = 1
)

type Estructura_resize struct {
	Pid       int `json:"pid"`
	NumFrames int `json:"frames"`
}

type Response struct {
	Respuesta string `json:"respuesta"`
}

type CopyStringStruct struct {
	Pid           int   `json:"pid"`
	Length        int   `json:"length"`
	NumFramesFrom []int `json:"numframeRead"`
	OffsetRead    int   `json:"offsetRead"`
	NumFramesTo   []int `json:"numframeCopy"`
	OffsetTo      int   `json:"offsetCopy"`
}

var result = 0

// Ejecuto -> sumo PC en dispatch?
func Execute(pcb *model.PCB, instruction *model.Instruction) int {

	switch instruction.Operation {
	case "SET":
		set(pcb, instruction)
		result = CONTINUE
	case "SUM":
		sum(pcb, instruction)
		result = CONTINUE
	case "SUB":
		sub(pcb, instruction)
		result = CONTINUE
	case "JNZ":
		jnz(pcb, instruction)
		result = CONTINUE
	case "WAIT":
		result = RETURN_CONTEXT
	case "SIGNAL":
		result = RETURN_CONTEXT
	case "MOV_IN":
		mov_in(pcb, instruction)
		result = CONTINUE
	case "MOV_OUT":
		mov_out(pcb, instruction)
		result = CONTINUE
	case "RESIZE":
		result = resize(pcb, instruction)
	case "COPY_STRING":
		result = copyString(pcb, instruction)
	case "IO_GEN_SLEEP":
		result = RETURN_CONTEXT
	case "IO_STDIN_READ", "IO_STDOUT_WRITE":
		ioStd(pcb, instruction)
		result = RETURN_CONTEXT
	case "IO_FS_CREATE", "IO_FS_DELETE":
		result = RETURN_CONTEXT
	case "IO_FS_TRUNCATE":
		ioFsTrunc(pcb, instruction)
		result = RETURN_CONTEXT
	case "IO_FS_WRITE", "IO_FS_READ":
		ioFs(pcb, instruction)
		result = RETURN_CONTEXT
	case "EXIT":
		result = RETURN_CONTEXT
	}

	global.Logger.Log(
		fmt.Sprintf("PID: %d - Ejecutando: %s - %+v",
			pcb.PID,
			instruction.Operation,
			instruction.Parameters,
		),
		log.INFO)

	pcb.Instruction = *instruction

	return result
}

func set(pcb *model.PCB, instruction *model.Instruction) {
	value, _ := strconv.Atoi(instruction.Parameters[1])
	setRegister(instruction.Parameters[0], value, pcb)
}

func sum(pcb *model.PCB, instruction *model.Instruction) {

	destinationValue := getRegister(instruction.Parameters[0], pcb)
	sourceValue := getRegister(instruction.Parameters[1], pcb)
	destinationValue = destinationValue + sourceValue
	setRegister(instruction.Parameters[0], destinationValue, pcb)
}

func sub(pcb *model.PCB, instruction *model.Instruction) {

	destinationValue := getRegister(instruction.Parameters[0], pcb)
	sourceValue := getRegister(instruction.Parameters[1], pcb)
	destinationValue = destinationValue - sourceValue
	setRegister(instruction.Parameters[0], destinationValue, pcb)
}

func jnz(pcb *model.PCB, instruction *model.Instruction) {
	value := getRegister(instruction.Parameters[0], pcb)
	if value != 0 {
		newPC, _ := strconv.Atoi(instruction.Parameters[1])
		pcb.PC = newPC
	}
}

func getRegister(register string, pcb *model.PCB) int {
	switch register {
	case "AX":
		return pcb.Registers.AX
	case "BX":
		return pcb.Registers.BX
	case "CX":
		return pcb.Registers.CX
	case "DX":
		return pcb.Registers.DX
	case "EAX":
		return pcb.Registers.EAX
	case "EBX":
		return pcb.Registers.EBX
	case "ECX":
		return pcb.Registers.ECX
	case "EDX":
		return pcb.Registers.EDX
	case "SI":
		return pcb.Registers.SI
	case "DI":
		return pcb.Registers.DI
	default:
		return -1
	}
}

func setRegister(register string, value int, pcb *model.PCB) {
	switch register {
	case "AX":
		pcb.Registers.AX = value
	case "BX":
		pcb.Registers.BX = value
	case "CX":
		pcb.Registers.CX = value
	case "DX":
		pcb.Registers.DX = value
	case "EAX":
		pcb.Registers.EAX = value
	case "EBX":
		pcb.Registers.EBX = value
	case "ECX":
		pcb.Registers.ECX = value
	case "EDX":
		pcb.Registers.EDX = value
	case "SI":
		pcb.Registers.SI = value
	case "DI":
		pcb.Registers.DI = value
	case "PC":
		pcb.PC = value
	}
}

func mov_in(pcb *model.PCB, instruction *model.Instruction) {
	dataValue := instruction.Parameters[0]
	LogAdress := getRegister(instruction.Parameters[1], pcb)

	size := internal.GetLength(dataValue)
	SendStruct := internal.CreateAdress(size, LogAdress, pcb.PID, getRegister(dataValue, pcb))

	// put a memoria para que devuelva el valor solicitado

	resp, err := requests.PutHTTPwithBody[internal.MemStruct, int](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "memIn", SendStruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}
	global.Logger.Log(fmt.Sprintf("Resp %+v", *resp), log.DEBUG)
	setRegister(dataValue, *resp, pcb)
	global.Logger.Log(fmt.Sprintf("PID: %d - Acción: LEER - Dirección Física: %d %d - Valor: %d", pcb.PID, SendStruct.NumFrames[0], SendStruct.Offset, *resp), log.INFO)
}

func mov_out(pcb *model.PCB, instruction *model.Instruction) {
	// Dato a escribir
	dataRegister := instruction.Parameters[1]
	dataValue := getRegister(dataRegister, pcb)
	global.Logger.Log(fmt.Sprintf("Registro %s - Valor: %d", dataRegister, dataValue), log.DEBUG)
	// Direccion donde quiero escribir
	LogAdress := getRegister(instruction.Parameters[0], pcb)

	size := internal.GetLength(dataRegister)
	SendStruct := internal.CreateAdress(size, LogAdress, pcb.PID, dataValue)

	global.Logger.Log(fmt.Sprintf("Struct de direccion creada %+v", SendStruct), log.DEBUG)
	// put a memoria para que guarde

	_, err := requests.PutHTTPwithBody[internal.MemStruct, interface{}](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "memOut", SendStruct)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.INFO)
		panic(1)
		// TODO: falta que memoria vea si puede escribir o no (?)
	}
	global.Logger.Log(fmt.Sprintf("PID: %d - Acción: ESCRIBIR - Dirección Física: %d %d - Valor: %d", pcb.PID, SendStruct.NumFrames[0], SendStruct.Offset, SendStruct.Content), log.INFO)
}

func resize(pcb *model.PCB, instruction *model.Instruction) int {
	newSize, _ := strconv.Atoi(instruction.Parameters[0])

	ceilSize := math.Ceil(float64(newSize) / float64(global.CPUConfig.Page_size))

	global.Logger.Log(fmt.Sprintf("Tamanio %f - Valor: %d", (float64(newSize)/float64(global.CPUConfig.Page_size)), int(ceilSize)), log.DEBUG)

	estructura_resize := Estructura_resize{
		Pid:       pcb.PID,
		NumFrames: int(ceilSize),
	}

	global.Logger.Log(fmt.Sprintf("Struct a resize: %+v", estructura_resize), log.DEBUG)

	res, err := requests.PutHTTPwithBody[Estructura_resize, interface{}](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "resize", estructura_resize)
	// global.Logger.Log(fmt.Sprintf("STATUS CODE DSP RESIZE: %d", resp.StatusCode), log.DEBUG)
	if err != nil {
		global.Logger.Log("Out of memory: "+ err.Error(), log.DEBUG)
		pcb.DisplaceReason = "FAILED RESIZE"
		return RETURN_CONTEXT
	}

	global.Logger.Log(fmt.Sprintf("Memoria respondio al resize: %+v", res), log.DEBUG)

	return CONTINUE
}

func copyString(pcb *model.PCB, instruction *model.Instruction) int {

	size, _ := strconv.Atoi(instruction.Parameters[0])
	copyFrom := getRegister("SI", pcb)
	copyTo := getRegister("DI", pcb)

	global.Logger.Log(fmt.Sprintf("Tam: %d, From: %d to %d", size, copyFrom, copyTo), log.DEBUG)

	SIPhysicalAddress := internal.CreateAdress(size, copyFrom, pcb.PID, "")
	DIPhysicalAddress := internal.CreateAdress(size, copyTo, pcb.PID, "")

	global.Logger.Log(fmt.Sprintf("SI - Frames: %+v - Offset: %d", SIPhysicalAddress.NumFrames, SIPhysicalAddress.Offset), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("DI - Frames: %+v - Offset: %d", DIPhysicalAddress.NumFrames, DIPhysicalAddress.Offset), log.DEBUG)

	copyString := CopyStringStruct{
		Pid:           pcb.PID,
		Length:        size,
		NumFramesFrom: SIPhysicalAddress.NumFrames,
		OffsetRead:    SIPhysicalAddress.Offset,
		NumFramesTo:   DIPhysicalAddress.NumFrames,
		OffsetTo:      DIPhysicalAddress.Offset,
	}

	global.Logger.Log(fmt.Sprintf("Struct a memoria: %+v", copyString), log.DEBUG)

	_, err := requests.PutHTTPwithBody[CopyStringStruct, any](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "copy_string", copyString)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al enviar la estructura copyString a memoria: %s", err.Error()), log.INFO)
		return RETURN_CONTEXT
	}

	return CONTINUE
}


func ioStd(pcb *model.PCB, instruction *model.Instruction) {

	// interfaceName := instruction.Parameters[0]
	logicAddress := getRegister(instruction.Parameters[1], pcb)
	size := getRegister(instruction.Parameters[2], pcb)

	physicalAddress := internal.CreateAdress(size, logicAddress, pcb.PID, "")

	instruction.NumFrames = physicalAddress.NumFrames
	instruction.Offset = physicalAddress.Offset
	instruction.Size = size
}

func ioFs(pcb *model.PCB, instruction *model.Instruction) {
	logicAddress := getRegister(instruction.Parameters[2], pcb)
	size := getRegister(instruction.Parameters[3], pcb)
	filePointer := getRegister(instruction.Parameters[4], pcb)

	physicalAddress := internal.CreateAdress(size, logicAddress, pcb.PID, "")

	instruction.NumFrames = physicalAddress.NumFrames
	instruction.Offset = physicalAddress.Offset
	instruction.Size = size
	instruction.FSPointer = filePointer
}

func ioFsTrunc(pcb *model.PCB, instruction *model.Instruction) {
	size := getRegister(instruction.Parameters[2], pcb)
	instruction.Size = size
}
package block

import (
	"strconv"

	"github.com/sisoputnfrba/tp-golang/utils/model"
)

type IO interface {
	GetName() string
	GetInstruction() string
}

type IOGen struct {
	Name        string `json:"nombre"`
	Instruction string `json:"instruccion"`
	Time        int    `json:"tiempo"`
	Pid         int    `json:"pid"`
}

func (io IOGen) GetName() string {
	return io.Name
}

func (io IOGen) GetInstruction() string {
	return io.Instruction
}

type IOStd struct {
	Pid         int    `json:"pid"`
	Instruction string `json:"instruccion"`
	Name        string `json:"name"`
	Length      int    `json:"length"`
	NumFrames   []int  `json:"numframe"`
	Offset      int    `json:"offset"`
}

func (io IOStd) GetName() string {
	return io.Name
}

func (io IOStd) GetInstruction() string {
	return io.Instruction
}

type IOFs struct {
	Pid         int    `json:"pid"`
	Instruction string `json:"instruccion"`
	IOName        string `json:"nombre"`
	FileName    string `json:"filename"`
}

func (io IOFs) GetName() string {
	return io.IOName
}

func (io IOFs) GetInstruction() string {
	return io.Instruction
}

type IOFsTruncate struct {
	Pid         int    `json:"pid"`
	Instruction string `json:"instruccion"`
	IOName        string `json:"nombre"`
	FileName    string `json:"filename"`
	Size int `json:"tamanio"`
}

func (io IOFsTruncate) GetName() string {
	return io.IOName
}

func (io IOFsTruncate) GetInstruction() string {
	return io.Instruction
}

type IOFsWR struct {
	Pid            int    `json:"pid"`
	Instruction    string `json:"instruccion"`
	IOName         string `json:"nombre"`
	FileName       string `json:"filename"`
	NumFrames      []int  `json:"numframe"`
	Offset         int    `json:"offset"`
	Size        int    `json:"tamanio"`
	FSPointer int    `json:"punteroArchivo"`
}

func (io IOFsWR) GetName() string {
	return io.IOName
}

func (io IOFsWR) GetInstruction() string {
	return io.Instruction
}

func factoryIO(pcb *model.PCB) IO {
	switch pcb.Instruction.Operation {
	case "IO_GEN_SLEEP":
		time, _ := strconv.Atoi(pcb.Instruction.Parameters[1])
		return IOGen{
			Name:        pcb.Instruction.Parameters[0],
			Instruction: pcb.Instruction.Operation,
			Time:        time,
			Pid:         pcb.PID,
		}

	case "IO_STDIN_READ", "IO_STDOUT_WRITE":
		return IOStd{
			Name:        pcb.Instruction.Parameters[0],
			Instruction: pcb.Instruction.Operation,
			Pid:         pcb.PID,
			Length:      pcb.Instruction.Size,
			NumFrames:   pcb.Instruction.NumFrames,
			Offset:      pcb.Instruction.Offset,
		}
	
	case "IO_FS_CREATE", "IO_FS_DELETE":
		return IOFs{
			Pid: pcb.PID,
			Instruction: pcb.Instruction.Operation,
			IOName: pcb.Instruction.Parameters[0],
			FileName: pcb.Instruction.Parameters[1],
		}

	case "IO_FS_TRUNCATE":
		return IOFsTruncate{
			Pid: pcb.PID,
			Instruction: pcb.Instruction.Operation,
			IOName: pcb.Instruction.Parameters[0],
			FileName: pcb.Instruction.Parameters[1],
			Size: pcb.Instruction.Size,
		}

	case "IO_FS_WRITE", "IO_FS_READ":
		return IOFsWR{
			Pid: pcb.PID,
			Instruction: pcb.Instruction.Operation,
			IOName: pcb.Instruction.Parameters[0],
			FileName: pcb.Instruction.Parameters[1],
			NumFrames: pcb.Instruction.NumFrames,
			Offset: pcb.Instruction.Offset,
			Size: pcb.Instruction.Size,
			FSPointer: pcb.Instruction.FSPointer,
		}
	}
	return nil
}

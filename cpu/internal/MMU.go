package internal

import (
	"fmt"
	"math"

	"github.com/sisoputnfrba/tp-golang/cpu/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
)

// handlers "github.com/sisoputnfrba/tp-golang/cpu/api/handlers"

type MemStruct struct {
	Pid       int   `json:"pid"`
	Content   any   `json:"content"`
	Length    int   `json:"length"`
	NumFrames []int `json:"numframe"`
	Offset    int   `json:"offset"`
}

type Frame struct {
	Pid        int `json:"pid"`
	PageNumber int `json:"page_number"`
}

func translateAddress(logicalAddress int) (int, int) {
	pageSize := global.CPUConfig.Page_size

	pageNumber := logicalAddress / pageSize
	offset := logicalAddress - (pageNumber * pageSize)

	return pageNumber, offset
}

func CreateAdress(size int, logicalAddress int, pid int, Content any) MemStruct {
	global.Logger.Log(fmt.Sprintf("Valor de direccion a traducir %d", logicalAddress), log.DEBUG)

	pageNumber, offset := translateAddress(logicalAddress)

	global.Logger.Log(fmt.Sprintf("Numero de pagina %d - Offset: %d", pageNumber, offset), log.DEBUG)

	address := MemStruct{Pid: pid, Content: Content, Length: size, Offset: offset}

	frame, hit := global.Tlb.Search(pid, pageNumber)
	if hit {
		global.Logger.Log(fmt.Sprintf("TLB HIT: PID %d, Página %d -> Marco %d", pid, pageNumber, frame), log.INFO)
	} else {
		global.Logger.Log(fmt.Sprintf("TLB Miss: PID %d, Página %d", pid, pageNumber), log.INFO)
		frame = consultMemory(pid, pageNumber)
		global.Tlb.AddEntry(pid, pageNumber, frame)
	}
	address.NumFrames = append(address.NumFrames, frame)

	numPages := math.Ceil(float64(offset+size) / float64(global.CPUConfig.Page_size))
	global.Logger.Log(fmt.Sprintf("Paginas necesarias %d", int(numPages)), log.DEBUG)

	for i := 1; i < int(numPages); i++ {
		pageNumber++
		frame, hit = global.Tlb.Search(pid, pageNumber)
		if hit {
			global.Logger.Log(fmt.Sprintf("TLB HIT: PID %d, Página %d -> Marco %d", pid, pageNumber, frame), log.INFO)
		} else {
			global.Logger.Log(fmt.Sprintf("TLB Miss: PID %d, Página %d", pid, pageNumber), log.INFO)
			frame = consultMemory(pid, pageNumber)
			global.Tlb.AddEntry(pid, pageNumber, frame)
		}
		address.NumFrames = append(address.NumFrames, frame)
	}

	global.Logger.Log(fmt.Sprintf("TLB DSP DEL TRANSLATE: %+v", global.Tlb), log.DEBUG)

	return address
}

func consultMemory(pid, pageNumber int) int {
	Page := Frame{Pid: pid, PageNumber: pageNumber}
	frame, err := requests.PutHTTPwithBody[Frame, int](global.CPUConfig.IPMemory, global.CPUConfig.PortMemory, "framenumber", Page)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error: %s", err.Error()), log.DEBUG)
		panic(1)
	}
	global.Logger.Log(fmt.Sprintf("PID: %d - OBTENER MARCO - Página: %d - Marco: %d", pid, pageNumber, *frame), log.INFO)
	return *frame
}

func GetLength(Register string) int {
	switch Register {
	case "AX", "BX", "CX", "DX":
		return 1
	case "EAX", "EBX", "ECX", "EDX", "SI", "SD":
		return 4
	}
	return -1
}

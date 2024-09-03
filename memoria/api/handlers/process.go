package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

type Response struct {
	Respuesta string `json:"respuesta"`
}


func CodeReciever(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	var pPath internal.ProcessPath
	err := serialization.DecodeHTTPBody(r, &pPath)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("Me enviaron %+v", pPath), log.DEBUG)

	//leo el archivo con las instrucciones y las guardo en un array de string
	ListInstructions, err := internal.ReadTxt(pPath.Path)
	if err != nil {
		global.Logger.Log("Error al leer el archivo "+err.Error(), log.ERROR)
		http.Error(w, "Error al leer archivo", http.StatusBadRequest)
		return
	}
	//escribo en el map el pid y su lista de instrucciones y crea tabla de paginas
	internal.InstructionStorage(ListInstructions, pPath.Pid)
	global.Logger.Log(fmt.Sprintf("%+v\n", global.DictProcess), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Pid: %d -Tamaño:%d\n", pPath.Pid, len(global.DictProcess[pPath.Pid].PageTable.Pages)), log.INFO)
	w.WriteHeader(http.StatusOK)
}

func SendInstruction(w http.ResponseWriter, r *http.Request) {
	var ProcessAssets internal.ProcessAssets
	global.Logger.Log(fmt.Sprintf("Me enviaron %+v", ProcessAssets), log.DEBUG)

	err := serialization.DecodeHTTPBody(r, &ProcessAssets)
	if err != nil {
		http.Error(w, "Error al decodear el PC", http.StatusBadRequest)
		return
	}

	Instruction := ProcessAssets.Pc
	//me traigo toda la lista de instrucciones del pid correspondiente
	ListInstructions := global.DictProcess[ProcessAssets.Pid].Instructions
	//de aca en adelante la logica es la misma
	if Instruction > len(ListInstructions) { //esto chequea si la intruccion esta dentro del rango
		global.Logger.Log("out of memory: ", log.ERROR)
		http.Error(w, "out of memory", http.StatusForbidden)
		return
	}
	if Instruction == len(ListInstructions) { //esto chequea que no lea memoria q no le corresponde
		w.WriteHeader(http.StatusNoContent)
		return
	}
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond) //genera el delay response de la respuesta
	err = serialization.EncodeHTTPResponse(w, ListInstructions[Instruction], http.StatusOK)
	if err != nil {
		global.Logger.Log("Error al encodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al encodear el body", http.StatusBadRequest)
		return
	}
}

func DeleteProcess(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	pid, _ := strconv.Atoi(r.PathValue("pid"))
	for i := 0; i < len(global.DictProcess[pid].PageTable.Pages); i++ {
		global.BitMap[global.DictProcess[pid].PageTable.Pages[len(global.DictProcess[pid].PageTable.Pages)-1-i]] = 0
	}
	global.Logger.Log(fmt.Sprintf("Me enviaron proceso %d", pid), log.DEBUG)
	global.DictProcess[pid].PageTable.Pages = global.DictProcess[pid].PageTable.Pages[:0]

	global.Logger.Log(fmt.Sprintf("Se elimino el proceso con el PID : %d ", pid), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Pid: %d -Tamaño:%d\n", pid, len(global.DictProcess[pid].PageTable.Pages)), log.INFO)
	delete(global.DictProcess, pid)

	global.Logger.Log(fmt.Sprintf("Page table %d %+v", pid, global.DictProcess[pid].PageTable), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)

	resp:=fmt.Sprintf("Se borro el proceso %d correctamente",pid)
	serialization.EncodeHTTPResponse(w, resp, 200)
	
}

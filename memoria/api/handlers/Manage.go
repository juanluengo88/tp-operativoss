package handlers

import (
	"fmt"
	"net/http"
	"time"
	global "github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// recibo tamaño en frames
func Resize(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	var Process internal.Resize

	err := serialization.DecodeHTTPBody(r, &Process)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}

	global.Logger.Log(fmt.Sprintf("Me enviaron: %+v", Process), log.DEBUG)
	tablaPag := global.DictProcess[Process.Pid].PageTable.Pages

	if len(tablaPag) < Process.Frames {
		global.Logger.Log(fmt.Sprintf("Frames a aumentar %d", Process.Frames-len(tablaPag)), log.DEBUG)
		for i := 0; i < Process.Frames-len(tablaPag); i++ {

			if internal.AddPage(Process.Pid) == -1 {
				global.Logger.Log("Error memoria llena", log.DEBUG)
				http.Error(w, "Out of memory", http.StatusForbidden)
				return
			}
		}
		global.Logger.Log(fmt.Sprintf("PID: %d - Tamaño Actual: %d- Tamaño a Ampliar: %d", Process.Pid, len(tablaPag), Process.Frames), log.INFO)
		global.Logger.Log(fmt.Sprintf("Page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)

		w.WriteHeader(http.StatusNoContent)
		//Reduzco tamaño
	} else if len(global.DictProcess[Process.Pid].PageTable.Pages) > Process.Frames && Process.Frames != 0 {
		difTam := len(global.DictProcess[Process.Pid].PageTable.Pages) - Process.Frames
		global.Logger.Log(fmt.Sprintf("PID: %d - Tamaño Actual: %d- Tamaño a Reducir: %d", Process.Pid, len(global.DictProcess[Process.Pid].PageTable.Pages), Process.Frames), log.INFO)
		for i := 0; i < difTam; i++ {
			global.BitMap[global.DictProcess[Process.Pid].PageTable.Pages[len(global.DictProcess[Process.Pid].PageTable.Pages)-1-i]] = 0
		}
		global.DictProcess[Process.Pid].PageTable.Pages = global.DictProcess[Process.Pid].PageTable.Pages[:Process.Frames]

		global.Logger.Log(fmt.Sprintf("Page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
		w.WriteHeader(http.StatusNoContent)
	} else if Process.Frames == 0 {

		for i := 0; i < len(global.DictProcess[Process.Pid].PageTable.Pages); i++ {
			global.BitMap[global.DictProcess[Process.Pid].PageTable.Pages[len(global.DictProcess[Process.Pid].PageTable.Pages)-1-i]] = 0
		}
		global.DictProcess[Process.Pid].PageTable.Pages = global.DictProcess[Process.Pid].PageTable.Pages[:0]
		//tablaPag=tablaPag[:0]
		global.Logger.Log("Vaciando tabla de paginas", log.DEBUG)
		global.Logger.Log(fmt.Sprintf("PID: %d - Tamaño Actual: %d- Tamaño a Reducir: %d", Process.Pid, len(global.DictProcess[Process.Pid].PageTable.Pages), Process.Frames), log.INFO)

		global.Logger.Log(fmt.Sprintf("Page table %d %+v", Process.Pid, global.DictProcess[Process.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)

		w.WriteHeader(http.StatusNoContent)
	}
}

func PageTableAccess(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	var PageNumber internal.Page
	err := serialization.DecodeHTTPBody(r, &PageNumber)
	global.Logger.Log(fmt.Sprintf("Me enviaron %+v", PageNumber), log.DEBUG)

	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	frame := internal.GetFrame(PageNumber.PageNumber, PageNumber.Pid)
	if frame == -1 {
		global.Logger.Log("Error no existe la pagina", log.DEBUG)
		http.Error(w, "Invalid page", http.StatusForbidden)
	} else {
		global.Logger.Log(fmt.Sprintf("PID: %d- Pagina: %d - Marco: %d", PageNumber.Pid, PageNumber.PageNumber, frame), log.INFO)

		serialization.EncodeHTTPResponse(w, frame, http.StatusOK)
	}

}

// LEER EN MEMORIA
func MemoryAccessIn(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	var MemoryAccess internal.MemStruct
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("Me enviaron: %+v", MemoryAccess), log.DEBUG)

	MemoryAccess.Content = internal.MemIn(MemoryAccess.NumFrames, MemoryAccess.Offset, MemoryAccess.Pid, MemoryAccess.Length)
	serialization.EncodeHTTPResponse(w, MemoryAccess.Content, http.StatusOK)
	//global.Logger.Log(fmt.Sprintf("PID: %d - Accion: LEER - Direccion fisica: %+v + %d - Tamaño: %d Bytes  A LEER", MemoryAccess.Pid, MemoryAccess.NumFrames,MemoryAccess.Offset, MemoryAccess.Length,), log.INFO)

}

// ESCRIBE EN MEMORIA
func MemoryAccessOut(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	var MemoryAccess internal.MemStruct
	err := serialization.DecodeHTTPBody(r, &MemoryAccess)

	global.Logger.Log(fmt.Sprintf("Me enviaron: %+v", MemoryAccess), log.DEBUG)

	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	if internal.MemOut(MemoryAccess.NumFrames, MemoryAccess.Offset, MemoryAccess.Content, MemoryAccess.Pid, MemoryAccess.Length) {

		global.Logger.Log(fmt.Sprintf("Page table %d %+v", MemoryAccess.Pid, global.DictProcess[MemoryAccess.Pid].PageTable), log.DEBUG)
		global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
		//global.Logger.Log(fmt.Sprintf("PID: %d - Accion: ESCRIBIR - Direccion fisica: %+v + %d - Tamaño: %d Bytes  A ESCRIBIR", MemoryAccess.Pid, MemoryAccess.NumFrames,MemoryAccess.Offset, MemoryAccess.Length,), log.INFO)

		//internal.PrintMemoryTable(global.Memory.Spaces, global.MemoryConfig.PageSize)

		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)

	}

}

func Copy_string(w http.ResponseWriter, r *http.Request) {
	DelayResponse := time.Duration(global.MemoryConfig.DelayResponse)
	time.Sleep(DelayResponse * time.Millisecond)
	var MemoryCopy internal.MemCopyString
	err := serialization.DecodeHTTPBody(r, &MemoryCopy)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó este mensaje : %+v", MemoryCopy), log.DEBUG)

	//LEE
	
	Content := internal.ReadInMemory(MemoryCopy.Length, MemoryCopy.NumFramesRead, MemoryCopy.OffsetRead)
	str := string(Content)

	global.Logger.Log(fmt.Sprintf("Page table %d %+v", MemoryCopy.Pid, global.DictProcess[MemoryCopy.Pid].PageTable), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Se leyó el string %s y se procedera a copiar en la direccion indicada", str), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("PID: %d - Accion: LEER - Direccion fisica: %+v + %d - Tamaño: %d Bytes  A LEER", MemoryCopy.Pid, MemoryCopy.NumFramesRead,MemoryCopy.OffsetRead, MemoryCopy.Length,), log.INFO)

	//COPIA
	global.Logger.Log(fmt.Sprintf("Largo de lo que se va a copiar %d", len(Content)), log.DEBUG)
	internal.WriteInMemory(Content,MemoryCopy.Length,MemoryCopy.NumFramesCopy,MemoryCopy.OffsetCopy)
	//internal.PrintMemoryTable(global.Memory.Spaces, global.MemoryConfig.PageSize)
	global.Logger.Log(fmt.Sprintf("PID: %d - Accion: ESCRIBIR - Direccion fisica: %+v + %d - Tamaño: %d Bytes  A ESCRIBIR", MemoryCopy.Pid, MemoryCopy.NumFramesCopy,MemoryCopy.OffsetCopy, MemoryCopy.Length,), log.INFO)

	resp := "Fue copiado correctamente"

	serialization.EncodeHTTPResponse(w, resp, 200)
	if err != nil {
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}

}

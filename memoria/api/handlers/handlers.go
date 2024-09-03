package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	internal "github.com/sisoputnfrba/tp-golang/memoria/internal"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

// LEE DE MEMORIA
func Stdout_write(w http.ResponseWriter, r *http.Request) {
	var MemoryAccessIO internal.MemStdIO
	err := serialization.DecodeHTTPBody(r, &MemoryAccessIO)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó ésta mensaje: %+v", MemoryAccessIO), log.INFO)
	global.Logger.Log(fmt.Sprintf("Page table %d %+v", MemoryAccessIO.Pid, global.DictProcess[MemoryAccessIO.Pid].PageTable), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
	//global.Logger.Log(fmt.Sprintf("Memoria  %+v", global.Memory), log.DEBUG)
	str := string(internal.ReadInMemory(MemoryAccessIO.Length, MemoryAccessIO.NumFrames, MemoryAccessIO.Offset))
	global.Logger.Log(fmt.Sprintf("PID: %d - Accion: LEER - Direccion fisica: %+v + %d - Tamaño: %d Bytes A LEER", MemoryAccessIO.Pid, MemoryAccessIO.NumFrames,MemoryAccessIO.Offset, MemoryAccessIO.Length,), log.INFO)

	serialization.EncodeHTTPResponse(w, str, 200)

}

// ESCRIBE EN MEMORIA
func Stdin_read(w http.ResponseWriter, r *http.Request) {
	//var estructura estructura_write
	var MemoryAccessIO internal.MemStdIO
	err := serialization.DecodeHTTPBody(r, &MemoryAccessIO)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
	}
	global.Logger.Log(fmt.Sprintf("Me llegó este mensaje : %+v", MemoryAccessIO), log.DEBUG)

	byteArray := []byte(MemoryAccessIO.Content)
	global.Logger.Log(fmt.Sprintf("largo %+v", len(byteArray)), log.DEBUG)

	internal.WriteInMemory(byteArray, MemoryAccessIO.Length, MemoryAccessIO.NumFrames, MemoryAccessIO.Offset)
	str := "Se escribió correctamente"
	global.Logger.Log(fmt.Sprintf("Page table %d %+v", MemoryAccessIO.Pid, global.DictProcess[MemoryAccessIO.Pid].PageTable), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Bit Map  %+v", global.BitMap), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("PID: %d - Accion: ESCRIBIR - Direccion fisica: %+v + %d - Tamaño: %d Bytes A ESCRIBIR", MemoryAccessIO.Pid, MemoryAccessIO.NumFrames,MemoryAccessIO.Offset, MemoryAccessIO.Length,), log.INFO)

	internal.PrintMemoryTable(global.Memory.Spaces, global.MemoryConfig.PageSize)
	serialization.EncodeHTTPResponse(w, str, 200)
	if err != nil {
		http.Error(w, "Error encodeando respuesta", http.StatusInternalServerError)
		return
	}
}

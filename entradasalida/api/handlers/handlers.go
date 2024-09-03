package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/requests"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func Sleep(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.Estructura_sleep
	err := serialization.DecodeHTTPBody[*global.Estructura_sleep](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s>", estructura.Pid, estructura.Instruction), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("a punto de dormir: %+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("durmiendo: %+v", dispositivo), log.DEBUG)

	time.Sleep(time.Duration(estructura.Time*global.IOConfig.UnitWorkTime) * time.Millisecond)

	global.Logger.Log(fmt.Sprintf("terminé de dormir: %+v", dispositivo), log.DEBUG)

	w.WriteHeader(http.StatusNoContent)
	dispositivo.InUse = false
}

func Stdin_read(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOStd

	err := serialization.DecodeHTTPBody[*global.KernelIOStd](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: %d - Operacion: <%s>", estructura.Pid, estructura.Instruction), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	global.Logger.Log(fmt.Sprintf("Ingrese un valor (tamaño máximo %d): ", estructura.Length), log.INFO)

	global.Estructura_actualizada.Pid = estructura.Pid
	global.Estructura_actualizada.NumFrames = estructura.NumFrames
	global.Estructura_actualizada.Offset = estructura.Offset

	// fmt.Scanf("%s", global.Texto)
	reader := bufio.NewReader(os.Stdin)
	global.Texto, _ = reader.ReadString('\n')

	global.Logger.Log("De consola escribi: "+global.Texto, log.DEBUG)

	global.VerificacionTamanio(global.Texto, estructura.Length)

	global.Estructura_actualizada.Length = len(global.Texto) - 1

	global.Logger.Log(fmt.Sprintf("Estructura actualizada para mandar a memoria: %+v", global.Estructura_actualizada), log.DEBUG)

	// PUT a memoria de la estructura
	_, err = requests.PutHTTPwithBody[global.MemStdIO, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdin_read", global.Estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	dispositivo.InUse = false

	w.WriteHeader(http.StatusNoContent)
}

func Stdout_write(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true
	var estructura global.KernelIOStd
	var estructura_actualizada global.MemStdIO
	err := serialization.DecodeHTTPBody[*global.KernelIOStd](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	estructura_actualizada.Pid = estructura.Pid
	estructura_actualizada.NumFrames = estructura.NumFrames
	estructura_actualizada.Offset = estructura.Offset
	estructura_actualizada.Length = estructura.Length

	global.Logger.Log(fmt.Sprintf("Intentando leer con %s", estructura.Name), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// PUT a memoria (le paso un registro y me devuelve el valor)

	resp, err := requests.PutHTTPwithBody[global.MemStdIO, string](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdout_write", estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.ERROR)
		http.Error(w, "Error al enviar a memoria el valor a escribir", http.StatusBadRequest)
		return
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y me devuelva el valor del registro
	}
	global.Logger.Log(fmt.Sprintf("Memoria devolvió este valor: %s", *resp), log.INFO)

	dispositivo.InUse = false

	w.WriteHeader(http.StatusNoContent)
}

func Fs_create(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_CD

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_CD](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("PID: <%d> - Crear Archivo: <%s>", estructura.Pid, estructura.FileName), log.INFO)

	global.Logger.Log(fmt.Sprintf("Estructura: %+v", estructura), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Dispositivo: %+v", dispositivo), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// implementación

	// 1) busco en mi bitmap el primer bloque libre, uso ese dato para asignarlo como initial_block del archivo metadata estructura.Filename

	firstFreeBlock := getFirstFreeBlock()

	// 2) creo el archivo metadata, de nombre estructura.Filename, con size = 0 e initial_block = al valor hallado en 2)
	filename := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + estructura.FileName
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		global.Logger.Log(fmt.Sprint("Error al crear el archivo:", err), log.ERROR)
		return
	}

	defer file.Close()

	// 3) actualizo el bitmap, tanto el slice bytes como el archivo (podría hacerlo en el paso 1))

	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s antes de crear el nuevo archivo: %+v", global.Dispositivo.Name, global.Bitmap), log.DEBUG)
	global.UpdateBitmap(1, firstFreeBlock, 1)
	global.Logger.Log(fmt.Sprintf("Bitmap del FS %s luego de crear el nuevo archivo: %+v", global.Dispositivo.Name, global.Bitmap), log.DEBUG)

	var filestruct global.File

	filestruct.CurrentBlocks = 0
	filestruct.Initial_block = -1
	filestruct.Size = -1
	global.Logger.Log(fmt.Sprintf("Datos del archivo antes de ser creado (%s): %+v ", filename, filestruct), log.DEBUG)
	filestruct.CurrentBlocks = 1
	filestruct.Initial_block = firstFreeBlock
	filestruct.Size = 0
	global.Logger.Log(fmt.Sprintf("Datos del archivo luego de ser creado (%s): %+v ", filename, filestruct), log.DEBUG)
	global.FilesMap[estructura.FileName] = filestruct

	newSizemap := map[string]interface{}{
		"initial_block": filestruct.Initial_block,
		"size":          filestruct.Size,
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(newSizemap)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al encodear el nuevo size en el archivo %s: %s ", filename, err.Error()), log.ERROR)
		http.Error(w, "Error al encodear el nuevo size en el archivo", http.StatusBadRequest)
		return
	}

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
	global.Logger.Log(fmt.Sprintf("FilesMap al crear %+v", global.FilesMap), log.DEBUG)
}

func Fs_delete(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_CD

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_CD](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("PID: <%d> - Eliminar Archivo: <%s>", estructura.Pid, estructura.FileName), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// implementación

	filestruct := global.FilesMap[estructura.FileName]

	// actualizar el bitmap!

	global.Logger.Log(fmt.Sprintf("Bitmap antes de eliminar archivo: %+v", global.Bitmap), log.DEBUG)
	global.UpdateBitmap(0, filestruct.Initial_block, filestruct.CurrentBlocks)
	global.Logger.Log(fmt.Sprintf("Bitmap luego de eliminar archivo: %+v", global.Bitmap), log.DEBUG)

	//actualizar la cerpeta de archivos
	metadatapath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + estructura.FileName

	// Eliminar el archivo metadata
	err = os.Remove(metadatapath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al eliminar el archivo: "+err.Error()), log.ERROR)
		http.Error(w, "Error al eliminar el archivo", http.StatusInternalServerError)
	}

	// Eliminar el elemento del map

	delete(global.FilesMap, estructura.FileName)

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
	global.Logger.Log(fmt.Sprintf("FilesMap al eliminar %+v", global.FilesMap), log.DEBUG)

}

func Fs_truncate(w http.ResponseWriter, r *http.Request) {

	// decodeo el json que me acaba de llegar para los logs obligatorios

	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_Truncate](r, &global.Estructura_truncate)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", global.Estructura_truncate.Pid, global.Estructura_truncate.Instruction+">"), log.INFO)
	
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Truncar Archivo: <%s> - Tamaño: <%d>", global.Estructura_truncate.Pid, global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio), log.INFO)

	global.Logger.Log(fmt.Sprintf("Dispositivo: %+v", dispositivo), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Instrucción: %+v", global.Estructura_truncate), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// obtengo los datos del archivo metadata

	filestruct := global.FilesMap[global.Estructura_truncate.FileName]

	metadatapath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + global.Estructura_truncate.FileName

	metadatafile, err := os.Open(metadatapath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", metadatapath, err.Error()), log.DEBUG)
		http.Error(w, "Error al abrir el archivo", http.StatusBadRequest)
		return
	}

	defer metadatafile.Close()

	global.Logger.Log(fmt.Sprintf("Filestruct recién decodeado: %+v", filestruct), log.DEBUG)

	currentBlocks := global.GetCurrentBlocks(global.Estructura_truncate.FileName)
	freeContiguousBlocks := global.GetFreeContiguousBlocks(global.Estructura_truncate.FileName)
	neededBlocks := global.GetNeededBlocks(global.Estructura_truncate)
	totalFreeBlocks := global.GetTotalFreeBlocks()

	if currentBlocks == neededBlocks {
		global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, neededBlocks)
		global.Logger.Log(fmt.Sprintf("No es necesario truncar pero actualicé el size: %+v", global.Estructura_truncate), log.DEBUG)

	} else if !(totalFreeBlocks >= neededBlocks-currentBlocks) {
		global.Logger.Log(fmt.Sprintf("No es posible agrandar el archivo: %+v", global.Estructura_truncate), log.ERROR)

	} else if currentBlocks > neededBlocks {
		global.Logger.Log(fmt.Sprintf("Trunco a menos %+v", global.Estructura_truncate), log.DEBUG)

		global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, neededBlocks)
		global.PrintBitmap()
		global.UpdateBitmap(0, filestruct.Initial_block+neededBlocks, currentBlocks-neededBlocks)
		global.PrintBitmap()

	} else if neededBlocks-currentBlocks <= freeContiguousBlocks {
		global.Logger.Log(fmt.Sprintf("Trunco a más %+v", global.Estructura_truncate), log.DEBUG)

		global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, neededBlocks)
		global.PrintBitmap()
		global.UpdateBitmap(1, filestruct.Initial_block+currentBlocks, neededBlocks-currentBlocks)
		global.PrintBitmap()

	} else {
		global.Logger.Log(fmt.Sprintf("Es necesario compactar: %+v", global.Estructura_truncate), log.DEBUG)

		// actualizar bitamp y archivos metadata
		compactar(global.Estructura_truncate.FileName, totalFreeBlocks)

		global.PrintBloques()

	}

	global.Logger.Log(fmt.Sprintf("FilesMap al truncar %+v", global.FilesMap), log.DEBUG)

	w.WriteHeader(http.StatusNoContent)
	dispositivo.InUse = false
}

func Fs_write(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_WR
	var estructura_actualizada global.MemStdIO

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_WR](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// implementación

	estructura_actualizada.Pid = estructura.Pid
	estructura_actualizada.NumFrames = estructura.NumFrames
	estructura_actualizada.Offset = estructura.Offset
	estructura_actualizada.Length = estructura.Tamanio

	// hago una request a memoria para obtener un valor

	resp, err := requests.PutHTTPwithBody[global.MemStdIO, string](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdout_write", estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria el valor a escribir %s", err.Error()), log.ERROR)
		http.Error(w, "Error al enviar a memoria el valor a escribir", http.StatusBadRequest)
		return
		// TODO: memoria falta que entienda el mensaje (hacer el endpoint) y me devuelva el valor del registro
	}
	global.Logger.Log(fmt.Sprintf("Memoria devolvió este valor: %s", *resp), log.INFO)

	// convierto la response en un slice de bytes

	valor := []byte(*resp)
	//revisar este log
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Escribir Archivo: <%s> - Tamaño a Leer: <%d> - Puntero Archivo: <%d>", estructura.Pid, estructura.FileName, len(valor), estructura.PunteroArchivo), log.INFO)

	global.Logger.Log(fmt.Sprintf("Conversión de la respuesta de memoria en un slice de bytes: %v", valor), log.DEBUG)

	// TODO: chequear que donde escribo pertenece al archivo

	global.UpdateBlocksFile(valor, estructura.FileName, estructura.PunteroArchivo)

	global.Logger.Log("Datos escritos exitosamente en el archivo bloques.dat", log.INFO)
	global.PrintBloques()

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func Fs_read(w http.ResponseWriter, r *http.Request) {
	dispositivo := global.Dispositivo
	dispositivo.InUse = true

	var estructura global.KernelIOFS_WR

	err := serialization.DecodeHTTPBody[*global.KernelIOFS_WR](r, &estructura)
	if err != nil {
		global.Logger.Log("Error al decodear: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear", http.StatusBadRequest)
		return
	}
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Operacion: <%s", estructura.Pid, estructura.Instruction+">"), log.INFO)
	global.Logger.Log(fmt.Sprintf("PID: <%d> - Leer Archivo: <%s> - Tamaño a Leer: <%d> - Puntero Archivo: <%d>", estructura.Pid, estructura.FileName, estructura.Tamanio, estructura.PunteroArchivo), log.INFO)

	global.Logger.Log(fmt.Sprintf("%+v", dispositivo), log.DEBUG)

	time.Sleep(time.Duration(global.IOConfig.UnitWorkTime) * time.Millisecond)

	// implementación

	filestruct := global.FilesMap[estructura.FileName]

	ubicacionDeseada := global.IOConfig.DialFSBlockSize*filestruct.Initial_block + estructura.PunteroArchivo // probar si esta ubicación es la correcta

	data := make([]byte, estructura.Tamanio)

	for i := 0; i < estructura.Tamanio; i++ {
		data[i] = global.Bloques[ubicacionDeseada+i]
	}

	global.Logger.Log(fmt.Sprintf("Del slice Bloques leí: %+v ", data), log.DEBUG)

	// armo la estructura para mandar a memoria

	global.Estructura_actualizada.Pid = estructura.Pid
	global.Estructura_actualizada.Content = string(data)
	global.Estructura_actualizada.NumFrames = estructura.NumFrames
	global.Estructura_actualizada.Offset = estructura.Offset
	global.Estructura_actualizada.Length = len(data)

	global.Logger.Log(fmt.Sprintf("String a mandar a memoria: %+v", string(data)), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("Estructura a mandar a memoria: %+v", global.Estructura_actualizada), log.DEBUG)

	// Put a memoria de la estructura
	_, err = requests.PutHTTPwithBody[global.MemStdIO, interface{}](global.IOConfig.IPMemory, global.IOConfig.PortMemory, "stdin_read", global.Estructura_actualizada)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("NO se pudo enviar a memoria la estructura %s", err.Error()), log.ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
		// TODO: falta que memoria vea si puede escribir o no (?)
	}

	dispositivo.InUse = false
	w.WriteHeader(http.StatusNoContent)
}

func compactar(file string, totalfreeblocks int) {

	//sacar el truncado
	global.Logger.Log("Compactando...", log.INFO)
	time.Sleep(time.Duration(global.IOConfig.CompactionDelay) * time.Millisecond)

	filestruct := global.FilesMap[file]
	global.UpdateBitmap(0, filestruct.Initial_block, filestruct.CurrentBlocks)
	global.PrintBitmap()
	
	//actualizar bitmap (mover todos los 1 a la izquierda)
	totalUsedBlocks := global.IOConfig.DialFSBlockCount - totalfreeblocks
	global.UpdateBitmap(1, 0, totalUsedBlocks)
	global.PrintBitmap()
	global.UpdateBitmap(0, totalUsedBlocks, totalfreeblocks)
	global.PrintBitmap()

	//actualizar los initial block de los archivos de metadata
	updateMetadataFiles(file)
	newCurrentBlocksTruncatedFile := global.GetNeededBlocks(global.Estructura_truncate)
	global.UpdateSize(global.Estructura_truncate.FileName, global.Estructura_truncate.Tamanio, newCurrentBlocksTruncatedFile)
	global.RebuildFilesMap(global.IOConfig)

	filestruct = global.FilesMap[file]
	global.UpdateBitmap(1, filestruct.Initial_block, filestruct.CurrentBlocks)
	global.PrintBitmap()

}

func updateMetadataFiles(filename string) {

	var fileNames []string

	dirPath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name

	// Leer los contenidos del directorio
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		global.Logger.Log(fmt.Sprintf("No se pudo leer el directorio que contiene los metadata %s", err.Error()), log.ERROR)
	}
	// Iterar sobre los archivos y agregar sus nombres al slice
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "txt") {
			fileNames = append(fileNames, entry.Name())
		}
	}

	// función para mover el file a la última posición de fileNames

	fileNames = moveToLastPosition(fileNames, filename)

	// Imprimir los nombres de los archivos
	global.Logger.Log(fmt.Sprintf("Nombre de los archivos (el truncado último): %+v", fileNames), log.DEBUG)

	currentInitialBlock := 0

	oldBloques := make([]byte, len(global.Bloques))
	copy(oldBloques, global.Bloques)

	for i := 0; i < len(fileNames); i++ {

		filestruct := global.FilesMap[fileNames[i]]

		oldPosition := filestruct.Initial_block * global.IOConfig.DialFSBlockSize
		oldCantidadDeBytes := filestruct.CurrentBlocks * global.IOConfig.DialFSBlockSize

		slicePortion := oldBloques[oldPosition : oldPosition+oldCantidadDeBytes]

		global.UpdateInitialBlock(fileNames[i], currentInitialBlock)

		global.UpdateBlocksFile(slicePortion, fileNames[i], 0)

		currentInitialBlock = currentInitialBlock + global.GetCurrentBlocks(fileNames[i])
	}

	printMetadataFiles(fileNames)

}

func getFirstFreeBlock() int {

	var firstFreeBlock int

	found := false
	for i, v := range global.Bitmap {
		if v == byte(0) {
			global.Logger.Log(fmt.Sprintf("FirstFreeBlock: %d", i), log.DEBUG)
			firstFreeBlock = i
			found = true
			break
		}
	}
	if !found {
		global.Logger.Log("No hay bloques libres", log.DEBUG)
	}
	global.PrintBitmap()
	return firstFreeBlock
}

func moveToLastPosition(list []string, target string) []string {
	var result []string
	var found bool

	for _, item := range list {
		if item == target {
			found = true
		} else {
			result = append(result, item)
		}
	}

	if found {
		result = append(result, target)
	}

	return result
}

func printMetadataFiles(fileNames []string) {

	var filestruct global.File

	for i := 0; i < len(fileNames); i++ {
		filepath := global.IOConfig.DialFSPath + "/" + global.Dispositivo.Name + "/" + fileNames[i]

		file, err := os.Open(filepath)
		if err != nil {
			global.Logger.Log(fmt.Sprintf("Error al abrir el archivo %s: %s ", filepath, err.Error()), log.DEBUG)
			return
		}

		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&filestruct)
		if err != nil {
			global.Logger.Log(fmt.Sprintf("Error al decodear el archivo %s: %s ", filepath, err.Error()), log.ERROR)

			return
		}

		global.Logger.Log(fmt.Sprintf("Filestruct de %s: %+v ", fileNames[i], filestruct), log.DEBUG)

	}

}

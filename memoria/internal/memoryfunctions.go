package internal

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/sisoputnfrba/tp-golang/memoria/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

var NumPages int



func InstructionStorage(data []string, pid int) {
	
	pagetable := global.NewPageTable()
	//le asigno al map la lista de instrucciones y la tabla de paginas del proceso q pase por id
	global.DictProcess[pid] = global.ListInstructions{Instructions: data, PageTable: pagetable}

	global.Logger.Log(fmt.Sprintf("Contenido pagetable %+v", pagetable), log.DEBUG)
}



func ReadTxt(Path string) ([]string, error) {
	Data, err := os.ReadFile(Path)
	if err != nil {
		global.Logger.Log("Error al leer el archivo "+err.Error(), log.ERROR)
		return nil, err
	}
	ListInstructions := strings.Split(string(Data), "\n")

	return ListInstructions, nil
}

// se le envia un contenido y una direccion para escribir en memoria
func MemOut(NumFrames []int, Offset int, content int, Pid int, Largo int) bool {
	var Slicebytes []byte
	accu := 0
	global.Logger.Log(fmt.Sprintf("largo %d", Largo), log.DEBUG)

	if Offset >= global.MemoryConfig.PageSize {
		global.Logger.Log("Memoria inaccesible", log.ERROR)
		return false
	}
	global.Logger.Log(fmt.Sprintf("PID: %d - Accion: ESCRIBIR - Direccion fisica: %+v + %d - Tamaño: %d Bytes A ESCRIBIR", Pid, NumFrames,Offset, Largo,), log.INFO)

	if Largo == 4 {
		Slicebytes = EncodeContent(uint32(content))

		global.Logger.Log(fmt.Sprintf("largo %+v", Slicebytes), log.DEBUG)
		
		for i := 0; i < Largo; i++ {
			if i+Offset < global.MemoryConfig.PageSize {
				MemFrame := NumFrames[0]*global.MemoryConfig.PageSize + Offset + i
				global.Memory.Spaces[MemFrame] = Slicebytes[i]

			} else {
				
				MemFrame := NumFrames[1]*global.MemoryConfig.PageSize + accu
				global.Memory.Spaces[MemFrame] = Slicebytes[i]
				accu++
			}
		}
	} else if Largo == 1 {
		global.Memory.Spaces[NumFrames[0]*global.MemoryConfig.PageSize+Offset] = byte(content)
		//global.Logger.Log(fmt.Sprintf("PID: %d - Accion: ESCRIBIR - Direccion fisica: %+v + %d - Tamaño: %d Bytes A ESCRIBIR", Pid, NumFrames,Offset, Largo,), log.INFO)

	}


	return true

}

// le paso un valor y me devuelve un slice de bytes en hexa
func EncodeContent(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return bytes
}

func DecodeContent(slice []byte) uint32 {
	return binary.LittleEndian.Uint32(slice)
}

func MemIn(NumFrame []int, Offset int, Pid int, Largo int) int {
	var Content []byte
	var ContentByte byte
	global.Logger.Log(fmt.Sprintf("PID: %d - Accion: LEER - Direccion fisica: %+v + %d - Tamaño: %d Bytes A LEER", Pid, NumFrame,Offset, Largo,), log.INFO)

	if Largo == 4 {
		accu := 0
		for i := 0; i < 4; i++ {
			if Offset+i < global.MemoryConfig.PageSize {
				MemFrame := NumFrame[0]*global.MemoryConfig.PageSize + Offset + i
				ContentByte = global.Memory.Spaces[MemFrame]
				Content = append(Content, ContentByte)
			} else {
				
				MemFrame := NumFrame[1]*global.MemoryConfig.PageSize + accu
				ContentByte = global.Memory.Spaces[MemFrame]
				Content = append(Content, ContentByte)
				accu++
			}
		}

		return int(DecodeContent(Content))
	} else {
		MemFrame := NumFrame[0]*global.MemoryConfig.PageSize + Offset
		ContentByte = global.Memory.Spaces[MemFrame]
		return int(ContentByte)
	}

}

func PageCheck(PageNumber int, Pid int, Offset int) bool {

	global.Logger.Log("La pagina esta bien", log.DEBUG)
	global.Logger.Log(fmt.Sprintf(" %+v", global.DictProcess[Pid]), log.DEBUG)
	

	if checkCompletedPage(PageNumber-1, Pid) {
		global.Logger.Log("estoy dentro de la addpage del else", log.DEBUG)
		AddPage(Pid)
		return true
	}
	return false
}

func checkCompletedPage(PageNumber int, Pid int) bool {
	for i := 0; i < 16; i++ {
		if global.Memory.Spaces[global.DictProcess[Pid].PageTable.Pages[PageNumber]+i] == 0 {
			return false
		}
	}
	return true
}

func GetFrame(PageNumber int, Pid int) int {

	if CheckIfValid(PageNumber, Pid) {
		return global.DictProcess[Pid].PageTable.Pages[PageNumber]
	}
	
	return -1
}

func CheckIfValid(PageNumber int, Pid int) bool {
	if process, ok := global.DictProcess[Pid]; ok && process.PageTable != nil {
		if len(process.PageTable.Pages) > 0 {
			for pageNum := range process.PageTable.Pages {
				if PageNumber == pageNum {
					return true
				}
			}
		}
	}
	return false
}

func AddPage(Pid int) int {
	for i := 0; i < len(global.BitMap); i++ {
		if global.BitMap[i] == 0 {
			global.DictProcess[Pid].PageTable.Pages = append(global.DictProcess[Pid].PageTable.Pages, i)
			global.BitMap[i] = 1
			return i
		}
	}
	return -1
}

func WriteInMemory (byteArray []byte,Length int,NumFrames []int,Offset int){
	
	accu:=0
	
	j:=1
		for i := 0; i < Length; i++ {
			
			if i+Offset < global.MemoryConfig.PageSize {
				MemFrame := NumFrames[0]*global.MemoryConfig.PageSize + Offset + i
				global.Memory.Spaces[MemFrame] = byteArray[i]
				
			} else {
				if accu>=global.MemoryConfig.PageSize {
					accu=0
					j++
				}
				
				MemFrame := NumFrames[j]*global.MemoryConfig.PageSize + accu
				global.Memory.Spaces[MemFrame] = byteArray[i]
				accu++
			}
		}
}

func ReadInMemory (Length int,NumFrames []int,Offset int) []byte{
	var Content []byte
	var ContentByte byte

		j:=1
		accu := 0
		for i := 0; i < Length; i++ {
			if Offset+i < global.MemoryConfig.PageSize {
				MemFrame := NumFrames[0]*global.MemoryConfig.PageSize + Offset + i
				ContentByte = global.Memory.Spaces[MemFrame]
				Content = append(Content, ContentByte)
			} else {
				if accu>=global.MemoryConfig.PageSize {
					accu=0
					j++
				}
				
				MemFrame := NumFrames[j]*global.MemoryConfig.PageSize + accu
				ContentByte = global.Memory.Spaces[MemFrame]
				Content = append(Content, ContentByte)
				accu++
			}
		}
		
		return Content
}


func PrintMemoryTable(memory []byte, cols int) {
	// Imprimir encabezados de columna
	fmt.Print("Addr\t")
	for i := 0; i < cols; i++ {
		fmt.Printf("%02X ", i)
	}
	fmt.Println()

	// Imprimir separador de encabezado
	fmt.Print("----\t")
	for i := 0; i < cols; i++ {
		fmt.Print("---")
	}
	fmt.Println()

	// Imprimir contenido de la memoria
	for i := 0; i < len(memory); i += cols {
		// Imprimir dirección (índice de fila)
		fmt.Printf("%04X\t", i)

		// Imprimir los bytes en la fila
		for j := 0; j < cols; j++ {
			if i+j < len(memory) {
				fmt.Printf("%02X ", memory[i+j])
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println()
	}
}

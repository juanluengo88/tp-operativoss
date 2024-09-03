package global

import (
	"fmt"
	"os"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/internal/tlb"
	config "github.com/sisoputnfrba/tp-golang/utils/config"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
)

const CPULOG = "./cpu.log"

type Config struct {
	Port             int    `json:"port"`
	IPKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	IPMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	NumberFellingTLB int    `json:"number_felling_tlb"`
	AlgorithmTLB     string `json:"algorithm_tlb"`
	Page_size        int    `json:"page_size"`
}

var CPUConfig *Config
var Execute bool
var InterruptReason string
var Logger *log.LoggerStruct
var Tlb *tlb.TLB

// mutex
var ExecuteMutex sync.Mutex

func InitGlobal() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("ARGS: ENV=dev|prod CONFIG=config_path")
		os.Exit(1)
	}
	env := args[0]
	configFile := args[1]

	Logger = log.ConfigureLogger(CPULOG, env)
	CPUConfig = config.LoadConfiguration[Config](configFile)
	Tlb = tlb.NewTLB(CPUConfig.NumberFellingTLB, CPUConfig.AlgorithmTLB)
}

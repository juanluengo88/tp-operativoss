package model

import (
	"container/list"
	"sync"
)

type PCB struct {
	PID              int
	State            string
	PC               int
	CPUTime          int
	Quantum          int
	RemainingQuantum int
	DisplaceReason   string
	Registers        CPURegister
	Instruction      Instruction
}

type CPURegister struct {
	AX  int
	BX  int
	CX  int
	DX  int
	EAX int
	EBX int
	ECX int
	EDX int
	SI  int
	DI  int
}

type Instruction struct {
	Operation  string
	Parameters []string
	NumFrames  []int
	Offset     int
	Size       int
	FSPointer  int
}

type ProcessInstruction struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}

type IOSTD struct {
	Pid       int    `json:"pid"`
	Name      string `json:"name"`
	Length    int    `json:"length"`
	NumFrames []int  `json:"numframe"`
	Offset    int    `json:"offset"`
}

type IoDevice struct {
	Port int    `json:"port"`
	Name string `json:"name"`
	Type string `json:"type"`
	Sem  chan int
}

type Resource struct {
	Name        string
	Count       int
	BlockedList *list.List
	MutexList   sync.Mutex
	PidList     []int
}

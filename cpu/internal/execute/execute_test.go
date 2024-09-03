package execute

import (
	"testing"

	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func TestSet(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{},
	}
	instruction := &model.Instruction{
		Operation:  "SET",
		Parameters: []string{"AX", "100"},
	}

	set(pcb, instruction)

	expected := 100
	if pcb.Registers.AX != expected {
		t.Errorf("Set failed: expected AX = %d, got %d", expected, pcb.Registers.AX)
	}
}

func TestSum(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{
			AX: 10,
			BX: 5,
		},
	}
	instruction := &model.Instruction{
		Operation:  "SUM",
		Parameters: []string{"AX", "BX"},
	}

	sum(pcb, instruction)

	expected := 15
	if pcb.Registers.AX != expected {
		t.Errorf("Sum failed: expected AX = %d, got %d", expected, pcb.Registers.AX)
	}
}

func TestSub(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{
			AX: 10,
			BX: 5,
		},
	}
	instruction := &model.Instruction{
		Operation:  "SUB",
		Parameters: []string{"AX", "BX"},
	}

	sub(pcb, instruction)

	expected := 5
	if pcb.Registers.AX != expected {
		t.Errorf("Sub failed: expected AX = %d, got %d", expected, pcb.Registers.AX)
	}
}

func TestJnz(t *testing.T) {
	pcb := &model.PCB{
		Registers: model.CPURegister{
			AX: 5, // AX es distinto de cero inicialmente
		},
		PC: 0,
	}
	instruction := &model.Instruction{
		Operation:  "JNZ",
		Parameters: []string{"AX", "10"},
	}

	t.Run("JNZ sin 0", func(t *testing.T) {
		jnz(pcb, instruction)

		expected := 10
		if pcb.PC != expected {
			t.Errorf("Jnz failed: expected PC = %d, got %d", expected, pcb.PC)
		}
	})

	// Reestablecer PCB para el siguiente test
	pcb.Registers.AX = 0
	pcb.PC = 0

	// Test con AX cero
	t.Run("JNZ con 0", func(t *testing.T) {
		jnz(pcb, instruction)
		expected := 0
		if pcb.PC != expected {
			t.Errorf("Jnz failed on zero condition: expected PC to remain %d, got %d", expected, pcb.PC)
		}
	})
}

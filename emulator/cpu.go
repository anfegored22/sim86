package emulator

import (
	"fmt"
)

// Trap, Direction, Interrupt, Overflow, Sign,
// Zero, Auxiliary Carry, Parity, Carry
type Flags struct {
	TF, DF, IF, OF, SF, ZF, AF, PF, CF bool
}

type CPU struct {
	Memory
	Registers
	Flags
}

func NewCPU() CPU {
	var regs = map[byte]uint16{
		0b000: 0, 0b001: 0, 0b010: 0, 0b011: 0,
		0b100: 0, 0b101: 0, 0b110: 0, 0b111: 0,
	}
	var segRegs = map[byte]uint16{
		0b00: 0, 0b01: 0,
		0b10: 0, 0b11: 0,
	}
	r := Registers{R: regs, SR: segRegs,
		RNames:  [8]string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"},
		SRNames: [4]string{"ES", "CS", "SS", "DS"}}
	return CPU{Memory: Memory{}, Registers: r, Flags: Flags{}}
}

func (c *CPU) Add(reg byte, w bool, data uint16) error {
	old := c.Registers.Load(reg, w)
	mask := uint16(0xff)
	signBit := uint16(0x80) // 1000 000
	if w {
		mask = 0xffff
		signBit = 0x8000
	}
	signOld := old&mask&signBit != 0
	signData := data&mask&signBit != 0
	result := old&mask + data&mask
	signResult := result&mask&signBit != 0
	// Now we need to set the flags
	c.Flags.ZF = result&mask == 0
	c.Flags.SF = signResult
	// Carrie Flag
	if w {
		c.Flags.CF = result < old
	} else {
		c.Flags.CF = result&0x100 != 0
	}
	c.Flags.OF = signData == signOld && signResult != signOld
	return c.Registers.Store(reg, w, result&mask)
}

func (c *CPU) Execute(inst Instruction) error {
	return fmt.Errorf("We don't have that opcode: %s", inst[Opcode])
}

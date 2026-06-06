package emulator

import (
	"fmt"
)

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
	result := old&mask + data&mask
	c.Flags.SF = (result&mask)&signBit == 1
	if w {
		c.Flags.OF = result < old
	} else {
		c.Flags.OF = result&0x100 == 1
	}
	return nil
}

func (c *CPU) Execute(inst Instruction) error {
	return fmt.Errorf("We don't have that opcode: %s", inst[Opcode])
}

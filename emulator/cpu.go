package emulator

import (
	"fmt"
	"strings"
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

func (c *CPU) Execute(inst Instruction) error {
	switch inst[Operation] {
	case Mov:
		return ExecuteMov(inst, c)
	case Add:
		return ExecuteAdd(inst, c)
	case Sub:
		return ExecuteSub(inst, c)
	case Cmp:
		return ExecuteCmp(inst, c)
	}
	return fmt.Errorf("Don't understand operation %b", inst[Opcode])
}

// It stores to the register but can be an address!
// We should divided into set flags
func (c *CPU) Add(old uint16, w bool, data uint16) uint16 {
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
	return result
}

func (c *CPU) Sub(old uint16, w bool, data uint16) uint16 {
	mask := uint16(0xff)
	signBit := uint16(0x80) // 1000 000
	if w {
		mask = 0xffff
		signBit = 0x8000
	}
	signOld := old&mask&signBit != 0
	signData := data&mask&signBit != 0
	result := old&mask - data&mask
	signResult := result&mask&signBit != 0
	// Now we need to set the flags
	c.Flags.ZF = result&mask == 0
	c.Flags.SF = signResult
	// Carrie Flag
	c.Flags.CF = old&mask < data&mask // Borrow flag
	c.Flags.OF = signOld != signData && signResult != signOld
	return result
}

func (c *CPU) Cmp(old uint16, w bool, data uint16) {
	mask := uint16(0xff)
	signBit := uint16(0x80) // 1000 000
	if w {
		mask = 0xffff
		signBit = 0x8000
	}
	signOld := old&mask&signBit != 0
	signData := data&mask&signBit != 0
	result := old&mask - data&mask
	signResult := result&mask&signBit != 0
	// Now we need to set the flags
	c.Flags.ZF = result&mask == 0
	c.Flags.SF = signResult
	// Carrie Flag
	c.Flags.CF = old&mask < data&mask // Borrow flag
	c.Flags.OF = signOld != signData && signResult != signOld
}

func (c *CPU) ReadFlags() {
	switch {
	case c.Flags.SF != c.Flags.CF:
		print("a < b")
	case c.Flags.ZF:
		print("a = b")
	default:
		print("a > b")
	}
}

func (c *CPU) Print() {
	for i, name := range c.Registers.RNames {
		value := uint16(c.Registers.R[byte(i)])
		high := byte(value >> 8)
		low := byte(value)
		signed := int16(value)

		fmt.Print("\033[2K") // clear current line before overwriting shorter old values
		fmt.Printf("%s: %08b | %08b (%d / %d)\n", name, high, low, value, signed)
	}
	for i, name := range c.Registers.SRNames {
		value := uint16(c.Registers.SR[byte(i)])
		fmt.Print("\033[2K") // clear current line before overwriting shorter old values
		fmt.Printf("%s: %16b (%d)\n", name, value, c.Registers.SR[byte(i)])
	}

	flagNames := []string{"TF", "DF", "IF", "OF", "SF", "ZF", "AF", "PF", "CF"}
	flagValues := []bool{c.Flags.TF, c.Flags.DF, c.Flags.IF, c.Flags.OF, c.Flags.SF, c.Flags.ZF, c.Flags.AF, c.Flags.PF, c.Flags.CF}

	fmt.Print("\033[2K")
	fmt.Printf("%s  ", strings.Join(flagNames, " "))
	fmt.Println()

	fmt.Print("\033[2K")
	for _, value := range flagValues {
		if value {
			fmt.Print("1  ")
		} else {
			fmt.Print("0  ")
		}
	}
	fmt.Println()
}

func (c *CPU) Update(inst Instruction) {
	fmt.Printf("\033[%dA", len(c.Registers.RNames)+len(c.Registers.SRNames)+4)
	c.Print()
	fmt.Printf("Instruction: opcode=%06b d=%d w=%d mod=%02b reg=%03b rm=%03b\n",
		inst[Opcode], inst[D], inst[W], inst[Mod], inst[Reg], inst[RM])
}

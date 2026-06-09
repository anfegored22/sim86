package emulator

import "fmt"

func ExecuteJump(inst Instruction, c *CPU) error {
	switch inst[Opcode] {
	case 0x74:
		if c.Flags.ZF {
			c.Registers.IP += int(int8(byte(inst[IPinc8]))) * 8
		}
		return nil
	case 0x75:
		if !c.Flags.ZF {
			c.Registers.IP += int(int8(byte(inst[IPinc8]))) * 8
		}
		return nil
	}
	return fmt.Errorf("No jump instruction")
}

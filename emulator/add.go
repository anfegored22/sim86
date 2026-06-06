package emulator

import "fmt"

func ExecuteAdd(inst Instruction, c *CPU) error {
	switch {
	case inst[Opcode] == 0b0:
		wide := inst[W] == 1
		reg := byte(inst[Reg])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		if inst[D] == 0 {
			src := c.Registers.Load(reg, wide)
			if mod == 0b11 {
				dst := c.Registers.Load(rm, wide)
				result := src + dst
				// I need to add the Flag!
				return c.Registers.Store(rm, wide, result)
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			return c.Memory.Store(addr, src, wide)
		} else {
			if mod == 0b11 {
				src := c.Registers.Load(rm, wide)
				return c.Registers.Store(reg, wide, src)
			} else {
				addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
				src, err := c.Memory.Load(addr, wide)
				if err != nil {
					return err
				}
				return c.Registers.Store(reg, wide, src)
			}
		}
	case inst[Opcode] == 0b100000:
		wide := inst[W] == 1
		s := inst[S]
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		src := inst[Data]
		if mod == 0b11 {
			return c.Registers.Store(rm, wide, src)
		}
		addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
		return c.Memory.Store(addr, src, wide)
	case inst[Opcode] == 0b0000010:
		wide := inst[W] == 1
		data := inst[Data]
		result := c.Registers.Load(0, wide) + data
		if wide {
			result - 0x100 // 100000000 00000000
		}

	}
	return fmt.Errorf("Opcode %b Not found", inst[Opcode])
}

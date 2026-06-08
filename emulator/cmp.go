package emulator

import "fmt"

// TODO: We need an Add that saves to memory instead of the register!

func ExecuteCmp(inst Instruction, c *CPU) error {
	switch {
	case inst[Opcode] == 0b001110:
		wide := inst[W] == 1
		reg := byte(inst[Reg])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		if inst[D] == 0 {
			src := c.Registers.Load(reg, wide)
			if mod == 0b11 {
				old := c.Registers.Load(rm, wide)
				c.Sub(old, wide, src)
				return nil
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			old, err := c.Memory.Load(addr, wide) // Is the address a signed int? Why?
			if err != nil {
				return err
			}
			c.Sub(old, wide, src)
			return nil
		} else {
			if mod == 0b11 {
				src := c.Registers.Load(rm, wide)
				old := c.Registers.Load(reg, wide)
				c.Sub(old, wide, src)
				return nil
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			src, err := c.Memory.Load(addr, wide)
			if err != nil {
				return err
			}
			old := c.Registers.Load(reg, wide)
			c.Sub(old, wide, src)
			return nil
		}
	case inst[Opcode] == 0b100000 || inst[Opcode] == 0b0011110:
		wide := inst[W] == 1
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		src := inst[Data]
		if mod == 0b11 {
			old := c.Registers.Load(rm, wide)
			c.Sub(old, wide, src) // The Add set flags
			return nil
		}
		addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
		old, err := c.Memory.Load(addr, wide)
		if err != nil {
			return err
		}
		c.Sub(old, wide, src)
		return nil

	}
	return fmt.Errorf("Opcode %b Not found", inst[Opcode])
}

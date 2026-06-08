package emulator

import "fmt"

// TODO: We need an Add that saves to memory instead of the register!

func ExecuteSub(inst Instruction, c *CPU) error {
	switch {
	case inst[Opcode] == 0b001010:
		wide := inst[W] == 1
		reg := byte(inst[Reg])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		if inst[D] == 0 {
			src := c.Registers.Load(reg, wide)
			if mod == 0b11 {
				old := c.Registers.Load(rm, wide)
				result := c.Sub(old, wide, src)
				return c.Registers.Store(rm, wide, result)
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			old, err := c.Memory.Load(addr, wide) // Is the address a signed int? Why?
			if err != nil {
				return err
			}
			result := c.Sub(old, wide, src)
			return c.Memory.Store(addr, result, wide)
		} else {
			if mod == 0b11 {
				src := c.Registers.Load(rm, wide)
				old := c.Registers.Load(reg, wide)
				result := c.Sub(old, wide, src)
				return c.Registers.Store(reg, wide, result)
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			src, err := c.Memory.Load(addr, wide)
			if err != nil {
				return err
			}
			old := c.Registers.Load(reg, wide)
			result := c.Sub(old, wide, src)
			return c.Registers.Store(reg, wide, result)
		}
	case inst[Opcode] == 0b100000 || inst[Opcode] == 0b0010110:
		wide := inst[W] == 1
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		src := inst[Data]
		if mod == 0b11 {
			old := c.Registers.Load(rm, wide)
			result := c.Sub(old, wide, src) // The Add set flags
			return c.Registers.Store(rm, wide, result)
		}
		addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
		old, err := c.Memory.Load(addr, wide)
		if err != nil {
			return err
		}
		result := c.Sub(old, wide, src)
		return c.Memory.Store(addr, result, wide)

	}
	return fmt.Errorf("Opcode %b Not found", inst[Opcode])
}

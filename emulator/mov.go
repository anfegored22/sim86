package emulator

import "fmt"

func ExecuteMov(inst Instruction, c *CPU) error {
	switch {
	case inst[Opcode] == 0b100010 || (inst[Opcode]>>1) == 0b101000:
		wide := inst[W] == 1
		reg := byte(inst[Reg])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		if inst[D] == 0 {
			src := c.Registers.Load(reg, wide)
			if mod == 0b11 {
				return c.Registers.Store(rm, wide, src)
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
	case inst[Opcode] == 0b1100011:
		wide := inst[W] == 1
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		src := inst[Data]
		if mod == 0b11 {
			return c.Registers.Store(rm, wide, src)
		}
		addr := c.Registers.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
		return c.Memory.Store(addr, src, wide)
	case inst[Opcode] == 0b1011:
		return c.Registers.Store(byte(inst[Reg]), inst[W] == 1, inst[Data])
	case inst[Opcode]>>2 == 0b100011:
		sr := byte(inst[SR])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		disp := int16(inst[Disp])
		if inst[D] == 0 {
			src := c.Registers.SR[sr]
			if mod == 0b11 {
				return c.Registers.Store(rm, true, src)
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, disp, mod)
			return c.Memory.Store(addr, src, true)
		} else {
			if mod == 0b11 {
				src := c.Registers.Load(rm, true)
				c.Registers.SR[sr] = src
				return nil
			}
			addr := c.Registers.DecodeEffectiveAddr(rm, disp, mod)
			src, err := c.Memory.Load(addr, true)
			if err != nil {
				return fmt.Errorf("error loading address %b: %w", addr, err)
			}
			if _, ok := c.Registers.SR[sr]; !ok {
				return fmt.Errorf("No segment register found: %0b", sr)
			}
			c.Registers.SR[sr] = src
			return nil
		}
	}
	return fmt.Errorf("Opcode %b Not found", inst[Opcode])
}

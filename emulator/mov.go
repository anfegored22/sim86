package emulator

import "fmt"

func ExecuteMov(inst Instruction, regs *Registers, memory *Memory) error {
	switch {
	case inst[Opcode] == 0b100010 || (inst[Opcode]>>1) == 0b101000:
		wide := inst[W] == 1
		reg := byte(inst[Reg])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		if inst[D] == 0 {
			src := regs.Load(reg, wide)
			if mod == 0b11 {
				return regs.Store(rm, wide, src)
			}
			addr := regs.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			return memory.Store(addr, src, wide)
		} else {
			if mod == 0b11 {
				src := regs.Load(rm, wide)
				return regs.Store(reg, wide, src)
			} else {
				addr := regs.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
				src, err := memory.Load(addr, wide)
				if err != nil {
					return err
				}
				return regs.Store(reg, wide, src)
			}
		}
	case inst[Opcode] == 0b1100011:
		wide := inst[W] == 1
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		src := inst[Data]
		if mod == 0b11 {
			return regs.Store(rm, wide, src)
		}
		addr := regs.DecodeEffectiveAddr(rm, int16(inst[Disp]), mod)
		return memory.Store(addr, src, wide)
	case inst[Opcode] == 0b1011:
		return regs.Store(byte(inst[Reg]), inst[W] == 1, inst[Data])
	case inst[Opcode]>>2 == 0b100011:
		sr := byte(inst[SR])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		disp := int16(inst[Disp])
		if inst[D] == 0 {
			src := regs.SR[sr]
			if mod == 0b11 {
				return regs.Store(rm, true, src)
			}
			addr := regs.DecodeEffectiveAddr(rm, disp, mod)
			return memory.Store(addr, src, true)
		} else {
			if mod == 0b11 {
				src := regs.Load(rm, true)
				regs.SR[sr] = src
				return nil
			}
			addr := regs.DecodeEffectiveAddr(rm, disp, mod)
			src, err := memory.Load(addr, true)
			if err != nil {
				return fmt.Errorf("error loading address %b: %w", addr, err)
			}
			if _, ok := regs.SR[sr]; !ok {
				return fmt.Errorf("No segment register found: %0b", sr)
			}
			regs.SR[sr] = src
			return nil
		}
	}
	return fmt.Errorf("Opcode %b Not found", inst[Opcode])
}

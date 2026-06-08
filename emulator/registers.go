package emulator

import "fmt"

var effectiveAddress = map[byte]byte{
	0b000: 0b011110,
	0b001: 0b011111,
	0b010: 0b101110,
	0b011: 0b101111,
	0b100: 0b110000,
	0b101: 0b111000,
	0b110: 0b101000,
	0b111: 0b011000,
}

type Registers struct {
	R       map[byte]uint16
	SR      map[byte]uint16
	RNames  [8]string
	SRNames [4]string
}

func (r *Registers) Load(reg byte, w bool) uint16 {
	v, _ := r.R[reg] // Load never fails!
	if w {
		return v
	}
	if reg >= 4 {
		return r.R[reg-0b100] >> 8
	} else {
		return r.R[reg] & 0x00ff
	}
}

// The data has been decoded in little-indian. We dont shift anything in anything
func (r *Registers) Store(reg byte, w bool, data uint16) error {
	if _, ok := r.R[reg]; !ok {
		return fmt.Errorf("No register found")
	}
	if w {
		r.R[reg] = data
		return nil
	}
	if reg >= 4 {
		xReg := r.R[reg-0x4]
		r.R[reg-0x4] = xReg&0x00ff | ((data & 0x00ff) << 8)
		return nil
	} else {
		xReg := r.R[reg]
		r.R[reg] = xReg&0xff00 | (data & 0x00ff)
		return nil
	}
}

func (r *Registers) DecodeEffectiveAddr(rm byte, disp int16, mod byte) int16 {
	if rm == 0b110 && mod == 0b00 {
		return disp
	}
	regs, _ := effectiveAddress[rm]
	reg1 := byte(regs >> 3)
	reg2 := byte(regs & 0b111)
	if reg2 != 0 {
		return int16(r.Load(reg1, true)) + int16(r.Load(reg2, true)) + disp
	}
	return int16(r.Load(reg1, true)) + disp
}

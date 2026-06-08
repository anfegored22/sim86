package emulator

type Instruction map[FieldKind]uint16

func Decoder(data []byte) []Instruction {
	bitPos := 0
	var instructions []Instruction
	for bitPos < len(data)*8 {
		for _, inst := range instructionPatterns {
			opcode := ReadBits(data, bitPos, int(inst.Bits))
			if opcode != inst.Value {
				continue
			}
			instValue := make(map[FieldKind]uint16)
			instValue[Opcode] = opcode
			instValue[Operation] = inst.Mnemonic
			pos := bitPos + int(inst.Bits)
			for _, field := range inst.Fields {
				if field.Bits == 0 {
					instValue[field.Kind] = field.Value
					continue
				}
				if field.Kind == DataW {
					if instValue[W] == 1 && instValue[S] == 0 {
						dataHi := ReadBits(data, pos, int(field.Bits))
						instValue[Data] = dataHi<<8 | instValue[Data]
						pos += int(field.Bits)
					}
					if instValue[W] == 1 && instValue[S] == 1 && instValue[Data]&0x0080 != 0 {
						instValue[Data] = instValue[Data] | 0xff00
					}
					continue
				}
				if field.Kind == Disp && field.Bits == 16 {
					instValue[Disp] = ReadU16LE(data, pos)
					pos += 16
					continue
				}
				instValue[field.Kind] = ReadBits(data, pos, int(field.Bits))
				pos += int(field.Bits)
				if field.Kind == RM {
					disp, inc := checkMod(data, instValue[Mod], instValue[RM], pos)
					pos += inc
					instValue[Disp] = disp
				}
			}
			bitPos = pos
			instValue[Operation] = resolveOp(instValue)
			instructions = append(instructions, instValue)
			break
		}
	}
	return instructions
}

func checkMod(data []byte, mod uint16, rm uint16, bitPos int) (uint16, int) {
	switch mod {
	case 0b11:
		return 0, 0
	case 0b00:
		if rm == 0b110 {
			lo := ReadBits(data, bitPos, 8)
			hi := ReadBits(data, bitPos+8, 8)
			disp := hi<<8 | lo
			return disp, 16
		}
		return 0, 0
	case 0b01:
		// For 8 bits displacement we need to check the sign and padd the other 8 bits
		disp := ReadBits(data, bitPos, 8)
		if disp&0x80 != 0 {
			disp = disp | 0xff00
		}
		return disp, 8
	case 0b10:
		lo := ReadBits(data, bitPos, 8)
		hi := ReadBits(data, bitPos+8, 8)
		disp := hi<<8 | lo
		return disp, 16
	}
	return 0, 0
}

func resolveOp(ins map[FieldKind]uint16) Mnemonic {
	switch ins[Opcode] {
	case 0b100000:
		switch ins[Reg] {
		case 0b000:
			return Add
		case 0b101:
			return Sub
		case 0b111:
			return Cmp
		}
	default:
		return ins[Operation]
	}
	return ins[Operation]
}

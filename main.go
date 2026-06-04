package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var memory [1024 * 1024]byte

func loadMemory(addr int16, wide bool) (uint16, error) {
	if addr < 0 || int(addr) >= len(memory) {
		return 0, fmt.Errorf("no memory address found")
	}
	lo := memory[addr]
	if !wide {
		return uint16(lo), nil
	}
	if int(addr+1) >= len(memory) {
		return 0, fmt.Errorf("no memory address found")
	}
	hi := memory[addr+1]
	return uint16(hi)<<8 | uint16(lo), nil
}

func storeMemory(addr int16, value uint16, w bool) error {
	if addr < 0 || int(addr) >= len(memory) {
		return fmt.Errorf("no memory address found")
	}
	if w && int(addr+1) >= len(memory) {
		return fmt.Errorf("Out of memory")
	}
	memory[addr] = byte(value & 0x00ff)
	if w {
		memory[addr+1] = byte(value >> 8)
		return nil
	}
	return nil
}

var registers = map[byte]uint16{
	0b000: 0, //AX
	0b001: 0, //CX
	0b010: 0, //DX
	0b011: 0, //BX
	0b100: 0, //SP
	0b101: 0, //BP
	0b110: 0, //SI
	0b111: 0, //DI
}

func loadFromRegister(reg byte, w bool) uint16 {
	v, _ := registers[reg]
	if w {
		return v
	}
	if reg >= 4 {
		return registers[reg-0b100] >> 8
	} else {
		return registers[reg] & 0x00ff
	}
}

func storeToRegister(reg byte, w bool, data uint16) error {
	if _, ok := registers[reg]; !ok {
		return fmt.Errorf("No register found")
	}
	if w {
		registers[reg] = data
		return nil
	}
	if reg >= 4 {
		xReg := registers[reg-0x4]
		registers[reg-0x4] = xReg&0x00ff | ((data & 0x00ff) << 8)
		return nil
	} else {
		xReg := registers[reg]
		registers[reg] = xReg&0xff00 | (data & 0x00ff)
		return nil
	}
}

var segmentRegisters = map[byte]uint16{
	0b00: 0, // ES
	0b01: 0, // CS
	0b10: 0, // SS
	0b11: 0, // DS
}

var registerNames = []string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"}
var segmentsNames = []string{"ES", "CS", "SS", "DS"}

// Encoded as regreg
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

func decodeEffectiveAddr(rm byte, disp int16, mod byte) int16 {
	if rm == 0b110 && mod == 0b00 {
		return disp
	}
	regs, _ := effectiveAddress[rm]
	reg1 := byte(regs >> 3)
	reg2 := byte(regs & 0b111)
	if reg2 != 0 {
		return int16(registers[reg1]) + int16(registers[reg2]) + disp
	}
	return int16(registers[reg1]) + disp
}

// Should I do something like wregreg
func PrintRegisters(registers map[byte]uint16) {
	for i, name := range registerNames {
		value := uint16(registers[byte(i)])
		high := byte(value >> 8)
		low := byte(value)

		fmt.Print("\033[2K") // clear current line before overwriting shorter old values
		fmt.Printf("%s: %08b | %08b (%d)\n", name, high, low, registers[byte(i)])
	}
	for i, name := range segmentsNames {
		value := uint16(segmentRegisters[byte(i)])
		fmt.Print("\033[2K") // clear current line before overwriting shorter old values
		fmt.Printf("%s: %16b (%d)\n", name, value, segmentRegisters[byte(i)])
	}
}

func UpdateRegisters(registers map[byte]uint16) {
	fmt.Printf("\033[%dA", len(registerNames)+1)
	PrintRegisters(registers)
}

func readBits(data []byte, bitPos int, n int) uint16 {
	var value uint16
	for range n {
		byteIndex := bitPos / 8
		bitIndex := 7 - (bitPos % 8)
		bit := (data[byteIndex] >> bitIndex) & 1
		value = (value << 1) | uint16(bit)
		bitPos++
	}
	return value
}

func readU16LE(data []byte, bitPos int) uint16 {
	lo := readBits(data, bitPos, 8)
	hi := readBits(data, bitPos+8, 8)
	return hi<<8 | lo
}

type Instruction map[FieldKind]uint16

func checkMod(data []byte, mod uint16, rm uint16, bitPos int) (uint16, int) {
	switch mod {
	case 0b11:
		return 0, 0
	case 0b00:
		if rm == 0b110 {
			lo := readBits(data, bitPos, 8)
			hi := readBits(data, bitPos+8, 8)
			disp := hi<<8 | lo
			return disp, 16
		}
		return 0, 0
	case 0b01:
		return readBits(data, bitPos, 8), 8
	case 0b10:
		lo := readBits(data, bitPos, 8)
		hi := readBits(data, bitPos+8, 8)
		disp := hi<<8 | lo
		return disp, 16
	}
	return 0, 0
}

func DecodeInstructions(data []byte) []Instruction {
	bitPos := 0
	var instructions []Instruction
	for bitPos < len(data)*8 {
		for _, inst := range instructionPatterns {
			opcode := readBits(data, bitPos, int(inst.Bits))
			if opcode != inst.Value {
				continue
			}
			instValue := make(map[FieldKind]uint16)
			instValue[Opcode] = opcode
			pos := bitPos + int(inst.Bits)
			for _, field := range inst.Fields {
				if field.Bits == 0 {
					instValue[field.Kind] = field.Value
					continue
				}
				if field.Kind == DataW {
					if instValue[W] == 1 {
						dataHi := readBits(data, pos, int(field.Bits))
						instValue[Data] = dataHi<<8 | instValue[Data]
						pos += int(field.Bits)
					}
					continue
				}
				if field.Kind == Disp && field.Bits == 16 {
					instValue[Disp] = readU16LE(data, pos)
					pos += 16
					continue
				}
				instValue[field.Kind] = readBits(data, pos, int(field.Bits))
				pos += int(field.Bits)
				if field.Kind == RM {
					disp, inc := checkMod(data, instValue[Mod], instValue[RM], pos)
					pos += inc
					instValue[Disp] = disp
				}
			}
			bitPos = pos
			instructions = append(instructions, instValue)
			break
		}
	}
	return instructions
}

func ExecuteMov(inst Instruction) error {
	switch {
	case inst[Opcode] == 0b100010 || (inst[Opcode]>>1) == 0b101000:
		wide := inst[W] == 1
		reg := byte(inst[Reg])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		if inst[D] == 0 {
			src := loadFromRegister(reg, wide)
			if mod == 0b11 {
				return storeToRegister(rm, wide, src)
			}
			addr := decodeEffectiveAddr(rm, int16(inst[Disp]), mod)
			return storeMemory(addr, src, wide)
		} else {
			if mod == 0b11 {
				src := loadFromRegister(rm, wide)
				return storeToRegister(reg, wide, src)
			} else {
				addr := decodeEffectiveAddr(rm, int16(inst[Disp]), mod)
				src, err := loadMemory(addr, wide)
				if err != nil {
					return err
				}
				return storeToRegister(reg, wide, src)
			}
		}
	case inst[Opcode] == 0b1100011:
		wide := inst[W] == 1
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		src := inst[Data]
		if mod == 0b11 {
			return storeToRegister(rm, wide, src)
		}
		addr := decodeEffectiveAddr(rm, int16(inst[Disp]), mod)
		return storeMemory(addr, src, wide)
	case inst[Opcode] == 0b1011:
		return storeToRegister(byte(inst[Reg]), inst[W] == 1, inst[Data])
	case inst[Opcode]>>2 == 0b100011:
		sr := byte(inst[SR])
		rm := byte(inst[RM])
		mod := byte(inst[Mod])
		disp := int16(inst[Disp])
		if inst[D] == 0 {
			src := segmentRegisters[sr]
			if mod == 0b11 {
				return storeToRegister(rm, true, src)
			}
			addr := decodeEffectiveAddr(rm, disp, mod)
			return storeMemory(addr, src, true)
		} else {
			if mod == 0b11 {
				src := loadFromRegister(rm, true)
				segmentRegisters[sr] = src
				return nil
			}
			addr := decodeEffectiveAddr(rm, disp, mod)
			src, err := loadMemory(addr, true)
			if err != nil {
				return fmt.Errorf("error loading address %b: %w", addr, err)
			}
			if _, ok := segmentRegisters[sr]; !ok {
				return fmt.Errorf("No segment register found: %0b", sr)
			}
			segmentRegisters[sr] = src
			return nil
		}
	}
	return fmt.Errorf("Opcode %b Not found", inst[Opcode])
}

func main() {
	fn := flag.String("file", "", "Name of the assembly file")
	flag.Parse()
	if *fn == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: -file")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.ReadFile(*fn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	PrintRegisters(registers)
	instructions := DecodeInstructions(f)
	reader := bufio.NewReader(os.Stdin)
	for _, inst := range instructions {
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ExecuteMov(inst)
		UpdateRegisters(registers)
	}
}

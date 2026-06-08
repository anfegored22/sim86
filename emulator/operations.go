package emulator

type FieldKind int

const (
	D FieldKind = iota
	S
	W
	V
	Z
	Mod
	Reg
	Opcode
	Operation
	RM
	SR
	Fix
	Data
	DataW
	Disp
	AddrLo
	AddrHi
)

type Mnemonic = uint16

const (
	Mov Mnemonic = iota
	Add
	Sub
	Cmp
)

type Field struct {
	Kind  FieldKind
	Bits  uint8
	Value uint16
}

type InstructionPattern struct {
	Mnemonic
	Desc   string
	Bits   uint8
	Value  uint16
	Fields []Field
}

var instructionPatterns = []InstructionPattern{
	// MOV
	{Mnemonic: Mov, Desc: "Register/memory to/from register",
		Bits: 6, Value: 0b100010,
		Fields: []Field{{Kind: D, Bits: 1}, {Kind: W, Bits: 1}, {Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3}, {Kind: RM, Bits: 3}},
	},
	{Mnemonic: Mov, Desc: "Immediate to register/memory",
		Bits: 7, Value: 0b1100011,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3}, {Kind: RM, Bits: 3},
			{Kind: D, Value: 1, Bits: 0}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	{Mnemonic: Mov, Desc: "Immediate to register",
		Bits: 4, Value: 0b1011,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: Reg, Bits: 3}, {Kind: D, Bits: 0, Value: 1},
			{Kind: Mod, Bits: 0, Value: 0b11}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	{Mnemonic: Mov, Desc: "Memory to accumulator",
		Bits: 7, Value: 0b1010000,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: Disp, Bits: 16},
			{Kind: Reg, Value: 0b000}, {Kind: D, Bits: 0, Value: 1},
			{Kind: Mod, Bits: 0, Value: 0b00}, {Kind: RM, Bits: 0, Value: 0b110}}},
	{Mnemonic: Mov, Desc: "Accumulator to memory",
		Bits: 7, Value: 0b1010001,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: D, Bits: 0, Value: 0}, {Kind: RM, Bits: 0, Value: 0b110},
			{Kind: Reg, Bits: 0, Value: 0b000}, {Kind: Mod, Bits: 0, Value: 0b00}, {Kind: Disp, Bits: 16}},
	},
	{Mnemonic: Mov, Desc: "Register/memory to segment register",
		Bits: 8, Value: 0b10001110,
		Fields: []Field{{Kind: Mod, Bits: 2}, {Kind: Fix, Value: 0, Bits: 1}, {Kind: SR, Bits: 2}, {Kind: RM, Bits: 3},
			{Kind: D, Bits: 0, Value: 1}},
	},
	{Mnemonic: Mov, Desc: "Segment register to register memory",
		Bits: 8, Value: 0b10001100,
		Fields: []Field{{Kind: Mod, Bits: 2}, {Kind: Fix, Value: 0, Bits: 1}, {Kind: SR, Bits: 2}, {Kind: RM, Bits: 3},
			{Kind: D, Bits: 0, Value: 0}},
	},
	// ADD
	{Mnemonic: Add, Desc: "Reg/memory with register to either",
		Bits: 6, Value: 0b0,
		Fields: []Field{{Kind: D, Bits: 1}, {Kind: W, Bits: 1},
			{Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3}, {Kind: RM, Bits: 3}},
	},
	{Mnemonic: Add, Desc: "Immediate to register/memory",
		Bits: 6, Value: 0b100000,
		Fields: []Field{{Kind: S, Bits: 1}, {Kind: W, Bits: 1},
			{Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3, Value: 0b000}, {Kind: RM, Bits: 3}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	{Mnemonic: Add, Desc: "Immediate to accumulator",
		Bits: 7, Value: 0b0000010,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: RM, Bits: 0, Value: 0b000},
			{Kind: Mod, Bits: 0, Value: 0b11}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	// SUB
	{Mnemonic: Sub, Desc: "Reg/memory and register to either",
		Bits: 6, Value: 0b001010,
		Fields: []Field{{Kind: D, Bits: 1}, {Kind: W, Bits: 1},
			{Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3}, {Kind: RM, Bits: 3}},
	},
	{Mnemonic: Sub, Desc: "Immediate from register/memory",
		Bits: 6, Value: 0b100000,
		Fields: []Field{{Kind: S, Bits: 1}, {Kind: W, Bits: 1},
			{Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3, Value: 0b101}, {Kind: RM, Bits: 3}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	{Mnemonic: Sub, Desc: "Immediate from accumulator",
		Bits: 7, Value: 0b0010110,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: Reg, Bits: 0, Value: 0b000}, {Kind: D, Bits: 0, Value: 1},
			{Kind: Mod, Bits: 0, Value: 0b11}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	// CMP
	{Mnemonic: Cmp, Desc: "Reg/memory and register",
		Bits: 6, Value: 0b001110,
		Fields: []Field{{Kind: D, Bits: 1}, {Kind: W, Bits: 1},
			{Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3}, {Kind: RM, Bits: 3}},
	},
	// Oo: The Value is the same for ADD! I don't know if I'm doing an add or a SUB! Fuck!
	{Mnemonic: Cmp, Desc: "Immediate from register/memory",
		Bits: 6, Value: 0b100000,
		Fields: []Field{{Kind: S, Bits: 1}, {Kind: W, Bits: 1},
			{Kind: Mod, Bits: 2}, {Kind: Reg, Bits: 3, Value: 0b111}, {Kind: RM, Bits: 3}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
	{Mnemonic: Cmp, Desc: "Immediate from accumulator",
		Bits: 7, Value: 0b0011110,
		Fields: []Field{{Kind: W, Bits: 1}, {Kind: Reg, Bits: 0, Value: 0b000}, {Kind: D, Bits: 0, Value: 1},
			{Kind: Mod, Bits: 0, Value: 0b11}, {Kind: Data, Bits: 8}, {Kind: DataW, Bits: 8}},
	},
}

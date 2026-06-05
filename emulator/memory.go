package emulator

import "fmt"

type Memory struct {
	segment [1024 * 1024]byte
}

func (m *Memory) Load(addr int16, wide bool) (uint16, error) {
	if addr < 0 || int(addr) >= len(m.segment) {
		return 0, fmt.Errorf("no memory address found")
	}
	lo := m.segment[addr]
	if !wide {
		return uint16(lo), nil
	}
	if int(addr+1) >= len(m.segment) {
		return 0, fmt.Errorf("no memory address found")
	}
	hi := m.segment[addr+1]
	return uint16(hi)<<8 | uint16(lo), nil
}

func (m *Memory) Store(addr int16, value uint16, w bool) error {
	if addr < 0 || int(addr) >= len(m.segment) {
		return fmt.Errorf("no memory address found")
	}
	if w && int(addr+1) >= len(m.segment) {
		return fmt.Errorf("Out of memory")
	}
	m.segment[addr] = byte(value & 0x00ff)
	if w {
		m.segment[addr+1] = byte(value >> 8)
		return nil
	}
	return nil
}

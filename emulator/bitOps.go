package emulator

func ReadBits(data []byte, bitPos int, n int) uint16 {
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

func ReadU16LE(data []byte, bitPos int) uint16 {
	lo := ReadBits(data, bitPos, 8)
	hi := ReadBits(data, bitPos+8, 8)
	return hi<<8 | lo
}

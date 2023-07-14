package main

func CalculateTableCRC8() []byte {
	const generator byte = 0x1D
	crcTable := make([]byte, 256)
	for dividend := 0; dividend < 256; dividend++ {
		currByte := byte(dividend)
		for bit := 0; bit < 8; bit++ {
			if (currByte & 0x80) != 0 {
				currByte <<= 1
				currByte ^= generator
			} else {
				currByte <<= 1
			}
		}
		crcTable[dividend] = currByte
	}
	return crcTable
}
func ComputeCRC8(bytes []byte) byte {
	crc := byte(0)
	for _, b := range bytes {
		data := b ^ crc
		crc = CrcLookup[data]
	}
	return crc
}

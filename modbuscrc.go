package modbusCRC

import (
	"bytes"
)

const polynomial = 0xA001

var crc16table [256]uint16

func makeCrc16Table() {
	var crc uint16
	for indexTable := 0; indexTable < 256; indexTable++ {
		crc = 0
		i := indexTable
		for j := 0; j < 8; j++ {
			if (crc^uint16(i))&0x0001 == 1 {
				crc = (crc >> 1) ^ polynomial
			} else {
				crc = crc >> 1
			}
			i = i >> 1
		}
		crc16table[indexTable] = crc
	}
}

func Calculate(data []byte) []byte {
	if crc16table[255] == 0x0000 {
		makeCrc16Table()
	}

	var CRC uint16 = 0xFFFF

	for _, v := range data {
		n := uint8(uint16(v) ^ CRC)
		CRC = CRC >> 8
		CRC = CRC ^ crc16table[n]
	}
	uint16ToBytes := make([]byte, 2)
	uint16ToBytes[1], uint16ToBytes[0] = uint8(CRC>>8), uint8(CRC&0xff)
	return uint16ToBytes
}

func Check(data []byte) bool {

	calculate := Calculate(data[:len(data)-2])

	request := data[len(data)-2:]

	if bytes.Equal(calculate, request) {
		return true
	} else {
		return false
	}
}

func Add(data []byte) []byte {
	CRC := Calculate(data)
	answer := append(data, CRC...)
	return answer
}

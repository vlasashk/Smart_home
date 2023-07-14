package main

import (
	"encoding/base64"
	"fmt"
)

func (data *Packet) Unmarshal(rawSrc []byte) (int, error) {
	var index uint = 1
	var err error = nil
	data.Length = rawSrc[0]
	data.Crc8 = rawSrc[data.Length+1]
	computeCrc := ComputeCRC8(rawSrc[1 : data.Length+1])
	if computeCrc == data.Crc8 {
		data.Payload.Src, index = unmarshalVaruint(rawSrc, uint(data.Length), index)
		data.Payload.Dst, index = unmarshalVaruint(rawSrc, uint(data.Length), index)
		data.Payload.Serial, index = unmarshalVaruint(rawSrc, uint(data.Length), index)
		data.Payload.DevType = rawSrc[index]
		data.Payload.Cmd = rawSrc[index+1]
		index += 2
		data.Payload.CmdBody = ParseCmdBody(data.Payload.DevType, data.Payload.Cmd, rawSrc[index:data.Length+1])
	} else {
		err = fmt.Errorf("ERROR: CRC8 values do not match. Calculated = %d, Received = %d", computeCrc, data.Crc8)
	}

	return int(data.Length + 2), err
}

func (data *Packet) Marshal() []byte {
	payloadRaw := make([]byte, 0, 8)
	marshalVaruint(data.Payload.Src, &payloadRaw)
	marshalVaruint(data.Payload.Dst, &payloadRaw)
	marshalVaruint(data.Payload.Serial, &payloadRaw)
	payloadRaw = append(payloadRaw, data.Payload.DevType)
	payloadRaw = append(payloadRaw, data.Payload.Cmd)
	if data.Payload.CmdBody != nil {
		data.Payload.CmdBody.MarshalInfo(&payloadRaw)
	}

	payloadLen := len(payloadRaw)
	rawData := make([]byte, payloadLen+2)
	rawData[0] = byte(payloadLen)
	copy(rawData[1:], payloadRaw)
	rawData[payloadLen+1] = ComputeCRC8(rawData[1 : payloadLen+1])
	return rawData
}

func marshalVaruint(src Varuint, rawData *[]byte) {
	octal := Varuint(0xFF)
	temp := make([]byte, 0, 8)
	tempLen := 0
	temp = append(temp, byte(src&octal))
	src >>= 8
	tempLen++
	for src != 0 {
		temp = append(temp, byte(src&octal))
		src >>= 8
		tempLen++
	}
	for i := len(temp) - 1; i >= 0; i-- {
		*rawData = append(*rawData, temp[i])
	}
}

func unmarshalVaruint(rawSrc []byte, length uint, index uint) (res Varuint, newIndex uint) {
	stopULEB := byte(0x80)
	newIndex = index
	for i := index; i < length; i++ {
		newIndex++
		res |= Varuint(rawSrc[i])
		if rawSrc[i]&stopULEB == 0 {
			break
		} else {
			res <<= 8
		}
	}
	return
}

func Base64UrlDecoder(rawSrc []byte) ([]Packet, error) {
	data := make([]Packet, 0)
	decoded, err := base64.RawURLEncoding.DecodeString(string(rawSrc[:]))
	totalLen := len(decoded)
	if err == nil {
		var processedSize, packetSize int
		for i := 0; totalLen > processedSize; i++ {
			data = append(data, Packet{})
			packetSize, err = data[i].Unmarshal(decoded[processedSize:])
			processedSize += packetSize
			if err != nil {
				break
			}
		}
	}
	return data, err
}

func Base64UrlEncoder(data []Packet) string {
	length := len(data)
	test := data[0].Marshal()
	for i := 1; i < length; i++ {
		test = append(test, data[i].Marshal()...)
	}
	encode := base64.RawURLEncoding.EncodeToString(test)
	return encode
}

func DecodeULEB128(value Varuint) (result uint64) {
	var maskULEB uint64 = 0x7F
	for {
		result |= maskULEB & uint64(value)
		value >>= 8
		if value == 0 {
			break
		} else {
			result <<= 7
		}
	}
	return result
}

func EncodeULEB128(value uint64) (result Varuint) {
	var maskULEB uint64 = 0x7F
	var highULEB uint64 = 0x80
	for {
		result |= Varuint(maskULEB & value)
		value >>= 7
		if value == 0 {
			break
		} else {
			result |= Varuint(highULEB)
			result <<= 8
		}
	}
	return result
}

package main

import (
	"encoding/base64"
	"fmt"
)

func (data *Packet) Unmarshal(rawSrc []byte) error {
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
		data.Payload.CmdBody = rawSrc[index : data.Length+1]
	} else {
		err = fmt.Errorf("ERROR: CRC8 values do not match. Calculated = %d, Received = %d", computeCrc, data.Crc8)
	}

	return err
}

func (data *Packet) Marshal() ([]byte, error) {
	rawData := make([]byte, data.Length+2)
	rawData[0] = data.Length
	var index uint = 1
	index = marshalVaruint(data.Payload.Src, rawData, index)
	index = marshalVaruint(data.Payload.Dst, rawData, index)
	index = marshalVaruint(data.Payload.Serial, rawData, index)
	rawData[index] = data.Payload.DevType
	rawData[index+1] = data.Payload.Cmd
	index += 2
	for i := 0; i < len(data.Payload.CmdBody); i++ {
		rawData[index] = data.Payload.CmdBody[i]
		index++
	}
	rawData[data.Length+1] = ComputeCRC8(rawData[1 : data.Length+1])
	return rawData, nil
}

func marshalVaruint(src Varuint, rawData []byte, index uint) uint {
	octal := Varuint(0xFF)
	temp := make([]byte, 0, 8)
	tempLen := 0
	temp = append(temp, rawData[index]|byte(src&octal))
	src >>= 8
	tempLen++
	for src != 0 {
		temp = append(temp, rawData[index]|byte(src&octal))
		src >>= 8
		tempLen++
	}
	for i := len(temp) - 1; i >= 0; i-- {
		rawData[index] = temp[i]
		index++
	}
	return index
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

func Base64UrlDecoder(rawSrc []byte) (Packet, error) {
	data := Packet{}
	decoded, err := base64.RawURLEncoding.DecodeString(string(rawSrc[:]))

	if err == nil {
		_ = data.Unmarshal(decoded)
	}
	return data, err
}

func Base64UrlEncoder(data Packet) string {
	test, _ := data.Marshal()
	encode := base64.RawURLEncoding.EncodeToString(test)
	return encode
}

package main

import (
	"encoding/base64"
)

/*
 =========================================
|                                         |
|              Packet struct              |
|     BASE64URL, ULEB128 (en)decoding     |
|                                         |
 =========================================
*/

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

package main

type Varuint uint64
type Bytes []byte

type DeviceInfo interface {
	Unmarshal() int
}

type String struct {
	Length byte
	Value  string
}

type Payload struct {
	Src     Varuint
	Dst     Varuint
	Serial  Varuint
	DevType byte
	Cmd     byte
	CmdBody DeviceInfo
}

type Device struct {
	DevName  String
	DevProps Bytes
}

type Array struct {
	Length   byte
	Elements []interface{} // This can hold any type
}

type Packet struct {
	Length  byte
	Payload Payload
	Crc8    byte
}

type TimerCmdBody struct {
	timestamp Varuint
}

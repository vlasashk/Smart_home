package main

type Varuint uint64
type Bytes []byte
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
	CmdBody Bytes
}

type Device struct {
	dev_name  String
	dev_props Bytes
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

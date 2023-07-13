package main

type Varuint uint64

type DeviceInfo interface {
	UnmarshalInfo(rawSrc *[]byte) int
	//MarshalInfo() int
}

type String struct {
	Length byte
	Value  string
}

type Packet struct {
	Length  byte
	Payload Payload
	Crc8    byte
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
	DevProps DeviceInfo
}

type Array struct {
	Length   byte
	Elements []interface{} // This can hold any type
}

type TimerCmdBody struct {
	Timestamp Varuint
}

type EnvSensorStatus struct {
	Values []Varuint
}

type switchOnOff struct {
	Status byte
}

type EnvSensorProps struct {
	Sensors  byte
	Triggers []struct {
		Op    byte
		Value Varuint
		Name  string
	}
}

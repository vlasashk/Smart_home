package main

type Varuint uint64

type DeviceInfo interface {
	UnmarshalInfo(rawSrc []byte)
	MarshalInfo(rawSrc *[]byte)
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
	DevName  string
	DevProps DeviceInfo
}
type PropsString struct {
	Length byte
	Name   []string
}

type TimerCmdBody struct {
	Timestamp Varuint
}

type EnvSensorStatus struct {
	Values []Varuint
}

type SwitchOnOff struct {
	Status byte
}

type EnvSensorProps struct {
	Sensors  byte
	Triggers []TriggersT
}

type TriggersT struct {
	Op    byte
	Value Varuint
	Name  string
}

var packetSerial uint64

var CrcLookup = CalculateTableCRC8()

const (
	WhoIsHereCMD = iota + 1
	IamHereCMD
	GetStatusCMD
	StatusCMD
	SetStatusCMD
	TickCMD
)
const (
	SmartHubDev = iota + 1
	EnvSensorDev
	SwitchDev
	LampDev
	SocketDev
	ClockDev
)

package main

/*
 =========================================
|                                         |
|             Project Structure           |
|                                         |
 =========================================
*/

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

type DeviceAddr struct {
	Address    Varuint
	DevType    byte
	Controlled []string
}

type SmartHub struct {
	PacketsQueue     *Queue
	CurrTime         uint64
	HubAddress       Varuint
	PacketSerial     uint64
	HubName          string
	ActiveDevices    map[string]DeviceAddr
	DeviceNames      map[Varuint]string
	AwaitingResponse map[string]uint64
	HubTriggers      map[Varuint]EnvSensorProps
}

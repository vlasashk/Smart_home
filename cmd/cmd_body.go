package main

const (
	WHOISHERE_CMD = iota + 1
	IAMHERE_CMD
	GETSTATUS_CMD
	STATUS_CMD
	SETSTATUS_CMD
	TICK_CMD
)
const (
	SmartHub_Dev = iota + 1
	EnvSensor_Dev
	Switch_Dev
	Lamp_Dev
	Socket_Dev
	Clock_Dev
)

func ParseCmdBody(device, cmd byte, rawSrc *[]byte) DeviceInfo {
	var res DeviceInfo
	switch {
	case device == SmartHub_Dev:
		res = &Device{}
	case device == EnvSensor_Dev && (cmd == WHOISHERE_CMD || cmd == IAMHERE_CMD):
		res = &Device{}
	case device == EnvSensor_Dev && cmd == STATUS_CMD:
		res = &EnvSensorStatus{}
	case device == Switch_Dev && (cmd == WHOISHERE_CMD || cmd == IAMHERE_CMD):
		res = &Device{}
	case device == Switch_Dev && cmd == STATUS_CMD:
		res = &switchOnOff{}
	case (device == Lamp_Dev || device == Socket_Dev) && (cmd == WHOISHERE_CMD || cmd == IAMHERE_CMD):
		res = &Device{}
	case (device == Lamp_Dev || device == Socket_Dev) && (cmd == STATUS_CMD || cmd == SETSTATUS_CMD):
		res = &switchOnOff{}
	case cmd == TICK_CMD:
		res = &TimerCmdBody{}
	}
	return res
}

func (timerType *TimerCmdBody) UnmarshalInfo(rawSrc *[]byte) int {
	return 6
}

func (deviceType *Device) UnmarshalInfo(rawSrc *[]byte) int {
	return 1
}
func (sensorStatusType *EnvSensorStatus) UnmarshalInfo(rawSrc *[]byte) int {
	return 1
}

func (onOffType *switchOnOff) UnmarshalInfo(rawSrc *[]byte) int {
	return 1
}

//func (sensorPropsType *EnvSensorProps) UnmarshalInfo(rawSrc *[]byte) int {
//	return 1
//}

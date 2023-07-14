package main

func ParseCmdBody(device, cmd byte, rawSrc []byte) DeviceInfo {
	var res DeviceInfo
	switch {
	case device == SmartHubDev:
		res = &Device{}
	case device == EnvSensorDev && (cmd == WhoIsHereCMD || cmd == IamHereCMD):
		res = &Device{
			DevProps: &EnvSensorProps{},
		}
	case device == EnvSensorDev && cmd == StatusCMD:
		res = &EnvSensorStatus{}
	case device == SwitchDev && (cmd == WhoIsHereCMD || cmd == IamHereCMD):
		res = &Device{
			DevProps: &PropsString{},
		}
	case device == SwitchDev && cmd == StatusCMD:
		res = &SwitchOnOff{}
	case (device == LampDev || device == SocketDev) && (cmd == WhoIsHereCMD || cmd == IamHereCMD):
		res = &Device{}
	case (device == LampDev || device == SocketDev) && (cmd == StatusCMD || cmd == SetStatusCMD):
		res = &SwitchOnOff{}
	case cmd == TickCMD:
		res = &TimerCmdBody{}
	case device == ClockDev && cmd == IamHereCMD:
		res = &Device{}
	}
	if res != nil {
		res.UnmarshalInfo(rawSrc)
	}
	return res
}

func (timerType *TimerCmdBody) UnmarshalInfo(rawSrc []byte) {
	timerType.Timestamp, _ = unmarshalVaruint(rawSrc, uint(len(rawSrc)), 0)
	return
}

func ParseString(rawSrc []byte) (res string, newIndex uint) {
	length := uint(rawSrc[0])
	newIndex = length + 1
	res = string(rawSrc[1:newIndex])
	return
}

func (deviceType *Device) UnmarshalInfo(rawSrc []byte) {
	var index uint
	deviceType.DevName, index = ParseString(rawSrc)
	switch deviceType.DevProps.(type) {
	case *EnvSensorProps:
		deviceType.DevProps.UnmarshalInfo(rawSrc[index:])
	case *PropsString:
		deviceType.DevProps.UnmarshalInfo(rawSrc[index:])
	}
}

func (sensorStatusType *EnvSensorStatus) UnmarshalInfo(rawSrc []byte) {
	length := rawSrc[0]
	var index uint = 1
	sensorStatusType.Values = make([]Varuint, length)
	for i := 0; i < int(length); i++ {
		sensorStatusType.Values[i], index = unmarshalVaruint(rawSrc, uint(len(rawSrc)), index)
	}

}

func (onOffType *SwitchOnOff) UnmarshalInfo(rawSrc []byte) {
	onOffType.Status = rawSrc[0]
}

func (sensorPropsType *EnvSensorProps) UnmarshalInfo(rawSrc []byte) {
	sensorPropsType.Sensors = rawSrc[0]
	length := rawSrc[1]
	var index uint = 2
	sensorPropsType.Triggers = make([]TriggersT, length)
	for i := 0; i < int(length); i++ {
		var processed uint
		sensorPropsType.Triggers[i], processed = ParseTriggers(rawSrc[index:])
		index += processed
	}
}

func ParseTriggers(rawSrc []byte) (res TriggersT, index uint) {
	res.Op = rawSrc[0]
	res.Value, index = unmarshalVaruint(rawSrc, uint(len(rawSrc)), 1)
	strSize := uint(rawSrc[index])
	res.Name, _ = ParseString(rawSrc[index:])
	index += strSize + 1
	return
}

func (propsStrType *PropsString) UnmarshalInfo(rawSrc []byte) {
	propsStrType.Length = rawSrc[0]
	var arrIndex uint = 1
	propsStrType.Name = make([]string, propsStrType.Length)
	for i := uint(0); i < uint(propsStrType.Length); i++ {
		var processed uint
		propsStrType.Name[i], processed = ParseString(rawSrc[arrIndex:])
		arrIndex += processed
	}
}

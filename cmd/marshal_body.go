package main

func (timerType *TimerCmdBody) MarshalInfo(rawSrc *[]byte) {
	marshalVaruint(timerType.Timestamp, rawSrc)
}

func (deviceType *Device) MarshalInfo(rawSrc *[]byte) {
	MarshalString(deviceType.DevName, rawSrc)
	switch deviceType.DevProps.(type) {
	case *EnvSensorProps:
		deviceType.DevProps.MarshalInfo(rawSrc)
	case *PropsString:
		deviceType.DevProps.MarshalInfo(rawSrc)
	}
}

func (sensorStatusType *EnvSensorStatus) MarshalInfo(rawSrc *[]byte) {
	length := len(sensorStatusType.Values)
	*rawSrc = append(*rawSrc, byte(length))
	for i := 0; i < length; i++ {
		marshalVaruint(sensorStatusType.Values[i], rawSrc)
	}
}

func (onOffType *SwitchOnOff) MarshalInfo(rawSrc *[]byte) {
	*rawSrc = append(*rawSrc, onOffType.Status)
}

func (sensorPropsType *EnvSensorProps) MarshalInfo(rawSrc *[]byte) {
	*rawSrc = append(*rawSrc, sensorPropsType.Sensors)
	length := len(sensorPropsType.Triggers)
	*rawSrc = append(*rawSrc, byte(length))
	for i := 0; i < length; i++ {
		*rawSrc = append(*rawSrc, sensorPropsType.Triggers[i].Op)
		marshalVaruint(sensorPropsType.Triggers[i].Value, rawSrc)
		MarshalString(sensorPropsType.Triggers[i].Name, rawSrc)
	}
}

func (propsStrType *PropsString) MarshalInfo(rawSrc *[]byte) {
	*rawSrc = append(*rawSrc, propsStrType.Length)
	for i := 0; i < int(propsStrType.Length); i++ {
		MarshalString(propsStrType.Name[i], rawSrc)
	}
}

func MarshalString(src string, rawData *[]byte) {
	temp := []byte(src)
	length := byte(len(temp))
	if length > 0 {
		*rawData = append(*rawData, length)
		*rawData = append(*rawData, temp...)
	}
}

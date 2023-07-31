package main

import "fmt"

/*
 =========================================
|                                         |
|      Packet struct (UN)MARSHALING       |
|                                         |
 =========================================
*/

func (data *Packet) Unmarshal(rawSrc []byte) (int, error) {
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
		data.Payload.CmdBody = ParseCmdBody(data.Payload.DevType, data.Payload.Cmd, rawSrc[index:data.Length+1])
	} else {
		err = fmt.Errorf("ERROR: CRC8 values do not match. Calculated = %d, Received = %d", computeCrc, data.Crc8)
	}

	return int(data.Length + 2), err
}

func (data *Packet) Marshal() []byte {
	payloadRaw := make([]byte, 0, 8)
	marshalVaruint(data.Payload.Src, &payloadRaw)
	marshalVaruint(data.Payload.Dst, &payloadRaw)
	marshalVaruint(data.Payload.Serial, &payloadRaw)
	payloadRaw = append(payloadRaw, data.Payload.DevType)
	payloadRaw = append(payloadRaw, data.Payload.Cmd)
	if data.Payload.CmdBody != nil {
		data.Payload.CmdBody.MarshalInfo(&payloadRaw)
	}

	payloadLen := len(payloadRaw)
	rawData := make([]byte, payloadLen+2)
	rawData[0] = byte(payloadLen)
	copy(rawData[1:], payloadRaw)
	rawData[payloadLen+1] = ComputeCRC8(rawData[1 : payloadLen+1])
	return rawData
}

func marshalVaruint(src Varuint, rawData *[]byte) {
	octal := Varuint(0xFF)
	temp := make([]byte, 0, 8)
	tempLen := 0
	temp = append(temp, byte(src&octal))
	src >>= 8
	tempLen++
	for src != 0 {
		temp = append(temp, byte(src&octal))
		src >>= 8
		tempLen++
	}
	for i := len(temp) - 1; i >= 0; i-- {
		*rawData = append(*rawData, temp[i])
	}
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

func MarshalString(src string, rawData *[]byte) {
	temp := []byte(src)
	length := byte(len(temp))
	if length > 0 {
		*rawData = append(*rawData, length)
		*rawData = append(*rawData, temp...)
	}
}

/*
 =========================================
|                                         |
|    CMD_BODY INTERFACE (UN)MARSHALING    |
|                                         |
 =========================================
*/

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

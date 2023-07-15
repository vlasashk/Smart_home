package main

func (hub *SmartHub) RequestHandler(req []byte) {

}

func (hub *SmartHub) SendSetGetPacket(packets []Packet) string {
	length := len(packets)
	for i := 0; i < length; i++ {
		devName, ok := hub.DeviceNames[packets[i].Payload.Dst]
		if ok {
			_, waiting := hub.AwaitingResponse[devName]
			if !waiting {
				hub.AwaitingResponse[devName] = 0
			}
		}
	}
	return Base64UrlEncoder(packets)
}

func (hub *SmartHub) ReceiveIamHere(data Packet) {
	switch dev := data.Payload.CmdBody.(type) {
	case *Device:
		temp := DeviceAddr{
			Address:    data.Payload.Src,
			DevType:    data.Payload.DevType,
			Controlled: nil,
		}
		if data.Payload.DevType == SwitchDev {
			switch body := dev.DevProps.(type) {
			case *PropsString:
				temp.Controlled = body.Name
			}
		}
		hub.ActiveDevices[dev.DevName] = temp
		hub.DeviceNames[data.Payload.Src] = dev.DevName
		packet, ok := hub.makeGetStatusPack(dev.DevName)
		if ok {
			packets := make([]Packet, 1)
			packets[0] = packet
			hub.PacketsQueue.AddPack(packets)
		}
	}
}

func (hub *SmartHub) ReceiveTick(clockTick TimerCmdBody) {
	hub.CurrTime = DecodeULEB128(clockTick.Timestamp)
	for device, sentTime := range hub.AwaitingResponse {
		if sentTime == 0 {
			sentTime = hub.CurrTime
		} else if _, ok := hub.ActiveDevices[device]; ok && (hub.CurrTime-sentTime > 300) {
			delete(hub.DeviceNames, hub.ActiveDevices[device].Address)
			delete(hub.ActiveDevices, device)
		}
	}
}

func (hub *SmartHub) ReceiveStatus(pack Packet) {
	srcName, found := hub.DeviceNames[pack.Payload.Src]
	if found {
		_, waiting := hub.AwaitingResponse[srcName]
		if waiting {
			delete(hub.AwaitingResponse, srcName)
		}
	}
	switch pack.Payload.DevType {
	case SwitchDev:
		switch dev := pack.Payload.CmdBody.(type) {
		case *SwitchOnOff:
			packets := hub.makeSwitchPacks(dev.Status, pack.Payload.Src)
			if len(packets) > 0 {
				hub.PacketsQueue.AddPack(packets)
			}
		}
	case EnvSensorDev:
		switch dev := pack.Payload.CmdBody.(type) {
		case *EnvSensorStatus:
			if len(hub.HubTriggers.Triggers) > 0 {
				hub.TriggerResponseAction(*dev)
			}
		}
	}
}

func (hub *SmartHub) makeSwitchPacks(onOff byte, switchAddr Varuint) []Packet {
	switchName := hub.DeviceNames[switchAddr]
	dev, ok := hub.ActiveDevices[switchName]
	var packets []Packet
	if ok {
		packets = make([]Packet, 0, len(dev.Controlled))
		for i := 0; i < len(dev.Controlled); i++ {
			temp, success := hub.makeSetStatusPack(onOff, dev.Controlled[i])
			if success {
				packets = append(packets, temp)
			}
		}
	}
	return packets
}

func (hub *SmartHub) makeGetStatusPack(dstName string) (Packet, bool) {
	dev, ok := hub.ActiveDevices[dstName]
	getPacket := Packet{}
	if dev.DevType == SmartHubDev || dev.DevType == ClockDev {
		ok = false
	}
	if ok {
		getPacket.Payload.Src = hub.HubAddress
		getPacket.Payload.Dst = dev.Address
		getPacket.Payload.Serial = EncodeULEB128(hub.PacketSerial)
		hub.PacketSerial++
		getPacket.Payload.DevType = dev.DevType
		getPacket.Payload.Cmd = GetStatusCMD
	}
	return getPacket, ok
}

func (hub *SmartHub) makeSetStatusPack(onOff byte, dstName string) (Packet, bool) {
	dev, ok := hub.ActiveDevices[dstName]
	setPacket := Packet{}
	if dev.DevType != LampDev && dev.DevType != SocketDev {
		ok = false
	}
	if ok {
		setPacket.Payload.Src = hub.HubAddress
		setPacket.Payload.Dst = dev.Address
		setPacket.Payload.Serial = EncodeULEB128(hub.PacketSerial)
		hub.PacketSerial++
		setPacket.Payload.DevType = dev.DevType
		setPacket.Payload.Cmd = SetStatusCMD
		setPacket.Payload.CmdBody = &SwitchOnOff{
			Status: onOff,
		}
	}
	return setPacket, ok
}

func (hub *SmartHub) TriggerResponseAction(value EnvSensorStatus) {
	trigger := hub.HubTriggers
	envInfo := parseEnvValues(trigger, value)
	length := len(trigger.Triggers)
	typeMask := byte(0xC)
	for i := 0; i < length; i++ {
		packets := make([]Packet, 1)
		var success bool
		turnOnOff := trigger.Triggers[i].Op & 1
		greater := trigger.Triggers[i].Op & 2 >> 1
		sensType := trigger.Triggers[i].Op & typeMask >> 2
		triggerVal := DecodeULEB128(trigger.Triggers[i].Value)
		if greater != 0 {
			if v, ok := envInfo[sensType]; ok && v > triggerVal {
				packets[0], success = hub.makeSetStatusPack(turnOnOff, trigger.Triggers[i].Name)
				if success {
					hub.PacketsQueue.AddPack(packets)
				}
			}
		} else {
			if v, ok := envInfo[sensType]; ok && v < triggerVal {
				packets[0], success = hub.makeSetStatusPack(turnOnOff, trigger.Triggers[i].Name)
				if success {
					hub.PacketsQueue.AddPack(packets)
				}
			}
		}
	}
}

func parseEnvValues(trigger EnvSensorProps, value EnvSensorStatus) map[byte]uint64 {
	temp := byte(0x1)
	humidity := byte(0x2)
	light := byte(0x4)
	air := byte(0x8)
	envInfo := make(map[byte]uint64)
	if trigger.Sensors != 0 {
		index := 0
		if (trigger.Sensors & temp) != 0 {
			envInfo[0] = DecodeULEB128(value.Values[index])
			index++
		}
		if (trigger.Sensors & humidity) != 0 {
			envInfo[1] = DecodeULEB128(value.Values[index])
			index++
		}
		if (trigger.Sensors & light) != 0 {
			envInfo[2] = DecodeULEB128(value.Values[index])
			index++
		}
		if (trigger.Sensors & air) != 0 {
			envInfo[3] = DecodeULEB128(value.Values[index])
			index++
		}
	}
	return envInfo
}

func (hub *SmartHub) SendWhoIsHere() string {
	res := make([]Packet, 1)
	res[0].Payload.Src = hub.HubAddress
	res[0].Payload.Dst = EncodeULEB128(0x3FFF)
	res[0].Payload.Serial = EncodeULEB128(hub.PacketSerial)
	hub.PacketSerial++
	res[0].Payload.DevType = SmartHubDev
	res[0].Payload.Cmd = WhoIsHereCMD
	res[0].Payload.CmdBody = &Device{
		DevName: hub.HubName,
	}
	return Base64UrlEncoder(res)
}

func (hub *SmartHub) SendIamHere() string {
	res := make([]Packet, 1)
	res[0].Payload.Src = hub.HubAddress
	res[0].Payload.Dst = EncodeULEB128(0x3FFF)
	res[0].Payload.Serial = EncodeULEB128(hub.PacketSerial)
	hub.PacketSerial++
	res[0].Payload.DevType = SmartHubDev
	res[0].Payload.Cmd = IamHereCMD
	res[0].Payload.CmdBody = &Device{
		DevName: hub.HubName,
	}
	return Base64UrlEncoder(res)
}

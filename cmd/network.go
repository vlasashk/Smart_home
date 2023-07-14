package main

func InitialPacket(hubAddress uint64) string {
	res := make([]Packet, 1)
	res[0].Payload.Src = EncodeULEB128(hubAddress)
	res[0].Payload.Dst = EncodeULEB128(0x3FFF)
	res[0].Payload.Serial = EncodeULEB128(packetSerial)
	res[0].Payload.DevType = 1
	res[0].Payload.Cmd = 1
	res[0].Payload.CmdBody = &Device{
		DevName: "SmartHub",
	}
	return Base64UrlEncoder(res)
}

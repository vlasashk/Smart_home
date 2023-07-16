package main

import (
	"net/http"
	"os"
	"strconv"
)

func main() {
	argsToRun := os.Args
	if len(argsToRun) < 3 {
		os.Exit(99)
	}
	url := os.Args[1]
	srcAddress, _ := strconv.ParseUint(os.Args[2], 16, 64)
	hub := InitHub(srcAddress)
	hub.SendWhoIsHere()
	client := &http.Client{}
	for {
		resp := hub.SendHandler(client, url)
		hub.ResponseHandler(resp)
	}
}

func InitHub(address uint64) SmartHub {
	return SmartHub{
		PacketsQueue:     &Queue{},
		CurrTime:         0,
		HubAddress:       EncodeULEB128(address),
		PacketSerial:     1,
		HubName:          "SmartHub",
		ActiveDevices:    make(map[string]DeviceAddr),
		DeviceNames:      make(map[Varuint]string),
		AwaitingResponse: make(map[string]uint64),
		HubTriggers:      make(map[Varuint]EnvSensorProps),
	}
}

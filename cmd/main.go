package main

import (
	"bytes"
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

func MakeHttpReq(url string, data []byte) *http.Request {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		os.Exit(99)
	}
	req.Header.Set("Content-Type", "text/plain")
	return req
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
		HubTriggers:      EnvSensorProps{},
	}
}

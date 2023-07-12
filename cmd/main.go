package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	argsToRun := os.Args
	if len(argsToRun) < 2 {
		fmt.Println("Please provide a URL")
		os.Exit(1)
	}
	url := os.Args[1]
	data := []byte("C7MG_383AQEDSFVCuQ")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Panic("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response Status: %s\n", body)
	packet, _ := Base64UrlDecoder(body)
	fmt.Println(packet)
	//encoded := Base64UrlEncoder(packet)
	//fmt.Printf("%X %X %X %X %X %X\n", packet.Payload.Src, packet.Payload.Dst, packet.Payload.Serial, packet.Payload.DevType, packet.Payload.Cmd, packet.Payload.CmdBody)
	//fmt.Println(string(body))
	//fmt.Println(encoded)
}

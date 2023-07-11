package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Varuint uint64 // Assuming the maximum value is within 32 bits
type Bytes []byte
type String struct {
	Length byte
	Value  string
}

type Payload struct {
	Src     Varuint
	Dst     Varuint
	Serial  Varuint
	DevType byte
	Cmd     byte
	CmdBody Bytes
}

// Define the array type
type Array struct {
	Length   byte
	Elements []interface{} // This can hold any type
}

// Define the packet and payload structures
type Packet struct {
	Length  byte
	Payload Payload
	Src8    byte
}

func main() {
	url := "http://localhost:9998"
	data := []byte("")

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
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	data := []byte(InitialPacket(srcAddress))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		os.Exit(99)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		os.Exit(99)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	packet, _ := Base64UrlDecoder(body)

	for {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			os.Exit(99)
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err = client.Do(req)
		if err != nil {
			os.Exit(99)
		}

		if resp.StatusCode == http.StatusNoContent {
			os.Exit(0)
		} else if resp.StatusCode != http.StatusOK {
			os.Exit(99)
		}
		body, err = ioutil.ReadAll(resp.Body)

		packet, _ = Base64UrlDecoder(body)
		encoded := Base64UrlEncoder(packet)
		fmt.Printf("Response Status: \n%s\n", body)
		fmt.Println(encoded)
		fmt.Println()
	}
}

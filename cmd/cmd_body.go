package main

func ParseCmdBody(device, cmd byte, rawSrc []byte) DeviceInfo {
	var res DeviceInfo
	switch cmd {
		
	}
	return res
}

func (timerType TimerCmdBody) GetType() int {
	return 6
}

func (deviceType Device) GetType() int {
	return 1
}

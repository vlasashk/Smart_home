# Smart home simulation
## Build
1. Clone project:
```
git clone https://github.com/vlasashk/Smart_home.git
cd Smart_home
```
2. Run demo-server `smarthome` to simulate a network:
```
./smarthome -s -V -S 2
```
3. Run hub with 2 args (demo-server address:port, hub hexadecimal 14-bit index):
```
go run cmd/* http://localhost:9998 ef0
```
## Project information
Simulation of smart home hub functions. A smart house contains 
- sensors - measure environmental parameters
- actuators - affect the environment in the house 
- timer - transmits current time readings
- hub - controls all other devices in the smart house.

All the devices are in a common communication environment (network).
Information on the network is transmitted using packets. 

In this task, receiving data from the network and sending data to the network is modeled using an HTTP POST request to a special server that simulates the functioning of the other smart home devices. (Provided by Tinkoff education, more info about server functions here: https://github.com/blackav/smart-home-binary)

To send data to the network, POST request must be executed, passing in the body of the request the data packets sent to the network in the form of a URL-encoded Base64 string.
In response POST-request will return the data packets received from the network.
Each portion of the input data is encoded using URL-encoded unpadded Base64.

The device must respond to a broadcast `WHOISHERE` request within 300 ms at the latest. The device must respond to `GETSTATUS` or `SETSTATUS` requests addressed to this device within 300 ms at the latest. The `TICK` command does not require a response.
### Project Features
- Bitstream marshaling/unmarshaling
- Interfaces usage
- Queue implementation
- ULEB128 encoding/decoding
- Base64URL
- CRC8 checksum
- HTTP POST requests

### Data format
Each packet transmitted over the network is described as follows:
```
type packet struct {
    length byte
    payload bytes
    crc8 byte
};
```
Where:
- `length` is the size of the payload field in octets (bytes);
- `payload` - data transmitted in the packet, the specific data format for each packet type is described below;
  ```
    type payload struct {
      src varuint
      dst varuint
      serial varuint
      dev_type byte
      cmd byte
      cmd_body bytes
    };
  ```
- `varuint` - unsigned integer in ULEB128 format.
- `crc8` - checksum of the payload field calculated using the cyclic redundancy check 8 algorithm.

### Functionality
- 0x01 - WHOISHERE - sent by a device wishing to discover its neighbors on the network.
The dst distribution address must be broadcast 0x3FFF.
The cmd_body field describes the characteristics of the device itself in the form of a structure:
    ```
    type device struct {
        dev_name string
        dev_props bytes
    };
    ```
    > The content of dev_props is defined depending on the device type.

- 0x02 - `IAMHERE` - sent by the device that received the WHOISHERE command and contains information about the device itself in the cmd_body field. The IAMHERE command is sent strictly in response to WHOISHERE. The command is sent to a broadcast address.
- 0x03 - `GETSTATUS` - sent by the hub to some device to read the device status. If the device does not support the GETSTATUS command (for example, a timer), the command is ignored.
- 0x04 - `STATUS` - sent by the device to the hub both as a response to GETSTATUS, SETSTATUS requests and independently when the device state changes. For example, a switch sends a STATUS message at the moment of switching. In this case, the recipient address is the device that last sent the GETSTATUS command to this device. If no such command has yet been received, the STATUS message is not sent to anyone.
- 0x05 - `SETSTATUS` - sent by the hub to some device to make the device change its state, for example, to turn on a lamp. If the device does not support state change (such as a timer), the command is ignored.
- 0x06 - `TICK` - timer tick, sent by the timer. The frequency of sending is not guaranteed, but if an event is scheduled at some point in time, the event should be triggered when the time transmitted in the TICK message becomes greater than or equal to the scheduled time. The cmd_body field contains the following data:
    ```
    type timer_cmd_body struct {
      timestamp varuint
    };
    ```
### Device list and it's available commands along with `dev_props` content

- 0x01 - SmartHub - device that program simulates
  - WHOISHERE, IAMHERE — `dev_props`: empty
- 0x02 - EnvSensor - sensor of environmental characteristics (temperature, humidity, light, air pollution);
  - WHOISHERE, IAMHERE — `dev_props`:
  ```
  type env_sensor_props struct {
    sensors byte
    triggers [] struct {
        op byte
        value varuint
        name string
    }
  };
  ```
  > - The `sensors` field contains a bit mask of supported sensors, where the values of each bit indicate the following:
  >  - 0x1 - there is a temperature sensor (sensor 0);
  >  - 0x2 - there is a humidity sensor (sensor 1); 
  >  - 0x4 - there is a light sensor (sensor 2); 
  >  - 0x8 - there is an air pollution sensor (sensor 3).
  > - The `triggers` field is an array of sensor thresholds to trigger. 
  >   - Here `op` operation has the following format:
  >    - Bit 0 (low) - enable or disable the device;
  >    - Bit 1 - compare by condition less (0) or more (1); 
  >    - Bits 2-3 - sensor type 
  >   - `value` - this is the threshold value of the sensor; 
  >   - `name` - name of the device to be enabled or disabled.
  
  - STATUS - cmd_body field contains readings of all sensors supported by the device as an array of integers.
  ```
  type env_sensor_status_cmd_body struct {
    values []varuint
  };
  ```
- 0x03 - Switch; 
  - WHOISHERE, IAMHERE - `dev_props` is an array of strings. Each string is the name (dev_name) of a device that is connected to this switch. Turning the switch on should turn all devices on, and turning it off should turn them off.
  - STATUS - The `cmd_body` field is 1 byte in size and contains a value of 0 if the switch is in the OFF position and 1 if the switch is in the ON position.
- 0x04 - Lamp;
  - WHOISHERE, IAMHERE — `dev_props`: empty
  - STATUS - The cmd_body field is 1 byte in size and contains a value of 0 if the switch is in the OFF position and 1 if the switch is in the ON position.
  - SETSTATUS - the cmd_body field must be 1 byte in size and contain 0 to turn the device off and 1 to turn the device on.
- 0x05 - Socket;
  - WHOISHERE, IAMHERE — `dev_props`: empty
  - STATUS - The cmd_body field is 1 byte in size and contains a value of 0 if the switch is in the OFF position and 1 if the switch is in the ON position.
  - SETSTATUS - the cmd_body field must be 1 byte in size and contain 0 to turn the device off and 1 to turn the device on.
- 0x06 - Clock - the clock that broadcasts TICK messages. The `cmd_body` contains the following data:
  ```
  type timer_cmd_body struct {
    timestamp varuint
  };
  ```
  

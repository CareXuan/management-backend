package common

import (
	"encoding/json"
	"net"
)

type TcpCommonData struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type ElectricVehicleTcpData struct {
	DeviceId int    `json:"device_id"`
	Msg      string `json:"msg"`
}

func TcpRequest(host string, port string, data TcpCommonData) error {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()
	ss, _ := json.Marshal(data)
	_, err = conn.Write(ss)
	if err != nil {
		return err
	}
	return nil
}

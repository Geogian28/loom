package main

import (
	"fmt"
	"log"
	"time"

	arista "github.com/aristanetworks/goeapi"
)

type deviceInfo struct {
	hostname string
	model    string
	port     []portInfo
	mac      string
}

type portInfo struct {
	number string
	name   string
	status string
	vendor string
	model  string
	serial string
	media  string
	speed  string
}

type lldpNeighbor struct {
	name   string
	mac    string
	portID string
}

type switchScanner interface {
	getHostname() (string, error)
	getModel() (string, error)
	getPorts() ([]portInfo, error)
	getSystemStats() error
}

type aristaScanner struct {
	address       string
	port          int
	username      string
	password      string
	node          *arista.Node
	conn          *arista.EapiConnection
	handle        *arista.EapiReqHandle
	runningConfig string
}

// Layout defines the physical arrangement of ports for a specific model
type Layout struct {
	Model      string      `json:"model"`
	Groups     []PortGroup `json:"groups"`
	UplinkType string      `json:"uplinkType"` // "QSFP" or "SFP"
}

type PortGroup struct {
	StartPort int `json:"start"`
	EndPort   int `json:"end"`
}

var ChassisTemplates = map[string]Layout{
	"DCS-7050S-64": {
		Model: "DCS-7050S-64",
		Groups: []PortGroup{
			{StartPort: 1, EndPort: 16},
			{StartPort: 17, EndPort: 32},
			{StartPort: 33, EndPort: 48},
		},
		UplinkType: "QSFP", // Represents ports 49-64
	},
	"DCS-7050TX-48": {
		Model: "DCS-7050TX-48",
		Groups: []PortGroup{
			{StartPort: 1, EndPort: 48},
		},
		UplinkType: "SFP",
	},
}

func detectSwitchType(address string, port int, username string, password string) (switchScanner, error) {
	return &aristaScanner{
		address:  address,
		port:     port,
		username: username,
		password: password,
	}, nil
}

func getDeviceInfo(address string, port int, username string, password string) (deviceInfo, error) {
	mySwitch, err := detectSwitchType(address, port, username, password)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	hostname, err := mySwitch.getHostname()
	fmt.Println("getHostname Duration: ", time.Since(startTime))
	if err != nil {
		log.Fatal(err)
		return deviceInfo{}, err
	}

	startTime = time.Now()
	model, err := mySwitch.getModel()
	fmt.Println("getModel Duration: ", time.Since(startTime))
	if err != nil {
		log.Fatal(err)
		return deviceInfo{}, err
	}

	startTime = time.Now()
	ports, err := mySwitch.getPorts()
	fmt.Println("getPorts Duration: ", time.Since(startTime))
	if err != nil {
		log.Fatal(err)
		return deviceInfo{}, err
	}

	di := deviceInfo{
		hostname: hostname,
		model:    model,
		port:     ports,
		mac:      "de:ad:be:ef:00:00",
	}
	fmt.Println(di.hostname)
	return di, nil
}

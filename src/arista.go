package main

import (
	"fmt"
	"strconv"
	"strings"

	arista "github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
)

type ShowHostnameResp struct {
	HostName string
	FQDN     string
}

func (s *ShowHostnameResp) GetCmd() string {
	return "show hostname"
}

type ShowVersionResp struct {
	ModelName        string
	InternalVersion  string
	SystemMacAddress string
	SerialNumber     string
	MemTotal         int
	BootupTimestamp  float64
	MemFree          int
	Version          string
	Architecture     string
	InternalBuildID  string
	HardwareRevision string
}

func (s *ShowVersionResp) GetCmd() string {
	return "show version"
}

var speedMap = map[int]string{
	100000000:     "100meg",
	1000000000:    "1gig",
	2500000000:    "2.5gig",
	5000000000:    "4gig",
	10000000000:   "10gig",
	25000000000:   "25gig",
	40000000000:   "40gig",
	50000000000:   "50gig",
	100000000000:  "100gig",
	200000000000:  "200gig",
	400000000000:  "400gig",
	800000000000:  "800gig",
	1600000000000: "1600gig",
}

func (a *aristaScanner) establishAristaConnection() error {
	var err error
	a.node, err = arista.Connect(
		"http",
		a.address,
		a.username,
		a.password,
		a.port,
	)
	if err != nil {
		fmt.Println("EstablishAristaConnection error: ", err)
		return err
	}
	a.handle, err = arista.GetHandle(a.node, "json")
	if err != nil {
		fmt.Println("EstablishAristaConnection error: ", err)
		return err
	}
	return nil
}

func (a *aristaScanner) getRunningConfig() error {
	a.runningConfig = a.node.RunningConfig()
	if a.runningConfig == "" {
		return fmt.Errorf("received empty configuration from switch")
	}
	fmt.Println(a.runningConfig)

	return nil
}

func (a *aristaScanner) getHostname() (string, error) {
	err := a.establishAristaConnection()
	if err != nil {
		fmt.Println("Unable to establish connection: ", err)
		return "", err
	}
	sh := ShowHostnameResp{}
	handle, _ := a.node.GetHandle("json")
	handle.AddCommand(&sh)
	fmt.Println("")
	if err := handle.Call(); err != nil {
		panic(err)
	}
	fmt.Println("show hostname: ", sh.HostName)
	fmt.Println("show FQDN: ", sh.FQDN)
	return sh.HostName, nil
}

// func (a *aristaScanner) getHostname() (string, error) {
// 	err := a.establishAristaConnection()
// 	if err != nil {
// 		fmt.Println("Unable to establish connection: ", err)
// 		return "", err
// 	}

// 	// hostname := string(module.System(a.node).Get().HostName())
// 	// return hostname, nil

// 	cmds := []string{"show hostname"}
// 	res, err := a.node.RunCommands(cmds, "json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	var sv module.SystemConfig
// 	mapstructure.Decode(res, &sv)
// 	return sv.HostName(), nil
// }

func (a *aristaScanner) getModel() (string, error) {
	err := a.establishAristaConnection()
	if err != nil {
		fmt.Println("Unable to establish connection: ", err)
		return "", err
	}
	sv := module.ShowVersion{}
	handle, _ := a.node.GetHandle("json")
	handle.AddCommand(&sv)
	fmt.Println("")
	if err := handle.Call(); err != nil {
		panic(err)
	}
	fmt.Println("ModelName: ", sv.ModelName)
	return sv.ModelName, nil
}

// func (a *aristaScanner) getModel() (string, error) {
// 	err := a.establishAristaConnection()
// 	if err != nil {
// 		fmt.Println("Unable to establish connection: ", err)
// 		return "", err
// 	}

// 	var sv module.ShowVersion
// 	handle, _ := a.node.GetHandle("json")
// 	handle.AddCommand(&sv)

// 	cmds := []string{"show version"}
// 	res, err := a.node.RunCommands(cmds, "json")
// 	fmt.Println("")
// 	// if model1, ok := res.Result[0]["modelName"].(map[string]interface{})["Ethernet40"].(map[string]interface{}); ok {
// 	mapstructure.Decode(res.Result[0], &sv)
// 	fmt.Println("ModelName: ", sv.ModelName)
// 	return sv.ModelName, nil
// }

func (a *aristaScanner) ShowLLDPNeighbors() error {
	// if a.conn == nil {
	err := a.establishAristaConnection()
	if err != nil {
		fmt.Println("getHostname error: ", err)
		return err
	}
	// }
	if a.runningConfig == "" {
		err := a.getRunningConfig()
		if err != nil {
			fmt.Println("getHostname error: ", err)
			return err
		}
	}
	if a.node == nil {
		fmt.Println("a.node is nil")
		return err
	}
	// fmt.Println(module.Show(a.node).Version())
	var sv module.ShowLLDPNeighbors
	handle, _ := a.node.GetHandle("json")
	handle.AddCommand(&sv)
	err = handle.Call()
	if err != nil {
		fmt.Println("getHostname error: ", err)
		return err
	}
	for k, v := range sv.LLDPNeighbors {
		fmt.Printf("Interface: %s, Mode: %s \n", strconv.Itoa(k), v.NeighborDevice)
	}

	fmt.Println("getHostname no error")
	return nil
}

// func (a *aristaScanner) getPorts() ([]portInfo, error) {
// 	err := a.establishAristaConnection()
// 	if err != nil {
// 		fmt.Println("Unable to establish connection: ", err)
// 		return nil, err
// 	}

// 	// resp := module.SwitchPort(a.node).GetAll()
// 	// for _, switchport := range resp {
// 	// 	fmt.Println("switchport.Name: ", switchport.Name())
// 	// 	fmt.Println("switchport.Mode: ", switchport.Mode())
// 	// }
// 	var sv module.ShowInterface
// 	handle, _ := a.node.GetHandle("json")
// 	handle.AddCommand(&sv)

// 	cmds := []string{"show interfaces"}
// 	res, err := a.node.RunCommands(cmds, "json")
// 	// fmt.Println("res: ", res)
// 	fmt.Println("")

// 	interfacesRaw, ok := res.Result[0]["interfaces"].(map[string]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("could not find interfaces key in response")
// 	}
// 	var allPortInfos []portInfo

// 	for name := range interfacesRaw {
// 		var switchportEt1 module.SwitchInterface
// 		if ethernet1, ok := res.Result[0]["interfaces"].(map[string]interface{})[name].(map[string]interface{}); ok {
// 			mapstructure.Decode(ethernet1, &switchportEt1)
// 		}
// 		p := portInfo{
// 			number: strings.TrimPrefix(name, "Ethernet"),
// 			name:   name,
// 			speed:  fmt.Sprintf("%v", switchportEt1.Bandwidth),
// 		}
// 		if p.speed, ok = speedMap[switchportEt1.Bandwidth]; !ok {
// 			p.speed = fmt.Sprintf("%v", switchportEt1.Bandwidth)
// 		}
// 		allPortInfos = append(allPortInfos, p)
// 	}

// 	// for _, port := range allPortInfos {
// 	// 	fmt.Println("port.name: ", port.name)
// 	// 	fmt.Println("port.number: ", port.number)
// 	// 	fmt.Println("port.speed: ", port.speed)
// 	// 	fmt.Println("")
// 	// }
// 	fmt.Println("getPorts no error")
// 	return allPortInfos, nil

// 	// var allPortInfos []portInfo
// 	// for name, data := range interfacesRaw {
// 	// 	var config module.InterfaceConfig
// 	// 	// Decode the specific port data into the module struct
// 	// 	err := mapstructure.Decode(data, &config)
// 	// 	if err != nil {
// 	// 		fmt.Println("Error decoding port data:", err)
// 	// 		continue // Skip ports that fail to decode
// 	// 	}
// 	// 	fmt.Println("config: ", config)
// 	// 	// 4. Map to your custom internal struct
// 	// 	// Note: Verify if module.InterfaceConfig allows map-style access or struct-style
// 	// 	p := portInfo{
// 	// 		name:   name,
// 	// 		number: strings.TrimPrefix(name, "Ethernet"),
// 	// 		// Status: config.Status,                                                 // Usually "connected", "notconnect", or "disabled"
// 	// 		speed: fmt.Sprintf("%v", data.(map[string]interface{})["bandwidth"]), // eAPI often provides bandwidth in bits
// 	// 		// Add other fields from your previous 'inspect' requirements
// 	// 	}
// 	// 	allPortInfos = append(allPortInfos, p)
// 	// }

// 	// var switchportEt1 module.InterfaceConfig
// 	// if ethernet1, ok := res.Result[0]["interfaces"].(map[string]interface{})["Ethernet40"].(map[string]interface{}); ok {
// 	// 	mapstructure.Decode(ethernet1, &switchportEt1)
// 	// }

//		// fmt.Println(switchportEt1)
//		// portInfo := portInfo{
//		// 	name:   switchportEt1["Name"],
//		// 	number: strings.TrimPrefix(switchportEt1["Name"], "Ethernet"),
//		// }
//		// fmt.Println(portInfo.name)
//		// fmt.Println(portInfo.number)
//		// for _, port := range allPortInfos {
//		// 	fmt.Println(port.name)
//		// 	fmt.Println(port.number)
//		// 	fmt.Println(port.speed)
//		// 	fmt.Println("")
//		// }
//		// return nil, nil
//	}
func (a *aristaScanner) getPorts() ([]portInfo, error) {
	err := a.establishAristaConnection()
	if err != nil {
		return nil, err
	}

	// Use RunCommands to get the raw JSON map which is more detailed
	res, err := a.node.RunCommands([]string{"show interfaces status"}, "json")
	if err != nil {
		return nil, err
	}

	interfacesRaw, ok := res.Result[0]["interfaceStatuses"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("loom: interfaceStatuses missing")
	}

	var allPortInfos []portInfo
	for name, data := range interfacesRaw {
		if !strings.HasPrefix(name, "Ethernet") {
			continue
		}

		raw := data.(map[string]interface{})

		// Map eAPI bandwidth to our human labels
		bw := 0
		if val, ok := raw["bandwidth"].(float64); ok {
			bw = int(val)
		}

		p := portInfo{
			name:   name,
			number: strings.TrimPrefix(name, "Ethernet"),
			status: fmt.Sprintf("%v", raw["linkStatus"]), // Arista returns 'connected', 'notconnect', etc.
			speed:  formatBw(bw),
		}
		allPortInfos = append(allPortInfos, p)
	}

	return allPortInfos, nil
}

func formatBw(bw int) string {
	if label, ok := speedMap[bw]; ok {
		return label
	}
	if bw >= 1000000000 {
		return fmt.Sprintf("%dgig", bw/1000000000)
	}
	return fmt.Sprintf("%dmeg", bw/1000000)
}
func (a *aristaScanner) getSystemStats() error {
	if a.conn == nil {
		err := a.establishAristaConnection()
		if err != nil {
			fmt.Println("getSystemStats error: ", err)
			return err
		}
	}
	fmt.Println("getSystemStats not yet implemented")
	return nil
}

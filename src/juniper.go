package main

// import (
// 	"context"
// 	"encoding/xml"
// 	"log"
// 	"time"

// 	"golang.org/x/crypto/ssh"
// 	"nemith.io/netconf"
// 	ncssh "nemith.io/netconf/transport/ssh"
// )

// type juniperAPI struct {
// 	session netconf.Session
// }

// // 1. Define the Vendor-Specific RPC Request
// // This tells Junos to run `show interfaces diagnostics optics`
// type GetOptics struct {
// 	XMLName xml.Name `xml:"get-interface-optics-diagnostics-information"`
// }

// // (Optional) Request for `show chassis hardware` if you need SFP serial numbers instead
// type GetInventory struct {
// 	XMLName xml.Name `xml:"get-chassis-inventory"`
// }

// // 2. Define a generic struct to catch the raw XML payload
// // This prevents you from having to map out the entire Junos XML tree immediately
// type RawReply struct {
// 	Data []byte `xml:",innerxml"`
// }

// func startJuniper(server string, port int, username string, password string) (juniperAPI, error) {
// 	// Configure SSH client
// 	config := &ssh.ClientConfig{
// 		User: "sam",
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password("ubec1924SJB!"),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Don't use in production
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	transport, err := ncssh.Dial(ctx, "tcp", "192.168.0.1:22", config)
// 	if err != nil {
// 		log.Fatalf("failed to connect: %v", err)
// 	}
// 	defer transport.Close()

// 	juniperVar, err := netconf.NewSession(transport)
// 	if err != nil {
// 		log.Fatalf("failed to create session: %v", err)
// 	}
// 	defer session.Close(context.Background())
// 	juniper := juniperAPI{session: juniperVar}
// 	return juniper, nil

// 	// // --- Execute the Operational Command ---
// 	// req := GetOptics{}
// 	// var reply RawReply

// 	// // session.Exec automatically wraps your struct in an <rpc> tag
// 	// // and unmarshals the <rpc-reply> body into your reply struct.
// 	// err = juniper.session.Exec(ctx, req, &reply)
// 	// if err != nil {
// 	// 	log.Fatalf("rpc failed: %v", err)
// 	// }

// 	// // Print the raw XML returned by the router
// 	// fmt.Printf("Transceiver Diagnostics XML:\n%s\n", string(reply.Data))

// 	// // Example of parsing specific fields out of the raw XML:
// 	// type OpticsInfo struct {
// 	// 	XMLName    xml.Name `xml:"interface-information"`
// 	// 	Interfaces []struct {
// 	// 		Name    string `xml:"name"`
// 	// 		RxPower string `xml:"optics-diagnostics>receiver-signal-average-optical-power"`
// 	// 		TxPower string `xml:"optics-diagnostics>laser-output-power"`
// 	// 	} `xml:"physical-interface"`
// 	// }

// 	// var parsedData OpticsInfo
// 	// if err := xml.Unmarshal(reply.Data, &parsedData); err != nil {
// 	// 	log.Fatalf("failed to parse xml: %v", err)
// 	// }

// 	// for _, intf := range parsedData.Interfaces {
// 	// 	// Junos will return empty structs for interfaces without optics, so filter them
// 	// 	if intf.RxPower != "" {
// 	// 		fmt.Printf("Interface: %s | Rx: %s | Tx: %s\n", intf.Name, intf.RxPower, intf.TxPower)
// 	// 	}
// 	// }
// }

// func getDeviceInfo() {}

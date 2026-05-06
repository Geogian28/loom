package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	ac := parseFlags()

	// ac := appConfig{
	// 	address:  "192.168.0.2",
	// 	port:     80,
	// 	username: "loomsvc",
	// 	password: "loomsvc",
	// }

	_, err := getDeviceInfo(ac.address, ac.port, ac.username, ac.password)
	if err != nil {
		log.Fatal(err)
	}

	setupRoutes(ac)

	fmt.Println("loom running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

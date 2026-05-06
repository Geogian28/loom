package main

import (
	"fmt"
	"log"
)

func main() {

	ac := parseFlags()

	_, err := getDeviceInfo(ac.address, ac.port, ac.username, ac.password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("loom running on port 8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

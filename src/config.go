package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type appConfig struct {
	port     int
	username string
	password string
	address  string
}

func parseFlags() appConfig {
	var ac appConfig
	flag.Parse()

	var argumentErrors []string

	val, exists := os.LookupEnv("LOOM_PORT")
	if !exists {
		ac.port = 80
	} else if port, err := strconv.Atoi(val); err != nil {
		argumentErrors = append(argumentErrors, "invalid port: "+val)
	} else {
		ac.port = port
	}
	ac.username = os.Getenv("LOOM_USERNAME")
	ac.password = os.Getenv("LOOM_PASSWORD")
	ac.address = os.Getenv("LOOM_ADDRESS")

	if ac.username == "" {
		argumentErrors = append(argumentErrors, "missing username")
	}
	if ac.address == "" {
		argumentErrors = append(argumentErrors, "missing address")
	}
	if ac.password == "" {
		argumentErrors = append(argumentErrors, "missing password")
	}
	if len(argumentErrors) > 0 {
		fmt.Printf("Issues with arguments:\n  - %s\n", strings.Join(argumentErrors, "\n  - "))
		os.Exit(1)
	}
	return ac
}

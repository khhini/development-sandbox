package main

import "log"

func main() {
	svr, err := InitializeServer()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := svr.Listen("0.0.0.0:8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

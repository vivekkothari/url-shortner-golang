package main

import (
	"fmt"
	"os"
	"url-shortner/server"
)

func main() {
	s, err := server.NewServer(":8080")
	if err != nil {
		fmt.Println("Error starting server")
		os.Exit(1)
	}
	s.Start()
}

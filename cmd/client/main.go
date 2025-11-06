package main

import (
	"log"

	"example.com/user/internal/client"
)

func main() {
	c := client.New()
	if err := c.RunExamples(); err != nil {
		log.Fatalf("Client examples failed: %v", err)
	}
}
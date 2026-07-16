package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Worker iniciado.")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Nenhum evento pendente no MVP.")
	}
}

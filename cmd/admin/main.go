package main

import (
	"example/API_Gateway/pkg/services"
	"log"
)

func main() {
	log.Println("Starting Admin Service...")
	services.StartAdminService()
}

package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Startings")
	scanner, err := CreateInMemoryFileScanner()
	fmt.Println("Generated scanner")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Startings server")
	log.Fatal(RunServer(scanner))
}

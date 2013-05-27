package main

import (
	"log"
	"os"
)

func ExitMsg(msg string) {
	log.Println("Fatal:", msg)
	os.Exit(1)
}



package main

import (
	"github.com/g3co/img-slicer/cmd"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd.Execute()
}

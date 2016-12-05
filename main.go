package main

import (
	"flag"

	"github.com/brodo/go-opc/shader"
)

func main() {
	flag.Parse()

	shader.Start(flag.Args()[0])
}

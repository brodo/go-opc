package main

import (
	"flag"
	"fmt"

	"github.com/brodo/go-opc/shader"
)

func main() {
	target := flag.String("target", "balldachin", "the opc target")
	port := flag.Uint("port", 7890, "the opc port")

	flag.Parse()

	shader.Start(flag.Args()[0], fmt.Sprintf("%s:%d", *target, *port))
}

package main

import (
	"flag"
	"fmt"

	"github.com/tkmktzw/golang-training/ch07/ex06/tempflag"
)

var temp = tempflag.CelsiusFlag("temp", 20.0, "the temprature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

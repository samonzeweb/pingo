package main

import (
	"fmt"
	"os"
	"time"

	"github.com/samonzeweb/pingo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a hostname as argument")
		os.Exit(1)
	}

	t, err := pingo.SimplePing(os.Args[1], pingo.IP, time.Second)
	if err != nil {
		if err == pingo.ErrTimeOut {
			fmt.Printf("Time out : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Error : %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Got a response from %s, in %d ms\n", os.Args[1], t/time.Millisecond)
}

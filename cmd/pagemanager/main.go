package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bokwoon95/pagemanager"
)

func main() {
	flag.Parse()
	cfg, err := pagemanager.DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", cfg)
}

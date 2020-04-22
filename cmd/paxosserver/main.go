package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	ps "github.com/selabhvl/singlepaxos"
)

func main() {
	/* get port of each server from input */
	var (
		port   = flag.Int("port", 8080, "port to listen on")
		saddrs = flag.String("addrs", "", "server addresses separated by ','")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	addrs := strings.Split(*saddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no server addresses provided")
	}
	q := (len(addrs)-1)/2 + 1

	// start a single server
	ps.ServerStart(*port, addrs, q)
}

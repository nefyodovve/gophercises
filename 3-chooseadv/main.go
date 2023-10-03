package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var dataPath *string = flag.String("path", "gopher.json", "Path for data file in json format, default \"gopher.json\"")
	var isLocal *bool = flag.Bool("local", false, "If true, run in command line")
	flag.Parse()
	data, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatal(err)
	}
	m, err := parse(data)
	if err != nil {
		log.Fatal(err)
	}
	if *isLocal {
		local(m)
	} else {
		startServer(m)
	}
}

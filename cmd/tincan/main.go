package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	config := flag.String("config", "", "config directory")

	debug := flag.Bool("debug", false, "print the --config value on stderr")

	flag.Parse()

	if *debug && *config != "" {
		fmt.Fprintf(os.Stderr, "config: %s\n", *config)
	}

	fmt.Println("tincan: hello")
	fmt.Println(*config)
	fmt.Println(*debug)
}

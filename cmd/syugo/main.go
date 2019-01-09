package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nissy/syugo"
)

var filename = flag.String("c", "syugo.toml", "set configuration file.")

type Config struct {
	Collects []*syugo.Collect
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func run() (err error) {
	flag.Parse()

	cfg := &Config{}
	if _, err := toml.DecodeFile(*filename, cfg); err != nil {
		return err
	}

	syu, err := syugo.NewSyugo(cfg.Collects)
	if err != nil {
		return err
	}

	return syu.Run()
}

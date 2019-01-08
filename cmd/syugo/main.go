package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nissy/syugo"
	"gopkg.in/go-playground/validator.v9"
)

var filename = flag.String("c", "syugo.toml", "set configuration file.")

type Config struct {
	Collects syugo.Collects `toml:"collect" validate:"required,dive,required"`
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
	if err := validator.New().Struct(cfg); err != nil {
		return err
	}

	return cfg.Collects.Run()
}

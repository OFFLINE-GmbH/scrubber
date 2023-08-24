package main

import (
	"flag"
	"log"
	"os"
	"scrubber"

	"github.com/BurntSushi/toml"
)

func main() {
	cfgFile := flag.String("config", "scrubber.config.toml", "Path to the config file")
	pretend := flag.Bool("pretend", false, "Print out actions that would be executed but do nothing")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Printf("Loading configuration file %s", *cfgFile)

	var conf scrubber.TomlConfig
	if _, err := toml.DecodeFile(*cfgFile, &conf); err != nil {
		logger.Fatalf("Could not decode config file: %s", err)
		return
	}

	logger.Println("Beginning to scrub...")

	fs := scrubber.OSFilesystem{}

	s := scrubber.New(&conf, fs, logger, *pretend)
	err := s.Scrub()
	if err != nil {
		logger.Fatalf("error while scrubbing files: %s", err)
	}
}

package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Read a configuration from fname or die

func ReadGlobalConfig(fname string) *SCORSHmaster {

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Error while reading file: ", err)
	}

	var cfg *SCORSHmaster

	cfg = new(SCORSHmaster)

	// Unmarshal the YAML configuration file into a SCORSHcfg structure
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatal("Error while reading configuration: ", err)
	}

	fmt.Printf("%s", cfg)
	
	if cfg.Logfile != "" {
		f, err := os.OpenFile(cfg.Logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			log.Fatal("Error opening logfile: ", cfg.Logfile, err)
		} else {
			log.SetOutput(io.Writer(f))
		}
	}

	if cfg.LogPrefix != "" {
		log.SetPrefix(cfg.LogPrefix+ " ")
	}

	// If the user has not set a spooldir, crash loudly
	if cfg.Spooldir == "" {
		log.Fatal("No spooldir defined in ", fname, ". Exiting\n")
	}

	// Check if the user has set a custom logprefix

	// Check if the user wants to redirect the logs to a file

	// If we got so far, then there is some sort of config in cfg
	log.Printf("----- Starting SCORSH -----\n")
	log.Printf("Successfully read config from %s\n", fname)

	return cfg

}


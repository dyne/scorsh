package main

import (
	"bytes"
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
	
	// If the user has not set a spooldir, crash loudly
	if cfg.Spooldir == "" {
		log.Fatal("No spooldir defined in ", fname, ". Exiting\n")
	}

	// Check if the user has set a custom logprefix

	// Check if the user wants to redirect the logs to a file
	if cfg.Logfile != "" {
		log.Printf("Opening log file: %s\n", cfg.Logfile)
		f, err := os.OpenFile(cfg.Logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			log.SetOutput(io.Writer(f))
		} else {
			log.Fatal("Error opening logfile: ", cfg.Logfile, err)
		}
	}

	if cfg.LogPrefix != "" {
		log.SetPrefix(cfg.LogPrefix)
	}

	// If we got so far, then there is some sort of config in cfg
	log.Printf("Successfully read config from %s\n", fname)

	return cfg

}




func (cfg *SCORSHmaster) String() string {

	var buff bytes.Buffer

	buff.WriteString("spooldir: ")
	buff.WriteString(cfg.Spooldir)
	buff.WriteString("\nlogfile: ")
	buff.WriteString(cfg.Logfile)
	buff.WriteString("\nlogprefix: ")
	buff.WriteString(cfg.LogPrefix)
	buff.WriteString("\nWorkers: \n")

	for _, w := range cfg.Workers {
		buff.WriteString("---\n  name: ")
		buff.WriteString(w.Name)
		buff.WriteString("\n  repos: ")
		for _, r := range w.Repos {
			buff.WriteString("\n    ")
			buff.WriteString(r)
		}
		buff.WriteString("\n  folder: ")
		buff.WriteString(w.Folder)
		buff.WriteString("\n  logfile: ")
		buff.WriteString(w.Logfile)
		buff.WriteString("\n  tagfile: ")
		buff.WriteString(w.Tagfile)
		buff.WriteString("\n  keyrings: ")
		for _, k := range w.Keyrings {
			buff.WriteString("\n    ")
			buff.WriteString(k)
		}
		buff.WriteString("\n...\n")

	}

	return buff.String()
}

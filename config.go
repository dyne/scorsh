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

type SCORSHWorker_cfg struct {
	Name     string   `yaml:"w_name"`
	Repos    []string `yaml:"w_repos"`
	Folder   string   `yaml:"w_folder"`
	Logfile  string   `yaml:"w_logfile"`
	Tagfile  string   `yaml:"w_tagfile"`
	Keyrings []string `yaml:"w_keyrings"`
}

type SCORSHcfg struct {
	Spooldir  string             `yaml:"s_spooldir"`
	Logfile   string             `yaml:"s_logfile"`
	LogPrefix string             `yaml:"s_logprefix"`
	Workers   []SCORSHWorker_cfg `yaml:"s_workers"`
}

// Read a configuration from fname or die

func ReadGlobalConfig(fname string) *SCORSHcfg {

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Error while reading file: ", err)
	}

	var cfg *SCORSHcfg
	cfg = new(SCORSHcfg)

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
	if cfg.LogPrefix != "" {
		log.SetPrefix(cfg.LogPrefix)
	}

	// Check if the user wants to redirect the logs to a file
	if cfg.Logfile != "" {
		f, err := os.Open(cfg.Logfile)
		if err != nil {
			log.SetOutput(io.Writer(f))
		} else {
			log.Printf("Error opening logfile: \n", err)
		}
	}

	// If we got so far, then there is some sort of config in cfg
	log.Printf("Successfully read config from %s\n", fname)

	return cfg

}

func (cfg *SCORSHcfg) String() string {

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

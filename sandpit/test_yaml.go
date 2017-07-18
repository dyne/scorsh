package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
	"strings"
)

type STag struct {
	S_tag  string
	S_args []string
}

type SCmd struct {
	S_cmd  string
	S_hash string
}

type STagConfig struct {
	S_tag      string
	S_commands []SCmd
}

type SCORSHmsg struct {
	S_msg []STag
}

type SCORSHcfg struct {
	S_cfg []STagConfig
}

var msg_str = `
s_msg:
   [
    { s_tag: BUILD,
      s_args: [ suites/jessie,  suites/ascii]
    },
    {
     s_tag: REMOVE,
     s_args: [file1]
    },
    {
     s_tag: CUSTOM,
     s_args: [first, second, third]
    }
]
`

var other_msg = `
this is my comment...

---
s_msg: [
        {s_tag: "BUILD", s_args: [suites/jessie, suites/ascii]},
        {s_tag: "REMOVE", s_args: [file1]},
        {s_tag: "CUSTOM", s_args: [first, second, third]}
       ]
`

var cfg_str = `
some stuff
---
s_cfg:
  - s_tag: BUILD
    s_commands:
     - s_cmd: file:///bin/ls
       s_hash: 12345
     - s_cmd: file:///home/katolaz/script.sh
       s_hash: abc123df
     - s_cmd: http://myserver.org/build.php?name=\1
       s_hash: 
  - s_tag: REMOVE
    s_commands:
     - s_cmd: file:///bin/rm 
  - s_tag: CUSTOM
    s_commands: [
                 {s_cmd: "file:///home/user/script/sh", s_hash: "1234567890abcdef"}, 
                 {s_cmd: "http://my.server.net/submit.php", s_hash: "0987654321abce"}
                ]
...
`

func main() {

	var c SCORSHmsg

	var conf SCORSHcfg

	sep := "\n---\n"

	//log.Printf("%s\n", test_str)

	scorsh_idx := strings.Index(other_msg, sep)
	if scorsh_idx >= 0 {

		err := yaml.Unmarshal([]byte(other_msg[scorsh_idx:]), &c)

		if err != nil {
			log.Fatal("error: ", err)
		}

		for _, item := range c.S_msg {
			fmt.Printf("Record: \n")
			fmt.Printf("  s_tag: %s\n", item.S_tag)
			fmt.Printf("  s_args:\n")

			for _, a := range item.S_args {
				fmt.Printf("    %s\n", a)
			}
		}
	}

	fmt.Println("----------------------------")

	scorsh_idx = strings.Index(cfg_str, sep)
	if scorsh_idx >= 0 {
		
		err := yaml.Unmarshal([]byte(cfg_str[scorsh_idx:]), &conf)
		if err != nil {
			log.Fatal("error: ", err)
		}

		for _, cfg_item := range conf.S_cfg {
			fmt.Printf("Config record:\n")
			fmt.Printf("  s_tag: %s\n", cfg_item.S_tag)
			fmt.Printf("  s_commands:\n")

			for _, c := range cfg_item.S_commands {
				fmt.Printf("    s_cmd: %s\n", c.S_cmd)
				fmt.Printf("    s_hash: %s\n", c.S_hash)
				fmt.Println("    ---")
			}
			fmt.Println("-+-+-")

		}
	}
}

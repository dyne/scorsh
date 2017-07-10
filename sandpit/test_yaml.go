package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
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
   - s_tag: BUILD
     s_args: 
      -   suites/jessie 
      -   suites/ascii
   - s_tag: REMOVE
     s_args: 
      - file1
   - s_tag: CUSTOM
     s_args: [first, second, third]
`

var other_msg = `
s_msg: [
        {s_tag: "BUILD", s_args: [suites/jessie, suites/ascii]},
        {s_tag: "REMOVE", s_args: [file1]},
        {s_tag: "CUSTOM", s_args: [first, second, third]}
       ]
`


var cfg_str = `
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
`

func main() {

	var c SCORSHmsg

	var conf SCORSHcfg

	//log.Printf("%s\n", test_str)

	err := yaml.Unmarshal([]byte(other_msg), &c)
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

	fmt.Println("----------------------------")

	err = yaml.Unmarshal([]byte(cfg_str), &conf)
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

package main

import(
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
	"strings"
)


var worker_cfg = `
---
w_tags:
    [
     {
      t_name: "BUILD",
      t_keyrings: ["build_keyring.asc", "general_keyring.asc"],
      t_commands: [
                    {
                     c_url: "file:///home/user/bin/script.sh $1 $2",
                     c_hash: "12da324fb76s924acbce"
                    },
                    {
                     c_url: "http://my.server.net/call.pl?branch=$1"
                    }
                   ]
     },
     {
      t_name: "PUBLISH",
      t_keyrings: ["web_developers.asc"],
      t_commands: [
                    {
                     c_url: "file:///usr/local/bin/publish.py $repo $branch",
                     c_hash: "3234567898765432345678"
                    }
                   ]
      }
   ]
...

`


func main(){
	
	var w *SCORSHworker
	w = new(SCORSHworker)


	sep := "\n---\n"
	
	idx := strings.Index(worker_cfg, sep)

	err := yaml.Unmarshal([]byte(worker_cfg[idx:]), w)

		
	if err != nil{
		log.Printf("Error parsing message: %s", err)
	}

	fmt.Printf("%s\n", w)
	
}

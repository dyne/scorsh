package main

import(
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
)


var orig_msg= `
---
m_id: 123456
m_repo: master
m_branch: test_branch
m_oldrev: a1b2c3d4e5f6
m_newrev: 9a8b7c6d5e4f
...

`


func main(){
	
	var msg *SCORSHmsg
	msg = new(SCORSHmsg)


	fname := "spool/test_2"

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("Unable to open file: %s\n", fname)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(data), msg)
	if err != nil{
		log.Printf("Error parsing message: %s", err)
	}

	fmt.Printf("%s\n", msg)
	
}

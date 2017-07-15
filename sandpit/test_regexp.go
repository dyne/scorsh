package main

import(
	"regexp"
	"log"
)


func main (){

	pattern := ".*"
	str := "my_string"

	matched, err := regexp.MatchString(pattern, str)

	if err != nil {
		log.Fatal("Error matching string: ", err)
	}

	if matched {
		log.Printf("Yes! '%s' matched '%s'\n", str, pattern)
	} else {
		log.Printf("Bad luck!\n")
	}

}

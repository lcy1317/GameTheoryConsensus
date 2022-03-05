package main

import (
	"colorout"
	"log"
)

func main() {

	if err := configInitial(); err != nil {
		log.Fatalln(colorout.Red("ReadInConfig error:" + err.Error()))
	}
}

package main

import (
	"log"
	"net"
)


func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}
	// }

}

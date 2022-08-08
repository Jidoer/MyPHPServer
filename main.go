package main

import "log"

var Serv *Server

func main() {
	Serv := NewServer()
	log.Println(Serv.Config)
	Serv.Start()
}


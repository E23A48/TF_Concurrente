package consenso

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"../entities"
)

var Prediction int = -1

const (
	Cnum = iota
	Opa  = 1
	Opb  = 0
)

type Tmsg struct {
	Code    int
	Addr    string
	Op      int
	Pacient entities.Pacient
}

var LocalAddr string = "192.168.0.7:8100"

var Addrs = []string{
	"192.168.0.7:8200",
}

var chInfo chan map[string]int

func GoSV() {

	chInfo = make(chan map[string]int)

	go func() { chInfo <- map[string]int{} }()

	go server()

}

func server() {
	if ln, err := net.Listen("tcp", LocalAddr); err != nil {
		log.Panicln("Can't start listener on", LocalAddr)
	} else {
		defer ln.Close()
		fmt.Println("Listeing on", LocalAddr)
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Println("Can't accept", conn.RemoteAddr())
			} else {
				go handle(conn)
			}
		}
	}
}
func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg Tmsg
	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case Cnum:
			concensus(conn, msg)
		}
	}
}

func concensus(conn net.Conn, msg Tmsg) {
	info := <-chInfo
	info[msg.Addr] = msg.Op
	fmt.Println(info)
	if len(info) == len(Addrs) {
		ca, cb := 0, 0
		for _, op := range info {
			if op == Opa {
				ca++
			} else {
				cb++
			}
		}
		if ca > cb {
			Prediction = 1
		} else {
			Prediction = 0
		}
		info = map[string]int{}
	}
	go func() { chInfo <- info }()
}

func Send(remoteAddr string, msg Tmsg) {

	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr, err)
	} else {
		defer conn.Close()
		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}

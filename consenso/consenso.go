package consenso

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const localAddr = "localhost:8004" // su propia IP aquí
const (
	cnum = iota // iota genera valores en secuencia y se reinicia en cada bloque const
	opa  = 1
	opb  = 0
)

type Tmsg struct {
	Code int
	Addr string
	Op   int
}

// Las IP de los demás participantes acá, todos deberían usar el puerto 8000
var addrs = []string{"localhost:8001", "localhost:8002", "localhost:8003"}

var chInfo chan map[string]int

func GoServer() {

	chInfo = make(chan map[string]int)

	go func() { chInfo <- map[string]int{} }()

	go server()

}

func server() {
	if ln, err := net.Listen("tcp", localAddr); err != nil {
		log.Panicln("Can't start listener on", localAddr)
	} else {
		defer ln.Close()
		fmt.Println("Listeing on", localAddr)
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
		case cnum:
			concensus(conn, msg)
		}
	}
}

func concensus(conn net.Conn, msg Tmsg) {
	info := <-chInfo
	info[msg.Addr] = msg.Op
	if len(info) == len(addrs) {
		ca, cb := 0, 0
		for _, op := range info {
			if op == opa {
				ca++
			} else {
				cb++
			}
		}
		if ca > cb {
			fmt.Println("GO A!")
		} else {
			fmt.Println("GO B!")
		}
		info = map[string]int{}
	}
	go func() { chInfo <- info }()
}

func send(remoteAddr string, msg Tmsg) {

	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr, err)
	} else {
		defer conn.Close()
		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}

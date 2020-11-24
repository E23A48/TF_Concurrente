package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	"./entities"
	rforest "./random_forest"
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

var Pacient entities.Pacient = entities.Pacient{
	0,
	0,
	0,
	0,
	0,
	0.0,
	0.0,
	0,
}

var Msg Tmsg

var MakePrediction bool = false

var LocalAddr string = "192.168.0.7:8200"

var Addrs = []string{"192.168.0.7:8100"}

var chInfo chan map[string]int

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

	if err := dec.Decode(&Msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(Msg)
		MakePrediction = true
		switch Msg.Code {
		case Cnum:
			concensus(conn, Msg)
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

func main() {

	chInfo = make(chan map[string]int)

	go func() { chInfo <- map[string]int{} }()

	go server()

	for {

		if MakePrediction {

			url := "https://raw.githubusercontent.com/Malvodio/TF_Funda_Videojuegos/master/diabetes.csv"
			df, _ := rforest.LoadCSV(url)

			feature := []string{
				strconv.Itoa(Pacient.Pregnancies),
				strconv.Itoa(Pacient.Glucose),
				strconv.Itoa(Pacient.BloodPressure),
				strconv.Itoa(Pacient.SkinThickness),
				strconv.Itoa(Pacient.Insulin),
				strconv.FormatFloat(Pacient.BMI, 'f', 6, 64),
				strconv.FormatFloat(Pacient.DiabetesPedigreeFunction, 'f', 6, 64),
				strconv.Itoa(Pacient.Age),
			}

			_, _, X_train, y_train, _, _ := rforest.TrainTestSplit(df, 0.15)

			forest := rforest.BuildForest(X_train, y_train, rand.Intn(30), 500, len(X_train[0]))

			X := make([]interface{}, 0)

			for _, x := range feature {
				X = append(X, x)
			}

			var result string = forest.Predicate(X)
			result_int, _ := strconv.Atoi(result)

			fmt.Println("Resultado de prediccion: ", result_int)

			msg := Tmsg{Cnum, LocalAddr, result_int, Pacient}
			for _, addr := range Addrs {
				Send(addr, msg)
			}

			MakePrediction = false

		}
	}

}

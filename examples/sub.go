package main

import (
	"log"
	"runtime"

	server "github.com/hemerajs/go-hemera/server"
	nats "github.com/nats-io/go-nats"
)

type Result struct {
	Result float64 `json:"result"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	hemera := server.Hemera{Conn: nc}

	pattern := server.Pattern{"topic": "math", "cmd": "add"}
	hemera.Add(pattern, func(req server.Pattern, reply server.Reply) {
		r := req["a"].(float64) + req["b"].(float64)
		reply(Result{Result: r})
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on \n")

	runtime.Goexit()
}

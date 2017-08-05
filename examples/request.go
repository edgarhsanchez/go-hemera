// +build ignore

package main

import (
	"fmt"
	"log"
	"runtime"

	server "github.com/hemerajs/go-hemera"
	nats "github.com/nats-io/go-nats"
)

type MathPattern struct {
	Topic string `json:"topic"`
	Cmd string `json:"cmd"`
}

type Delegate struct {
		Query string `json:"query"`
}

type Meta struct {
		Token string `json:"token"`
}

type RequestPattern struct {
	Topic string `json:"topic" mapstructure:"topic"`
	Cmd string `json:"cmd" mapstructure:"cmd"`
	A int `json:"a" mapstructure:"a"`
	B int `json:"b" mapstructure:"b"`
	Meta_ Meta `json:"meta"`
	Delegate_ Delegate `json:"meta"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	hemera, _ := server.NewHemera(nc)

	pattern := MathPattern{ Topic: "math", Cmd: "add" }

	hemera.Add(pattern, func(context server.Context, req *RequestPattern, reply server.Reply) {
		fmt.Printf("Request: %+v\n", context)
		reply.Send(req.A + req.B)
	})

	requestPattern := RequestPattern{
		Topic: "math",
		Cmd: "add",
		A: 1,
		B: 2,
		Delegate_: Delegate{ Query : "DEF" },
		Meta_: Meta{ Token : "ABC" },
	}
	
	hemera.Act(requestPattern, func(resp server.ClientResult) {
		fmt.Printf("Response: %+v\n", resp)
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on \n")

	runtime.Goexit()
}

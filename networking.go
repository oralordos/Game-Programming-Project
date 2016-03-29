package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/nu7hatch/gouuid"
)

func StartNetworkListener() {
	ln, err := net.Listen("tcp", ":10328")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		NewNetworkFrontend(conn)
	}
}

type NetworkFrontend struct {
	id      *uuid.UUID
	conn    net.Conn
	eventCh chan events.Event
	close   chan struct{}
}

func NewNetworkFrontend(conn net.Conn) *NetworkFrontend {
	n := NetworkFrontend{
		conn:    conn,
		eventCh: make(chan events.Event),
		close:   make(chan struct{}),
	}
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalln(err)
	}
	n.id = id
	events.AddListener(n.eventCh, events.DirFront, 0)
	go n.readloop()
	go n.mainloop()
	return &n
}

func (n *NetworkFrontend) readloop() {
	for {
		data := make(map[string]interface{})
		err := json.NewDecoder(n.conn).Decode(&data)
		if err != nil {
			log.Println(err)
			break
		}
		ev := events.DecodeJSON(data)
		if ev == nil {
			log.Println("Unable to parse event")
		}
		log.Printf("%T\n", ev)
		// events.SendEvent(ev)
	}
}

func (n *NetworkFrontend) mainloop() {
loop:
	for {
		select {
		case ev := <-n.eventCh:
			n.handleEvent(ev)
		case _, ok := <-n.close:
			if !ok {
				break loop
			}
		}
	}
	events.RemoveListener(n.eventCh, events.DirFront, 0)
}

func (n *NetworkFrontend) Destroy() {
	close(n.close)
}

func (n *NetworkFrontend) handleEvent(ev events.Event) {
	err := json.NewEncoder(n.conn).Encode(ev)
	if err != nil {
		log.Println(err)
	}
}

type NetworkBackend struct {
	conn    net.Conn
	eventCh chan events.Event
	close   chan struct{}
}

func NewNetworkBackend(address string) *NetworkBackend {
	n := NetworkBackend{
		eventCh: make(chan events.Event),
		close:   make(chan struct{}),
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	n.conn = conn
	events.AddListener(n.eventCh, events.DirSystem, 0)
	go n.readloop()
	go n.mainloop()
	return &n
}

func (n *NetworkBackend) readloop() {
	for {
		data := make(map[string]interface{})
		err := json.NewDecoder(n.conn).Decode(&data)
		if err != nil {
			log.Println(err)
			break
		}
		ev := events.DecodeJSON(data)
		if ev == nil {
			log.Println("Unable to parse event")
		}
		log.Printf("%T\n", ev)
		// events.SendEvent(ev)
	}
}

func (n *NetworkBackend) mainloop() {
loop:
	for {
		select {
		case ev := <-n.eventCh:
			n.handleEvent(ev)
		case _, ok := <-n.close:
			if !ok {
				break loop
			}
		}
	}
	events.RemoveListener(n.eventCh, events.DirSystem, 0)
}

func (n *NetworkBackend) handleEvent(ev events.Event) {
	err := json.NewEncoder(n.conn).Encode(ev)
	if err != nil {
		log.Println(err)
	}
}

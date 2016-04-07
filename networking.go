package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/Oralordos/Game-Programming-Project/events"
	"github.com/nu7hatch/gouuid"
)

const port = ":10328"

type Network struct {
	conn    net.Conn
	decoder *json.Decoder
	dir     int
	eventCh chan events.Event
	close   chan struct{}
}

func (n *Network) readloop() {
	for {
		data := make(map[string]interface{})
		err := n.decoder.Decode(&data)
		if err != nil {
			log.Println(err)
			break
		}
		ev := events.DecodeJSON(data)
		if ev == nil {
			log.Println("Unable to parse event")
			continue
		}
		ev.SetDuplicate(true)
		events.SendEvent(ev)
	}
}

func (n *Network) mainloop() {
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
	events.RemoveListener(n.eventCh, n.dir, 0)
	n.conn.Close()
}

func (n *Network) Destroy() {
	close(n.close)
}

func (n *Network) handleEvent(ev events.Event) {
	if ev.HasDuplicate() {
		return
	}
	err := json.NewEncoder(n.conn).Encode(ev)
	if err != nil {
		log.Println(err)
	}
}

func StartNetworkListener() {
	ln, err := net.Listen("tcp", port)
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
	Network
	id *uuid.UUID
}

func NewNetworkFrontend(conn net.Conn) *NetworkFrontend {
	n := NetworkFrontend{
		Network: Network{
			conn:    conn,
			dir:     events.DirFront,
			decoder: json.NewDecoder(conn),
			eventCh: make(chan events.Event),
			close:   make(chan struct{}),
		},
	}
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalln(err)
	}
	n.id = id
	events.AddListener(n.eventCh, events.DirFront, 0)
	go n.readloop()
	go n.mainloop()
	events.SendEvent(events.ReloadLevel{})
	return &n
}

type NetworkBackend struct {
	Network
}

func NewNetworkBackend(address string) *NetworkBackend {
	n := NetworkBackend{
		Network: Network{
			dir:     events.DirSystem,
			eventCh: make(chan events.Event),
			close:   make(chan struct{}),
		},
	}
	conn, err := net.Dial("tcp", address+port)
	if err != nil {
		log.Fatalln(err)
	}
	n.conn = conn
	n.decoder = json.NewDecoder(n.conn)
	events.AddListener(n.eventCh, events.DirSystem, 0)
	go n.readloop()
	go n.mainloop()
	return &n
}

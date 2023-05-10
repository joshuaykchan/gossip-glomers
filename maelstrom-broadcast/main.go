package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Server struct {
	n *maelstrom.Node

	rMu     sync.Mutex
	readset []int

	nMu       sync.Mutex
	neighbors []string
}

func main() {
	n := maelstrom.NewNode()

	s := Server{
		n:         n,
		readset:   []int{},
		neighbors: []string{},
	}

	n.Handle("broadcast", s.broadcastHandler)
	n.Handle("read", s.readHandler)
	n.Handle("topology", s.topologyHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) broadcastHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.rMu.Lock()
	s.readset = append(s.readset, int(body["message"].(float64)))
	s.rMu.Unlock()

	return s.n.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

func (s *Server) readHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	return s.n.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": s.readset,
	})
}

func (s *Server) topologyHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	type topologyStruct struct {
		Topology map[string][]string `json:"topology"`
	}
	var topo topologyStruct
	if err := json.Unmarshal(msg.Body, &topo); err != nil {
		return err
	}

	s.nMu.Lock()
	s.neighbors = topo.Topology[s.n.ID()]
	s.nMu.Unlock()

	body["type"] = "topology_ok"

	return s.n.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}

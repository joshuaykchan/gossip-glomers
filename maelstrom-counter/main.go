package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Server struct {
	n *maelstrom.Node

	kvMu sync.Mutex
	kv   *maelstrom.KV
}

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)
	s := Server{
		n:  n,
		kv: kv,
	}

	n.Handle("init", s.initHandler)
	n.Handle("add", s.addHandler)
	n.Handle("read", s.readHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) initHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	ctx := context.Background()
	s.kvMu.Lock()
	if err := s.kv.Write(ctx, "count", 0); err != nil {
		return err
	}
	s.kvMu.Unlock()

	return s.n.Reply(msg, map[string]any{
		"type":        "init_ok",
		"in_reply_to": body["msg_id"],
	})
}

func (s *Server) addHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	delta := int(body["delta"].(float64))

	s.kvMu.Lock()
	for {
		ctx, cancel := context.WithCancel(context.Background())
		prev, err := s.kv.ReadInt(ctx, "count")
		if err != nil {
			prev = 0
		}
		if err := s.kv.CompareAndSwap(ctx, "count", prev, prev+delta, true); err == nil {
			cancel()
			break
		}
		cancel()
	}
	s.kvMu.Unlock()

	return s.n.Reply(msg, map[string]any{
		"type": "add_ok",
	})
}

func (s *Server) readHandler(msg maelstrom.Message) error {
	count := 0
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.kvMu.Lock()
	for {
		ctx, cancel := context.WithCancel(context.Background())
		if kvVal, err := s.kv.ReadInt(ctx, "count"); err == nil {
			count = kvVal
			cancel()
			break
		}
		cancel()
	}
	s.kvMu.Unlock()

	return s.n.Reply(msg, map[string]any{
		"type":  "read_ok",
		"value": count,
	})
}

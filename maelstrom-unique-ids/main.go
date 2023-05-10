package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maelstrom-unique-ids/counter"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Server struct {
	n *maelstrom.Node
	c *counter.Counter
}

func (s *Server) genId(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	uid := fmt.Sprintf("%d-%s-%d", time.Now().Unix(), s.n.ID(), s.c.Incr())

	body["type"] = "generate_ok"
	body["id"] = uid

	return s.n.Reply(msg, body)
}

func main() {
	n := maelstrom.NewNode()
	c := counter.NewCounter()

	s := Server{
		n: n,
		c: c,
	}

	n.Handle("generate", s.genId)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

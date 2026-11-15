package main

import (
	"sync"
	"treenity/protocol"
)

type Topic struct {
	mu       sync.RWMutex
	name     string
	storage  []Message
	filter   PrefixFilter
	messages chan protocol.ProduceReq
	requests chan protocol.Request
}

func (t *Topic) run() {
	for {
		req := <-t.requests
		switch r := req.(type) {
		case protocol.ProduceReq:
			t.mu.Lock()
			t.storage = append(t.storage,
				Message{key: r.Key, value: r.Value})
			t.mu.Unlock()
			consumers := t.filter.Match(r.Key)
			for _, c := range consumers {
				select {
				case c.notify <- struct{}{}:
				default:
				}
			}
		case protocol.SubscribeReq:
		case protocol.DisconnectReq:
		}
	}
}

func CreateTopic(request protocol.CreateReq, topics map[string]*Topic) error {
	if _, exists := topics[request.Topic]; exists {
		Reply(request.PID, 0, "error: topic already exists")
	}
	t := &Topic{
		name:     request.Topic,
		filter:   PrefixFilter{root: newNode()},
		messages: make(chan protocol.ProduceReq)}
	topics[request.Topic] = t
	go t.run()
	if err := Reply(request.PID, 0, "topic created"); err != nil {
		return err
	}
	return nil
}

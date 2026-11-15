package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
	"treenity/protocol"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

const ServerFifoPath = "/tmp/treenity.server."
const ReplyFifoPath = "/tmp/treenity.reply."

type Server struct {
	topics  map[string]*Topic
	clients *Hashtable
	IPC     string
}

type Message struct {
	key   string
	value string
}

type Consumer struct {
	name   string
	prefix string
	fifo   string
	notify chan struct{}
}

func Reply(pid uint32, status int, message string) error {
	clientPath := fmt.Sprintf("%s%d", ReplyFifoPath, pid)
	f, err := os.OpenFile(clientPath, os.O_WRONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return err
	}
	buff := make([]byte, 5+len(message))
	buff[0] = byte(status)
	binary.LittleEndian.PutUint32(buff[1:5], uint32(len(message)))
	copy(buff[5:], []byte(message))
	f.Write(buff)
	f.Close()
	return nil
}

func waitRequest(ch chan protocol.Request, f io.Reader) {
	for {
		req, err := protocol.ReadRequest(f)
		if err != nil {
			continue
		}
		ch <- req
	}
}

func handleSubscribe(r protocol.SubscribeReq, s *Server) {
	if _, exists := s.topics[r.Topic]; !exists {
		Reply(r.PID,
			int(protocol.StatusClient),
			"topic: not found")
		return
	}
	clientIPC := fmt.Sprintf("%s.%s", s.IPC, r.Name)
	client := s.clients.Get(r.Name)
	if client == nil {
		client = &ClientMeta{
			Name:      r.Name,
			Topic:     r.Topic,
			Prefix:    r.Prefix,
			IPC:       clientIPC,
			Connected: true,
			LastSeen:  time.Now(),
		}
		if r.HassOffset {
			client.Offset = r.Offset
		}
		s.clients.Insert(client)
	} else {
		if client.Connected {
			Reply(r.PID,
				int(protocol.StatusClient),
				"client already connected")
			return
		} else {
			offset := client.Offset
			if r.HassOffset {
				offset = r.Offset
			} else {
				if r.Topic != client.Topic {
					offset = 0
				}
			}
			client.Prefix = r.Prefix
			client.Topic = r.Topic
			client.Offset = offset
			client.Connected = true
			client.LastSeen = time.Now()
		}
	}
}

func run() error {
	server := Server{topics: make(map[string]*Topic)}
	server.IPC = fmt.Sprintf(
		"%s%d", ServerFifoPath, uint32(os.Getpid()))
	os.Remove(server.IPC)
	defer os.Remove(server.IPC)
	syscall.Mkfifo(server.IPC, 0666)
	fmt.Println(server.IPC)
	sigCh := make(chan os.Signal, 1)
	reqCh := make(chan protocol.Request, 1)
	if f, err := os.OpenFile(server.IPC, os.O_RDWR, 0); err != nil {
		return err
	} else {
		defer f.Close()
		go waitRequest(reqCh, f)
	}
loop:
	for {
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		select {
		case req := <-reqCh:
			switch r := req.(type) {
			case protocol.CreateReq:
				if err := CreateTopic(r, server.topics); err != nil {
					fmt.Fprintln(os.Stderr, "create:", err)
				}
			case protocol.ProduceReq:
				t, exists := server.topics[r.Topic]
				if !exists {
					fmt.Fprintln(os.Stderr, "produce: topic not found")
					break
				}
				t.messages <- r
			case protocol.SubscribeReq:
				handleSubscribe(r, &server)
			default:
			}
		case <-sigCh:
			fmt.Printf("\nserver stopped safely\n")
			break loop
		default:
		}
	}
	return nil
}

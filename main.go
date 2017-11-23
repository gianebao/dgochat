package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gianebao/disgo"
	"github.com/gianebao/disgo/log"
)

var replyChannel = map[string]string{}

func Reader(m disgo.Message) string {

	str := strings.ToLower(strings.TrimRight(m.Content, "\r\n"))
	switch str {

	case "hi":
		return "Hello\r\n"

	case "exit":
		if _, keyExists := replyChannel[m.Worker.ID]; keyExists {
			m.Worker.Swarm.Workers[replyChannel[m.Worker.ID]].WriteString("Worker [" + m.Worker.ID + "] has disconnected.\r\n")
			delete(replyChannel, replyChannel[m.Worker.ID])
			delete(replyChannel, m.Worker.ID)
		}

		m.Worker.Die()

	case "list":
		l := []string{}

		for i := range m.Worker.Swarm.Workers {
			l = append(l, i)
		}
		return strings.Join(l, "\r\n") + "\r\n"

	case "": // fix slice error. not elegant :)
		return ""

	default:
		if len(str) > 2 {
			return ""
		}

		switch str[0:2] {
		case "/s": // starting a conversation
			if i := strings.Index(str, " "); -1 != i {
				id := str[2:i]
				if w, keyExists := m.Worker.Swarm.Workers[id]; keyExists {
					replyChannel[id] = m.Worker.ID
					replyChannel[m.Worker.ID] = id
					w.WriteString(m.Worker.ID + "> " + m.Content[i+1:len(m.Content)] + "/r to reply.\r\n")
					return ""
				}

				return "Worker [" + id + "] does not exist.\r\n"
			}

		case "/r": // replying to a conversation
			if id, rExists := replyChannel[m.Worker.ID]; rExists {
				if w, keyExists := m.Worker.Swarm.Workers[id]; keyExists {
					w.WriteString(m.Worker.ID + "> " + m.Content[3:len(m.Content)])
					return ""
				}

				delete(replyChannel, id)
				delete(replyChannel, m.Worker.ID)
				return "Cannot Send message to disconnected Worker [" + id + "].\r\n"
			}

			return "No conversation was set. Use /s<user> <message> to start a conversation.\r\n"
		}
	}

	return "Unkown command!\r\n"
}

func main() {
	var (
		port     = flag.Int("port", 60217, "listening port")
		portStr  = strconv.Itoa(*port)
		swarm    *disgo.Swarm
		listener net.Listener
		conn     net.Conn
		err      error
	)

	if listener, err = net.Listen("tcp", "0.0.0.0:"+portStr); err != nil {
		fmt.Printf("Failed to listen to port [:%s] with error [%v]. Exit!\n", portStr, err)
		os.Exit(1)
	}

	fmt.Printf("Server now listening to [:%d]. Waiting for incoming connections.\n", *port)

	logchan := log.NewChannel()

	go func(l *log.Channel) {
		var msg string
		for {
			select {
			case msg = <-l.Fatal:
				fmt.Println("[FATAL] ", msg)
				return

			case msg = <-l.Info:
				fmt.Println("[INFO] ", msg)

			case msg = <-l.Warning:
				fmt.Println("[WARNING] ", msg)

			case msg = <-l.Error:
				fmt.Println("[ERROR] ", msg)

			case msg = <-l.Message:
				fmt.Println("[MESSAGE] ", msg)
			}
		}
	}(logchan)

	swarm = disgo.NewSwarm(logchan).
		HandleNewConnections(nil).
		Reader(Reader)

	for {
		if conn, err = listener.Accept(); err != nil {
			fmt.Printf("Connection attempt failed with error [%v].\n", err)
			conn.Close()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		swarm.NewConnection <- conn
	}
}

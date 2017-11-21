package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gianebao/dgochat/hive"
)

func main() {

	var (
		port     = flag.Int("port", 60217, "listening port")
		portStr  = strconv.Itoa(*port)
		swarm    *hive.Swarm
		listener net.Listener
		conn     net.Conn
		err      error
	)

	if listener, err = net.Listen("tcp", "0.0.0.0:"+portStr); err != nil {
		fmt.Printf("Failed to listen to port [:%s] with error [%v]. Exit!\n", portStr, err)
		os.Exit(1)
	}

	fmt.Printf("Server now listening to [:%d]. Waiting for incoming connections.\n", *port)

	logchan := hive.NewLogchan()

	go func(l *hive.Logchan) {
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

	swarm = hive.NewSwarm(logchan).
		HandleNewConnections(nil).
		HandleRead(func(m hive.Message) string {
			switch strings.ToUpper(strings.TrimRight(m.Content, "\r\n")) {
			case "HI":
				return "Hello\n"
			case "EXIT":
				m.Worker.Die()
			}

			return ""
		})

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

package main

import (
	"fmt"
	"github.com/nicwaller/udp_pubsub"
	"log"
	"net"
	"time"
)

func main() {
	// this is the only address that works for broadcast
	// as far as I know
	const broadcastAddress = "255.255.255.255"

	addr := fmt.Sprintf("%s:%d", broadcastAddress, udp_pubsub.PORT)
	fmt.Printf("publishing to %s\n", addr)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatalf("failed binding to broadcast address: %v", err)
		return
	} else {
		defer func(conn net.Conn) {
			if err := conn.Close(); err != nil {
				log.Fatalf("failed to close connection: %v", err)
			}
		}(conn)
	}

	for i := 0; ; i++ {
		event := fmt.Sprintf("event #%d", i)
		fmt.Printf("publishing %s...", event)
		_, err := conn.Write([]byte(event))

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("sent")
		}
		time.Sleep(time.Second)
	}
}

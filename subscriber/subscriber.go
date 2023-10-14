package main

import (
	"context"
	"fmt"
	"github.com/nicwaller/udp_pubsub"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"syscall"
)

func main() {
	fmt.Println("starting up")
	for datagram := range listenUdpBroadcast(udp_pubsub.PORT) {
		fmt.Printf("%s\n", datagram)
	}
	fmt.Println("exiting")
}

func listenUdpBroadcast(port int) chan []byte {
	datagramChannel := make(chan []byte)

	conn, err := listenUdpReuse(port)
	if err != nil {
		log.Printf("listen failed: %v", err)
		close(datagramChannel)
		return datagramChannel
	}

	const udpMaxDatagramSize = 65515
	packet := make([]byte, udpMaxDatagramSize)
	go func() {
		for {
			n, _, err := conn.ReadFromUDP(packet)
			if err != nil {
				fmt.Println(err)
				close(datagramChannel)
				break
			} else {
				datagramChannel <- packet[:n]
			}
		}
	}()

	return datagramChannel
}

func listenUdpReuse(port int) (*net.UDPConn, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			err := c.Control(func(fd uintptr) {
				fmt.Printf("enabling socket reuse (SO_REUSEPORT) on %s\n", address)
				opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			})
			if err != nil {
				return err
			}
			return opErr
		},
	}

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("listening on %s\n", addr)
	lp, err := lc.ListenPacket(context.Background(), "udp", addr)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}

	return lp.(*net.UDPConn), nil
}

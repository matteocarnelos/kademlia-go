package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
	ListenPort int
}

func Listen(ip string, port int) {
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	conn, _ := net.ListenUDP("udp", &addr)
	go func() {
		buf := make([]byte, 1024)
		for {
			_, addr, _ := conn.ReadFromUDP(buf)
			fmt.Printf("Received message from %v: %s\n", addr.IP, buf)
		}
	}()
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO (M1.a)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	addr := net.UDPAddr{
		IP:   net.ParseIP(contact.Address),
		Port: network.ListenPort,
	}
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, "FIND %s", contact.ID)
	conn.Close()
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

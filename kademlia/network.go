package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
	ListenPort int
	RoutingTable *RoutingTable
}

func Listen(ip string, port int) {
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	conn, _ := net.ListenUDP("udp", &addr)
	buf := make([]byte, 1024)
	for {
		_, addr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("%s -> %s\n", addr.IP, buf)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// "create" address
	addr := net.UDPAddr{
		IP:   net.ParseIP(contact.Address),
		Port: network.ListenPort,
	}
	// create connection
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, "PING")
	fmt.Printf("PING -> %s\n", contact.Address)
	conn.Close()

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	for _ , c := range network.RoutingTable.FindClosestContacts(contact.ID, 3) {
		addr := net.UDPAddr{
			IP:   net.ParseIP(c.Address),
			Port: network.ListenPort,
		}
		conn, _ := net.DialUDP("udp", nil, &addr)
		fmt.Fprintf(conn, "FIND_NODE %s", contact.ID)
		fmt.Printf("FIND_NODE %s -> %s\n", contact.ID, c.Address)
		conn.Close()
		network.SendPingMessage(&c)
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

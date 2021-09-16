package kademlia

import (
	"fmt"
	"net"
	"strings"
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
	for {
		buf := make([]byte, 1024)
		_, addr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("%s -> %s\n", addr.IP, buf)
		cmdLine := strings.Fields(string(buf))
		id := cmdLine[0]
		cmd := cmdLine[1]
		var _ []string
		if len(cmdLine) > 2 {
			_ = cmdLine[2:]
		}
		switch cmd {
		case "PING":
			addr.Port = port
			conn, _ := net.DialUDP("udp", nil, addr)
			fmt.Fprintf(conn, "%s PINGREPLY", id)
			fmt.Printf("%s PINGREPLY -> %s\n", id, addr.IP)
			conn.Close()
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// "create" address
	addr := net.UDPAddr{
		IP:   net.ParseIP(contact.Address),
		Port: network.ListenPort,
	}
	// create RPC id
	id := NewRandomKademliaID()

	// create connection
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, "%s PING" , id)
	fmt.Printf("%s PING -> %s\n", id, contact.Address)
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

package kademlia

import (
	"fmt"
	"net"
	"strings"
)

type Network struct {
	ListenIP net.IP
	ListenPort int
}

func (n *Network) listen(handler *Kademlia) {
	addr := net.UDPAddr{
		IP:   n.ListenIP,
		Port: n.ListenPort,
	}
	conn, _ := net.ListenUDP("udp", &addr)
	buf := make([]byte, 1024)
	for {
		size, addr, _ := conn.ReadFromUDP(buf)
		msg := string(buf[:size])
		cmdLine := strings.Fields(msg)
		id := cmdLine[0]
		var cmd string
		var args []string
		if len(cmdLine) > 1 {
			cmd = cmdLine[1]
		}
		if len(cmdLine) > 2 {
			args = cmdLine[2:]
		}
		fmt.Printf("%s -> %s\n", addr.IP, msg)
		resp := handler.handleRPC(id, cmd, args)
		if resp != "" {
			addr.Port = n.ListenPort
			conn, _ := net.DialUDP("udp", nil, addr)
			msg := fmt.Sprintf("%s %s", id, resp)
			fmt.Fprintf(conn, msg)
			fmt.Printf("%s -> %s\n", msg, addr.IP)
			conn.Close()
		}
	}
}

func (n *Network) SendPingMessage(contact *Contact) {
	// TODO (M1.a)
}

func (n *Network) SendFindContactMessage(target *Contact, recipient *Contact) {
	addr := net.UDPAddr{
		IP: net.ParseIP(recipient.Address),
		Port: n.ListenPort,
	}
	id := NewRandomKademliaID()
	conn, _ := net.DialUDP("udp", nil, &addr)
	msg := fmt.Sprintf("%s FIND_NODE %s", id, target.ID)
	fmt.Fprintf(conn, msg)
	fmt.Printf("%s -> %s\n", msg, recipient.Address)
	conn.Close()
}

func (n *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (n *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

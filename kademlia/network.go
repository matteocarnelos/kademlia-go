package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

type Network struct {
	RPC map[KademliaID]chan []string
	RT *RoutingTable
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
		h := sha1.New()
		h.Write(addr.IP)
		n.RT.AddContact(NewContact(NewKademliaID(hex.EncodeToString(h.Sum(nil))), addr.IP.String()))
		msg := string(buf[:size])
		cmdLine := strings.Fields(msg)
		id := NewKademliaID(cmdLine[0])
		fmt.Printf("%s -> %s\n", addr.IP, msg)
		if n.RPC[*id] != nil {
			n.RPC[*id] <- cmdLine[1:]
			continue
		}
		cmd := cmdLine[1]
		var args []string
		if len(cmdLine) > 2 {
			args = cmdLine[2:]
		}
		resp := handler.handleRPC(cmd, args)
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

func (n *Network) SendFindContactMessage(target *Contact, recipient *Contact) *KademliaID {
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
	n.RPC[*id] = make(chan []string)
	return id
}

func (n *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (n *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Network struct {
	RPC sync.Map
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
		h.Write(addr.IP.To4())
		msg := string(buf[:size])
		cmdLine := strings.Fields(msg)
		id := NewKademliaID(cmdLine[0])
		fmt.Printf("%s -> %s\n", addr.IP, msg[41:])
		n.updateRoutingTable(NewContact(NewKademliaID(hex.EncodeToString(h.Sum(nil))), addr.IP.String()))
		if ch, b := n.RPC.Load(*id); b {
			ch.(chan []string) <- cmdLine[1:]
			continue
		}
		cmd := cmdLine[1]
		var args []string
		if len(cmdLine) > 2 {
			args = cmdLine[2:]
		}
		resp := handler.handleRPC(cmd, args)
		addr.Port = n.ListenPort
		conn, _ := net.DialUDP("udp", nil, addr)
		msg = fmt.Sprintf("%s %s", id, resp)
		fmt.Fprintf(conn, msg)
		fmt.Printf("%s -> %s\n", msg[41:], addr.IP)
		conn.Close()
	}
}

func (n *Network) sendUDP(destination net.IP, msg string) {
	addr := net.UDPAddr{
		IP: destination,
		Port: n.ListenPort,
	}
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, msg)
	fmt.Printf("%s -> %s\n", msg[41:], destination)
	conn.Close()
}

func (n *Network) updateRoutingTable(contact Contact) {
	if n.RT.buckets[n.RT.getBucketIndex(contact.ID)].Len() < bucketSize {
		n.RT.AddContact(contact)
	} else {
		// TODO: Ping oldest contact
	}
}

func (n *Network) SendPingMessage(recipient *Contact) *KademliaID {
	id := NewRandomKademliaID()
	msg := fmt.Sprintf("%s PING", id)
	n.sendUDP(net.ParseIP(recipient.Address), msg)
	n.RPC.Store(*id, make(chan []string))
	return id
}

func (n *Network) SendFindContactMessage(target *Contact, recipient *Contact) *KademliaID {
	id := NewRandomKademliaID()
	msg := fmt.Sprintf("%s FIND_NODE %s", id, target.ID)
	n.sendUDP(net.ParseIP(recipient.Address), msg)
	n.RPC.Store(*id, make(chan []string))
	return id
}

func (n *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (n *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

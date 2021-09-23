package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
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
			ch.(chan []string) <-cmdLine[1:]
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

func (n *Network) sendRPC(recipient *Contact, request string) *KademliaID {
	addr := net.UDPAddr{
		IP: net.ParseIP(recipient.Address),
		Port: n.ListenPort,
	}
	id := NewRandomKademliaID()
	msg := fmt.Sprintf("%s %s", id, request)
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, msg)
	fmt.Printf("%s -> %s\n", msg[41:], recipient.Address)
	conn.Close()
	n.RPC.Store(*id, make(chan []string))
	return id
}

func (n *Network) updateRoutingTable(contact Contact) {
	bucket := n.RT.buckets[n.RT.getBucketIndex(contact.ID)]
	if bucket.Len() < bucketSize {
		n.RT.AddContact(contact)
	} else {
		lrs := bucket.list.Back().Value.(Contact)
		ch, _ := n.RPC.Load(*n.SendPingMessage(&lrs))
		select {
		case <-ch.(chan []string):
			bucket.list.MoveToFront(bucket.list.Back())
		case <-time.After(2 * time.Second):
			bucket.list.Remove(bucket.list.Back())
			bucket.list.PushFront(contact)
		}
	}
}

func (n *Network) SendPingMessage(recipient *Contact) *KademliaID {
	return n.sendRPC(recipient, "PING")
}

func (n *Network) SendFindContactMessage(target *Contact, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("FIND_NODE %s", target.ID)
	return n.sendRPC(recipient, req)
}

func (n *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (n *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

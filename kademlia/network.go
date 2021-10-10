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

const pingTimeoutSec = 3
const findTimeoutSec = 20
const storeTimeoutSec = 10
const bufferSize = 8192

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
	buf := make([]byte, bufferSize)
	for {
		size, addr, _ := conn.ReadFromUDP(buf)
		h := sha1.New()
		h.Write(addr.IP.To4())
		msg := string(buf[:size])
		cmdLine := strings.Fields(msg)
		id := NewKademliaID(cmdLine[0])
		fmt.Printf("%s -> %s\n", addr.IP, msg[41:])
		contact := NewContact(NewKademliaID(hex.EncodeToString(h.Sum(nil))), addr.IP.String())
		n.updateRoutingTable(contact)
		//handler.updateStorage(contact)
		if ch, ok := n.RPC.Load(*id); ok {
			ch.(chan []string) <-cmdLine[1:]
			close(ch.(chan []string))
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
	n.RPC.Store(*id, make(chan []string, 10))
	msg := fmt.Sprintf("%s %s", id, request)
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, msg)
	fmt.Printf("%s -> %s\n", msg[41:], recipient.Address)
	conn.Close()
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
		case <-time.After(pingTimeoutSec * time.Second):
			bucket.list.Remove(bucket.list.Back())
			bucket.list.PushFront(contact)
		}
	}
}

func (n *Network) SendPingMessage(recipient *Contact) *KademliaID {
	return n.sendRPC(recipient, "PING")
}

func (n *Network) SendFindContactMessage(target *KademliaID, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("FIND_NODE %s", target)
	return n.sendRPC(recipient, req)
}

func (n *Network) SendFindDataMessage(hash string, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("FIND_VALUE %s", hash)
	return n.sendRPC(recipient, req)
}

func (n *Network) SendStoreMessage(data []byte, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("STORE %s", data)
	return n.sendRPC(recipient, req)
}

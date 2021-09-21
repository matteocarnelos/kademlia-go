package kademlia

import (
	"fmt"
	"net"
	"strings"
)

type Kademlia struct {
	rpc map[KademliaID]chan []string
	Net Network
}

func NewKademlia(me Contact) *Kademlia {
	return &Kademlia{
		rpc: make(map[KademliaID]chan []string),
		Net: Network{RT: NewRoutingTable(me)},
	}
}

func (k *Kademlia) StartListen(ip string, port int) {
	k.Net.ListenIP = net.ParseIP(ip)
	k.Net.ListenPort = port
	go k.Net.listen(k)
}

func (k *Kademlia) handleRPC(id *KademliaID, args []string) string {
	if k.rpc[*id] != nil {
		k.rpc[*id] <- args
		return ""
	}
	switch args[0] {
	case "FIND_NODE":
		resp := ""
		for _, c := range k.Net.RT.FindClosestContacts(NewKademliaID(args[1]), 3) {
			resp += fmt.Sprintf("%s,%d,%s ", c.Address, k.Net.ListenPort, c.ID)
		}
		return strings.TrimSpace(resp)
	}
	return ""
}

func (k *Kademlia) LookupContact(target *Contact) {
	for _, c := range k.Net.RT.FindClosestContacts(target.ID, 3) {
		id := k.Net.SendFindContactMessage(target, &c)
		k.rpc[*id] = make(chan []string)
		
	}
}

func (k *Kademlia) LookupData(hash string) {
	// TODO (M2.b)
}

func (k *Kademlia) Store(data []byte) {
	// TODO (M2.a)
}

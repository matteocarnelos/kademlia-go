package kademlia

import (
	"fmt"
	"net"
	"strings"
)

const concurrencyParam = 3
const replicationParam = 20

type Kademlia struct {
	Net Network
}

func NewKademlia(me Contact) *Kademlia {
	return &Kademlia{
		Net: Network{
			RPC: make(map[KademliaID]chan []string),
			RT: NewRoutingTable(me),
		},
	}
}

func (k *Kademlia) StartListen(ip string, port int) {
	k.Net.ListenIP = net.ParseIP(ip)
	k.Net.ListenPort = port
	go k.Net.listen(k)
}

func (k *Kademlia) handleRPC(cmd string, args []string) string {
	switch cmd {
	case "FIND_NODE":
		resp := ""
		for _, c := range k.Net.RT.FindClosestContacts(NewKademliaID(args[0]), replicationParam) {
			resp += fmt.Sprintf("%s,%d,%s ", c.Address, k.Net.ListenPort, c.ID)
		}
		return strings.TrimSpace(resp)
	}
	return ""
}

func (k *Kademlia) LookupContact(target *Contact) {
	var id []KademliaID
	for _, c := range k.Net.RT.FindClosestContacts(target.ID, concurrencyParam) {
		id = append(id, *k.Net.SendFindContactMessage(target, &c))
	}
}

func (k *Kademlia) LookupData(hash string) {
	// TODO (M2.b)
}

func (k *Kademlia) Store(data []byte) {
	// TODO (M2.a)
}

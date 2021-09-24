package kademlia

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const concurrencyParam = 3
const replicationParam = 20

type Kademlia struct {
	Net Network
}

func NewKademlia(me Contact) *Kademlia {
	return &Kademlia{
		Net: Network{
			RPC: sync.Map{},
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
	case "PING":
		return ""
	case "FIND_NODE":
		resp := ""
		for _, c := range k.Net.RT.FindClosestContacts(NewKademliaID(args[0]), replicationParam) {
			resp += fmt.Sprintf("%s,%d,%s ", c.Address, k.Net.ListenPort, c.ID)
		}
		return strings.TrimSpace(resp)
	}
	return ""
}

func (k *Kademlia) LookupContact(target *Contact) []Contact {
	var closest ContactCandidates
	queried := make(map[string]bool)
	for _, c := range k.Net.RT.FindClosestContacts(target.ID, replicationParam) {
		c.CalcDistance(target.ID)
		queried[c.Address] = false
		closest.Append([]Contact{c})
	}
	for {
		var ids []KademliaID
		closest.Sort()
		for _, c := range closest.GetContacts(replicationParam) {
			if queried[c.Address] { continue }
			ids = append(ids, *k.Net.SendFindContactMessage(target, &c))
			queried[c.Address] = true
			if len(ids) == concurrencyParam { break }
		}
		if len(ids) == 0 {
			return closest.GetContacts(replicationParam)
		}
		for _, id := range ids {
			ch, _ := k.Net.RPC.Load(id)
			select {
			case resp := <-ch.(chan []string):
				for _, t := range resp {
					info := strings.Split(t, ",")
					contact := NewContact(NewKademliaID(info[2]), info[0])
					contact.CalcDistance(target.ID)
					if _, b := queried[contact.Address]; !b {
						closest.Append([]Contact{contact})
						queried[contact.Address] = contact.ID.Equals(k.Net.RT.me.ID)
					}
				}
			case <-time.After(findTimeoutSec * time.Second):
			}
		}
	}
}

func (k *Kademlia) LookupData(hash string) []byte {
	// TODO (M2.b)
	return []byte{}
}

func (k *Kademlia) Store(data []byte) string {
	// TODO (M2.a)
	return ""
}

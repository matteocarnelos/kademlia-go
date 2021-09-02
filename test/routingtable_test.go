package test

import (
	"fmt"
	"github.com/matteocarnelos/kademlia-go/kademlia"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}

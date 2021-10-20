package kademlia

import (
	"testing"
)

const contact1 = "ffffffff00000000000000000000000000000000"
const contact2 = "1111111100000000000000000000000000000000"
const contact3 = "1111111200000000000000000000000000000000"
const contact4 = "1111111300000000000000000000000000000000"
const contact5 = "1111111400000000000000000000000000000000"
const contact6 = "2111111400000000000000000000000000000000"

func TestRoutingTable(t *testing.T) {
	// Sample contacts
	contacts := map[string]struct{}{
		contact1: {},
		contact2: {},
		contact3: {},
		contact4: {},
		contact5: {},
		contact6: {},
	}
	// Test routing table creation
	rt := NewRoutingTable(NewContact(NewKademliaID(contact1), "localhost:8000"))
	// Test routing table population
	rt.AddContact(NewContact(NewKademliaID(contact1), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID(contact2), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(contact3), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(contact4), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(contact5), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID(contact6), "localhost:8002"))
	// Test routing table search
	for _, c := range rt.FindClosestContacts(NewKademliaID(contact6), 20) {
		_ = c.String()
		if _, ok := contacts[c.ID.String()]; !ok {
			t.Error("FindClosestContacts failed: wrong contact list returned")
		}
	}
}

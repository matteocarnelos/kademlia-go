package kademlia

import (
	"fmt"
	"testing"
)

const localAddr = "127.0.0.1"
const listenIP = "0.0.0.0"
const listenPort = 62000

const contactID = "fea50412207cb0a45715ed4f1b3c4b4a6f68ed57"
const contactAddr = "142.250.74.46"

const objContent = "hello"
const objHash = "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"

const nullID = "0000000000000000000000000000000000000000"

var kdm *Kademlia
var contact Contact

func TestNewKademlia(t *testing.T) {
	// Test kademlia initialization
	kdm = NewKademlia(NewContact(NewRandomKademliaID(), localAddr))
	// Add contact to routing table (for later use)
	contact = NewContact(NewKademliaID(contactID), contactAddr)
	kdm.Net.RT.AddContact(contact)
}

func TestStartListen(t *testing.T) {
	// Start listening on default address and port
	kdm.StartListen(listenIP, listenPort)
}

func TestPingRPC(t *testing.T) {
	// Simple PING should return empty string
	if kdm.handleRPC("PING", []string{}) != "" {
		t.Error("PING RPC failed: empty string not returned")
	}
}

func TestStoreRPC(t *testing.T) {
	// New object, should store it and return empty string
	if kdm.handleRPC("STORE", []string{objContent}) != "" {
		t.Error("STORE RPC failed: empty string not returned")
	}
	// Already existent object, should refresh it and return empty string
	if kdm.handleRPC("STORE", []string{objContent}) != "" {
		t.Error("STORE RPC failed: empty string not returned")
	}
}

func TestFindNodeRPC(t *testing.T) {
	// Since we have just one node in the routing table, just that node is expected to be returned
	expected := fmt.Sprintf("%s,%d,%s", contactAddr, listenPort, contactID)
	if kdm.handleRPC("FIND_NODE", []string{nullID}) != expected {
		t.Error("FIND_NODE failed: wrong or no node list returned")
	}
}

func TestFindValueRPC(t *testing.T) {
	// Object contained in the local hash table, should return its content
	if kdm.handleRPC("FIND_VALUE", []string{objHash}) != objContent {
		t.Error("FIND_VALUE failed: wrong or no object returned")
	}
	// Object not contained in the local hash table, should return the closest nodes
	expected := fmt.Sprintf("%s,%d,%s", contactAddr, listenPort, contactID)
	if kdm.handleRPC("FIND_VALUE", []string{nullID}) != expected {
		t.Error("FIND_VALUE failed: wrong or no node list returned")
	}
}

func TestUnknownRPC(t *testing.T) {
	// Unknown RPCs should return the empty string
	if kdm.handleRPC("FOO", []string{}) != "" {
		t.Error("Unknown RPC failed: empty string not returned")
	}
}

func TestUpdateStorage(t *testing.T) {
	// New contact further from test object, should not transfer any data
	kdm.updateStorage(NewContact(NewKademliaID(nullID), contactAddr))
	// New contact closer to test object, should transfer the test object
	kdm.updateStorage(NewContact(NewKademliaID(objHash), contactAddr))
}

func TestLookupContact(t *testing.T) {
	// Perform lookup on random contact
	kdm.LookupContact(NewRandomKademliaID())
}

func TestLookupData(t *testing.T) {
	// Object contained in the local hash table, should trigger the quick lookup and return its content immediately
	if data, ok := kdm.LookupData(objHash); !ok || data != objContent {
		t.Error("LookupData failed: wrong or no object found")
	}
	// Object not contained in the local hash table, should start the proper lookup
	if _, ok := kdm.LookupData(nullID); ok {
		t.Error("LookupData failed: wrong object found")
	}
}

func TestStore(t *testing.T) {
	// Store the test object in the network
	kdm.Store([]byte(objContent))
}

func TestForget(t *testing.T) {
	// Object contained in the local hash table, should stop refreshing it and return true
	if !kdm.ForgetData(objHash) {
		t.Error("ForgetData failed: expected true, false returned instead")
	}
	// Object not contained in the local hash table, should return false
	if kdm.ForgetData(nullID) {
		t.Error("ForgetData failed: expected false, true returned instead")
	}
}

func TestSendPingMessage(t *testing.T) {
	// Simple ping message to the test contact
	kdm.Net.SendPingMessage(&contact)
}

func TestUpdateRoutingTable(t *testing.T) {
	// Update routing table with the test contact
	kdm.Net.updateRoutingTable(contact)
}

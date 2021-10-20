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

const pingTimeoutSec = 3 // Timeout for the PING RPC
const findTimeoutSec = 20 // Timeout for the FIND_NODE and FIND_VALUE RPCs
const storeTimeoutSec = 10 // Timeout for the STORE RPC
const bufferSize = 8192

type Network struct {
	RPC sync.Map // Channel map for communicating with the service layer
	RT *RoutingTable
	ListenIP net.IP
	ListenPort int
}

// listen accounts for incoming messages to the node communicating with
// the service layer and responding to the RPC calls
func (n *Network) listen(handler *Kademlia) {
	addr := net.UDPAddr{ // myListenAddress
		IP:   n.ListenIP,
		Port: n.ListenPort,
	}
	// Announce the listen address to the local network
	conn, _ := net.ListenUDP("udp", &addr)
	buf := make([]byte, bufferSize)
	for {
		size, addr, _ := conn.ReadFromUDP(buf) // Listen for incoming messages
		h := sha1.New()
		h.Write(addr.IP.To4())
		msg := string(buf[:size]) // Obtain the plain text string message
		cmdLine := strings.Fields(msg) // Divide its fields
		id := NewKademliaID(cmdLine[0]) // ID of the RPC
		fmt.Printf("%s -> %s\n", addr.IP, msg[41:])
		// Create a new contact from the sender's address
		contact := NewContact(NewKademliaID(hex.EncodeToString(h.Sum(nil))), addr.IP.String())
		if n.updateRoutingTable(contact) { // If the routing table is updated
			// Update the storage by sending the appropriate values to the new known node
			handler.updateStorage(contact)
		}
		if ch, ok := n.RPC.Load(*id); ok { // If we receive a response
			ch.(chan []string) <-cmdLine[1:] // Send it to the service layer
			close(ch.(chan []string)) // Close the channel
			continue
		}
		// If it's not a response, it's an RPC
		cmd := cmdLine[1] // RPC type
		var args []string
		if len(cmdLine) > 2 {
			args = cmdLine[2:]
		}
		// Call for the handling of the RPC
		resp := handler.handleRPC(cmd, args)
		addr.Port = n.ListenPort
		conn, _ := net.DialUDP("udp", nil, addr)
		msg = fmt.Sprintf("%s %s", id, resp) // Create the message
		fmt.Fprintf(conn, msg) // Send the response back
		fmt.Printf("%s -> %s\n", msg[41:], addr.IP)
		conn.Close()
	}
}

// sendRPC sends the request message to the contact specified
// in the parameters
func (n *Network) sendRPC(recipient *Contact, request string) *KademliaID {
	addr := net.UDPAddr{
		IP: net.ParseIP(recipient.Address),
		Port: n.ListenPort,
	}
	id := NewRandomKademliaID() // Generate an ID for the RPC
	// Store a channel for sending the response to the service layer
	n.RPC.Store(*id, make(chan []string, 10))
	msg := fmt.Sprintf("%s %s", id, request) // Create the message
	conn, _ := net.DialUDP("udp", nil, &addr)
	fmt.Fprintf(conn, msg) // Send the message
	fmt.Printf("%s -> %s\n", msg[41:], recipient.Address)
	conn.Close()
	return id
}

// updateRoutingTable updates the necessary k-bucket of the routing table
// with the contact received as a parameter. It returns true if the table is
// updated and false otherwise
func (n *Network) updateRoutingTable(contact Contact) bool {
	// Obtain the k-bucket associated with the contact's ID
	bucket := n.RT.buckets[n.RT.getBucketIndex(contact.ID)]
	if bucket.Len() < bucketSize { // If the k-bucket is not full
		return n.RT.AddContact(contact) // The contact is added
	}
	// If not obtain the LeastRecentlySeen contact of the k-bucket
	lrs := bucket.list.Back().Value.(Contact)
	ch, _ := n.RPC.Load(*n.SendPingMessage(&lrs)) // We check its availability
	select {
	case <-ch.(chan []string): // If the LeastRecentlySeen node responds
		// Move it to the front of the list
		bucket.list.MoveToFront(bucket.list.Back())
		return false
	case <-time.After(pingTimeoutSec * time.Second): // If the LeastRecentlySeen node does not respond
		bucket.list.Remove(bucket.list.Back()) // Remove it from the k-bucket
		bucket.list.PushFront(contact) // Add the new contact
		return true
	}
}

// SendPingMessage sends a PING RPC to the recipient specified
func (n *Network) SendPingMessage(recipient *Contact) *KademliaID {
	return n.sendRPC(recipient, "PING")
}

// SendFindContactMessage sends a FIND_NODE RPC for the target to the recipient specified
func (n *Network) SendFindContactMessage(target *KademliaID, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("FIND_NODE %s", target)
	return n.sendRPC(recipient, req)
}

// SendFindDataMessage sends a FIND_VALUE RPC for the hash to the recipient specified
func (n *Network) SendFindDataMessage(hash string, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("FIND_VALUE %s", hash)
	return n.sendRPC(recipient, req)
}

// SendStoreMessage sends a STORE RPC for the data to the recipient specified
func (n *Network) SendStoreMessage(data []byte, recipient *Contact) *KademliaID {
	req := fmt.Sprintf("STORE %s", data)
	return n.sendRPC(recipient, req)
}

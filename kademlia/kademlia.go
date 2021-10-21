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

const concurrencyParam = 3   // Alpha definition
const replicationParam = 20  // K definition
const republishDelayHr = 12  // Delay for the republishing routines
const expirationDelayHr = 24 // Delay for the expiration routines

type Kademlia struct {
	hashTable    sync.Map // String map that stores the data
	refreshTable sync.Map // Channel map for communicating with the refreshing routines
	forgetTable  sync.Map // Channel map for communicating with the deleting routines
	Net          Network
}

// NewKademlia creates and returns a new Kademlia object based on the
// information of the contact
func NewKademlia(me Contact) *Kademlia {
	return &Kademlia{
		hashTable:    sync.Map{},
		refreshTable: sync.Map{},
		forgetTable:  sync.Map{},
		Net: Network{
			RPC: sync.Map{},
			RT:  NewRoutingTable(me),
		},
	}
}

// StartListen associates the ip and port specified to the listen parameters
// of the Network and  calls for the network layer to start listening
func (k *Kademlia) StartListen(ip string, port int) {
	k.Net.ListenIP = net.ParseIP(ip)
	k.Net.ListenPort = port
	go k.Net.listen(k)
}

// ForgetData stops the updating routine of the refresher node. Returns true if the node holds
// the data and false otherwise
func (k *Kademlia) ForgetData(hash string) bool {
	if ch, ok := k.forgetTable.Load(hash); ok { // If the node is the refresher of the data
		ch.(chan interface{}) <- nil // Stop the updating routine
		k.forgetTable.Delete(hash)   // Delete the channel
		return true
	}
	return false
}

// handleRPC executes the code associated with the handling of the RPC
// specified in the parameters and returns the generated response
func (k *Kademlia) handleRPC(cmd string, args []string) string {
	switch cmd {
	case "PING":
		return ""
	case "STORE":
		// Obtain the hash of the data
		h := sha1.New()
		h.Write([]byte(args[0]))
		key := hex.EncodeToString(h.Sum(nil))
		// If the value is loaded then it is a refresh STORE
		if ch, ok := k.refreshTable.LoadOrStore(key, make(chan interface{})); ok {
			ch.(chan interface{}) <- nil // Notify the refreshing routine
		} else { // If the value is not stored
			k.hashTable.Store(key, args[0])   // Store the value
			ch, _ := k.refreshTable.Load(key) // Obtain a channel for the refreshing routine
			go func() {                       // Create an anonymous parallel function
				for {
					select {
					case <-ch.(chan interface{}): // If it receives a "notification" it restarts the timeout
					case <-time.After(expirationDelayHr * time.Hour): // If the timeout is completed
						k.refreshTable.Delete(key) // The channel is deleted
						k.hashTable.Delete(key)    // The data is deleted
						return
					}
				}
			}()
		}
		return ""
	case "FIND_VALUE":
		key := args[0]
		if data, ok := k.hashTable.Load(key); ok { // If the data is present in the hash table
			ch, _ := k.refreshTable.Load(key) // Obtain the channel associated with that value
			ch.(chan interface{}) <- nil      // Refresh the timeout
			return data.(string)              // Return the value
		}
		fallthrough // If not execute the following case clause
	case "FIND_NODE":
		resp := ""
		// Look for the k-closest contacts to the hash
		for _, c := range k.Net.RT.FindClosestContacts(NewKademliaID(args[0]), replicationParam) {
			resp += fmt.Sprintf("%s,%d,%s ", c.Address, k.Net.ListenPort, c.ID) // Format the information
		}
		return strings.TrimSpace(resp)
	}
	return ""
}

// updateStorage checks for each value stored in the hash table if the necessary
// requirements for data transfer to the new contact are met
func (k *Kademlia) updateStorage(contact Contact) {
	k.hashTable.Range(func(hash, value interface{}) bool { // For each element of the hashTable
		key := NewKademliaID(hash.(string))
		// Calculate the distance of the contact to the key
		contact.CalcDistance(key)
		// Calculate my distance to the key
		k.Net.RT.me.CalcDistance(key)
		if contact.Less(&k.Net.RT.me) { // If the contact is closer
			// For each of the k-closest contacts to the key
			for _, c := range k.Net.RT.FindClosestContacts(key, replicationParam) {
				if c.ID.Equals(contact.ID) {
					continue
				} // If it is the same contact, continue to the next
				// Calculate its distance to the key
				c.CalcDistance(key)
				if c.Less(&k.Net.RT.me) { // If its distance is closer to the key than me
					// Only one node sends the STORE message
					return true // Continue to the next value
				}
			}
			// Send the STORE RPC to the contact with the data
			k.Net.SendStoreMessage([]byte(value.(string)), &contact)
		}
		return true // Continue to the next value
	})
}

// LookupContact returns a list of the k-closest contacts to the target
func (k *Kademlia) LookupContact(target *KademliaID) []Contact {
	var closest ContactCandidates
	queried := make(map[string]bool)
	// For each contact of the k closest to the target
	for _, c := range k.Net.RT.FindClosestContacts(target, replicationParam) {
		c.CalcDistance(target) // Calculate the distance to the target
		queried[c.Address] = false
		closest.Append([]Contact{c}) // Add it to the closest list
	}
	for {
		var ids []KademliaID
		closest.Sort()                                            // Sort the contacts by their distance
		for _, c := range closest.GetContacts(replicationParam) { // For each contact of the k-closest
			if queried[c.Address] {
				continue
			} // If it has already been queried, continue to the next
			ids = append(ids, *k.Net.SendFindContactMessage(target, &c)) // Send a FIND_NODE RPC
			queried[c.Address] = true
			if len(ids) == concurrencyParam {
				break
			} // If it has reached alpha contacts then finish
		}
		if len(ids) == 0 { // If all contacts were queried
			return closest.GetContacts(replicationParam)
		}
		for _, id := range ids { // For each of the alpha contacts with the FIND_NODE RPC
			ch, _ := k.Net.RPC.Load(id) // Obtain the channel for communicating with the network layer
			select {
			case resp := <-ch.(chan []string): // If the node responds
				for _, t := range resp { // For each string of the message
					triple := strings.Split(t, ",") // Split it by commas
					// Create the new contact with the information received from the node
					contact := NewContact(NewKademliaID(triple[2]), triple[0])
					contact.CalcDistance(target)
					if _, ok := queried[contact.Address]; !ok { // If the "new" contact was not queried
						closest.Append([]Contact{contact}) // Append it to the 'closest' struct
						// Check whether the contact received is me and if so, consider it queried
						queried[contact.Address] = contact.ID.Equals(k.Net.RT.me.ID)
					}
				}
			case <-time.After(findTimeoutSec * time.Second): // If the node does not respond continue
			}
		}
	}
}

// LookupData returns the data associated with the hash if it is in the hashTable
// or a list of the k-closest contacts to the hash otherwise
func (k *Kademlia) LookupData(hash string) (interface{}, bool) {
	if data, ok := k.hashTable.Load(hash); ok { // If the data is stored
		// Obtain the channel associated with the refreshing routine
		ch, _ := k.refreshTable.Load(hash)
		ch.(chan interface{}) <- nil // Refresh the data
		return data.(string), true
	}
	target := NewKademliaID(hash)
	var closest ContactCandidates
	queried := make(map[string]bool)
	// For each contact of the k closest to the target
	for _, c := range k.Net.RT.FindClosestContacts(target, replicationParam) {
		c.CalcDistance(target) // Calculate the distance to the target
		queried[c.Address] = false
		closest.Append([]Contact{c}) // Add it to the closest list
	}
	for {
		var ids []KademliaID
		closest.Sort()                                            // Sort the contacts by their distance
		for _, c := range closest.GetContacts(replicationParam) { // For each contact of the k-closest
			if queried[c.Address] {
				continue
			} // If it has already been queried, continue to the next
			ids = append(ids, *k.Net.SendFindDataMessage(target.String(), &c)) // Send a FIND_VALUE RPC
			queried[c.Address] = true
			if len(ids) == concurrencyParam {
				break
			} // If it has reached alpha contacts then finish
		}
		if len(ids) == 0 { // If all contacts were queried
			return closest.GetContacts(replicationParam), false
		}
		for _, id := range ids { // For each of the alpha contacts with the FIND_VALUE RPC
			ch, _ := k.Net.RPC.Load(id) // Obtain the channel for communicating with the network layer
			select {
			case resp := <-ch.(chan []string): // If the node responds
				for _, t := range resp { // For each string of the message
					triple := strings.Split(t, ",") // Split it by commas
					if len(triple) == 1 {           // If the message contains only one string
						// We return the data
						return triple[0], true
					}
					// Create the new contact with the information received from the node
					contact := NewContact(NewKademliaID(triple[2]), triple[0])
					contact.CalcDistance(target)
					if _, ok := queried[contact.Address]; !ok { // If the "new" contact was not queried
						closest.Append([]Contact{contact}) // Append it to the 'closest' struct
						// Check whether the contact received is me and if so, consider it queried
						queried[contact.Address] = contact.ID.Equals(k.Net.RT.me.ID)
					}
				}
			case <-time.After(findTimeoutSec * time.Second): // If the node does not respond continue
			}
		}
	}
}

// Store puts the data in the hashTable if I am one of the closest contacts and
// sends STORE RPCs to the rest of the k-closest
func (k *Kademlia) Store(data []byte) string {
	// Obtain the hash from the data
	h := sha1.New()
	h.Write(data)
	key := hex.EncodeToString(h.Sum(nil))
	var ids []KademliaID
	for _, c := range k.LookupContact(NewKademliaID(key)) { // For each of the k-closest contacts to the hash
		if c.ID.Equals(k.Net.RT.me.ID) { // If I am one of the closest, I store the value
			k.handleRPC("STORE", []string{string(data)})
		} else { // If not send a STORE RPC to that contact
			ids = append(ids, *k.Net.SendStoreMessage(data, &c))
		}
	}
	for _, id := range ids { // For each of the alpha contacts with the STORE RPC
		ch, _ := k.Net.RPC.Load(id) // Obtain the channel for communicating with the network layer
		select {
		case <-ch.(chan []string): // If the node responds
		case <-time.After(storeTimeoutSec * time.Second): // If the node does not respond continue
		}
	}
	ch, _ := k.forgetTable.LoadOrStore(key, make(chan interface{}))
	go func() { // Create an anonymous parallel function
		select {
		case <-ch.(chan interface{}): // If a forget command was sent
			return
		case <-time.After(republishDelayHr * time.Hour): // If the node does not respond continue
			k.Store(data) // Refresh the data with the topology of the network
		}
	}()
	return key
}

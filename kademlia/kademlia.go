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

const concurrencyParam = 3
const replicationParam = 20
const republishDelayHr = 24
const expirationDelayHr = 24

type Kademlia struct {
	hashTable sync.Map
	refreshTable sync.Map
	Net       Network
}

func NewKademlia(me Contact) *Kademlia {
	return &Kademlia{
		hashTable: sync.Map{},
		refreshTable: sync.Map{},
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
	case "STORE":
		h := sha1.New()
		h.Write([]byte(args[0]))
		key := hex.EncodeToString(h.Sum(nil))
		if ch, ok := k.refreshTable.LoadOrStore(key, make(chan interface{})); ok {
			ch.(chan interface{}) <-nil
		} else {
			k.hashTable.Store(key, args[0])
			ch, _ := k.refreshTable.Load(key)
			go func() {
				for {
					select {
					case <-ch.(chan interface{}):
					case <-time.After(expirationDelayHr * time.Hour):
						k.refreshTable.Delete(key)
						k.hashTable.Delete(key)
						return
					}
				}
			}()
		}
		return ""
	case "FIND_VALUE":
		key := args[0]
		if data, ok := k.hashTable.Load(key); ok {
			ch, _ := k.refreshTable.Load(key)
			ch.(chan interface{}) <-nil
			return data.(string)
		}
		fallthrough
	case "FIND_NODE":
		resp := ""
		for _, c := range k.Net.RT.FindClosestContacts(NewKademliaID(args[0]), replicationParam) {
			resp += fmt.Sprintf("%s,%d,%s ", c.Address, k.Net.ListenPort, c.ID)
		}
		return strings.TrimSpace(resp)
	}
	return ""
}

func (k *Kademlia) updateStorage(contact Contact) {
	k.hashTable.Range(func(hash, value interface{}) bool {
		key := NewKademliaID(hash.(string))
		contact.CalcDistance(key)
		k.Net.RT.me.CalcDistance(key)
		if contact.Less(&k.Net.RT.me) {
			k.Net.SendStoreMessage([]byte(value.(string)), &contact)
		}
		return true
	})
}

func (k *Kademlia) LookupContact(target *KademliaID) []Contact {
	var closest ContactCandidates
	queried := make(map[string]bool)
	for _, c := range k.Net.RT.FindClosestContacts(target, replicationParam) {
		c.CalcDistance(target)
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
					triple := strings.Split(t, ",")
					contact := NewContact(NewKademliaID(triple[2]), triple[0])
					contact.CalcDistance(target)
					if _, ok := queried[contact.Address]; !ok {
						closest.Append([]Contact{contact})
						queried[contact.Address] = contact.ID.Equals(k.Net.RT.me.ID)
					}
				}
			case <-time.After(findTimeoutSec * time.Second):
			}
		}
	}
}

func (k *Kademlia) LookupData(hash string) (interface{}, bool) {
	if data, ok := k.hashTable.Load(hash); ok {
		ch, _ := k.refreshTable.Load(hash)
		ch.(chan interface{}) <-nil
		return data.(string), true
	}
	target := NewKademliaID(hash)
	var closest ContactCandidates
	queried := make(map[string]bool)
	for _, c := range k.Net.RT.FindClosestContacts(target, replicationParam) {
		c.CalcDistance(target)
		queried[c.Address] = false
		closest.Append([]Contact{c})
	}
	for {
		var ids []KademliaID
		closest.Sort()
		for _, c := range closest.GetContacts(replicationParam) {
			if queried[c.Address] { continue }
			ids = append(ids, *k.Net.SendFindDataMessage(target.String(), &c))
			queried[c.Address] = true
			if len(ids) == concurrencyParam { break }
		}
		if len(ids) == 0 {
			return closest.GetContacts(replicationParam), false
		}
		for _, id := range ids {
			ch, _ := k.Net.RPC.Load(id)
			select {
			case resp := <-ch.(chan []string):
				for _, t := range resp {
					triple := strings.Split(t, ",")
					if len(triple) == 1 {
						return triple[0], true
					}
					contact := NewContact(NewKademliaID(triple[2]), triple[0])
					contact.CalcDistance(target)
					if _, ok := queried[contact.Address]; !ok {
						closest.Append([]Contact{contact})
						queried[contact.Address] = contact.ID.Equals(k.Net.RT.me.ID)
					}
				}
			case <-time.After(findTimeoutSec * time.Second):
			}
		}
	}
}

func (k *Kademlia) Store(data []byte) string {
	h := sha1.New()
	h.Write(data)
	key := hex.EncodeToString(h.Sum(nil))
	var ids []KademliaID
	for _, c := range k.LookupContact(NewKademliaID(key)) {
		if c.ID.Equals(k.Net.RT.me.ID) {
			k.handleRPC("STORE", []string{string(data)})
		} else {
			ids = append(ids, *k.Net.SendStoreMessage(data, &c))
		}
	}
	for _, id := range ids {
		ch, _ := k.Net.RPC.Load(id)
		select {
		case <-ch.(chan []string):
		case <-time.After(storeTimeoutSec * time.Second):
		}
	}
	go func() {
		time.Sleep(republishDelayHr * time.Hour)
		k.Store(data)
	}()
	return key
}

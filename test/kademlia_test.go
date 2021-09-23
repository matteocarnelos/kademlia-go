package test

import (
	"github.com/matteocarnelos/kadlab/kademlia"
	"testing"
)

func TestKademlia(t *testing.T) {
	kdm := kademlia.NewKademlia(kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost:8000"))
	kdm.StartListen("0.0.0.0", 8000)
}

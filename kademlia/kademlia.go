package kademlia

type Kademlia struct {
	Network Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	go kademlia.Network.SendFindContactMessage(target)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO (M2.b)
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO (M2.a)
}

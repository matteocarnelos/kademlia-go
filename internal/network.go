package internal

type Network struct {
}

func Listen(ip string, port int) {
	// TODO (M1.b [create])
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO (M1.a)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO (M1.b [join]) (M1.c)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO (M2.b)
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO (M2.a)
}

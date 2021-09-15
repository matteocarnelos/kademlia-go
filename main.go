package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/matteocarnelos/kadlab/kademlia"
	"net"
	"os"
	"strings"
)

const BNIp = "172.17.0.3"
const ListenPort = 62000
const CLIPrefix = ">>>"

func main() {
	iface, _ := net.InterfaceByName("eth0")
	addrs, _ := iface.Addrs()
	ip := addrs[0].(*net.IPNet).IP
	h := sha1.New()
	h.Write(ip)
	id := hex.EncodeToString(h.Sum(nil))
	me := kademlia.NewContact(kademlia.NewKademliaID(id), ip.String())
	fmt.Printf("IP Address: %s\n", ip)
	fmt.Printf("Kademlia ID: %s\n", id)
	kad := kademlia.Kademlia{
		Network: kademlia.Network{
			ListenPort: ListenPort,
			RoutingTable: kademlia.NewRoutingTable(me),
		},
	}
	kademlia.Listen("0.0.0.0", ListenPort)
	if ip.String() == BNIp {
		fmt.Println("Bootstrap Node: Yes")
	} else {
		fmt.Println("Bootstrap Node: No")
		h = sha1.New()
		h.Write([]byte(BNIp))
		id = hex.EncodeToString(h.Sum(nil))
		kad.Network.RoutingTable.AddContact(kademlia.NewContact(kademlia.NewKademliaID(id), BNIp))
		kad.LookupContact(&me)
	}

	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(CLIPrefix + " ")
		if !scanner.Scan() {
			break
		}
		cmdLine := strings.Fields(scanner.Text())
		cmd := ""
		var args []string
		if len(cmdLine) > 0 {
			cmd = cmdLine[0]
		}
		if len(cmdLine) > 1 {
			args = cmdLine[1:]
		}
		switch cmd {
		case "udplisten":
			conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: ListenPort})
			fmt.Printf("Listening on port %d...\n", ListenPort)
			buf := make([]byte, 1024)
			_, addr, _ := conn.ReadFromUDP(buf)
			fmt.Printf("Received message from %v: %s\n", addr.IP, buf)
			conn.Close()
		case "udpsend":
			if len(args) < 2 {
				fmt.Println("udpsend: Too few arguments given")
				fmt.Println("usage: udpsend <dest> <msg>")
				break
			}
			addr := net.UDPAddr{
				IP:   net.ParseIP(args[0]),
				Port: ListenPort,
			}
			conn, _ := net.DialUDP("udp", nil, &addr)
			fmt.Fprintf(conn, args[1])
			fmt.Println("Message sent!")
			conn.Close()
		case "put":
		case "get":
		case "":
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Unsupported command: %s\n", cmd)
		}
	}
}

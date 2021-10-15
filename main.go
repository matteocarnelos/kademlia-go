package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/matteocarnelos/kadlab/kademlia"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const BNHost = 3 // Bootstrap Node Identifier

const ListenPort = 62000
const ListenIP = "0.0.0.0"
const ListenDelaySec = 5

const CLIPrefix = ">>>"

func main() {
	iface, _ := net.InterfaceByName("eth0") // Obtain the interface
	addrs, _ := iface.Addrs()
	ip := addrs[0].(*net.IPNet).IP.To4() // Obtain one address of the interface
	isBN := ip[3] == BNHost
	rand.Seed(int64(ip[3]))
	h := sha1.New()
	h.Write(ip)
	id := kademlia.NewKademliaID(hex.EncodeToString(h.Sum(nil))) // Obtain the ID of the node

	fmt.Printf("IP Address: %s", ip)
	if isBN {
		fmt.Print(" (Bootstrap Node)")
	}
	fmt.Printf("\nKademlia ID: %s\n", id)
	fmt.Println()

	// Create the kademlia object that defines the logic of the service
	me := kademlia.NewContact(id, ip.String())
	kdm := kademlia.NewKademlia(me)
	kdm.StartListen(ListenIP, ListenPort) // Start listening
	delay := time.Duration(ListenDelaySec + rand.Intn(5))
	time.Sleep(delay * time.Second)

	if !isBN { // If it is not the Bootstrap Node
		fmt.Println("Joining network...")
		BNIp := net.IP{ ip[0], ip[1], ip[2], BNHost } // Define the Bootstrap Node's IP
		h = sha1.New()
		h.Write(BNIp)
		BNId := kademlia.NewKademliaID(hex.EncodeToString(h.Sum(nil)))
		kdm.Net.RT.AddContact(kademlia.NewContact(BNId, BNIp.String())) // Add the BN to the routing table
		kdm.LookupContact(me.ID) // Initiate a lookup
		fmt.Println("Network joined!")
		fmt.Println()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() { // CLI interface
		r := csv.NewReader(strings.NewReader(scanner.Text()))
		r.Comma = ' '
		cmdLine, _ := r.Read()
		var cmd string
		var args []string
		if len(cmdLine) > 0 { cmd = cmdLine[0] }
		if len(cmdLine) > 1 { args = cmdLine[1:] }
		switch cmd {
		case "put":
			if len(args) != 1 {
				fmt.Println("Incorrect syntax")
				fmt.Println("Usage: put <data>")
				break
			}
			if len(args[0]) > 255 {
				fmt.Println("Invalid object size, maximum size is 255 bytes")
				break
			}
			fmt.Println("Storing object...")
			hash := kdm.Store([]byte(args[0])) // Store the data
			fmt.Println("Object stored!")
			fmt.Println()
			fmt.Printf("Object hash: %s\n\n", hash)
		case "get":
			if len(args) != 1 {
				fmt.Println("Incorrect syntax")
				fmt.Println("Usage: get <hash>")
				break
			}
			if len(args[0]) != 40 {
				fmt.Println("Invalid hash, please provide a 160-bit data hash")
				break
			}
			fmt.Println("Finding object...")
			if data, ok := kdm.LookupData(args[0]); ok { // Look for the data
				fmt.Println("Object found!")
				fmt.Println()
				fmt.Printf("Object content: %s\n\n", data)
			} else {
				fmt.Println("Object not found")
			}
		case "":
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Command not found: %s\n", cmd)
		}
		fmt.Print(CLIPrefix + " ")
	}
}

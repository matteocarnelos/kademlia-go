package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/matteocarnelos/kadlab/kademlia"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const BNHost = 3

const ListenPort = 62000
const ListenIP = "0.0.0.0"
const ListenDelaySec = 5

const CLIPrefix = ">>>"

func main() {
	iface, _ := net.InterfaceByName("eth0")
	addrs, _ := iface.Addrs()
	ip := addrs[0].(*net.IPNet).IP.To4()
	isBN := ip[3] == BNHost
	rand.Seed(int64(ip[3]))
	h := sha1.New()
	h.Write(ip)
	id := kademlia.NewKademliaID(hex.EncodeToString(h.Sum(nil)))

	fmt.Printf("IP Address: %s", ip)
	if isBN {
		fmt.Print(" (Bootstrap Node)")
	}
	fmt.Printf("\nKademlia ID: %s\n", id)
	fmt.Println()

	me := kademlia.NewContact(id, ip.String())
	kdm := kademlia.NewKademlia(me)
	kdm.StartListen(ListenIP, ListenPort)
	delay := time.Duration(ListenDelaySec + rand.Intn(5))
	time.Sleep(delay * time.Second)

	if !isBN {
		fmt.Println("Joining network...")
		BNIp := net.IP{ ip[0], ip[1], ip[2], BNHost }
		h = sha1.New()
		h.Write(BNIp)
		BNId := kademlia.NewKademliaID(hex.EncodeToString(h.Sum(nil)))
		kdm.Net.RT.AddContact(kademlia.NewContact(BNId, BNIp.String()))
		kdm.LookupContact(&me)
		fmt.Println("Network joined!")
		fmt.Println()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print(CLIPrefix + " ")
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

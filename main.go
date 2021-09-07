package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const CLPrefix = ">>>"
const ListenPort = 62000

func main() {
	fmt.Println("Kademlia node started!")
	int, _ := net.InterfaceByName("eth0")
	addrs, _ := int.Addrs()
	fmt.Printf("IP Address: %s\n", addrs[0].(*net.IPNet).IP)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(CLPrefix + " ")
		if ! scanner.Scan() { break }
		cmdLine := strings.Fields(scanner.Text())
		cmd := ""
		var args []string
		if len(cmdLine) > 0 { cmd = cmdLine[0] }
		if len(cmdLine) > 1 { args = cmdLine[1:] }
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
			addr := net.UDPAddr {
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

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Kademlia node started!")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")
	for scanner.Scan() {
		switch scanner.Text() {
		case "test_listen":
			// Listen to the IP and Port
			listen("127.0.0.1", 62000)

		case "test_send": // test_send <ip_address or hostname>
			// Obtain the IP
			fmt.Print("Introduce an IP: ")
			scanner.Scan()
			inputIP := scanner.Text()

			// Obtain the Port
			fmt.Print("Introduce a port number: ")
			scanner.Scan()
			inputPort, _ := strconv.Atoi(scanner.Text())

			// Obtain the Message
			fmt.Print("Introduce the message: ")
			scanner.Scan()
			inputMessage := scanner.Text()

			send(inputIP, inputPort, []byte(inputMessage))
		case "put":
			fmt.Println("put command received")
		case "get":
			fmt.Println("get command received")
		case "exit":
			os.Exit(0)
		case "":
		default:
			fmt.Println("usage: [put/get/exit] ...")
		}
		fmt.Print(">>> ")
	}
}

func listen(ip string, port int) {
	// Make the UDPAddr type
	var nodeAddr net.UDPAddr
	var netIP net.IP
	netIP = net.ParseIP(ip) // Parse the IP to a series of bytes
	nodeAddr = net.UDPAddr{netIP,port, ""} // Create the object
	fmt.Printf("Listening in port %d\n", nodeAddr.Port)

	// Make the node listen to the IP and port assigned
	nodeCon, err := net.ListenUDP("udp", &nodeAddr)
	if err != nil {
		panic(err) // Stop the execution of this function
	}

	// Receive the message once a connection has been created
	message := make([]byte, 1024)
	n, _, _ := nodeCon.ReadFrom(message)

	messageBuff := bytes.NewBuffer(message)
	messageBuff.Truncate(n)

	// Show the message
	fmt.Printf("Message Received: %s\n", messageBuff)

	// Close the connection
	nodeCon.Close()
}

func send(ip string, port int, message []byte) {
	// Create a nil connection
	nodeCon, _ := net.ListenPacket("udp",":0")

	// Make the UDPAddr type
	var nodeAddr net.UDPAddr
	var netIP net.IP
	netIP = net.ParseIP(ip) // Parse the IP to a series of bytes
	nodeAddr = net.UDPAddr{netIP,port, ""} // Create the object

	// Send the message
	nodeCon.WriteTo(message, &nodeAddr)

	// Close the connection
	nodeCon.Close()
}

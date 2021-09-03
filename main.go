package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Kademlia node started!")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">>> ")
	for scanner.Scan() {
		switch scanner.Text() {
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

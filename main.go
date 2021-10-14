package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/matteocarnelos/kadlab/kademlia"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const BNHost = 3

const ListenPort = 62000
const ListenIP = "0.0.0.0"
const ListenDelaySec = 5

const CLIPrefix = ">>>"

var kdm *kademlia.Kademlia

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("\n%s -> [%s %s %s] %s\n", ip, r.Method, r.URL, r.Proto, body)
	var msg string
	var code int
	switch r.Method {
	case "GET":
		hash := strings.Split(r.URL.String(), "/")[2]
		if len(hash) != 40 {
			code = http.StatusBadRequest
			msg = "Invalid hash, please provide a valid 160-bit data hash"
			break
		}
		if content, ok := load(hash); ok {
			code = http.StatusOK
			msg = content
		} else {
			code = http.StatusNotFound
			msg = "Object not found"
		}
	case "POST":
		if len(body) > 255 {
			code = http.StatusBadRequest
			msg = "Invalid object size, maximum size is 255 bytes"
			break
		}
		hash := store(string(body))
		w.Header().Set("Location", "/objects/" + hash)
		code = http.StatusCreated
		msg = "Object stored!"
	}
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
	fmt.Printf("[%s %d %s] %s -> %s\n\n", r.Proto, code, http.StatusText(code), msg, ip)
}

func store(content string) string {
	fmt.Println("Storing object...")
	hash := kdm.Store([]byte(content))
	fmt.Println("Object stored!")
	fmt.Println()
	return hash
}

func load(hash string) (string, bool) {
	fmt.Println("Finding object...")
	if data, ok := kdm.LookupData(hash); ok {
		fmt.Println("Object found!")
		fmt.Println()
		return data.(string), true
	}
	return "", false
}

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
	kdm = kademlia.NewKademlia(me)
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
		kdm.LookupContact(me.ID)
		fmt.Println("Network joined!")
		fmt.Println()
	}

	http.HandleFunc("/objects", handleRequest)
	http.HandleFunc("/objects/", handleRequest)
	go http.ListenAndServe(":80", nil)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
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
			hash := store(args[0])
			fmt.Printf("Object hash: %s\n\n", hash)
		case "get":
			if len(args) != 1 {
				fmt.Println("Incorrect syntax")
				fmt.Println("Usage: get <hash>")
				break
			}
			if len(args[0]) != 40 {
				fmt.Println("Invalid hash, please provide a valid 160-bit data hash")
				break
			}
			if content, ok := load(args[0]); ok {
				fmt.Printf("Object content: %s\n\n", content)
			} else {
				fmt.Printf("Object not found\n\n")
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

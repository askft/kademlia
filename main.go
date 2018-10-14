package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}
	port := os.Args[1]
	if !validPort(port) {
		printUsage()
		return
	}

	peer, err := NewPeer(&Options{
		id:        GenerateRandomNodeID(),
		host:      getLocalIP(),
		port:      port,
		store:     NewLocalStore(),
		networkID: "v1",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	wg.Add(3)
	{
		go RunServer(peer)
		go HandleInput(peer)
		go RunUI()
	}
	wg.Wait()
}

// HandleInput reads user input from an input channel
// and dispatches commands that depend on the input.
func HandleInput(peer *Peer) {
	defer wg.Done()
	for {
		input := <-uiInputChannel
		action, rest, err := parseUIInput(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch action {
		case "store":
			key, err := peer.store.Put([]byte(rest))
			if err != nil {
				fmt.Println("this shouldnt happen")
				panic(err)
			}
			fmt.Printf("Stored data. Key: [ %s ].", key)
		case "get":
			data, err := peer.store.Get(rest)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Data for key [ %s ] is:\n%s\n", rest, data)
		case "bootstrap":
			peer.Bootstrap("4000")
		default:
			fmt.Println(`  usage:
    store [string]
    get   [key]
    bootstrap`)
		}
	}
}

func getLocalIP() net.IP {
	return net.ParseIP("127.0.0.1")
}

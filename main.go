package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"p2p/encoding"
	"p2p/node"
	"p2p/peer"
	"p2p/store"
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

	p, err := peer.NewPeer(&peer.Options{
		Key:       node.GenerateRandomKey(),
		Host:      getLocalIP(),
		Port:      port,
		Store:     store.NewLocalStore(),
		NetworkID: "v1",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO Temporary workaround hack for bootstrap node.
	// 4000 is the bootstrap node port when running locally!
	if port == "4000" {
		p.Key = node.Key{}
	}

	wg.Add(3)
	{
		go peer.RunServer(p, &wg)
		go HandleInput(p)
		go RunUI()
	}
	wg.Wait()
}

// HandleInput reads user input from an input channel
// and dispatches commands that depend on the input.
func HandleInput(peer *peer.Peer) {
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
			data := []byte(rest)
			key := encoding.HashData(data)
			keyStr := encoding.EncodeHash(key)
			peer.IterativeStore(key, data)
			if _, err := peer.Get(keyStr); err != nil {
				if _, err := peer.Put(data); err != nil {
					panic(err)
				}
			}
			fmt.Printf("Stored data. Key: [ %s ].", keyStr)

		case "get":
			// TODO look first in own store
			keyStr := rest
			key, err := encoding.DecodeKeyStr(keyStr)
			if err != nil {
				fmt.Printf("Could not decode key [ %s ].\n", keyStr)
				panic(err)
			}
			data, contacts := peer.IterativeFindValue(key)
			if data != nil {
				fmt.Printf("Data for key [ %s ] is:\n%s\n", keyStr, string(data))
			} else if contacts != nil {
				fmt.Printf("Data for key [ %s ] could not be found.\n", key)
			} else {
				panic(errors.New("this should not happen"))
			}

		case "bootstrap":
			peer.Bootstrap(node.Contact{
				Key:  node.Key{},
				Host: getLocalIP(),
				Port: "4000",
			})

		case "table":
			peer.PrintAllContacts()

		default:
			fmt.Println(`  usage:
    store [string]
    get   [key]
    bootstrap`)
		}
	}
}

func validPort(data string) bool {
	if port, err := strconv.Atoi(data); err != nil || port < 4000 || port > 5000 {
		return false
	}
	return true
}

func getLocalIP() net.IP {
	return net.ParseIP("127.0.0.1")
}

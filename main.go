package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/pkg/errors"

	"github.com/askft/kademlia/encoding"
	"github.com/askft/kademlia/node"
	"github.com/askft/kademlia/peer"
	"github.com/askft/kademlia/store"
)

var wg sync.WaitGroup

var bootstrapContact = node.Contact{
	Key:  node.Key{},
	Host: getLocalIP(),
	Port: "4000",
}

func main() {
	if len(os.Args) != 2 {
		printUsageAndExit()
	}

	port := os.Args[1]
	if !validPort(port) {
		printUsageAndExit()
	}

	p, err := peer.NewPeer(&peer.Options{
		Key:       node.GenerateRandomKey(),
		Host:      getLocalIP(),
		Port:      port,
		Store:     store.NewMemStore(),
		NetworkID: "v1",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO Temporary workaround hack for bootstrap node.
	// 4000 is the bootstrap node port when running locally!
	if port == "4000" {
		p.Contact.Key = node.Key{}
	}

	ui := NewCommandLineUI()

	server, err := peer.NewServer(p)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create server"))
	}

	wg.Add(3)
	go server.Run(&wg)
	go ui.Run(&wg)
	go handleInput(ui, p)
	wg.Wait()
}

// handleInput reads a message from a user interface
// and dispatches a command depending on the message.
func handleInput(ui UI, peer *peer.Peer) {
	defer wg.Done()

	fmt.Println(uiUsage)

	for {
		message := ui.Get()
		action, rest, err := message.Parse()
		if err != nil {
			log.Println(errors.Wrapf(err, "could not parse message %s", string(message)))
			continue
		}

		switch action {

		case ActionStore:
			data := []byte(rest)
			key := encoding.HashData(data)
			keyStr := encoding.EncodeHash(key)
			peer.IterativeStore(key, data)
			if _, err := peer.Get(keyStr); err != nil {
				if _, err := peer.Put(data); err != nil {
					panic(err)
				}
			}
			log.Printf("Stored data. Key: [ %s ].", keyStr)

		case ActionGet:
			// TODO look first in own store
			keyStr := rest
			key, err := encoding.DecodeKeyStr(keyStr)
			if err != nil {
				log.Printf("Could not decode key [ %s ].\n", keyStr)
				panic(err)
			}
			data, contacts := peer.IterativeFindValue(key)
			if data != nil {
				log.Printf("Data for key [ %s ] is:\n%s\n", keyStr, string(data))
			} else if contacts != nil {
				log.Printf("Data for key [ %s ] could not be found.\n", key)
			} else {
				panic(errors.New("this should not happen"))
			}

		case ActionBootstrap:
			peer.Bootstrap(bootstrapContact)

		case ActionTable:
			peer.PrintAllContacts()

		default:
			fmt.Println(uiUsage)
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

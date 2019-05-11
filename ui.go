package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Message string

type Action string

const (
	ActionStore     = Action("store")
	ActionGet       = Action("get")
	ActionBootstrap = Action("bootstrap")
	ActionTable     = Action("table")
)

func (m Message) Parse() (Action, string, error) {
	sep := " "
	tokens := strings.Split(string(m), sep)
	action := strings.ToLower(tokens[0])
	data := strings.Join(tokens[1:], sep)
	// fmt.Printf("action=[%s], data=[%s]\n", action, data)
	return Action(action), data, nil
}

const uiUsage = `
  usage:
    store [string]  (store a value and returns its key)
    get   [key]     (get a value by its key)
    bootstrap       (connect to the network via the bootstrap node)
`

// UI is a user interface that sends user input to the input channel.
// This could be a GUI, another TCP server, or whatever you like.
type UI interface {
	Get() Message
	Run(*sync.WaitGroup)
}

type CommandLineUI struct {
	c chan Message
}

func NewCommandLineUI() *CommandLineUI {
	return &CommandLineUI{
		make(chan Message),
	}
}

func (ui *CommandLineUI) Get() Message {
	return <-ui.c
}

func (ui *CommandLineUI) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ui.c <- Message(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func printUsageAndExit() {
	fmt.Printf("usage: %s [port]\nport must be in range [4000, 5000]\n", os.Args[0])
	os.Exit(0)
}

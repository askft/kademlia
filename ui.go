package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var uiInputChannel = make(chan string)

// RunUI is a user interface that sends user input to the input channel.
// This could also be a GUI, another TCP server, or whatever you like.
// As of now, the console is used.
func RunUI() {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		uiInputChannel <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func parseUIInput(input string) (string, string, error) {
	sep := " "
	tokens := strings.Split(input, sep)
	action := strings.ToLower(tokens[0])
	data := strings.Join(tokens[1:], sep)
	// fmt.Printf("action=[%s], data=[%s]\n", action, data)
	return action, data, nil
}

func printUsage() {
	fmt.Printf("usage: %s [port]\n", os.Args[0])
}

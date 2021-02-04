package main

import (
	"bufio"
	"fmt"
	"github.com/rcapraro/chat/client"
	"github.com/rcapraro/chat/internal/message"
	"log"
	"os"
)

func main() {
	c := client.NewClient()
	defer c.Close()
	err := c.Connect("")
	if err != nil {
		log.Fatalf("Impossible to Connect to the server...exiting")
	}
	go c.Start()
	err = c.Send(message.NameMessage{
		Name: c.Name,
	})
	if err != nil {
		log.Fatalf("Impossible to Send message to the server...exiting")
	}

	fmt.Printf("Connected as \033[36m%s\033[0m\n", c.Name)

	go func() {
		for msg := range c.MessagesChan {
			//Only display messages from other clients to Stdout
			if msg.ClientName != c.Name {
				fmt.Printf("\n\033[35m# %v says: %v\n\033[32m>\033[0m ", msg.ClientName, msg.Message)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[32m>\033[0m ")
		//Scans a line from Stdin(Console)
		scanner.Scan()
		msg := scanner.Text()
		if len(msg) != 0 {
			_ = c.Send(message.ClientMessage{Message: msg})
		} else {
			//empty string
			continue
		}
	}
}

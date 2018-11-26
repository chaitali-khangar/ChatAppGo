package main

import (
	"bufio"
	"flag"
	"fmt"
	"errors"
	"os"
	"strings"
	"sync"
	"github.com/gautamrege/gochat/api"
)

const helpStr = `Commands
1. /users :- Get list of live users
2. @{user} message :- send message to specified user
3. /exit :- Exit the Chat
4. /all :- Send message to all the users [TODO]`

var (
	name      = flag.String("name", "", "The name you want to chat as")
	port      = flag.Int("port", 12345, "Port that your server will run on.")
	host      = flag.String("host", "", "Host IP that your server is running on.")
	stdReader = bufio.NewReader(os.Stdin)
)

var MyHandle api.Handle
var USERS = PeerHandleMapSync{
	PeerHandleMap: make(map[string]api.Handle),
}

func invalidUser(name string, host string) error{
	if(name == "" || host == ""){
		return errors.New("Please enter name and host")
	}
	return nil
}

func main() {
	// Parse flags for host, port and name
	flag.Parse()

	// TODO-WORKSHOP-STEP-1: If the name and host are empty, return an error with help message
   	error := invalidUser(*name,*host)

	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}


	// TODO-WORKSHOP-STEP-2: Initialize global MyHandle of type api.Handle
	
	MyHandle = api.Handle{
		Host: *host,
		Port: int32(*port),
		Name: *name,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// Broadcast for is-alive on 33333 with own UserHandle.
	go broadcastOwnHandle(&wg)

	// Listener for is-alive broadcasts from other hosts. Listening on 33333
	go listenAndRegisterUsers(&wg)

	// gRPC listener
	go startServer(&wg)

	for {
		fmt.Printf("> ")
		textInput, _ := stdReader.ReadString('\n')
		// convert CRLF to LF
		textInput = strings.Replace(textInput, "\n", "", -1)
		parseAndExecInput(textInput)
	}

	wg.Wait()
}

// Handle the input chat messages as well as help commands
func parseAndExecInput(input string) {
	// Split the line into 2 tokens (cmd and message)
	tokens := strings.SplitN(input, " ", 2)
	cmd := tokens[0]

	switch {
	case cmd == "":
		break
	case cmd == "?":
		fmt.Printf(helpStr)
		break
	case strings.ToLower(cmd) == "/users":
		fmt.Println(USERS)
		break
	case strings.ToLower(cmd) == "/exit":
		os.Exit(1)
		break
	case cmd[0] == '@':
		// TODO-WORKSHOP-STEP-9: Write code to sendChat. Example
		// "@gautam hello golang" should send a message to handle with name "gautam" and message "hello golang"
		// Invoke sendChat to send the  message
		_,ok := USERS.Get(cmd[1:])
		if(ok){
			sendChat(USERS.PeerHandleMap[cmd[1:]],tokens[1])
		}else{
			fmt.Println("User not exist")
		}
		
		break
	case strings.ToLower(cmd) == "/all":
	  for key, _ := range USERS.PeerHandleMap {
	    sendChat(USERS.PeerHandleMap[key],tokens[1])
	  }
	case strings.ToLower(cmd) == "/help":
		fmt.Println(helpStr)
		break
	default:
		fmt.Println(helpStr)
	}
}

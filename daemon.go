package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var eventSocFile = fmt.Sprintf("/tmp/hypr/%v/.socket2.sock", his)
var socFile = "/tmp/hyprman.socket"

// Attempts to check if there is another instance of the daemon running
// by connecting to it. Successful connection means there's a conflict
// and an error is returned. If the connection was unsuccessful, but the `socFile`
// exists, then it is assumed that the previous instance errored out, and the
// socFile is removed.
func conflictCheck(f string) (confict bool) {
	_, err := os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal("error performing sockfile lookup: ", err)
	}
    conn, err := net.Dial("unix", f)
	if err != nil {
		os.Remove(f)
		return
	}
    conn.Write([]byte("ping"))
	fmt.Println("hyprman daemon is already running")
	os.Exit(1)
	return true
}


func handleEvent(s StateManager, msg string){
    if strings.Contains(msg, "activewindow"){
        s.SetActive() 
    }
}

func eventListen(state *State) {
	eventSoc, err := net.Dial("unix", eventSocFile)
	if err != nil {
		log.Fatal("failed open connection to hyprland", err)
	}
	defer eventSoc.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(socFile)
		os.Exit(1)
	}()

	for {
		buf := make([]byte, 1024)
		_, err := eventSoc.Read(buf)
		if err != nil {
			log.Fatal("Unable to read socket", err)
			break
		}
		m := string(buf)
        handleEvent(state, m)
        // msg := strings.Split(m, ">>")
        fmt.Println("\nEvent:", m)

	}
}

func foo() {
	fmt.Println("called foo")
}

var simpleCommands = map[string]func(){
	"fullscreen": foo,
}

func handleCommand(state StateManager, conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 100)
	_, err := conn.Read(buf)
	if err != nil && err.Error() != "EOF" {
		log.Fatal("Unable to read from socket connection ", err)
		return
	}
	rcmd := strings.Split(string(buf), " ")
	args := []string{}
	cmd := rcmd[0]
	if len(rcmd) > 1 {
		args = rcmd[1:]
	}
	cmd = string(bytes.Trim([]byte(cmd), "\x00"))
	if cmd == "ping" {
		conn.Write([]byte("pong"))
        return
	}
	fn, set := EnabledCmmands[cmd]
	if !set {
		fmt.Printf("Unable to find handler for command:'%v'\n", cmd)
		return
	}
	fn(state, args...)
}

// Listens for commands on the socket and attempts to execute them
func commandListen(state *State) {
	socFile := "/tmp/hyprman.socket"
	sock, err := net.Listen("unix", socFile)
	if err != nil {
		log.Fatal("unable to start command socket", err)
	}
	defer func() {
		sock.Close()
		os.Remove(socFile)
	}()
	for {
		conn, err := sock.Accept()
		if err != nil && err.Error() != "EOF" {
			log.Fatal("Unable to read from socket connection ", err)
		}
		go handleCommand(state, conn)
	}
}

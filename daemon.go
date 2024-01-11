package main

import (
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
			log.Fatal("Unabel to read socket", err)
			break
		}
		m := string(buf)
		strings.Split(m, ">>")
	}
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
	cmd := rcmd[0]
	args := []string{}
	if len(rcmd) > 1 {
		args = rcmd[1:]
	}
    fmt.Printf("Got '%v', expect 'fullscreen'\n", cmd)
	fn, set := EnabledCmmands[cmd]
	if !set {
		fmt.Printf("Unable to find handler for command:'%v'\n", cmd)
		return
	}
	fn(state, args...)
	fmt.Printf("buf: %v\n", rcmd)
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

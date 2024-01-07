package main

import (
	"fmt"
	"io"
	// "io/fs"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var state = State{windows: make(map[string]*Window)}

func eventListen() {
	eventSoc, err := net.Dial("unix", eventSocFile)
	if err != nil {
		log.Fatal("failed open connection to hyprland", err)
	}
	defer eventSoc.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
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
		msg := strings.Split(m, ">>")
		if msg[0] == "activewindow" {
			state.SetActive(msg[1])
		}
		fmt.Println("message from hyprland", m)
	}
}

// Listens for commands on the socket and attempts to perform them
func commandListen() {
    // os.Mkdir("/tmp/hyprman/", fs.FileMode(os.O_RDWR))
    socFile := "/tmp/hyprman/command.socket"
	sock, err := net.Listen("unix",socFile )
	   if err != nil{
	       log.Fatal("unable to start command socket", err)
	}
    defer func(){
        os.Remove(socFile)
        sock.Close()
    }()
	command := Router{
		state: &state,
		cmds:  make(map[string]func(StateManager, io.Writer, ...string)),
	}
	command.Register("minimize", Minimize)
	command.Register("restore", Restore)
    for {
        conn, err := sock.Accept()
        if err != nil{
            log.Fatal("Unable to make socket connection", err)
        }
        defer conn.Close()
        buf := make([]byte, 1024)
        _, err = conn.Read(buf)
        if err != nil{
            log.Fatal("Unable to read from socket connection", err)
        }
        fmt.Printf("buf: %v\n", buf)
    }
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Window struct {
	Address        string
	Mapped         bool
	Hidden         bool
	At             [2]int
	Size           [2]int
	Workspace      int
	Floating       bool
	Monitor        int
	Class          string
	Title          string
	InitialClass   string
	InitialTitle   string
	Pid            int
	Xwayland       bool
	Pinned         bool
	Fullscreen     bool
	FullscreenMode int
	FakeFullscreen bool
	Minimzied      bool
}

type Router struct {
	state *State
	cmds  map[string]func(StateManager, io.Writer, ...string)
}

func (r *Router) Register(cmd string, f func(StateManager, io.Writer, ...string)) {
	r.cmds[cmd] = f
}

var his, _ = os.LookupEnv("HYPRLAND_INSTANCE_SIGNATURE")
var ctrlSocFile = fmt.Sprintf("/tmp/hypr/%v/.socket.sock", his)
var eventSocFile = fmt.Sprintf("/tmp/hypr/%v/.socket2.sock", his)

type ConnWriter struct {
	Conn net.Conn
}

func (cw *ConnWriter) Write(p []byte) (n int, err error) {
	return cw.Conn.Write(p)
}

type Args struct {
	flags map[string]string
	args  []string
}

func main() {
	var command string
	daemonMode := flag.Bool("daemon", false, "Run daemon")
	flag.BoolVar(daemonMode, "d", false, "Run daemon")
	flag.StringVar(&command, "command", "", "Command to run")
	flag.Parse()
	go eventListen()
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	ctrlSoc, err := net.Dial("unix", ctrlSocFile)
	if err != nil {
		log.Fatal("failed open connection to hyprland", err)
	}
	defer ctrlSoc.Close()
	_, err = ctrlSoc.Write([]byte(ctrlCommands["float"]))
	if err != nil {
		log.Fatal("Fialed to get current active window", err)
	}
	commandListen()
}

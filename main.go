package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func getHis() string {
	p := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if p == "" {
		log.Fatal("unable to find hyprland instance signature")
	}
	return p
}

var his = getHis()

func main() {
	var command string
	daemonMode := flag.Bool("daemon", false, "Run daemon")
	flag.BoolVar(daemonMode, "d", false, "Run daemon")
	flag.StringVar(&command, "command", "", fmt.Sprint(availableCmds()))
	flag.Parse()
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}
	if command == "" {
		*daemonMode = true
	}
	if *daemonMode {
		conflictCheck(socFile)
		state := StateInit()
		go eventListen(state)
		commandListen(state)
		return
	}
	_, set := EnabledCmmands[command]
	if !set {
		fmt.Println("unkown command", command)
		flag.Usage()
		os.Exit(1)
	}
	conn := daemonConnect()
	defer conn.Close()
	_, err := conn.Write([]byte(command))
	if err != nil {
		log.Fatal("Unable to send command to daemon ", err)
	}

}

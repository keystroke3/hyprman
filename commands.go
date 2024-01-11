package main

import (
	"fmt"
	"log"
	"net"
)

type Command func(StateManager, ...string)

type Executor interface {
	Exec(string) (string, error)
}

var EnabledCmmands map[string]Command = map[string]Command{
	"minimize":   Minimize,
	"restore":    Restore,
	"fullscreen": ForceFullScreen,
}

var hyprctlCmds map[string]string = map[string]string{
	"active":     "dispatch activewindow",
	"minimize":   "dispatch movetoworkspacesilent special hyprman,address:%v",
	"restore":    "dispatch movetoworkspacesilent %v,address:%v",
	"fullscreen": "dispatch fullscreen %v address:%v",
	"pin":        "dispatch pin address:%v",
	"float":      "dispatch togglefloating",
}

func availableCmds() string {
	s := ""
	for cmd := range EnabledCmmands {
		s = s + cmd + "\n"
	}
	return s
}

// toggle pin property on window with `addr`
func togglePin(e Executor, addr string) {
	cmd := fmt.Sprintf(hyprctlCmds["pin"], addr)
	_, err := e.Exec(cmd)
	if err != nil {
		log.Fatal("Unable to toggle pin", err)
	}
}

// toggles fullscreen of type t for window at address addr
// t is one of "0" (real fullscreen) or "1" (monocle)
func toggleFull(e Executor, w *Window, t string) {
	cmd := fmt.Sprintf(hyprctlCmds["fullscreen"], t, w.Address)
	_, err := e.Exec(cmd)
	if err != nil {
		log.Fatal("Unable to toggle fullscreen", err)
	}
	w.Fullscreen = !w.Fullscreen
}


// Provides a way for command functions to access state
type StateManager interface {
	SetActive(addr string)
	ActiveWindow() *Window
	Client() *Client
}

// Sends active window to 'hyprman' workspace and remvoes pin attribte
func Minimize(s StateManager, _ ...string) {
	w := s.ActiveWindow()
	cmd := fmt.Sprintf(hyprctlCmds["minimize"], w.Address)
	_, err := s.Client().Exec(cmd)
	if err != nil {
		log.Println("Unable to minimize", err)
		return
	}
}

// Returns a minimized window to the workspace it was on and reattaches the pin attribute
func Restore(s StateManager, _ ...string) {
	w := s.ActiveWindow()
	cmd := fmt.Sprintf(hyprctlCmds["restore"], fmt.Sprint(w.Workspace), w.Address)
	_, err := s.Client().Exec(cmd)
	if err != nil {
		log.Println("Unable to restore", err)
		return
	}
}

// Forces a window to become fullscreen even if it is pinned
func ForceFullScreen(s StateManager, args ...string) {
	fscreenType := ""
	if len(args) > 0 {
		fscreenType = args[0]
	}
	c := s.Client()
	w := s.ActiveWindow()
	if w.Fullscreen {
		toggleFull(c, w, fscreenType)
		if w.Pinned {
			togglePin(s.Client(), w.Address)
		}
		return
	}
	togglePin(s.Client(), w.Address)
	w.Fullscreen = true
}

func daemonConnect() net.Conn {
    conn, err:=  net.Dial("unix", socFile) 
    if err != nil{
        log.Fatal("Unable to connect to daemon. Is hyprmand is it running?")
    }
    return conn
}

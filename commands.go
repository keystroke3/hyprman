package main

import (
	"fmt"
	"log"
	"net"
)

type StateManager interface {
	SetActive(args ...string)
	ActiveWindow() *Window
	Filter(field string, value any) map[string]*Window
	Client() Executor
    IsPinned(addr string) bool
    SavePinned(addr string)
    RemovePinned(addr string)
}

type Command func(StateManager, ...string)

var EnabledCmmands map[string]Command = map[string]Command{
	"minimize":   Minimize,
	"restore":    Restore,
	"fullscreen": ForceFullScreen,
	"float":      ToggleFloating,
	"pin":        Pin,
}

var hyprctlCmds map[string]string = map[string]string{
	"active":     "hyprctl dispatch activewindow",
	"minimize":   "hyprctl dispatch movetoworkspacesilent special hyprman,address:%v",
	"restore":    "hyprctl dispatch movetoworkspacesilent %v,address:%v",
	"fullscreen": "hyprctl dispatch fullscreen %v address:%v",
	"pin":        "hyprctl dispatch pin address:%v",
	"float":      "hyprctl dispatch togglefloating",
}

func availableCmds() string {
	s := ""
	for cmd := range EnabledCmmands {
		s = s + cmd + "\n"
	}
	return s
}

// toggles fullscreen of type t for window at address addr.
//
// t is one of "0" (real fullscreen) or "1" (monocle)
func toggleFull(e Executor, addr string, t string) error {
	cmd := fmt.Sprintf(hyprctlCmds["fullscreen"], t, addr)
	out, err := e.Exec(cmd)
	if err != nil {
		log.Fatal("Unable to toggle fullscreen", err)
	}
	fmt.Println("fullscreen output:", out)
	return nil
}

// Provides a way for command functions to access state

// Sends active window to 'hyprman' workspace and remvoes pin attribte
func Minimize(s StateManager, _ ...string) {
	w := s.ActiveWindow()
	cmd := fmt.Sprintf(hyprctlCmds["minimize"], w.Address)
	fmt.Println(cmd)
	resp, err := s.Client().Exec(cmd)
	if err != nil {
		log.Println("Unable to minimize", err)
		return
	}
	log.Println("Miminize command:", resp)
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
	if !w.Fullscreen {
		if w.Pinned {
			fmt.Println("Unpinning")
			s.SavePinned(w.Address)
			Pin(s, w.Address)
		} 
		toggleFull(c, w.Address, fscreenType)
		w.Fullscreen = true
		return
	}
	toggleFull(c, w.Address, fscreenType)
	w.Fullscreen = false
	if s.IsPinned(w.Address) {
		fmt.Println("Re-pinning")
		Pin(s, w.Address)
        s.RemovePinned(w.Address)
		return
	}
}

// Toggles pinned status of the active window
func ToggleFloating(s StateManager, _ ...string) {
	c := s.Client()
	cmd := hyprctlCmds["float"]
	resp, err := c.Exec(cmd)
	if err != nil {
		log.Println("unable to toggle floating:", err)
	}
	fmt.Println("Float response:", resp)
}

// Toggles pinned status of the active window or window with given `address`
func Pin(s StateManager, address ...string) {
	var w *Window
	addr := ""
	if len(address) > 0 {
		addr = address[0]
	}
	if addr == "" {
		w = s.ActiveWindow()
	} else {
		w = s.Filter("address", addr)[addr]
	}
	c := s.Client()
	cmd := fmt.Sprintf(hyprctlCmds["pin"], w.Address)
	out, err := c.Exec(cmd)
	if err != nil {
		log.Println("Error pinning window:", out, err)
		return
	}
	fmt.Println(addr, "is pinned?", w.Pinned)
}

func daemonConnect() net.Conn {
	conn, err := net.Dial("unix", socFile)
	if err != nil {
		log.Fatal("Unable to connect to daemon. Is hyprmand running?")
	}
	return conn
}

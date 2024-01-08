package main

import (
	"fmt"
	"log"
	"os/exec"
)

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
	for cmd := range hyprctlCmds {
		s = s + cmd + "\n"
	}
	return s
}

var hyprctlPath = "/usr/bin/hyprctl"

func run(file string, cmd string) error {
	path, err := exec.LookPath(file)
	if err != nil {
		return fmt.Errorf("unable to find executable path", file)
	}
	exec.Command(path, cmd)
	return nil
}

// toggle pin property on window with `addr`
func togglePin(addr string) {
	cmd := fmt.Sprintf(hyprctlCmds["pin"], addr)
	err := run(hyprctlPath, cmd)
	if err != nil {
		log.Fatal("Unable to toggle pin", err)
	}
}

// toggles fullscreen of type t for window at address addr
// t is one of "0" (real fullscreen) or "1" (monocle)
func toggleFull(w *Window, t string) {
	cmd := fmt.Sprintf(hyprctlCmds["fullscreen"], t, w.Address)
	err := run(hyprctlPath, cmd)
	if err != nil {
		log.Fatal("Unable to toggle fullscreen", err)
	}
	w.Fullscreen = !w.Fullscreen
}

// func cmdBytes(c string, v ...string) []byte {
// 	for _, i := range v {
// 		c = strings.Replace(c, `%v`, i, 1)
// 	}
// 	return []byte(c)
// }

type StateManager interface {
	SetActive(addr string)
	GetActive() *Window
}

func Minimize(s StateManager, _ ...string) {
	w := s.GetActive()
	cmd := fmt.Sprintf(hyprctlCmds["minimize"], w.Address)
	err := run("hyprland", cmd)
	if err != nil {
		log.Println("Unable to minimize", err)
		return
	}
}

func Restore(s StateManager, _ ...string) {
	w := s.GetActive()
	cmd := fmt.Sprintf(hyprctlCmds["restore"], fmt.Sprint(w.Workspace), w.Address)
	err := run("hyprland", cmd)
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
	w := s.GetActive()
	if w.Fullscreen {
		toggleFull(w, fscreenType)
		if w.Pinned {
			togglePin(w.Address)
		}
		return
	}
	togglePin(w.Address)
	w.Fullscreen = true

}

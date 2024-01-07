package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

var ctrlCommands map[string]string = map[string]string{
	"active":     "hyprctl dispatch activewindow",
	"minimize":   "hyprctl dispatch movetoworkspacesilent special hyprman,address:%v",
	"restore":    "hyprctl dispatch movetoworkspacesilent %v,address:%v",
	"fullscreen": "hyprctl dispatch fullscreen address:%v",
	"pin":        "hyprctl dispatch pin address:%v",
	"float":      "hyprctl dispatch togglefloating",
}

func cmdBytes(c string, v ...string) []byte {
	for _, i := range v {
		c = strings.Replace(c, `%v`, i, 1)
	}
	return []byte(c)
}

type StateManager interface {
	SetActive(addr string)
	GetActive() *Window
}

func Minimize(s StateManager, c io.Writer, _ ...string) {
	w := s.GetActive()
	cmd := cmdBytes(ctrlCommands["minimize"], w.Address)
	_, err := c.Write(cmd)
	if err != nil {
		log.Println("Unable to minimize", err)
		return
	}
}

func Restore(s StateManager, c io.Writer, _ ...string) {
	w := s.GetActive()
	cmd := cmdBytes(ctrlCommands["restore"], fmt.Sprint(w.Workspace), w.Address)
	_, err := c.Write(cmd)
	if err != nil {
		log.Println("Unable to restore", err)
		return
	}
}

func ToggleFullscreen(s StateManager, c io.Writer, _ ...string) {
	w := s.GetActive()
	if w.Pinned {
	}
	cmd := []byte(fmt.Sprintf(ctrlCommands["fullscreen"], w.Address))
	_, err := c.Write(cmd)
	if err != nil {
		log.Println("Unable to toggle fullscreen", err)
	}
}

package main

import (
	"fmt"
	"log"
	"strings"
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

func StateInit() *State {
	client := NewClient()
	resp, err := client.Exec("dispatch activewindow")
	if err != nil {
		log.Fatal("unable to create new state", err)
	}
    fmt.Println(resp)
	return &State{
		activeWindow: nil,
		windows:      make(map[string]*Window),
		client:       client,
	}
}

type State struct {
	activeWindow *Window
	windows      map[string]*Window
	client       *Client
}

func (s *State) ActiveWindow() *Window {
	return s.activeWindow
}

func (s *State) SetActive(addr string) {}

func (s *State) Client() *Client {
	return s.client
}

// Filters windows in `State` that have value `v` in field `f`
func (s *State) Filter(f string, v any) map[string]*Window {
	windows := make(map[string]*Window)
	for k, w := range s.windows {
		switch strings.ToLower(f) {
		case "address":
			if w.Address == v {
				windows[k] = w
			}
		case "mapped":
			if w.Mapped == v {
				windows[k] = w
			}
		case "hidden":
			if w.Hidden == v {
				windows[k] = w
			}
		case "at":
			if w.At == v {
				windows[k] = w
			}
		case "size":
			if w.Size == v {
				windows[k] = w
			}
		case "workspace":
			if w.Workspace == v {
				windows[k] = w
			}
		case "floating":
			if w.Floating == v {
				windows[k] = w
			}
		case "monitor":
			if w.Monitor == v {
				windows[k] = w
			}
		case "class":
			if w.Class == v {
				windows[k] = w
			}
		case "title":
			if w.Title == v {
				windows[k] = w
			}
		case "initialclass":
			if w.InitialClass == v {
				windows[k] = w
			}
		case "initialtitle":
			if w.InitialTitle == v {
				windows[k] = w
			}
		case "pid":
			if w.Pid == v {
				windows[k] = w
			}
		case "xwayland":
			if w.Xwayland == v {
				windows[k] = w
			}
		case "pinned":
			if w.Pinned == v {
				windows[k] = w
			}
		case "fullscreen":
			if w.Fullscreen == v {
				windows[k] = w
			}
		case "fullscreenmode":
			if w.FullscreenMode == v {
				windows[k] = w
			}
		case "fakefullscreen":
			if w.FakeFullscreen == v {
				windows[k] = w
			}
		case "minimzied":
			if w.Minimzied == v {
				windows[k] = w
			}
		}
	}
	return windows
}
func (s *State) AddWindow(w *Window) {
	s.windows[w.Address] = w
}

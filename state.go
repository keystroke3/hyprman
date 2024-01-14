package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type Workspace struct {
	Id   int
	Name string
}

type Window struct {
	Address        string    `json:"address"`
	Mapped         bool      `json:"mapped"`
	Hidden         bool      `json:"hidden"`
	At             [2]int    `json:"at"`
	Size           [2]int    `json:"size"`
	Workspace      Workspace `json:"workspace"`
	Floating       bool      `json:"floating"`
	Monitor        int       `json:"monitor"`
	Class          string    `json:"class"`
	Title          string    `json:"title"`
	InitialClass   string    `json:"initial_class"`
	InitialTitle   string    `json:"initial_title"`
	Pid            int       `json:"pid"`
	Xwayland       bool      `json:"xwayland"`
	Pinned         bool      `json:"pinned"`
	Fullscreen     bool      `json:"fullscreen"`
	FullscreenMode int       `json:"fullscreen_mode"`
	FakeFullscreen bool      `json:"fake_fullscreen"`
	Minimzied      bool      `json:"minimzied"`
	PreviouslyPinned bool
}

func StateInit() *State {
	// client := NewClient()
	client := &Shell{}
	_, err := client.Exec("hyprctl activewindow")
	if err != nil {
		log.Fatal("unable to create new state: ", err)
	}
	state := &State{
		activeWindow: nil,
		windows:      make(map[string]*Window),
		client:       client,
		pinned:      make(map[string]bool),
	}
	state.SetActive()
	return state
}

type State struct {
	activeWindow *Window
	windows      map[string]*Window
	client       Executor
	pinned       map[string]bool
}

func (s *State) ActiveWindow() *Window {
	return s.activeWindow
}

func (s *State) SavePinned(addr string){
    s.pinned[addr] = true
}

func (s *State) RemovePinned(addr string){
    delete(s.pinned, addr)
}
func (s *State) IsPinned(addr string) bool {
    return s.pinned[addr]
}


func (s *State) SetActive(addr ...string) {
	var w Window
	var a string
	if len(addr) > 0 {
		a = addr[0]
	}
	win, set := s.windows[a]
	if set {
		s.activeWindow = win
		return
	}
	wJson, err := exec.Command("hyprctl", "activewindow", "-j").Output()
	if err != nil {
		log.Println("unable to query active window")
		return
	}
	err = json.Unmarshal(wJson, &w)
	if err != nil {
		log.Println("unable to unmarshal command output:", err)
	}
	_, set = s.windows[a]
	s.AddWindow(&w)
	s.activeWindow = &w
	log.Println("Set active window to", s.activeWindow.Address)
}

func (s *State) Client() Executor {
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
			if w.Workspace.Id == v {
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

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var ctrlSocFile = fmt.Sprintf("%v/hypr/%v/.socket.sock", os.Getenv("XDG_RUNTIME_DIR"), his)

type Executor interface {
	Exec(string) (string, error)
}

func NewClient() *Client {
	conn, err := net.Dial("unix", ctrlSocFile)
	if err != nil {
		log.Fatal("Unable to connect to control socket ", err)
	}
	if err != nil {
		log.Fatal("Unable to accept control socket packets ", err)
	}
	return &Client{Conn: conn}
}

// Provides a connection to a hyprland client
type Client struct {
	Conn net.Conn
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Write(s string) (resp string, err error) {
	writer := bufio.NewWriter(c.Conn)
	_, err = writer.Write([]byte(s))
	if err != nil {
		return
	}
	var res []byte
	reader := bufio.NewReaderSize(c.Conn, 8192)
	var e error
	var buf []byte
	for e != io.EOF {
		if e != bufio.ErrBufferFull {
			err = e
			break
		}
		buf, e = reader.ReadSlice('\n')
		res = append(res, buf...)
	}
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *Client) Exec(s string) (resp string, err error) {
	return c.Write(s)
}

type Shell struct{}

func (s *Shell) Exec(c string) (string, error) {
	args := []string{}
	input := strings.Split(c, " ")
	cmd := input[0]
	if len(input) > 1 {
		args = input[1:]
	}
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

package main

import (
	"fmt"
	//	"io"
	"bytes"
	"container/list"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

type broad struct {
	who     int
	message string
}

type user struct {
	Name  string
	Get   chan string
	Send  chan string
	Quit  chan bool
	conn  net.Conn
	users *list.List
}

func (u *user) Read(buffer []byte) bool {
	_, err := u.conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		u.Close()
		return false
	}
	return true
}

func (u *user) Close() {
	u.Quit <- true
	u.conn.Close()
	u.RemoveMe()
}

func (u *user) Equal(other *user) bool {
	if u.Name == other.Name {
		if u.conn == other.conn {
			return true
		}
	}
	return false
}

func (u *user) RemoveMe() {
	for entry := u.users.Front(); entry != nil; entry = entry.Next() {
		user := entry.Value.(*user)
		if u.Equal(user) {
			u.users.Remove(entry)
		}
	}
}

type chatServer struct {
	send  chan string
	users *list.List
}

var (
	server *chatServer
	count  int
)

func init() {
	server = &chatServer{
		users: list.New(),
		send:  make(chan string),
	}
}

//从用户的get channel中拿到消息并且返回conn
func sender(u *user) {
	for {
		//fmt.Println("sender loop")
		select {
		case buffer := <-u.Get:
			//fmt.Println(u.Name, "get", string(buffer))
			u.conn.Write([]byte(buffer))
		case <-u.Quit:
			u.conn.Close()
		}
	}
	return
}

//从send中拿到消息分发到各个用户的get channel中
func (s *chatServer) loop() {
	for msg := range s.send {
		//fmt.Println("broadcat:", msg)
		for e := s.users.Front(); e != nil; e = e.Next() {
			u := e.Value.(*user)
			//fmt.Println("broadcast to ", u.Name)
			u.Get <- msg
		}
	}
}

//读取conn中的消息并且发到send channel中
func receiver(u *user) {
	buffer := make([]byte, 1024)
	for u.Read(buffer) {
		if bytes.Equal(buffer, []byte("/quit")) {
			u.Close()
			break
		}
		send := u.Name + ">" + string(buffer)
		//fmt.Println("read", send)
		u.Send <- send
		//fmt.Println("after send")
		buffer = make([]byte, 0)
	}
	u.Send <- u.Name + "has left chat"
	return
}

func (s *chatServer) handle(conn net.Conn) {
	//获取姓名
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("connection error", err)
	}
	name := string(buffer[0:n])
	u := &user{name, make(chan string), s.send, make(chan bool), conn, s.users}
	go sender(u)
	go receiver(u)
	s.users.PushBack(u)
	//fmt.Println(name, "join the chat room")
	s.send <- string(name + " has joined the chat")
}

func main() {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	go server.loop()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.handle(c)
	}
}

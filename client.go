package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strings"

	"github.com/tncardoso/gocurses"
)

var (
	windFriends *gocurses.Window
	windRoom    *gocurses.Window
	windCommand *gocurses.Window
	windChat    *gocurses.Window

	users  *list.List
	conn   net.Conn
	readCh chan string
	sendCh chan string
)

const address = "127.0.0.1:4000"

func clear(win *gocurses.Window) {

	row, col := win.Getmaxyx()
	for i := 0; i < row; i++ {
		win.Mvaddstr(i, 0, strings.Repeat(" ", col))
	}

	win.Box(0, 0)
	//	win.Refresh()
}
func init() {
	gocurses.Initscr()
	gocurses.Cbreak()
	//gocurses.Noecho()
	gocurses.Stdscr.Keypad(true)
	gocurses.Refresh()
	row, col := gocurses.Getmaxyx()
	windFriends = gocurses.NewWindow(row, col/4, 0, 0)
	windFriends.Box(0, 0)
	windFriends.Mvaddstr(1, 1, "OnLine:")
	windFriends.Refresh()
	windRoom = gocurses.NewWindow(row/2, col*3/4, 0, col/4)
	windRoom.Box(0, 0)
	windRoom.Refresh()
	windCommand = gocurses.NewWindow(row/4, col*3/4, row/2, col/4)
	windCommand.Box(0, 0)
	windCommand.Refresh()
	windChat = gocurses.NewWindow(row/4, col*3/4, row*3/4, col/4)
	windChat.Box(0, 0)
	windChat.Refresh()

	var err error
	conn, err = net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	readCh = make(chan string, 1)
	sendCh = make(chan string, 1)

	users = list.New()
}

func getChatLine() func() int {
	chatline := -1
	return func() int {
		if chatline > 10 {
			chatline = 0
			clear(windRoom)
		} else {
			chatline++
		}
		return chatline
	}
}

func readString(win *gocurses.Window) string {
	var buffer []byte
	for c := win.Getch(); c != int('\n'); c = win.Getch() {
		buffer = append(buffer, byte(c))
	}
	if len(buffer) > 0 {
		return string(buffer)
	} else {
		return ""
	}
}

func updateOnline(users *list.List, win *gocurses.Window) {
	clear(win)
	i := 2
	win.Mvaddstr(1, 1, "Online:\n")
	for e := users.Front(); e != nil; e = e.Next() {
		win.Mvaddstr(i, 0, e.Value.(string))
		i++
	}
	win.Refresh()
}
func readRoutine(ch chan<- string) {
	reader := bufio.NewReader(conn)
	cl := getChatLine()
	for {
		res, err := reader.ReadString('\n')
		c := cl()
		index := strings.Index(res, ":")
		windCommand.Addstr(res)

		cmd := res[:index]
		windCommand.Addstr(cmd)
		switch cmd {
		case "join":
			users.PushBack(res[index+1:])
			updateOnline(users, windFriends)
			if err != nil {
				windRoom.Mvaddstr(c, 1, fmt.Sprintf("%s %d\n", err.Error(), c))
			} else {
				windRoom.Mvaddstr(c, 1, fmt.Sprintf("%s %d\n", res, c))
			}
			windRoom.Refresh()
		case "leave":
			for e := users.Front(); e != nil; e = e.Next() {
				name := e.Value.(string)
				windCommand.Addstr("\n\nLEAVE:-----" + name + "------\n")
				if name == res[index+1:] {
					windCommand.Addstr("\n\nLEAVE:" + name + "\n\n")
					users.Remove(e)
				}
			}
			updateOnline(users, windFriends)
			if err != nil {
				windRoom.Mvaddstr(c, 1, fmt.Sprintf("%s %d\n", err.Error(), c))
			} else {
				windRoom.Mvaddstr(c, 1, fmt.Sprintf("%s %d\n", res, c))
			}
			windRoom.Refresh()
		case "users":
			uss := strings.Split(res[index+1:], ";")
			for _, u := range uss {
				users.PushBack(u + "\n")
			}
			updateOnline(users, windFriends)
		default:
			if err != nil {
				windRoom.Mvaddstr(c, 1, fmt.Sprintf("%s %d\n", err.Error(), c))
			} else {
				windRoom.Mvaddstr(c, 1, fmt.Sprintf("%s %d\n", res, c))
			}
			windRoom.Refresh()
		}
	}
}

func sendRoutine(ch <-chan string) {
	for send := range ch {
		windCommand.Addstr(send)
		windCommand.Refresh()
		fmt.Fprint(conn, send)
	}
}

func main() {
	go readRoutine(readCh)
	go sendRoutine(sendCh)
	for {
		input := readString(windChat)
		if input == "" {
			continue
		} else {
			if input == "quit" {
				gocurses.End()
			} else {
				sendCh <- input
			}
		}
	}
}

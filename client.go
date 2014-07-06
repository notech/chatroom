package main

import (
	"container/list"
	"fmt"

	"github.com/tncardoso/gocurses"
)

var (
	windFriends *gocurses.Window
	windRoom    *gocurses.Window
	windCommand *gocurses.Window
	windChat    *gocurses.Window

	users *list.List
)

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

func main() {

	for {
		input := readString(windChat)
		windRoom.Addstr("input is " + input)
		if input == "" {
			continue
		} else {
			windRoom.Mvaddstr(1, 1, fmt.Sprintf("%s\n", input))
			windRoom.Refresh()
		}
	}
	gocurses.Getch()
	gocurses.End()
}

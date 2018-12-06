package main

import (
	. "github.com/lxn/walk/declarative"
	"strconv"
)

type TestItem struct {
	Name  string
	Value int
}

func getTestList() []TestItem {
	var dst = []int{110, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 38400, 43000, 56000, 57600, 115200, 128000, 256000}
	portList := make([]TestItem, len(dst))
	for i, v := range dst {
		portList[i] = TestItem{strconv.Itoa(v), v}
	}
	return portList
}
func main() {
	MainWindow{
		Title:   "Label Example",
		MinSize: Size{320, 240},
		Size:    Size{320, 240},
		Layout:  VBox{},
		Children: []Widget{
			ComboBox{
				Name:          "Name",
				MaxSize:       Size{Width: 80, Height: 0},
				MinSize:       Size{Width: 80, Height: 0},
				BindingMember: "Value",
				CurrentIndex:  0,
				DisplayMember: "Name",
				Model:         getTestList(),
			},
			TextEdit{
				Name: "txtInput",
				Text: "hello",
			},
		},
	}.Run()
}

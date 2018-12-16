package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	t := &walk.TextEdit{}
	MainWindow{
		Title:   "Label Test",
		MinSize: Size{Width: 300, Height: 200},
		Layout:  VBox{},
		Children: []Widget{
			TextEdit{},
			TextEdit{
				Name:       "txtSerialSend",
				AssignTo:   &t,
				Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
				//Font:     Font{Family: "Consolas", PointSize: 10},
				ReadOnly: true,
				VScroll:  true,
			},
			PushButton{
				Text: "发送",
				OnClicked: func() {
					t.AppendText("我有一只小毛驴")
				},
			},
		},
	}.Run()
}

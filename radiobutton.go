package main

import (
	. "github.com/lxn/walk/declarative"
)

func main() {
	MainWindow{
		Title:   "Label Example",
		MinSize: Size{320, 240},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: Bind("lbSync3.Text"),
			},
			Label{
				Name: "lbSync3",
				Text: Bind("lbSync2.Text"),
			},
			Label{
				Name: "lbSync2",
				Text: Bind("lbSync1.Text"),
			},
			Label{
				Name: "lbSync1",
				Text: Bind("txtInput.Text"),
			},
			TextEdit{
				Name: "txtInput",
				Text: "hello",
			},
		},
	}.Run()
}

package main

import (
	. "github.com/lxn/walk/declarative"
)

func main() {

	MainWindow{
		Title:   "Label Test",
		MinSize: Size{Width: 300, Height: 200},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: "Default",
			},
			Label{
				Text:          "AlignNear",
				TextAlignment: AlignNear,
			},
			Label{
				Text:          "AlignCenter",
				TextAlignment: AlignCenter,
			},
			Label{
				Text:          "AlignFar",
				TextAlignment: AlignFar,
			},
		},
	}.Run()
}

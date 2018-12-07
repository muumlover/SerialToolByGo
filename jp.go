package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"log"
)

type MyLoadWindow struct {
	*walk.MainWindow
	progressBar *walk.ProgressBar
}

func main() {
	mw := &MyLoadWindow{}

	// 画面情報設定
	if err := (MainWindow{
		AssignTo: &mw.MainWindow, // Widgetを実体に割り当て
		Title:    "コンピュータの情報を取得中",
		Size:     Size{Width: 300, Height: 100},
		Visible:  false,
		Font:     Font{PointSize: 12},
		Layout:   VBox{},

		Children: []Widget{ // ウィジェットを入れるスライス

			ProgressBar{
				AssignTo:    &mw.progressBar,
				MarqueeMode: true,
			},
		},
	}).Create(); err != nil {
		log.Print(err)
	}

	screenX := int(win.GetSystemMetrics(win.SM_CXSCREEN))
	screenY := int(win.GetSystemMetrics(win.SM_CYSCREEN))
	if err := mw.MainWindow.SetBounds(walk.Rectangle{
		X:      (screenX - mw.MainWindow.Width()) / 2,
		Y:      (screenY - mw.MainWindow.Height()) / 2,
		Width:  mw.MainWindow.Width(),
		Height: mw.MainWindow.Height(),
	}); err != nil {
		log.Print(err)
		return
	}

	mw.MainWindow.SetVisible(true)
	mw.MainWindow.Run()

}

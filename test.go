package main

import (
	"github.com/muumlover/walk"
	. "github.com/muumlover/walk/declarative"
	"log"
)

func main() {
	if _, err := (MainWindow{
		Title:   "Label Test",
		MinSize: Size{Width: 300, Height: 200},
		Layout:  HBox{},
		Children: []Widget{
			Composite{
				Layout:  VBox{MarginsZero: true},
				MaxSize: Size{Width: 160, Height: 160},
				MinSize: Size{Width: 160, Height: 160},
				Children: []Widget{
					GroupBox{
						Title:  "设置",
						Layout: VBox{},
						Children: []Widget{
							Composite{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text:          "端口号",
										TextAlignment: AlignNear,
									},
									ComboBox{
										Enabled: Bind("SerialState.Text=='OFF'"),
										Name:    "Name",
										MaxSize: Size{Width: 80, Height: 0},
										MinSize: Size{Width: 80, Height: 0},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "波特率",
									},
									ComboBox{
										Enabled: Bind("SerialState.Text=='OFF'"),
										Name:    "Baud",
										MaxSize: Size{Width: 80, Height: 0},
										MinSize: Size{Width: 80, Height: 0},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "数据位",
									},
									ComboBox{
										Enabled: Bind("SerialState.Text=='OFF'"),
										Name:    "Size",
										MaxSize: Size{Width: 80, Height: 0},
										MinSize: Size{Width: 80, Height: 0},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "校验位",
									},
									ComboBox{
										Enabled: Bind("SerialState.Text=='OFF'"),
										Name:    "Parity",
										MaxSize: Size{Width: 80, Height: 0},
										MinSize: Size{Width: 80, Height: 0},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "停止位",
									},
									ComboBox{
										Enabled: Bind("SerialState.Text=='OFF'"),
										Name:    "StopBits",
										MaxSize: Size{Width: 80, Height: 0},
										MinSize: Size{Width: 80, Height: 0},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Name: "SerialOpen",
										Text: Bind("SerialState.Text=='OFF'?'打开串口':'关闭串口'"),
										OnClicked: func() {
										},
									},
									Label{
										Name:          "SerialState",
										Background:    SolidColorBrush{Color: walk.RGB(255, 0, 0)},
										MaxSize:       Size{Width: 50, Height: 0},
										MinSize:       Size{Width: 50, Height: 0},
										TextAlignment: AlignCenter,
										TextColor:     walk.RGB(0, 0, 0),
										Text:          "OFF",
									},
								},
							},
						},
					},
					VSpacer{},
				},
			}, Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TextEdit{
						Name:       "txtSerialRecv",
						Background: SolidColorBrush{Color: walk.RGB(0, 0, 0)},
						TextColor:  walk.RGB(0, 255, 0),
						//Enabled:    false,
						//Font:     Font{Family: "Courier New", PointSize: 10},
						Font:     Font{Family: "Consolas", PointSize: 10},
						ReadOnly: true,
						VScroll:  true,
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							Label{
								Text: "编码方式",
							},
							ComboBox{
								CurrentIndex: 0,
								Model:        []string{"gbk", "utf8"},
							},
							CheckBox{
								Text: "定时发送",
							},
							HSpacer{},
						},
					}, Composite{
						Layout:  HBox{MarginsZero: true},
						MaxSize: Size{Width: 100, Height: 100},
						MinSize: Size{Width: 100, Height: 100},
						Children: []Widget{
							TextEdit{
								Name: "txtSerialSend",
								Font: Font{Family: "Consolas", PointSize: 10},
							},
							Composite{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text: "发送",
										OnClicked: func() {
										},
									},
								},
							},
							//HSpacer{},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							CheckBox{
								Text: "定时发送",
							},
							Label{
								Text: "周期",
							},
							NumberEdit{},
							Label{
								Text: "ms",
							},
							HSpacer{},
						},
					},
					//VSpacer{},
				},
			},
			//HSpacer{},
		},
	}).Run(); err != nil {
		log.Print(err)
	}
}

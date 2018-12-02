package main

import (
	"github.com/StackExchange/wmi"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
)

type ComboBoxItem struct {
	Name  string
	Value string
}
type ComboBoxItemInt struct {
	Name  string
	Value int
}
type ComboBoxItemByte struct {
	Name  string
	Value byte
}
type ComboBoxItemParity struct {
	Name  string
	Value serial.Parity
}
type ComboBoxItemStopBits struct {
	Name  string
	Value serial.StopBits
}

func getPortList() []ComboBoxItem {
	type Win32_PnPEntity struct {
		Name string
	}
	var dst []Win32_PnPEntity
	q := wmi.CreateQuery(&dst, "WHERE Name LIKE '%(COM%)'")
	err := wmi.Query(q, &dst)
	if err != nil {
		log.Fatal(err)
	}
	portList := make([]ComboBoxItem, len(dst))
	for i, v := range dst {
		portList[i] = ComboBoxItem{v.Name, strings.Split(strings.Split(v.Name, "(")[1], ")")[0]}
	}
	return portList
}
func getBaudRate() []ComboBoxItemInt {
	var dst = []int{110, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 38400, 43000, 56000, 57600, 115200, 128000, 256000}
	portList := make([]ComboBoxItemInt, len(dst))
	for i, v := range dst {
		portList[i] = ComboBoxItemInt{strconv.Itoa(v), v}
	}
	return portList
}
func getDataBits() []ComboBoxItemByte {
	var dst = []byte{8, 7, 6, 5}
	portList := make([]ComboBoxItemByte, len(dst))
	for i, v := range dst {
		portList[i] = ComboBoxItemByte{strconv.Itoa(int(v)), v}
	}
	return portList
}
func getParity() []ComboBoxItemParity {
	return []ComboBoxItemParity{
		{"None", serial.ParityNone},
		{"Odd", serial.ParityOdd},
		{"Even", serial.ParityEven},
		{"Mark", serial.ParityMark},
		{"Space", serial.ParitySpace},
	}
}
func getStopBits() []ComboBoxItemStopBits {
	return []ComboBoxItemStopBits{
		{"1", serial.Stop1},
		{"1.5", serial.Stop1Half},
		{"2", serial.Stop2},
	}
}

type SerialConfigItems struct {
	PortList []ComboBoxItem
	BaudRate []ComboBoxItemInt
	DataBits []ComboBoxItemByte
	Parity   []ComboBoxItemParity
	StopBits []ComboBoxItemStopBits
}

func getSerialConfigItems() SerialConfigItems {
	return SerialConfigItems{
		getPortList(),
		getBaudRate(),
		getDataBits(),
		getParity(),
		getStopBits(),
	}
}

type mwMainWindow struct {
	*walk.MainWindow
	cbSerialPort,
	cbBaudRate,
	cbDataBits,
	cbParity,
	cbStopBits *walk.ComboBox

	scItems SerialConfigItems
	sc      *serial.Config
}

func (mw *mwMainWindow) openSerial() {
	println(mw.sc.Name)
	println(mw.sc.Baud)
	println(mw.sc.Size)
	println(mw.sc.Parity)
	println(mw.sc.StopBits)
	print(mw.cbSerialPort.BindingMember())
	//spc := new(serial.Config)
	s, err := serial.OpenPort(mw.sc)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}

	go func(s *serial.Port) {
		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("%q", buf[:n])
	}(s)
}

func main() {
	var db *walk.DataBinder
	//serialConfig := new(serial.Config)

	mw := &mwMainWindow{}
	mw.scItems = getSerialConfigItems()
	mw.sc = new(serial.Config)

	partLeftTop := GroupBox{
		Title:  "设置",
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "端口号",
					},
					ComboBox{
						AssignTo:      &mw.cbSerialPort,
						MaxSize:       Size{Width: 100, Height: 0},
						MinSize:       Size{Width: 100, Height: 0},
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         mw.scItems.PortList,
						Value:         Bind("Name"),
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
						AssignTo:      &mw.cbBaudRate,
						MaxSize:       Size{Width: 100, Height: 0},
						MinSize:       Size{Width: 100, Height: 0},
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         mw.scItems.BaudRate,
						Value:         Bind("Baud"),
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
						AssignTo:      &mw.cbDataBits,
						MaxSize:       Size{Width: 100, Height: 0},
						MinSize:       Size{Width: 100, Height: 0},
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         mw.scItems.DataBits,
						Value:         Bind("Size"),
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
						AssignTo:      &mw.cbStopBits,
						MaxSize:       Size{Width: 100, Height: 0},
						MinSize:       Size{Width: 100, Height: 0},
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         mw.scItems.StopBits,
						Value:         Bind("StopBits"),
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
						AssignTo:      &mw.cbParity,
						MaxSize:       Size{Width: 100, Height: 0},
						MinSize:       Size{Width: 100, Height: 0},
						BindingMember: "Value",
						DisplayMember: "Name",
						Model:         mw.scItems.Parity,
						Value:         Bind("Parity"),
					},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text: "打开串口",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							mw.openSerial()
						},
					},
					Label{
						Text: "OFF",
					},
				},
			},
			VSpacer{},
		},
	}

	partLeft := Composite{
		Layout:  VBox{MarginsZero: true},
		MaxSize: Size{Width: 180, Height: 10},
		MinSize: Size{Width: 180, Height: 10},
		Children: []Widget{
			partLeftTop,
			VSpacer{},
		},
	}

	partRight := Composite{
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			TextEdit{},
			VSpacer{},
			TextEdit{},
		},
	}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		//Title:    "SerialTool By Golang",
		Title: Bind("'Animal Details' + (sc.Name == '' ? '' : ' - ' + sc.Name)"),
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "sc",
			DataSource:     mw.sc,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{Width: 600, Height: 400},
		Layout:  HBox{},
		Children: []Widget{
			partLeft,
			HSpacer{},
			partRight,
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}
	mw.cbSerialPort.SetCurrentIndex(0)
	mw.cbBaudRate.SetCurrentIndex(6)
	mw.cbDataBits.SetCurrentIndex(0)
	mw.cbParity.SetCurrentIndex(0)
	mw.cbStopBits.SetCurrentIndex(0)
	mw.MainWindow.Run()

	/*
		MainWindow{
			AssignTo: &mw.MainWindow,
			Title:    "SerialTool By Golang",
			MinSize:  Size{600, 400},
			Layout:   VBox{},
			Children: []Widget{
				Composite{
					Layout: HBox{MarginsZero: true},
					Children: []Widget{
						Composite{
							Layout: VBox{MarginsZero: true},
							Children: []Widget{
								Label{
									Text: "SerialPort:",
								},
								ComboBox{
									AssignTo:      &mw.cbSerialPort,
									BindingMember: "Value",
									DisplayMember: "Name",
									Model:         getPortList(),
								},
								Composite{
									Layout: HBox{MarginsZero: true},
									Children: []Widget{
										Label{
											Text: "BaudRate:",
										},
										ComboBox{
											AssignTo:      &mw.cbBaudRate,
											BindingMember: "Value",
											DisplayMember: "Name",
											Model:         getBaudRate(),
										},
									},
								},
								Composite{
									Layout: HBox{MarginsZero: true},
									Children: []Widget{
										Label{
											Text: "DataBits:",
										},
										ComboBox{
											AssignTo:      &mw.cbDataBits,
											BindingMember: "Value",
											DisplayMember: "Name",
											Model:         getDataBits(),
										},
									},
								},
								Composite{
									Layout: HBox{MarginsZero: true},
									Children: []Widget{
										Label{
											Text: "Parity:  ",
										},
										ComboBox{
											AssignTo:      &mw.cbParity,
											BindingMember: "Value",
											DisplayMember: "Name",
											Model:         getParity(),
										},
									},
								},
								Composite{
									Layout: HBox{MarginsZero: true},
									Children: []Widget{
										Label{
											Text: "StopBits:",
										},
										ComboBox{
											AssignTo:      &mw.cbStopBits,
											BindingMember: "Value",
											DisplayMember: "Name",
											Model:         getStopBits(),
										},
									},
								},
								Composite{
									Layout: HBox{MarginsZero: true},
									Children: []Widget{
										Label{
											Text: "SerialOperate:",
										},
										PushButton{
											Text: "OpenSerial",
										},
									},
								},
								VSpacer{},
							},
						},
						HSpacer{},
						TextEdit{},
					},
				},
				VSpacer{},
				TextEdit{},
			},
		}.Run()
	*/
}

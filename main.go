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

type PortNameItem struct {
	Name  string
	Value string
}

func getPortNameList() []PortNameItem {
	// noinspection ALL
	type Win32_PnPEntity struct {
		Name string
	}
	var dst []Win32_PnPEntity
	q := wmi.CreateQuery(&dst, "WHERE Name LIKE '%(COM%)'")
	err := wmi.Query(q, &dst)
	if err != nil {
		log.Fatal(err)
	}
	portList := make([]PortNameItem, len(dst))
	for i, v := range dst {
		portList[i] = PortNameItem{v.Name, strings.Split(strings.Split(v.Name, "(")[1], ")")[0]}
	}
	return portList
}

type BaudRateItem struct {
	Name  string
	Value int
}

func getBaudRateList() []BaudRateItem {
	var dst = []int{110, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 38400, 43000, 56000, 57600, 115200, 128000, 256000}
	portList := make([]BaudRateItem, len(dst))
	for i, v := range dst {
		portList[i] = BaudRateItem{strconv.Itoa(v), v}
	}
	return portList
}

type DataBitsItem struct {
	Name  string
	Value byte
}

func getDataBitsList() []DataBitsItem {
	var dst = []byte{8, 7, 6, 5}
	portList := make([]DataBitsItem, len(dst))
	for i, v := range dst {
		portList[i] = DataBitsItem{strconv.Itoa(int(v)), v}
	}
	return portList
}

type ParityItem struct {
	Name  string
	Value serial.Parity
}

func getParityList() []ParityItem {
	return []ParityItem{
		{"None", serial.ParityNone},
		{"Odd", serial.ParityOdd},
		{"Even", serial.ParityEven},
		{"Mark", serial.ParityMark},
		{"Space", serial.ParitySpace},
	}
}

type StopBitItem struct {
	Name  string
	Value serial.StopBits
}

func getStopBitsList() []StopBitItem {
	return []StopBitItem{
		{"1", serial.Stop1},
		{"1.5", serial.Stop1Half},
		{"2", serial.Stop2},
	}
}

type SerialConfigItems struct {
	PortList []PortNameItem
	BaudRate []BaudRateItem
	DataBits []DataBitsItem
	Parity   []ParityItem
	StopBits []StopBitItem
}

func getSerialConfigItems() SerialConfigItems {
	return SerialConfigItems{
		getPortNameList(),
		getBaudRateList(),
		getDataBitsList(),
		getParityList(),
		getStopBitsList(),
	}
}

type mwMainWindow struct {
	*walk.MainWindow
	cbSerialPort   *walk.ComboBox
	cbBaudRate     *walk.ComboBox
	cbDataBits     *walk.ComboBox
	cbParity       *walk.ComboBox
	cbStopBits     *walk.ComboBox
	btnSerialOpen  *walk.PushButton
	txtSerialState *walk.Label

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
		// noinspection ALL
		log.Print("%q", buf[:n])
	}(s)
}

func main() {
	mw := mwMainWindow{}
	mw.scItems = getSerialConfigItems()
	mw.sc = new(serial.Config)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    Bind("'Animal Details' + (sc.Name == '' ? '' : ' - ' + sc.Name)"),
		MinSize:  Size{Width: 600, Height: 400},
		Layout:   HBox{},
		Children: []Widget{
			Composite{
				Layout:  VBox{MarginsZero: true},
				MaxSize: Size{Width: 150, Height: 10},
				MinSize: Size{Width: 150, Height: 10},
				Children: []Widget{
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							GroupBox{
								Title:  "设置",
								Layout: VBox{},
								Children: []Widget{
									Composite{
										Layout: VBox{MarginsZero: true},
										Children: []Widget{
											Label{
												Text: "端口号",
											},
											ComboBox{
												Enabled:       Bind("SerialState.Text=='OFF'"),
												Name:          "Name",
												MaxSize:       Size{Width: 80, Height: 0},
												MinSize:       Size{Width: 80, Height: 0},
												AssignTo:      &mw.cbSerialPort,
												BindingMember: "Value",
												DisplayMember: "Name",
												Model:         mw.scItems.PortList,
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
												Enabled:       Bind("SerialState.Text=='OFF'"),
												Name:          "Baud",
												MaxSize:       Size{Width: 80, Height: 0},
												MinSize:       Size{Width: 80, Height: 0},
												AssignTo:      &mw.cbBaudRate,
												BindingMember: "Value",
												DisplayMember: "Name",
												Model:         mw.scItems.BaudRate,
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
												Enabled:       Bind("SerialState.Text=='OFF'"),
												Name:          "Size",
												MaxSize:       Size{Width: 80, Height: 0},
												MinSize:       Size{Width: 80, Height: 0},
												AssignTo:      &mw.cbDataBits,
												BindingMember: "Value",
												DisplayMember: "Name",
												Model:         mw.scItems.DataBits,
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
												Enabled:       Bind("SerialState.Text=='OFF'"),
												Name:          "StopBits",
												MaxSize:       Size{Width: 80, Height: 0},
												MinSize:       Size{Width: 80, Height: 0},
												AssignTo:      &mw.cbStopBits,
												BindingMember: "Value",
												DisplayMember: "Name",
												Model:         mw.scItems.StopBits,
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
												Enabled:       Bind("SerialState.Text=='OFF'"),
												Name:          "Parity",
												MaxSize:       Size{Width: 80, Height: 0},
												MinSize:       Size{Width: 80, Height: 0},
												AssignTo:      &mw.cbParity,
												BindingMember: "Value",
												DisplayMember: "Name",
												Model:         mw.scItems.Parity,
											},
										},
									},
									Composite{
										Layout: HBox{MarginsZero: true},
										Children: []Widget{
											PushButton{
												AssignTo: &mw.btnSerialOpen,
												Name:     "SerialOpen",
												Text:     Bind("SerialState.Text=='OFF'?'打开串口':'关闭串口'"),
												OnClicked: func() {
													if mw.btnSerialOpen.Text() == "关闭串口" {
														bg, err := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
														if err != nil {
															log.Print(err)
														}
														mw.txtSerialState.SetBackground(bg)
														if err := mw.txtSerialState.SetText("OFF"); err != nil {
															log.Print(err)
															return
														}
													} else {
														bg, err := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
														if err != nil {
															log.Print(err)
														}
														mw.txtSerialState.SetBackground(bg)
														if err := mw.txtSerialState.SetText("ON"); err != nil {
															log.Print(err)
															return
														}
													}
												},
											},
											HSpacer{},
											Label{
												AssignTo:   &mw.txtSerialState,
												Name:       "SerialState",
												Background: SolidColorBrush{Color: walk.RGB(255, 0, 0)},
												MaxSize:    Size{Width: 50, Height: 0},
												MinSize:    Size{Width: 50, Height: 0},
												Alignment:  AlignCenter,
												TextColor:  walk.RGB(0, 0, 0),
												Text:       "OFF",
											},
										},
									},
								},
							},
							Label{},
							VSpacer{},
						},
					},
					VSpacer{},
				},
			},
			HSpacer{},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TextEdit{
						Name: "textRecv",
						Text: Bind("textSend.Text"),
					},
					VSpacer{},
					TextEdit{
						Name: "textSend",
						Text: "",
					},
				},
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"rgb": func(args ...interface{}) (interface{}, error) {
				return walk.RGB(byte(args[0].(float64)), byte(args[1].(float64)), byte(args[2].(float64))), nil
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}
	if err := mw.cbSerialPort.SetCurrentIndex(0); err != nil {
		log.Print(err)
		return
	}
	if err := mw.cbBaudRate.SetCurrentIndex(6); err != nil {
		log.Print(err)
		return
	}
	if err := mw.cbDataBits.SetCurrentIndex(0); err != nil {
		log.Print(err)
		return
	}
	if err := mw.cbParity.SetCurrentIndex(0); err != nil {
		log.Print(err)
		return
	}
	if err := mw.cbStopBits.SetCurrentIndex(0); err != nil {
		log.Print(err)
		return
	}
	mw.MainWindow.Run()
}

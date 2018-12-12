package main

import (
	"github.com/StackExchange/wmi"
	"github.com/axgle/mahonia"
	"github.com/lxn/win"
	"github.com/muumlover/walk"
	. "github.com/muumlover/walk/declarative"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
	"time"
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
		log.Print(err)
	}
	portList := make([]PortNameItem, len(dst))
	for i, v := range dst {
		name := strings.Split(v.Name, "(")[0]
		port := strings.Split(strings.Split(v.Name, "(")[1], ")")[0]
		portList[i] = PortNameItem{port + " : " + name, port}
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

func UseNewEncoder(src string, oldEncoder string, newEncoder string) string {
	srcDecoder := mahonia.NewDecoder(oldEncoder)
	desDecoder := mahonia.NewDecoder(newEncoder)
	resStr := srcDecoder.ConvertString(src)
	_, resBytes, _ := desDecoder.Translate([]byte(resStr), true)
	return string(resBytes)
}

func (mw *mwMainWindow) openSerial() error {
	sc := serial.Config{}
	sc.Name = mw.cbSerialPort.Value().(string)
	sc.Baud = mw.cbBaudRate.Value().(int)
	sc.Size = mw.cbDataBits.Value().(byte)
	sc.Parity = mw.cbParity.Value().(serial.Parity)
	sc.StopBits = mw.cbStopBits.Value().(serial.StopBits)
	sc.ReadTimeout = time.Millisecond
	s, err := serial.OpenPort(&sc)
	if err != nil {
		return err
	}
	mw.sp = s
	n, err := s.Write([]byte("test"))
	if err != nil {
		return err
	}

	go func(s *serial.Port) {
		isHalf := false
		half := byte(0)
		for {
			buf := make([]byte, 10240)
			n, err = s.Read(buf)
			if err != nil {
				log.Print(err)
				break
			}
			if n == 0 && !isHalf {
				continue
			}
			if isHalf {
				buf = append([]byte{half}, buf[:n]...)
				n += 1
				isHalf = false
			} else if n == 1 {
				half = buf[:n][0]
				isHalf = true
				continue
			}
			// noinspection ALL0
			log.Print(buf[:n])
			//str := string(buf[:n])
			//log.Print(str)
			decoder := mahonia.NewDecoder("gbk")
			_, cdata, _ := decoder.Translate(buf[:n], true)
			str := string(cdata)
			log.Print(str)
			str = strings.Replace(str, "\x00", "", -1)
			mw.txtSerialRecv.SetSuspended(true)
			mw.txtSerialRecv.AppendText(str)
			mw.txtSerialRecv.SetSuspended(false)
			if err != nil {
				log.Print(err)
				break
			}
		}
	}(s)

	return nil
}

func (mw *mwMainWindow) closeSerial() error {
	if err := mw.sp.Close(); err != nil {
		return err
	}
	return nil
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
	txtSerialRecv  *walk.TextEdit
	txtSerialSend  *walk.TextEdit

	scItems SerialConfigItems
	sc      *serial.Config
	sp      *serial.Port
}

func main() {
	mw := mwMainWindow{}
	mw.scItems = getSerialConfigItems()
	mw.sc = new(serial.Config)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Icon:     "favicon.ico",
		Title:    "SerialTool",
		MinSize:  Size{Width: 600, Height: 400},
		Size:     Size{Width: 600, Height: 400},
		Visible:  false,
		Layout:   HBox{},
		Children: []Widget{
			Composite{
				Layout:  VBox{MarginsZero: true},
				MaxSize: Size{Width: 160, Height: 10},
				MinSize: Size{Width: 160, Height: 10},
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
												Text:          "端口号",
												TextAlignment: AlignNear,
											},
											ComboBox{
												Enabled:       Bind("SerialState.Text=='OFF'"),
												Name:          "Name",
												MaxSize:       Size{Width: 80, Height: 0},
												MinSize:       Size{Width: 80, Height: 0},
												AssignTo:      &mw.cbSerialPort,
												BindingMember: "Value",
												CurrentIndex:  0,
												DisplayMember: "Name",
												Model:         mw.scItems.PortList,
												OnDropDown: func() {
													go func() {
														if err := mw.cbSerialPort.SetModel(getPortNameList()); err != nil {
															log.Print(err)
															return
														}
														//Todo 列表减少时刷新显示
													}()
												},
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
												CurrentIndex:  6,
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
												CurrentIndex:  0,
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
												CurrentIndex:  0,
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
												CurrentIndex:  0,
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
													if mw.txtSerialState.Text() == "ON" {
														if err := mw.closeSerial(); err != nil {
															log.Print(err)
															return
														}
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
														if err := mw.openSerial(); err != nil {
															log.Print(err)
															return
														}
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
												AssignTo:      &mw.txtSerialState,
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
					},
					VSpacer{},
				},
			},
			HSpacer{},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TextEdit{
						Name:       "txtSerialRecv",
						AssignTo:   &mw.txtSerialRecv,
						Background: SolidColorBrush{Color: walk.RGB(0, 0, 0)},
						TextColor:  walk.RGB(0, 255, 0),
						//Enabled:    false,
						//Font:     Font{Family: "Courier New", PointSize: 10},
						Font:     Font{Family: "Consolas", PointSize: 10},
						ReadOnly: true,
						VScroll:  true,
					},
					VSpacer{},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							TextEdit{
								Name:     "txtSerialSend",
								AssignTo: &mw.txtSerialSend,
								Font:     Font{Family: "Consolas", PointSize: 10},
							},
							HSpacer{},
							PushButton{
								Text: "SCREAM",
								OnClicked: func() {
									_, err := mw.sp.Write([]byte(mw.txtSerialSend.Text()))
									if err != nil {
										log.Print(err)
									}
								},
							},
						},
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
		log.Print(err)
	}
	screenX := int(win.GetSystemMetrics(win.SM_CXSCREEN))
	screenY := int(win.GetSystemMetrics(win.SM_CYSCREEN))
	if err := mw.MainWindow.SetBounds(walk.Rectangle{
		X:      (screenX - mw.Width()) / 2,
		Y:      (screenY - mw.Height()) / 2,
		Width:  mw.Width(),
		Height: mw.Height(),
	}); err != nil {
		log.Print(err)
		return
	}
	mw.MainWindow.SetVisible(true)
	mw.MainWindow.Run()
}

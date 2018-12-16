package main

import (
	"bytes"
	"github.com/StackExchange/wmi"
	"github.com/axgle/mahonia"
	"github.com/lxn/win"
	"github.com/muumlover/walk"
	. "github.com/muumlover/walk/declarative"
	"github.com/tarm/serial"
	"log"
	"math/rand"
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

func (mw *myWindow) openSerial() error {
	sc := serial.Config{}
	sc.Name = mw.cbSerialPort.Value().(string)
	sc.Baud = mw.cbBaudRate.Value().(int)
	sc.Size = mw.cbDataBits.Value().(byte)
	sc.Parity = mw.cbParity.Value().(serial.Parity)
	sc.StopBits = mw.cbStopBits.Value().(serial.StopBits)
	sc.ReadTimeout = time.Millisecond
	p, err := serial.OpenPort(&sc)
	if err != nil {
		return err
	}
	msp.Port = p
	var b bytes.Buffer
	go func() {
		for {
			buf := make([]byte, 1)
			n, err := msp.Read(buf)
			if err != nil {
				log.Print(err)
				break
			}
			b.Write(buf[:n])
			if err != nil {
				log.Print(err)
				break
			}
		}
	}()
	go func() {
		h := false
		ht := false
		d := rand.Float32()
		//c := make([]byte, 6)
		//c := bytes.Buffer{}
		for {
			if b.Len() > 0 {
				if msp.dataRecvEncoding == "gbk" {
					o := b.Bytes()[0]
					if o < 0x80 { //单字节
						d = rand.Float32()
						str := string(b.Next(1))
						//str = strings.Replace(str, "\x00", "", -1)
						//mw.txtSerialRecv.SetSuspended(true)
						mw.txtSerialRecv.AppendText(str)
						//mw.txtSerialRecv.SetSuspended(false)
					} else {
						//cb = false
						//c.WriteByte(o)
						if b.Len() >= 2 {
							h = false
							ht = false
							d = rand.Float32()
							decoder := mahonia.NewDecoder(msp.dataRecvEncoding)
							_, cdata, _ := decoder.Translate(b.Next(b.Len()/2*2), true)
							str := string(cdata)
							//str = strings.Replace(str, "\x00", "", -1)
							//mw.txtSerialRecv.SetSuspended(true)
							mw.txtSerialRecv.AppendText(str)
							//mw.txtSerialRecv.SetSuspended(false)
						} else {
							if ht {
								ht = false
								str := string(b.Next(1))
								str = strings.Replace(str, "\x00", "", -1)
								//mw.txtSerialRecv.SetSuspended(true)
								mw.txtSerialRecv.AppendText(str)
							} else if !h {
								h = true
								dn := rand.Float32()
								d = dn
								go func(ds float32) {
									time.Sleep(time.Millisecond * 50)
									if ds == d {
										ht = true
									}
								}(dn)
							}
						}
					}
					//else if o > 127 {
					//	cb = true
					//	c.WriteByte(o)
					//}
				} else if msp.dataRecvEncoding == "utf8" {
					o := b.Bytes()[0]
					hl := 0
					if o < 0x80 { //单字节
						hl = 0
						d = rand.Float32()
						str := string(b.Next(1))
						str = strings.Replace(str, "\x00", "", -1)
						//mw.txtSerialRecv.SetSuspended(true)
						mw.txtSerialRecv.AppendText(str)
						//mw.txtSerialRecv.SetSuspended(false)
					} else if o < 0xC0 { //多字节补充
						hl = 0
						d = rand.Float32()
						str := string(b.Next(1))
						str = strings.Replace(str, "\x00", "", -1)
						//mw.txtSerialRecv.SetSuspended(true)
						mw.txtSerialRecv.AppendText(str)
						//mw.txtSerialRecv.SetSuspended(false)
					} else if o < 0xE0 { //双字节头
						hl = 2
					} else if o < 0xF0 { //三字节头
						hl = 3
					} else if o < 0xF8 { //四字节头
						hl = 4
					} else if o < 0xFC { //五字节头
						hl = 5
					} else if o < 0xFE { //六字节头
						hl = 6
					}
					if hl > 0 {
						if b.Len() >= hl {
							h = false
							ht = false
							d = rand.Float32()
							decoder := mahonia.NewDecoder(msp.dataRecvEncoding)
							_, cdata, _ := decoder.Translate(b.Next(hl), true)
							str := string(cdata)
							//str = strings.Replace(str, "\x00", "", -1)
							d = rand.Float32()
							//mw.txtSerialRecv.SetSuspended(true)
							mw.txtSerialRecv.AppendText(str)
							//mw.txtSerialRecv.SetSuspended(false)
						} else {
							if ht {
								ht = false
								d = rand.Float32()
								str := string(b.Next(1))
								str = strings.Replace(str, "\x00", "", -1)
								//mw.txtSerialRecv.SetSuspended(true)
								mw.txtSerialRecv.AppendText(str)
							} else if !h {
								h = true
								dn := rand.Float32()
								d = dn
								go func(ds float32) {
									time.Sleep(time.Millisecond * 50)
									if ds == d {
										ht = true
									}
								}(dn)
							}
						}
					}
				}
			} else {
				time.Sleep(time.Microsecond)
			}
			//if n == 0 && !isHalf {
			//	continue
			//}
			//if isHalf {
			//	buf = append([]byte{half}, buf[:n]...)
			//	n += 1
			//	isHalf = false
			//} else if n == 1 {
			//	//} else if n%2 == 1 {
			//	half = buf[:n][0]
			//	isHalf = true
			//	continue
			//}
			//// noinspection ALL0
			//log.Print(buf[:n])
			////str := string(buf[:n])
			////log.Print(str)
			//decoder := mahonia.NewDecoder(sp.dataEncoding)
			//_, cdata, _ := decoder.Translate(buf[:n], true)
			//str := string(cdata)
			//log.Print(str)
			//str = strings.Replace(str, "\x00", "", -1)
			//mw.txtSerialRecv.SetSuspended(true)
			//mw.txtSerialRecv.AppendText(str)
			//mw.txtSerialRecv.SetSuspended(false)
		}
	}()
	return nil
}

func (mw *myWindow) closeSerial() error {
	if err := msp.Close(); err != nil {
		return err
	}
	return nil
}

type myWindow struct {
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
	cbRecvEncoding *walk.ComboBox
	cbSendEncoding *walk.ComboBox
}

type mySerialPort struct {
	*serial.Port
	dataRecvEncoding string
	dataSendEncoding string
}

var msp = mySerialPort{
	dataRecvEncoding: "gbk",
	dataSendEncoding: "gbk",
}

func main() {
	mw := myWindow{}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Icon:     "favicon.ico",
		Title:    "SerialTool",
		MinSize:  Size{Width: 800, Height: 600},
		Size:     Size{Width: 800, Height: 600},
		Visible:  false,
		Layout:   HBox{},
		Children: []Widget{
			Composite{
				Layout:  VBox{MarginsZero: true},
				MaxSize: Size{Width: 160, Height: 160},
				MinSize: Size{Width: 160, Height: 160},
				Children: []Widget{
					GroupBox{
						Title:  "通讯设置",
						Layout: VBox{},
						Children: []Widget{
							Composite{
								Layout: VBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text:          "端口号",
										TextAlignment: AlignNear,
										Visible:       false,
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
										Model:         getPortNameList(),
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
										Model:         getBaudRateList(),
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
										Model:         getDataBitsList(),
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
										Model:         getParityList(),
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
										Model:         getStopBitsList(),
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
					GroupBox{
						Title:  "接收设置",
						Layout: VBox{},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "编码方式",
									},
									ComboBox{
										AssignTo:     &mw.cbRecvEncoding,
										CurrentIndex: 0,
										Model:        []string{"gbk", "utf8"},
										OnCurrentIndexChanged: func() {
											msp.dataRecvEncoding = mw.cbRecvEncoding.Value().(string)
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									CheckBox{
										Text: "十六进制",
									},
									HSpacer{},
								},
							},
						},
					},
					GroupBox{
						Title:  "发送设置",
						Layout: VBox{},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "编码方式",
									},
									ComboBox{
										AssignTo:     &mw.cbSendEncoding,
										CurrentIndex: 0,
										Model:        []string{"gbk", "utf8"},
										OnCurrentIndexChanged: func() {
											msp.dataSendEncoding = mw.cbSendEncoding.Value().(string)
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									CheckBox{
										Text: "十六进制",
									},
									HSpacer{},
								},
							},
						},
					},
					VSpacer{},
				},
			},
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
						//Font:     Font{Family: "Consolas", PointSize: 10},
						ReadOnly: true,
						VScroll:  true,
					},
					Composite{
						Layout:  HBox{MarginsZero: true},
						MaxSize: Size{Width: 130, Height: 130},
						MinSize: Size{Width: 130, Height: 130},
						Children: []Widget{
							TextEdit{
								Name:     "txtSerialSend",
								AssignTo: &mw.txtSerialSend,
								Font:     Font{Family: "Consolas", PointSize: 10},
								Text:     "我有一只小毛驴",
							},
							Composite{
								Layout:  VBox{MarginsZero: true},
								MaxSize: Size{Width: 100, Height: 100},
								MinSize: Size{Width: 100, Height: 100},
								Children: []Widget{
									Composite{
										Layout: HBox{MarginsZero: true},
										Children: []Widget{
											CheckBox{
												Text: "定时发送",
											},
											HSpacer{},
										},
									},
									Composite{
										Layout: HBox{MarginsZero: true},
										Children: []Widget{
											Label{
												Text: "周期",
											},
											NumberEdit{},
											Label{
												Text: "ms",
											},
										},
									},
									PushButton{
										Text: "发送",
										OnClicked: func() {
											encoder := mahonia.NewEncoder(msp.dataSendEncoding)
											result := encoder.ConvertString(mw.txtSerialSend.Text())
											_, err := msp.Write([]byte(result))
											if err != nil {
												log.Print(err)
											}
										},
									},
									VSpacer{},
								},
							},
							//HSpacer{},
						},
					},
					//VSpacer{},
				},
			},
			//HSpacer{},
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

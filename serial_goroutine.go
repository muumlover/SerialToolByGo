package main

import (
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {
	loop := true

	c := &serial.Config{Name: "COM6", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}

	var recv []byte
	go func(s *serial.Port) {
		for loop {
			buf := make([]byte, 1024)
			n, err = s.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			recv = append(recv, buf[:n]...)
		}
	}(s)
	go func() {
		for loop {
			if len(recv) > 0 {
				time.Sleep(10)
				print(string(recv))
				recv = *new([]byte)
			}
		}
	}()

	for loop {
		a := 1
		a++
	}
	print("exit")
}

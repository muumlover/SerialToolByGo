package main

import (
	"github.com/tarm/serial"
	"log"
)

func main() {
	c := &serial.Config{Name: "COM1", Baud: 115200}
	s, err := serial.OpenPort(c)
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
	print("exit")
}

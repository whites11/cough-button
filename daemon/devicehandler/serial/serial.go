package serial

import (
	"fmt"
	"io"

	rfcomm "github.com/jacobsa/go-serial/serial"
)

type serial struct {
	port io.ReadWriteCloser
}

func New() (*serial, error) {
	options := rfcomm.OpenOptions{
		PortName:        "/dev/ttyACM0",
		BaudRate:        38400,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := rfcomm.Open(options)
	if err != nil {
		return nil, err
	}

	return &serial{
		port: port,
	}, nil
}

func (s *serial) Cleanup() error {
	return s.port.Close()
}

func (s *serial) NotifyMicrophoneMutedState(muted bool) error {
	var data byte
	{
		if muted {
			data = 0x00
		} else {
			data = 0x01
		}
	}

	_, err := s.port.Write([]byte{data})

	return err
}

func (s *serial) ListenForMuteToggleCommand(c chan<- bool) error {
	go func() {
		buffer := make([]byte, 20)
		for true {
			n, err := s.port.Read(buffer)
			if err != nil {
				fmt.Println(err)
				return
			}

			s := string(buffer[0:n])

			switch s {
			case "toggle":
				c <- true
			default:
				fmt.Printf("Unknown message coming from serial device: %q\n", s)
			}
		}
	}()

	return nil
}

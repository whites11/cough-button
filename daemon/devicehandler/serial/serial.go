package serial

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	rfcomm "github.com/3nueves/serial"
)

type serial struct {
	port io.ReadWriteCloser
}

func New() (*serial, error) {
	devices := []string{
		"/dev/ttyACM0",
		"/dev/ttyACM1",
	}

	bauds := []int{
		//57600,
		38400,
		//19200,
		//9600,
		//115200,
	}

	for _, d := range devices {
		for _, b := range bauds {
			c := &rfcomm.Config{Name: d, Baud: b}
			port, err := rfcomm.OpenPort(c)
			if err != nil {
				fmt.Printf("Error opening device %s\n", d)
				continue
			}

			fmt.Printf("Probing device %s (baud %d).\n", d, b)

			_ = port.Flush()

			ch := make(chan bool, 10)
			defer close(ch)

			go func() {
				buffer := make([]byte, 1024)
				for true {
					n, err := port.Read(buffer)
					if err != nil {
						panic(err)
					}

					s := string(buffer[0:n])

					if strings.Trim(s, "\r\n") == "pong" {
						ch <- true
						return
					} else {
						fmt.Println(s)
					}
				}
			}()

			time.Sleep(3 * time.Second)

			fmt.Println("Sending ping.")

			_, err = port.Write([]byte("ping\n"))
			if err != nil {
				fmt.Println("Ping failed")
				continue
			}

			fmt.Println("Ping sent, waiting for response.")
			time.Sleep(5 * time.Second)

			select {
			case <-ch:
				// Got pong response.
			case <-time.After(time.Second * 10):
				// report timeout
				fmt.Println("Didn't get a ping response in time")
				_ = port.Close()
				continue
			}

			fmt.Printf("Found suitable device %s\n", d)

			return &serial{
				port: port,
			}, nil
		}
	}

	return nil, errors.New("no suitable device found")
}

func (s *serial) Cleanup() error {
	return s.port.Close()
}

func (s *serial) NotifyMicrophoneMutedState(muted bool) error {
	var data string
	{
		if muted {
			data = "muted\n"
		} else {
			data = "unmuted\n"
		}
	}

	_, err := s.port.Write([]byte(data))

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

			switch strings.Trim(s, "\r\n") {
			case "toggle":
				c <- true
			default:
				fmt.Printf("Unknown message coming from serial device: %q\n", s)
			}
		}
	}()

	return nil
}

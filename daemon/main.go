package main

import (
	"log"
	"time"

	"github.com/whites11/coughbutton/devicehandler"
	"github.com/whites11/coughbutton/devicehandler/serial"
	"github.com/whites11/coughbutton/microphonehandler"
	"github.com/whites11/coughbutton/microphonehandler/pulsedbus"
)

func main() {
	var err error
	var micHandler microphonehandler.Interface
	{
		micHandler, err = pulsedbus.New()
		testFatal(err, "Init pulse dbus")

		defer micHandler.Cleanup()
	}

	var deviceHandler devicehandler.Interface
	{
		deviceHandler, err = serial.New()
		testFatal(err, "Init device handler")

		defer deviceHandler.Cleanup()
	}

	muted, err := micHandler.IsMuted()
	testFatal(err, "Check if mic is muted at startup")

	err = deviceHandler.NotifyMicrophoneMutedState(muted)
	testFatal(err, "Initial state notify")

	// Listen for external changes in the mic state.
	c := make(chan bool, 10)
	err = micHandler.ListenMuteChanges(c)
	testFatal(err, "Listen for mute changes")

	// Listen for button triggers in the serial device.
	d := make(chan bool)
	err = deviceHandler.ListenForMuteToggleCommand(d)
	testFatal(err, "Listen for external button triggers")

	for true {
		select {
		case muted := <-c:
			// Mic mute state changed externally.
			err = deviceHandler.NotifyMicrophoneMutedState(muted)
		case _ = <-d:
			// External device triggered mute toggle.
			err = micHandler.ToggleMuted()
		default:
			// Nothing happened, just wait.f
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func testFatal(e error, msg string) {
	if e != nil {
		log.Fatalln(msg+":", e)
	}
}

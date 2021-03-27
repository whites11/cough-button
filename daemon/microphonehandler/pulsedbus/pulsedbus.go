package pulsedbus

import (
	"fmt"

	"github.com/godbus/dbus"
)

type pulsedbus struct {
	conn *dbus.Conn
}

func New() (*pulsedbus, error) {
	dbusAddr, err := retrieveDBusAddress()
	if err != nil {
		return nil, err
	}

	conn, err := dbus.Dial(dbusAddr)
	if err != nil {
		return nil, err
	}

	err = conn.Auth(nil)
	if err != nil {
		return nil, err
	}

	return &pulsedbus{
		conn: conn,
	}, nil
}

func (p *pulsedbus) Cleanup() error {
	p.StopListeningMuteChanges()

	return p.conn.Close()
}

func (p *pulsedbus) IsMuted() (bool, error) {
	activeSource, err := getActiveSource(p.conn)
	if err != nil {
		return false, err
	}

	sourceObj := p.conn.Object("org.PulseAudio.Core1.Device", *activeSource)

	muted, err := sourceObj.GetProperty("org.PulseAudio.Core1.Device.Mute")
	if err != nil {
		return false, err
	}

	return muted.Value().(bool), nil
}

func (p *pulsedbus) ListenMuteChanges(notify chan<- bool) error {
	go func() {
		core := p.conn.Object("org.PulseAudio.Core1", "/org/pulseaudio/core1")

		core.Call(
			"org.PulseAudio.Core1.ListenForSignal",
			0,
			"org.PulseAudio.Core1.Device.MuteUpdated", []dbus.ObjectPath{})

		fmt.Println("Listening for signals")

		c := make(chan *dbus.Signal, 10)
		p.conn.Signal(c)
		for v := range c {
			switch v.Name {
			case "org.PulseAudio.Core1.Device.MuteUpdated":
				muted, err := p.IsMuted()
				if err != nil {
					// TODO handle error
				}
				notify <- muted
			}
		}
	}()

	return nil
}

func (p *pulsedbus) StopListeningMuteChanges() error {
	fmt.Println("Stopping signals listening")

	core := p.conn.Object("org.PulseAudio.Core1", "/org/pulseaudio/core1")

	core.Call(
		"org.PulseAudio.Core1.StopListeningForSignal",
		0,
		"org.PulseAudio.Core1.Device.MuteUpdated")

	return nil
}

func (p *pulsedbus) SetMuted(mute bool) error {
	activeSource, err := getActiveSource(p.conn)
	if err != nil {
		return err
	}

	sourceObj := p.conn.Object("org.PulseAudio.Core1.Device", *activeSource)

	return sourceObj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.PulseAudio.Core1.Device", "Mute", dbus.MakeVariant(mute)).Err
}

func (p *pulsedbus) ToggleMuted() error {
	m, err := p.IsMuted()
	if err != nil {
		return err
	}

	return p.SetMuted(!m)
}

func getActiveSource(conn *dbus.Conn) (*dbus.ObjectPath, error) {
	obj := conn.Object("org.PulseAudio.Core1", "/org/pulseaudio/core1")

	fallbackSourcePath, err := obj.GetProperty("org.PulseAudio.Core1.FallbackSource")
	if err != nil {
		return nil, err
	}

	fallbackSource := fallbackSourcePath.Value().(dbus.ObjectPath)
	if fallbackSource.IsValid() {
		return &fallbackSource, nil
	}

	return nil, nil
}

func retrieveDBusAddress() (string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return "", err
	}

	addr, err := conn.Object("org.PulseAudio1", "/org/pulseaudio/server_lookup1").GetProperty("org.PulseAudio.ServerLookup1.Address")
	if err != nil {
		return "", err
	}

	return addr.Value().(string), nil
}

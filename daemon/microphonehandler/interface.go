package microphonehandler

type Interface interface {
	Cleanup() error
	IsMuted() (bool, error)
	SetMuted(bool) error
	ToggleMuted() error
	ListenMuteChanges(chan<- bool) error
	StopListeningMuteChanges() error
}

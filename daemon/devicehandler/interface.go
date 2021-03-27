package devicehandler

type Interface interface {
	Cleanup() error
	NotifyMicrophoneMutedState(bool) error
	ListenForMuteToggleCommand(chan<- bool) error
}

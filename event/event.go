package event

type Event int

const (
	EVENT_KEYDOWN_LEFT Event = iota
	EVENT_KEYDOWN_RIGHT
	EVENT_KEYDOWN_SPACE
)

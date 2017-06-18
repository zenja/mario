package event

type Event int

const (
	EVENT_KEYDOWN_LEFT Event = iota
	EVENT_KEYDOWN_RIGHT
	EVENT_KEYDOWN_UP
	EVENT_KEYDOWN_DOWN
	EVENT_KEYDOWN_SPACE
	EVENT_KEYDOWN_F

	// for debug use
	EVENT_KEYDOWN_F1
	EVENT_KEYDOWN_F2
	EVENT_KEYDOWN_F3
	EVENT_KEYDOWN_F4
	EVENT_KEYDOWN_F5
)

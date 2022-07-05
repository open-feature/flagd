package sync

type INotify interface {
	GetEvent() Event
}

type EEventType int32

const (
	EEventTypeCreate = iota
	EEventTypeModify = 1
	EEventTypeDelete = 2
)

type Event struct {
	EventType EEventType
}

type Notifier struct {
	Event Event
}

func (w *Notifier) GetEvent() Event {
	return w.Event
}

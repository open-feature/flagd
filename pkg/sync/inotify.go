package sync

type INotify interface {
	GetEvent() Event
}

type E_EVENT_TYPE int32

const (
	E_EVENT_TYPE_CREATE = iota
	E_EVENT_TYPE_MODIFY = 1
	E_EVENT_TYPE_DELETE = 2
)

type Event struct {
	EventType E_EVENT_TYPE
}

type Notifier struct {
	Event Event
}

func (w *Notifier) GetEvent() Event {
	return w.Event
}

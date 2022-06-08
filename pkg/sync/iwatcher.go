package sync

type IWatcher interface {
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

type Watcher struct {
	Event Event
}

func (w *Watcher) GetEvent() Event {
	return w.Event
}

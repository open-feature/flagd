package sync

type INotify interface {
	GetEvent() Event[DefaultEventType]
}

type Notifier struct {
	Event Event[DefaultEventType]
}

func (w *Notifier) GetEvent() Event[DefaultEventType] {
	return w.Event
}

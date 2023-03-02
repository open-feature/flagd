package kubernetes

type DefaultEventType int32

const (
	DefaultEventTypeCreate = iota
	DefaultEventTypeModify = 1
	DefaultEventTypeDelete = 2
	DefaultEventTypeReady  = 3
)

// Event is a struct that represents a single event.
// It is a generic struct that can be cast to a more specific struct.
type Event[T any] struct {
	EventType T
}

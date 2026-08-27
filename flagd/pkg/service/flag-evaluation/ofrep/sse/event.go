package sse

import "encoding/json"

// refetchEventType is the only event data type defined by OpenFeature protocol ADR-0008.
const refetchEventType = "refetchEvaluation"

// eventName is the SSE `event:` field, using the ADR-0008 "message" envelope.
const eventName = "message"

type refetchPayload struct {
	Type         string `json:"type"`
	Etag         string `json:"etag,omitempty"` // config version token, not an HTTP entity-tag
	LastModified int64  `json:"lastModified,omitempty"`
}

// refetchEvent implements eventsource.Event, telling subscribed OFREP clients that the flag
// configuration for their channel changed and they should re-run bulk evaluation.
type refetchEvent struct {
	id   string
	data string
}

func newRefetchEvent(id, etag string, lastModified int64) refetchEvent {
	payload := refetchPayload{
		Type:         refetchEventType,
		Etag:         etag,
		LastModified: lastModified,
	}
	data, _ := json.Marshal(payload)
	return refetchEvent{id: id, data: string(data)}
}

func (e refetchEvent) Id() string    { return e.id }
func (e refetchEvent) Event() string { return eventName }
func (e refetchEvent) Data() string  { return e.data }

package application

import "sync"

// RequestEvent is broadcast to connected clients when a new request is tracked.
type RequestEvent struct {
	Path      string `json:"path"`
	ID        string `json:"id"`
	Method    string `json:"method"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`
}

// EventHub is a minimal in-memory pub/sub used to push live updates to SSE
// clients. Publish never blocks: if a subscriber's buffer is full the event is
// dropped for that subscriber rather than stalling request tracking.
type EventHub struct {
	mu   sync.Mutex
	subs map[chan RequestEvent]struct{}
}

func NewEventHub() *EventHub {
	return &EventHub{subs: make(map[chan RequestEvent]struct{})}
}

func (h *EventHub) Subscribe() chan RequestEvent {
	ch := make(chan RequestEvent, 16)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *EventHub) Unsubscribe(ch chan RequestEvent) {
	h.mu.Lock()
	delete(h.subs, ch)
	h.mu.Unlock()
}

func (h *EventHub) Publish(event RequestEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs {
		select {
		case ch <- event:
		default:
		}
	}
}

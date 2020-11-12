package coordinator

type EventRaiser interface {
	AddListener(eventName string, f func(interface{})) //interfaces allow to send multiple types of data, not just one concrete type
}

type EventAggregator struct {
	listeners map[string][]func(interface{}) //wow, a slice of functions
}

func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}
	return &ea
}

func (ea *EventAggregator) AddListener(name string, f func(interface{})) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

func (ea *EventAggregator) PublishEvent(name string, eventData interface{}) {
	if ea.listeners[name] == nil {
		for _, r := range ea.listeners[name] {
			r(eventData)
		}
	}
}

type EventData struct {
	Item  string
	Count int
}

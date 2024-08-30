package roomServer

type EventType int

const (
	Tree_DIE = iota
	MONSTER_DIE
)

type IEvent interface {
	Type() EventType
	Data() interface{}
}

type Event struct {
	typename EventType
	data     interface{}
}

func (e *Event) Type() EventType {
	return e.typename
}
func (e *Event) Data() interface{} {
	return e.data
}

type EventCallback func(IEvent)
type EventCallBackSlice []EventCallback

type EventBus struct {
	Subscribers map[EventType]EventCallBackSlice
}

func (this *EventBus) Publish(event IEvent) {
	if val, ok := this.Subscribers[event.Type()]; ok {
		for _, eventFunc := range val {
			eventFunc(event)
		}
	}
}

func (this *EventBus) Subscribe(eventType EventType, eventFunc EventCallback) {
	if this.Subscribers == nil {
		this.Subscribers = make(map[EventType]EventCallBackSlice)
	}
	this.Subscribers[eventType] = append(this.Subscribers[eventType], eventFunc)
}

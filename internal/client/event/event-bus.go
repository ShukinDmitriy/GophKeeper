package event

type EventName string

// Event Структура события
type Event struct {
	Name EventName
	Data interface{}
}

// Subscriber - подписчик на события
type Subscriber struct {
	ch         chan *Event
	observable *Observable
}

// Unsubscribe - отписаться от событий
func (s *Subscriber) Unsubscribe() {
	close(s.ch)
	s.observable.Unsubscribe(s.ch)
}

// Observable - наблюдаемый
type Observable struct {
	subscribers []chan *Event
	ch          chan *Event
}

// NewObservable - конструктор для наблюдаемого
func NewObservable() *Observable {
	return &Observable{
		ch:          make(chan *Event, 10),
		subscribers: make([]chan *Event, 0),
	}
}

// Subscribe - подписаться на получение событий
func (o *Observable) Subscribe(next func(event *Event)) *Subscriber {
	subChan := make(chan *Event, 10)

	subscriber := &Subscriber{
		ch:         subChan,
		observable: o,
	}

	o.subscribers = append(o.subscribers, subChan)

	go func() {
		for event := range subChan {
			next(event)
		}
	}()

	return subscriber
}

// Unsubscribe - Удалить слушателя
func (o *Observable) Unsubscribe(ch chan *Event) {
	for i, subCh := range o.subscribers {
		if subCh == ch {
			o.subscribers = append(o.subscribers[:i], o.subscribers[i+1:]...)
			return
		}
	}
}

// Next - Распространить событие
func (o *Observable) Next(event *Event) {
	for _, subCh := range o.subscribers {
		if subCh != nil {
			subCh <- event
		}
	}
}

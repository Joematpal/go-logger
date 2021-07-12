package events

import (
	"sync"
	"unsafe"
)

type FactoryBus interface {
	Subscriber
	Unsubscriber
	Publisher
}

type FactoryEvent interface {
	Close() error
	Data() chan []byte
	Done() <-chan struct{}
	Error() error
}

type Unsubscriber interface {
	Unsubscribe(topic string, ptr unsafe.Pointer) error
}

type Subscriber interface {
	Subscribe(topic string) FactoryEvent
}

type Publisher interface {
	Publish(topic string, Data []byte)
}

type DataChannelSlice []FactoryEvent

type Events struct {
	sync.RWMutex
	subscribers map[string]DataChannelSlice
}

func NewEvents() *Events {
	return &Events{
		subscribers: map[string]DataChannelSlice{},
	}
}

type iterator struct {
	topic string
	done  chan struct{}
	ch    chan []byte
	err   chan error
	unsub Unsubscriber
}

func NewIterator(topic string, unsub Unsubscriber) *iterator {
	return &iterator{
		topic: topic,
		done:  make(chan struct{}),
		ch:    make(chan []byte),
		err:   make(chan error),
		unsub: unsub,
	}
}

func (i *iterator) Done() <-chan struct{} {
	return i.done
}

func (i *iterator) Data() chan []byte {
	return i.ch
}

func (i *iterator) Error() error {
	select {
	case err := <-i.err:
		return err
	default:
		return nil
	}
}

func (i *iterator) Close() error {
	select {
	case i.done <- struct{}{}:
	default:
	}
	close(i.ch)
	err := i.Error()
	close(i.err)
	close(i.done)

	uerr := i.unsub.Unsubscribe(i.topic, unsafe.Pointer(i))
	if err != nil {
		return err
	}
	return uerr
}

func (e *Events) Subscribe(topic string) FactoryEvent {
	e.Lock()

	iter := NewIterator(topic, e)
	if prev, found := e.subscribers[topic]; found {
		e.subscribers[topic] = append(prev, iter)
	} else {
		e.subscribers[topic] = []FactoryEvent{iter}
	}
	e.Unlock()

	return iter
}

func (e *Events) Unsubscribe(topic string, iter unsafe.Pointer) error {
	e.Lock()
	defer e.Unlock()

	var loc int

	val, ok := e.subscribers[topic]
	if !ok || len(val) == 0 {
		return nil
	}
	for i, v := range val {
		val, ok := v.(*iterator)
		if ok && unsafe.Pointer(val) == iter {
			loc = i

		}
	}

	e.subscribers[topic] = remove(val, loc)
	return nil
}

func remove(slice []FactoryEvent, s int) []FactoryEvent {
	return append(slice[:s], slice[s+1:]...)
}

func (e *Events) Publish(topic string, data []byte) {
	e.RLock()
	if chans, found := e.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		channels := append(DataChannelSlice{}, chans...)
		go func(data []byte, dataChannelSlices DataChannelSlice) {
			for _, itr := range dataChannelSlices {
				itr.Data() <- data
			}
		}(data, channels)
	}
	e.RUnlock()
}

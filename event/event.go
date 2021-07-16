package event

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
	Unsubscribe(ptr unsafe.Pointer) error
}

type Subscriber interface {
	Subscribe() FactoryEvent
}

type Publisher interface {
	Publish(Data []byte)
}

type DataChannelSlice []FactoryEvent

type Event struct {
	sync.RWMutex
	subscribers DataChannelSlice
}

func New() *Event {
	return &Event{
		subscribers: DataChannelSlice{},
	}
}

type iterator struct {
	topic string
	done  chan struct{}
	ch    chan []byte
	err   chan error
	unsub Unsubscriber
}

func NewIterator(unsub Unsubscriber) *iterator {
	return &iterator{
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

	uerr := i.unsub.Unsubscribe(unsafe.Pointer(i))
	if err != nil {
		return err
	}
	return uerr
}

func (e *Event) Subscribe() FactoryEvent {
	e.Lock()
	defer e.Unlock()

	iter := NewIterator(e)
	e.subscribers = append(e.subscribers, iter)

	return iter
}

func (e *Event) Unsubscribe(iter unsafe.Pointer) error {
	e.Lock()
	defer e.Unlock()

	var loc int

	val := e.subscribers

	for i, v := range val {
		val, ok := v.(*iterator)
		if ok && unsafe.Pointer(val) == iter {
			loc = i
		}
	}

	e.subscribers = remove(val, loc)
	return nil
}

func remove(slice []FactoryEvent, s int) []FactoryEvent {
	return append(slice[:s], slice[s+1:]...)
}

func (e *Event) Publish(data []byte) {
	e.RLock()
	defer e.RUnlock()
	// this is done because the slices refer to same array even though they are passed by value
	// thus we are creating a new slice with our elements thus preserve locking correctly.
	channels := append(DataChannelSlice{}, e.subscribers...)
	go func(data []byte, dataChannelSlices DataChannelSlice) {
		for _, itr := range dataChannelSlices {
			itr.Data() <- data
		}
	}(data, channels)

}

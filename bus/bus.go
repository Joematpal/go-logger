package bus

// import "unsafe"

// type FactoryBus interface {
// 	Subscriber
// 	Unsubscriber
// 	Publisher
// }

// type FactoryEvent interface {
// 	Close() error
// 	Data() chan interface{}
// 	Done() <-chan struct{}
// 	Error() error
// }

// type Unsubscriber interface {
// 	Unsubscribe(topic string, ptr unsafe.Pointer) error
// }

// type Subscriber interface {
// 	Subscribe(topic string) FactoryEvent
// }

// type Publisher interface {
// 	Publish(topic string, Data interface{})
// }

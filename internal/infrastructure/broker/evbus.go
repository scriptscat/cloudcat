package event

import (
	"sync"

	evbus "github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
)

var bus evbus.Bus
var once sync.Once

type EvBusBroker struct {
}

func NewEvBusBroker() Broker {
	ret := &EvBusBroker{}
	ret.Init()
	return ret
}

func (e *EvBusBroker) Init(option ...Option) error {
	once.Do(func() {
		bus = evbus.New()
	})
	return nil
}

func (e *EvBusBroker) Options() Options {
	return Options{}
}

func (e *EvBusBroker) Publish(topic string, data *Message, opt ...PublishOption) error {
	options := NewPublishOptions(opt...)
	bus.Publish(topic, data, options)
	return nil
}

func (e *EvBusBroker) Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := NewSubscribeOptions(opts...)
	ret := NewEvBusSubscriber(topic, options)
	if err := bus.Subscribe(topic, ret.handle(h)); err != nil {
		logrus.Errorf("event bus broker subscribe %v err: %v", topic, err)
		return nil, err
	}
	return ret, nil
}

type EvBusSubscriber struct {
	opts    SubscribeOptions
	topic   string
	_handle func(data *Message, opts PublishOptions)
}

func NewEvBusSubscriber(topic string, options SubscribeOptions) *EvBusSubscriber {
	return &EvBusSubscriber{
		opts:  options,
		topic: topic,
	}
}

func (n *EvBusSubscriber) Options() SubscribeOptions {
	return n.opts
}

func (n *EvBusSubscriber) Topic() string {
	return n.topic
}

func (n *EvBusSubscriber) Unsubscribe() error {
	return bus.Unsubscribe(n.topic, n._handle)
}

func (n *EvBusSubscriber) handle(h Handler) func(data *Message, opts PublishOptions) {
	n._handle = func(data *Message, opts PublishOptions) {
		go func() {
			if err := h(data); err != nil {
				logrus.Errorf("event bus subscriber handler %v err: %v", n.topic, err)
			}
		}()
	}
	return n._handle
}

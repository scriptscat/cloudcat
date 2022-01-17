package event

import "context"

var DefaultBroker = NewEvBusBroker()

type Ids map[string]int64

type Options struct {
	Context context.Context
}

type PublishOptions struct {
	Context context.Context
}

type SubscribeOptions struct {
	Queue string

	Context context.Context
}

type Message struct {
	Header map[string]string
	Body   []byte
}

type PublishOption func(options *PublishOptions)
type SubscribeOption func(options *SubscribeOptions)
type Handler func(*Message) error

type Option func(*Options)
type Broker interface {
	Init(...Option) error
	Options() Options
	Publish(topic string, data *Message, opt ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
}

type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

func NewOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func NewPublishOptions(opts ...PublishOption) PublishOptions {
	opt := PublishOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

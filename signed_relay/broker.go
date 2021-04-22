package signedrelay

import (
	"fmt"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	pb "github.com/psychobummer/pbrelay/rpc/relay"
	"github.com/rs/zerolog/log"
)

// The Broker exposes a means of adding/removing producers, managing subscriptions to those producers, and
// delivering messages from a producer to its subscribed clients. It effectively provides a higher-level
// aggregation around the functionality exposed by the Producer and Subscriber types.
// TODO: this locking feels overly aggressive, but some of these operations are kind of two-phase-commit-like
// and i'm not yet sure if it's sufficient to rely on the locking within cmap
type Broker struct {
	mutex sync.Mutex
	cache cmap.ConcurrentMap
}

// NewBroker returns an initailized Broker.
func NewBroker() *Broker {
	broker := &Broker{
		mutex: sync.Mutex{},
		cache: cmap.New(),
	}
	return broker
}

// Broadcast will broadcast a given message from a producer to its subscribed clients
func (b *Broker) Broadcast(message *pb.StreamMessage) {
	producer, err := b.Producer(message.GetId())
	if err != nil {
		log.Error().Msg(err.Error())
	}

	for _, sub := range producer.Subscribers() {
		select {
		case sub.queue <- message:
		default:
			log.Trace().Msgf("dropping message for %s; it can't keep up", sub.id)
		}
	}
}

// RegisterProducer will create a new Producer identified by ID and register it with in the broker.
// Any messages from this ID will be sent to its subscribed clients.
func (b *Broker) RegisterProducer(producerID string) error {
	if b.cache.Has(producerID) {
		return fmt.Errorf("producer %s already exists", producerID)
	}
	b.cache.Set(producerID, NewProducer(producerID))
	return nil

}

// DeregisterProducer will shutdown all subscribers to a given producerID and purge the producer
// from the broker.
// OPTIMIZE per-producer lock?
func (b *Broker) DeregisterProducer(producerID string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	producer, err := b.Producer(producerID)
	if err != nil {
		return err
	}

	// we don't call DeregisterSubscriber() here as it's faster to just do this all in bulk
	for _, sub := range producer.Subscribers() {
		sub.Shutdown()
	}
	b.cache.Remove(producerID)

	return nil
}

// RegisterSubscriber will register a subscriber to a producer, and return a channel
// of messages representing that subscription.
// OPTIMIZE per-producer lock?
func (b *Broker) RegisterSubscriber(producerID string, subscriberID string) (<-chan *pb.StreamMessage, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	producer, err := b.Producer(producerID)
	if err != nil {
		return nil, err
	}

	sub := NewSubscriber(subscriberID)
	if err := producer.RegisterSubscriber(sub); err != nil {
		sub.Shutdown()
		return nil, err
	}
	b.cache.Set(producerID, producer)
	return sub.queue, nil
}

// DeregisterSubscriber will deregister a subscriber from a producer, and free any resources
// bound to the subscriber.
func (b *Broker) DeregisterSubscriber(producerID string, subscriberID string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	producer, err := b.Producer(producerID)
	if err != nil {
		return err
	}

	if err := producer.DeregisterSubscriber(subscriberID); err != nil {
		return err
	}

	b.cache.Set(producerID, producer)
	return nil
}

// Producer will return the Producer identified by the provided ID
func (b *Broker) Producer(producerID string) (Producer, error) {
	producer, found := b.cache.Get(producerID)
	if !found {
		return Producer{}, fmt.Errorf("producer %s doesn't exist", producerID)
	}
	return producer.(Producer), nil
}

// A Producer exposes means of lifecycling its subscribers. That is, adding them, removing them, etc
type Producer struct {
	id          string
	subscribers cmap.ConcurrentMap
}

// NewProducer returns a Producer identified by the specificied ID
func NewProducer(producerID string) Producer {
	producer := Producer{
		id:          producerID,
		subscribers: cmap.New(),
	}
	return producer
}

// DeregisterSubscriber deregisters the specified subscriber from this producer
func (p *Producer) DeregisterSubscriber(id string) error {
	sub, err := p.Subscriber(id)
	if err != nil {
		return fmt.Errorf("producer %s has no subscriber %s", p.id, id)
	}
	sub.Shutdown()
	p.subscribers.Remove(id)
	return nil
}

// RegisterSubscriber registers the provided subscriber with this produer
func (p *Producer) RegisterSubscriber(sub Subscriber) error {
	if p.subscribers.Has(sub.id) {
		return fmt.Errorf("pubscriber %s already registered to producer %s", sub.id, p.id)
	}
	p.subscribers.Set(sub.id, sub)
	return nil
}

// Subscriber returns the Subscriber for this producer identified by provided ID
func (p *Producer) Subscriber(subscriberID string) (Subscriber, error) {
	sub, found := p.subscribers.Get(subscriberID)
	if !found {
		return Subscriber{}, fmt.Errorf("producer %s does not have a subscriber %s", p.id, subscriberID)
	}
	return sub.(Subscriber), nil
}

// Subscribers returns a slice containing all of this producers subscribers
func (p *Producer) Subscribers() []Subscriber {
	subs := []Subscriber{}
	for sub := range p.subscribers.IterBuffered() {
		subs = append(subs, sub.Val.(Subscriber))
	}
	return subs
}

// A Subscriber represents, ultimately, a connected client, with the queue containing messages
// destinated for it.
type Subscriber struct {
	id    string
	queue chan *pb.StreamMessage
}

// NewSubscriber returns an initialized Subscriber identified by the provided ID
func NewSubscriber(subscriberID string) Subscriber {
	sub := Subscriber{
		id:    subscriberID,
		queue: make(chan *pb.StreamMessage, 100),
	}
	return sub
}

// Shutdown will close the subscribers read queue, causing the handler to return.
// I'm exposing this through a specific method in case we want to do anything else.
func (s Subscriber) Shutdown() {
	close(s.queue)
}

package signedrelay

import (
	"testing"
	"time"

	pb "github.com/psychobummer/pbrelay/rpc/relay"
	"github.com/stretchr/testify/assert"
)

func TestRegisterProducer(t *testing.T) {
	broker := NewBroker()
	assert.NoError(t, broker.RegisterProducer("foo"))
	assert.Error(t, broker.RegisterProducer("foo")) // dupe; should error
}

func TestDeregisterProducer(t *testing.T) {
	broker := NewBroker()
	assert.NoError(t, broker.RegisterProducer("foo"))
	assert.NoError(t, broker.DeregisterProducer("foo"))
	assert.Error(t, broker.DeregisterProducer("foo")) // should be gone
}

func TestRegisterSubscriber(t *testing.T) {
	pid := "foo"
	broker := NewBroker()
	assert.NoError(t, broker.RegisterProducer(pid))

	_, err := broker.RegisterSubscriber(pid, "sub")
	assert.NoError(t, err)

	_, err = broker.RegisterSubscriber("thisbrokerdoesntexist", "sub")
	assert.Error(t, err)

	producer, err := broker.Producer(pid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(producer.Subscribers()))
}

func TestDeregisterSubscriber(t *testing.T) {
	pid := "foo"
	broker := NewBroker()
	assert.NoError(t, broker.RegisterProducer(pid))
	queue, err := broker.RegisterSubscriber(pid, "somesub")
	assert.NoError(t, err)

	done := make(chan bool)
	go func() {
		select {
		case <-queue: // the queue chan is closed as soon as the sub is deregistered
		case <-time.After(100 * time.Millisecond):
			t.Errorf("TestDeregisterSubscriber timed out waiting for subscription deregistration")
		}
		done <- true
	}()
	assert.NoError(t, broker.DeregisterSubscriber(pid, "somesub"))
	<-done
	producer, err := broker.Producer(pid)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(producer.Subscribers()))
}

func TestBroadcast(t *testing.T) {
	pid := "foo"
	broker := NewBroker()
	assert.NoError(t, broker.RegisterProducer(pid))
	queue, err := broker.RegisterSubscriber(pid, "somesub")
	assert.NoError(t, err)

	done := make(chan bool)
	go func() {
		select {
		case m := <-queue:
			assert.Equal(t, pid, string(m.GetData()))
		case <-time.After(100 * time.Millisecond):
			t.Errorf("TestBroadcast timed out waiting for messag to arrive")
		}
		done <- true
	}()

	message := &pb.StreamMessage{
		Id:   pid,
		Data: []byte(pid),
	}
	broker.Broadcast(message)
	<-done
}

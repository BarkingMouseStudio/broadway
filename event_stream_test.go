package broadway

import (
	"testing"
)

func TestEventStream(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	r := make(chan interface{})
	a := s.ActorOf(&echoBehaviour{r}, "Test")
	s.Events.Subscribe(a, "TEST")
	s.Events.Publish("TEST", echoMessage{})
	<-r
	s.Events.Unsubscribe(a, "TEST")
	s.Events.Publish("TEST", echoMessage{})

	select {
	case <-r:
		t.Error("Received after unsubscribed")
	default:
	}
}

func TestEventStream_UnsubscribeAll(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	r := make(chan interface{})
	a := s.ActorOf(&echoBehaviour{r}, "A")
	b := s.ActorOf(&echoBehaviour{r}, "B")

	s.Events.Subscribe(a, "A")
	s.Events.Subscribe(a, "B")
	s.Events.Subscribe(a, "B") // 2x subscribe

	// Multi-subscribers
	s.Events.Subscribe(b, "A")
	s.Events.Subscribe(b, "B")
	s.Events.Unsubscribe(b, "B")
	s.Events.Unsubscribe(b, "B")

	s.Events.Publish("A", echoMessage{})
	<-r
	<-r

	s.Events.Publish("B", echoMessage{})
	<-r

	s.Events.UnsubscribeAll(a)

	s.Events.Publish("A", echoMessage{})
	s.Events.Publish("B", echoMessage{})

	select {
	case <-r:
		t.Error("Received after unsubscribed")
	default:
	}
}

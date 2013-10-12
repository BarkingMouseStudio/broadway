package broadway

import (
	"testing"
)

type deadLetterBehaviour struct {
	results chan<- interface{}
}

func (e *deadLetterBehaviour) Receive(message interface{}, sender ActorRef, context *Actor) {
	switch message := message.(type) {
	case DeadLetter:
		e.results <- message
	}
}

func TestDeadLetters(t *testing.T) {
	s := NewActorSystem("Test", Config{})

	r := make(chan interface{})
	a := s.ActorOf(&deadLetterBehaviour{r}, "Test")
	s.Events.Subscribe(a, "DeadLetter")
	s.DeadLetters.Tell(DeadLetter{}, nil)

	<-r
}

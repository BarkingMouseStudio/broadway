package broadway

import (
	"testing"
)

type loggerBehaviour struct {
	results chan<- interface{}
}

func (e *loggerBehaviour) Receive(message interface{}, sender ActorRef, context *Actor) {
	switch message := message.(type) {
	case echoMessage:
		context.System.Logger.Logf("Logger test: %v", message)
		context.System.Logger.Log("Logger test")
		e.results <- message
	}
}

func TestLogger(t *testing.T) {
	s := NewActorSystem("Test", Config{
		Logging: NewLoggingConfig(),
	})
	r := make(chan interface{})
	a := s.ActorOf(&loggerBehaviour{r}, "Test")
	a.Tell(echoMessage{}, nil)
	<-r
	s.Shutdown()
}

package broadway

import (
	"os"
	"sync"
	"testing"
)

type echoBehaviour struct {
	results chan<- interface{}
}

type echoMessage struct{}

func (e *echoBehaviour) Receive(message interface{}, sender ActorRef, context *Actor) {
	switch message := message.(type) {
	case echoMessage:
		e.results <- message
	default:
		context.System.DeadLetters.Tell(message, sender)
	}
}

func TestNewActorSystem(t *testing.T) {
	name := "Test"
	conf := Config{
		Mailbox: NewMailboxConfig(),
		Logging: LoggingConfig{
			LogLifecycle: false,
			LogReceive:   false,
			Logger:       os.Stdout,
		},
	}
	s := NewActorSystem(name, conf)
	if s.guardian == nil {
		t.Error("Guardian nil")
	}
	if s.DeadLetters == nil {
		t.Error("DeadLetters nil")
	}
	if &s.Logger == nil {
		t.Error("Logger nil")
	}
	if &s.Events == nil {
		t.Error("Events nil")
	}
	if s.config != conf {
		t.Error("Config mismatch")
	}
	if s.name != name {
		t.Errorf("Invalid name, expected %v, got %v", name, s.name)
	}
	s.Shutdown()
}

func TestActorSystemActorOf(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	r := s.ActorOf(nil, "Actor")
	expected := "$guardian/Actor"
	path := r.Path()
	if path.String() != expected {
		t.Errorf("Invalid child path, expected %v, got %v", expected, path.String())
	}
	s.Shutdown()
}

func TestActorSystemShutdown(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	s.ActorOf(nil, "Actor")
	s.Shutdown()
}

func BenchmarkActorSystem(b *testing.B) {
	b.StopTimer()

	var wg sync.WaitGroup

	s := NewActorSystem("Benchmark", Config{
		Mailbox: MailboxConfig{
			BufferSize: 10,
		},
		Logging: LoggingConfig{},
	})
	results := make(chan interface{})
	actor := s.ActorOf(&echoBehaviour{results}, "Actor") // Echoes back messages sent to it

	go func() {
		for _ = range results {
			wg.Done()
		}
	}()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		actor.Tell(echoMessage{}, nil) // Send messages to actor
	}

	wg.Wait()
	b.StopTimer()

	s.Shutdown()
}

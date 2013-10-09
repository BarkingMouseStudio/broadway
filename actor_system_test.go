package broadway

import (
	"sync"
	"testing"
)

type echo struct {
	results chan<- interface{}
}

func (e *echo) Receive(message interface{}, sender ActorRef, context *Actor) {
	e.results <- message
}

func TestNewActorSystem(t *testing.T) {
	name := "Test"
	conf := Config{}
	s := NewActorSystem(name, conf)
	if s.guardian == nil {
		t.Error("Guardian nil")
	}
	if s.DeadLetters == nil {
		t.Error("DeadLetters nil")
	}
	if s.config != conf {
		t.Error("Config mismatch")
	}
	if s.name != name {
		t.Errorf("Invalid name, expected %v, got %v", name, s.name)
	}
}

func TestActorSystemActorOf(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	r := s.ActorOf(nil, "Actor")
	expected := "$guardian/Actor"
	path := r.Path()
	if path.String() != expected {
		t.Errorf("Invalid child path, expected %v, got %v", expected, path.String())
	}
}

func TestActorSystemShutdown(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	s.ActorOf(nil, "Actor")
	s.Shutdown()
}

func BenchmarkActorSystem(b *testing.B) {
	b.StopTimer()

	var wg sync.WaitGroup

	config := NewConfig()
	config.mailbox = MailboxConfig{
		bufferSize:     100,
		overflowPolicy: BlockOnOverflow,
	}

	system := NewActorSystem("Benchmark", config)
	results := make(chan interface{})
	actor := system.ActorOf(&echo{results}, "ActorA") // Echoes back messages sent to it

	go func() {
		for _ = range results {
			wg.Done()
		}
	}()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		actor.Tell("test", nil) // Send messages to actor
	}

	wg.Wait()
}

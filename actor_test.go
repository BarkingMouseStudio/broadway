package broadway

import (
	"testing"
)

func TestActor(t *testing.T) {
	name := "Test"
	s := NewActorSystem("Test", Config{})
	a := newActor(&echoBehaviour{make(chan interface{})}, name, s, s.guardian.Self)
	if a.System == nil {
		t.Error("System nil")
	}
	if a.Parent == nil {
		t.Error("Parent nil")
	}
	if a.Self == nil {
		t.Error("Self nil")
	}
	if a.name != name {
		t.Error("Bad name", a.name)
	}
	if a.Path.String() != "$guardian/Test" {
		t.Error("Bad path", a.Path.String())
	}
	if a.children == nil {
		t.Error("Children nil")
	}
	if a.receiver == nil {
		t.Error("Invalid receiver")
	}
	if a.mailbox == nil {
		t.Error("Mailbox nil")
	}

	// ActorOf
	r1 := s.ActorOf(nil, "A")
	r2 := s.ActorOf(nil, "A")
	if r1 != r2 {
		t.Error("Did not receive the same actor back")
	}

	// Stop
	a.Self.Stop(nil)
	r1.Stop(nil)

	s.Shutdown()
}

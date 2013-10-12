package broadway

import (
	"testing"
)

func TestActorRefString(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	ref := s.ActorOf(nil, "Actor")
	path := ref.Path()
	if ref.String() != path.String() {
		t.Error("Invalid string", ref.String(), path.String())
	}
	s.Shutdown()
}

func TestActorRefEquals(t *testing.T) {
	s := NewActorSystem("Test", Config{})
	a := s.ActorOf(nil, "Actor")
	b := s.ActorOf(nil, "Actor")
	if !a.Equals(b) {
		t.Error("Did not compare paths properly", a.String(), b.String())
	}
	s.Shutdown()
}

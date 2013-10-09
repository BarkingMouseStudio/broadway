package broadway

import (
	"testing"
)

func TestActorPath(t *testing.T) {
	p := NewActorPath("B", NewActorPath("A", nil))
	if p[0] != "A" {
		t.Errorf("Expected %v to be %v", p[0], "A")
	}
	if p[1] != "B" {
		t.Errorf("Expected %v to be %v", p[1], "B")
	}
}

func TestActorPathName(t *testing.T) {
	p := NewActorPath("A", nil)
	name := p.Name()
	expected := "A"
	if name != expected {
		t.Errorf("Expected %v to be %v", name, expected)
	}
}

func TestActorPathString(t *testing.T) {
	p := NewActorPath("B", NewActorPath("A", nil))
	result := p.String()
	expected := "A/B"
	if result != expected {
		t.Errorf("Expected %v to be %v", result, expected)
	}
}

func TestActorPathEquals(t *testing.T) {
	a := NewActorPath("A", nil)
	b := NewActorPath("B", a)
	c := NewActorPath("B", a)
	if a.Equals(b) {
		t.Errorf("Expected %v not to equal %v", a, b)
	}
	if !b.Equals(c) {
		t.Errorf("Expected %v to equal %v", b, c)
	}
}

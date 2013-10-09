package broadway

import (
	"fmt"
)

// ActorRef enables location transparency, representing
// an instance of the running actor in-process or on a
// remote machine.
type ActorRef interface {
	Tell(message interface{}, sender ActorRef)
	Path() ActorPath
	Equals(ActorRef) bool

	fmt.Stringer
}

type LocalActorRef struct {
	path    ActorPath
	mailbox *Mailbox
}

// Notifies the actor of a message
func (r *LocalActorRef) Tell(message interface{}, sender ActorRef) {
	r.mailbox.Enqueue(Envelope{Message{message, sender}, r})
}

func (r *LocalActorRef) Path() ActorPath {
	return r.path
}

func (r *LocalActorRef) String() string {
	return r.path.String()
}

// Do the two ActorRefs refer to the same underlying actors (comparing ActorPaths.)
func (r *LocalActorRef) Equals(s ActorRef) bool {
	path := r.Path()
	other := s.Path()
	return path.Equals(other)
}

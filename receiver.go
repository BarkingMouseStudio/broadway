package broadway

// Defines the body of the Actor. Also known as the Actor's Behaviour
// but we're following idiomatic interface naming conventions here.
type Receiver interface {

	// The body of the Actor. Takes a message, the sender of the
	// message and a context which contains state of the actor
	// relevant to handling the message.
	Receive(message interface{}, sender ActorRef, context *Actor)
}

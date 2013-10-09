package broadway

type DeadLetter struct {
	Message  interface{}
	Sender   ActorRef
	Receiver ActorRef
}

// Receiver

// If you choose not to pass in a sender reference into the Tell method then a reference to the 'dead-letter' Actor will be used.
// The 'dead-letter' Actor is where all unhandled messages end up, and you can use Akka's Event Bus to subscribe on them.
// When a message is sent to an Actor that is terminated before receiving the message, it will be sent as a DeadLetter to the ActorSystem's EventStream
type DeadLetters struct{}

func (d *DeadLetters) Receive(message interface{}, sender ActorRef, context *Actor) {
	switch message := message.(type) {
	case DeadLetter:
		context.System.Events.Publish("DeadLetter", message)
	}
}

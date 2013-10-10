package broadway

/*
 * An actor system is a hierarchical group of actors which share common
 * configuration, e.g. dispatchers, deployments, remote capabilities and
 * addresses. It is also the entry point for creating or looking up actors.
 */
type ActorSystem struct {
	// The name of this actor system used for remote discovery.
	name string

	// The top-level supervisor of all actors created using system.ActorOf(...).
	guardian *Actor

	// Default recipient
	deadLetters ActorRef

	// Logger reference for the system
	Logger LoggingActorRef

	// Event stream reference for the system
	Events EventStreamActorRef

	// Config used when creating new actors in the system.
	config Config
}

func NewActorSystem(name string, config Config) *ActorSystem {
	s := &ActorSystem{
		name:   name,
		config: config,
	}
	s.guardian = newActor(nil, "$guardian", s, nil)
	s.deadLetters = s.guardian.ActorOf(&DeadLetters{}, "$deadLetters")
	s.Logger = LoggingActorRef{s.guardian.ActorOf(&Logger{wc: s.config.Logging.Logger}, "$logger")}
	s.Events = EventStreamActorRef{s.guardian.ActorOf(&EventStream{}, "$events")}
	return s
}

// Create new actor as child of the system with the given name.
func (s *ActorSystem) ActorOf(receiver Receiver, name string) ActorRef {
	return s.guardian.ActorOf(receiver, name) // Actor has lock on children making this safe
}

// Stop this actor system. This will stop the guardian actor,
// which in turn will recursively stop all its child actors.
func (s *ActorSystem) Shutdown() {
	s.guardian.Self.Stop(s.guardian.Self)
}

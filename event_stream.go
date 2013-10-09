package broadway

// Messages

type Event struct {
	message interface{}
	channel string
}

type SubscribeMessage struct {
	subscriber ActorRef
	channel    string
}

type UnsubscribeMessage struct {
	subscriber ActorRef
	channel    string
}

type UnsubscribeAllMessage struct {
	subscriber ActorRef
}

// Receiver

type EventStream struct {
	channels map[string][]ActorRef
}

func add(refs []ActorRef, newRef ActorRef) []ActorRef {
	for _, ref := range refs {
		if ref.Equals(newRef) {
			return refs // Already exists
		}
	}
	return append(refs, newRef)
}

func remove(refs []ActorRef, oldRef ActorRef) []ActorRef {
	for i, ref := range refs {
		if ref.Equals(oldRef) {
			copy(refs[i:], refs[i+1:]) // Shift left by 1
			return refs[:len(refs)-1]  // Shrink by 1
		}
	}
	return refs
}

func (e *EventStream) Publish(event Event, self ActorRef) {
	if refs, ok := e.channels[event.channel]; ok {
		for _, ref := range refs {
			ref.Tell(event.message, self)
		}
	}
}

func (e *EventStream) Subscribe(subscriber ActorRef, channel string) {
	if refs, ok := e.channels[channel]; ok {
		e.channels[channel] = add(refs, subscriber)
	} else {
		// First ref for channel, make the slice
		e.channels[channel] = []ActorRef{subscriber}
	}
}

func (e *EventStream) Unsubscribe(subscriber ActorRef, channel string) {
	if refs, ok := e.channels[channel]; ok {
		e.channels[channel] = remove(refs, subscriber)
	}
}

func (e *EventStream) UnsubscribeAll(subscriber ActorRef) {
	for channel, refs := range e.channels {
		e.channels[channel] = remove(refs, subscriber)
	}
}

func (e *EventStream) Receive(message interface{}, sender ActorRef, context *Actor) {
	switch message := message.(type) {
	case Event:
		e.Publish(message, context.Self)
	case SubscribeMessage:
		e.Subscribe(message.subscriber, message.channel)
	case UnsubscribeMessage:
		e.Unsubscribe(message.subscriber, message.channel)
	case UnsubscribeAllMessage:
		e.UnsubscribeAll(message.subscriber)
	}
}

// ActorRef

type EventStreamActorRef struct {
	ActorRef
}

func (e *EventStreamActorRef) Publish(channel string, message interface{}) {
	e.Tell(Event{
		message: message,
		channel: channel,
	}, nil)
}

func (e *EventStreamActorRef) Subscribe(subscriber ActorRef, channel string) {
	e.Tell(SubscribeMessage{
		subscriber: subscriber,
		channel:    channel,
	}, nil)
}

func (e *EventStreamActorRef) Unsubscribe(subscriber ActorRef, channel string) {
	e.Tell(UnsubscribeMessage{
		subscriber: subscriber,
		channel:    channel,
	}, nil)
}

func (e *EventStreamActorRef) UnsubscribeAll(subscriber ActorRef) {
	e.Tell(UnsubscribeAllMessage{
		subscriber: subscriber,
	}, nil)
}

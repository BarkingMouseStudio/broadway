package broadway

import (
	"sync"
)

// Stops the actor goroutine
type StopMessage struct{}

type Actor struct {
	System *ActorSystem // The system that the actor belongs to.
	Parent ActorRef
	Self   ActorRef
	Path   ActorPath

	name     string
	children map[string]ActorRef
	receiver Receiver
	mailbox  *Mailbox

	sync.RWMutex
}

func newActor(receiver Receiver, name string, system *ActorSystem, parent ActorRef) *Actor {
	var parentPath ActorPath
	if parent != nil {
		parentPath = parent.Path()
	}
	path := NewActorPath(name, parentPath)
	mailbox := NewMailbox(system.config.Mailbox)
	self := &LocalActorRef{
		mailbox: mailbox,
		path:    path,
	}
	a := &Actor{
		System: system,
		Parent: parent,
		Self:   self,
		Path:   path,

		name:     name,
		children: make(map[string]ActorRef),
		receiver: receiver,
		mailbox:  mailbox,
	}
	go a.run()
	return a
}

func (a *Actor) ActorOf(receiver Receiver, name string) ActorRef {
	a.Lock()
	defer a.Unlock()

	if c, ok := a.children[name]; ok {
		// Return the existing actor
		return c
	}

	c := newActor(receiver, name, a.System, a.Self)
	a.children[name] = c.Self
	return c.Self
}

func (a *Actor) stopping() {
	a.Lock()
	defer a.Unlock()

	if !a.Path.Equals(a.System.Logger.Path()) && a.System.config.Logging.LogLifecycle {
		a.System.Logger.Log(a.Path.String(), "stopping")
	}

	// Stop all children then return (who will, in turn,
	// stop their children and return)
	for _, child := range a.children {
		child.Stop(a.Self)
	}
}

// Implements a select over message types
func (a *Actor) run() {
	if !a.Path.Equals(a.System.Logger.Path()) && a.System.config.Logging.LogLifecycle {
		a.System.Logger.Log(a.Path.String(), "starting")
	}

	for {
		envelope := a.mailbox.Dequeue() // Blocks until message ready

		if !a.Path.Equals(a.System.Logger.Path()) && a.System.config.Logging.LogReceive {
			a.System.Logger.Logf("%s <- %#v", a.name, envelope)
		}

		switch envelope.message.body.(type) {
		case StopMessage:
			a.stopping()
			return
		default:
			if a.receiver != nil {
				if envelope.message.sender == nil {
					a.receiver.Receive(envelope.message.body, a.System.DeadLetters, a)
				} else {
					a.receiver.Receive(envelope.message.body, envelope.message.sender, a)
				}
			}
		}
	}
}

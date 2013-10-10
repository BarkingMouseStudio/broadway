package broadway

import (
	"errors"
)

var OverflowError = errors.New("Mailbox overflowed")

type MailboxConfig struct {
	BufferSize, OverflowPolicy uint
}

func NewMailboxConfig() MailboxConfig {
	return MailboxConfig{}
}

const (
	BlockOnOverflow = iota // Default channel behaviour
	PanicOnOverflow
	DropOnOverflow
)

// Message encapsulates a message body and its sender (used for replies).
type Message struct {
	body   interface{}
	sender ActorRef
}

// Envelope encapsulates a message and its recipient (used for dispatching).
type Envelope struct {
	message   Message
	recipient ActorRef
}

// Simple mailbox implemented as an channel with the given buffer size. Can define an overflow policy to handle full buffers.
type Mailbox struct {
	mailbox        chan Envelope
	overflowPolicy uint
}

func NewMailbox(config MailboxConfig) *Mailbox {
	return &Mailbox{
		mailbox:        make(chan Envelope, config.BufferSize),
		overflowPolicy: config.OverflowPolicy,
	}
}

func (m *Mailbox) Enqueue(envelope Envelope) {
	switch m.overflowPolicy {
	case PanicOnOverflow, DropOnOverflow:
		select {
		case m.mailbox <- envelope: // Attempt to send
		default:
			if m.overflowPolicy == PanicOnOverflow {
				panic(OverflowError) // Supervisor can handle this
			}
		}
	case BlockOnOverflow:
		m.mailbox <- envelope
	}
}

func (m *Mailbox) Dequeue() Envelope {
	return <-m.mailbox
}

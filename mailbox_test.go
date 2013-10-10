package broadway

import (
	"testing"
	"time"
)

/* case ThrowOnOverflow, DropOnOverflow:
  select {
  case m.mailbox <- envelope:
  default:
    if m.overflowPolicy == ThrowOnOverflow {
      panic(errors.New("Mailbox overflowed")) // Supervisor can handle this
    }
  }
case BlockOnOverflow:
  m.mailbox <- envelope */

func TestMailbox(t *testing.T) {
	m := NewMailbox(MailboxConfig{})
	ch := make(chan struct{})

	go func() {
		m.Dequeue()
		ch <- struct{}{}
	}()

	m.Enqueue(Envelope{})

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Error("Failed to dequeue")
	}
}

func TestMailbox_DropOnOverflow(t *testing.T) {
	m := NewMailbox(MailboxConfig{OverflowPolicy: DropOnOverflow})
	m.Enqueue(Envelope{})
	m.Enqueue(Envelope{}) // This will be dropped and should not deadlock
}

func TestMailbox_ThrowOnOverflow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != OverflowError {
			t.Error("Unexpected error", r)
		}
	}()

	m := NewMailbox(MailboxConfig{OverflowPolicy: PanicOnOverflow})
	m.Enqueue(Envelope{})
	m.Enqueue(Envelope{}) // This will error
}

func TestMailbox_BlockOnOverflow(t *testing.T) {
	m := NewMailbox(MailboxConfig{OverflowPolicy: BlockOnOverflow})

	go func() {
		m.Enqueue(Envelope{})
		m.Enqueue(Envelope{}) // This block
		t.Error("Enqueue did not block")
	}()

	<-time.After(time.Second) // Give it time to block
}

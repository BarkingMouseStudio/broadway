package broadway

import (
	"testing"
)

func TestConfig(t *testing.T) {
	c := NewConfig()

	if c.Mailbox.BufferSize != 0 || c.Mailbox.OverflowPolicy != 0 {
		t.Error("Mailbox config improperly initialized")
	}

	if c.Logging.LogLifecycle != true || c.Logging.LogReceive != true {
		t.Error("Logging config improperly intialized")
	}
}

package broadway

import (
	"fmt"
	"io"
	"os"
)

// Config

type LoggingConfig struct {
	LogLifecycle bool
	LogReceive   bool
	Logger       io.Writer
}

func NewLoggingConfig() LoggingConfig {
	return LoggingConfig{
		LogLifecycle: true,
		LogReceive:   true,
		Logger:       os.Stdout,
	}
}

// Messages

type LogMessage struct {
	v []interface{}
}

type LogfMessage struct {
	format string
	v      []interface{}
}

// Receiver

// Handles writing LogMessages from the event stream to the log destination
type Logger struct {
	wc io.Writer
}

func NewLogger(config LoggingConfig) Logger {
	return Logger{config.Logger}
}

func (l *Logger) Receive(message interface{}, sender ActorRef, context *Actor) {
	switch message := message.(type) {
	case LogMessage:
		fmt.Fprintln(l.wc, message.v...)
	case LogfMessage:
		fmt.Fprintln(l.wc, fmt.Sprintf(message.format, message.v...))
	}
}

// ActorRef

type LoggingActorRef struct {
	ActorRef
}

func (r *LoggingActorRef) Log(v ...interface{}) {
	r.Tell(LogMessage{v}, nil)
}

func (r *LoggingActorRef) Logf(format string, v ...interface{}) {
	r.Tell(LogfMessage{format, v}, nil)
}

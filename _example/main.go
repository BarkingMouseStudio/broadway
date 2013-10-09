package main

import (
	"github.com/FreeFlow/broadway"
	"time"
)

type Greeting struct {
	who string
}

type GreetingActor struct{}

func (a *GreetingActor) Receive(message interface{}, sender broadway.ActorRef, context *broadway.Actor) {
	switch message := message.(type) {
	case Greeting:
		context.System.Logger.Log("Hello", message.who)
	}
}

func main() {
	system := broadway.NewActorSystem("MySystem", broadway.NewConfig())
	greeter := system.ActorOf(&GreetingActor{}, "greeter")
	greeter.Tell(Greeting{"Charlie Parker"}, nil)

	time.Sleep(100 * time.Millisecond)
}

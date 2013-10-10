package main

import (
	"github.com/FreeFlow/broadway"
)

type Greeting struct {
	who string
}

type GreetingActor struct {
	done chan struct{}
}

func (a *GreetingActor) Receive(message interface{}, sender broadway.ActorRef, context *broadway.Actor) {
	switch message := message.(type) {
	case Greeting:
		context.System.Logger.Log("Hello", message.who)
		a.done <- struct{}{}
	}
}

func main() {
	done := make(chan struct{})
	system := broadway.NewActorSystem("MySystem", broadway.NewConfig())
	greeter := system.ActorOf(&GreetingActor{done}, "greeter")
	greeter.Tell(Greeting{"Charlie Parker"}, nil)
	<-done
}

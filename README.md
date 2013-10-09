Broadway
===

Actor library for Go. Designed to provide minimal
abstractions on top of Goroutines and channels. Based
loosely on Akka.

Refs, Selections, Paths
===

Akka uses a lot of terms to describe its system. Here are the breakdowns of some of those terms:
ActorRef - handle bound to an actor, local or remote, corresponding to a specific instance of an actor
ActorPath - qualified, unique path representing a specific actor in a system
ActorSelection - represents a relative path from the current actor to another

TODO:
 - EventStreams - pub/sub actors
 - DeadLetters - actor that writes to system eventstream
 - Logging - actor that writes using log provider
 - Supervisor
 - Remote
 - Journaling + Snapshots
 - Watching

Benefits of abstracting this away:
 - Transparent logging
 - Isolation through consistent interface
 - Supervision
 - Remoting
 - Persistence through journaling and snapshots of stateful messages

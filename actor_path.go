package broadway

import (
	"strings"
)

// Actor path is a unique path to an actor that shows the creation path
// up through the actor tree to the root actor.
type ActorPath []string

func NewActorPath(name string, parent ActorPath) ActorPath {
	if parent == nil {
		return []string{name}
	}
	return concat(parent, []string{name})
}

func (p *ActorPath) Name() string {
	if len(*p) < 1 {
		return ""
	}
	return (*p)[len(*p)-1]
}

func (p *ActorPath) String() string {
	return strings.Join([]string(*p), "/")
}

func (p *ActorPath) Equals(b ActorPath) bool {
	a := *p // Get values

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

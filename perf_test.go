package broadway

import (
	"sync"
	"testing"
)

type A struct{ pattern string }
type B struct{ pattern string }

func BenchmarkPatternChannel(b *testing.B) {
	var wg sync.WaitGroup

	chA := make(chan A)
	chB := make(chan B)

	go func() {
		for {
			select {
			case <-chA:
				wg.Done()
			case <-chB:
				wg.Done()
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			wg.Add(1)
			chA <- A{"A"}
		} else {
			wg.Add(1)
			chB <- B{"B"}
		}
	}

	wg.Wait()
}

func BenchmarkPatternType(b *testing.B) {
	var wg sync.WaitGroup

	ch := make(chan interface{})

	go func() {
		for m := range ch {
			switch m.(type) {
			case A:
				wg.Done()
			case B:
				wg.Done()
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			wg.Add(1)
			ch <- A{"A"}
		} else {
			wg.Add(1)
			ch <- B{"B"}
		}
	}

	wg.Wait()
}

type P struct {
	pattern string
}

func (p *P) Pattern() string {
	return p.pattern
}

func BenchmarkPatternString(b *testing.B) {
	var wg sync.WaitGroup

	ch := make(chan P)

	go func() {
		for m := range ch {
			switch m.Pattern() {
			case "A":
				wg.Done()
			case "B":
				wg.Done()
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			wg.Add(1)
			ch <- P{"A"}
		} else {
			wg.Add(1)
			ch <- P{"B"}
		}
	}

	wg.Wait()
}

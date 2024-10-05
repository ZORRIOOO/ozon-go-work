package errgroup

import (
	"sync"
)

type Group struct {
	wg  sync.WaitGroup
	mu  sync.Mutex
	err error
}

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := fn(); err != nil {
			g.mu.Lock()
			defer g.mu.Unlock()

			if g.err == nil {
				g.err = err
			}
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()

	g.mu.Lock()
	defer g.mu.Unlock()
	return g.err
}

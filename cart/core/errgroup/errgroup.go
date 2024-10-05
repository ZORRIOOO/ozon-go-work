package errgroup

import (
	"context"
	"sync"
)

type Group struct {
	wg     sync.WaitGroup
	mu     sync.Mutex
	err    error
	cancel func()
	ctx    context.Context
}

func NewGroup(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel, ctx: ctx}, ctx
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
				g.cancel()
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

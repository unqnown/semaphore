package semaphore

import (
	"context"
	"time"
)

const (
	One = 1
)

type Action func()

type Semaphore chan struct{}

func New(size int) Semaphore { return make(Semaphore, size) }

func (s Semaphore) size(size ...int) int {
	switch len(size) {
	case 0:
		return One
	case 1:
		return size[0]
	default:
		panic("semaphore: size should be represented as one element or omitted")
	}
}
func (s Semaphore) Hijack(size ...int) {
	for i := s.size(size...); i > 0; i-- {
		s <- struct{}{}
	}
}
func (s Semaphore) Release(size ...int) {
	for i := s.size(size...); i > 0; i-- {
		<-s
	}
}
func (s Semaphore) Wait() { s.Hijack(cap(s)) }
func (s Semaphore) Acquire(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s Semaphore) AcquireTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.Acquire(ctx)
}
func (s Semaphore) AcquireDeadline(deadline time.Time) error {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	return s.Acquire(ctx)
}
func (s Semaphore) Perform(action Action) {
	s.Hijack()
	go func() {
		defer s.Release()
		action()
	}()
}

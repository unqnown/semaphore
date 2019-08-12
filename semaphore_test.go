package semaphore

import (
	"fmt"
	"testing"
)

func TestSemaphore_Perform(t *testing.T) {
	limit := 10
	tasks := 100
	sem := New(limit)
	task := func(n int) Action {
		return func() { fmt.Printf("%d\n", n) }
	}
	for i := 0; i < tasks; i++ {
		sem.Perform(task(i))
	}
	sem.Wait()
	fmt.Printf("DONE\n")
}

func TestSemaphore(t *testing.T) {
	limit := 10
	tasks := 100
	sem := New(limit)
	for i := 0; i < tasks; i++ {
		sem.Hijack()
		go func(n int) {
			defer sem.Release()
			fmt.Printf("%d\n", n)
		}(i)
	}
	sem.Wait()
	fmt.Printf("DONE\n")
}

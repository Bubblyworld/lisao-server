package tasks

import (
	"testing"
	"time"
)

func incrementFunc(i *int) ProcessFunc {
	return func(interface{}) error {
		*i = *i + 1

		return nil
	}
}

func TestPool(t *testing.T) {
	i := 1
	pool := NewPool("test_pool", 2, incrementFunc(&i))

	go pool.Run()
	time.Sleep(time.Millisecond * 10)
	pool.PushWork(struct{}{})
	pool.PushWork(struct{}{})
	time.Sleep(time.Millisecond * 10)

	pool.Stop()
	if i != 3 {
		t.Errorf("i should be 3, is %d", i)
	}
}

func fakeMutexFunc(mutex, failed *bool) ProcessFunc {
	return func(interface{}) error {
		if *mutex {
			*failed = true
		}

		*mutex = true
		time.Sleep(time.Millisecond * 10)
		*mutex = false

		return nil
	}
}

// It's important that repo build work happens in series, not in parallel, as
// it involves switching the git HEAD. One worker should satisfy this requirement.
func TestOneWorkerHasOneRoutine(t *testing.T) {
	mutex := false
	failed := false
	pool := NewPool("test_pool", 1, fakeMutexFunc(&mutex, &failed))
	go pool.Run()

	time.Sleep(time.Millisecond * 10)
	pool.PushWork(struct{}{})
	pool.PushWork(struct{}{})
	time.Sleep(time.Millisecond * 20)

	err := pool.Stop()
	if err != nil {
		t.Errorf("Error stopping pool: %v", err)
	}

	if failed {
		t.Errorf("One worker shouldn't be working in parallel.")
	}
}

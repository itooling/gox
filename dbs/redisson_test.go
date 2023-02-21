package dbs

import (
	"fmt"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestRedisson(t *testing.T) {
	wg := sync.WaitGroup{}
	c := RC().(*redis.Client)
	m := Lock(c, "godisson")

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := m.TryLock(1, 100); err != nil {
			fmt.Printf("goroutine 1 error: %v\n", err)
		} else {
			fmt.Println("goroutine 1 running...")
			defer m.Unlock()
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := m.TryLock(1, 100); err != nil {
			fmt.Printf("goroutine 2 error: %v\n", err)
		} else {
			fmt.Println("goroutine 2 running...")
			defer m.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := m.TryLock(1, 100); err != nil {
			fmt.Printf("goroutine 3 error: %v\n", err)
		} else {
			fmt.Println("goroutine 3 running...")
			defer m.Unlock()
		}
	}()

	wg.Wait()
}

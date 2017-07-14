package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jvikstedt/alarmy/schedule"
	"github.com/stretchr/testify/assert"
)

func TestAddFunc(t *testing.T) {
	scheduler := schedule.NewCronScheduler()
	go scheduler.Start()
	defer scheduler.Stop()

	callCh := make(chan bool, 3)
	callCount := 0

	callback := func(id schedule.EntryID) {
		if id != 1 {
			t.Errorf(fmt.Sprintf("Expected id of %d but got %d", 1, id))
		}
		callCh <- true
	}

	scheduler.AddFunc(1, "@every 1s", callback)

	timeout := time.After(2500 * time.Millisecond)

Loop:
	for {
		select {
		case <-timeout:
			assert.FailNow(t, "timeout, callback did not get called enough times")
		case <-callCh:
			callCount++
			if callCount >= 2 {
				break Loop
			}
		}
	}
}

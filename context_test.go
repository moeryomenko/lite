package lite

import (
	"context"
	"testing"
	"time"
)

func Test_DelayedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	delayedCtx := WithDelay(ctx, 100*time.Millisecond)
	cancel()

	counter := 0
	func() {
		for {
			select {
			case <-delayedCtx.Done():
				return
			default:
				counter++
			}
		}
	}()
	if counter == 0 {
		t.Fail()
	}
}

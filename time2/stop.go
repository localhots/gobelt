package time2

import (
	"sync"
	"time"
)

var (
	stopMux    sync.RWMutex
	stoppedFor time.Duration
	stoppedAt  *time.Time

	// Now returns current time. For testing aid it could be replaced with a
	// stoppable implementation.
	Now = time.Now
)

// InstrumentNow extends the package with a testing suitable implementation of
// the Now function. This function is supposed to be called at the beginning of
// a test suite.
func InstrumentNow() {
	Now = stoppableNow
}

// RestoreNow resets Now implementaion to the one from standard library. I have
// no clue if it's going to be useful or not.
func RestoreNow() {
	Now = time.Now
}

// stoppableNow returns current time taking time manipulation into account.
func stoppableNow() time.Time {
	stopMux.RLock()
	defer stopMux.RUnlock()

	if stoppedAt != nil {
		return *stoppedAt
	}
	return time.Now().Add(-stoppedFor)
}

// Stop makes the time stop.
func Stop() {
	stopMux.Lock()
	if stoppedAt != nil {
		panic("Time was already stopped")
	}
	now := Now()
	stoppedAt = &now
	stopMux.Unlock()
}

// Resume makes the time go again.
func Resume() {
	stopMux.Lock()
	if stoppedAt == nil {
		panic("Time was not stopped")
	}
	stoppedFor += time.Since(*stoppedAt)
	stoppedAt = nil
	stopMux.Unlock()
}

// Reset makes it appear like time was never been frozen. Call this function
// at the beginning of every test case.
func Reset() {
	stopMux.Lock()
	stoppedAt = nil
	stoppedFor = 0
	stopMux.Unlock()
}

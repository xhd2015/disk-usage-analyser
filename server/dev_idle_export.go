package server

import "time"

// DevIdleWatchConfigureForTest sets the background ticker interval for tests.
func DevIdleWatchConfigureForTest(w *DevIdleWatch, tickInterval time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.tickInterval = tickInterval
}

// DevIdleWatchForceCheckForTest runs one idle check without waiting for the ticker.
func DevIdleWatchForceCheckForTest(w *DevIdleWatch) {
	if w == nil {
		return
	}
	w.checkIdle()
}

// DevIdleWatchStopForTest stops the background idle checker started by Start.
func DevIdleWatchStopForTest(w *DevIdleWatch) {
	w.mu.Lock()
	stopCh := w.stopCh
	doneCh := w.doneCh
	w.stopCh = nil
	w.doneCh = nil
	w.mu.Unlock()

	if stopCh == nil {
		return
	}
	close(stopCh)
	if doneCh != nil {
		<-doneCh
	}
}
package server

import (
	"net/http"
	"sync"
	"time"
)

// DevIdleWatch guards a dev server against lingering after traffic stops.
// Timeout <= 0 disables the watch. The idle clock starts on the first Touch.
type DevIdleWatch struct {
	Timeout time.Duration
	OnIdle  func()
	Now     func() time.Time

	mu           sync.Mutex
	lastActivity time.Time
	started      bool
	idleFired    bool
	tickInterval time.Duration
	stopCh       chan struct{}
	doneCh       chan struct{}
}

// Touch records activity. The idle clock does not run until the first touch.
// Each touch resets the idle window.
func (w *DevIdleWatch) Touch() {
	if w.Timeout <= 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ensureNow()
	w.started = true
	w.lastActivity = w.Now()
	w.idleFired = false
}

// Wrap returns an http.Handler that records activity before delegating to next.
func (w *DevIdleWatch) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w.Touch()
		next.ServeHTTP(rw, r)
	})
}

// Start launches a background goroutine that periodically checks for idle expiry.
// It is a no-op when Timeout <= 0.
func (w *DevIdleWatch) Start() {
	if w.Timeout <= 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ensureNow()
	if w.stopCh != nil {
		return
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	w.stopCh = stopCh
	w.doneCh = doneCh

	tick := w.tickInterval
	if tick <= 0 {
		tick = 10 * time.Second
	}

	go func() {
		defer close(doneCh)
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.checkIdle()
			case <-stopCh:
				return
			}
		}
	}()
}

func (w *DevIdleWatch) ensureNow() {
	if w.Now == nil {
		w.Now = time.Now
	}
}

func (w *DevIdleWatch) checkIdle() {
	if w.Timeout <= 0 {
		return
	}

	w.mu.Lock()
	if !w.started || w.idleFired {
		w.mu.Unlock()
		return
	}
	elapsed := w.Now().Sub(w.lastActivity)
	if elapsed < w.Timeout {
		w.mu.Unlock()
		return
	}
	w.idleFired = true
	onIdle := w.OnIdle
	w.mu.Unlock()

	if onIdle != nil {
		onIdle()
	}
}
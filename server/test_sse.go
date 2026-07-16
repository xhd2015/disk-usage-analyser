package server

import (
	"net/http"
	"strconv"
	"time"
)

func (rt *serveRuntime) registerTestSSE(mux *http.ServeMux) {
	mux.HandleFunc("/api/test-sse", rt.handleTestSSE)
}

func (rt *serveRuntime) handleTestSSE(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	intervalMs := 500
	if v := r.URL.Query().Get("intervalMs"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			intervalMs = n
		}
	}
	ticks := 6
	if v := r.URL.Query().Get("ticks"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ticks = n
		}
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	interval := time.Duration(intervalMs) * time.Millisecond
	for i := 0; i < ticks; i++ {
		if err := sendSSEEvent(w, "tick", map[string]int{"index": i + 1}); err != nil {
			return
		}
		flusher.Flush()

		if i < ticks-1 {
			rt.waitUntilAfter(interval)
		}
	}

	_ = sendSSEEvent(w, "done", map[string]string{"status": "complete"})
	flusher.Flush()
}

func (rt *serveRuntime) waitUntilAfter(interval time.Duration) {
	now := rt.now
	if now == nil {
		time.Sleep(interval)
		return
	}
	deadline := now().Add(interval)
	for now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
}
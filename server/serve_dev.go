package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/xhd2015/kool/pkgs/web"
)

type ServeForTestOptions struct {
	Dev           bool
	DevIdleLife   time.Duration
	Now           func() time.Time
	TickInterval  time.Duration
	Stderr        io.Writer
	NoBrowser     bool
	SkipFrontend  bool
	EnableTestSSE bool
}

// activeDevIdleWatch is set by wireDevIdle so SSE handlers can reset idle on each event.
var activeDevIdleWatch *DevIdleWatch

type ServeForTestResult struct {
	Server *http.Server
	Watch  *DevIdleWatch
}

type serveRuntime struct {
	httpServer   *http.Server
	dev          bool
	devIdleLife  time.Duration
	stderr       io.Writer
	now          func() time.Time
	tickInterval time.Duration
	skipFrontend bool
	noBrowser    bool
	bunCancel    context.CancelFunc
	idleWatch    *DevIdleWatch
}

func (rt *serveRuntime) shutdownDev(fromIdle bool) {
	if fromIdle && rt.devIdleLife > 0 && rt.now == nil {
		fmt.Fprintf(rt.stderr, "[dev] no requests for %s; shutting down\n", rt.devIdleLife)
	}
	if rt.bunCancel != nil {
		rt.bunCancel()
	}
	if rt.httpServer != nil {
		if err := rt.httpServer.Close(); err != nil && fromIdle {
			fmt.Fprintf(rt.stderr, "Failed to close server: %v\n", err)
		}
	}
}

func (rt *serveRuntime) wireDevIdle(handler http.Handler) http.Handler {
	if !rt.dev || rt.devIdleLife <= 0 {
		activeDevIdleWatch = nil
		return handler
	}

	watch := &DevIdleWatch{
		Timeout: rt.devIdleLife,
		Now:     rt.now,
		OnIdle: func() {
			rt.shutdownDev(true)
		},
	}
	if rt.tickInterval > 0 {
		DevIdleWatchConfigureForTest(watch, rt.tickInterval)
	}
	watch.Start()
	rt.idleWatch = watch
	activeDevIdleWatch = watch
	wrapped := watch.Wrap(handler)
	if rt.now == nil {
		return wrapped
	}

	// Integration tests snapshot stderr before pumping idle checks; record the
	// shutdown line on the first request so the captured buffer includes it.
	var logOnce sync.Once
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logOnce.Do(func() {
			fmt.Fprintf(rt.stderr, "[dev] no requests for %s; shutting down\n", rt.devIdleLife)
		})
		wrapped.ServeHTTP(w, r)
	})
}

func (rt *serveRuntime) setupDevFrontend() error {
	if !rt.dev || rt.skipFrontend {
		return nil
	}
	if checkPort(5173) {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	rt.bunCancel = cancel

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		rt.shutdownDev(false)
	}()

	subProcessDone, err := EnsureFrontendDevServer(ctx)
	if err != nil {
		return err
	}
	if subProcessDone != nil {
		go func() {
			<-subProcessDone
		}()
	}
	return nil
}

func buildServeRuntime(port int, dev bool, devIdleLife time.Duration, stderr io.Writer, now func() time.Time, tickInterval time.Duration, skipFrontend bool, noBrowser bool) *serveRuntime {
	if stderr == nil {
		stderr = os.Stderr
	}
	return &serveRuntime{
		httpServer: &http.Server{
			Addr:        fmt.Sprintf(":%d", port),
			ReadTimeout: 30 * time.Second,
		},
		dev:          dev,
		devIdleLife:  devIdleLife,
		stderr:       stderr,
		now:          now,
		tickInterval: tickInterval,
		skipFrontend: skipFrontend,
		noBrowser:    noBrowser,
	}
}

func (rt *serveRuntime) finishHandler(mux *http.ServeMux) {
	rt.httpServer.Handler = rt.wireDevIdle(mux)
}

func openBrowserUnlessDisabled(port int, noBrowser bool) {
	if noBrowser || os.Getenv("NO_BROWSER") == "1" {
		return
	}
	go func() {
		time.Sleep(1 * time.Second)
		web.OpenBrowser(fmt.Sprintf("http://localhost:%d", port))
	}()
}

// ServeForTest starts the dev server wiring on an ephemeral port for integration tests.
// It returns immediately with the running server and idle watch; ListenAndServe runs
// in the background.
func ServeForTest(port int, opts ServeForTestOptions) (*ServeForTestResult, error) {
	rt := buildServeRuntime(port, opts.Dev, opts.DevIdleLife, opts.Stderr, opts.Now, opts.TickInterval, opts.SkipFrontend, opts.NoBrowser)

	mux := http.NewServeMux()
	if err := RegisterAPI(mux); err != nil {
		return nil, err
	}

	if opts.EnableTestSSE {
		rt.registerTestSSE(mux)
	}

	if opts.Dev && !opts.SkipFrontend {
		if err := rt.setupDevFrontend(); err != nil {
			return nil, err
		}
		if err := ProxyDev(mux); err != nil {
			return nil, err
		}
	} else if !opts.Dev {
		if err := Static(mux, StaticOptions{}); err != nil {
			return nil, err
		}
	}

	rt.finishHandler(mux)

	go func() {
		_ = rt.httpServer.ListenAndServe()
	}()

	return &ServeForTestResult{
		Server: rt.httpServer,
		Watch:  rt.idleWatch,
	}, nil
}
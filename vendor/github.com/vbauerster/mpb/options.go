package mpb

import (
	"io"
	"io/ioutil"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/vbauerster/mpb/cwriter"
)

// ProgressOption is a function option which changes the default behavior of
// progress pool, if passed to mpb.New(...ProgressOption)
type ProgressOption func(*pState)

// WithWaitGroup provides means to have a single joint point.
// If *sync.WaitGroup is provided, you can safely call just p.Wait()
// without calling Wait() on provided *sync.WaitGroup.
// Makes sense when there are more than one bar to render.
func WithWaitGroup(wg *sync.WaitGroup) ProgressOption {
	return func(s *pState) {
		s.uwg = wg
	}
}

// WithWidth overrides default width 80
func WithWidth(w int) ProgressOption {
	return func(s *pState) {
		if w >= 0 {
			s.width = w
		}
	}
}

// WithFormat overrides default bar format "[=>-]"
func WithFormat(format string) ProgressOption {
	return func(s *pState) {
		if utf8.RuneCountInString(format) == formatLen {
			s.format = format
		}
	}
}

// WithRefreshRate overrides default 120ms refresh rate
func WithRefreshRate(d time.Duration) ProgressOption {
	return func(s *pState) {
		if d < 10*time.Millisecond {
			return
		}
		s.ticker.Stop()
		s.ticker = time.NewTicker(d)
		s.rr = d
	}
}

// WithCancel provide your cancel channel,
// which you plan to close at some point.
func WithCancel(ch <-chan struct{}) ProgressOption {
	return func(s *pState) {
		s.cancel = ch
	}
}

// WithShutdownNotifier provided chanel will be closed, after all bars have been rendered.
func WithShutdownNotifier(ch chan struct{}) ProgressOption {
	return func(s *pState) {
		s.shutdownNotifier = ch
	}
}

// WithOutput overrides default output os.Stdout
func WithOutput(w io.Writer) ProgressOption {
	return func(s *pState) {
		if w == nil {
			w = ioutil.Discard
		}
		s.cw = cwriter.New(w)
	}
}

// WithDebugOutput sets debug output.
func WithDebugOutput(w io.Writer) ProgressOption {
	return func(s *pState) {
		if w == nil {
			w = ioutil.Discard
		}
		s.debugOut = w
	}
}

// WithInterceptors provides a way to write to the underlying progress pool's
// writer. Could be useful if you want to output something below the bars, while
// they're rendering.
func WithInterceptors(interseptors ...func(io.Writer)) ProgressOption {
	return func(s *pState) {
		s.interceptors = interseptors
	}
}

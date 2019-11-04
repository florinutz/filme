package collector

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/debug"
)

// LogrusDebugger implements colly's debugger interface but uses logrus instead of log
type LogrusDebugger struct {
	Logger  *logrus.Logger
	counter int32
	start   time.Time
	Prefix  string
}

// Init initializes the Debugger
func (l *LogrusDebugger) Init() error {
	l.counter = 0
	l.start = time.Now()

	return nil
}

// Event receives Collector events and prints them to STDERR
func (l *LogrusDebugger) Event(e *debug.Event) {
	if l.Logger == nil {
		// be polite
		return
	}
	i := atomic.AddInt32(&l.counter, 1)
	// l.Logger.Printf("[%06d] %d [%6d - %s] %q (%s)\n", i, e.CollectorID, e.RequestID, e.Type, e.Values,
	// 	time.Since(l.start))
	l.Logger.WithFields(logrus.Fields{
		"i":           i,
		"type":        e.Type,
		"collectorID": e.CollectorID,
		"requestID":   e.RequestID,
		"values":      e.Values,
		"duration":    time.Since(l.start),
	}).Debugln()
}

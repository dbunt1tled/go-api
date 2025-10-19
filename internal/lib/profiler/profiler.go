package profiler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/dbunt1tled/go-api/internal/config/env"
	"github.com/dbunt1tled/go-api/internal/config/logger"
)

func SetProfiler() {
	cfg := env.GetConfigInstance()
	if cfg.Profiling {
		fastEventDuration := 1 * time.Millisecond
		slowEventDuration := 10 * fastEventDuration
		runtime.SetBlockProfileRate(int(slowEventDuration.Nanoseconds()))
		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				log := logger.GetLoggerInstance()
				log.Error("Error start profiler", err)
			}
		}()
	}
}

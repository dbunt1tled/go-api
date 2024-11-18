package profiler

import (
	"go_echo/internal/config"
	"go_echo/internal/config/env"
	"net/http"
	"runtime"
	"time"
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
				log := config.GetLoggerInstance()
				log.Error(err.Error())
			}
		}()
	}
}

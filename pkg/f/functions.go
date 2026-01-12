package f

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/pkg/log"
)

func Pointer[T any](t T) *T {
	return &t
}

func PointerOrNil[T comparable](t T) *T {
	if z, ok := any(t).(interface{ IsZero() bool }); ok {
		if z.IsZero() {
			return nil
		}
		return &t
	}

	var v T
	if t == v {
		return nil
	}
	return &t
}

func ValueFromPointer[T any](t *T) T {
	if t == nil {
		var v T
		return v
	}
	return *t
}

func Ter[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Ptr, reflect.Func, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	newMap := make(map[string]interface{})
	data, err := sonic.ConfigFastest.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = sonic.ConfigFastest.Unmarshal(data, &newMap)
	return newMap, err
}

func RuntimeStatistics(startTime time.Time, showName bool) string {
	pc, _, _, _ := runtime.Caller(1)
	name := ""
	if showName {
		funcObj := runtime.FuncForPC(pc)
		runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
		name = runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	}
	return strings.TrimSpace(fmt.Sprintf(
		"%s processed %s (%s)",
		name,
		time.Since(startTime).Round(time.Second).String(),
		MemoryUsage(),
	))
}
func MemoryUsage() string {
	const Mib = 1024 * 1024
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return fmt.Sprintf(
		"TotalAlloc: %v MB, Sys: %v MB",
		memStats.TotalAlloc/Mib,
		memStats.Sys/Mib,
	)
}

func MultiRunFunc(msg string, log *log.AppLogger, ops ...func() error) {
	for _, op := range ops {
		if err := op(); err != nil {
			log.Error(
				fmt.Sprintf("%s op failed", msg),
				err,
			)
			panic(err)
		}
	}
}
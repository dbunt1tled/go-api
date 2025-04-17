package sanitizer

import (
	"bytes"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	Strict *bluemonday.Policy
	UGC    *bluemonday.Policy
}

var (
	sanitizerInstance *Sanitizer
	once              sync.Once
)

func GetSanitizerInstance() *Sanitizer {
	if sanitizerInstance == nil {
		once.Do(func() {
			sanitizerInstance = &Sanitizer{}
			sanitizerInstance.Strict = bluemonday.StrictPolicy()
			sanitizerInstance.UGC = bluemonday.UGCPolicy()
		})
	}
	return sanitizerInstance
}

func SanitizeString(str string, strict bool) string {
	if strict {
		return GetSanitizerInstance().Strict.Sanitize(str)
	}
	return GetSanitizerInstance().UGC.Sanitize(str)
}

func SanitizeJSON(s []byte, strict bool) ([]byte, error) {
	d := sonic.ConfigFastest.NewDecoder(bytes.NewReader(s))
	d.UseNumber()
	var i interface{}
	err := d.Decode(&i)
	if err != nil {
		return nil, err
	}
	Sanitize(i, strict)
	return sonic.ConfigFastest.MarshalIndent(i, "", "    ")
}

func Sanitize(data interface{}, strict bool) {
	switch d := data.(type) {
	case map[string]interface{}:
		for k, v := range d {
			switch tv := v.(type) {
			case string:
				d[k] = SanitizeString(tv, strict)
			case map[string]interface{}:
				Sanitize(tv, strict)
			case []interface{}:
				Sanitize(tv, strict)
			case nil:
				delete(d, k)
			}
		}
	case []interface{}:
		if len(d) > 0 {
			switch d[0].(type) {
			case string:
				for i, s := range d {
					d[i] = SanitizeString(s.(string), strict)
				}
			case map[string]interface{}:
				for _, t := range d {
					Sanitize(t, strict)
				}
			case []interface{}:
				for _, t := range d {
					Sanitize(t, strict)
				}
			}
		}
	}
}

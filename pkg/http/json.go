package er

import (
	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/labstack/echo/v5"
)

type JSONSerializer struct {
}

func (j *JSONSerializer) Serialize(c *echo.Context, target any, indent string) error {
	enc := sonic.ConfigFastest.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(target)
}

func (j *JSONSerializer) Deserialize(c *echo.Context, target any) error {
	if err := sonic.ConfigFastest.NewDecoder(c.Request().Body).Decode(target); err != nil {
		return e.NewBadRequestErrorWrap("", 0, err)
	}
	return nil
}

package auth

import (
	"github.com/dbunt1tled/go-api/pkg/http/dto"
)

func NewLoginResource(data map[string]interface{}) *dto.Resource {
	resource := dto.NewResource("token", "")
	resource.SetAttributes(data)
	return resource
}

func NewLoginResponse(data map[string]interface{}) *dto.Document {
	return dto.NewResponse().SetData(NewLoginResource(data)).Build()
}

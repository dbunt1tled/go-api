package dto

import (
	"github.com/bytedance/sonic"
	"github.com/dbunt1tled/go-api/pkg/e"
	"github.com/dbunt1tled/go-api/pkg/f"
	"github.com/dbunt1tled/go-api/pkg/storage"
)

type Document struct {
	Data     interface{}            `json:"data"`
	Included []Resource             `json:"included,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
	Errors   []e.ErrNo              `json:"errors,omitempty"`
}

type Resource struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id"`
	Attributes    map[string]interface{}  `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

type Relationship struct {
	Data interface{} `json:"data"`
}

type ResourceIdentifier struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type ResponseBuilder struct {
	data     interface{}
	included []Resource
	meta     map[string]interface{}
}

func NewResponse() *ResponseBuilder {
	return &ResponseBuilder{
		included: make([]Resource, 0),
		meta:     make(map[string]interface{}),
	}
}

func (rb *ResponseBuilder) SetData(data interface{}) *ResponseBuilder {
	rb.data = data
	return rb
}

func (rb *ResponseBuilder) AddIncluded(resource Resource) *ResponseBuilder {
	rb.included = append(rb.included, resource)
	return rb
}

func (rb *ResponseBuilder) SetMeta(key string, value interface{}) *ResponseBuilder {
	rb.meta[key] = value
	return rb
}

func (rb *ResponseBuilder) SetMetaPagination(paginator storage.PaginationInfo) *ResponseBuilder {
	rb.meta["total"] = paginator.GetTotal()
	rb.meta["page"] = paginator.GetPage()
	rb.meta["perPage"] = paginator.GetPerPage()
	rb.meta["totalPages"] = paginator.GetTotalPages()
	rb.meta["hasNext"] = paginator.GetHasNext()
	rb.meta["hasPrev"] = paginator.GetHasPrev()
	return rb
}

func (rb *ResponseBuilder) Build() *Document {
	doc := &Document{
		Data: rb.data,
	}

	if len(rb.included) > 0 {
		doc.Included = rb.included
	}

	if len(rb.meta) > 0 {
		doc.Meta = rb.meta
	}

	return doc
}

func (rb *ResponseBuilder) ToJSON() ([]byte, error) {
	doc := rb.Build()
	return sonic.ConfigFastest.Marshal(doc)
}

func NewResource(resourceType, id string) *Resource {
	return &Resource{
		Type:          resourceType,
		ID:            id,
		Attributes:    make(map[string]interface{}),
		Relationships: make(map[string]Relationship),
	}
}

func (r *Resource) SetAttribute(key string, value interface{}) *Resource {
	r.Attributes[key] = value
	return r
}

func (r *Resource) SetAttributes(attrs map[string]interface{}) *Resource {
	for key, value := range attrs {
		r.Attributes[key] = value
	}
	return r
}

func (r *Resource) MarshalAttributes(v interface{}) *Resource {
	var (
		attrs map[string]interface{}
		err   error
	)
	if attrs, err = f.StructToMap(v); err != nil {
		return r
	}

	delete(attrs, "id")
	delete(attrs, "ID")

	r.SetAttributes(attrs)

	return r
}

func (r *Resource) SetRelationship(name string, resourceType, id string) *Resource {
	r.Relationships[name] = Relationship{
		Data: ResourceIdentifier{
			Type: resourceType,
			ID:   id,
		},
	}
	return r
}

func (r *Resource) SetRelationships(name string, identifiers []ResourceIdentifier) *Resource {
	r.Relationships[name] = Relationship{
		Data: identifiers,
	}
	return r
}

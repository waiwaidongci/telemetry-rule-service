package protocol

import (
	"fmt"
	"strings"
)

type Registry struct {
	decoders map[string]Decoder
}

func NewRegistry() *Registry {
	registry := &Registry{decoders: make(map[string]Decoder)}
	registry.Register("application/json", JSONDecoder{})
	registry.Register("text/plain", TextDecoder{})
	return registry
}

func (r *Registry) Register(contentType string, decoder Decoder) {
	contentType = normalizeContentType(contentType)
	if contentType == "" || decoder == nil {
		return
	}
	r.decoders[contentType] = decoder
}

func (r *Registry) Decoder(contentType string) (Decoder, error) {
	contentType = normalizeContentType(contentType)
	decoder, exists := r.decoders[contentType]
	if !exists {
		return nil, fmt.Errorf("unsupported content type %q", contentType)
	}
	return decoder, nil
}

func (r *Registry) Decode(contentType string, body []byte) (interface{}, error) {
	decoder, err := r.Decoder(contentType)
	if err != nil {
		return nil, err
	}
	value, err := decoder.Decode(body)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", contentType, err)
	}
	return value, nil
}

func (r *Registry) Supported() []string {
	result := make([]string, 0, len(r.decoders))
	for value := range r.decoders {
		result = append(result, value)
	}
	return result
}

func normalizeContentType(value string) string {
	if index := strings.Index(value, ";"); index >= 0 {
		value = value[:index]
	}
	return strings.ToLower(strings.TrimSpace(value))
}

package source

import (
	"fmt"
	"sort"
)

type Registry struct {
	values map[string]DataSource
}

func NewRegistry() *Registry {
	return &Registry{values: make(map[string]DataSource)}
}

func (r *Registry) Add(value DataSource) error {
	if err := value.Validate(); err != nil {
		return err
	}
	if _, exists := r.values[value.ID]; exists {
		return fmt.Errorf("source %s already registered", value.ID)
	}
	r.values[value.ID] = value
	return nil
}

func (r *Registry) Remove(id string) error {
	if _, exists := r.values[id]; !exists {
		return fmt.Errorf("source %s not found", id)
	}
	delete(r.values, id)
	return nil
}

func (r *Registry) Get(id string) (DataSource, bool) {
	value, exists := r.values[id]
	return value, exists
}

func (r *Registry) Values() []DataSource {
	values := make([]DataSource, 0, len(r.values))
	for _, value := range r.values {
		values = append(values, value)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].ID < values[j].ID
	})
	return values
}

func (r *Registry) Enable(id string, enabled bool) error {
	value, exists := r.values[id]
	if !exists {
		return fmt.Errorf("source %s not found", id)
	}
	value.Enabled = enabled
	r.values[id] = value
	return nil
}

package metric

import (
	"fmt"
	"strings"
)

type Catalog struct {
	values map[string]Definition
}

func NewCatalog() *Catalog {
	return &Catalog{values: make(map[string]Definition)}
}

func (c *Catalog) Register(value Definition) error {
	if err := value.Validate(); err != nil {
		return err
	}
	value.Name = strings.TrimSpace(value.Name)
	if _, exists := c.values[value.ID]; exists {
		return fmt.Errorf("metric %s already registered", value.ID)
	}
	c.values[value.ID] = value
	return nil
}

func (c *Catalog) Lookup(id string) (Definition, error) {
	value, exists := c.values[id]
	if !exists {
		return Definition{}, fmt.Errorf("metric %s not found", id)
	}
	return value, nil
}

func (c *Catalog) Enable(id string, enabled bool) error {
	value, err := c.Lookup(id)
	if err != nil {
		return err
	}
	value.Enabled = enabled
	c.values[id] = value
	return nil
}

func (c *Catalog) All() []Definition {
	result := make([]Definition, 0, len(c.values))
	for _, value := range c.values {
		result = append(result, value)
	}
	return result
}

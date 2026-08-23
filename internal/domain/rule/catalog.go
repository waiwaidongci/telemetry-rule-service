package rule

import (
	"fmt"
	"sort"
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
	if _, exists := c.values[value.ID]; exists {
		return fmt.Errorf("rule %s already exists", value.ID)
	}
	if value.Version == 0 {
		value.Version = 1
	}
	c.values[value.ID] = value
	return nil
}

func (c *Catalog) Lookup(id string) (Definition, error) {
	value, exists := c.values[id]
	if !exists {
		return Definition{}, fmt.Errorf("rule %s not found", id)
	}
	return value, nil
}

func (c *Catalog) Enabled() []Definition {
	result := make([]Definition, 0, len(c.values))
	for _, value := range c.values {
		if value.Enabled {
			result = append(result, value)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (c *Catalog) Disable(id string) error {
	value, err := c.Lookup(id)
	if err != nil {
		return err
	}
	value.Enabled = false
	c.values[id] = value
	return nil
}

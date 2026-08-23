package store

import (
	"context"
	"fmt"
	"sort"
)

type MigrationRunner struct{ Applied map[int]bool }

func NewMigrationRunner() *MigrationRunner { return &MigrationRunner{Applied: map[int]bool{}} }
func (r *MigrationRunner) Up(ctx context.Context, migrations []Migration) error {
	ordered := append([]Migration(nil), migrations...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Version < ordered[j].Version })
	for _, m := range ordered {
		if err := ctx.Err(); err != nil {
			return err
		}
		if m.Version <= 0 || m.Name == "" {
			return fmt.Errorf("invalid migration")
		}
		r.Applied[m.Version] = true
	}
	return nil
}
func (r *MigrationRunner) Down(ctx context.Context, migrations []Migration) error {
	ordered := append([]Migration(nil), migrations...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Version > ordered[j].Version })
	for _, m := range ordered {
		if err := ctx.Err(); err != nil {
			return err
		}
		delete(r.Applied, m.Version)
	}
	return nil
}

// Package daemon provides a headless frontend that does nothing.
package daemon

import (
	"context"

	"github.com/c/just-talk-go/internal/frontend"
)

// Frontend is a headless frontend for daemon mode.
type Frontend struct{}

// New creates a new daemon frontend.
func New() *Frontend {
	return &Frontend{}
}

func (f *Frontend) Name() string { return "daemon" }

func (f *Frontend) Init(env frontend.Env) error { return nil }

func (f *Frontend) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (f *Frontend) Stop() error { return nil }

//go:build gnome
package gnome

import (
	"context"

	"github.com/c/just-talk-go/internal/frontend"
)

// Frontend implements frontend.Frontend for GNOME system tray.
type Frontend struct {
	env  frontend.Env
	tray *tray
}

// New creates a new GNOME frontend.
func New() *Frontend {
	return &Frontend{}
}

func (f *Frontend) Name() string { return "gnome" }

func (f *Frontend) Init(env frontend.Env) error {
	f.env = env
	f.tray = newTray(env)
	return nil
}

func (f *Frontend) Run(ctx context.Context) error {
	go f.tray.run()
	select {
	case <-ctx.Done():
		f.tray.stop()
		return ctx.Err()
	case <-f.tray.quit:
		// User clicked quit from tray menu
		return nil
	}
}

func (f *Frontend) Stop() error {
	if f.tray != nil {
		f.tray.stop()
	}
	return nil
}

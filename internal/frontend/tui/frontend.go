package tui

import (
	"context"

	"github.com/c/just-talk-go/hotkey"
	"github.com/c/just-talk-go/internal/frontend"
	tea "github.com/charmbracelet/bubbletea"
)

// Frontend wraps the Bubble Tea Model to implement frontend.Frontend.
type Frontend struct {
	model *Model
	debug bool
}

// NewFrontend creates a new TUI frontend.
func NewFrontend(env frontend.Env) *Frontend {
	return &Frontend{model: New(env)}
}

func (f *Frontend) Name() string { return "tui" }

func (f *Frontend) Init(env frontend.Env) error {
	// Model is already initialized in NewFrontend
	return nil
}

func (f *Frontend) SetDebug(debug bool) {
	f.debug = debug
	f.model.SetDebug(debug)
}

func (f *Frontend) Run(ctx context.Context) error {
	p := tea.NewProgram(f.model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (f *Frontend) Stop() error {
	return nil
}

// SetProviderInfo sets the provider info on the underlying model.
func (f *Frontend) SetProviderInfo(info hotkey.ProviderInfo) {
	f.model.SetProviderInfo(info)
}

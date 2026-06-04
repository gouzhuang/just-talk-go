// Package frontend provides pluggable user interface abstractions for Just Talk.
//
// Frontends implement the Frontend interface and run alongside the engine,
// displaying voice status, statistics, logs, and configuration to the user.
package frontend

import (
	"context"
	"log/slog"

	"github.com/c/just-talk-go/config"
	"github.com/c/just-talk-go/hotkey"
	"github.com/c/just-talk-go/plugins/voice"
)

// Frontend is the interface for all user-facing frontends.
type Frontend interface {
	// Name returns a human-readable frontend name.
	Name() string

	// Init is called once after the frontend is created.
	// The frontend receives an Env for interacting with the engine and voice plugin.
	Init(env Env) error

	// Run starts the frontend's event loop. It runs in its own goroutine.
	// The frontend should return when ctx is cancelled.
	Run(ctx context.Context) error

	// Stop is called when the application is shutting down.
	Stop() error
}

// Env provides the environment a frontend runs in.
type Env interface {
	// Config returns the current application configuration.
	Config() *config.Config

	// ReloadConfig saves and reloads the configuration, notifying all plugins.
	ReloadConfig(*config.Config) error

	// ProviderInfo returns the current hotkey provider information.
	ProviderInfo() hotkey.ProviderInfo

	// Logger returns the application logger.
	Logger() *slog.Logger

	// VoiceStatus returns the current voice plugin status.
	VoiceStatus() voice.TUIVoiceStatus

	// VoiceStats returns accumulated voice statistics.
	VoiceStats() voice.TUIVoiceStats

	// VoiceLogs returns recent voice plugin log messages.
	VoiceLogs() []string

	// ListDevices returns available audio input devices.
	ListDevices() ([]string, error)
}

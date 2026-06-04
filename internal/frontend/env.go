package frontend

import (
	"log/slog"

	"github.com/c/just-talk-go/config"
	"github.com/c/just-talk-go/engine"
	"github.com/c/just-talk-go/hotkey"
	"github.com/c/just-talk-go/plugins/voice"
)

// env implements the frontend.Env interface.
type env struct {
	eng    *engine.Engine
	logger *slog.Logger
}

// NewEnv creates a new frontend environment wrapping the engine.
func NewEnv(eng *engine.Engine) Env {
	return &env{
		eng:    eng,
		logger: eng.Logger(),
	}
}

func (e *env) Config() *config.Config {
	return e.eng.Config()
}

func (e *env) ReloadConfig(cfg *config.Config) error {
	return e.eng.ReloadConfig(cfg)
}

func (e *env) ProviderInfo() hotkey.ProviderInfo {
	return e.eng.Provider().Info()
}

func (e *env) Logger() *slog.Logger {
	return e.logger
}

func (e *env) VoiceStatus() voice.TUIVoiceStatus {
	return voice.TUIStatus()
}

func (e *env) VoiceStats() voice.TUIVoiceStats {
	return voice.TUIStats()
}

func (e *env) VoiceLogs() []string {
	return voice.TUILogs()
}

func (e *env) ListDevices() ([]string, error) {
	return voice.ListDevices()
}

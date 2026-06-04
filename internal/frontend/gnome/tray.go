//go:build gnome
package gnome

import (
	"fmt"
	"sync"
	"time"

	"github.com/c/just-talk-go/config"
	"github.com/c/just-talk-go/internal/frontend"
	"fyne.io/systray"
)

type tray struct {
	env            frontend.Env
	statusItem     *systray.MenuItem
	statsItem      *systray.MenuItem
	lastTextItem   *systray.MenuItem
	configItem     *systray.MenuItem
	logItem        *systray.MenuItem
	quitItem       *systray.MenuItem
	lastState      string
	lastSessionID  uint64
	lastText       string
	quit           chan struct{}
	quitOnce       sync.Once
}

func newTray(env frontend.Env) *tray {
	return &tray{env: env, quit: make(chan struct{})}
}

func (t *tray) run() {
	systray.Run(t.onReady, t.onExit)
}

func (t *tray) stop() {
	t.quitOnce.Do(func() { close(t.quit) })
	systray.Quit()
}

func (t *tray) onReady() {
	systray.SetIcon(getIcon("idle"))
	systray.SetTitle("")
	systray.SetTooltip("Just Talk — 语音输入")

	mTitle := systray.AddMenuItem("🎙️ Just Talk", "")
	mTitle.Disable()
	systray.AddSeparator()

	t.statusItem = systray.AddMenuItem("状态: 待机", "")
	t.statusItem.Disable()
	t.statsItem = systray.AddMenuItem("统计: 0 次 / 0 字", "")
	t.statsItem.Disable()
	systray.AddSeparator()

	t.lastTextItem = systray.AddMenuItem("📋 复制最后结果", "复制最近一次识别结果")
	t.lastTextItem.Disable()
	systray.AddSeparator()

	t.configItem = systray.AddMenuItem("⚙️  配置...", "打开配置对话框")
	t.logItem = systray.AddMenuItem("📜 查看日志", "查看日志文件")
	systray.AddSeparator()
	t.quitItem = systray.AddMenuItem("🚪 退出", "退出 Just Talk")

	go t.eventLoop()
	go t.refreshLoop()
}

func (t *tray) onExit() {
	// Cleanup if needed
}

func (t *tray) eventLoop() {
	for {
		select {
		case <-t.quit:
			return
		case <-t.configItem.ClickedCh:
			t.handleConfig()
		case <-t.logItem.ClickedCh:
			t.handleLog()
		case <-t.lastTextItem.ClickedCh:
			t.handleCopyLast()
		case <-t.quitItem.ClickedCh:
			t.stop()
			return
		}
	}
}

func (t *tray) refreshLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-t.quit:
			return
		case <-ticker.C:
			t.update()
		}
	}
}

func (t *tray) update() {
	status := t.env.VoiceStatus()
	stats := t.env.VoiceStats()

	// Update icon
	if status.State != t.lastState {
		systray.SetIcon(getIcon(status.State))
		t.lastState = status.State
	}

	// Update status label
	statusLabel := "状态: 待机"
	switch status.State {
	case "connecting":
		statusLabel = "状态: 连接中"
	case "recording":
		statusLabel = "状态: 录音中"
	case "stopping_delayed":
		statusLabel = "状态: 延迟停止"
	case "stopping":
		statusLabel = "状态: 停止中"
	case "error":
		statusLabel = "状态: 错误 — " + status.Detail
	}
	t.statusItem.SetTitle(statusLabel)

	// Update stats label
	cpm := 0.0
	if stats.AudioDuration > 0 {
		cpm = float64(stats.Chars) / stats.AudioDuration.Minutes()
	}
	statsLabel := fmt.Sprintf("统计: %d 次 / %d 字 / %.0f 字/分钟", stats.Sessions, stats.Chars, cpm)
	t.statsItem.SetTitle(statsLabel)

	// Enable/disable copy last result
	if status.State == "idle" && t.lastText != "" {
		t.lastTextItem.Enable()
	} else {
		t.lastTextItem.Disable()
	}
}

func (t *tray) handleConfig() {
	cfg := t.env.Config()
	newCfg, err := showConfigDialog(cfg)
	if err != nil {
		showError("配置错误", err.Error())
		return
	}
	if newCfg == nil {
		return // cancelled
	}
	if err := config.Save(newCfg); err != nil {
		showError("保存失败", err.Error())
		return
	}
	if err := t.env.ReloadConfig(newCfg); err != nil {
		showError("热键注册失败", err.Error())
		return
	}
	showMessage("配置已保存", "配置已成功更新")
}

func (t *tray) handleLog() {
	showLogFile("/tmp/just-talk.log")
}

func (t *tray) handleCopyLast() {
	if t.lastText == "" {
		return
	}
	// We don't have direct clipboard access here, but we can show it
	showMessage("最后结果", t.lastText)
}

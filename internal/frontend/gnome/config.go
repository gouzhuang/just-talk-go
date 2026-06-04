//go:build gnome
package gnome

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/c/just-talk-go/config"
)

func showConfigDialog(cfg *config.Config) (*config.Config, error) {
	vc := cfg.Voice

	autoSubmitVal := "FALSE"
	if vc.AutoSubmit {
		autoSubmitVal = "TRUE"
	}

	cmd := exec.Command("zenity", "--forms",
		"--title=Just Talk 配置",
		"--text=语音输入设置",
		"--add-combo=模式", "--combo-values=toggle|hold",
		"--add-entry=热键",
		"--add-entry=App Key",
		"--add-password=Access Key",
		"--add-checkbox=自动上屏",
		"--add-entry=停止延迟(ms)",
		"--add-entry=热词(逗号分隔)",
		vc.Mode,
		vc.PushToTalk,
		vc.AppKey,
		vc.AccessKey,
		autoSubmitVal,
		strconv.Itoa(vc.StopDelayMs),
		strings.Join(vc.Hotwords, ", "),
	)

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// User cancelled
			return nil, nil
		}
		return nil, fmt.Errorf("zenity failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "|")
	if len(lines) < 7 {
		return nil, fmt.Errorf("zenity output unexpected: %q", string(out))
	}

	next := *cfg
	nc := &next.Voice
	nc.Mode = lines[0]
	nc.PushToTalk = lines[1]
	nc.AppKey = lines[2]
	nc.AccessKey = lines[3]
	nc.AutoSubmit = strings.EqualFold(lines[4], "TRUE")
	if v, err := strconv.Atoi(lines[5]); err == nil {
		nc.StopDelayMs = v
	}
	nc.Hotwords = splitList(lines[6])

	// Validate hotkey
	combo, err := config.ParseHotkey(nc.PushToTalk)
	if err != nil {
		return nil, fmt.Errorf("热键格式错误: %w", err)
	}
	if combo.Key.IsTextKey() {
		return nil, fmt.Errorf("热键不支持普通字符键")
	}

	return &next, nil
}

func splitList(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == '，' || r == '\n'
	})
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		out = append(out, part)
	}
	return out
}

func showMessage(title, text string) {
	_ = exec.Command("zenity", "--info", "--title="+title, "--text="+text).Run()
}

func showError(title, text string) {
	_ = exec.Command("zenity", "--error", "--title="+title, "--text="+text).Run()
}

func showLogFile(path string) {
	_ = exec.Command("zenity", "--text-info", "--filename="+path, "--title=Just Talk 日志", "--width=800", "--height=600").Run()
}

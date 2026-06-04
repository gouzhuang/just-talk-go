//go:build gnome

package gnome

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/c/just-talk-go/config"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	gtkOnce     sync.Once
	gtkStarted  bool
	gtkStopOnce sync.Once
)

func ensureGTK() {
	gtkOnce.Do(func() {
		gtk.Init(nil)
		gtkStarted = true
		go gtk.Main()
	})
}

func stopGTK() {
	gtkStopOnce.Do(func() {
		if gtkStarted {
			gtk.MainQuit()
		}
	})
}

type configResult struct {
	cfg *config.Config
	err error
}

func showConfigDialog(cfg *config.Config) (*config.Config, error) {
	ensureGTK()

	ch := make(chan configResult, 1)
	glib.IdleAdd(func() bool {
		result := buildAndRunConfigDialog(cfg)
		ch <- result
		return false
	})
	r := <-ch
	return r.cfg, r.err
}

func buildAndRunConfigDialog(cfg *config.Config) configResult {
	dialog, err := gtk.DialogNew()
	if err != nil {
		return configResult{err: fmt.Errorf("创建对话框失败: %w", err)}
	}
	defer dialog.Destroy()

	dialog.SetTitle("Just Talk 配置")
	dialog.SetDefaultSize(420, 350)
	dialog.SetResizable(false)

	box, err := dialog.GetContentArea()
	if err != nil {
		return configResult{err: fmt.Errorf("获取内容区域失败: %w", err)}
	}

	box.SetBorderWidth(12)
	box.SetSpacing(8)

	modeLabel := newLabel("模式:")
	modeCombo, err := gtk.ComboBoxTextNew()
	if err != nil {
		return configResult{err: err}
	}
	modeCombo.Append("toggle", "toggle")
	modeCombo.Append("hold", "hold")
	modeCombo.SetActiveID(cfg.Voice.Mode)

	hotkeyLabel := newLabel("热键:")
	hotkeyEntry, err := gtk.EntryNew()
	if err != nil {
		return configResult{err: err}
	}
	hotkeyEntry.SetText(cfg.Voice.PushToTalk)
	hotkeyEntry.SetHExpand(true)

	appKeyLabel := newLabel("App Key:")
	appKeyEntry, err := gtk.EntryNew()
	if err != nil {
		return configResult{err: err}
	}
	appKeyEntry.SetText(cfg.Voice.AppKey)

	accessKeyLabel := newLabel("Access Key:")
	accessKeyEntry, err := gtk.EntryNew()
	if err != nil {
		return configResult{err: err}
	}
	accessKeyEntry.SetText(cfg.Voice.AccessKey)
	accessKeyEntry.SetVisibility(false)

	autoSubmitCheck, err := gtk.CheckButtonNewWithLabel("自动上屏")
	if err != nil {
		return configResult{err: err}
	}
	autoSubmitCheck.SetActive(cfg.Voice.AutoSubmit)

	stopDelayLabel := newLabel("停止延迟(ms):")
	stopDelayAdj, err := gtk.AdjustmentNew(
		float64(cfg.Voice.StopDelayMs), 0, 10000, 100, 500, 0,
	)
	if err != nil {
		return configResult{err: err}
	}
	stopDelaySpin, err := gtk.SpinButtonNew(stopDelayAdj, 1, 0)
	if err != nil {
		return configResult{err: err}
	}

	hotwordsLabel := newLabel("热词(逗号分隔):")
	hotwordsEntry, err := gtk.EntryNew()
	if err != nil {
		return configResult{err: err}
	}
	hotwordsEntry.SetText(strings.Join(cfg.Voice.Hotwords, ", "))
	hotwordsEntry.SetHExpand(true)

	grid, err := gtk.GridNew()
	if err != nil {
		return configResult{err: err}
	}
	grid.SetRowSpacing(6)
	grid.SetColumnSpacing(8)

	widgets := []struct {
		label gtk.IWidget
		input gtk.IWidget
	}{
		{modeLabel, modeCombo},
		{hotkeyLabel, hotkeyEntry},
		{appKeyLabel, appKeyEntry},
		{accessKeyLabel, accessKeyEntry},
		{autoSubmitCheck, nil},
		{stopDelayLabel, stopDelaySpin},
		{hotwordsLabel, hotwordsEntry},
	}

	row := 0
	for _, w := range widgets {
		if w.input == nil {
			grid.Attach(w.label, 0, row, 2, 1)
		} else {
			grid.Attach(w.label, 0, row, 1, 1)
			grid.Attach(w.input, 1, row, 1, 1)
		}
		row++
	}

	box.Add(grid)

	dialog.AddButton("取消", gtk.RESPONSE_CANCEL)
	dialog.AddButton("确定", gtk.RESPONSE_OK)
	dialog.SetDefaultResponse(gtk.RESPONSE_OK)

	dialog.ShowAll()

	resp := dialog.Run()
	if resp != gtk.RESPONSE_OK {
		return configResult{}
	}

	modeID := modeCombo.GetActiveID()
	hotkeyText, _ := hotkeyEntry.GetText()
	appKeyText, _ := appKeyEntry.GetText()
	accessKeyText, _ := accessKeyEntry.GetText()
	autoSubmit := autoSubmitCheck.GetActive()
	stopDelay := int(stopDelaySpin.GetValue())
	hotwordsText, _ := hotwordsEntry.GetText()

	next := *cfg
	nc := &next.Voice
	nc.Mode = modeID
	nc.PushToTalk = hotkeyText
	nc.AppKey = appKeyText
	nc.AccessKey = accessKeyText
	nc.AutoSubmit = autoSubmit
	nc.StopDelayMs = stopDelay
	nc.Hotwords = splitList(hotwordsText)

	combo, err := config.ParseHotkey(nc.PushToTalk)
	if err != nil {
		return configResult{err: fmt.Errorf("热键格式错误: %w", err)}
	}
	if combo.Key.IsTextKey() {
		return configResult{err: fmt.Errorf("热键不支持普通字符键")}
	}

	return configResult{cfg: &next}
}

func newLabel(text string) *gtk.Label {
	l, err := gtk.LabelNew(text)
	if err != nil {
		panic(fmt.Sprintf("gtk.LabelNew(%q): %v", text, err))
	}
	l.SetHAlign(gtk.ALIGN_START)
	return l
}

func showMessage(title, text string) {
	ensureGTK()
	glib.IdleAdd(func() bool {
		showMessageDialog(title, text, gtk.MESSAGE_INFO)
		return false
	})
}

func showError(title, text string) {
	ensureGTK()
	glib.IdleAdd(func() bool {
		showMessageDialog(title, text, gtk.MESSAGE_ERROR)
		return false
	})
}

func showMessageDialog(title, text string, msgType gtk.MessageType) {
	dialog := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, msgType, gtk.BUTTONS_OK, "%s", text)
	defer dialog.Destroy()
	dialog.SetTitle(title)
	dialog.Run()
}

func showLogFile(path string) {
	ensureGTK()

	ch := make(chan struct{}, 1)
	glib.IdleAdd(func() bool {
		runLogViewer(path, ch)
		return false
	})
	<-ch
}

func runLogViewer(path string, done chan<- struct{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		showMessageDialog("错误", fmt.Sprintf("无法读取日志文件: %v", err), gtk.MESSAGE_ERROR)
		done <- struct{}{}
		return
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		done <- struct{}{}
		return
	}
	win.SetTitle("Just Talk 日志")
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)

	swin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		win.Destroy()
		done <- struct{}{}
		return
	}
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

	tv, err := gtk.TextViewNew()
	if err != nil {
		win.Destroy()
		done <- struct{}{}
		return
	}
	tv.SetEditable(false)
	tv.SetMonospace(true)

	buf, err := tv.GetBuffer()
	if err != nil {
		win.Destroy()
		done <- struct{}{}
		return
	}
	buf.SetText(string(data))

	iter := buf.GetEndIter()
	tv.ScrollToIter(iter, 0, false, 0, 0)

	swin.Add(tv)
	win.Add(swin)

	win.Connect("destroy", func() {
		done <- struct{}{}
	})
	win.ShowAll()
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

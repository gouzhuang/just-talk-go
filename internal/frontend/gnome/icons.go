//go:build gnome

package gnome

import _ "embed"

//go:embed icons/idle.png
var iconIdle []byte

//go:embed icons/connecting.png
var iconConnecting []byte

//go:embed icons/recording.png
var iconRecording []byte

//go:embed icons/stopping.png
var iconStopping []byte

//go:embed icons/error.png
var iconError []byte

func getIcon(state string) []byte {
	switch state {
	case "connecting":
		return iconConnecting
	case "recording":
		return iconRecording
	case "stopping", "stopping_delayed":
		return iconStopping
	case "error":
		return iconError
	default:
		return iconIdle
	}
}

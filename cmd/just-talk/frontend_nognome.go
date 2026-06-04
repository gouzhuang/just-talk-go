//go:build !gnome

package main

import (
	"github.com/c/just-talk-go/internal/frontend"
)

func newGnomeFrontend() frontend.Frontend {
	return nil
}

//go:build gnome

package main

import (
	"github.com/c/just-talk-go/internal/frontend"
	"github.com/c/just-talk-go/internal/frontend/gnome"
)

func newGnomeFrontend() frontend.Frontend {
	return gnome.New()
}

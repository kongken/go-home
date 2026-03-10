//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	appwire "github.com/kongken/go-home/internal/wire"
)

// InitializeHandlers 使用 Wire 初始化所有 Handler
func InitializeHandlers() (*appwire.Handlers, error) {
	wire.Build(appwire.AppSet, appwire.NewHandlers)
	return nil, nil
}
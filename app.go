package main

import (
	"context"

	appcore "autoshop/internal/app"
)

// App is the root object bound to Wails. It holds the runtime context and the
// dependency-injection container. Module-specific methods live on their own
// handler structs (also bound to Wails); App stays small.
type App struct {
	ctx       context.Context
	container *appcore.Container
}

// NewApp creates the root App with its assembled container.
func NewApp(c *appcore.Container) *App {
	return &App{container: c}
}

// startup runs when the webview is ready. We keep the context for calling Wails
// runtime APIs (dialogs, events) later.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown runs when the app is closing; release the database connection.
func (a *App) shutdown(ctx context.Context) {
	if a.container != nil {
		_ = a.container.Close()
	}
}

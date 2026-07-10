package main

import (
	"embed"
	"os"

	appcore "autoshop/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Assemble the application (open DB, migrate, seed, wire services) before
	// the UI launches. If this fails the app cannot run, so exit clearly.
	container, err := appcore.Bootstrap()
	if err != nil {
		println("Fatal: failed to initialise application:", err.Error())
		os.Exit(1)
	}

	application := NewApp(container)

	// Bind the root App plus every module handler the container exposes.
	bind := append([]interface{}{application}, container.Handlers()...)

	err = wails.Run(&options.App{
		Title:  "Auto Spare Parts Manager",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        application.startup,
		OnShutdown:       application.shutdown,
		Bind:             bind,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	g "xabbo.b7c.io/goearth"
)

//go:embed all:frontend/dist
var assets embed.FS

var ext = g.NewExt(g.ExtInfo{
	Title:       "G-WallMover",
	Description: "Move wall items in Origins",
	Version:     "1.0",
	Author:      "QDave",
})

var app *App

func main() {
	app = NewApp(ext, assets)

	setupExtension()

	err := wails.Run(&options.App{
		Title:  "G-WallMover",
		Width:  355,
		Height: 738,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		StartHidden:       true,
		HideWindowOnClose: true,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func setupExtension() {
	ext.Initialized(func(e g.InitArgs) {
		log.Printf("initialized (connected=%t)", e.Connected)
	})

	ext.Activated(func() {
		log.Printf("activated")
		app.ShowWindow()
	})

	ext.Connected(func(e g.ConnectArgs) {
		log.Printf("connected (%s:%d)", e.Host, e.Port)
		log.Printf("client %s (%s)", e.Client.Identifier, e.Client.Version)
	})

	ext.Disconnected(func() {
		log.Printf("connection lost")
	})
}

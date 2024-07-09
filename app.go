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
<<<<<<< HEAD
	setupExt()
	err := wails.Run(&options.App{
		Title:  "G-WallMover",
		Width:  355,
		Height: 481,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 44, G: 62, B: 80, A: 1},
=======

	setupExtension()

	err := wails.Run(&options.App{
		Title:  "G-WallMover",
		Width:  355,
		Height: 738,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		StartHidden:       true,
		HideWindowOnClose: true,
<<<<<<< HEAD
		MinWidth:          355,
		MaxWidth:          355,
		MinHeight:         481,
		MaxHeight:         481,
=======
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
	})

	if err != nil {
		log.Fatal(err)
	}
}

<<<<<<< HEAD
func setupExt() {
=======
func setupExtension() {
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
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

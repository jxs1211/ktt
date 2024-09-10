package main

import (
	"context"
	"embed"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"

	"ktt/backend/cli"
	"ktt/backend/client"
	"ktt/backend/consts"
	"ktt/backend/services"
	"ktt/backend/utils/log"
	strutil "ktt/backend/utils/string"
	"ktt/backend/watch"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"
var gaMeasurementID, gaSecretKey string

func init() {
	log.Init(filepath.Join(strutil.RootPath(), "logs"))
}

func main() {
	clientSvc := client.NewClientService()
	watcherManager := watch.NewWatcherManager()
	terminalSvc := cli.NewTerminalService()
	// Create an instance of the app structure
	sysSvc := services.System()
	connSvc := services.Connection()
	browserSvc := services.Browser()
	cliSvc := services.Cli()
	monitorSvc := services.Monitor()
	pubsubSvc := services.Pubsub()
	prefSvc := services.Preferences()
	prefSvc.SetAppVersion(version)
	windowWidth, windowHeight, maximised := prefSvc.GetWindowSize()
	windowStartState := options.Normal
	if maximised {
		windowStartState = options.Maximised
	}

	// menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:                    "KT",
		Width:                    windowWidth,
		Height:                   windowHeight,
		MinWidth:                 consts.MIN_WINDOW_WIDTH,
		MinHeight:                consts.MIN_WINDOW_HEIGHT,
		WindowStartState:         windowStartState,
		Frameless:                runtime.GOOS != "darwin",
		Menu:                     appMenu,
		EnableDefaultContextMenu: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: options.NewRGBA(27, 38, 54, 0),
		StartHidden:      true,
		OnStartup: func(ctx context.Context) {
			clientSvc.Start(ctx)
			watcherManager.Start(ctx)
			terminalSvc.Start(ctx)
			sysSvc.Start(ctx, version)
			connSvc.Start(ctx)
			browserSvc.Start(ctx)
			cliSvc.Start(ctx)
			monitorSvc.Start(ctx)
			pubsubSvc.Start(ctx)

			services.GA().SetSecretKey(gaMeasurementID, gaSecretKey)
			services.GA().Startup(version)
		},
		OnDomReady: func(ctx context.Context) {
			x, y := prefSvc.GetWindowPosition(ctx)
			runtime2.WindowSetPosition(ctx, x, y)
			runtime2.WindowShow(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			x, y := runtime2.WindowGetPosition(ctx)
			prefSvc.SaveWindowPosition(x, y)
			return false
		},
		OnShutdown: func(ctx context.Context) {
			browserSvc.Stop()
			cliSvc.CloseAll()
			monitorSvc.StopAll()
			pubsubSvc.StopAll()
		},
		Bind: []interface{}{
			clientSvc,
			sysSvc,
			terminalSvc,
			connSvc,
			browserSvc,
			cliSvc,
			monitorSvc,
			pubsubSvc,
			prefSvc,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "KT " + version,
				Message: "A modern lightweight cross-platform Redis desktop client.\n\nCopyright © 2024",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
		},
		Linux: &linux.Options{
			ProgramName:         "KT",
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

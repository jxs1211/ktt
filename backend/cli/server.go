package cli

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sorenisanerd/gotty/backend/localcommand"
	"github.com/sorenisanerd/gotty/pkg/homedir"
	"github.com/sorenisanerd/gotty/server"
	"github.com/sorenisanerd/gotty/utils"
	cli "github.com/urfave/cli/v2"
)

type CliServer struct {
	app        *cli.App
	ctx        context.Context
	ctxCancel  context.CancelFunc
	gCtx       context.Context
	gCtxCancel context.CancelFunc
	errs       chan error
}

// func NewServer() (*CliServer, error) {

// 	appOptions := &server.Options{}

// 	if err := utils.ApplyDefaultValues(appOptions); err != nil {
// 		return nil, err
// 	}
// 	backendOptions := &localcommand.Options{}
// 	if err := utils.ApplyDefaultValues(backendOptions); err != nil {
// 		return nil, err
// 	}

// 	cliFlags, flagMappings, err := utils.GenerateFlags(appOptions, backendOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	args := c.Args()
// 	factory, err := localcommand.NewFactory(args.First(), args.Tail(), backendOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	hostname, _ := os.Hostname()
// 	appOptions.TitleVariables = map[string]interface{}{
// 		"command":  args.First(),
// 		"argv":     args.Tail(),
// 		"hostname": hostname,
// 	}

// 	srv, err := server.New(factory, appOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// }

func NewCliServer() (*CliServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	gCtx, gCtxCancel := context.WithCancel(context.Background())
	cliServer := &CliServer{
		ctx:        ctx,
		ctxCancel:  cancel,
		gCtx:       gCtx,
		gCtxCancel: gCtxCancel,
		errs:       make(chan error, 1),
	}
	app := cli.NewApp()
	app.Name = "gotty"
	app.Version = "unknown_version"
	app.Usage = "Share your terminal as a web application"
	app.HideHelpCommand = true
	appOptions := &server.Options{}

	if err := utils.ApplyDefaultValues(appOptions); err != nil {
		return nil, err
	}
	backendOptions := &localcommand.Options{}
	if err := utils.ApplyDefaultValues(backendOptions); err != nil {
		return nil, err
	}

	cliFlags, flagMappings, err := utils.GenerateFlags(appOptions, backendOptions)
	if err != nil {
		return nil, err
	}

	app.Flags = append(
		cliFlags,
		&cli.StringFlag{
			Name:    "config",
			Value:   "~/.gotty",
			Usage:   "Config file path",
			EnvVars: []string{"GOTTY_CONFIG"},
		},
	)

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			msg := "error: No command given."
			cli.ShowAppHelp(c)
			return errors.New(msg)
		}

		configFile := c.String("config")
		_, err := os.Stat(homedir.Expand(configFile))
		if configFile != "~/.gotty" || !os.IsNotExist(err) {
			if err := utils.ApplyConfigFile(configFile, appOptions, backendOptions); err != nil {
				return err
			}
		}

		utils.ApplyFlags(cliFlags, flagMappings, c, appOptions, backendOptions)

		if appOptions.Quiet {
			log.SetFlags(0)
			log.SetOutput(io.Discard)
		}

		if c.IsSet("credential") {
			appOptions.EnableBasicAuth = true
		}
		if c.IsSet("tls-ca-crt") {
			appOptions.EnableTLSClientAuth = true
		}

		err = appOptions.Validate()
		if err != nil {
			return err
		}

		args := c.Args()
		factory, err := localcommand.NewFactory(args.First(), args.Tail(), backendOptions)
		if err != nil {
			return err
		}

		hostname, _ := os.Hostname()
		appOptions.TitleVariables = map[string]interface{}{
			"command":  args.First(),
			"argv":     args.Tail(),
			"hostname": hostname,
		}

		srv, err := server.New(factory, appOptions)
		if err != nil {
			return err
		}

		log.Printf("GoTTY is starting with command: %s", strings.Join(args.Slice(), " "))

		go func() {
			cliServer.errs <- srv.Run(cliServer.ctx, server.WithGracefullContext(cliServer.gCtx))
		}()
		// err = waitSignals(cliServer.errs, cliServer.ctxCancel, cliServer.gCtxCancel)

		select {
		case err := <-cliServer.errs:
			return err
		case <-cliServer.gCtx.Done():
			log.Println("g shutdown")
		}

		if err != nil && err != context.Canceled {
			return err
		}

		return nil
	}
	return &CliServer{
		app: app,
	}, nil
}

func (s *CliServer) Start() error {
	return s.app.Run(os.Args)
}

func (s *CliServer) Close() error {
	s.gCtxCancel()
	return <-s.errs
}

func (s *CliServer) Restart() error {
	return nil
}

// func waitSignals(errs chan error, cancel context.CancelFunc, gracefullCancel context.CancelFunc) error {
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(
// 		sigChan,
// 		syscall.SIGINT,
// 		syscall.SIGTERM,
// 	)
// 	select {
// 	case err := <-errs:
// 		return err
// 	case s := <-sigChan:
// 		switch s {
// 		case syscall.SIGINT:
// 			gracefullCancel()
// 			fmt.Println("C-C to force close")
// 			select {
// 			case err := <-errs:
// 				return err
// 			case <-sigChan:
// 				fmt.Println("Force closing...")
// 				cancel()
// 				return <-errs
// 			}
// 		default:
// 			cancel()
// 			return <-errs
// 		}
// 	}
// }

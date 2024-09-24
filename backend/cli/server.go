package cli

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sorenisanerd/gotty/backend/localcommand"
	"github.com/sorenisanerd/gotty/pkg/homedir"
	"github.com/sorenisanerd/gotty/server"
	"github.com/sorenisanerd/gotty/utils"
	cli "github.com/urfave/cli/v2"

	logutil "ktt/backend/utils/log"
)

type arguments []string

func (a *arguments) Get(n int) string {
	if len(*a) > n {
		return (*a)[n]
	}
	return ""
}

func (a *arguments) First() string {
	return a.Get(0)
}

func (a *arguments) Tail() []string {
	if a.Len() >= 2 {
		tail := []string((*a)[1:])
		ret := make([]string, len(tail))
		copy(ret, tail)
		return ret
	}
	return []string{}
}

func (a *arguments) Len() int {
	return len(*a)
}

func (a *arguments) Present() bool {
	return a.Len() != 0
}

func (a *arguments) Slice() []string {
	ret := make([]string, len(*a))
	copy(ret, *a)
	return ret
}

type CliServer struct {
	app        *cli.App
	args       []string
	ctx        context.Context
	ctxCancel  context.CancelFunc
	gCtx       context.Context
	gCtxCancel context.CancelFunc
	errs       chan error
}

func NewCliServer(pCtx context.Context, addr, port string, cmds []string) (*CliServer, error) {
	ctx, cancel := context.WithCancel(pCtx)
	gCtx, gCtxCancel := context.WithCancel(pCtx)
	cliServer := &CliServer{
		ctx:        ctx,
		ctxCancel:  cancel,
		gCtx:       gCtx,
		gCtxCancel: gCtxCancel,
		errs:       make(chan error, 1),
		args:       cmds,
	}
	app := cli.NewApp()
	app.Name = "KT-Console"
	app.Version = "unknown_version"
	app.Usage = "Share terminal as a web application"
	app.HideHelpCommand = true
	appOptions := &server.Options{}

	if err := utils.ApplyDefaultValues(appOptions); err != nil {
		return nil, err
	}
	appOptions.Address = addr
	appOptions.Port = port
	appOptions.EnableReconnect = true
	appOptions.PermitWrite = true
	// os.Args = append(os.Args, cmds...)
	// log.Printf("options: %+v", appOptions)
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
		logutil.Info("NewCliServer", "c.NArg", c.NArg())
		if len(cmds) == 0 {
			// if c.NArg() == 0 {
			msg := "error: no command given"
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

		args := arguments(cmds)
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

		log.Printf("starting with command: %s", strings.Join(args.Slice(), " "))

		go func() {
			cliServer.errs <- srv.Run(cliServer.ctx, server.WithGracefullContext(cliServer.gCtx))
		}()
		// err = waitSignals(cliServer.errs, cliServer.ctxCancel, cliServer.gCtxCancel)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(
			sigChan,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		select {
		case err := <-cliServer.errs:
			return err
		case <-cliServer.gCtx.Done():
			logutil.Info("NewCliServer", "gCtxCancel shutdown")
		case sig := <-sigChan:
			logutil.Info("NewCliServer", "shutdown by sig", sig.String())
		}

		// if err != nil && err != context.Canceled {
		// 	return err
		// }
		return nil
	}
	cliServer.app = app
	return cliServer, nil
}

func (s *CliServer) Start() error {
	go s.app.Run(s.args)
	return nil
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

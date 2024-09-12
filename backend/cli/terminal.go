package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/creack/pty"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"

	"ktt/backend/types"
	"ktt/backend/utils/log"
)

type TerminalService struct {
	ctx       context.Context
	pty       *os.File
	cmd       *exec.Cmd
	tempInput string
}

func NewTerminalService() *TerminalService {
	return &TerminalService{}
}

func (s *TerminalService) Start(ctx context.Context) {
	s.ctx = ctx
}

func (s *TerminalService) StartTerminal() types.JSResp {
	err := s.startTerminal()
	if err != nil {
		log.Error("start terminal failed", "msg", err)
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) startTerminal() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe")
	default:
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		cmd = exec.Command(shell)
	}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %v", err)
	}

	s.pty = ptmx
	s.cmd = cmd
	log.Info("terminal", "cmd", s.cmd.String())
	go s.readOutput()

	// runtime2.EventsOn(s.ctx, "terminal:input", func(data ...interface{}){
	// 	if len(data) > 0 {
	// 		if str, ok := data[0].(string); ok {
	// 			_, err := s.pty.Write([]byte(str))
	// 			if err != nil {
	// 				log.Info("input:", "err", err)
	// 			}
	// 		}
	// 	}
	// 	})
	return nil
}

func (s *TerminalService) readOutput() {
	for {
		buf := make([]byte, 1024)
		n, err := s.pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			// runtime2.LogErrorf(s.ctx, "Error reading from pty: %v", err)
			log.Error("terminal", "Error reading from pty", err)
			continue
		}
		data := string(buf[:n])
		if data == s.tempInput {
			log.Log.Warn("data should be ignored: ", data)
			continue
		}

		log.Info("terminal", "data", data)
		runtime2.EventsEmit(s.ctx, "terminal:output", data)
	}
}

func (s *TerminalService) WriteInput(input string) error {
	s.tempInput = input
	_, err := s.pty.Write([]byte(input))
	return err
}

func (s *TerminalService) Resize(rows, cols uint16) types.JSResp {
	err := pty.Setsize(s.pty, &pty.Winsize{Rows: rows, Cols: cols})
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) CloseTerminal() error {
	if s.pty != nil {
		s.pty.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}

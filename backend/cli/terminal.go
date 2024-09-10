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
)

type TerminalService struct {
	ctx context.Context
	pty *os.File
	cmd *exec.Cmd
}

func NewTerminalService() *TerminalService {
	return &TerminalService{}
}

func (s *TerminalService) Start(ctx context.Context) {
	s.ctx = ctx
}

func (s *TerminalService) StartTerminal() error {
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

	go s.readOutput()

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
			runtime2.LogErrorf(s.ctx, "Error reading from pty: %v", err)
			continue
		}
		runtime2.EventsEmit(s.ctx, "terminal:output", string(buf[:n]))
	}
}

func (s *TerminalService) WriteInput(input string) error {
	_, err := s.pty.Write([]byte(input))
	return err
}

func (s *TerminalService) Resize(rows, cols uint16) error {
	return pty.Setsize(s.pty, &pty.Winsize{Rows: rows, Cols: cols})
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

package cli

import (
	"bufio"
	"context"
	"fmt"
	"ktt/backend/utils/log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
	"syscall"

	"github.com/creack/pty"
)

type terminal interface {
	Start() error
	Close() error
}

type Terminal struct {
	pty   *os.File
	cmd   *exec.Cmd
	mutex sync.Mutex
	done  chan struct{}
	ctx   context.Context
}

func NewTerminal(ctx context.Context) *Terminal {
	return &Terminal{
		done: make(chan struct{}),
		ctx:  ctx,
	}
}

func (t *Terminal) Start() error {
	var cmd *exec.Cmd
	shell := os.Getenv("SHELL")

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe")
	case "darwin", "linux":
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.Command(shell)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	t.pty = f
	return nil
}

func (t *Terminal) Write(p []byte) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.pty.Write(p)
}

func (t *Terminal) Read(p []byte) (int, error) {
	return t.pty.Read(p)
}

func (t *Terminal) Close() error {
	close(t.done)
	t.cmd.Process.Kill()
	t.pty.Close()
	return nil
}

func (t *Terminal) SetSize(rows, cols uint16) error {
	return pty.Setsize(t.pty, &pty.Winsize{Rows: rows, Cols: cols})
}

type TtydTerminal struct {
	mu          sync.Mutex
	ctx         context.Context
	cmd         *exec.Cmd
	readySignal chan struct{} // Add a channel for readiness signal

}

func NewTtydTerminal(ctx context.Context) *TtydTerminal {
	return &TtydTerminal{
		ctx:         ctx,
		readySignal: make(chan struct{}), // Initialize the channel
	}
}

func (t *TtydTerminal) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	port := "9999"
	switch runtime.GOOS {
	case "windows":
		t.cmd = exec.Command("powershell.exe")
	default:
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		t.cmd = exec.Command("ttyd", "--port", port, "-W", shell)
		log.Info("Starting ttyd command", "command", t.cmd.String())

		// Start the command
		stdout, err := t.cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe: %v", err)
		}
		stderr, err := t.cmd.StderrPipe() // Capture stderr
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe: %v", err)
		}
		err = t.cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start terminal: %v", err)
		}
		log.Info("Started ttyd process", "pid", t.cmd.Process.Pid)

		stdoutReader := bufio.NewReader(stdout)
		stderrReader := bufio.NewReader(stderr)        // Read stderr
		go t.waitForBootup(stdoutReader, stderrReader) // Pass stderr reader
	}

	// Wait for ttyd to be ready
	if err := t.WaitForReady(); err != nil {
		return fmt.Errorf("ttyd did not become ready: %v", err)
	}

	url := fmt.Sprintf("http://localhost:%v", port)
	log.Info("terminal", "url", url)
	return nil
}

func (t *TtydTerminal) WaitForReady() error {
	<-t.readySignal
	return nil
}

func (t *TtydTerminal) Close() error {
	log.Info("Close", "t.Cmd.Process", t.cmd.Process)
	if t.cmd != nil {
		if err := t.cmd.Process.Signal(syscall.SIGHUP); err != nil {
			return fmt.Errorf("failed to kill ttyd process: %v", err)
		}
	}
	log.Info("Close done")
	return nil
}

func (t *TtydTerminal) waitForBootup(stdout *bufio.Reader, stderr *bufio.Reader) error {
	portRegex := regexp.MustCompile(`Listening on port: (\d+)`)

	// Read stdout
	go func() {
		for {
			line, err := stdout.ReadString('\n')
			if err != nil {
				break
			}
			log.Info("ttyd stdout", "line", line)
			if portRegex.MatchString(line) {
				matches := portRegex.FindStringSubmatch(line)
				if len(matches) > 1 {
					port := matches[1]
					log.Info("ttyd started on port", "port", port)
					close(t.readySignal) // Signal that ttyd is ready
				}
			}
		}
	}()

	// Read stderr
	go func() {
		for {
			line, err := stderr.ReadString('\n')
			if err != nil {
				break
			}
			log.Error("ttyd stderr", "line", line) // Log stderr output
		}
	}()

	return nil
}

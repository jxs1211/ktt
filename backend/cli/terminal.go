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
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
)

type terminal interface {
	Start() error
	Close() error
	isExited() bool
	Restart() error
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

func (t *Terminal) isExited() bool {
	return true
}

func (t *Terminal) Restart() error {
	return nil
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
	readySignal chan string // Add a channel for readiness signal

}

func NewTtydTerminal(ctx context.Context) *TtydTerminal {
	return &TtydTerminal{
		ctx:         ctx,
		readySignal: make(chan string), // Initialize the channel
	}
}

func (t *TtydTerminal) isExited() bool {
	return t.cmd.ProcessState.Exited()
}

func (t *TtydTerminal) Restart() error {
	return nil
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
	p := t.WaitForReady()
	url := fmt.Sprintf("http://localhost:%v", p)
	log.Info("terminal", "url", url)
	return nil
}

func (t *TtydTerminal) WaitForReady() string {
	return <-t.readySignal
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
					log.Info("ttyd started", "port", port)
					t.readySignal <- port // Signal that ttyd is ready
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

type GottyTerminal struct {
	mu          sync.Mutex
	ctx         context.Context
	cmd         *exec.Cmd
	readySignal chan string
}

func NewGottyTerminal(ctx context.Context) *GottyTerminal {
	return &GottyTerminal{
		ctx:         ctx,
		readySignal: make(chan string),
	}
}
func (t *GottyTerminal) Restart() error {
	err := t.Close()
	if err != nil {
		return err
	}
	return t.Start()
}

func (t *GottyTerminal) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var cmd *exec.Cmd
	port := "8887" // Default port for GoTTY

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("gotty", "powershell.exe")
	default:
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		cmd = exec.Command("gotty", "--port", port, "-w", shell)
	}
	// Start the command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start GoTTY: %v", err)
	}

	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)
	go t.waitForBootup(stdoutReader, stderrReader)

	// Wait for GoTTY to be ready
	p := t.WaitForReady()
	url := fmt.Sprintf("http://localhost:%v", p)
	log.Info("terminal", "url", url)
	t.cmd = cmd
	runtime2.EventsEmit(t.ctx, "terminal:url", url)
	return nil
}

func (t *GottyTerminal) WaitForReady() string {
	return <-t.readySignal
}

func (t *GottyTerminal) isExited() bool {
	return t.cmd.ProcessState.Exited()
}

func (t *GottyTerminal) Close() error {
	if t.cmd != nil {
		if err := t.cmd.Process.Signal(syscall.SIGHUP); err != nil {
			return fmt.Errorf("failed to kill GoTTY process: %v", err)
		}
	}
	log.Info("Close done")
	return nil
}

func match(pattern, content string) ([]string, bool) {
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(content) {
		return []string{}, false
	}
	return regex.FindStringSubmatch(content), true
}

func (t *GottyTerminal) waitForBootup(stdout *bufio.Reader, stderr *bufio.Reader) error {
	pattern := `HTTP server is listening at: http://[^\s]+:(\d+)/`
	// portRegex := regexp.MustCompile(`HTTP server is listening at: http://[^\s]+:(\d+)/`)

	// Read stdout
	go func() {
		for {
			line, err := stdout.ReadString('\n')
			if err != nil {
				break
			}
			// log.Info("GoTTY stdout", "line", line)
			matches, ok := match(pattern, line)
			if !ok {
				continue
			}
			if len(matches) > 1 {
				port := matches[1]
				log.Info("GoTTY started", "port", port)
				t.readySignal <- port // Signal that GoTTY is ready
				break
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
			// log.Error("GoTTY stderr", "line", line) // Log stderr output
			matches, ok := match(pattern, line)
			if !ok {
				continue
			}
			if len(matches) > 1 {
				port := matches[1]
				log.Info("GoTTY started", "port", port)
				t.readySignal <- port // Signal that GoTTY is ready
				break
			}
		}
	}()

	return nil
}

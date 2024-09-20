package cli

import (
	"context"
	"fmt"
	"sync"

	"ktt/backend/types"
	"ktt/backend/utils/log"
)

type TerminalService struct {
	ctx      context.Context
	terminal terminal
	mutex    sync.Mutex
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.terminal != nil {
		return fmt.Errorf("terminal is already running")
	}

	// terminal := NewTerminal(s.ctx)
	terminal := NewTtydTerminal(s.ctx)
	err := terminal.Start()
	if err != nil {
		return err
	}
	s.terminal = terminal
	// go s.readLoop()
	return nil
}

func (s *TerminalService) CloseTerminal() error {
	log.Info("CloseTerminal")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.terminal == nil {
		return fmt.Errorf("no terminal is running")
	}

	err := s.terminal.Close()
	if err != nil {
		return err
	}
	s.terminal = nil
	return err
}

// func (s *TerminalService) WriteToTerminal(input string) error {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()

// 	if s.terminal == nil {
// 		return fmt.Errorf("no terminal is running")
// 	}

// 	_, err := s.terminal.Write([]byte(input))
// 	return err
// }

// func (s *TerminalService) SetTerminalSize(rows, cols uint16) error {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()

// 	if s.terminal == nil {
// 		return fmt.Errorf("no terminal is running")
// 	}

// 	return s.terminal.SetSize(rows, cols)
// }

// func (s *TerminalService) readLoop() {
// 	buffer := make([]byte, 1024)
// 	for {
// 		n, err := s.terminal.Read(buffer)
// 		if err != nil {
// 			if err != io.EOF {
// 				fmt.Println("Error reading from terminal:", err)
// 			}
// 			break
// 		}
// 		data := buffer[:n]
// 		output := base64.StdEncoding.EncodeToString(buffer[:n])
// 		log.Info("terminal", "output", output)
// 		log.Info("terminal", "data", string(data))
// 		// Emit the output to the frontend
// 		runtime2.EventsEmit(s.ctx, "terminal:output", string(data))
// 	}
// }

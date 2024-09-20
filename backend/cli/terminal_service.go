package cli

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"ktt/backend/types"
	"ktt/backend/utils/log"
)

var ErrTerminalAlreadyRunning = errors.New("terminal is already running")
var ErrRestartTerminalFailed = errors.New("restart terminal failed")
var ErrCloseTerminalFailed = errors.New("close terminal failed")

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

func (s *TerminalService) restartTerminal() error {
	err := s.terminal.Restart()
	if err != nil {
		return err
	}
	return nil
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

	// terminal := NewTerminal(s.ctx)
	// terminal := NewTtydTerminal(s.ctx)
	terminal := NewGottyTerminal(s.ctx)
	err := terminal.Restart()
	// err := terminal.Start()
	if err != nil {
		return err
	}
	s.terminal = terminal
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
	return nil
}

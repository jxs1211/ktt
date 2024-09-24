package cli

import (
	"context"
	"errors"
	"fmt"
	"sync"

	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"

	"ktt/backend/types"
	"ktt/backend/utils/log"
)

var (
	ErrTerminalAlreadyRunning = errors.New("terminal is already running")
	ErrRestartTerminalFailed  = errors.New("restart terminal failed")
	ErrCloseTerminalFailed    = errors.New("close terminal failed")
	ErrTerminalNotExist       = errors.New("terminal not exist")
)

type TerminalService struct {
	ctx         context.Context
	terminalMap map[string]terminal
	mutex       sync.Mutex
}

func NewTerminalService() *TerminalService {
	return &TerminalService{
		terminalMap: make(map[string]terminal),
	}
}

func (s *TerminalService) Start(ctx context.Context) {
	s.ctx = ctx
}

// func (s *TerminalService) restartTerminal() error {
// 	err := s.terminal.Restart()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *TerminalService) StartTerminal(address, port string, cmds []string) types.JSResp {
	err := s.startTerminal(address, port, cmds)
	if err != nil {
		log.Error("start terminal failed", "msg", err)
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) terminalMapKey(address, port string) string {
	return fmt.Sprintf("%s:%s", address, port)
}

func (s *TerminalService) startTerminal(address, port string, cmds []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := s.terminalMapKey(address, port)
	terminal, ok := s.terminalMap[key]
	if !ok {
		srv, err := NewCliServer(s.ctx, address, port, cmds)
		if err != nil {
			return err
		}
		terminal = srv
	}

	err := terminal.Start()
	if err != nil {
		return err
	}
	s.terminalMap[key] = terminal
	url := fmt.Sprintf("http://%s:%s", address, port)
	runtime2.EventsEmit(s.ctx, "terminal:url", url)
	log.Info("startTerminal", "emit url", url)
	return nil
}

func (s *TerminalService) CloseTerminal(address, port string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := s.terminalMapKey(address, port)
	ter, ok := s.terminalMap[key]
	if !ok {
		return ErrTerminalNotExist
	}

	err := ter.Close()
	if err != nil {
		return err
	}
	delete(s.terminalMap, key)
	log.Info("CloseTerminal", "result", "done", "deleted item", address+port)
	return nil
}

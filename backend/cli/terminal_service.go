package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"

	"ktt/backend/db/store/session"
	"ktt/backend/kubeconfig"
	"ktt/backend/types"
	"ktt/backend/utils/log"
)

type DBEvent string

const (
	DBCreation DBEvent = "Creation"
	DBDeletion DBEvent = "Deletion"
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
	q           session.Queries
}

func NewTerminalService(db *sql.DB) *TerminalService {
	return &TerminalService{
		terminalMap: make(map[string]terminal),
		q:           *session.New(db),
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

func (s *TerminalService) StartTerminal2(address, port, cmds string) types.JSResp {
	err := s.startTerminal2(address, port, cmds)
	if err != nil {
		log.Error("start terminal failed", "msg", err)
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) startTerminal2(address, port, cmds string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	runtime2.EventsEmit(s.ctx, "terminal2:url", "")
	key := s.terminalMapKey2(address, port, cmds)
	_, ok := s.terminalMap[key]
	if !ok {
		srv, err := s.startCliServer(address, port, cmds)
		if err != nil {
			return err
		}

		s.terminalMap[key] = srv
		log.Info("startTerminal", "new terminal", key, "reason", "start cli server due to server not found in db")
	}
	url := fmt.Sprintf("http://%s:%s", address, port)
	runtime2.EventsEmit(s.ctx, "terminal2:url", url)
	log.Info("startTerminal", "emit url", url, "save item", key)
	return nil
}

func (s *TerminalService) StartTerminal(clusterName, address, port, cmds string) types.JSResp {
	if len(clusterName) == 0 && len(address) == 0 && len(port) == 0 && len(cmds) == 0 {
		return types.FailedResp(fmt.Sprintf("all of them must be not empty: %s,%s,%s,%s", clusterName, address, port, cmds))
	}
	err := s.startTerminal(clusterName, address, port, cmds)
	if err != nil {
		log.Error("start terminal failed", "msg", err)
		return types.FailedResp(err.Error())
	}
	// switch ctx on user's kubeconfig
	contxt, err := kubeconfig.SwitchContext(clusterName)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true, Data: contxt}
}

func (s *TerminalService) terminalMapKey(clusterName, address, port, cmds string) string {
	return fmt.Sprintf("%s:%s:%s:%s", clusterName, address, port, cmds)
}

func (s *TerminalService) terminalMapKey2(address, port, cmds string) string {
	return fmt.Sprintf("%s:%s:%s", address, port, cmds)
}

func (s *TerminalService) startTerminal(clusterName, address, port, cmds string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	runtime2.EventsEmit(s.ctx, "terminal:url", "")
	key := s.terminalMapKey(clusterName, address, port, cmds)
	_, ok := s.terminalMap[key]
	if !ok {
		srv, err := s.startCliServer(address, port, cmds)
		if err != nil {
			return err
		}

		s.terminalMap[key] = srv
		log.Info("startTerminal", "new terminal", key, "reason", "start cli server due to server not found in db")
	}
	// p, err := strconv.Atoi(port)
	// if err != nil {
	// 	return err
	// }
	// isOpened, err := tool.IsPortOpen(address, p)
	// if err != nil {
	// 	return err
	// }
	// if ktt restarted, the server exists in db but not exists in runtime
	// if !isOpened {
	// 	srv, err := s.startCliServer(address, port, cmds)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	s.terminalMap[key] = srv
	// 	log.Info("startTerminal", "new terminal", key, "reason", "server not found in runtime")
	// }
	url := fmt.Sprintf("http://%s:%s", address, port)
	runtime2.EventsEmit(s.ctx, "terminal:url", url)
	log.Info("startTerminal", "emit url", url, "save item", key)
	return nil
}

func (s *TerminalService) startCliServer(address, port, cmds string) (*CliServer, error) {
	srv, err := NewCliServer(s.ctx, address, port, strings.Split(cmds, " "))
	if err != nil {
		return nil, err
	}
	err = srv.Start()
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *TerminalService) CloseTerminal2(address, port, cmds string) types.JSResp {
	err := s.closeTerminal2(address, port, cmds)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) closeTerminal2(address, port, cmds string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := s.terminalMapKey2(address, port, cmds)
	if ter, ok := s.terminalMap[key]; ok {
		if err := ter.Close(); err != nil {
			return err
		}
		delete(s.terminalMap, key)
	}
	log.Info("CloseTerminal", "result", "done", "deleted item", key)
	return nil
}

func (s *TerminalService) CloseTerminal(id int64, clusterName, address, port, cmds string) types.JSResp {
	err := s.closeTerminal(id, clusterName, address, port, cmds)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) closeTerminal(id int64, clusterName, address, port, cmds string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := s.terminalMapKey(clusterName, address, port, cmds)
	if ter, ok := s.terminalMap[key]; ok {
		if err := ter.Close(); err != nil {
			return err
		}
		delete(s.terminalMap, key)
	}
	err := s.q.DeleteSession(s.ctx, id)
	if err != nil {
		return err
	}

	log.Info("CloseTerminal", "result", "done", "deleted item", key)
	return nil
}

func (s *TerminalService) StopAll() types.JSResp {
	err := s.closeAllTerminals()
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) closeAllTerminals() error {
	if len(s.terminalMap) == 0 {
		return nil
	}
	errs := make([]error, 0, len(s.terminalMap))
	for key, ter := range s.terminalMap {
		err := ter.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s", key, err))
		}
	}
	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	s.terminalMap = nil
	return nil
}

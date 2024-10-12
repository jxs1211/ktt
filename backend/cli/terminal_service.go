package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"ktt/backend/db/store/session"
	"ktt/backend/kubeconfig"
	"ktt/backend/types"
	"ktt/backend/utils/log"

	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
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

type TerminalServiceOptions struct {
	EventEmit bool
}

type TerminalService struct {
	ctx         context.Context
	terminalMap map[string]terminal
	mutex       sync.Mutex
	q           session.Queries
	EventEmit   bool
}

func NewTerminalService(db *sql.DB, eventEmit bool) *TerminalService {
	return &TerminalService{
		terminalMap: make(map[string]terminal),
		q:           *session.New(db),
		EventEmit:   eventEmit,
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

func (s *TerminalService) cacheStartCliServer(clusterName, address, port, cmds string) error {
	key := s.terminalMapKey(clusterName, address, port, cmds)
	_, ok := s.terminalMap[key]
	if !ok {
		srv, err := s.startCliServer(address, port, cmds)
		log.Info("cacheStartCliServer", "created srv", srv)
		if err != nil {
			return err
		}
		s.terminalMap[key] = srv
		log.Info("cacheStartCliServer", "save item", key, "termMap", s.terminalMap, "reason", "start cli server due to server not found in db")
	}
	return nil
}

func (s *TerminalService) startTerminal(clusterName, address, port, cmds string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.EventEmit {
		runtime2.EventsEmit(s.ctx, "terminal:url", "")
	}
	err := s.cacheStartCliServer(clusterName, address, port, cmds)
	if err != nil {
		return err
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
	if s.EventEmit {
		runtime2.EventsEmit(s.ctx, "terminal:url", url)
	}
	log.Info("startTerminal", "emit url", url)
	return nil
}

func (s *TerminalService) startCliServer(address, port, cmds string) (*CliServer, error) {
	srv, err := NewCliServer(s.ctx, address, port, strings.Split(cmds, " "))
	if err != nil {
		return nil, err
	}
	log.Info("startCliServer", "srv", srv)
	err = srv.Start()
	if err != nil {
		return nil, err
	}
	return srv, nil
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
	if ter := s.terminalMap[key]; ter != nil {
		log.Info("closeTerminal", "term", fmt.Sprintf("%+v", ter))
		if err := ter.Close(); err != nil {
			log.Info("closeTerminal", "terminal in cache, but close error", key)
			return err
		}
	}
	delete(s.terminalMap, key)
	err := s.q.DeleteSession(s.ctx, id)
	if err != nil {
		return err
	}
	log.Info("CloseTerminal", "result", "done", "deleted item", key)
	return nil
}

func (s *TerminalService) CloseAllTerminals() types.JSResp {
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
	s.terminalMap = map[string]terminal{}
	return nil
}

func (s *TerminalService) EditTerminal(id int64, clusterName, address, port, cmds string) types.JSResp {
	err := s.editTerminal(id, clusterName, address, port, cmds)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *TerminalService) editTerminal(id int64, clusterName, address, port, cmds string) error {
	oldOne, err := s.q.GetSession(s.ctx, id)
	if err != nil {
		return err
	}
	newOne, err := s.q.UpdateSession(s.ctx, session.UpdateSessionParams{
		ClusterName: clusterName, Address: address, Port: port, Cmds: cmds, ID: id,
	})
	if err != nil {
		return err
	}
	// delete cache if srv's port or cmds updated
	if oldOne.Cmds != newOne.Cmds || oldOne.Port != newOne.Port {
		key := s.terminalMapKey(clusterName, address, oldOne.Port, oldOne.Cmds)
		ter, ok := s.terminalMap[key].(*CliServer)
		log.Info("editTerminal", "to delete", key, "termMap", s.terminalMap, "parsed obj", ter)
		if !ok {
			return fmt.Errorf("get terminal from cache failed: %+v", ter)
		}
		if ter != nil {
			if err := ter.Close(); err != nil {
				log.Info("closeTerminal", "terminal in cache, but close error", key)
				return err
			}
			log.Info("editTerminal", "delete old key", key)
			delete(s.terminalMap, key)
		}
		// term, ok := s.terminalMap[key]
		// if ok {
		// 	delete(s.terminalMap, key)
		// 	log.Info("editTerminal", "delete key", key)
		// 	if err := term.Close(); err != nil {
		// 		log.Info("editTerminal", "err", err)
		// 		return err
		// 	}
		// }
		// err = s.cacheStartCliServer(clusterName, address, newOne.Port, newOne.Cmds)
		// if err != nil {
		// 	return err
		// }
		// log.Info("EditTerminal", "start server ok", fmt.Sprintf("%s:%s:%s:%s", clusterName, address, newOne.Port, newOne.Cmds))
	}
	log.Info("editTerminal", "msg", "ok")
	return nil
}

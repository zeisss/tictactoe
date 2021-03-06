package lock

import (
	"log"
	"os"
	"sync"
	"time"
)

func New(verbose bool) *localLockManager {
	return &localLockManager{
		logger:    log.New(os.Stderr, "[lock] ", log.LstdFlags),
		verbose:   verbose,
		Global:    &sync.Mutex{},
		GameLocks: make(map[string]*sync.Mutex),
	}
}

type localLockManager struct {
	logger  *log.Logger
	verbose bool

	Global    *sync.Mutex
	GameLocks map[string]*sync.Mutex
}

func (lck *localLockManager) LockGame(gameId string, timeout time.Duration) {
	if lck.verbose {
		lck.logger.Printf("Lock:   %s", gameId)
	}

	lck.Global.Lock()
	defer lck.Global.Unlock()

	mutex, ok := lck.GameLocks[gameId]
	if !ok {
		mutex = &sync.Mutex{}
		lck.GameLocks[gameId] = mutex
	}

	mutex.Lock()
}

func (lck *localLockManager) UnlockGame(gameId string) {
	if lck.verbose {
		lck.logger.Printf("Unlock: %s", gameId)
	}

	lck.Global.Lock()
	defer lck.Global.Unlock()

	mutex, ok := lck.GameLocks[gameId]
	if ok {
		mutex.Unlock()
	}
}

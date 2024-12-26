package services

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/weekend-dev-labs/lancer/db/repo"
)

type TaskManager struct {
	mu    sync.Mutex
	tasks map[string]context.CancelFunc
	repo  *repo.Repo
}

func NewTaskManager(repo *repo.Repo) *TaskManager {
	return &TaskManager{
		tasks: make(map[string]context.CancelFunc),
		repo:  repo,
	}
}

func BaseTask(folder_path string, ctx context.Context) error {
	return os.RemoveAll(folder_path)
}

func (tm *TaskManager) AddTask(id string, duration time.Duration, task func(ctx context.Context)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancel, exists := tm.tasks[id]; exists {
		cancel()
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	tm.tasks[id] = cancel

	go func() {

		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				fmt.Println("Task cancelled:", id)
			} else if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Task timed out:", id)
				task(ctx)
			}
		}

		defer func() {
			tm.mu.Lock()
			delete(tm.tasks, id)
			tm.mu.Unlock()
		}()

	}()
}

func (tm *TaskManager) CancelTask(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancel, exists := tm.tasks[id]; exists {
		cancel()
		tm.repo.DeleteSession(id)
		delete(tm.tasks, id)
	}
}

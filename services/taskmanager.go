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
	tasks map[string]*Task
	repo  *repo.Repo
}

type Task struct {
	ctx       context.Context
	cancel    context.CancelFunc
	extension chan time.Duration
	execute   chan bool
}

func NewTaskManager(repo *repo.Repo) *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
		repo:  repo,
	}
}

func BaseTask(folder_path string, ctx context.Context) error {
	return os.RemoveAll(folder_path)
}

func (tm *TaskManager) AddTask(id string, duration time.Duration, task func(ctx context.Context)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if existingTask, exists := tm.tasks[id]; exists {
		existingTask.cancel()
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	t := &Task{
		ctx:       ctx,
		cancel:    cancel,
		extension: make(chan time.Duration, 1),
		execute:   make(chan bool, 1),
	}
	tm.tasks[id] = t

	go func() {
		defer func() {
			tm.mu.Lock()
			delete(tm.tasks, id)
			tm.mu.Unlock()
		}()

		for {
			select {
			case ext := <-t.extension:
				cancel()
				ctx, cancel := context.WithTimeout(context.Background(), ext)
				t.ctx = ctx
				t.cancel = cancel

			case <-t.execute:
				task(t.ctx)
				return

			case <-t.ctx.Done():
				if t.ctx.Err() == context.Canceled {
				} else if t.ctx.Err() == context.DeadlineExceeded {
					task(t.ctx)
				}
				return
			}

		}
	}()
}

func (tm *TaskManager) ExtendDuration(id string, extension time.Duration) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[id]; exists {
		select {
		case task.extension <- extension:
			fmt.Println("Duartion extended for the task: ", id)
		default:
			fmt.Println("failed to extend duration, Channel Busy : ", id)
		}
	} else {
		fmt.Println("Task Not found : ", id)
	}
}

func (tm *TaskManager) Execute(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[id]; exists {
		select {
		case task.execute <- true:
			return
		default:
			return
		}
	}
}

func (tm *TaskManager) CancelWithBaseTask(path string, id string) error {
	// tm.CancelTask(id)

	return BaseTask(path, context.Background())
}

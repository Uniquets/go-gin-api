package cron_server

import (
	"sync"

	"github.com/xinliangnote/go-gin-api/internal/api/repository/db_repo/cron_task_repo"
	"github.com/xinliangnote/go-gin-api/internal/pkg/cache"
	"github.com/xinliangnote/go-gin-api/internal/pkg/db"
	"github.com/xinliangnote/go-gin-api/pkg/errors"

	"github.com/jakecoffman/cron"
	"go.uber.org/zap"
)

var _ Server = (*server)(nil)

type taskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *taskCount) i() {}

func (tc *taskCount) Add() {
	tc.wg.Add(1)
}

func (tc *taskCount) Done() {
	tc.wg.Done()
}

func (tc *taskCount) Exit() {
	tc.wg.Done()
	<-tc.exit
}

func (tc *taskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

type server struct {
	logger    *zap.Logger
	db        db.Repo
	cache     cache.Repo
	cron      *cron.Cron
	taskCount *taskCount
}

type Server interface {
	i()
	Start()
	Stop()
	AddTask(task *cron_task_repo.CronTask)
	AddJob(task *cron_task_repo.CronTask) cron.FuncJob
	RemoveTask(taskId int)
}

func New(logger *zap.Logger, db db.Repo, cache cache.Repo) (Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	if db == nil {
		return nil, errors.New("db required")
	}

	if cache == nil {
		return nil, errors.New("cache required")
	}

	return &server{
		logger: logger,
		db:     db,
		cache:  cache,
		cron:   cron.New(),
		taskCount: &taskCount{
			wg:   sync.WaitGroup{},
			exit: make(chan struct{}),
		},
	}, nil
}

func (s *server) i() {}

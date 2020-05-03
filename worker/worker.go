package worker

import (
	"github.com/hibiken/asynq"
	"github.com/my1562/crawler/config"
	"github.com/my1562/crawler/tasks"
	"github.com/my1562/queue"
)

type Worker struct {
	server *asynq.Server
	tasks  *tasks.Tasks
}

func NewWorker(config *config.Config, tasks *tasks.Tasks) *Worker {
	return &Worker{
		server: asynq.NewServer(
			asynq.RedisClientOpt{
				Addr: config.Redis,
			},
			asynq.Config{
				Concurrency: 1,
			},
		),
		tasks: tasks,
	}
}

func (w *Worker) Listen() error {
	mux := asynq.NewServeMux()
	mux.Handle(
		queue.TaskTypePriorityCheck,
		queue.NewPriorityCheckHandler(w.tasks),
	)
	if err := w.server.Run(mux); err != nil {
		return err
	}
	return nil
}

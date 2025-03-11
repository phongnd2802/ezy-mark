package worker

import (
	"context"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/phongnd2802/ezy-mark/internal/pkg/email"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskRemoveOldAvatar(ctx context.Context, task *asynq.Task) error
}

type redisTaskProcessor struct {
	server *asynq.Server
	sender email.EmailSender
}


// Start implements TaskProcessor.
func (processor *redisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(TaskRemoveOldAvatar, processor.ProcessTaskRemoveOldAvatar)
	return processor.server.Start(mux)
}

func (processor *redisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, sender email.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			slog.Info("process task failed", "type", task.Type(), "payload", task.Payload(), "error", err.Error())
		}),
		Logger: NewLogger(),
	})

	return &redisTaskProcessor{
		server: server,
		sender: sender,
	}

}

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

// DistributeTaskSendVerifyEmail implements TaskDistributor.
func (distributor *redisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	slog.Info("Enqueued Task", "type", task.Type(), "payload", task.Payload(), "queue", info.Queue, "max_retry", info.MaxRetry)
	return nil
}

// ProcessTaskSendVerifyEmail implements TaskProcessor.
func (processor *redisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Send Email
	// subject := "Welcome to Simple Bank Service"
	// content := fmt.Sprintf(`Hello,<br/> 
	// Thank you for registering with us!<br/>
	// OTP is <strong> %s </strong>`, payload.Otp)

	// to := []string{payload.Email}
	// err := processor.sender.SendEmail(subject, content, to, nil, nil, nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to send otp email: %w", err)
	// }

	slog.Info("Processed Task", "type", task.Type(), "payload", task.Payload())
	return nil
}

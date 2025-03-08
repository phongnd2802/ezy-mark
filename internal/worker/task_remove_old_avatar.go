package worker

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
	"github.com/minio/minio-go/v7"
	"github.com/phongnd2802/ezy-mark/internal/consts"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/rs/zerolog/log"
)

const TaskRemoveOldAvatar = "task:remove_old_avatar"

type PayloadRemoveOldAvatar struct {
	ObjectName string `json:"object_name"`
}

// DistributeTaskRemoveOldAvatar implements TaskDistributor.
func (distributor *redisTaskDistributor) DistributeTaskRemoveOldAvatar(ctx context.Context, payload *PayloadRemoveOldAvatar, opts ...asynq.Option) error {
	jsonPayload, err := sonic.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskRemoveOldAvatar, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).
				Str("payload", string(task.Payload())).
				Str("queue", info.Queue).
				Int("max_retry", info.MaxRetry).
				Msg("Enqueued Task")
	return nil
}

// ProcessTaskRemoveOldAvatar implements TaskProcessor.
func (processor *redisTaskProcessor) ProcessTaskRemoveOldAvatar(ctx context.Context, task *asynq.Task) error {
	var payload PayloadRemoveOldAvatar
	if err := sonic.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	err := global.Minio.RemoveObject(ctx, consts.BucketUserAvatar, payload.ObjectName, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	})
	if err != nil {
		log.Error().Err(err).Str("object_name", payload.ObjectName).Msg("failed to remove old avatar")
		return err
	}

	log.Info().Str("object_name", payload.ObjectName).Msg("removed old avatar")
	return nil
}

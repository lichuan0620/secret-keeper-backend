package service

import (
	"context"

	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/standard"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
)

func Sync(ctx context.Context, req *models.SyncRequest) (models.SyncResponse, standard.Error) {
	GetQueue(ctx).Sync(models.QueueItem{
		Id:    req.Id,
		Score: req.Score,
	})
	return models.SyncResponse{}, nil
}

func Dequeue(ctx context.Context) (*models.DequeueResponse, standard.Error) {
	id, err := GetQueue(ctx).Dequeue()
	if err != nil {
		log.FromContext(ctx).Error(err, "dequeue Box")
		return nil, standard.InternalServiceError()
	}
	return &models.DequeueResponse{Id: id}, nil
}

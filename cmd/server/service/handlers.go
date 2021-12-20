package service

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/lichuan0620/secret-keeper-backend/pkg/mongo"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/standard"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
)

var location *time.Location

func init() {
	location, _ = time.LoadLocation("Asia/Shanghai")
}

func CreateBox(ctx context.Context, req *models.CreateBoxRequest) (*models.CreateBoxResponse, standard.Error) {
	logger := log.FromContext(ctx)
	db := mongo.DB()
	defer db.Session.Close()
	now := time.Now().In(location)
	box := models.Box{
		Id:         uuid.New().String(),
		CreatedAt:  &now,
		Body:       req.Body,
		Viewed:     []time.Time{now},
		LastViewed: &now,
	}
	if err := db.C(mongo.CollectionBox).Insert(&box); err != nil {
		logger.Error(err, "unexpected database error")
		return nil, standard.InternalServiceError()
	}
	qc := GetQueueClient(ctx)
	go func() {
		if err := qc.Sync(context.TODO(), &models.SyncRequest{
			Id:    box.Id,
			Score: now.UnixNano(),
		}); err != nil {
			logger.Error(err, "sync created Box")
		}
	}()
	return (*models.CreateBoxResponse)(&box), nil
}

func AddBoxEmoji(ctx context.Context, req *models.AddBoxEmojiRequest) (*models.AddBoxEmojiResponse, standard.Error) {
	db := mongo.DB()
	defer db.Session.Close()
	if len(req.EmojiFeedbacks) > 0 {
		incOpt := make(bson.M, len(req.EmojiFeedbacks))
		for k, v := range req.EmojiFeedbacks {
			incOpt["EmojiFeedbacks."+k] = v
		}
		if err := db.C(mongo.CollectionBox).UpdateId(req.Id, bson.M{"$inc": incOpt}); err != nil {
			if err == mgo.ErrNotFound {
				return nil, standard.ResourceNotFound(req.Id)
			}
			log.FromContext(ctx).Error(err, "unexpected database error")
			return nil, standard.InternalServiceError()
		}
	}
	var box models.Box
	if err := db.C(mongo.CollectionBox).FindId(req.Id).One(&box); err != nil {
		if err == mgo.ErrNotFound {
			return nil, standard.ResourceNotFound(req.Id)
		}
		log.FromContext(ctx).Error(err, "unexpected database error")
		return nil, standard.InternalServiceError()
	}
	return &models.AddBoxEmojiResponse{
		Id:             req.Id,
		EmojiFeedbacks: box.EmojiFeedbacks,
	}, nil
}

func ViewBox(ctx context.Context) (*models.ViewBoxResponse, standard.Error) {
	resp, err := GetQueueClient(ctx).Dequeue(ctx)
	if err != nil {
		log.FromContext(ctx).Error(err, "view item from queue")
		return nil, standard.InternalServiceError()
	}
	db := mongo.DB()
	defer db.Session.Close()
	var box models.Box
	if err = db.C(mongo.CollectionBox).FindId(resp.Id).One(&box); err != nil {
		log.FromContext(ctx).Error(err, "unexpected database error")
		return nil, standard.InternalServiceError()
	}
	return (*models.ViewBoxResponse)(&box), nil
}

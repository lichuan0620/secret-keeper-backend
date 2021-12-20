package service

import (
	"net/http"

	"github.com/lichuan0620/secret-keeper-backend/internal/queue"
	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/middlewares"
	servicemodel "github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
)

func Build(q queue.Interface) (http.Handler, error) {
	logger := log.New().WithName("handlers")
	return (&service.Builder{
		GlobalMiddlewares: []servicemodel.Middleware{middlewares.WithLogger(logger)},
	}).AddActionGroup(servicemodel.ActionGroup{
		Mutator: func(action *servicemodel.Action) {
			action.Version = models.Version
		},
		Middlewares: []servicemodel.Middleware{WithQueue(q)},
		Actions: []servicemodel.Action{
			{
				Name: "Sync",
				Parameters: []servicemodel.Parameter{{
					Source: servicemodel.ParameterSourceBody,
					Name:   "body",
				}},
				Handler: Sync,
			},
			{
				Name:    "Dequeue",
				Handler: Dequeue,
			},
		},
	}).Build()
}

package service

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/lichuan0620/secret-keeper-backend/internal/queueclient"
	"github.com/lichuan0620/secret-keeper-backend/pkg/models"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/middlewares"
	servicemodel "github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
)

func Build(qc queueclient.Interface) (http.Handler, error) {
	logger := log.New().WithName("handlers")
	return (&service.Builder{
		GlobalMiddlewares: []servicemodel.Middleware{
			middlewares.WithLogger(logger),
			middlewares.RequestLog(logger),
		},
	}).AddActionGroup(servicemodel.ActionGroup{
		Mutator: func(action *servicemodel.Action) {
			action.Version = models.Version
		},
		Middlewares: []servicemodel.Middleware{WithQueueClient(qc)},
		Actions: []servicemodel.Action{
			buildStandardActionFromHandler(CreateBox),
			buildStandardActionFromHandler(AddBoxEmoji),
			{
				Name:    "ViewBox",
				Handler: ViewBox,
			},
		},
	}).Build()
}

func buildStandardActionFromHandler(handler interface{}) servicemodel.Action {
	handlerV := reflect.ValueOf(handler)
	if handlerV.Kind() != reflect.Func {
		panic("handler must be a function")
	}
	slice := strings.Split(runtime.FuncForPC(handlerV.Pointer()).Name(), ".")
	return servicemodel.Action{
		Name: slice[len(slice)-1],
		Parameters: []servicemodel.Parameter{{
			Source: servicemodel.ParameterSourceBody,
			Name:   "Body",
		}},
		Handler: handler,
	}
}

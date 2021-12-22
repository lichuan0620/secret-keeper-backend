package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/textproto"
	"reflect"
	"strings"
	"sync"

	"github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/standard"
	"github.com/pkg/errors"
)

type parameter struct {
	source       model.ParameterSource
	name         string
	defaultValue interface{}
	optional     bool
	targetType   reflect.Type
}

type middleware struct {
	self model.Middleware
	next *middleware
}

func parseMiddlewares(middlewares []model.Middleware) *middleware {
	var ret *middleware
	if len(middlewares) > 0 {
		ret = &middleware{
			self: middlewares[0],
		}
		for i, cursor := 1, ret; i < len(middlewares); i++ {
			cursor.next = &middleware{
				self: middlewares[i],
			}
			cursor = cursor.next
		}
	}
	return ret
}

func (m *middleware) execute(ctx context.Context, f func(context.Context)) {
	if m == nil {
		f(ctx)
	} else if m.next == nil {
		m.self(ctx, f)
	} else {
		m.self(ctx, func(ctx context.Context) {
			m.next.execute(ctx, f)
		})
	}
}

type actionHandler struct {
	middlewareLn *middleware
	parameters   []parameter
	handler      reflect.Value
	respPool     sync.Pool
}

// NewActionHandler builds a http.Handler that handles requests for one Action. In most cases, you
// should use Builder instead to build a single handler for all Actions.
func NewActionHandler(action *model.Action, middlewares ...model.Middleware) (http.Handler, error) {
	handler := reflect.ValueOf(action.Handler)
	if handler.Kind() != reflect.Func {
		return nil, errors.New("handler must be a function")
	}
	handlerType := handler.Type()

	if cHandlerIn, cParam := handlerType.NumIn(), len(action.Parameters); cHandlerIn != cParam+1 {
		return nil, errors.Errorf("handler expects %d parameters but %d is given", cHandlerIn-1, cParam)
	}
	parameters := make([]parameter, 0, len(action.Parameters))
	for i := range action.Parameters {
		param, handlerIn := &action.Parameters[i], handlerType.In(i+1)
		if err := validateParameter(param, handlerIn); err != nil {
			return nil, err
		}
		parameters = append(parameters, parameter{
			source:       param.Source,
			name:         param.Name,
			defaultValue: param.Default,
			optional:     param.Optional,
			targetType:   handlerIn,
		})
	}

	if handlerType.NumOut() != 2 {
		return nil, errors.New("handler must produce exactly two outputs")
	}
	if !isJSONType(handlerType.Out(0)) {
		return nil, errors.New("handler produces invalid result type")
	}
	if handlerType.Out(1) != reflect.TypeOf((*standard.Error)(nil)).Elem() {
		return nil, errors.New("handler produces invalid error type")
	}

	ret := actionHandler{
		middlewareLn: parseMiddlewares(middlewares),
		parameters:   parameters,
		handler:      handler,
	}
	ret.respPool.New = func() interface{} {
		return &model.Response{
			Metadata: model.ResponseMetadata{
				Action:  action.Name,
				Version: action.Version,
			},
		}
	}
	return &ret, nil
}

func (exec *actionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	exec.middlewareLn.execute(req.Context(), func(ctx context.Context) {
		// prepare response
		response := exec.respPool.Get().(*model.Response)
		defer func() {
			exec.respPool.Put(response)
		}()

		// parse parameters
		paramValues, err := exec.parseParameters(req)
		if err != nil {
			writeError(w, response, err)
			req.WithContext(context.WithValue(req.Context(), contextKeyError, err))
			return
		}
		paramValues[0] = reflect.ValueOf(ctx)

		// execute handler
		if out := exec.handler.Call(paramValues); out[1].IsNil() {
			writeSuccess(w, response, out[0].Interface())
		} else {
			writeError(w, response, out[1].Interface().(standard.Error))
			req.WithContext(context.WithValue(req.Context(), contextKeyError, err))
		}
	})
}

func (exec *actionHandler) parseParameters(req *http.Request) ([]reflect.Value, standard.Error) {
	paramValues := make([]reflect.Value, 0, len(exec.parameters)+1)
	paramValues = append(paramValues, reflect.ValueOf(req.Context()))
	var err standard.Error
	for i := range exec.parameters {
		var parsed interface{}
		param := &exec.parameters[i]
		switch param.source {
		case model.ParameterSourceQuery:
			values := GetQueryValues(req.Context())
			if parsed, err = parseHeaderOrQuery(param, values[param.name]); err != nil {
				return nil, err
			}
		case model.ParameterSourceHeader:
			values := req.Header[textproto.CanonicalMIMEHeaderKey(param.name)]
			if parsed, err = parseHeaderOrQuery(param, values); err != nil {
				return nil, err
			}
		case model.ParameterSourceBody:
			if !strings.HasPrefix(req.Header.Get(model.HeaderContentType), model.ContentTypeJSON) {
				return nil, standard.UnsupportedContentType()
			}
			value := reflect.New(param.targetType)
			if json.NewDecoder(req.Body).Decode(value.Interface()) != nil {
				return nil, standard.MalformedParameter("body")
			}
			parsed = value.Elem().Interface()
		}
		if parsed == nil {
			if param.defaultValue != nil {
				parsed = param.defaultValue
			} else if param.optional {
				parsed = reflect.Zero(param.targetType).Interface()
			} else {
				return nil, standard.MissingParameter(param.name)
			}
		}
		paramValues = append(paramValues, reflect.ValueOf(parsed))
	}
	return paramValues, nil
}

func isJSONType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Struct:
	case reflect.Map:
	case reflect.Slice:
	case reflect.Ptr:
		if typ.Elem().Kind() != reflect.Struct {
			return false
		}
	default:
		return false
	}
	return true
}

func validateParameter(param *model.Parameter, target reflect.Type) error {
	if len(param.Name) == 0 {
		return errors.New("empty parameter name")
	}
	switch param.Source {
	case model.ParameterSourceQuery, model.ParameterSourceHeader:
		if _, err := ConverterFor(target); err != nil {
			return errors.Wrapf(err, "invalid parameter %s", param.Name)
		}
	case model.ParameterSourceBody:
		if !isJSONType(target) {
			return errors.Errorf("invalid type for body parameter %s", param.Name)
		}
	default:
		return errors.Errorf("invalid source for parameter %s", param.Name)
	}
	if param.Default != nil {
		if defaultType := reflect.ValueOf(param.Default).Type(); !defaultType.AssignableTo(target) {
			return errors.Errorf("default value of type %s is not assignable to %s", defaultType, target)
		}
	}
	return nil
}

func parseHeaderOrQuery(param *parameter, values []string) (interface{}, standard.Error) {
	if len(values) > 0 && len(values[0]) > 0 {
		ret, err := MustConverterFor(param.targetType)(values)
		if err != nil {
			return nil, standard.MalformedParameter(param.name)
		}
		return ret, nil
	}
	return nil, nil
}

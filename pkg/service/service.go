package service

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"

	"github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/standard"
	"github.com/pkg/errors"
)

var (
	contextKeyQueryValue interface{} = new(byte)
	contextKeyAction     interface{} = new(byte)
	contextKeyVersion    interface{} = new(byte)
	contextKeyError      interface{} = new(byte)
)

// GetQueryValues return query values parsed from the given context. When using the http.Handler built
// from this package, the query values are parsed when the request was first received.
func GetQueryValues(ctx context.Context) url.Values {
	values := ctx.Value(contextKeyQueryValue)
	if values == nil {
		return nil
	}
	return values.(url.Values)
}

// GetHandlingInfo returns information about the handling of a request. Some information is only available
// after the handler has returned (so only middlewares can use them).
func GetHandlingInfo(ctx context.Context) HandlingInfo {
	ret := HandlingInfo{
		Action:  ctx.Value(contextKeyAction).(string),
		Version: ctx.Value(contextKeyVersion).(string),
	}
	if err := ctx.Value(contextKeyError); err != nil {
		if stdErr, ok := err.(standard.Error); ok {
			ret.Error = stdErr
		}
	}
	return ret
}

// HandlingInfo contains information about the handling of a request.
type HandlingInfo struct {
	Error           standard.Error
	Action, Version string
}

// record describes a registered Action and is used to build http.Handler later.
type record struct {
	Middlewares []model.Middleware
	Action      *model.Action
}

// Builder is used during the service initialization stage for registering Actions and later
// building a handler for all registered Action.
type Builder struct {
	// GlobalMiddlewares are middlewares that affects all HTTP requests. Note that the middlewares
	// registered with the Actions can only affect a known action; requests that cannot match a
	// registered action can only be covered by the global middlewares.
	GlobalMiddlewares []model.Middleware

	records []record
}

// Build construct a fasthttp router that handles all registered Actions.
func (builder *Builder) Build() (http.Handler, error) {
	handlers := make(map[uint64]http.Handler, len(builder.records))
	for i := range builder.records {
		r := &builder.records[i]
		handler, err := NewActionHandler(r.Action, r.Middlewares...)
		if err != nil {
			return nil, errors.Wrapf(err,
				"invalid definition for action %s version %s",
				r.Action.Name, r.Action.Version,
			)
		}
		index := indexAction(r.Action.Name, r.Action.Version)
		if _, exists := handlers[index]; exists {
			return nil, errors.Errorf("duplicated Action: name %s version %s", r.Action.Name, r.Action.Version)
		}
		handlers[index] = handler
	}
	globalMiddleware := parseMiddlewares(builder.GlobalMiddlewares)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		queryValues := req.URL.Query()
		parse := func(key string) string {
			parsed := strings.Join(queryValues[key], ",")
			if parsed == "" {
				return "<UNKNOWN>"
			}
			return parsed
		}
		action := parse(model.QueryParameterAction)
		version := parse(model.QueryParameterVersion)
		reqCtx := req.Context()
		reqCtx = context.WithValue(reqCtx, contextKeyQueryValue, queryValues)
		reqCtx = context.WithValue(reqCtx, contextKeyAction, action)
		reqCtx = context.WithValue(reqCtx, contextKeyVersion, version)
		globalMiddleware.execute(
			reqCtx,
			func(ctx context.Context) {
				handler, exists := handlers[indexAction(action, version)]
				if exists {
					handler.ServeHTTP(w, req.WithContext(ctx))
				} else {
					writeError(w, &model.Response{
						Metadata: model.ResponseMetadata{
							Action:  action,
							Version: version,
						},
					}, standard.InvalidActionOrVersion(action, version))
				}
			},
		)
	}), nil
}

// AddActionGroup registers a group of Actions to the Builder. Note that the added objects should
// not be modified or reused afterwards.
func (builder *Builder) AddActionGroup(group ...model.ActionGroup) *Builder {
	for i := range group {
		builder.addActionGroup(&group[i])
	}
	return builder
}

func (builder *Builder) addActionGroup(group *model.ActionGroup) {
	for i := range group.Actions {
		action := &group.Actions[i]
		if group.Mutator != nil {
			group.Mutator(action)
		}
		builder.records = append(builder.records, record{
			Middlewares: group.Middlewares,
			Action:      action,
		})
	}
	for i := range group.Subgroups {
		subgroup := &group.Subgroups[i]
		subgroup.Middlewares = append(group.Middlewares, subgroup.Middlewares...)
		subgroupMutator := subgroup.Mutator
		subgroup.Mutator = func(action *model.Action) {
			if group.Mutator != nil {
				group.Mutator(action)
			}
			if subgroupMutator != nil {
				subgroupMutator(action)
			}
		}
		builder.addActionGroup(subgroup)
	}
}

func indexAction(action, version string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(action))
	_, _ = h.Write([]byte{255})
	_, _ = h.Write([]byte(version))
	return h.Sum64()
}

func writeError(w http.ResponseWriter, resp *model.Response, err standard.Error) {
	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(int(err.GetHTTPCode()))
	resp.Result = nil
	if resp.Error == nil {
		resp.Error = new(model.Error)
	}
	resp.Error.Code, resp.Error.Message, resp.Error.Data = err.GetCode(), err.GetMessage(), err.GetData()
	_ = json.NewEncoder(w).Encode(resp)
}

func writeSuccess(w http.ResponseWriter, resp *model.Response, result interface{}) {
	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	resp.Result = result
	resp.Error = nil
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(resp)
}

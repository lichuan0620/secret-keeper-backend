package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/standard"
)

const region, service = "LF", "VMP"
const fakeURL = "http://localhost/api"

func TestBuilder(t *testing.T) {
	type RequestTC struct {
		Request      *http.Request
		ExpectStatus int
		ExpectBody   model.Response
	}
	tests := []*struct {
		Actions        model.ActionGroup
		ActionsInvalid bool
		RequestTCs     []*RequestTC
	}{
		{
			Actions: model.ActionGroup{
				Actions: []model.Action{
					{
						Name:    "DoNothing",
						Version: "20211128",
						Handler: func(context.Context) (*struct{}, standard.Error) { return &struct{}{}, nil },
					},
					{
						Name:    "DoNothing",
						Version: "20211128",
						Handler: func(context.Context) (*struct{}, standard.Error) { return &struct{}{}, nil },
					},
				},
			},
			ActionsInvalid: true,
		},
		{
			Actions: model.ActionGroup{
				Actions: []model.Action{
					{
						Name:    "DoNothing",
						Version: "20211128",
						Handler: func(context.Context) (*struct{}, standard.Error) { return &struct{}{}, nil },
					},
				},
				Subgroups: []model.ActionGroup{
					{
						Actions: []model.Action{
							{
								Name:    "DoNothing",
								Version: "20211128",
								Handler: func(context.Context) (*struct{}, standard.Error) { return &struct{}{}, nil },
							},
						},
					},
				},
			},
			ActionsInvalid: true,
		},
		{
			Actions: model.ActionGroup{
				Subgroups: []model.ActionGroup{
					{
						Mutator: func(action *model.Action) {
							action.Name = "EchoQuery"
							action.Handler = func(_ context.Context, data string) (map[string]interface{}, standard.Error) {
								return map[string]interface{}{"Data": data}, nil
							}
						},
						Actions: []model.Action{
							{
								Version: "20211128",
								Parameters: []model.Parameter{{
									Source: model.ParameterSourceQuery,
									Name:   "Data",
								}},
							},
							{
								Version: "20211129",
								Parameters: []model.Parameter{{
									Source:  model.ParameterSourceQuery,
									Name:    "Data",
									Default: "Hello World",
								}},
							},
						},
					},
					{
						Actions: []model.Action{
							{
								Name:    "EchoBody",
								Version: "20211128",
								Parameters: []model.Parameter{{
									Source: model.ParameterSourceBody,
									Name:   "Data",
								}},
								Handler: func(_ context.Context, data map[string]interface{}) (map[string]interface{}, standard.Error) {
									return data, nil
								},
							},
							{
								Name:    "EchoHeader",
								Version: "20211128",
								Parameters: []model.Parameter{{
									Source: model.ParameterSourceHeader,
									Name:   "Data",
								}},
								Handler: func(_ context.Context, data string) (map[string]interface{}, standard.Error) {
									return map[string]interface{}{"Data": data}, nil
								},
							},
						},
					},
				},
			},
			RequestTCs: []*RequestTC{
				{
					Request: func() *http.Request {
						ret, _ := http.NewRequest(
							http.MethodPost,
							fakeURL+"?Action=Unknown&Version=20211128",
							bytes.NewBuffer(nil),
						)
						ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
						return ret
					}(),
					ExpectStatus: http.StatusNotFound,
					ExpectBody: model.Response{
						Metadata: model.ResponseMetadata{
							Action:  "Unknown",
							Version: "20211128",
						},
						Error: func() *model.Error {
							err := standard.InvalidActionOrVersion("Unknown", "20211128")
							return &model.Error{
								Code:    err.GetCode(),
								Message: err.GetMessage(),
								Data:    err.GetData(),
							}
						}(),
					},
				},
				{
					Request: func() *http.Request {
						ret, _ := http.NewRequest(
							http.MethodPost,
							fakeURL+"?Action=EchoHeader&Version=20211128",
							bytes.NewBuffer(nil),
						)
						ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
						ret.Header.Set("Data", "test")
						return ret
					}(),
					ExpectStatus: http.StatusOK,
					ExpectBody: model.Response{
						Metadata: model.ResponseMetadata{
							Action:  "EchoHeader",
							Version: "20211128",
						},
						Result: map[string]interface{}{"Data": "test"},
					},
				},
				{
					Request: func() *http.Request {
						ret, _ := http.NewRequest(
							http.MethodPost,
							fakeURL+"?Action=EchoQuery&Version=20211128&Data=test",
							bytes.NewBuffer(nil),
						)
						ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
						return ret
					}(),
					ExpectStatus: http.StatusOK,
					ExpectBody: model.Response{
						Metadata: model.ResponseMetadata{
							Action:  "EchoQuery",
							Version: "20211128",
						},
						Result: map[string]interface{}{"Data": "test"},
					},
				},
				{
					Request: func() *http.Request {
						ret, _ := http.NewRequest(
							http.MethodPost,
							fakeURL+"?Action=EchoQuery&Version=20211129",
							bytes.NewBuffer(nil),
						)
						ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
						return ret
					}(),
					ExpectStatus: http.StatusOK,
					ExpectBody: model.Response{
						Metadata: model.ResponseMetadata{
							Action:  "EchoQuery",
							Version: "20211129",
						},
						Result: map[string]interface{}{"Data": "Hello World"},
					},
				},
				{
					Request: func() *http.Request {
						ret, _ := http.NewRequest(
							http.MethodPost,
							fakeURL+"?Action=EchoQuery&Version=20211128",
							bytes.NewBuffer(nil),
						)
						ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
						return ret
					}(),
					ExpectStatus: http.StatusBadRequest,
					ExpectBody: model.Response{
						Metadata: model.ResponseMetadata{
							Action:  "EchoQuery",
							Version: "20211128",
						},
						Error: func() *model.Error {
							err := standard.MissingParameter("Data")
							return &model.Error{
								Code:    err.GetCode(),
								Message: err.GetMessage(),
								Data:    err.GetData(),
							}
						}(),
					},
				},
				{
					Request: func() *http.Request {
						ret, _ := http.NewRequest(
							http.MethodPost,
							fakeURL+"?Action=EchoBody&Version=20211128",
							bytes.NewBuffer([]byte(`{"Data":"test"}`)),
						)
						ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
						return ret
					}(),
					ExpectStatus: http.StatusOK,
					ExpectBody: model.Response{
						Metadata: model.ResponseMetadata{
							Action:  "EchoBody",
							Version: "20211128",
						},
						Result: map[string]interface{}{"Data": "test"},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run("build", func(tt *testing.T) {
			handler, err := (&Builder{}).AddActionGroup(tc.Actions).Build()
			if tc.ActionsInvalid == (err == nil) {
				tt.Fatalf("expecting invalid action: %v; add action got error: %v", tc.ActionsInvalid, err)
			}
			if tc.ActionsInvalid {
				return
			}
			for _, reqTC := range tc.RequestTCs {
				tt.Run("request", func(ttt *testing.T) {
					rr := httptest.NewRecorder()
					handler.ServeHTTP(rr, reqTC.Request)
					if rr.Code != reqTC.ExpectStatus {
						ttt.Fatalf("expecting response status %d; got %d", reqTC.ExpectStatus, rr.Code)
					}
					var buf model.Response
					if err = json.Unmarshal(rr.Body.Bytes(), &buf); err != nil {
						ttt.Fatalf("unmarshal response, body: %s, err: %s", rr.Body.String(), err)
					}
					if !reflect.DeepEqual(reqTC.ExpectBody, buf) {
						ttt.Fatalf("expecting response body: %+v got: %+v", reqTC.ExpectBody, buf)
					}
				})
			}
		})
	}
}

func TestGlobalMiddleware(t *testing.T) {
	builder := Builder{
		GlobalMiddlewares: []model.Middleware{
			func(ctx context.Context, f func(context.Context)) {
				f(ctx)
				info := GetHandlingInfo(ctx)
				if info.Error == nil {
					t.Fatal("didn't get response error")
				}
				if info.Error.GetHTTPCode() == http.StatusOK {
					t.Fatal("expecting error status code, got 200")
				}
			},
		},
	}
	const action, version = "EchoError", "20211206"
	builder.AddActionGroup(model.ActionGroup{
		Actions: []model.Action{{
			Name:    action,
			Version: version,
			Handler: func(ctx context.Context) (*struct{}, standard.Error) {
				return nil, standard.InternalServiceError()
			},
		}},
	})
	h, err := builder.Build()
	if err != nil {
		t.Fatalf("build error: %v", err)
	}
	h.ServeHTTP(httptest.NewRecorder(), func() *http.Request {
		ret, _ := http.NewRequest(
			http.MethodPost,
			fakeURL+"?Action=Unknown&Version=20211206",
			bytes.NewBuffer(nil),
		)
		ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
		return ret
	}())
	h.ServeHTTP(httptest.NewRecorder(), func() *http.Request {
		ret, _ := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s?Action=%s&Version=%s", fakeURL, action, version),
			bytes.NewBuffer(nil),
		)
		ret.Header.Set(model.HeaderContentType, model.ContentTypeJSON)
		return ret
	}())
}

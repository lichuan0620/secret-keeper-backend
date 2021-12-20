package model

import (
	"context"
)

// Action describes a RPC action: its name and version, parameters, and a handler function.
type Action struct {
	// Name is the name of the RPC action such as CreateInstance.
	Name string
	// Version is the version of the RPC action such as 20211123
	Version string
	// Parameters describes the parameters of Handler.
	Parameters []Parameter
	// Handler is a function that handles Action requests. It should have at least one parameter, the
	// first one being a context.Context, and the rest are described by Parameters; it should have
	// exactly two return values, the first one being the response body struct (pointer or value), the
	// second a standard.Error.
	Handler interface{}
}

// ActionGroup describes a set of Actions that share some common properties.
type ActionGroup struct {
	// Mutator will be called on Actions from this ActionGroup and all of the Subgroups during building
	// process. The outer ActionMutator is called first (so the inner one takes priority).
	Mutator func(*Action)
	// Middlewares are added to the Actions from this ActionGroup and all of the Subgroups. The order of
	// their position in the slice is important: the leftmost Middleware is called first (so it can be
	// viewed as the outer most Middleware).
	Middlewares []Middleware
	// Subgroups are nested ActionGroups. An ActionGroup shares its ActionMutator and the Middlewares
	// with all of its children.
	Subgroups []ActionGroup
	// Actions is a list of Action belonging to this group. Note that no two Actions under the a ActionGroup
	// and all of its Subgroups should have the same name-version set.
	Actions []Action
}

// Middleware defines optional wrappers for Action handlers. It provide a way to add pre-run and
// post-run code to Action handlers.
type Middleware func(context.Context, func(context.Context))

// Parameter describes a function parameter.
type Parameter struct {
	// Source the place from which this parameter is parsed.
	Source ParameterSource
	// Name is the name of this parameter. ex. a query name, a header key, etc.
	Name string
	// Default value is used when a request does not provide a value
	// for the parameter.
	Default interface{}
	// Optional, if set to true, makes the a header or query Parameter optional. If a request is
	// missing a parameter that is required and has no Default, the request would fail with a
	// MissingParameter error.
	Optional bool
}

// ParameterSource indicates the place from which the value of an Action parameter is parsed.
type ParameterSource string

const (
	// ParameterSourceQuery means value is from URL query string.
	ParameterSourceQuery ParameterSource = "Query"

	// ParameterSourceHeader means value is from request header.
	ParameterSourceHeader ParameterSource = "Header"

	// ParameterSourceBody means value is from request body. Only JSON format is supported per
	// our API convention. If a Body parameter is added, a request with content type other than
	// application/json would cause a 415 error.
	ParameterSourceBody ParameterSource = "Body"
)

// Error is the standard RPC response format.
type Error struct {
	Code    string            `json:"Code"`
	Message string            `json:"Message"`
	Data    map[string]string `json:"Data"`
}

// Response is the standard RPC response format.
type Response struct {
	Metadata ResponseMetadata `json:"ResponseMetadata"`
	Result   interface{}      `json:"Result,omitempty"`
	Error    *Error           `json:"Error,omitempty"`
}

// ResponseMetadata is the standard RPC response metadata format.
type ResponseMetadata struct {
	Action  string `json:"Action"`
	Version string `json:"Version"`
}

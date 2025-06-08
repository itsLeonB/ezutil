package types

type ParamType string

const (
	ParamTypeQuery ParamType = "query"
	ParamTypePath  ParamType = "path"
)

type BindType string

const (
	BindTypeJSON  BindType = "json"
	BindTypeForm  BindType = "form"
	BindTypeQuery BindType = "query"
)

package surrealdb

import (
	"github.com/tai-kun/surrealdb.go/pkg/codec"
	"github.com/tai-kun/surrealdb.go/pkg/engines"
)

type (
	Engine  = func(fmt codec.Formatter) engines.Engine
	Engines = map[string]Engine
)

var (
	HTTPEngine Engine = func(fmt codec.Formatter) engines.Engine {
		return engines.NewHTTPEngine(fmt)
	}
)

package gqlgen

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fwojciec/asr"
)

func NewQueryHandler(data []*asr.Service) http.Handler {
	c := NewCache(data)
	es := NewExecutableSchema(Config{Resolvers: &Resolver{Cache: c}})
	srv := handler.New(es)
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New(100)})
	return srv
}

func NewPlaygroundHandler(title, endpoint string) http.Handler {
	return playground.Handler(title, endpoint)
}

package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/cache"
	"github.com/fwojciec/asr/gqlgen"
)

//go:embed data.json
var data []byte

var muxAdapter *httpadapter.HandlerAdapter

func init() {
	log.Println("initializing")
	var d []*asr.Service
	_ = json.Unmarshal(data, &d)

	c := cache.NewCache(d)
	es := gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &gqlgen.Resolver{Cache: c}})
	srv := handler.New(es)
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	r := http.NewServeMux()
	r.Handle("/query", srv)
	r.Handle("/playground", playground.Handler("GraphQL", fmt.Sprintf("/%s/query", os.Getenv("STAGE"))))
	muxAdapter = httpadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rsp, err := muxAdapter.Proxy(req)
	if err != nil {
		log.Println(err)
	}
	return rsp, err
}

func main() {
	lambda.Start(Handler)
}

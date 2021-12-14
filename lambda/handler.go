package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/fwojciec/asr"
	"github.com/fwojciec/asr/gqlgen"
)

//go:embed data.json
var data []byte

var muxAdapter *httpadapter.HandlerAdapter

func init() {
	var d []*asr.Service
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatal(err)
	}
	r := http.NewServeMux()
	r.Handle("/query", gqlgen.NewQueryHandler(d))
	r.Handle("/playground", gqlgen.NewPlaygroundHandler("GraphQL", path.Join("/", os.Getenv("STAGE"), "query")))
	muxAdapter = httpadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxAdapter.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

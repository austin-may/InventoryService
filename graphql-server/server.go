package main

import (
	"log"
	"my-go-apps/InventoryService/graph"
	"my-go-apps/InventoryService/graph/generated"
	"my-go-apps/InventoryService/cors"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", cors.MiddlewareHandler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playgroundaustin", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

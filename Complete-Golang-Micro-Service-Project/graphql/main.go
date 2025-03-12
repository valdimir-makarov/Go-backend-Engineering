package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"

	"github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/graphql/graph"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var configvar AppConfig
	err := envconfig.Process("", &configvar)
	if err != nil {
		log.Fatal(err)
	}

	s, err := graph.NewGraphQLServer(configvar.AccountURL, configvar.CatalogURL, configvar.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	// ✅ Create a GraphQL handler with explicit transport configuration
	// srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: s}))
	srv := handler.New(s.ToExecutableSchema())

	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Println("🚀 Server running on http://localhost:8080/playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package naos

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/friendsofgo/graphiql"
	"github.com/julienschmidt/httprouter"
	"github.com/Dophin2009/nao/internal/graphql"
	"github.com/Dophin2009/nao/pkg/web"
)

// NewGraphQLHandler returns a POST endpoint handler for the GraphQL API.
func NewGraphQLHandler(path []string, ds *graphql.DataService) web.Handler {
	cfg := graphql.Config{
		Resolvers: &graphql.Resolver{},
	}
	gqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(cfg))

	ctx := context.WithValue(context.Background(), graphql.DataServiceKey, ds)
	return web.Handler{
		Method: http.MethodPost,
		Path:   path,
		Func: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			r = r.WithContext(ctx)
			gqlHandler.ServeHTTP(w, r)
		},
	}
}

// NewGraphiQLHandler returns a new GET endpoint handler for rendering a
// GraphiQL page for the given GraphQL API.
func NewGraphiQLHandler(path []string, graphqlPath string) (web.Handler, error) {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler(graphqlPath)
	if err != nil {
		return web.Handler{}, nil
	}
	return web.Handler{
		Method: http.MethodGet,
		Path:   path,
		Func: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			graphiqlHandler.ServeHTTP(w, r)
		},
	}, nil
}
